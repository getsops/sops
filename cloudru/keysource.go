package cloudru

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	kmsV1 "github.com/cloudru-tech/key-manager-sdk/api/v1"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/getsops/sops/v3/logging"
)

const (
	// EnvDiscoveryURL is the environment variable used to set the Cloudru.
	EnvDiscoveryURL = "CLOUDRU_DISCOVERY_URL"
	// EnvAccessKeyID is the environment variable used to set the Cloudru API access key ID.
	EnvAccessKeyID = "CLOUDRU_ACCESS_KEY_ID"
	// EnvAccessKeySecret is the environment variable used to set the Cloudru API access key secret.
	EnvAccessKeySecret = "CLOUDRU_ACCESS_KEY_SECRET"

	// KeyTypeIdentifier is the string used to identify a Cloudru KMS MasterKey.
	KeyTypeIdentifier = "cloudru_kms"

	// DiscoveryURL is the default Cloudru API discovery URL, which is used to
	// retrieve the actual API addresses of Cloudru products.
	DiscoveryURL = "https://api.cloud.ru/endpoints"
)

var (
	// ErrInvalidKeyID is returned when the key ID format is invalid.
	ErrInvalidKeyID = errors.New("key id is invalid: it should be a UUID with an optional version number separated by a colon, e.g.: 123e4567-e89b-12d3-a456-426614174000:1")
)

func init() {
	logger = logging.NewLogger("CLOUDRU_KMS")
}

var (
	// log is the global logger for any Vault Transit MasterKey.
	logger *logrus.Logger
)

// AccessKey is used to retrieve the cloudru API access token.
type AccessKey struct {
	KeyID  string
	Secret string
}

// MasterKey is a Cloudru KMS backend path, which is used to Encrypt and Decrypt SOPS data key.
type MasterKey struct {
	// KeyID is the identifier is used to refer to the cloudru kms key.
	KeyID string
	// Cipher is the value returned after encrypting with Cloudru KMS.
	Cipher []byte
	// CreatedAt is the date and time when the MasterKey was created.
	// Used for NeedsRotation.
	CreatedAt time.Time

	// credentials is the AccessKey used for authenticating towards the Cloudru.
	credentials AccessKey
	// conn is the gRPC connection to the Cloudru KMS server.
	conn *grpc.ClientConn
}

// NewMasterKeyFromKeyID creates a new MasterKey with the provided KMS key ID.
func NewMasterKeyFromKeyID(kid string) (*MasterKey, error) {
	parts := strings.SplitN(kid, ":", 2)
	if len(parts) == 2 {
		_, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, ErrInvalidKeyID
		}
	}

	parsed, err := uuid.Parse(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid KMS key id: %s", err)
	}
	if parsed == uuid.Nil {
		return nil, errors.New("invalid KMS key id: key id should be a valid and non-zero UUID")
	}

	return &MasterKey{
		KeyID:     kid,
		CreatedAt: time.Now().UTC(),
	}, nil
}

// NewMasterKeysFromKeyIDs creates a new MasterKey with the provided KMS key IDs.
func NewMasterKeysFromKeyIDs(kids string) ([]*MasterKey, error) {
	if kids == "" {
		return nil, nil
	}

	var keys []*MasterKey
	for _, kid := range strings.Split(kids, ",") {
		key, err := NewMasterKeyFromKeyID(kid)
		if err != nil {
			return nil, fmt.Errorf("parse key '%s': %s", kid, err)
		}
		keys = append(keys, key)
	}

	return keys, nil
}

// Encrypt takes a SOPS data key, encrypts it with Vault Transit, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	c, err := provideClient()
	if err != nil {
		return fmt.Errorf("initialize cloud.ru API: %w", err)
	}
	defer func() {
		if closeErr := c.Close(); closeErr != nil {
			logger.Errorf("failed to close cloud.ru KMS client connection: %s", closeErr)
		}
	}()

	kid, version, err := parseKeyID(key.KeyID)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &kmsV1.EncryptRequest{
		KeyId:        kid,
		KeyVersionId: version,
		Plaintext:    wrapperspb.Bytes(dataKey),
	}
	resp, err := c.KMS.Encrypt(ctx, req)
	if err != nil {
		return fmt.Errorf("encrypt data with key_id '%s': %s", key.KeyID, err)
	}

	key.Cipher = resp.Ciphertext.GetValue()
	logger.WithField("key_id", key.KeyID).Info("Encryption successful")
	return nil
}

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been
// encrypted yet.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if len(key.Cipher) == 0 {
		return key.Encrypt(dataKey)
	}
	return nil
}

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte { return key.Cipher }

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) { key.Cipher = enc }

// Decrypt decrypts the EncryptedKey field with Vault Transit and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	c, err := provideClient()
	if err != nil {
		return nil, fmt.Errorf("initialize cloud.ru API: %w", err)
	}
	defer func() {
		if closeErr := c.Close(); closeErr != nil {
			logger.Errorf("failed to close cloud.ru KMS client connection: %s", closeErr)
		}
	}()

	kid, _, err := parseKeyID(key.KeyID)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &kmsV1.DecryptRequest{
		KeyId:      kid,
		Ciphertext: wrapperspb.Bytes(key.Cipher),
	}
	resp, err := c.KMS.Decrypt(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("decrypt data with key_id '%s': %s", key.KeyID, err)
	}

	logger.WithField("key_id", key.KeyID).Info("Decryption successful")
	return resp.Plaintext.GetValue(), nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	// TODO: research this value
	return time.Since(key.CreatedAt) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string { return key.KeyID }

// ToMap converts the MasterKey to a map for serialization purposes.
func (key *MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["key_id"] = key.KeyID
	out["created_at"] = key.CreatedAt.UTC().Format(time.RFC3339)
	out["enc"] = string(key.Cipher)
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string { return KeyTypeIdentifier }

func parseKeyID(kid string) (key string, version int32, err error) {
	parts := strings.SplitN(kid, ":", 2)
	key = parts[0]
	if len(parts) == 2 {
		var parsed int64
		parsed, err = strconv.ParseInt(parts[1], 10, 32)
		if err != nil {
			return "", 0, ErrInvalidKeyID
		}

		if parsed > 0 && parsed <= math.MaxInt32 {
			version = int32(parsed)
		}
	}
	return key, version, nil
}
