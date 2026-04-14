package acskms

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testARN1 = "acs:kms:cn-shanghai:123456789012:key/key-abc123def456"
	testARN2 = "acs:kms:cn-hangzhou:123456789012:key/key-xyz789uvw012"
	testARN3 = "acs:kms:cn-beijing:123456789012:alias/my-alias"
)

func TestNewMasterKey(t *testing.T) {
	tests := []struct {
		name          string
		arn           string
		expectErr     bool
		expectRegion  string
	}{
		{
			name:         "valid key ARN",
			arn:          testARN1,
			expectRegion: "cn-shanghai",
		},
		{
			name:         "valid key ARN different region",
			arn:          testARN2,
			expectRegion: "cn-hangzhou",
		},
		{
			name:         "valid alias ARN",
			arn:          testARN3,
			expectRegion: "cn-beijing",
		},
		{
			name:         "valid ARN with leading/trailing whitespace",
			arn:          "  " + testARN1 + "  ",
			expectRegion: "cn-shanghai",
		},
		{
			name:      "invalid ARN - wrong prefix",
			arn:       "aws:kms:cn-shanghai:123456789012:key/key-abc123",
			expectErr: true,
		},
		{
			name:      "invalid ARN - missing region",
			arn:       "acs:kms::123456789012:key/key-abc123",
			expectErr: true,
		},
		{
			name:      "invalid ARN - empty string",
			arn:       "",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := NewMasterKey(tt.arn)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, key)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, key)
			assert.Equal(t, tt.expectRegion, key.region)
		})
	}
}

func TestMasterKeysFromARNString(t *testing.T) {
	tests := []struct {
		name      string
		arnList   string
		expectN   int
		expectErr bool
	}{
		{
			name:    "empty string returns empty slice",
			arnList: "",
			expectN: 0,
		},
		{
			name:    "single ARN",
			arnList: testARN1,
			expectN: 1,
		},
		{
			name:    "multiple ARNs comma-separated",
			arnList: testARN1 + "," + testARN2,
			expectN: 2,
		},
		{
			name:    "multiple ARNs with spaces",
			arnList: testARN1 + " , " + testARN2,
			expectN: 2,
		},
		{
			name:      "invalid ARN in list",
			arnList:   testARN1 + ",invalid-arn",
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, err := MasterKeysFromARNString(tt.arnList)
			if tt.expectErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, keys, tt.expectN)
		})
	}
}

func TestMasterKey_ToString(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)
	assert.Equal(t, testARN1, key.ToString())
}

func TestMasterKey_TypeToIdentifier(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
	assert.Equal(t, "acs_kms", key.TypeToIdentifier())
}

func TestMasterKey_ToMap(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)
	key.EncryptedKey = "some-encrypted-value"

	m := key.ToMap()
	assert.Equal(t, testARN1, m["arn"])
	assert.Equal(t, "some-encrypted-value", m["enc"])
	assert.NotEmpty(t, m["created_at"])
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)
	key.EncryptedKey = "test-encrypted-key"
	assert.Equal(t, []byte("test-encrypted-key"), key.EncryptedDataKey())
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)
	key.SetEncryptedDataKey([]byte("new-encrypted-key"))
	assert.Equal(t, "new-encrypted-key", key.EncryptedKey)
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)
	key.EncryptedKey = "already-encrypted"

	// Should not call Encrypt when EncryptedKey is already set.
	err = key.EncryptIfNeeded([]byte("data-key"))
	assert.NoError(t, err)
	assert.Equal(t, "already-encrypted", key.EncryptedKey)
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key, err := NewMasterKey(testARN1)
	require.NoError(t, err)

	// Fresh key does not need rotation.
	assert.False(t, key.NeedsRotation())

	// Key older than 6 months needs rotation.
	key.CreationDate = time.Now().Add(-acskmsKeyTTL - time.Hour)
	assert.True(t, key.NeedsRotation())
}

func TestRegionFromARN(t *testing.T) {
	tests := []struct {
		arn           string
		expectRegion  string
		expectErr     bool
	}{
		{"acs:kms:cn-shanghai:123:key/key-abc", "cn-shanghai", false},
		{"acs:kms:cn-hangzhou:123:alias/myalias", "cn-hangzhou", false},
		{"acs:kms:ap-southeast-1:123:key/key-xyz", "ap-southeast-1", false},
		{"invalid", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		region, err := regionFromARN(tt.arn)
		if tt.expectErr {
			assert.Error(t, err, "expected error for ARN %q", tt.arn)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.expectRegion, region)
		}
	}
}

func TestLoadAliyunConfigCredentials_NoFile(t *testing.T) {
	t.Setenv("ALIBABA_CLOUD_CONFIG_FILE", "/nonexistent/path/config.json")
	_, err := loadAliyunConfigCredentials()
	assert.Error(t, err)
}

func TestBase64RoundTrip(t *testing.T) {
	// Verify the base64 encoding we use for the KMS Plaintext field is consistent.
	original := []byte("32-byte-sops-data-key-here!!!!!")
	encoded := base64.StdEncoding.EncodeToString(original)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	assert.Equal(t, original, decoded)
}
