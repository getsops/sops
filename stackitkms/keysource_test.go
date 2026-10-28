package stackitkms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	testResourceID1 = "projects/test-project-id/regions/eu01/keyRings/test-keyring-id/keys/test-key-id/versions/1"
	testResourceID2 = "projects/test-project-id/regions/eu01/keyRings/test-keyring-id/keys/test-key-id-2/versions/2"
)

func TestNewMasterKey(t *testing.T) {
	tests := []struct {
		name      string
		resID     string
		expectErr bool
		expectKey MasterKey
	}{
		{
			name:  "valid resource ID",
			resID: testResourceID1,
			expectKey: MasterKey{
				ResourceID:    testResourceID1,
				ProjectID:     "test-project-id",
				RegionID:      "eu01",
				KeyRingID:     "test-keyring-id",
				KeyID:         "test-key-id",
				VersionNumber: 1,
			},
		},
		{
			name:      "invalid format - too few parts",
			resID:     "projects/foo/regions/bar",
			expectErr: true,
		},
		{
			name:      "invalid format - wrong prefix",
			resID:     "proj/foo/regions/bar/keyRings/baz/keys/qux/versions/1",
			expectErr: true,
		},
		{
			name:      "invalid format - bad version number",
			resID:     "projects/foo/regions/bar/keyRings/baz/keys/qux/versions/abc",
			expectErr: true,
		},
		{
			name:      "invalid format - empty project",
			resID:     "projects//regions/bar/keyRings/baz/keys/qux/versions/1",
			expectErr: true,
		},
		{
			name:      "invalid format - empty string",
			resID:     "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewMasterKey(tt.resID)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, key)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectKey.ResourceID, key.ResourceID)
			assert.Equal(t, tt.expectKey.ProjectID, key.ProjectID)
			assert.Equal(t, tt.expectKey.RegionID, key.RegionID)
			assert.Equal(t, tt.expectKey.KeyRingID, key.KeyRingID)
			assert.Equal(t, tt.expectKey.KeyID, key.KeyID)
			assert.Equal(t, tt.expectKey.VersionNumber, key.VersionNumber)
			assert.NotNil(t, key.CreationDate)
		})
	}
}

func TestNewMasterKeyFromResourceIDString(t *testing.T) {
	tests := []struct {
		name           string
		resIDString    string
		expectErr      bool
		expectKeyCount int
	}{
		{
			name:           "single resource ID",
			resIDString:    testResourceID1,
			expectKeyCount: 1,
		},
		{
			name:           "multiple resource IDs",
			resIDString:    testResourceID1 + "," + testResourceID2,
			expectKeyCount: 2,
		},
		{
			name:           "multiple with spaces",
			resIDString:    " " + testResourceID1 + " , " + testResourceID2 + " ",
			expectKeyCount: 2,
		},
		{
			name:           "empty string",
			resIDString:    "",
			expectKeyCount: 0,
		},
		{
			name:        "invalid resource ID in list",
			resIDString: testResourceID1 + ",invalid-id",
			expectErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := NewMasterKeyFromResourceIDString(tt.resIDString)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, keys)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expectKeyCount, len(keys))
		})
	}
}

func TestMasterKey_ToString(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)
	assert.Equal(t, testResourceID1, key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)
	key.EncryptedKey = "test-encrypted-key"
	key.CreationDate = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)

	m := key.ToMap()
	assert.Equal(t, testResourceID1, m["resource_id"])
	assert.Equal(t, "test-encrypted-key", m["enc"])
	assert.Equal(t, "2025-01-01T12:00:00Z", m["created_at"])
}

func TestMasterKey_TypeToIdentifier(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)

	// New key should not need rotation
	assert.False(t, key.NeedsRotation())

	// Key older than TTL should need rotation
	key.CreationDate = time.Now().UTC().Add(-stackitKmsTTL - time.Hour)
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)
	key.EncryptedKey = "test-encrypted-data"
	assert.Equal(t, []byte("test-encrypted-data"), key.EncryptedDataKey())
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)
	key.SetEncryptedDataKey([]byte("test-encrypted-data"))
	assert.Equal(t, "test-encrypted-data", key.EncryptedKey)
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key, err := NewMasterKey(testResourceID1)
	assert.NoError(t, err)

	// Key with encrypted data should not attempt encryption
	key.EncryptedKey = "already-encrypted"
	err = key.EncryptIfNeeded([]byte("test-data"))
	assert.NoError(t, err)
	assert.Equal(t, "already-encrypted", key.EncryptedKey)
}
