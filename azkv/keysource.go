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
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"encoding/json"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	azidentitycache "github.com/Azure/azure-sdk-for-go/sdk/azidentity/cache"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/sirupsen/logrus"

	"github.com/getsops/sops/v3/logging"
)

const (
	// KeyTypeIdentifier is the string used to identify an Azure Key Vault MasterKey.
	KeyTypeIdentifier = "azure_kv"

	SopsAzureAuthMethodEnv = "SOPS_AZURE_AUTH_METHOD"

	cachedBrowserAuthRecordFileName    = "azure-auth-record-browser.json"
	cachedDeviceCodeAuthRecordFileName = "azure-auth-record-device-code.json"
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
	url = strings.TrimSpace(url)
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
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to get Azure token credential to encrypt data: %w", err)
	}

	c, err := azkeys.NewClient(key.VaultURL, token, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Encryption failed")
		return fmt.Errorf("failed to construct Azure Key Vault client to encrypt data: %w", err)
	}

	resp, err := c.Encrypt(context.Background(), key.Name, key.Version, azkeys.KeyOperationParameters{
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
func (key *MasterKey) Decrypt() ([]byte, error) {
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

	c, err := azkeys.NewClient(key.VaultURL, token, nil)
	if err != nil {
		log.WithFields(logrus.Fields{"key": key.Name, "version": key.Version}).Info("Decryption failed")
		return nil, fmt.Errorf("failed to construct Azure Key Vault client to decrypt data: %w", err)
	}

	resp, err := c.Decrypt(context.Background(), key.Name, key.Version, azkeys.KeyOperationParameters{
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

		authMethod := strings.ToLower(os.Getenv(SopsAzureAuthMethodEnv))
		switch authMethod {
		case "cached-browser":
			return cachedInteractiveBrowserCredentials()
		case "cached-device-code":
			return cachedDeviceCodeCredentials()
		case "azure-cli":
			return azidentity.NewAzureCLICredential(nil)
		case "msi":
			return azidentity.NewManagedIdentityCredential(nil)
		// If "DEFAULT" or not explicitly specified then use the default authentication chain.
		case "", "default":
			return azidentity.NewDefaultAzureCredential(nil)
		default:
			return nil, fmt.Errorf("Value `%s` is unsupported for environment variable `%s`, to resolve this either leave it unset or use one of `default`/`msi`/`azure-cli`/`cached-browser`/`cached-device-code`", authMethod, SopsAzureAuthMethodEnv)
		}
	}
	return key.tokenCredential, nil
}

func sopsCacheDir() (string, error) {
	userCacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(userCacheDir, "/sops")

	if err = os.MkdirAll(cacheDir, 0o700); err != nil {
		return "", err
	}

	return cacheDir, nil
}

type CachableTokenCredential interface {
	Authenticate(ctx context.Context, opts *policy.TokenRequestOptions) (azidentity.AuthenticationRecord, error)
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

func cacheStoreRecord(cachePath string, record azidentity.AuthenticationRecord) error {
	b, err := json.Marshal(record)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, b, 0600)
}

func cacheLoadRecord(cachePath string) (azidentity.AuthenticationRecord, error) {
	var record azidentity.AuthenticationRecord

	b, err := os.ReadFile(cachePath)
	if err != nil {
		return record, err
	}

	err = json.Unmarshal(b, &record)
	if err != nil {
		return record, err
	}

	return record, nil
}

func cacheTokenCredential(cachePath string, tokenCredentialFn func(cache azidentity.Cache, record azidentity.AuthenticationRecord) (CachableTokenCredential, error)) (azcore.TokenCredential, error) {
	cache, err := azidentitycache.New(nil)
	// Errors if persistent caching is not supported by the current runtime
	if err != nil {
		return nil, err
	}

	cachedRecord, cacheLoadErr := cacheLoadRecord(cachePath)

	credential, err := tokenCredentialFn(cache, cachedRecord)
	if err != nil {
		return nil, err
	}

	// If loading the authenticationRecord from the cachePath failed for any reason (validation, file doesn't exist, not encoded using json, etc.)
	if cacheLoadErr != nil {
		record, err := credential.Authenticate(context.Background(), nil)
		if err != nil {
			return nil, err
		}

		if err = cacheStoreRecord(cachePath, record); err != nil {
			return nil, err
		}
	}

	return credential, nil
}

func cachedInteractiveBrowserCredentials() (azcore.TokenCredential, error) {
	cacheDir, err := sopsCacheDir()
	if err != nil {
		return nil, err
	}
	return cacheTokenCredential(
		filepath.Join(cacheDir, cachedBrowserAuthRecordFileName),
		func(cache azidentity.Cache, record azidentity.AuthenticationRecord) (CachableTokenCredential, error) {
			return azidentity.NewInteractiveBrowserCredential(&azidentity.InteractiveBrowserCredentialOptions{
				AuthenticationRecord: record,
				Cache:                cache,
			})
		},
	)
}

func cachedDeviceCodeCredentials() (azcore.TokenCredential, error) {
	cacheDir, err := sopsCacheDir()
	if err != nil {
		return nil, err
	}

	return cacheTokenCredential(
		filepath.Join(cacheDir, cachedDeviceCodeAuthRecordFileName),
		func(cache azidentity.Cache, record azidentity.AuthenticationRecord) (CachableTokenCredential, error) {
			return azidentity.NewDeviceCodeCredential(&azidentity.DeviceCodeCredentialOptions{
				AuthenticationRecord: record,
				Cache:                cache,
			})
		},
	)
}
