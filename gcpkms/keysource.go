package gcpkms // import "go.mozilla.org/sops/v3/gcpkms"

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"os"
	"regexp"
	"strings"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"github.com/sirupsen/logrus"
	"google.golang.org/api/option"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"go.mozilla.org/sops/v3/logging"
)

const (
	// SopsGoogleCredentialsEnv can be set as an environment variable as either
	// a path to a credentials file, or directly as the variable's value in JSON
	// format.
	SopsGoogleCredentialsEnv = "GOOGLE_CREDENTIALS"
)

var (
	// gcpkmsTTL is the duration after which a MasterKey requires rotation.
	gcpkmsTTL = time.Hour * 24 * 30 * 6
	// log is the global logger for any GCP KMS MasterKey.
	log *logrus.Logger
)

func init() {
	log = logging.NewLogger("GCPKMS")
}

// MasterKey is a GCP KMS key used to encrypt and decrypt the SOPS
// data key.
type MasterKey struct {
	// ResourceID is the resource id used to refer to the gcp kms key.
	// It can be retrieved using the `gcloud` command.
	ResourceID string
	// EncryptedKey is the string returned after encrypting with GCP KMS.
	EncryptedKey string
	// CreationDate is the creation timestamp of the MasterKey. Used
	// for NeedsRotation.
	CreationDate time.Time

	// credentialJSON is the Service Account credentials JSON used for
	// authenticating towards the GCP KMS service.
	credentialJSON []byte
	// grpcConn can be used to inject a custom GCP client connection.
	// Mostly useful for testing at present, to wire the client to a mock
	// server.
	grpcConn *grpc.ClientConn
}

// NewMasterKeyFromResourceID creates a new MasterKey with the provided resource
// ID.
func NewMasterKeyFromResourceID(resourceID string) *MasterKey {
	k := &MasterKey{}
	resourceID = strings.Replace(resourceID, " ", "", -1)
	k.ResourceID = resourceID
	k.CreationDate = time.Now().UTC()
	return k
}

// MasterKeysFromResourceIDString takes a comma separated list of GCP KMS
// resource IDs and returns a slice of new MasterKeys for them.
func MasterKeysFromResourceIDString(resourceID string) []*MasterKey {
	var keys []*MasterKey
	if resourceID == "" {
		return keys
	}
	for _, s := range strings.Split(resourceID, ",") {
		keys = append(keys, NewMasterKeyFromResourceID(s))
	}
	return keys
}

// CredentialJSON is the Service Account credentials JSON used for authenticating
// towards the GCP KMS service.
type CredentialJSON []byte

// ApplyToMasterKey configures the CredentialJSON on the provided key.
func (c CredentialJSON) ApplyToMasterKey(key *MasterKey) {
	key.credentialJSON = c
}

// Encrypt takes a SOPS data key, encrypts it with GCP KMS, and stores the
// result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	service, err := key.newKMSClient()
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Encryption failed")
		return fmt.Errorf("cannot create GCP KMS service: %w", err)
	}
	defer func() {
		if err := service.Close(); err != nil {
			log.WithError(err).Error("failed to close GCP KMS client connection")
		}
	}()

	ctx := context.Background()
	purpose, err := key.purpose(ctx, service)
	if err != nil {
		return err
	}

	switch purpose {
	case kmspb.CryptoKey_ENCRYPT_DECRYPT:
		return key.encryptSymmetric(ctx, service, dataKey)
	case kmspb.CryptoKey_ASYMMETRIC_DECRYPT:
		return key.encryptAsymmetric(ctx, service, dataKey)
	default:
		log.WithField("resourceID", key.ResourceID).WithField("purpose", purpose.String()).Error("This key is not for encryption")
		return fmt.Errorf("this key is not for encryption, purpose: %v", purpose.String())
	}
}

func (key *MasterKey) purpose(ctx context.Context, service *kms.KeyManagementClient) (kmspb.CryptoKey_CryptoKeyPurpose, error) {
	req := &kmspb.GetCryptoKeyRequest{
		Name: key.resourceIDWithoutVersion(),
	}
	cryptoKey, err := service.GetCryptoKey(ctx, req)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Get key metadata failed")
		return kmspb.CryptoKey_CRYPTO_KEY_PURPOSE_UNSPECIFIED, fmt.Errorf("failed to get key metadata from GCP KMS service: %w", err)
	}

	return cryptoKey.GetPurpose(), nil
}

