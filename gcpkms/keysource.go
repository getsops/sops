package gcpkms //import "go.mozilla.org/sops/gcpkms"

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	cloudkms "google.golang.org/api/cloudkms/v1"
)

// MasterKey is a GCP KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	ResourceId   string
	EncryptedKey string
	CreationDate time.Time
}

// Encrypt takes a sops data key, encrypts it with GCP KMS and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	cloudkmsService, err := key.createCloudKMSService()
	if err != nil {
		return fmt.Errorf("Cannot create GCP KMS service: %v", err)
	}
	req := &cloudkms.EncryptRequest{
		Plaintext: base64.StdEncoding.EncodeToString(dataKey),
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Encrypt(key.ResourceId, req).Do()
	if err != nil {
		return fmt.Errorf("Failed to call GCP KMS encryption service: %v", err)
	}

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
		return nil, fmt.Errorf("Cannot create GCP KMS service: %v", err)
	}

	req := &cloudkms.DecryptRequest{
		Ciphertext: key.EncryptedKey,
	}
	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.Decrypt(key.ResourceId, req).Do()
	if err != nil {
		return nil, fmt.Errorf("Error decrypting key: %v", err)
	}
	return base64.StdEncoding.DecodeString(resp.Plaintext)
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.ResourceId
}

// NewMasterKeyFromResourceId takes a GCP KMS resource ID string and returns a new MasterKey for that
func NewMasterKeyFromResourceId(resourceId string) *MasterKey {
	k := &MasterKey{}
	resourceId = strings.Replace(resourceId, " ", "", -1)
	k.ResourceId = resourceId
	k.CreationDate = time.Now().UTC()
	return k
}

// MasterKeysFromResourceIdString takes a comma separated list of GCP KMS resourece IDs and returns a slice of new MasterKeys for them
func MasterKeysFromResourceIdString(resourceId string) []*MasterKey {
	var keys []*MasterKey
	if resourceId == "" {
		return keys
	}
	for _, s := range strings.Split(resourceId, ",") {
		keys = append(keys, NewMasterKeyFromResourceId(s))
	}
	return keys
}

func (key MasterKey) createCloudKMSService() (*cloudkms.Service, error) {
	re := regexp.MustCompile(`^projects/[^/]+/locations/[^/]+/keyRings/[^/]+/cryptoKeys/[^/]+$`)
	matches := re.FindStringSubmatch(key.ResourceId)
	if matches == nil {
		return nil, fmt.Errorf("No valid resoureceId found in %q", key.ResourceId)
	}

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}
	return cloudkmsService, nil
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["resource_id"] = key.ResourceId
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	return out
}
