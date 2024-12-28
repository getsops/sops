package ocikms

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/getsops/sops/v3/logging"

	"github.com/oracle/oci-go-sdk/v65/common"
	"github.com/oracle/oci-go-sdk/v65/keymanagement"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

const (
	// cryptoEndpointTemplate is the template for the OCI KMS crypto endpoint that is constructed using parts of the key OCID
	cryptoEndpointTemplate = "https://%s-crypto.kms.%s.oraclecloud.com"
	// ocidParts is the number of parts in an OCID, separated by ".", eg: "ocid1.key.oc1.uk-london-1.aaaalgz5aacmg.aaaailjtjbkbc5ufsorrihgv2agugpfe7wrtngukihgkybqxcoozz7sbh6lq"
	ocidParts = 6
	// KeyTypeIdentifier is the string used to identify an OCI KMS MasterKey.
	KeyTypeIdentifier = "oci_kms"
)

func init() {
	log = logging.NewLogger("OCIKMS")
}

// MasterKey is an Oracle Cloud KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	Ocid         string
	EncryptedKey string
	CreationDate time.Time
}

func NewMasterKeyFromOCID(ocid string) *MasterKey {
	return &MasterKey{
		Ocid:         strings.Replace(ocid, " ", "", -1),
		CreationDate: time.Now().UTC(),
	}
}

func MasterKeysFromOCIDString(ocids string) []*MasterKey {
	var keys []*MasterKey
	if ocids == "" {
		return keys
	}
	for _, s := range strings.Split(ocids, ",") {
		keys = append(keys, NewMasterKeyFromOCID(s))
	}
	return keys
}

// createKeyManagementClient creates a new OCI KMS client
func (key *MasterKey) createCryptoClient() (client keymanagement.KmsCryptoClient, err error) {
	region, vault_ref, err := extractRefs(key)
	if err != nil {
		log.WithField("ocid", key.Ocid).Errorf("Cannot extract region and vault_ref from OCID: %s", err)
	}

	endpoint := fmt.Sprintf(cryptoEndpointTemplate, vault_ref, region)
	log.WithField("endpoint", endpoint).Info("Creating OCI KMS client")
	// The client is created using the default OCI config provider, using the default profile in the default config file (~/.oci/config)
	// There is currently no straightforward way to pass a custom config provider to the client.
	// The oci-go-sdk provides a way to pass a custom config provider to the client, but there's no environment variable to feature-flag it.
	// Related: https://github.com/oracle/oci-go-sdk/issues/318

	// In order to use a custom provider, the client would need to be created like this:
	// client, err := keymanagement.NewKmsCryptoClientWithConfigurationProvider(common.CustomProfileConfigProvider("/home/<user>/.oci/config", "<profile>"), endpoint)
	// Sticking with the defaults for now.

	client, err = keymanagement.NewKmsCryptoClientWithConfigurationProvider(common.DefaultConfigProvider(), endpoint)
	if err != nil {
		return client, fmt.Errorf("Cannot create OCI KMS client: %w", err)
	}
	return client, nil
}

func extractRefs(key *MasterKey) (string, string, error) {
	parts := strings.Split(key.Ocid, ".")
	if len(parts) != ocidParts {
		return "", "", fmt.Errorf("OCID length is %s, expected %d", key.Ocid, ocidParts)
	}
	region := parts[3]
	vault_ref := parts[4]
	return region, vault_ref, nil
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Encrypt takes a sops data key, encrypts it with Key Vault and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	c, err := key.createCryptoClient()
	if err != nil {
		log.WithField("ocid", key.Ocid).Info("Encryption failed")
		return fmt.Errorf("cannot create OCI KMS service: %w", err)
	}
	data := base64.StdEncoding.EncodeToString(dataKey)

	res, err := c.Encrypt(context.TODO(), keymanagement.EncryptRequest{
		EncryptDataDetails: keymanagement.EncryptDataDetails{
			KeyId:     common.String(key.Ocid),
			Plaintext: &data,
		},
		RequestMetadata: common.RequestMetadata{},
	})

	if err != nil {
		log.WithError(err).WithField("ocid", key.Ocid).
			Error("Encryption failed")
		return fmt.Errorf("failed to encrypt data: %w", err)
	}

	key.EncryptedKey = *res.EncryptedData.Ciphertext
	log.WithField("ocid", key.Ocid).Info("Encryption succeeded")

	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with Azure Key Vault and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	c, err := key.createCryptoClient()
	if err != nil {
		log.WithField("ocid", key.Ocid).Info("Decryption failed")
		return nil, err
	}

	res, err := c.Decrypt(context.TODO(), keymanagement.DecryptRequest{
		DecryptDataDetails: keymanagement.DecryptDataDetails{
			Ciphertext: &key.EncryptedKey,
			KeyId:      &key.Ocid,
		},
	})

	if err != nil {
		log.WithError(err).WithField("ocid", key.Ocid).Error("Decryption failed")
		return nil, fmt.Errorf("error decrypting key: %w", err)
	}

	plaintext, err := base64.StdEncoding.DecodeString(*res.Plaintext)
	if err != nil {
		log.WithError(err).WithField("ocid", key.Ocid).Error("Decryption failed")
		return nil, err
	}

	log.WithField("ocid", key.Ocid).Info("Decryption succeeded")
	return plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.Ocid
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["ocid"] = key.Ocid
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}
