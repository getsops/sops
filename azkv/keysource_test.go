package azkv

import (
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/stretchr/testify/assert"
)

const (
	mockAzureURL = "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90"
)

func TestNewMasterKeyFromURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
		expectKey MasterKey
	}{
		{
			name: "URL",
			url:  "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90",
			expectKey: MasterKey{
				VaultURL: "https://test.vault.azure.net",
				Name:     "test-key",
				Version:  "a2a690a4fcc04166b739da342a912c90",
			},
		},
		{
			name:      "malformed URL",
			url:       "https://test.vault.azure.net/no-keys-here/test-key/a2a690a4fcc04166b739da342a912c90",
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewMasterKeyFromURL(tt.url)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, key)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectKey.VaultURL, key.VaultURL)
			assert.Equal(t, tt.expectKey.Name, key.Name)
			assert.Equal(t, tt.expectKey.Version, key.Version)
			assert.NotNil(t, key.CreationDate)
		})
	}
}

func TestMasterKeysFromURLs(t *testing.T) {
	tests := []struct {
		name           string
		urls           string
		expectErr      bool
		expectKeyCount int
		expectKeys     []MasterKey
	}{
		{
			name:           "single URL",
			urls:           "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90",
			expectKeyCount: 1,
			expectKeys: []MasterKey{
				{
					VaultURL: "https://test.vault.azure.net",
					Name:     "test-key",
					Version:  "a2a690a4fcc04166b739da342a912c90",
				},
			},
		},
		{
			name:           "multiple URLs",
			urls:           "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90,https://test2.vault.azure.net/keys/another-test-key/cf0021e8b743453bae758e7fbf71b60e",
			expectKeyCount: 2,
			expectKeys: []MasterKey{
				{
					VaultURL: "https://test.vault.azure.net",
					Name:     "test-key",
					Version:  "a2a690a4fcc04166b739da342a912c90",
				},
				{
					VaultURL: "https://test2.vault.azure.net",
					Name:     "another-test-key",
					Version:  "cf0021e8b743453bae758e7fbf71b60e",
				},
			},
		},
		{
			name:      "multiple URLs, one malformed",
			urls:      "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90,https://test.vault.azure.net/no-keys-here/test-key/a2a690a4fcc04166b739da342a912c90",
			expectErr: true,
		},
		{
			name:           "empty",
			urls:           "",
			expectErr:      false,
			expectKeyCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := MasterKeysFromURLs(tt.urls)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, keys)
				return
			}
			assert.NoError(t, err)
			assert.Len(t, keys, tt.expectKeyCount)
			for idx := range keys {
				assert.Equal(t, tt.expectKeys[idx].VaultURL, keys[idx].VaultURL)
				assert.Equal(t, tt.expectKeys[idx].Name, keys[idx].Name)
				assert.Equal(t, tt.expectKeys[idx].Version, keys[idx].Version)
				assert.NotNil(t, keys[idx].CreationDate)
			}
		})
	}
}

func TestTokenCredential_ApplyToMasterKey(t *testing.T) {
	credential, err := azidentity.NewUsernamePasswordCredential("tenant", "client", "username", "password", nil)
	assert.NoError(t, err)
	token := NewTokenCredential(credential)

	key := &MasterKey{}
	token.ApplyToMasterKey(key)
	assert.Equal(t, credential, key.tokenCredential)
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: "some key"}
	assert.EqualValues(t, key.EncryptedKey, key.EncryptedDataKey())
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	encryptedKey := []byte("encrypted")
	key := &MasterKey{}
	key.SetEncryptedDataKey(encryptedKey)
	assert.EqualValues(t, encryptedKey, key.EncryptedKey)
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	t.Run("not encrypted", func(t *testing.T) {
		key, err := NewMasterKeyFromURL(mockAzureURL)
		assert.NoError(t, err)

		err = key.Encrypt([]byte("some data"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to encrypt sops data key with Azure Key Vault key")
	})

	t.Run("already encrypted", func(t *testing.T) {
		encryptedKey := "encrypted"
		key, err := NewMasterKeyFromURL(mockAzureURL)
		assert.NoError(t, err)
		key.EncryptedKey = encryptedKey

		assert.NoError(t, key.EncryptIfNeeded([]byte("other data")))
		assert.Equal(t, encryptedKey, key.EncryptedKey)
	})
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := NewMasterKey("", "", "")
	assert.False(t, key.NeedsRotation())

	key.CreationDate = key.CreationDate.Add(-(azkvTTL + time.Second))
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKey("https://test.vault.azure.net", "key-name", "key-version")
	assert.Equal(t, "https://test.vault.azure.net/keys/key-name/key-version", key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		VaultURL:     "https://test.vault.azure.net",
		Name:         "test-key",
		Version:      "1",
		EncryptedKey: "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"vaultUrl":   key.VaultURL,
		"key":        key.Name,
		"version":    key.Version,
		"enc":        "this is encrypted",
		"created_at": "2016-10-31T10:00:00Z",
	}, key.ToMap())
}

func TestMasterKey_getTokenCredential(t *testing.T) {
	t.Run("with TokenCredential", func(t *testing.T) {
		credential, err := azidentity.NewUsernamePasswordCredential("tenant", "client", "username", "password", nil)
		assert.NoError(t, err)
		token := NewTokenCredential(credential)

		key := &MasterKey{}
		token.ApplyToMasterKey(key)

		got, err := key.getTokenCredential()
		assert.NoError(t, err)
		assert.Equal(t, credential, got)
	})

	t.Run("default", func(t *testing.T) {
		key := &MasterKey{}
		got, err := key.getTokenCredential()
		assert.NoError(t, err)
		assert.IsType(t, &azidentity.DefaultAzureCredential{}, got)
	})
}
