package hcvault

import (
	"bytes"
	"encoding/base64"
	"errors"
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
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify a Vault MasterKey.
	KeyTypeIdentifier = "hc_vault"
)

func init() {
	log = logging.NewLogger("VAULT_TRANSIT")
}

var (
	// log is the global logger for any Vault Transit MasterKey.
	log *logrus.Logger
	// vaultTTL is the duration after which a MasterKey requires rotation.
	vaultTTL = time.Hour * 24 * 30 * 6
	// defaultTokenFile is the name of the file in the user's home directory
	// where a Vault token is expected to be stored.
	defaultTokenFile = ".vault-token"
)

// Token used for authenticating towards a Vault server.
type Token string

// ApplyToMasterKey configures the token on the provided key.
func (t Token) ApplyToMasterKey(key *MasterKey) {
	key.token = string(t)
}

// MasterKey is a Vault Transit backend path used to Encrypt and Decrypt
// SOPS' data key.
type MasterKey struct {
	// VaultAddress is the address of the Vault server.
	VaultAddress string
	// EnginePath is the path to the Vault Transit Secret engine relative
	// to the VaultAddress.
	EnginePath string
	// KeyName is the name of the key in the Vault Transit engine.
	KeyName string
	// EncryptedKey contains the SOPS data key encrypted with the Vault Transit
	// key.
	EncryptedKey string
	// CreationDate of the MasterKey, used to determine if the EncryptedKey
	// needs rotation.
	CreationDate time.Time

	// token is the token used for authenticating against the VaultAddress
	// server. It can be injected by a (local) keyservice.KeyServiceServer
	// Token.ApplyToMasterKey. If empty, the default client configuration
	// is used, before falling back to the token stored in defaultTokenFile.
	token string
}

// NewMasterKeysFromURIs creates a list of MasterKeys from a list of Vault
// URIs.
func NewMasterKeysFromURIs(uris string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if uris == "" {
		return keys, nil
	}
	for _, uri := range strings.Split(uris, ",") {
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

// NewMasterKeyFromURI obtains the Vault address, Transit backend path and the
// key name from the full URI of the key.
func NewMasterKeyFromURI(uri string) (*MasterKey, error) {
	var key *MasterKey
	if uri == "" {
		return key, nil
	}
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("missing scheme in Vault URL (should be like this: +"+
			"https://vault.example.com:8200/v1/transit/keys/keyName), got: %v", uri)
	}
	enginePath, keyName, err := engineAndKeyFromPath(u.RequestURI())
	if err != nil {
		return nil, err
	}
	u.Path = ""
	return NewMasterKey(u.String(), enginePath, keyName), nil

}

// NewMasterKey creates a new MasterKey from a Vault address, Transit backend
// path and a key name.
func NewMasterKey(address, enginePath, keyName string) *MasterKey {
	key := &MasterKey{
		VaultAddress: address,
		EnginePath:   enginePath,
		KeyName:      keyName,
		CreationDate: time.Now().UTC(),
	}
	return key
}

// Encrypt takes a SOPS data key, encrypts it with Vault Transit, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	fullPath := key.encryptPath()

	client, err := vaultClient(key.VaultAddress, key.token)
	if err != nil {
		log.WithField("Path", fullPath).Info("Encryption failed")
		return err
	}

	secret, err := client.Logical().Write(fullPath, encryptPayload(dataKey))
	if err != nil {
		log.WithField("Path", fullPath).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key to Vault transit backend '%s': %w", fullPath, err)
	}
	encryptedKey, err := encryptedKeyFromSecret(secret)
	if err != nil {
		log.WithField("Path", fullPath).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key to Vault transit backend '%s': %w", fullPath, err)
	}

	key.EncryptedKey = encryptedKey
	log.WithField("Path", fullPath).Info("Encryption successful")
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

// EncryptedDataKey returns the encrypted data key this master key holds.
func (key *MasterKey) EncryptedDataKey() []byte {
	return []byte(key.EncryptedKey)
}

// SetEncryptedDataKey sets the encrypted data key for this master key.
func (key *MasterKey) SetEncryptedDataKey(enc []byte) {
	key.EncryptedKey = string(enc)
}

// Decrypt decrypts the EncryptedKey field with Vault Transit and returns the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	fullPath := key.decryptPath()

	client, err := vaultClient(key.VaultAddress, key.token)
	if err != nil {
		log.WithField("Path", fullPath).Info("Decryption failed")
		return nil, err
	}

	secret, err := client.Logical().Write(fullPath, decryptPayload(key.EncryptedKey))
	if err != nil {
		log.WithField("Path", fullPath).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key from Vault transit backend '%s': %w", fullPath, err)
	}
	dataKey, err := dataKeyFromSecret(secret)
	if err != nil {
		log.WithField("Path", fullPath).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key from Vault transit backend '%s': %w", fullPath, err)
	}

	log.WithField("Path", fullPath).Info("Decryption successful")
	return dataKey, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	// TODO: manage rewrapping https://www.vaultproject.io/api/secret/transit/index.html#rewrap-data
	return time.Since(key.CreationDate) > (vaultTTL)
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return fmt.Sprintf("%s/v1/%s/keys/%s", key.VaultAddress, key.EnginePath, key.KeyName)
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["vault_address"] = key.VaultAddress
	out["key_name"] = key.KeyName
	out["engine_path"] = key.EnginePath
	out["enc"] = key.EncryptedKey
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	return out
}

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// encryptPath returns the path for Encrypt requests.
func (key *MasterKey) encryptPath() string {
	return path.Join(key.EnginePath, "encrypt", key.KeyName)
}

