/*
Package barbican contains an implementation of the go.mozilla.org/sops/v3/keys.MasterKey interface that encrypts and decrypts the
data key using OpenStack Barbican using the gophercloud sdk.
*/
package barbican //import "go.mozilla.org/sops/v3/barbican"

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"
	"os/user"
	"strings"
	"time"

	"go.mozilla.org/sops/v3/logging"
	"gopkg.in/ini.v1"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/trusts"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("BARBICAN")
}

// MasterKey is a GCP KMS key used to encrypt and decrypt sops' data key.
type MasterKey struct {
	SecretHref   string
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

// Encrypt takes a sops data key, encrypts it with Barbican and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {

	km, err := newKeyManager()
	if err != nil {
		log.WithField("SecretHref", key.SecretHref).Error("Encryption failed")
		return fmt.Errorf("Cannot create KeyManager service: %v", err)
	}

	secretID, err := parseID(key.SecretHref)
	if err != nil {
		return fmt.Errorf("Failed to parse secret href: %v", err)
	}
	payload, err := secrets.GetPayload(km, secretID, nil).Extract()
	if err != nil {
		return fmt.Errorf("Failed to fetch master key payload: %v", err)
	}
	masterKey := strings.Split(string(payload), "\n")

	enc, err := encrypt(masterKey[0], masterKey[1], dataKey)
	if err != nil {
		return fmt.Errorf("Failed to encrypt data key: %v", err)
	}
	key.SetEncryptedDataKey(enc)

	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with Barbican and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	km, err := newKeyManager()
	if err != nil {
		log.WithField("SecretHref", key.SecretHref).Error("Encryption failed")
		return []byte{}, fmt.Errorf("Cannot create KeyManager service: %v", err)
	}

	secretID, err := parseID(key.SecretHref)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to parse secret href: %v", err)
	}
	payload, err := secrets.GetPayload(km, secretID, nil).Extract()
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to fetch master key payload: %v", err)
	}
	masterKey := strings.Split(string(payload), "\n")

	dec, err := decrypt(masterKey[0], masterKey[1], string(key.EncryptedDataKey()))
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to decrypt data key: %v", err)
	}

	return dec, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return key.SecretHref
}

// NewMasterKeyFromSecretHref takes a Barbican Secret Href and returns a new MasterKey for that
func NewMasterKeyFromSecretHref(secretHref string) *MasterKey {
	k := &MasterKey{}
	k.SecretHref = secretHref
	k.CreationDate = time.Now().UTC()
	return k
}

// MasterKeysFromSecretHref takes a comma separated list of Secret Hrefs and returns a slice of new MasterKeys for them
func MasterKeysFromSecretHref(secretHref string) []*MasterKey {
	var keys []*MasterKey
	if secretHref == "" {
		return keys
	}
	for _, s := range strings.Split(secretHref, ",") {
		keys = append(keys, NewMasterKeyFromSecretHref(s))
	}
	return keys
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["secret_href"] = key.SecretHref
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	return out
}

func newKeyManager() (*gophercloud.ServiceClient, error) {

	home, err := user.Current()
	if err != nil {
		log.WithError(err).Error("Failed to determine user home")
		return nil, err
	}

	// support loading config gophercloud style
	cfg, err := ini.LooseLoad(os.Getenv("GOPHERCLOUD_CONFIG"), fmt.Sprintf("%v/.cloud-config", home), "/etc/cloud-config")
	if err == nil {
		for _, k := range cfg.Section("Global").KeyStrings() {
			// we're iterating keys, safe not to check for err
			kv, _ := cfg.Section("Global").GetKey(k)
			envVar := fmt.Sprintf("OS_%s", strings.ToUpper(strings.ReplaceAll(k, "-", "_")))
			os.Setenv(envVar, kv.Value())
		}
	} else if !os.IsNotExist(err) {
		return nil, err
	}

	// fetch the client based on the env
	var provider *gophercloud.ProviderClient

	opts, err := clientconfig.AuthOptions(nil)
	if err != nil {
		log.WithError(err).Error("Failed to fetch auth options")
		return nil, err
	}
	trustID, present := os.LookupEnv("OS_TRUST_ID")
	if present {
		optsExt := trusts.AuthOptsExt{
			TrustID: trustID,
			AuthOptionsBuilder: &tokens.AuthOptions{
				IdentityEndpoint: opts.IdentityEndpoint,
				UserID:           os.Getenv("OS_USER_ID"),
				Password:         os.Getenv("OS_PASSWORD"),
			},
		}
		//FIXME(rochaporto): this is silly, but it's easy to get a client like this
		//we authenticate twice though
		provider, err = openstack.AuthenticatedClient(*opts)
		if err != nil {
			log.WithError(err).Error("failed to get openstack client")
			return nil, err
		}
		err = openstack.AuthenticateV3(provider, optsExt, gophercloud.EndpointOpts{})
		if err != nil {
			log.WithError(err).Error("failed to authenticate with v3")
			return nil, err
		}
	} else {
		provider, err = openstack.AuthenticatedClient(*opts)
		if err != nil {
			log.WithError(err).Error("failed to create authenticated client")
			return nil, err
		}
	}

	client, err := openstack.NewKeyManagerV1(provider,
		gophercloud.EndpointOpts{Region: "cern"})
	if err != nil {
		log.WithError(err).Error("failed to create key manager")
		return nil, err
	}
	return client, nil
}

func parseID(ref string) (string, error) {
	parts := strings.Split(ref, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("Could not parse %s", ref)
	}

	return parts[len(parts)-1], nil
}

func encrypt(b64key string, b64nonce string, payload []byte) ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(b64nonce)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	sealed := aesgcm.Seal(nil, nonce, payload, nil)
	result := make([]byte, base64.StdEncoding.EncodedLen(len(sealed)))
	base64.StdEncoding.Encode(result, sealed)
	return result, nil
}

func decrypt(b64key string, b64nonce string, b64payload string) ([]byte, error) {
	if b64payload == "" {
		return []byte{}, nil
	}
	key, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		return nil, err
	}
	nonce, err := base64.StdEncoding.DecodeString(b64nonce)
	if err != nil {
		return nil, err
	}
	payload, err := base64.StdEncoding.DecodeString(b64payload)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	plain, err := aesgcm.Open(nil, nonce, payload, nil)
	if err != nil {
		return nil, err
	}
	return plain, nil
}
