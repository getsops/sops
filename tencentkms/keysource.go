/*
Package tencentkms contains an implementation of the github.com/getsops/sops/v3/keys.MasterKey
interface that encrypts and decrypts the data key using Tencent Cloud KMS with the
Tencent Cloud SDK for Go.
*/
package tencentkms // import "github.com/getsops/sops/v3/tencentkms"

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/getsops/sops/v3/logging"
	"github.com/sirupsen/logrus"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
)

const (
	// KeyTypeIdentifier is the string used to identify a Tencent Cloud KMS MasterKey.
	KeyTypeIdentifier = "tencent_kms"
	// TencentKmsEnvVar is the environment variable name for Tencent Cloud KMS key IDs.
	TencentKmsEnvVar = "SOPS_TENCENT_KMS_IDS"

	// TencentSecretIdEnvVar is the environment variable name for Tencent Cloud SecretId.
	TencentSecretIdEnvVar = "TENCENTCLOUD_SECRET_ID"
	// TencentSecretKeyEnvVar is the environment variable name for Tencent Cloud SecretKey.
	TencentSecretKeyEnvVar = "TENCENTCLOUD_SECRET_KEY"
	// TencentTokenEnvVar is the environment variable name for Tencent Cloud Token (optional, for STS).
	TencentTokenEnvVar = "TENCENTCLOUD_TOKEN"
	// TencentRegionEnvVar is the environment variable name for Tencent Cloud region.
	TencentRegionEnvVar = "TENCENTCLOUD_REGION"
	// TencentKMSEndpointEnvVar is the environment variable name for Tencent Cloud kms service endpoint.
	TencentKMSEndpointEnvVar = "TENCENTCLOUD_KMS_ENDPOINT"
)

var (
	// log is the global logger for any Tencent Cloud KMS MasterKey.
	log *logrus.Logger
	// tencent kms TTL is the duration after which a MasterKey requires rotation.
	tencentkmsTTL = time.Hour * 24 * 30 * 6
)

func init() {
	log = logging.NewLogger("TENCENT_KMS")
}

// MasterKey is a Tencent Cloud KMS Key used to Encrypt and Decrypt SOPS' data key.
type MasterKey struct {
	// KeyID is the ID of the Tencent Cloud KMS key.
	KeyID string
	// Region is the region of the Tencent Cloud KMS key.
	Region string
	// EncryptedKey contains the SOPS data key encrypted with the Tencent Cloud KMS key.
	EncryptedKey string
	// CreationDate of the MasterKey, used to determine if the EncryptedKey needs rotation.
	CreationDate time.Time

	// secretId is the Tencent Cloud SecretId used for authentication.
	secretId string
	// secretKey is the Tencent Cloud SecretKey used for authentication.
	secretKey string
	// token is the Tencent Cloud STS token used for authentication (optional).
	token string
}

// NewMasterKeyFromKeyID creates a new MasterKey with the provided Key ID.
func NewMasterKeyFromKeyID(keyID string) *MasterKey {
	k := &MasterKey{}
	keyID = strings.Replace(keyID, " ", "", -1)
	k.KeyID = keyID
	k.CreationDate = time.Now().UTC()
	return k
}

// MasterKeysFromKeyIDString takes a comma separated list of Tencent KMS
// KeyIDs and returns a slice of new MasterKeys for them.
func MasterKeysFromKeyIDString(keyID string) []*MasterKey {
	var keys []*MasterKey
	if keyID == "" {
		return keys
	}
	for _, s := range strings.Split(keyID, ",") {
		keys = append(keys, NewMasterKeyFromKeyID(s))
	}
	return keys
}

// Encrypt takes a SOPS data key, encrypts it with Tencent Cloud KMS, and stores the result in the EncryptedKey field.
// Consider using EncryptContext instead.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a SOPS data key, encrypts it with Tencent Cloud KMS, and stores the result in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	client, err := key.createClient()
	if err != nil {
		log.WithFields(logrus.Fields{"keyId": key.KeyID}).WithError(err).Error("Failed to create Tencent Cloud KMS client")
		return fmt.Errorf("failed to create Tencent Cloud KMS client: %w", err)
	}

	// Create encryption request
	request := kms.NewEncryptRequest()
	request.KeyId = common.StringPtr(key.KeyID)
	request.Plaintext = common.StringPtr(base64.StdEncoding.EncodeToString(dataKey))

	// Send encryption request
	response, err := client.EncryptWithContext(ctx, request)
	if err != nil {
		log.WithFields(logrus.Fields{"keyId": key.KeyID}).WithError(err).Error("Failed to encrypt data key with Tencent Cloud KMS")
		return fmt.Errorf("failed to encrypt data key with Tencent Cloud KMS: %w", err)
	}

	// Store the encrypted key
	key.EncryptedKey = *response.Response.CiphertextBlob
	log.WithFields(logrus.Fields{"keyId": key.KeyID}).Info("Encryption successful")
	return nil
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been encrypted yet.
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
func (key *MasterKey) SetEncryptedDataKey(encrypted []byte) {
	key.EncryptedKey = string(encrypted)
}

