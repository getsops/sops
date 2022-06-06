/*
Package azkv contains an implementation of the go.mozilla.org/sops/v3/keys.MasterKey
interface that encrypts and decrypts the data key using Azure Key Vault with the
Azure Key Vault Keys client module for Go.
*/
package azkv //import "go.mozilla.org/sops/v3/azkv"

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto"
	"github.com/sirupsen/logrus"
	"go.mozilla.org/sops/v3/logging"
)

var (
	// log is the global logger for any Azure Key Vault MasterKey.
	log *logrus.Logger
	// azkvTTL is the duration after which a MasterKey requires rotation.
	azkvTTL = time.Hour * 24 * 30 * 6
)

func init() {
	log = logging.NewLogger("AZKV")
}

// MasterKey is an Azure Key Vault Key used to Encrypt and Decrypt SOPS'
// data key.
type MasterKey struct {
	// VaultURL of the Azure Key Vault. For example:
	// "https://myvault.vault.azure.net/".
	VaultURL string
	// Name of the Azure Key Vault key in the VaultURL.
	Name string
	// Version of the Azure Key Vault key. Can be empty.
	Version string
	// EncryptedKey contains the SOPS data key encrypted with the Azure Key
	// Vault key.
	EncryptedKey string
	// CreationDate of the MasterKey, used to determine if the EncryptedKey
	// needs rotation.
	CreationDate time.Time

	// tokenCredential contains the azcore.TokenCredential used by the Azure
	// client. It can be injected by a (local) keyservice.KeyServiceServer
	// using TokenCredential.ApplyToMasterKey.
	// If nil, azidentity.NewDefaultAzureCredential is used.
	tokenCredential azcore.TokenCredential
}

// NewMasterKey creates a new MasterKey from a URL, key name and version,
// setting the creation date to the current date.
func NewMasterKey(vaultURL string, keyName string, keyVersion string) *MasterKey {
	return &MasterKey{
		VaultURL:     vaultURL,
		Name:         keyName,
		Version:      keyVersion,
		CreationDate: time.Now().UTC(),
	}
}

// NewMasterKeyFromURL takes an Azure Key Vault key URL, and returns a new
// MasterKey. The URL format is {vaultUrl}/keys/{keyName}/{keyVersion}.
func NewMasterKeyFromURL(url string) (*MasterKey, error) {
	re := regexp.MustCompile("^(https://[^/]+)/keys/([^/]+)/([^/]+)$")
	parts := re.FindStringSubmatch(url)
	if parts == nil || len(parts) < 3 {
		return nil, fmt.Errorf("could not parse %q into a valid Azure Key Vault MasterKey", url)
	}
	return NewMasterKey(parts[1], parts[2], parts[3]), nil
}

// MasterKeysFromURLs takes a comma separated list of Azure Key Vault URLs,
// and returns a slice of new MasterKeys.
func MasterKeysFromURLs(urls string) ([]*MasterKey, error) {
	var keys []*MasterKey
	if urls == "" {
		return keys, nil
	}
	for _, s := range strings.Split(urls, ",") {
		k, err := NewMasterKeyFromURL(s)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// TokenCredential is an azcore.TokenCredential used for authenticating towards Azure Key
// Vault.
type TokenCredential struct {
	token azcore.TokenCredential
}

// NewTokenCredential creates a new TokenCredential with the provided azcore.TokenCredential.
func NewTokenCredential(token azcore.TokenCredential) *TokenCredential {
	return &TokenCredential{token: token}
}

// ApplyToMasterKey configures the TokenCredential on the provided key.
func (t TokenCredential) ApplyToMasterKey(key *MasterKey) {
	key.tokenCredential = t.token
}

// Encrypt takes a SOPS data key, encrypts it with Azure Key Vault, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	token, err := key.getTokenCredential()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Encryption failed")
		return fmt.Errorf("failed to get Azure token credential to encrypt data: %w", err)
	}
	c, err := crypto.NewClient(key.ToString(), token, nil)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Encryption failed")
		return fmt.Errorf("failed to construct Azure Key Vault crypto client to encrypt data: %w", err)
	}

	resp, err := c.Encrypt(context.Background(), crypto.EncryptionAlgRSAOAEP256, dataKey, nil)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with Azure Key Vault key '%s': %w", key.ToString(), err)
	}
	encodedEncryptedKey := base64.RawURLEncoding.EncodeToString(resp.Ciphertext)
	key.SetEncryptedDataKey([]byte(encodedEncryptedKey))
	log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption succeeded")
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

// EncryptIfNeeded encrypts the provided SOPS data key, if it has not been
// encrypted yet.
func (key *MasterKey) EncryptIfNeeded(dataKey []byte) error {
	if key.EncryptedKey == "" {
		return key.Encrypt(dataKey)
	}
	return nil
}

// Decrypt decrypts the EncryptedKey field with Azure Key Vault and returns
// the result.
func (key *MasterKey) Decrypt() ([]byte, error) {
	token, err := key.getTokenCredential()
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Decryption failed")
		return nil, fmt.Errorf("failed to get Azure token credential to decrypt: %w", err)
	}
	c, err := crypto.NewClient(key.ToString(), token, nil)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Decryption failed")
		return nil, fmt.Errorf("failed to construct Azure Key Vault crypto client to decrypt data: %w", err)
	}

	rawEncryptedKey, err := base64.RawURLEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Decryption failed")
		return nil, fmt.Errorf("failed to base64 decode Azure Key Vault encrypted key: %w", err)
	}
	resp, err := c.Decrypt(context.Background(), crypto.EncryptionAlgRSAOAEP256, rawEncryptedKey, nil)
	if err != nil {
		log.WithError(err).WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Error("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with Azure Key Vault key '%s': %w", key.ToString(), err)
	}
	log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption succeeded")
	return resp.Plaintext, nil
}

// NeedsRotation returns whether the data key needs to be rotated or not.
func (key *MasterKey) NeedsRotation() bool {
	return time.Since(key.CreationDate) > (azkvTTL)
}

// ToString converts the key to a string representation.
func (key *MasterKey) ToString() string {
	return fmt.Sprintf("%s/keys/%s/%s", key.VaultURL, key.Name, key.Version)
}

// ToMap converts the MasterKey to a map for serialization purposes.
func (key MasterKey) ToMap() map[string]interface{} {
	out := make(map[string]interface{})
	out["vaultUrl"] = key.VaultURL
	out["key"] = key.Name
	out["version"] = key.Version
	out["created_at"] = key.CreationDate.UTC().Format(time.RFC3339)
	out["enc"] = key.EncryptedKey
	return out
}

// getTokenCredential returns the tokenCredential of the MasterKey, or
// azidentity.NewDefaultAzureCredential.
func (key *MasterKey) getTokenCredential() (azcore.TokenCredential, error) {
	if key.tokenCredential == nil {
		return azidentity.NewDefaultAzureCredential(nil)
	}
	return key.tokenCredential, nil
}
