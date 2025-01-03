package yckms

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/getsops/sops/v3/logging"
	"github.com/sirupsen/logrus"
	yckms "github.com/yandex-cloud/go-genproto/yandex/cloud/kms/v1"
	ycsdk "github.com/yandex-cloud/go-sdk"
	"github.com/yandex-cloud/go-sdk/gen/kmscrypto"
	"github.com/yandex-cloud/go-sdk/iamkey"
	"google.golang.org/grpc"
	"os"
	"strings"
	"time"
)

const (
	// kmsTTL is the duration after which a MasterKey requires rotation.
	kmsTTL                     = time.Hour * 24 * 30 * 6
	SopsYandexCloudIAMTokenEnv = "YC_TOKEN"
	SopsYandexCloudSAFileEnv   = "YC_SERVICE_ACCOUNT_KEY_FILE"
	KeyTypeIdentifier          = "yc_kms"
)

var (
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("YCKMS")
}

// MasterKey is a YC KMS key used to encrypt and decrypt the SOPS
// data key.
type MasterKey struct {
	KeyID string
	// EncryptedKey is the string returned after encrypting with YC KMS.
	EncryptedKey string
	// CreationDate is the creation timestamp of the MasterKey. Used
	// for NeedsRotation.
	CreationDate time.Time

	credentials ycsdk.Credentials

	// grpcConn can be used to inject a custom YC KMS client connection.
	// Mostly useful for testing at present, to wire the client to a mock
	// server.
	grpcConn *grpc.ClientConn
}

func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

func NewMasterKeyFromKeyID(keyID string) *MasterKey {
	return &MasterKey{
		KeyID:        keyID,
		CreationDate: time.Now().UTC(),
	}
}

func NewMasterKeyFromKeyIDString(keyID string) []*MasterKey {
	var keys []*MasterKey
	if keyID == "" {
		return keys
	}
	for _, s := range strings.Split(keyID, ",") {
		keys = append(keys, NewMasterKeyFromKeyID(strings.TrimSpace(s)))
	}
	return keys
}

// YCCredentials is a ycsdk.Credentials used for authenticating towards YC KMS
type YCCredentials struct {
	credentials ycsdk.Credentials
}

// NewYCCredentials creates a new YCCredentials with the provided ycsdk.Credentials.
func NewYCCredentials(credentials ycsdk.Credentials) *YCCredentials {
	return &YCCredentials{credentials: credentials}
}

// ApplyToMasterKey configures the TokenCredential on the provided key.
func (c YCCredentials) ApplyToMasterKey(key *MasterKey) {
	key.credentials = c.credentials
}

func (key *MasterKey) Encrypt(dataKey []byte) (err error) {
	client, err := key.newKMSClient()
	if err != nil {
		log.WithError(err).WithField("keyID", key.KeyID).Error("Encryption failed")
		return fmt.Errorf("cannot create YC KMS service: %w", err)
	}

	ciphertextResponse, err := client.Encrypt(context.Background(), &yckms.SymmetricEncryptRequest{
		KeyId:     key.KeyID,
		Plaintext: dataKey,
	})
	if err != nil {
		log.WithError(err).WithField("keyID", key.KeyID).Error("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with YC KMS key: %w", err)
	}
	key.EncryptedKey = base64.StdEncoding.EncodeToString(ciphertextResponse.Ciphertext)
	log.WithField("resourceID", key.KeyID).Info("Encryption succeeded")
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

// Decrypt decrypts the EncryptedKey field with YC KMS and returns
// the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	client, err := key.newKMSClient()
	if err != nil {
		log.WithError(err).WithField("keyID", key.KeyID).Error("Decryption failed")
		return nil, fmt.Errorf("cannot create YC KMS service: %w", err)
	}

	decodedCipher, err := base64.StdEncoding.DecodeString(string(key.EncryptedDataKey()))
	plaintextResponse, err := client.Decrypt(context.Background(), &yckms.SymmetricDecryptRequest{
		KeyId:      key.KeyID,
		Ciphertext: decodedCipher,
	})
	if err != nil {
		log.WithError(err).WithField("keyID", key.KeyID).Error("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with YC KMS key: %w", err)
	}
	log.WithField("resourceID", key.KeyID).Info("Decryption succeeded")
	return plaintextResponse.Plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (kmsTTL)
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return key.KeyID
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key *MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["key_id"] = key.KeyID
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// newKMSClient returns a YC KMS client configured with the credentialsStore
// and/or grpcConn, falling back to environmental defaults.
// It returns an error if the ResourceID is invalid, or if the setup of the
// client fails.
func (key *MasterKey) newKMSClient() (*kms.SymmetricCryptoServiceClient, error) {
	var (
		cred ycsdk.Credentials
		err  error
	)

	switch {
	case key.credentials != nil:
		cred = key.credentials
	default:
		cred, err = getYandexCloudCredentials()
		if err != nil {
			return nil, err
		}
	}

	client, err := ycsdk.Build(context.Background(), ycsdk.Config{
		Credentials: cred,
	})
	if err != nil {
		return nil, err
	}

	if key.grpcConn != nil {
		return kms.NewKMSCrypto(func(ctx context.Context) (*grpc.ClientConn, error) {
			return key.grpcConn, nil
		}).SymmetricCrypto(), nil
	}

	return client.KMSCrypto().SymmetricCrypto(), nil
}

// getYandexCloudCredentials trying to locate credentials in the following order
// 1. Service account. Env variable contains either a path to or the contents of the Service Account file in JSON format.
// 2. IAM token. You can get it via `yc iam create-token`
// 3. Instance metadata
func getYandexCloudCredentials() (ycsdk.Credentials, error) {
	_, exists := os.LookupEnv(SopsYandexCloudSAFileEnv)
	if exists {
		key, err := getServiceAccountCredentials()
		if err != nil {
			return nil, err
		}
		saKey, err := iamkey.ReadFromJSONBytes(key)
		if err != nil {
			return nil, err
		}
		return ycsdk.ServiceAccountKey(saKey)
	}

	token, exists := os.LookupEnv(SopsYandexCloudIAMTokenEnv)
	if exists {
		return ycsdk.NewIAMTokenCredentials(token), nil
	}

	return ycsdk.InstanceServiceAccount(), nil
}

func getServiceAccountCredentials() ([]byte, error) {
	serviceAccount := os.Getenv(SopsYandexCloudSAFileEnv)
	if _, err := os.Stat(serviceAccount); err == nil {
		return os.ReadFile(serviceAccount)
	}
	return []byte(serviceAccount), nil
}
