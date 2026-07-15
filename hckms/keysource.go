/*
Package hckms contains an implementation of the github.com/getsops/sops/v3/keys.MasterKey
interface that encrypts and decrypts the data key using HuaweiCloud KMS with the SDK
for Go V3.
*/
package hckms // import "github.com/getsops/sops/v3/hckms"

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/provider"
	huaweikms "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kms/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kms/v2/model"
	kmsregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kms/v2/region"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify a HuaweiCloud KMS MasterKey.
	KeyTypeIdentifier = "hckms"
	// hckmsTTL is the duration after which a MasterKey requires rotation.
	hckmsTTL = time.Hour * 24 * 30 * 6
)

var (
	// log is the global logger for any HuaweiCloud KMS MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("HCKMS")
}

// MasterKey is a HuaweiCloud KMS key used to encrypt and decrypt SOPS' data key.
type MasterKey struct {
	// KeyID is the full key identifier in format "region:key-uuid"
	KeyID string
	// Region is the HuaweiCloud region (e.g., "tr-west-1")
	Region string
	// KeyUUID is the UUID of the KMS key
	KeyUUID string
	// EncryptedKey stores the data key in its encrypted form.
	EncryptedKey string
	// CreationDate is when this MasterKey was created.
	CreationDate time.Time

	// credentials contains the HuaweiCloud credentials used by the KMS client.
	// It can be injected by a (local) keyservice.KeyServiceServer using
	// Credentials.ApplyToMasterKey.
	// If nil, the default credential provider chain is used.
	credentials auth.ICredential
}

// NewMasterKey creates a new MasterKey from a region:key-id string, setting
// the creation date to the current date.
func NewMasterKey(keyID string) (*MasterKey, error) {
	region, keyUUID, err := parseKeyID(keyID)
	if err != nil {
		return nil, err
	}
	return &MasterKey{
		KeyID:        keyID,
		Region:       region,
		KeyUUID:      keyUUID,
		CreationDate: time.Now().UTC(),
	}, nil
}

// NewMasterKeyFromKeyIDString takes a comma separated list of HuaweiCloud KMS
// key IDs in format "region:key-uuid", and returns a slice of new MasterKeys.
func NewMasterKeyFromKeyIDString(keyID string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if keyID == "" {
		return keys, nil
	}
	for _, s := range strings.Split(keyID, ",") {
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

// parseKeyID parses a key ID in format "region:key-uuid" and returns the region and UUID.
func parseKeyID(keyID string) (string, string, error) {
	keyID = strings.TrimSpace(keyID)
	parts := strings.SplitN(keyID, ":", 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid key ID format: expected 'region:key-uuid', got %q", keyID)
	}
	region := strings.TrimSpace(parts[0])
	keyUUID := strings.TrimSpace(parts[1])
	if region == "" {
		return "", "", fmt.Errorf("region cannot be empty in key ID: %q", keyID)
	}
	if keyUUID == "" {
		return "", "", fmt.Errorf("key UUID cannot be empty in key ID: %q", keyID)
	}
	return region, keyUUID, nil
}

// Credentials is a wrapper around auth.ICredential used for authentication
// towards HuaweiCloud KMS.
type Credentials struct {
	credential auth.ICredential
}

// NewCredentials returns a Credentials object with the provided auth.ICredential.
func NewCredentials(c auth.ICredential) *Credentials {
	return &Credentials{credential: c}
}

// ApplyToMasterKey configures the credentials on the provided key.
func (c Credentials) ApplyToMasterKey(key *MasterKey) {
	key.credentials = c.credential
}

// Encrypt takes a SOPS data key, encrypts it with HuaweiCloud KMS and stores the result
// in the EncryptedKey field.
//
// Consider using EncryptContext instead.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a SOPS data key, encrypts it with HuaweiCloud KMS and stores the result
// in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	client, err := key.createKMSClient(ctx)
	if err != nil {
		log.WithField("keyID", key.KeyID).Info("Encryption failed")
		return fmt.Errorf("failed to create HuaweiCloud KMS client: %w", err)
	}

	plaintext := base64.StdEncoding.EncodeToString(dataKey)
	encryptAlgorithm := model.GetEncryptDataRequestBodyEncryptionAlgorithmEnum().SYMMETRIC_DEFAULT

	request := &model.EncryptDataRequest{
		Body: &model.EncryptDataRequestBody{
			KeyId:               key.KeyUUID,
			PlainText:           plaintext,
			EncryptionAlgorithm: &encryptAlgorithm,
		},
	}

	response, err := client.EncryptData(request)
	if err != nil {
		log.WithField("keyID", key.KeyID).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with HuaweiCloud KMS: %w", err)
	}

	if response.CipherText == nil {
		return fmt.Errorf("encryption response missing ciphertext")
	}
	key.EncryptedKey = *response.CipherText
	log.WithField("keyID", key.KeyID).Info("Encryption succeeded")
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

// Decrypt decrypts the EncryptedKey with HuaweiCloud KMS and returns the result.
//
// Consider using DecryptContext instead.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey with HuaweiCloud KMS and returns the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	client, err := key.createKMSClient(ctx)
	if err != nil {
		log.WithField("keyID", key.KeyID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to create HuaweiCloud KMS client: %w", err)
	}

	decryptAlgorithm := model.GetDecryptDataRequestBodyEncryptionAlgorithmEnum().SYMMETRIC_DEFAULT

	request := &model.DecryptDataRequest{
		Body: &model.DecryptDataRequestBody{
			CipherText:          key.EncryptedKey,
			EncryptionAlgorithm: &decryptAlgorithm,
			KeyId:               &key.KeyUUID,
		},
	}

	response, err := client.DecryptData(request)
	if err != nil {
		log.WithField("keyID", key.KeyID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with HuaweiCloud KMS: %w", err)
	}

	if response.PlainText == nil {
		return nil, fmt.Errorf("decryption response missing plaintext")
	}
	decrypted, err := base64.StdEncoding.DecodeString(*response.PlainText)
	if err != nil {
		log.WithField("keyID", key.KeyID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to base64 decode decrypted data key: %w", err)
	}

	log.WithField("keyID", key.KeyID).Info("Decryption succeeded")
	return decrypted, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > hckmsTTL
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return key.KeyID
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["key_id"] = key.KeyID
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// createKMSClient creates a HuaweiCloud KMS client with the appropriate credentials
// and region configuration.
func (key *MasterKey) createKMSClient(ctx context.Context) (*huaweikms.KmsClient, error) {
	var cred auth.ICredential
	var err error

	if key.credentials != nil {
		cred = key.credentials
	} else {
		// Use default credential provider chain (env -> profile -> metadata)
		credentialProviderChain := provider.BasicCredentialProviderChain()
		cred, err = credentialProviderChain.GetCredentials()
		if err != nil {
			return nil, fmt.Errorf("failed to get HuaweiCloud credentials: %w", err)
		}
	}

	// Get KMS region with endpoint
	reg, err := kmsregion.SafeValueOf(key.Region)
	if err != nil {
		return nil, fmt.Errorf("invalid region %q: %w", key.Region, err)
	}

	// Create HTTP client builder
	hcClientBuilder := core.NewHcHttpClientBuilder().
		WithCredential(cred).
		WithRegion(reg)

	hcClient := hcClientBuilder.Build()

	// Create KMS client
	kmsClient := huaweikms.NewKmsClient(hcClient)
	return kmsClient, nil
}
