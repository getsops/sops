package hcvault

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/v3/logging"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("VAULT_TRANSIT")
}

// MasterKey is a Vault Transit backend path used to encrypt and decrypt sops' data key.
type MasterKey struct {
	EncryptedKey string
	KeyName      string
	EnginePath  string
	VaultAddress string
	CreationDate time.Time
}

// NewMasterKeysFromURIs gets lots of keys from lots of URIs
func NewMasterKeysFromURIs(uris string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if uris == "" {
		return keys, nil
	}
	uriList := strings.Split(uris, ",")
	for _, uri := range uriList {
		if uri == "" {
			continue
		}
		key, err := NewMasterKeyFromURI(uri)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return keys, nil
}

// NewMasterKeyFromURI obtains the vaultAddress the transit backend path and the key name from the full URI of the key
func NewMasterKeyFromURI(uri string) (*MasterKey, error) {
	log.Debugln("Called NewMasterKeyFromURI with uri: ", uri)
	var key *MasterKey
	if uri == "" {
		return key, nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("missing scheme in vault URL (should be like this: https://vault.example.com:8200/v1/transit/keys/keyName), got: %v", uri)
	}
	enginePath, keyName, err := getBackendAndKeyFromPath(u.RequestURI())
	if err != nil {
		return nil, err
	}
	u.Path = ""
	return NewMasterKey(u.String(), enginePath, keyName), nil

}

func getBackendAndKeyFromPath(fullPath string) (enginePath, keyName string, err error) {
	// Running vault behind a reverse proxy with longer urls seems not to be supported
	// by the vault client api so we have a separate Error for that here.
	if re := regexp.MustCompile(`/[^/]+/v[\d]+/[^/]+/[^/]+/[^/]+`); re.Match([]byte(fullPath)) {
		return "", "", fmt.Errorf("running Vault with a prefixed url is not supported! (Format has to be like https://vault.example.com:8200/v1/transit/keys/keyName)")
	} else if re := regexp.MustCompile(`/v[\d]+/[^/]+/[^/]+/[^/]+`); re.Match([]byte(fullPath)) == false {
		return "", "", fmt.Errorf("vault path does not seem to be formatted correctly: (eg. https://vault.example.com:8200/v1/transit/keys/keyName)")
	}
	fullPath = strings.TrimPrefix(fullPath, "/")
	fullPath = strings.TrimSuffix(fullPath, "/")

	dirs := strings.Split(fullPath, "/")

	keyName = dirs[len(dirs)-1]
	enginePath = path.Join(dirs[1 : len(dirs)-2]...)
	err = nil
	return
}

// NewMasterKey creates a new MasterKey from a vault address, transit backend path and a key name and setting the creation date to the current date
func NewMasterKey(addess, enginePath, keyName string) *MasterKey {
	mk := &MasterKey{
		VaultAddress: addess,
		EnginePath:  enginePath,
		KeyName:      keyName,
		CreationDate: time.Now().UTC(),
	}
	log.Debugln("Created Vault Master Key: ", mk)
	return mk
}

// EncryptedDataKey returns the encrypted data key this master key holds
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

func vaultClient(address string) (*api.Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = address
	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("Cannot create Vault Client: %v", err)
	}
	if cli.Token() != "" {
		return cli, nil
	}
	homePath, err := homedir.Dir()
	if err != nil {
		panic(fmt.Sprintf("error getting user's home directory: %v", err))
	}
	tokenPath := filepath.Join(homePath, ".vault-token")
	f, err := os.Open(tokenPath)
	if os.IsNotExist(err) {
		return cli, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}
	cli.SetToken(strings.TrimSpace(buf.String()))
	return cli, nil
}

// Encrypt takes a sops data key, encrypts it with Vault Transit and stores the result in the EncryptedKey field
func (key *MasterKey) Encrypt(dataKey []byte) error {
	fullPath := path.Join(key.EnginePath, "encrypt", key.KeyName)
	cli, err := vaultClient(key.VaultAddress)
	if err != nil {
		return err
	}
	encoded := base64.StdEncoding.EncodeToString(dataKey)
	payload := make(map[string]interface{})
	payload["plaintext"] = encoded
	raw, err := cli.Logical().Write(fullPath, payload)
	if err != nil {
		log.WithField("Path", fullPath).Info("Encryption failed")
		return err
	}
	if raw == nil || raw.Data == nil {
		return fmt.Errorf("The transit backend %s is empty", fullPath)
	}
	encrypted, ok := raw.Data["ciphertext"]
	if !ok {
		return fmt.Errorf("there's not encrypted data")
	}
	encryptedKey, ok := encrypted.(string)
	if !ok {
		return fmt.Errorf("the ciphertext cannot be casted to string")
	}
	key.EncryptedKey = encryptedKey
	return nil
}

// EncryptIfNeeded encrypts the provided sops' data key and encrypts it if it hasn't been encrypted yet
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with Vault Transit and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	fullPath := path.Join(key.EnginePath, "decrypt", key.KeyName)
	cli, err := vaultClient(key.VaultAddress)
	if err != nil {
		return nil, err
	}
	payload := make(map[string]interface{})
	payload["ciphertext"] = key.EncryptedKey
	raw, err := cli.Logical().Write(fullPath, payload)
	if err != nil {
		log.WithField("Path", fullPath).Info("Encryption failed")
		return nil, err
	}
	if raw == nil || raw.Data == nil {
		return nil, fmt.Errorf("The transit backend %s is empty", fullPath)
	}
	decrypted, ok := raw.Data["plaintext"]
	if ok != true {
		return nil, fmt.Errorf("there's no decrypted data")
	}
	dataKey, ok := decrypted.(string)
	if ok != true {
		return nil, fmt.Errorf("the plaintest cannot be casted to string")
	}
	result, err := base64.StdEncoding.DecodeString(dataKey)
	if err != nil {
		return nil, fmt.Errorf("Couldn't decode base64 plaintext")
	}
	return result, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
// This is simply copied from GCPKMS
// TODO: handle key rotation on vault side
func (key *MasterKey) NeedsRotation() bool {
	//TODO: manage rewrapping https://www.vaultproject.io/api/secret/transit/index.html#rewrap-data
	return time.Since(key.CreationDate) > (time.Hour * 24 * 30 * 6)
}

// ToString converts the key to a string representation
func (key *MasterKey) ToString() string {
	return fmt.Sprintf("%s/v1/%s/keys/%s", key.VaultAddress, key.EnginePath, key.KeyName)
}

func (key *MasterKey) createVaultTransitAndKey() error {
	cli, err := vaultClient(key.VaultAddress)
	if err != nil {
		return err
	}
	if err != nil {
		return fmt.Errorf("Cannot create Vault Client: %v", err)
	}
	err = cli.Sys().Mount(key.EnginePath, &api.MountInput{
		Type:        "transit",
		Description: "backend transit used by SOPS",
	})
	if err != nil {
		return err
	}
	path := path.Join(key.EnginePath, "keys", key.KeyName)
	payload := make(map[string]interface{})
	payload["type"] = "rsa-4096"
	_, err = cli.Logical().Write(path, payload)
	if err != nil {
		return err
	}
	_, err = cli.Logical().Read(path)
	if err != nil {
		return err
	}
	return nil
}

// ToMap converts the MasterKey to a map for serialization purposes
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["vault_address"] = key.VaultAddress
	out["key_name"] = key.KeyName
	out["engine_path"] = key.EnginePath
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	return out
}
