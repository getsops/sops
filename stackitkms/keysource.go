/*
Package stackitkms contains an implementation of the github.com/getsops/sops/v3/keys.MasterKey
interface that encrypts and decrypts the data key using STACKIT KMS.
*/
package stackitkms // import "github.com/getsops/sops/v3/stackitkms"

import (
	"context"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	stackitconfig "github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/kms"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify a STACKIT KMS MasterKey.
	KeyTypeIdentifier = "stackit_kms"
	// stackitKmsTTL is the duration after which a MasterKey requires rotation.
	stackitKmsTTL = time.Hour * 24 * 30 * 6
)

var (
	// log is the global logger for any STACKIT KMS MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("STACKITKMS")
}

// MasterKey is a STACKIT KMS key used to encrypt and decrypt SOPS' data key.
type MasterKey struct {
	// ResourceID is the full resource identifier in format:
	// projects/<projectId>/regions/<regionId>/keyRings/<keyRingId>/keys/<keyId>/versions/<versionNumber>
	ResourceID string
	// ProjectID is the STACKIT project UUID.
	ProjectID string
	// RegionID is the STACKIT region (e.g., "eu01").
	RegionID string
	// KeyRingID is the key ring UUID.
	KeyRingID string
	// KeyID is the key UUID.
	KeyID string
	// VersionNumber is the key version number.
	VersionNumber int64
	// EncryptedKey stores the data key in its encrypted form.
	EncryptedKey string
	// CreationDate is when this MasterKey was created.
	CreationDate time.Time

	// configOpts holds STACKIT SDK configuration options.
	// They can be injected by a (local) keyservice.KeyServiceServer using
	// Credentials.ApplyToMasterKey.
	// If nil, the default credential provider chain is used.
	configOpts []stackitconfig.ConfigurationOption
}

// NewMasterKey creates a new MasterKey from a resource ID string, setting
// the creation date to the current date.
func NewMasterKey(resourceID string) (*MasterKey, error) {
	projectID, regionID, keyRingID, keyID, versionNumber, err := parseResourceID(resourceID)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		ResourceID:    resourceID,
		ProjectID:     projectID,
		RegionID:      regionID,
		KeyRingID:     keyRingID,
		KeyID:         keyID,
		VersionNumber: versionNumber,
		CreationDate:  time.Now().UTC(),
	}, nil
}

// NewMasterKeyFromResourceIDString takes a comma separated list of STACKIT KMS
// resource IDs and returns a slice of new MasterKeys.
func NewMasterKeyFromResourceIDString(resourceID string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if resourceID == "" {
		return keys, nil
	}
	for _, s := range strings.Split(resourceID, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		k, err := NewMasterKey(s)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// parseResourceID parses a resource ID in format:
// projects/<projectId>/regions/<regionId>/keyRings/<keyRingId>/keys/<keyId>/versions/<versionNumber>
func parseResourceID(resourceID string) (projectID, regionID, keyRingID, keyID string, versionNumber int64, err error) {
	resourceID = strings.TrimSpace(resourceID)
	parts := strings.Split(resourceID, "/")
	if len(parts) != 10 {
		return "", "", "", "", 0, fmt.Errorf("invalid STACKIT KMS resource ID format: expected 'projects/<projectId>/regions/<regionId>/keyRings/<keyRingId>/keys/<keyId>/versions/<versionNumber>', got %q", resourceID)
	}
	if parts[0] != "projects" || parts[2] != "regions" || parts[4] != "keyRings" || parts[6] != "keys" || parts[8] != "versions" {
		return "", "", "", "", 0, fmt.Errorf("invalid STACKIT KMS resource ID format: expected 'projects/<projectId>/regions/<regionId>/keyRings/<keyRingId>/keys/<keyId>/versions/<versionNumber>', got %q", resourceID)
	}
	projectID = parts[1]
	regionID = parts[3]
	keyRingID = parts[5]
	keyID = parts[7]
	versionNumber, err = strconv.ParseInt(parts[9], 10, 64)
	if err != nil {
		return "", "", "", "", 0, fmt.Errorf("invalid version number in STACKIT KMS resource ID %q: %w", resourceID, err)
	}
	if projectID == "" || regionID == "" || keyRingID == "" || keyID == "" {
		return "", "", "", "", 0, fmt.Errorf("all components must be non-empty in STACKIT KMS resource ID %q", resourceID)
	}
	return projectID, regionID, keyRingID, keyID, versionNumber, nil
}

// Credentials is a wrapper around STACKIT SDK configuration options used
// for authentication towards STACKIT KMS.
type Credentials struct {
	opts []stackitconfig.ConfigurationOption
}

// NewCredentials returns a Credentials object with the provided configuration options.
func NewCredentials(opts ...stackitconfig.ConfigurationOption) *Credentials {
	return &Credentials{opts: opts}
}

// ApplyToMasterKey configures the credentials on the provided key.
func (c Credentials) ApplyToMasterKey(key *MasterKey) {
	key.configOpts = c.opts
}

// Encrypt takes a SOPS data key, encrypts it with STACKIT KMS and stores the result
// in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a SOPS data key, encrypts it with STACKIT KMS and stores the result
// in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	client, err := key.createKMSClient()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Encryption failed")
		return fmt.Errorf("failed to create STACKIT KMS client: %w", err)
	}

	plaintext := base64.StdEncoding.EncodeToString(dataKey)
	plaintextBytes := []byte(plaintext)

	result, err := client.Encrypt(ctx, key.ProjectID, key.RegionID, key.KeyRingID, key.KeyID, key.VersionNumber).
		EncryptPayload(kms.EncryptPayload{
			Data: &plaintextBytes,
		}).Execute()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with STACKIT KMS: %w", err)
	}

	encryptedData := result.GetData()
	if len(encryptedData) == 0 {
		return fmt.Errorf("encryption response missing ciphertext")
	}
	key.EncryptedKey = base64.StdEncoding.EncodeToString(encryptedData)
	log.WithField("resourceID", key.ResourceID).Info("Encryption succeeded")
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

// Decrypt decrypts the EncryptedKey with STACKIT KMS and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey with STACKIT KMS and returns the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	client, err := key.createKMSClient()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to create STACKIT KMS client: %w", err)
	}

	encryptedBytes, err := base64.StdEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to base64 decode encrypted data key: %w", err)
	}

	result, err := client.Decrypt(ctx, key.ProjectID, key.RegionID, key.KeyRingID, key.KeyID, key.VersionNumber).
		DecryptPayload(kms.DecryptPayload{
			Data: &encryptedBytes,
		}).Execute()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with STACKIT KMS: %w", err)
	}

	decryptedData := result.GetData()
	if len(decryptedData) == 0 {
		return nil, fmt.Errorf("decryption response missing plaintext")
	}

	plaintext, err := base64.StdEncoding.DecodeString(string(decryptedData))
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to base64 decode decrypted data key: %w", err)
	}

	log.WithField("resourceID", key.ResourceID).Info("Decryption succeeded")
	return plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > stackitKmsTTL
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return key.ResourceID
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["resource_id"] = key.ResourceID
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// createKMSClient creates a STACKIT KMS client with the appropriate credentials.
func (key *MasterKey) createKMSClient() (*kms.APIClient, error) {
	client, err := kms.NewAPIClient(key.configOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create STACKIT KMS API client: %w", err)
	}
	return client, nil
}
