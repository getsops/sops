package ocikms

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/getsops/sops/v3/logging"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
	"github.com/sirupsen/logrus"
)

var (
	// log is the global logger for any OCI KMS MasterKey.
	log *logrus.Logger
	// ocikmsTTL is the duration after which a MasterKey requires rotation.
	ocikmsTTL = time.Hour * 24 * 30 * 6
)

const (
	// ocidParts is the number of parts in an OCID, separated by ".", eg: "ocid1.key.oc1.uk-london-1.aaaalgz5aacmg.aaaailjtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq"
	ocidParts = 6
)

func init() {
	log = logging.NewLogger(LoggerName)
}

// MasterKey is an Oracle Cloud KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	// Ocid is the Oracle Cloud Identifier for the KMS key
	Ocid string
	// EncryptedKey stores the SOPS data key in its encrypted form
	EncryptedKey string
	// CreationDate is when this MasterKey was created
	CreationDate time.Time

	// configProvider is used to configure the OCI client with credentials.
	// It can be injected by a (local) keyservice.KeyServiceServer using
	// ConfigurationProvider.ApplyToMasterKey. If nil, a fresh config
	// provider is created on each operation which tries multiple auth methods.
	configProvider common.ConfigurationProvider
	// httpClient is used to override the default HTTP client used by the OCI client.
	// Mostly useful for testing purposes.
	httpClient common.HTTPRequestDispatcher
}

func NewMasterKeyFromOCID(ocid string) *MasterKey {
	return &MasterKey{
		Ocid:         strings.Replace(ocid, " ", "", -1),
		CreationDate: time.Now().UTC(),
	}
}

func MasterKeysFromOCIDString(ocids string) []*MasterKey {
	var keys []*MasterKey
	if ocids == "" {
		return keys
	}
	for _, s := range strings.Split(ocids, ",") {
		keys = append(keys, NewMasterKeyFromOCID(s))
	}
	return keys
}

// createCryptoClient creates a new OCI KMS client. It uses the injected configProvider
// if available, otherwise creates a new one on each call. If httpClient is set, it uses
// that for HTTP requests (useful for testing).
func (key *MasterKey) createCryptoClient() (client keymanagement.KmsCryptoClient, err error) {
	region, vaultExt, err := extractRefs(key)
	if err != nil {
		log.WithField("ocid", key.Ocid).Errorf("Failed to extract region and vault from OCID: %s", err)
		return client, fmt.Errorf("failed to parse OCI KMS key OCID: %w", err)
	}

	cryptoEndpointTemplate := fmt.Sprintf("https://%s-crypto.kms.{region}.{secondLevelDomain}", vaultExt)
	cryptoEndpoint := common.StringToRegion(region).EndpointForTemplate("kms", cryptoEndpointTemplate)
	log.WithField("endpoint", cryptoEndpoint).Info("Creating OCI KMS client")

	// Use injected config provider if available, otherwise create a fresh one
	cfg := key.configProvider
	if cfg == nil {
		cfg, err = configurationProvider()
		if err != nil {
			return client, fmt.Errorf("failed to create OCI configuration provider: %w", err)
		}
	}

	client, err = keymanagement.NewKmsCryptoClientWithConfigurationProvider(cfg, cryptoEndpoint)
	if err != nil {
		return client, fmt.Errorf("failed to create OCI KMS client: %w", err)
	}

	// Inject custom HTTP client if provided (for testing)
	if key.httpClient != nil {
		client.HTTPClient = key.httpClient
	}

	return client, nil
}

func extractRefs(key *MasterKey) (string, string, error) {
	parts := strings.Split(key.Ocid, ".")
	if len(parts) != ocidParts {
		return "", "", fmt.Errorf("invalid OCID format '%s': expected %d parts, got %d", key.Ocid, ocidParts, len(parts))
	}
	region := parts[3]
	vaultExt := parts[4]
	return region, vaultExt, nil
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Encrypt takes a sops data key, encrypts it with OCI KMS and stores the result
// in the EncryptedKey field.
//
// Consider using EncryptContext instead.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a sops data key, encrypts it with OCI KMS and stores the result
// in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	c, err := key.createCryptoClient()
	if err != nil {
		log.WithField("ocid", key.Ocid).Info("Encryption failed")
		return fmt.Errorf("failed to create OCI KMS service: %w", err)
	}

	data := base64.StdEncoding.EncodeToString(dataKey)

	res, err := c.Encrypt(ctx, keymanagement.EncryptRequest{
		EncryptDataDetails: keymanagement.EncryptDataDetails{
			KeyId:     common.String(key.Ocid),
			Plaintext: &data,
		},
		RequestMetadata: common.RequestMetadata{},
	})

	if err != nil {
		log.WithError(err).WithField("ocid", key.Ocid).
			Error("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with OCI KMS key: %w", err)
	}

	key.EncryptedKey = *res.EncryptedData.Ciphertext
	log.WithField("ocid", key.Ocid).Info("Encryption succeeded")

	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with OCI KMS and returns the result.
//
// Consider using DecryptContext instead.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey field with OCI KMS and returns the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	c, err := key.createCryptoClient()
	if err != nil {
		log.WithField("ocid", key.Ocid).Info("Decryption failed")
		return nil, fmt.Errorf("failed to create OCI KMS service: %w", err)
	}

	res, err := c.Decrypt(ctx, keymanagement.DecryptRequest{
		DecryptDataDetails: keymanagement.DecryptDataDetails{
			Ciphertext: &key.EncryptedKey,
			KeyId:      &key.Ocid,
		},
	})

	if err != nil {
		log.WithError(err).WithField("ocid", key.Ocid).Error("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with OCI KMS key: %w", err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(*res.Plaintext)
	if err != nil {
		log.WithError(err).WithField("ocid", key.Ocid).Error("Decryption failed")
		return nil, fmt.Errorf("failed to base64 decode OCI KMS decrypted key: %w", err)
	}

	log.WithField("ocid", key.Ocid).Info("Decryption succeeded")
	return plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > ocikmsTTL
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.Ocid
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["ocid"] = key.Ocid
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// ConfigurationProvider is a wrapper around common.ConfigurationProvider used for
// authentication towards OCI KMS.
type ConfigurationProvider struct {
	provider common.ConfigurationProvider
}

// NewConfigurationProvider creates a new ConfigurationProvider with the provided
// common.ConfigurationProvider.
func NewConfigurationProvider(cp common.ConfigurationProvider) *ConfigurationProvider {
	return &ConfigurationProvider{provider: cp}
}

// ApplyToMasterKey configures the ConfigurationProvider on the provided key.
func (c ConfigurationProvider) ApplyToMasterKey(key *MasterKey) {
	key.configProvider = c.provider
}

// HTTPClient is a wrapper around common.HTTPRequestDispatcher used for
// configuring the OCI KMS client HTTP requests.
type HTTPClient struct {
	client common.HTTPRequestDispatcher
}

// NewHTTPClient creates a new HTTPClient with the provided common.HTTPRequestDispatcher.
func NewHTTPClient(hc common.HTTPRequestDispatcher) *HTTPClient {
	return &HTTPClient{client: hc}
}

// ApplyToMasterKey configures the HTTP client on the provided key.
func (h HTTPClient) ApplyToMasterKey(key *MasterKey) {
	key.httpClient = h.client
}
