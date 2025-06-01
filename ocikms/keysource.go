package ocikms

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/oracle/oci-go-sdk/v65/common"
	kms "github.com/oracle/oci-go-sdk/v65/keymanagement"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify a OCI KMS MasterKey.
	KeyTypeIdentifier = "oci_kms"
)

var (
	// ocikmsTTL is the duration after which a MasterKey requires rotation.
	ocikmsTTL = time.Hour * 24 * 30 * 6
	// log is the global logger for any OCI KMS MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("OCIKMS")
}

// MasterKey is a OCI KMS key used to encrypt and decrypt the SOPS
// data key.
type MasterKey struct {
	// KeyVersionId is OCID of key version used in combination of Key ID
	KeyVersionId string
	// Id is OCID of Master Ecnryption Key stored in OCI Vault
	Id string
	// CryptoEndpoint is URL of Crypto Endpoint, url can differ between regions and types of key storage
	CryptoEndpoint string
	// CreationDate is time of creation
	CreationDate time.Time
	// EncryptedKey is the string returned after encrypting with OCI KMS.
	EncryptedKey string
}

// NewMasterKeyFromResourceID creates a new MasterKey with the provided ocid
func NewMasterKey(cryptoEndpoint string, ocid string, version string) *MasterKey {
	return &MasterKey{
		Id:             ocid,
		CreationDate:   time.Now().UTC(),
		KeyVersionId:   version,
		CryptoEndpoint: cryptoEndpoint,
	}
}

// MasterKeysFromResourceIDString takes a comma separated list of OCI KMS
// resource IDs and returns a slice of new MasterKeys for them.
func NewMasterKeyFromURL(url string) (*MasterKey, error) {
	url = strings.TrimSpace(url)
	re := regexp.MustCompile(`^(https://[^/]+)/(ocid1\.key\.[^/]+)/(ocid1\.keyversion\.[^/]+)/?$`)
	parts := re.FindStringSubmatch(url)
	if len(parts) != 4 {
		return nil, fmt.Errorf("could not parse %q into a valid OCI Key Vault Master Encryption Key", url)
	}
	return NewMasterKey(parts[1], parts[2], parts[3]), nil
}

func MasterKeysFromURLs(urls string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if urls == "" {
		return keys, nil
	}
	for _, s := range strings.Split(urls, ",") {
		k, err := NewMasterKeyFromURL(s)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// Encrypt takes a SOPS data key, encrypts it with OCI KMS, and stores the
// result in the EncryptedKey field.
//
// Consider using EncryptContext instead.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a SOPS data key, encrypts it with OCI KMS, and stores the
// result in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	service, err := key.newKMSClient()
	if err != nil {
		log.WithField("ocid", key.Id).Info("Encryption failed")
		return fmt.Errorf("cannot create OCI KMS service: %w", err)
	}

	// OCI KMS reguires Plaintext is base64 encoded
	plainText := base64.StdEncoding.EncodeToString(dataKey)
	req := kms.EncryptRequest{
		EncryptDataDetails: kms.EncryptDataDetails{
			KeyId:        &key.Id,
			Plaintext:    &plainText,
			KeyVersionId: &key.KeyVersionId,
		},
	}
	resp, err := service.Encrypt(ctx, req)
	if err != nil {
		log.WithField("ocid", key.Id).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with OCI KMS key: %w", err)
	}

	key.EncryptedKey = *resp.Ciphertext
	log.WithField("ocid", key.Id).Info("Encryption succeeded")
	return nil
}

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been
// encrypted yet.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with OCI KMS and returns
// the result.
//
// Consider using DecryptContext instead.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey field with OCI KMS and returns
// the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	service, err := key.newKMSClient()
	if err != nil {
		log.WithField("ocid", key.Id).Info("Decryption failed")
		return nil, fmt.Errorf("cannot create OCI KMS service: %w", err)
	}

	cipher := string(key.EncryptedDataKey())

	req := kms.DecryptRequest{
		DecryptDataDetails: kms.DecryptDataDetails{
			KeyId:        &key.Id,
			Ciphertext:   &cipher,
			KeyVersionId: &key.KeyVersionId,
		},
	}
	resp, err := service.Decrypt(ctx, req)
	if err != nil {
		log.WithField("ocid", key.Id).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with OCI KMS key: %w", err)
	}
	decodedPlainText, err := base64.StdEncoding.DecodeString(string(*resp.Plaintext))
	if err != nil {
		log.WithField("ocid", key.Id).Info("Decryption failed")
		return nil, err
	}
	log.WithField("ocid", key.Id).Info("Decryption succeeded")
	return []byte(decodedPlainText), nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (ocikmsTTL)
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return fmt.Sprintf("%s/%s/%s", key.CryptoEndpoint, key.Id, key.KeyVersionId)
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["id"] = key.Id
	out["crypto_endpoint"] = key.CryptoEndpoint
	out["key_version"] = key.KeyVersionId
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

func (key *MasterKey) newKMSClient() (*kms.KmsCryptoClient, error) {
	// Right now only implements authentication using default profile
	configProvider := common.DefaultConfigProvider()
	client, err := kms.NewKmsCryptoClientWithConfigurationProvider(configProvider, key.CryptoEndpoint)
	if err != nil {
		return nil, err
	}

	return &client, nil
}
