/*
Package barbican contains an implementation of the github.com/getsops/sops/v3.MasterKey
interface that encrypts and decrypts the data key using OpenStack Barbican with the
OpenStack SDK for Go.
*/
package barbican // import "github.com/getsops/sops/v3/barbican"

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// secretRefRegex matches a Barbican secret reference in various formats:
	// UUID: "550e8400-e29b-41d4-a716-446655440000"
	// URI: "https://barbican.example.com:9311/v1/secrets/550e8400-e29b-41d4-a716-446655440000"
	// Regional: "region:sjc3:550e8400-e29b-41d4-a716-446655440000"
	secretRefRegex = `^(?:(?:https?://[^/]+/v1/secrets/)|(?:region:[^:]+:))?([0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`
	
	// barbicanTTL is the duration after which a MasterKey requires rotation.
	barbicanTTL = time.Hour * 24 * 30 * 6
	
	// KeyTypeIdentifier is the string used to identify a Barbican MasterKey.
	KeyTypeIdentifier = "barbican"
)

var (
	// log is the global logger for any Barbican MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("BARBICAN")
}

// MasterKey is an OpenStack Barbican secret used to encrypt and decrypt SOPS' data key.
type MasterKey struct {
	// SecretRef is the Barbican secret reference (UUID, URI, or regional format)
	SecretRef string
	
	// Region specifies the OpenStack region for multi-region support
	Region string
	
	// EncryptedKey stores the encrypted data key (Barbican secret UUID)
	EncryptedKey string
	
	// CreationDate tracks when this master key was created
	CreationDate time.Time
	
	// AuthConfig contains OpenStack authentication configuration
	AuthConfig *AuthConfig
	
	// Internal fields for client management
	client              *BarbicanClient
	credentialsProvider *CredentialsProvider
	httpClient          *http.Client
	baseEndpoint        string
	authManager         *AuthManager
}

// AuthConfig holds OpenStack authentication parameters
type AuthConfig struct {
	AuthURL                     string
	Region                      string
	ProjectID                   string
	ProjectName                 string
	DomainID                    string
	DomainName                  string
	Username                    string
	Password                    string
	ApplicationCredentialID     string
	ApplicationCredentialSecret string
	Token                       string
	Insecure                    bool
	CACert                      string
}