// assume key.ResourceID is in following format
//   - `projects/project-id/locations/location/keyRings/keyring/cryptoKeys/key`
//   - `projects/project-id/locations/location/keyRings/keyring/cryptoKeys/key/cryptoKeyVersions/version`
func (key MasterKey) resourceIDWithoutVersion() string {
	re := regexp.MustCompile(`^(projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+)(?:/cryptoKeyVersions/[^/]+)?$`)
	matches := re.FindStringSubmatch(key.ResourceID)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}

func (key *MasterKey) encryptSymmetric(ctx context.Context, service *kms.KeyManagementClient, dataKey []byte) error {
	req := &kmspb.EncryptRequest{
		Name:      key.ResourceID,
		Plaintext: dataKey,
	}
	resp, err := service.Encrypt(ctx, req)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with GCP KMS key: %w", err)
	}
	// NB: base64 encoding is for compatibility with SOPS <=3.8.x.
	// The previous GCP KMS client used to work with base64 encoded
	// strings.
	key.EncryptedKey = base64.StdEncoding.EncodeToString(resp.Ciphertext)
	log.WithField("resourceID", key.ResourceID).Info("Symmetric encryption succeeded")
	return nil
}

func (key *MasterKey) encryptAsymmetric(ctx context.Context, service *kms.KeyManagementClient, dataKey []byte) error {
	req := &kmspb.GetPublicKeyRequest{
		Name: key.ResourceID,
	}
	resp, err := service.GetPublicKey(ctx, req)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Get public key failed")
		return fmt.Errorf("failed to get public key from GCP KMS service: %w", err)
	}

	if resp.GetPemCrc32C().GetValue() != wrapperspb.Int64(int64(crc32c([]byte(resp.GetPem())))).Value {
		log.WithField("resourceID", key.ResourceID).Error("Get public key response corrupted in-transit")
		return errors.New("get public key response corrupted in-transit")
	}

	block, _ := pem.Decode([]byte(resp.Pem))
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Info("Failed to parse public key")
		return fmt.Errorf("Failed to parse public key: %w", err)
	}
	rsaKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		log.WithField("resourceID", key.ResourceID).Info("Public key is not RSA")
		return errors.New("public key is not RSA")
	}

	var hash hash.Hash

	switch resp.GetAlgorithm() {
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_2048_SHA256:
		hash = sha256.New()
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_3072_SHA256:
		hash = sha256.New()
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_4096_SHA256:
		hash = sha256.New()
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_4096_SHA512:
		hash = sha512.New()
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_2048_SHA1:
		hash = sha1.New()
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_3072_SHA1:
		hash = sha1.New()
	case kmspb.CryptoKeyVersion_RSA_DECRYPT_OAEP_4096_SHA1:
		hash = sha1.New()
	default:
		log.WithField("resourceID", key.ResourceID).WithField("algorithm", resp.GetAlgorithm().String()).Error("Unsupported algorithm")
		return fmt.Errorf("Key with unsupported algorithm: %s", resp.GetAlgorithm().String())
	}

	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, rsaKey, dataKey, nil)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("rsa.EncryptOAEP() error")
		return fmt.Errorf("rsa.EncryptOAEP: %w", err)
	}

	// NB: base64 encoding is for compatibility with SOPS <=3.8.x.
	// The previous GCP KMS client used to work with base64 encoded
	// strings.
	key.EncryptedKey = base64.StdEncoding.EncodeToString(ciphertext)
	log.WithField("resourceID", key.ResourceID).Info("Asymmetric encryption succeeded")
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

