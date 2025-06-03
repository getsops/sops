package ocikms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	mockOciURL = "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xxxxx/ocid1.keyversion.xxxxxx"
)

func TestNewMasterKeyFromURL(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		expectErr bool
		expectKey MasterKey
	}{
		{
			name: "URL with slash",
			url:  "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xxxxx/ocid1.keyversion.xxxxxx/",
			expectKey: MasterKey{
				CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
				Id:             "ocid1.key.xxxxx",
				KeyVersionId:   "ocid1.keyversion.xxxxxx",
			},
		},
		{
			name: "URL",
			url:  "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xxxxx/ocid1.keyversion.xxxxxx",
			expectKey: MasterKey{
				CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
				Id:             "ocid1.key.xxxxx",
				KeyVersionId:   "ocid1.keyversion.xxxxxx",
			},
		},
		{
			name:      "wrong keyversion URL",
			url:       "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xxxxx/ocid1.key.xxxxx",
			expectErr: true,
		},
		{
			name:      "missing keyversion URL",
			url:       "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xxxxx",
			expectErr: true,
		},
		{
			name:      "missing keyid URL",
			url:       "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/",
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
			assert.Equal(t, tt.expectKey.CryptoEndpoint, key.CryptoEndpoint)
			assert.Equal(t, tt.expectKey.Id, key.Id)
			assert.Equal(t, tt.expectKey.KeyVersionId, key.KeyVersionId)
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
			name:           "URL",
			urls:           "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xx/ocid1.keyversion.xx",
			expectKeyCount: 1,
			expectKeys: []MasterKey{
				{
					CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
					Id:             "ocid1.key.xx",
					KeyVersionId:   "ocid1.keyversion.xx",
				},
			},
		},
		{
			name:           "multiple URLs",
			urls:           "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xx/ocid1.keyversion.xx,https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.yy/ocid1.keyversion.yy",
			expectKeyCount: 2,
			expectKeys: []MasterKey{
				{
					CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
					Id:             "ocid1.key.xx",
					KeyVersionId:   "ocid1.keyversion.xx",
				},
				{
					CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
					Id:             "ocid1.key.yy",
					KeyVersionId:   "ocid1.keyversion.yy",
				},
			},
		},
		{
			name:           "multiple URLs with leading and trailing spaces",
			urls:           " https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xx/ocid1.keyversion.xx,https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.yy/ocid1.keyversion.yy ",
			expectKeyCount: 2,
			expectKeys: []MasterKey{
				{
					CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
					Id:             "ocid1.key.xx",
					KeyVersionId:   "ocid1.keyversion.xx",
				},
				{
					CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
					Id:             "ocid1.key.yy",
					KeyVersionId:   "ocid1.keyversion.yy",
				},
			},
		},
		{
			name:      "multiple URLs, one malformed",
			urls:      "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xx/ocid1.test.xx,https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.yy/ocid1.keyversion.yy",
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
				assert.Equal(t, tt.expectKeys[idx].CryptoEndpoint, keys[idx].CryptoEndpoint)
				assert.Equal(t, tt.expectKeys[idx].Id, keys[idx].Id)
				assert.Equal(t, tt.expectKeys[idx].KeyVersionId, keys[idx].KeyVersionId)
				assert.NotNil(t, keys[idx].CreationDate)
			}
		})
	}
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
		key, err := NewMasterKeyFromURL(mockOciURL)
		assert.NoError(t, err)

		err = key.Encrypt([]byte("some data"))
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to encrypt sops data key with OCI KMS key")
	})

	t.Run("already encrypted", func(t *testing.T) {
		encryptedKey := "encrypted"
		key, err := NewMasterKeyFromURL(mockOciURL)
		assert.NoError(t, err)
		key.EncryptedKey = encryptedKey

		assert.NoError(t, key.EncryptIfNeeded([]byte("other data")))
		assert.Equal(t, encryptedKey, key.EncryptedKey)
	})
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := NewMasterKey("", "", "")
	assert.False(t, key.NeedsRotation())

	key.CreationDate = key.CreationDate.Add(-(ocikmsTTL + time.Second))
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKey("https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com", "ocid1.key.xx", "ocid1.keyversion.xx")
	assert.Equal(t, "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com/ocid1.key.xx/ocid1.keyversion.xx", key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := MasterKey{
		CreationDate:   time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		CryptoEndpoint: "https://test-crypto.kms.eu-frankfurt-1.oraclecloud.com",
		Id:             "ocid1.key.xx",
		KeyVersionId:   "ocid1.keyversion.xx",
		EncryptedKey:   "this is encrypted",
	}
	assert.Equal(t, map[string]interface{}{
		"crypto_endpoint": key.CryptoEndpoint,
		"id":              key.Id,
		"key_version":     key.KeyVersionId,
		"enc":             "this is encrypted",
		"created_at":      "2016-10-31T10:00:00Z",
	}, key.ToMap())
}
