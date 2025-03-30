package gcpkms // import "github.com/getsops/sops/v3/gcpkms"

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	kms "cloud.google.com/go/kms/apiv1"
	"cloud.google.com/go/kms/apiv1/kmspb"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/grpc"

	"github.com/getsops/sops/v3/logging"
)

const (
	// SopsGoogleCredentialsEnv can be set as an environment variable as either
	// a path to a credentials file, or directly as the variable's value in JSON
	// format.
	SopsGoogleCredentialsEnv = "GOOGLE_CREDENTIALS"
	// SopsGoogleCredentialsOAuthTokenEnv is the environment variable used for the
	// GCP OAuth 2.0 Token.
	SopsGoogleCredentialsOAuthTokenEnv = "GOOGLE_OAUTH_ACCESS_TOKEN"
	// KeyTypeIdentifier is the string used to identify a GCP KMS MasterKey.
	KeyTypeIdentifier = "gcp_kms"
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

	// tokenSource contains the oauth2.TokenSource used by the GCP client.
	// It can be injected by a (local) keyservice.KeyServiceServer using
	// TokenSource.ApplyToMasterKey.
	// If nil, the remaining authentication methods are attempted.
	tokenSource oauth2.TokenSource
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

// TokenSource is an oauth2.TokenSource used for authenticating towards the
// GCP KMS service.
type TokenSource struct {
	source oauth2.TokenSource
}

// NewTokenSource creates a new TokenSource from the provided oauth2.TokenSource.
func NewTokenSource(source oauth2.TokenSource) TokenSource {
	return TokenSource{source: source}
}

// ApplyToMasterKey configures the TokenSource on the provided key.
func (t TokenSource) ApplyToMasterKey(key *MasterKey) {
	key.tokenSource = t.source
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
		log.WithField("resourceID", key.ResourceID).Info("Encryption failed")
		return fmt.Errorf("cannot create GCP KMS service: %w", err)
	}
	defer func() {
		if err := service.Close(); err != nil {
			log.Error("failed to close GCP KMS client connection")
		}
	}()

	req := &kmspb.EncryptRequest{
		Name:      key.ResourceID,
		Plaintext: dataKey,
	}
	ctx := context.Background()
	resp, err := service.Encrypt(ctx, req)
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with GCP KMS key: %w", err)
	}
	// NB: base64 encoding is for compatibility with SOPS <=3.8.x.
	// The previous GCP KMS client used to work with base64 encoded
	// strings.
	key.EncryptedKey = base64.StdEncoding.EncodeToString(resp.Ciphertext)
	log.WithField("resourceID", key.ResourceID).Info("Encryption succeeded")
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
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("cannot create GCP KMS service: %w", err)
	}
	defer func() {
		if err := service.Close(); err != nil {
			log.Error("failed to close GCP KMS client connection")
		}
	}()

	// NB: this is for compatibility with SOPS <=3.8.x. The previous GCP KMS
	// client used to work with base64 encoded strings.
	decodedCipher, err := base64.StdEncoding.DecodeString(string(key.EncryptedDataKey()))
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, err
	}

	req := &kmspb.DecryptRequest{
		Name:       key.ResourceID,
		Ciphertext: decodedCipher,
	}
	ctx := context.Background()
	resp, err := service.Decrypt(ctx, req)
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with GCP KMS key: %w", err)
	}

	log.WithField("resourceID", key.ResourceID).Info("Decryption succeeded")
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

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// newKMSClient returns a GCP KMS client configured with the tokenSource
// or credentialJSON, and/or grpcConn, falling back to environmental defaults.
// It returns an error if the ResourceID is invalid, or if the setup of the
// client fails.
func (key *MasterKey) newKMSClient() (*kms.KeyManagementClient, error) {
	re := regexp.MustCompile(`^projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+$`)
	matches := re.FindStringSubmatch(key.ResourceID)
	if matches == nil {
		return nil, fmt.Errorf("no valid resource ID found in %q", key.ResourceID)
	}

	var opts []option.ClientOption
	switch {
	case key.tokenSource != nil:
		opts = append(opts, option.WithTokenSource(key.tokenSource))
	case key.credentialJSON != nil:
		opts = append(opts, option.WithCredentialsJSON(key.credentialJSON))
	default:
		credentials, err := getGoogleCredentials()
		if err != nil {
			return nil, fmt.Errorf("credentials: failed to obtain credentials from %q: %w", SopsGoogleCredentialsEnv, err)
		}
		if credentials != nil {
			opts = append(opts, option.WithCredentialsJSON(credentials))
			break
		}

		if atCredentials := getGoogleOAuthTokenFromEnv(); atCredentials != nil {
			opts = append(opts, option.WithTokenSource(atCredentials))
			break
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
// JSON format.
// It returns an error and a nil byte slice if the file cannot be read.
func getGoogleCredentials() ([]byte, error) {
	if defaultCredentials, ok := os.LookupEnv(SopsGoogleCredentialsEnv); ok && len(defaultCredentials) > 0 {
		if _, err := os.Stat(defaultCredentials); err == nil {
			return os.ReadFile(defaultCredentials)
		}
		return []byte(defaultCredentials), nil
	}
	return nil, nil
}

// getGoogleOAuthTokenFromEnv returns the SopsGoogleCredentialsOauthTokenEnv variable,
// as the OAauth 2.0 token.
// It returns an error and a nil byte slice if the envrionment variable is not set.
func getGoogleOAuthTokenFromEnv() oauth2.TokenSource {
	if token, ok := os.LookupEnv(SopsGoogleCredentialsOAuthTokenEnv); ok && len(token) > 0 {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		return tokenSource
	}
	return nil
}