// SecretMetadata contains metadata for Barbican secrets
type SecretMetadata struct {
	Name        string            `json:"name"`
	Algorithm   string            `json:"algorithm"`
	BitLength   int               `json:"bit_length"`
	Mode        string            `json:"mode"`
	SecretType  string            `json:"secret_type"`
	ContentType string            `json:"payload_content_type"`
	Expiration  *time.Time        `json:"expiration,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}



// CredentialsProvider is a wrapper around authentication credentials used for
// authentication towards OpenStack Barbican.
type CredentialsProvider struct {
	config *AuthConfig
}

// NewCredentialsProvider returns a CredentialsProvider object with the provided
// AuthConfig.
func NewCredentialsProvider(config *AuthConfig) *CredentialsProvider {
	return &CredentialsProvider{
		config: config,
	}
}

// ApplyToMasterKey configures the credentials on the provided key.
func (c *CredentialsProvider) ApplyToMasterKey(key *MasterKey) {
	key.credentialsProvider = c
	key.AuthConfig = c.config
	
	// Initialize authentication manager if config is provided
	if c.config != nil {
		authManager, err := NewAuthManager(c.config)
		if err != nil {
			log.WithError(err).Warn("Failed to initialize authentication manager")
		} else {
			key.authManager = authManager
		}
	}
}

// HTTPClient is a wrapper around http.Client used for configuring the
// Barbican client.
type HTTPClient struct {
	hc *http.Client
}

// NewHTTPClient creates a new HTTPClient with the provided http.Client.
func NewHTTPClient(hc *http.Client) *HTTPClient {
	return &HTTPClient{hc: hc}
}

// ApplyToMasterKey configures the HTTP client on the provided key.
func (h *HTTPClient) ApplyToMasterKey(key *MasterKey) {
	key.httpClient = h.hc
}

// NewMasterKey creates a new MasterKey from a secret reference, setting
// the creation date to the current date.
func NewMasterKey(secretRef string) *MasterKey {
	return &MasterKey{
		SecretRef:    secretRef,
		CreationDate: time.Now().UTC(),
	}
}

// NewMasterKeyWithRegion creates a new MasterKey from a secret reference and region,
// setting the creation date to the current date.
func NewMasterKeyWithRegion(secretRef string, region string) *MasterKey {
	key := NewMasterKey(secretRef)
	key.Region = region
	return key
}

// NewMasterKeyFromSecretRef takes a Barbican secret reference string and returns a new
// MasterKey for that reference. The reference can be in UUID, URI, or regional format.
func NewMasterKeyFromSecretRef(secretRef string) (*MasterKey, error) {
	secretRef = strings.TrimSpace(secretRef)
	
	// Validate secret reference using security validator
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	if err := securityValidator.ValidateSecretRef(secretRef); err != nil {
		return nil, err
	}
	
	key := &MasterKey{
		SecretRef:    secretRef,
		CreationDate: time.Now().UTC(),
	}
	
	// Extract region from regional format
	if strings.HasPrefix(secretRef, "region:") {
		parts := strings.Split(secretRef, ":")
		if len(parts) >= 3 {
			key.Region = parts[1]
		}
	}
	
	return key, nil
}

// MasterKeysFromSecretRefString takes a comma separated list of Barbican secret
// references, and returns a slice of new MasterKeys for those references.
func MasterKeysFromSecretRefString(secretRefs string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if secretRefs == "" {
		return keys, nil
	}
	
	for _, s := range strings.Split(secretRefs, ",") {
		key, err := NewMasterKeyFromSecretRef(s)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// isValidSecretRef validates a Barbican secret reference format
func isValidSecretRef(secretRef string) bool {
	re := regexp.MustCompile(secretRefRegex)
	return re.MatchString(secretRef)
}

// extractUUIDFromSecretRef extracts the UUID from various secret reference formats
func extractUUIDFromSecretRef(secretRef string) (string, error) {
	re := regexp.MustCompile(secretRefRegex)
	matches := re.FindStringSubmatch(secretRef)
	if len(matches) < 2 {
		return "", fmt.Errorf("could not extract UUID from secret reference: %s", secretRef)
	}
	return matches[1], nil
}

// Encrypt takes a SOPS data key, encrypts it with Barbican and stores the result
// in the EncryptedKey field.
//
// Consider using EncryptContext instead.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a SOPS data key, encrypts it with Barbican and stores the result
// in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	// Get or create Barbican client
	client, err := key.getBarbicanClient(ctx)
	if err != nil {
		return WrapError(err, ErrorTypeConfig, "Failed to create Barbican client")
	}

	// Prepare secret metadata
	metadata := SecretMetadata{
		Name:        "SOPS Data Key",
		SecretType:  "opaque",
		ContentType: "application/octet-stream",
		Metadata: map[string]string{
			"created_by": "sops",
			"purpose":    "data_key_encryption",
		},
	}

	// Store the data key as a secret in Barbican
	secretRef, err := client.StoreSecret(ctx, dataKey, metadata)
	if err != nil {
		return WrapError(err, ErrorTypeAPI, "Failed to encrypt data key with Barbican")
	}

	// Store the secret reference as the encrypted key
	key.EncryptedKey = secretRef
	
	// Sanitize secret reference for logging
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	sanitizedRef := securityValidator.sanitizeValue("secret_ref", secretRef)
	log.WithField("secret_ref", sanitizedRef).Debug("Data key encrypted and stored in Barbican")
	
	return nil
}

// EncryptMultiRegion encrypts a data key across multiple regions in parallel
func EncryptMultiRegion(ctx context.Context, dataKey []byte, keys []*MasterKey) error {
	if len(keys) == 0 {
		return fmt.Errorf("no master keys provided for multi-region encryption")
	}

	// Use a channel to collect results from parallel operations
	type encryptResult struct {
		key   *MasterKey
		error error
	}
	
	resultChan := make(chan encryptResult, len(keys))
	
	// Start encryption operations in parallel
	for _, key := range keys {
		go func(k *MasterKey) {
			err := k.EncryptContext(ctx, dataKey)
			resultChan <- encryptResult{key: k, error: err}
		}(key)
	}
	
	// Collect results
	var errors []error
	successCount := 0
	
	for i := 0; i < len(keys); i++ {
		result := <-resultChan
		if result.error != nil {
			region := result.key.getEffectiveRegion()
			log.WithError(result.error).WithField("region", region).Warn("Failed to encrypt in region")
			errors = append(errors, fmt.Errorf("region %s: %w", region, result.error))
		} else {
			successCount++
			region := result.key.getEffectiveRegion()
			log.WithField("region", region).Debug("Successfully encrypted in region")
		}
	}
	
	// Require at least one successful encryption
	if successCount == 0 {
		return fmt.Errorf("failed to encrypt in any region: %v", errors)
	}
	
	// Log partial failures but don't fail the operation
	if len(errors) > 0 {
		log.WithField("failed_regions", len(errors)).WithField("successful_regions", successCount).Warn("Some regions failed during encryption")
	}
	
	return nil
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been
// encrypted yet.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Decrypt decrypts the EncryptedKey with Barbican and returns the result.
//
// Consider using DecryptContext instead.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey with Barbican and returns the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	if key.EncryptedKey == "" {
		return nil, NewValidationError("No encrypted key to decrypt")
	}

	// Get or create Barbican client
	client, err := key.getBarbicanClient(ctx)
	if err != nil {
		return nil, WrapError(err, ErrorTypeConfig, "Failed to create Barbican client")
	}

	// Retrieve the data key from Barbican
	dataKey, err := client.GetSecretPayload(ctx, key.EncryptedKey)
	if err != nil {
		return nil, WrapError(err, ErrorTypeAPI, "Failed to decrypt data key from Barbican")
	}

	// Sanitize secret reference for logging
	securityValidator := NewSecurityValidator(DefaultSecurityConfig())
	sanitizedRef := securityValidator.sanitizeValue("secret_ref", key.EncryptedKey)
	log.WithField("secret_ref", sanitizedRef).Debug("Data key decrypted from Barbican")
	
	return dataKey, nil
}

// DecryptMultiRegion attempts to decrypt using multiple master keys with failover logic
func DecryptMultiRegion(ctx context.Context, keys []*MasterKey) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no master keys provided for multi-region decryption")
	}

	// Try keys in order, implementing failover logic
	var lastError error
	
	for i, key := range keys {
		if key.EncryptedKey == "" {
			continue // Skip keys without encrypted data
		}
		
		region := key.getEffectiveRegion()
		log.WithField("region", region).WithField("attempt", i+1).Debug("Attempting decryption")
		
		dataKey, err := key.DecryptContext(ctx)
		if err != nil {
			log.WithError(err).WithField("region", region).Warn("Decryption failed in region")
			lastError = err
			continue
		}
		
		log.WithField("region", region).Debug("Successfully decrypted from region")
		return dataKey, nil
	}
	
	// If we get here, all regions failed
	return nil, fmt.Errorf("failed to decrypt from any region, last error: %w", lastError)
}

// DecryptMultiRegionParallel attempts to decrypt using multiple master keys in parallel
// Returns the first successful result
func DecryptMultiRegionParallel(ctx context.Context, keys []*MasterKey) ([]byte, error) {
	if len(keys) == 0 {
		return nil, fmt.Errorf("no master keys provided for multi-region decryption")
	}

	// Filter keys that have encrypted data
	var validKeys []*MasterKey
	for _, key := range keys {
		if key.EncryptedKey != "" {
			validKeys = append(validKeys, key)
		}
	}
	
	if len(validKeys) == 0 {
		return nil, fmt.Errorf("no master keys have encrypted data")
	}

	// Use a channel to collect results from parallel operations
	type decryptResult struct {
		dataKey []byte
		region  string
		error   error
	}
	
	resultChan := make(chan decryptResult, len(validKeys))
	
	// Create a context that can be cancelled when we get the first success
	decryptCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	
	// Start decryption operations in parallel
	for _, key := range validKeys {
		go func(k *MasterKey) {
			region := k.getEffectiveRegion()
			dataKey, err := k.DecryptContext(decryptCtx)
			resultChan <- decryptResult{
				dataKey: dataKey,
				region:  region,
				error:   err,
			}
		}(key)
	}
	
	// Wait for the first successful result or all failures
	var errors []error
	
	for i := 0; i < len(validKeys); i++ {
		result := <-resultChan
		
		if result.error != nil {
			log.WithError(result.error).WithField("region", result.region).Debug("Parallel decryption failed in region")
			errors = append(errors, fmt.Errorf("region %s: %w", result.region, result.error))
		} else {
			log.WithField("region", result.region).Debug("Successfully decrypted from region (parallel)")
			cancel() // Cancel remaining operations
			return result.dataKey, nil
		}
	}
	
	// If we get here, all regions failed
	return nil, fmt.Errorf("failed to decrypt from any region in parallel: %v", errors)
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > barbicanTTL
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	if key.Region != "" && !strings.HasPrefix(key.SecretRef, "region:") {
		return fmt.Sprintf("region:%s:%s", key.Region, key.SecretRef)
	}
	return key.SecretRef
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key *MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["secret_ref"] = key.SecretRef
	if key.Region != "" {
		out["region"] = key.Region
	}
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// getAuthToken retrieves an authentication token using the configured authentication manager
func (key *MasterKey) getAuthToken(ctx context.Context) (string, string, error) {
	if key.authManager == nil {
		// Try to initialize auth manager if we have config
		if key.AuthConfig != nil {
			authManager, err := NewAuthManager(key.AuthConfig)
			if err != nil {
				return "", "", WrapError(err, ErrorTypeAuthentication, "Failed to initialize authentication manager")
			}
			key.authManager = authManager
		} else {
			// Try to load from environment
			config := LoadConfigFromEnvironment()
			if err := ValidateConfig(config); err != nil {
				return "", "", WrapError(err, ErrorTypeConfig, "No valid authentication configuration found")
			}
			
			authManager, err := NewAuthManager(config)
			if err != nil {
				return "", "", WrapError(err, ErrorTypeAuthentication, "Failed to initialize authentication manager")
			}
			key.authManager = authManager
			key.AuthConfig = config
		}
	}
	
	return key.authManager.GetToken(ctx)
}

// getBarbicanClient gets or creates a Barbican client for this master key
func (key *MasterKey) getBarbicanClient(ctx context.Context) (*BarbicanClient, error) {
	if key.client != nil {
		return key.client, nil
	}

	// Ensure we have an auth manager
	if key.authManager == nil {
		_, _, err := key.getAuthToken(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Get Barbican endpoint for the specific region
	var endpoint string
	var err error
	
	if key.baseEndpoint != "" {
		endpoint = key.baseEndpoint
	} else {
		// Determine region from secret reference or key region
		region := key.getEffectiveRegion()
		
		// Try to discover endpoint from auth manager for the specific region
		endpoint, err = GetBarbicanEndpointForRegion(key.authManager, region)
		if err != nil {
			return nil, NewConfigError("Failed to get Barbican endpoint").
				WithRegion(region).
				WithCause(err).
				WithSuggestions(
					"Verify the region name is correct",
					"Check network connectivity to OpenStack services",
					"Ensure the Barbican service is available in the specified region",
				)
		}
		key.baseEndpoint = endpoint
	}

	// Create client config
	config := DefaultClientConfig()
	if key.AuthConfig != nil {
		config.Insecure = key.AuthConfig.Insecure
		config.CACert = key.AuthConfig.CACert
	}

	// Create the client
	client, err := NewBarbicanClient(endpoint, key.authManager, config)
	if err != nil {
		return nil, WrapError(err, ErrorTypeConfig, "Failed to create Barbican client")
	}

	key.client = client
	return client, nil
}

// getEffectiveRegion returns the region to use for this master key
func (key *MasterKey) getEffectiveRegion() string {
	// First check if the secret reference contains a region
	if strings.HasPrefix(key.SecretRef, "region:") {
		parts := strings.Split(key.SecretRef, ":")
		if len(parts) >= 3 {
			return parts[1]
		}
	}
	
	// Fall back to the key's region field
	if key.Region != "" {
		return key.Region
	}
	
	// Fall back to auth config region
	if key.AuthConfig != nil && key.AuthConfig.Region != "" {
		return key.AuthConfig.Region
	}
	
	// Default region
	return "RegionOne"
}

// GroupKeysByRegion groups master keys by their effective region
func GroupKeysByRegion(keys []*MasterKey) map[string][]*MasterKey {
	regionGroups := make(map[string][]*MasterKey)
	
	for _, key := range keys {
		region := key.getEffectiveRegion()
		regionGroups[region] = append(regionGroups[region], key)
	}
	
	return regionGroups
}

// GetRegionsFromKeys extracts unique regions from a list of master keys
func GetRegionsFromKeys(keys []*MasterKey) []string {
	regionSet := make(map[string]bool)
	
	for _, key := range keys {
		region := key.getEffectiveRegion()
		regionSet[region] = true
	}
	
	var regions []string
	for region := range regionSet {
		regions = append(regions, region)
	}
	
	return regions
}

// ValidateMultiRegionKeys validates that all keys in different regions are properly configured
func ValidateMultiRegionKeys(keys []*MasterKey) error {
	if len(keys) == 0 {
		return fmt.Errorf("no master keys provided")
	}

	regionGroups := GroupKeysByRegion(keys)
	
	// Validate each region group
	for region, regionKeys := range regionGroups {
		if len(regionKeys) == 0 {
			continue
		}
		
		// Check that all keys in the same region have compatible auth configs
		baseAuthConfig := regionKeys[0].AuthConfig
		for i, key := range regionKeys {
			if i == 0 {
				continue
			}
			
			// Basic validation - in a real implementation, you might want more sophisticated checks
			if key.AuthConfig != nil && baseAuthConfig != nil {
				if key.AuthConfig.AuthURL != baseAuthConfig.AuthURL {
					log.WithField("region", region).Warn("Keys in same region have different auth URLs")
				}
			}
		}
		
		log.WithField("region", region).WithField("key_count", len(regionKeys)).Debug("Validated region key group")
	}
	
	return nil
}