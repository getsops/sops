package gcpkms //import "go.mozilla.org/sops/v3/gcpkms"

import (
	"encoding/base64"
	"fmt"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"go.mozilla.org/sops/v3/logging"

	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("GCPKMS")
}

// MasterKey is a GCP KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	ResourceID   string
	EncryptedKey string
	CreationDate time.Time
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Encrypt takes a sops data key, encrypts it with GCP KMS and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	cloudkmsService, err := key.createCloudKMSService()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Encryption failed")
		return fmt.Errorf("Cannot create GCP KMS service: %w", err)
	}
	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(dataKey),
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(key.ResourceID, req).Do()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Encryption failed")
		return fmt.Errorf("Failed to call GCP KMS encryption service: %w", err)
	}
	log.WithField("resourceID", key.ResourceID).Info("Encryption succeeded")
	key.EncryptedKey = resp.Ciphertext
	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with CGP KMS and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	cloudkmsService, err := key.createCloudKMSService()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("Cannot create GCP KMS service: %w", err)
	}

	req := &cloudkms.DecryptRequest{
		Ciphertext: key.EncryptedKey,
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(key.ResourceID, req).Do()
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, fmt.Errorf("Error decrypting key: %w", err)
	}
	encryptedKey, err := base64.StdEncoding.DecodeString(resp.Plaintext)
	if err != nil {
		log.WithField("resourceID", key.ResourceID).Info("Decryption failed")
		return nil, err
	}
	log.WithField("resourceID", key.ResourceID).Info("Decryption succeeded")
	return encryptedKey, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.ResourceID
}

// NewMasterKeyFromResourceID takes a GCP KMS resource ID string and returns a new MasterKey for that
func NewMasterKeyFromResourceID(resourceID string) *MasterKey {
	k := &MasterKey{}
	resourceID = strings.Replace(resourceID, " ", "", -1)
	k.ResourceID = resourceID
	k.CreationDate = time.Now().UTC()
	return k
}

// MasterKeysFromResourceIDString takes a comma separated list of GCP KMS resource IDs and returns a slice of new MasterKeys for them
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

func (key MasterKey) createCloudKMSService() (*cloudkms.Service, error) {
	re := regexp.MustCompile(`^projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+$`)
	matches := re.FindStringSubmatch(key.ResourceID)
	if matches == nil {
		return nil, fmt.Errorf("No valid resourceId found in %q", key.ResourceID)
	}

	ctx := context.Background()

	creds, err := getDefaultApplicationCredentials()
	if err != nil {
		return nil, err
	}

	cloudkmsService, err := cloudkms.NewService(ctx, option.WithCredentialsJSON(creds))
	if err != nil {
		return nil, err
	}
	return cloudkmsService, nil
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["resource_id"] = key.ResourceID
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	return out
}

// getDefaultApplicationCredentials allows for passing GCP Service Account
// Credentials as either a path to a file, or directly as an environment variable
// in JSON format.
func getDefaultApplicationCredentials() (token []byte, err error) {
	var defaultCredentials = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")

	if _, err := os.Stat(defaultCredentials); err == nil {
		if token, err = ioutil.ReadFile(defaultCredentials); err != nil {
			return nil, err
		}
	} else {
		token = []byte(defaultCredentials)
	}
	return
}