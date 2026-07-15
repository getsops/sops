package hckms

import (
	"testing"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/stretchr/testify/assert"
)

const (
	testKeyID1 = "tr-west-1:12345678-1234-1234-1234-123456789abc"
	testKeyID2 = "tr-west-1:87654321-4321-4321-4321-cba987654321"
)

func TestNewMasterKey(t *testing.T) {
	tests := []struct {
		name      string
		keyID     string
		expectErr bool
		expectKey MasterKey
	}{
		{
			name:  "valid key ID",
			keyID: testKeyID1,
			expectKey: MasterKey{
				KeyID:   testKeyID1,
				Region:  "tr-west-1",
				KeyUUID: "12345678-1234-1234-1234-123456789abc",
			},
		},
		{
			name:      "invalid format - no colon",
			keyID:     "tr-west-1-12345678-1234-1234-1234-123456789abc",
			expectErr: true,
		},
		{
			name:      "invalid format - empty region",
			keyID:     ":12345678-1234-1234-1234-123456789abc",
			expectErr: true,
		},
		{
			name:      "invalid format - empty UUID",
			keyID:     "tr-west-1:",
			expectErr: true,
		},
		{
			name:      "invalid format - empty string",
			keyID:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewMasterKey(tt.keyID)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, key)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectKey.KeyID, key.KeyID)
			assert.Equal(t, tt.expectKey.Region, key.Region)
			assert.Equal(t, tt.expectKey.KeyUUID, key.KeyUUID)
			assert.NotNil(t, key.CreationDate)
		})
	}
}

func TestNewMasterKeyFromKeyIDString(t *testing.T) {
	tests := []struct {
		name           string
		keyIDString    string
		expectErr      bool
		expectKeyCount int
		expectKeys     []MasterKey
	}{
		{
			name:           "single key ID",
			keyIDString:    testKeyID1,
			expectKeyCount: 1,
			expectKeys: []MasterKey{
				{
					KeyID:   testKeyID1,
					Region:  "tr-west-1",
					KeyUUID: "12345678-1234-1234-1234-123456789abc",
				},
			},
		},
		{
			name:           "multiple key IDs",
			keyIDString:    testKeyID1 + "," + testKeyID2,
			expectKeyCount: 2,
			expectKeys: []MasterKey{
				{
					KeyID:   testKeyID1,
					Region:  "tr-west-1",
					KeyUUID: "12345678-1234-1234-1234-123456789abc",
				},
				{
					KeyID:   testKeyID2,
					Region:  "tr-west-1",
					KeyUUID: "87654321-4321-4321-4321-cba987654321",
				},
			},
		},
		{
			name:           "multiple key IDs with spaces",
			keyIDString:    " " + testKeyID1 + " , " + testKeyID2 + " ",
			expectKeyCount: 2,
			expectKeys: []MasterKey{
				{
					KeyID:   testKeyID1,
					Region:  "tr-west-1",
					KeyUUID: "12345678-1234-1234-1234-123456789abc",
				},
				{
					KeyID:   testKeyID2,
					Region:  "tr-west-1",
					KeyUUID: "87654321-4321-4321-4321-cba987654321",
				},
			},
		},
		{
			name:           "empty string",
			keyIDString:    "",
			expectKeyCount: 0,
		},
		{
			name:        "invalid key ID in list",
			keyIDString: testKeyID1 + ",invalid-key-id",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := NewMasterKeyFromKeyIDString(tt.keyIDString)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, keys)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectKeyCount, len(keys))
			for i, expectedKey := range tt.expectKeys {
				if i < len(keys) {
					assert.Equal(t, expectedKey.KeyID, keys[i].KeyID)
					assert.Equal(t, expectedKey.Region, keys[i].Region)
					assert.Equal(t, expectedKey.KeyUUID, keys[i].KeyUUID)
				}
			}
		})
	}
}

func TestMasterKey_ToString(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)
	assert.Equal(t, testKeyID1, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)
	key.EncryptedKey = "test-encrypted-key"
	key.CreationDate = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	m := key.ToMap()
	assert.Equal(t, testKeyID1, m["key_id"])
	assert.Equal(t, "test-encrypted-key", m["enc"])
	assert.Equal(t, "2025-01-01T12:00:00Z", m["created_at"])
}

func TestMasterKey_TypeToIdentifier(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)

	// New key should not need rotation
	assert.False(t, key.NeedsRotation())

	// Key older than TTL should need rotation
	key.CreationDate = time.Now().UTC().Add(-hckmsTTL - time.Hour)
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)
	key.EncryptedKey = "test-encrypted-data"
	assert.Equal(t, []byte("test-encrypted-data"), key.EncryptedDataKey())
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)
	key.SetEncryptedDataKey([]byte("test-encrypted-data"))
	assert.Equal(t, "test-encrypted-data", key.EncryptedKey)
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)

	// Key without encrypted data should attempt encryption
	// (will fail without credentials, but that's expected)
	err = key.EncryptIfNeeded([]byte("test-data"))
	// We expect an error because we don't have real credentials
	assert.Error(t, err)

	// Key with encrypted data should not attempt encryption
	key.EncryptedKey = "already-encrypted"
	err = key.EncryptIfNeeded([]byte("test-data"))
	assert.NoError(t, err)
	assert.Equal(t, "already-encrypted", key.EncryptedKey)
}

func TestCredentials_ApplyToMasterKey(t *testing.T) {
	basicCred := auth.NewBasicCredentialsBuilder().
		WithAk("test-ak").
		WithSk("test-sk").
		Build()
	cred := NewCredentials(basicCred)
	key, err := NewMasterKey(testKeyID1)
	assert.NoError(t, err)
	cred.ApplyToMasterKey(key)
	assert.Equal(t, cred.credential, key.credentials)
}