// Decrypt decrypts the EncryptedKey field with GCP KMS and returns
// the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	service, err := key.newKMSClient()
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Decryption failed")
		return nil, fmt.Errorf("cannot create GCP KMS service: %w", err)
	}
	defer func() {
		if err := service.Close(); err != nil {
			log.WithError(err).Error("failed to close GCP KMS client connection")
		}
	}()

	// NB: this is for compatibility with SOPS <=3.8.x. The previous GCP KMS
	// client used to work with base64 encoded strings.
	decodedCipher, err := base64.StdEncoding.DecodeString(string(key.EncryptedDataKey()))
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Decryption failed")
		return nil, err
	}

	ctx := context.Background()

	purpose, err := key.purpose(ctx, service)
	if err != nil {
		return nil, err
	}

	switch purpose {
	case kmspb.CryptoKey_ENCRYPT_DECRYPT:
		return key.decryptSymmetric(ctx, service, decodedCipher)
	case kmspb.CryptoKey_ASYMMETRIC_DECRYPT:
		return key.decryptAsymmetric(ctx, service, decodedCipher)
	default:
		log.WithField("resourceID", key.ResourceID).WithField("purpose", purpose.String()).Info("This key cannot be used for decryption")
		return nil, fmt.Errorf("This key cannot be used for decryption, purpose: %s", purpose.String())
	}
}

func (key *MasterKey) decryptSymmetric(ctx context.Context, service *kms.KeyManagementClient, decodedCipher []byte) ([]byte, error) {
	req := &kmspb.DecryptRequest{
		Name:       key.ResourceID,
		Ciphertext: decodedCipher,
	}
	resp, err := service.Decrypt(ctx, req)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Symmetric decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with GCP KMS key: %w", err)
	}

	log.WithField("resourceID", key.ResourceID).Info("Symmetric decryption succeeded")
	return resp.Plaintext, nil
}

func (key *MasterKey) decryptAsymmetric(ctx context.Context, service *kms.KeyManagementClient, decodedCipher []byte) ([]byte, error) {
	req := &kmspb.AsymmetricDecryptRequest{
		Name:       key.ResourceID,
		Ciphertext: decodedCipher,
	}
	resp, err := service.AsymmetricDecrypt(ctx, req)
	if err != nil {
		log.WithError(err).WithField("resourceID", key.ResourceID).Error("Asymmetric decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with GCP KMS key: %w", err)
	}

	log.WithField("resourceID", key.ResourceID).Info("Asymmetric decryption succeeded")
	return resp.Plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (gcpkmsTTL)
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

// newKMSClient returns a GCP KMS client configured with the credentialJSON
// and/or grpcConn, falling back to environmental defaults.
// It returns an error if the ResourceID is invalid, or if the setup of the
// client fails.
func (key *MasterKey) newKMSClient() (*kms.KeyManagementClient, error) {
	re := regexp.MustCompile(`^projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+(?:/cryptoKeyVersions/[^/]+)?$`)
	matches := re.FindStringSubmatch(key.ResourceID)
	if matches == nil {
		return nil, fmt.Errorf("no valid resource ID found in %q", key.ResourceID)
	}

	var opts []option.ClientOption
	switch {
	case key.credentialJSON != nil:
		opts = append(opts, option.WithCredentialsJSON(key.credentialJSON))
	default:
		credentials, err := getGoogleCredentials()
		if err != nil {
			return nil, err
		}
		if len(credentials) > 0 {
			opts = append(opts, option.WithCredentialsJSON(credentials))
		}
	}
	if key.grpcConn != nil {
		opts = append(opts, option.WithGRPCConn(key.grpcConn))
	}

	ctx := context.Background()
	client, err := kms.NewKeyManagementClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// getGoogleCredentials returns the SopsGoogleCredentialsEnv variable, as
// either the file contents of the path of a credentials file, or as value in
// JSON format. It returns an error if the file cannot be read, and may return
// a nil byte slice if no value is set.
func getGoogleCredentials() ([]byte, error) {
	defaultCredentials := os.Getenv(SopsGoogleCredentialsEnv)
	if _, err := os.Stat(defaultCredentials); err == nil {
		return os.ReadFile(defaultCredentials)
	}
	return []byte(defaultCredentials), nil
}

func crc32c(data []byte) uint32 {
	t := crc32.MakeTable(crc32.Castagnoli)
	return crc32.Checksum(data, t)
}