// decryptPath returns the path for Decrypt requests.
func (key *MasterKey) decryptPath() string {
	return path.Join(key.EnginePath, "decrypt", key.KeyName)
}

// encryptPayload returns the payload for an encrypt request of the dataKey.
func encryptPayload(dataKey []byte) map[string]interface{} {
	encoded := base64.StdEncoding.EncodeToString(dataKey)
	return map[string]interface{}{
		"plaintext": encoded,
	}
}

// encryptedKeyFromSecret attempts to extract the encrypted key from the data
// of the provided secret.
func encryptedKeyFromSecret(secret *api.Secret) (string, error) {
	if secret == nil || secret.Data == nil {
		return "", fmt.Errorf("transit backend is empty")
	}
	encrypted, ok := secret.Data["ciphertext"]
	if !ok {
		return "", fmt.Errorf("no encrypted data")
	}
	encryptedKey, ok := encrypted.(string)
	if !ok {
		return "", fmt.Errorf("encrypted ciphertext cannot be cast to string")
	}
	return encryptedKey, nil
}

// decryptPayload returns the payload for a decrypt request of the
// encryptedKey.
func decryptPayload(encryptedKey string) map[string]interface{} {
	return map[string]interface{}{
		"ciphertext": encryptedKey,
	}
}

// dataKeyFromSecret attempts to extract the data key from the data of the
// provided secret.
func dataKeyFromSecret(secret *api.Secret) ([]byte, error) {
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("transit backend is empty")
	}
	decrypted, ok := secret.Data["plaintext"]
	if !ok {
		return nil, fmt.Errorf("no decrypted data")
	}
	plaintext, ok := decrypted.(string)
	if !ok {
		return nil, fmt.Errorf("decrypted plaintext data cannot be cast to string")
	}
	dataKey, err := base64.StdEncoding.DecodeString(plaintext)
	if err != nil {
		return nil, fmt.Errorf("cannot decode base64 plaintext into data key bytes")
	}
	return dataKey, nil
}

// vaultClient returns a new Vault client, configured with the given address
// and token.
func vaultClient(address, token string) (*api.Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create Vault client: %w", err)
	}

	if token != "" {
		client.SetToken(token)
	}
	// Provided token takes precedence over the user's token.
	if client.Token() == "" {
		if token, err = userVaultToken(); err != nil {
			return nil, fmt.Errorf("cannot get Vault token: %w", err)
		}
		if token != "" {
			client.SetToken(token)
		}
	}

	return client, nil
}

// userVaultsToken returns the token from `$HOME/.vault-token` if the file
// exists. It returns an error if the file exists but cannot be read from.
// If the file does not exist, it returns an empty string.
func userVaultToken() (string, error) {
	homePath, err := homedir.Dir()
	if err != nil {
		return "", fmt.Errorf("error getting user's home directory: %w", err)
	}
	tokenPath := filepath.Join(homePath, defaultTokenFile)

	f, err := os.Open(tokenPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

// engineAndKeyFromPath returns the engine path and key name from the full
// path, or an error.
func engineAndKeyFromPath(fullPath string) (enginePath, keyName string, err error) {
	// Running vault behind a reverse proxy with longer URLs seems not to be
	// supported by the Vault client API. Check for this here.
	// TODO(hidde): this may no longer be necessary with newer Vault versions,
	//  but needs to be confirmed.
	if re := regexp.MustCompile(`/[^/]+/v[\d]+/[^/]+/[^/]+/[^/]+`); re.Match([]byte(fullPath)) {
		err = fmt.Errorf("running Vault with a prefixed URL is not supported! (Format has to be like " +
			"https://vault.example.com:8200/v1/transit/keys/keyName)")
		return
	} else if re := regexp.MustCompile(`/v[\d]+/[^/]+/[^/]+/[^/]+`); !re.Match([]byte(fullPath)) {
		err = fmt.Errorf("vault path does not seem to be formatted correctly: (eg. " +
			"https://vault.example.com:8200/v1/transit/keys/keyName)")
		return
	}

	fullPath = strings.Trim(fullPath, "/")
	dirs := strings.Split(fullPath, "/")

	keyName = dirs[len(dirs)-1]
	enginePath = path.Join(dirs[1 : len(dirs)-2]...)
	return
}