// Decrypt decrypts the EncryptedKey field with Tencent Cloud KMS and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey field with Tencent Cloud KMS and returns the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	if key.EncryptedKey == "" {
		return nil, fmt.Errorf("master key is empty")
	}

	client, err := key.createClient()
	if err != nil {
		log.WithFields(logrus.Fields{"keyId": key.KeyID}).WithError(err).Error("Failed to create Tencent Cloud KMS client")
		return nil, fmt.Errorf("failed to create Tencent Cloud KMS client: %w", err)
	}

	// Create decryption request
	request := kms.NewDecryptRequest()
	request.CiphertextBlob = common.StringPtr(key.EncryptedKey)

	// Send decryption request
	response, err := client.DecryptWithContext(ctx, request)
	if err != nil {
		log.WithFields(logrus.Fields{"keyId": key.KeyID}).WithError(err).Error("Failed to decrypt data key with Tencent Cloud KMS")
		return nil, fmt.Errorf("failed to decrypt data key with Tencent Cloud KMS: %w", err)
	}

	decodedCipher, err := base64.StdEncoding.DecodeString(*response.Response.Plaintext)
	if err != nil {
		log.WithField("keyId", key.KeyID).WithError(err).Error("Failed to decode decrypted plaintext")
		return nil, err
	}

	log.WithFields(logrus.Fields{"keyId": key.KeyID}).Info("Decryption successful")
	return decodedCipher, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (tencentkmsTTL)
}

// ToString converts the master key to a string representation.
func (key *MasterKey) ToString() string {
	return key.KeyID
}

// ToMap converts the master key to a map representation.
func (key *MasterKey) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"keyId":      key.KeyID,
		"created_at": key.CreationDate.UTC().Format(time.RFC3339),
		"enc":        key.EncryptedKey,
	}
}

// TypeToIdentifier returns the key type identifier.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// createClient creates a new Tencent Cloud KMS client with support for multiple authentication methods.
func (key *MasterKey) createClient() (*kms.Client, error) {
	credential, err := key.getCredential()
	if err != nil {
		log.WithFields(logrus.Fields{"keyId": key.KeyID}).WithError(err).Error("Failed to obtain Tencent Cloud credentials")
		return nil, fmt.Errorf("failed to obtain Tencent Cloud credentials: %w", err)
	}

	if credential == nil {
		credErr := fmt.Errorf("no valid credentials found. Please set TENCENTCLOUD_SECRET_ID and TENCENTCLOUD_SECRET_KEY environment variables")
		log.WithFields(logrus.Fields{"keyId": key.KeyID}).WithError(credErr).Error("Failed to obtain Tencent Cloud credentials")
		return nil, credErr
	}

	region := os.Getenv(TencentRegionEnvVar)
	if region == "" {
		region = "ap-guangzhou"
	}

	cpf := profile.NewClientProfile()
	endpoint := os.Getenv(TencentKMSEndpointEnvVar)
	if endpoint != "" {
		cpf.HttpProfile.Endpoint = endpoint
	}

	client, err := kms.NewClient(credential, region, cpf)
	if err != nil {
		return nil, fmt.Errorf("failed to create Tencent KMS client for region %s: %w", region, err)
	}

	return client, nil
}

// getCredential gets authentication credentials, supports multiple methods
func (key *MasterKey) getCredential() (common.CredentialIface, error) {
	if secretId, ok := os.LookupEnv(TencentSecretIdEnvVar); ok && len(secretId) > 0 {
		key.secretId = secretId
	}

	if secretKey, ok := os.LookupEnv(TencentSecretKeyEnvVar); ok && len(secretKey) > 0 {
		key.secretKey = secretKey
	}

	if key.secretId == "" && key.secretKey == "" {
		return nil, fmt.Errorf("environment variable TENCENTCLOUD_SECRET_ID or TENCENTCLOUD_SECRET_KEY is not set")
	}

	if token, ok := os.LookupEnv(TencentTokenEnvVar); ok && len(token) > 0 {
		key.token = token
	}

	if key.token != "" {
		return common.NewTokenCredential(key.secretId, key.secretKey, key.token), nil
	}

	return common.NewCredential(key.secretId, key.secretKey), nil
}