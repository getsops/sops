/*
Package azkv contains an implementation of the github.com/getsops/sops/v3/keys.MasterKey
interface that encrypts and decrypts the data key using Azure Key Vault with the
Azure Key Vault Keys client module for Go.
*/
package azkv // import "github.com/getsops/sops/v3/azkv"

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify an Azure Key Vault MasterKey.
	KeyTypeIdentifier = "azure_kv"
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
	// clientOptions contains the azkeys.ClientOptions used by the Azure client.
	clientOptions *azkeys.ClientOptions
}

// newMasterKey creates a new MasterKey from a URL, key name and version,
// setting the creation date to the current date.
func newMasterKey(vaultURL string, keyName string, keyVersion string) *MasterKey {
	return &MasterKey{
		VaultURL:     vaultURL,
		Name:         keyName,
		Version:      keyVersion,
		CreationDate: time.Now().UTC(),
	}
}

// NewMasterKey creates a new MasterKey from a URL, key name and (mandatory) version,
// setting the creation date to the current date.
func NewMasterKey(vaultURL string, keyName string, keyVersion string) *MasterKey {
	return newMasterKey(vaultURL, keyName, keyVersion)
}

// NewMasterKey creates a new MasterKey from a URL, key name and (optional) version,
// setting the creation date to the current date.
func NewMasterKeyWithOptionalVersion(vaultURL string, keyName string, keyVersion string) (*MasterKey, error) {
	key := newMasterKey(vaultURL, keyName, keyVersion)
	if err := key.ensureKeyHasVersion(context.Background()); err != nil {
		return nil, err
	}
	return key, nil
}

// NewMasterKeyFromURL takes an Azure Key Vault key URL, and returns a new
// MasterKey. The URL format is {vaultUrl}/keys/{keyName}/{keyVersion}.
func NewMasterKeyFromURL(url string) (*MasterKey, error) {
	url = strings.TrimSpace(url)
	re := regexp.MustCompile("^(https://[^/]+)/keys/([^/]+)(/[^/]*)?$")
	parts := re.FindStringSubmatch(url)
	if len(parts) < 3 {
		return nil, fmt.Errorf("could not parse %q into a valid Azure Key Vault MasterKey %v", url, parts)
	}
	// Blank key versions are supported in Azure Key Vault, as they default to the latest
	// version of the key. We need to put the actual version in the sops metadata block though
	var key *MasterKey
	if len(parts[3]) > 1 {
		key = newMasterKey(parts[1], parts[2], parts[3][1:])
	} else {
		key = newMasterKey(parts[1], parts[2], "")
	}
	err := key.ensureKeyHasVersion(context.Background())
	return key, err
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

// ClientOptions is a wrapper around azkeys.ClientOptions to allow
// configuration of the Azure Key Vault client.
type ClientOptions struct {
	o *azkeys.ClientOptions
}

// NewClientOptions creates a new ClientOptions with the provided
// azkeys.ClientOptions.
func NewClientOptions(o *azkeys.ClientOptions) *ClientOptions {
	return &ClientOptions{o: o}
}

// ApplyToMasterKey configures the ClientOptions on the provided key.
func (c ClientOptions) ApplyToMasterKey(key *MasterKey) {
	key.clientOptions = c.o
}

// Encrypt takes a SOPS data key, encrypts it with Azure Key Vault, and stores
// the result in the EncryptedKey field.
//
// Consider using EncryptContext instead.
func (key *MasterKey) Encrypt(dataKey []byte) error {
	return key.EncryptContext(context.Background(), dataKey)
}

func (key *MasterKey) ensureKeyHasVersion(ctx context.Context) error {
	if (key.Version != "") {
		// Nothing to do
		return nil
	}

	token, err := key.getTokenCredential()

	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to get Azure token credential to retrieve key version: %w", err)
	}

	c, err := azkeys.NewClient(key.VaultURL, token, key.clientOptions)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to construct Azure Key Vault client to retrieve key version: %w", err)
	}

	kdetail, err := c.GetKey(ctx, key.Name, key.Version, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to fetch Azure Key to retrieve key version: %w", err)
	}
	key.Version = kdetail.Key.KID.Version()

	log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Version fetch succeeded")
	return nil
}

// EncryptContext takes a SOPS data key, encrypts it with Azure Key Vault, and stores
// the result in the EncryptedKey field.
func (key *MasterKey) EncryptContext(ctx context.Context, dataKey []byte) error {
	token, err := key.getTokenCredential()
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to get Azure token credential to encrypt data: %w", err)
	}

	c, err := azkeys.NewClient(key.VaultURL, token, key.clientOptions)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to construct Azure Key Vault client to encrypt data: %w", err)
	}

	resp, err := c.Encrypt(ctx, key.Name, key.Version, azkeys.KeyOperationParameters{
		Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP256),
		Value:     dataKey,
	}, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to encrypt sops data key with Azure Key Vault key '%s': %w", key.ToString(), err)
	}

	encodedEncryptedKey := base64.RawURLEncoding.EncodeToString(resp.KeyOperationResult.Result)
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
//
// Consider using DecryptContext instead.
func (key *MasterKey) Decrypt() ([]byte, error) {
	return key.DecryptContext(context.Background())
}

// DecryptContext decrypts the EncryptedKey field with Azure Key Vault and returns
// the result.
func (key *MasterKey) DecryptContext(ctx context.Context) ([]byte, error) {
	token, err := key.getTokenCredential()
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption failed")
		return nil, fmt.Errorf("failed to get Azure token credential to decrypt: %w", err)
	}

	rawEncryptedKey, err := base64.RawURLEncoding.DecodeString(key.EncryptedKey)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption failed")
		return nil, fmt.Errorf("failed to base64 decode Azure Key Vault encrypted key: %w", err)
	}

	c, err := azkeys.NewClient(key.VaultURL, token, key.clientOptions)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption failed")
		return nil, fmt.Errorf("failed to construct Azure Key Vault client to decrypt data: %w", err)
	}

	resp, err := c.Decrypt(ctx, key.Name, key.Version, azkeys.KeyOperationParameters{
		Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP256),
		Value:     rawEncryptedKey,
	}, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption failed")
		return nil, fmt.Errorf("failed to decrypt sops data key with Azure Key Vault key '%s': %w", key.ToString(), err)
	}
	log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption succeeded")
	return resp.KeyOperationResult.Result, nil
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

// TypeToIdentifier returns the string identifier for the MasterKey type.
func (key *MasterKey) TypeToIdentifier() string {
	return KeyTypeIdentifier
}

// getTokenCredential returns the tokenCredential of the MasterKey, or
// azidentity.NewDefaultAzureCredential.
func (key *MasterKey) getTokenCredential() (azcore.TokenCredential, error) {
	if key.tokenCredential == nil {
		return azidentity.NewDefaultAzureCredential(nil)
	}
	return key.tokenCredential, nil
}
