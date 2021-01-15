/*
 * Package aliyunkms contains an implementation of go.mozilla.org/sops/v3/keys.MasterKey interface
 * that encrypts and decrypts the data key using Aliyun KMS with the Aliyun Go SDK.
 */
package aliyunkms

import (
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/kms"
	"github.com/sirupsen/logrus"

	"go.mozilla.org/sops/v3/logging"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("ALIYUNKMS")
}

// MasterKey is a Aliyun KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	Role            string
	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	KeyID           string
	EncryptedKey    string
	CreationDate    time.Time
}

// NewMasterKeyWithEcsRamRole creates a new MasterKey from a role and regionId, setting the creation date to the current date
func NewMasterKeyWithEcsRamRole(regionId string, roleName string, keyID string) *MasterKey {
	return &MasterKey{
		RegionID:     regionId,
		Role:         roleName,
		KeyID:        keyID,
		CreationDate: time.Now().UTC(),
	}
}

func (key MasterKey) createCloudKMSService() (*kms.Client, error) {
	cloudkmsService, err := kms.NewClientWithAccessKey(key.RegionID, key.AccessKeyID, key.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return cloudkmsService, nil
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Encrypt takes a sops data key, encrypts it with KMS and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	cloudkmsService, err := key.createCloudKMSService()
	if err != nil {
		log.WithField("role", key.Role).Info("Encryption failed")
		return fmt.Errorf("Cannot create Aliyun KMS service: %v", err)
	}
	request := kms.CreateEncryptRequest()
	request.KeyId = key.KeyID
	request.Scheme = "https"
	request.Plaintext = string(dataKey)
	response, err := cloudkmsService.Encrypt(request)
	if err != nil {
		log.WithField("role", key.Role).Info("Encryption failed")
		return fmt.Errorf("Failed to call Aliyun KMS encryption service: %v", err)
	}
	log.WithField("role", key.Role).Info("Encryption succeeded")
	key.EncryptedKey = response.CiphertextBlob
	return nil
}

// Decrypt decrypts the EncryptedKey field with Aliyun KMS and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	cloudkmsService, err := key.createCloudKMSService()
	if err != nil {
		log.WithField("role", key.Role).Info("Decryption failed")
		return nil, fmt.Errorf("Cannot create Aliyun KMS service: %v", err)
	}

	request := kms.CreateDecryptRequest()
	request.Scheme = "https"
	request.CiphertextBlob = string(key.EncryptedKey)
	response, err := cloudkmsService.Decrypt(request)
	if err != nil {
		log.WithField("role", key.Role).Info("Decryption failed")
		return nil, fmt.Errorf("Error decrypting key: %v", err)
	}
	log.WithField("role", key.Role).Info("Decryption succeeded")
	return []byte(response.Plaintext), nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.Role
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	if key.Role != "" {
		out["role"] = key.Role
	}
	out["key_id"] = key.KeyID
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}
