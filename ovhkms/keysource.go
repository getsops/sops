package ovhkms

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/ovh/okms-sdk-go"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify an OVH KMS MasterKey.
	KeyTypeIdentifier = "ovh_kms"

	// CertificateFileEnv is the environment variable containing the path to the certificate file.
	CertificateFileEnv = "OVH_CERTIFICATE_FILE"

	// CertificateKeyFileEnv is the environment variable containing the path to the certificate key file.
	CertificateKeyFileEnv = "OVH_CERTIFICATE_KEY_FILE"
)

var (
	// log is the global logger for any OVH KMS MasterKey.
	log *logrus.Logger

	// ovhKmsTTL is the duration after which a MasterKey requires rotation.
	ovhKmsTTL = time.Hour * 24 * 30 * 6

	// ovhEndpointRegex is used to validate and extract components from the OVH KMS endpoint.
	ovhEndpointRegex = regexp.MustCompile(`^([^/]+)/([^/]+)$`)
)

func init() {
	log = logging.NewLogger("OVH_KMS")
}

// MasterKey is an OVH KMS key used to encrypt and decrypt SOPS' data key.
type MasterKey struct {
	// Endpoint is the OVH KMS endpoint.
	Endpoint string

	// KeyID is the UUID of the key in OVH KMS.
	KeyID string

	// EncryptedKey contains the SOPS data key encrypted with the OVH KMS key.
	EncryptedKey string

	// CreationDate of the MasterKey, used to determine if the EncryptedKey
	// needs rotation.
	CreationDate time.Time

	// certificateFile is the path to the certificate file used for authentication.
	certificateFile string

	// certificateKeyFile is the path to the certificate key file used for authentication.
	certificateKeyFile string
}

// NewMasterKeyFromKeyID creates a new MasterKey from an OVH KMS key ID
// in the format <endpoint>/<key-id>.
func NewMasterKeyFromKeyID(keyID string) (*MasterKey, error) {
	matches := ovhEndpointRegex.FindStringSubmatch(keyID)
	if matches == nil || len(matches) != 3 {
		// If no match, assume it's just a key ID and endpoint not provider
		return nil, fmt.Errorf("not a vaild key (should be like this: +"+
			"example.okms.ovh.net/keyId), got: %v", keyID)
	}

	return &MasterKey{
		Endpoint:     matches[1],
		KeyID:        matches[2],
		CreationDate: time.Now().UTC(),
	}, nil
}

// MasterKeysFromResourceIDString creates a list of MasterKeys from a comma-separated
// string of OVH KMS key IDs.
func MasterKeysFromResourceIDString(resourceID string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if resourceID == "" {
		return keys, nil
	}

	for _, resourceID := range strings.Split(resourceID, ",") {
		if resourceID == "" {
			continue
		}
		mk, err := NewMasterKeyFromKeyID(resourceID)
		if err != nil {
			return nil, err
		}
		keys = append(keys, mk)
	}

	return keys, nil
}

// Encrypt takes a SOPS data key, encrypts it with OVH KMS, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

// EncryptContext takes a SOPS data key, encrypts it with OVH KMS, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	// UUID key
	KeyID, err := uuid.Parse(key.KeyID)
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Encryption failed: invalid UUID")
		return fmt.Errorf("failed to parse UUID '%s': %w", key.KeyID, err)
	}

	client, err := key.getClient()
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Encryption failed")
		return err
	}

	// Base64 encode the data key
	plaintext := base64.StdEncoding.EncodeToString(dataKey)

	// Encrypt the data key using OVH KMS
	resp, err := client.Encrypt(ctx, KeyID, "", []byte(plaintext))
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with OVH KMS key '%s': %w", key.KeyID, err)
	}

	// Store the encrypted key
	key.EncryptedKey = resp
	log.WithField("KeyID", key.KeyID).Info("Encryption successful")
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

// Decrypt decrypts the EncryptedKey field with OVH KMS and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey field with OVH KMS and returns the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	// UUID key
	KeyID, err := uuid.Parse(key.KeyID)
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Decryption failed: invalid UUID")
		return nil, fmt.Errorf("failed to parse UUID '%s': %w", key.KeyID, err)
	}

	client, err := key.getClient()
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Decryption failed")
		return nil, err
	}

	// Decrypt the data key using OVH KMS
	resp, err := client.Decrypt(ctx, KeyID, "", string(key.EncryptedKey))
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with OVH KMS key '%s': %w", key.KeyID, err)
	}

	// Decode the base64 plaintext
	dataKey, err := base64.StdEncoding.DecodeString(string(resp))
	if err != nil {
		log.WithField("KeyID", key.KeyID).Info("Decryption failed: invalid base64 plaintext")
		return nil, fmt.Errorf("failed to decode base64 plaintext: %w", err)
	}

	log.WithField("KeyID", key.KeyID).Info("Decryption successful")
	return dataKey, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (ovhKmsTTL)
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return fmt.Sprintf("%s/%s", key.Endpoint, key.KeyID)
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["endpoint"] = key.Endpoint
	out["key_id"] = key.KeyID
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// getClient returns an OVH KMS client configured with the certificate
// from environment variables.
func (key *MasterKey) getClient() (*okms.Client, error) {
	// Get certificate paths from environment variables if not set
	certFile := key.certificateFile
	certKeyFile := key.certificateKeyFile

	if certFile == "" {
		certFile = os.Getenv(CertificateFileEnv)
	}

	if certKeyFile == "" {
		certKeyFile = os.Getenv(CertificateKeyFileEnv)
	}

	if certFile == "" || certKeyFile == "" {
		return nil, fmt.Errorf("OVH KMS certificate file paths not provided. Set %s and %s environment variables",
			CertificateFileEnv, CertificateKeyFileEnv)
	}

	// Create OVH KMS client
	cert, err := tls.LoadX509KeyPair(certFile, certKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate file: %w", err)
	}

	httpClient := http.Client{
		Transport: &http.Transport{TLSClientConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}},
	}
	client, err := okms.NewRestAPIClientWithHttp(fmt.Sprintf("https://%s", key.Endpoint), &httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create OVH KMS client: %w", err)
	}

	return client, nil
}

// SetCertificateFiles sets the certificate file paths for authentication.
func (key *MasterKey) SetCertificateFiles(certFile, certKeyFile string) {
	key.certificateFile = certFile
	key.certificateKeyFile = certKeyFile
}
