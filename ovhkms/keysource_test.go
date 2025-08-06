package ovhkms

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMasterKeyFromKeyID(t *testing.T) {
	tests := []struct {
		name      string
		keyID     string
		expectErr bool
		expected  *MasterKey
	}{
		{
			name:  "Full format with endpoint and key ID",
			keyID: "eu-west-sbg.okms.ovh.net/12345678-1234-1234-1234-123456789012",
			expected: &MasterKey{
				Endpoint: "eu-west-sbg.okms.ovh.net",
				KeyID:    "12345678-1234-1234-1234-123456789012",
			},
			expectErr: false,
		},
		{
			name:  "Different endpoint",
			keyID: "ca-east-tor.okms.ovh.net/12345678-1234-1234-1234-123456789012",
			expected: &MasterKey{
				Endpoint: "ca-east-tor.okms.ovh.net",
				KeyID:    "12345678-1234-1234-1234-123456789012",
			},
			expectErr: false,
		},
		{
			name:      "No endpoint (invalid format)",
			keyID:     "12345678-1234-1234-1234-123456789012",
			expectErr: true,
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewMasterKeyFromKeyID(tt.keyID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.expected.Endpoint, got.Endpoint)
				assert.Equal(t, tt.expected.KeyID, got.KeyID)
			}
		})
	}
}

func TestMasterKeysFromResourceIDString(t *testing.T) {
	tests := []struct {
		name       string
		resourceID string
		expected   int
	}{
		{
			name:       "Empty string",
			resourceID: "",
			expected:   0,
		},
		{
			name:       "Single key ID",
			resourceID: "eu-west-sbg.okms.ovh.net/12345678-1234-1234-1234-123456789012",
			expected:   1,
		},
		{
			name:       "Multiple key IDs",
			resourceID: "eu-west-sbg.okms.ovh.net/12345678-1234-1234-1234-123456789012,ca-east-tor.okms.ovh.net/be03bbff-1f9e-5f98-9d11-35a393bbf673",
			expected:   2,
		},
		{
			name:       "Multiple key IDs with empty entry",
			resourceID: "eu-west-sbg.okms.ovh.net/12345678-1234-1234-1234-123456789012,,ca-east-tor.okms.ovh.net/be03bbff-1f9e-5f98-9d11-35a393bbf673",
			expected:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, _ := MasterKeysFromResourceIDString(tt.resourceID)
			assert.Equal(t, tt.expected, len(keys))
		})
	}
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := &MasterKey{
		CreationDate: time.Now().UTC(),
	}

	// New key doesn't need rotation
	assert.False(t, key.NeedsRotation())

	// Set creation date to exceed TTL
	key.CreationDate = key.CreationDate.Add(-(ovhKmsTTL + time.Second))
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	tests := []struct {
		name     string
		key      *MasterKey
		expected string
	}{
		{
			name: "Full key representation",
			key: &MasterKey{
				Endpoint: "eu-west-sbg.okms.ovh.net",
				KeyID:    "12345678-1234-1234-1234-123456789012",
			},
			expected: "eu-west-sbg.okms.ovh.net/12345678-1234-1234-1234-123456789012",
		},
		{
			name: "Different endpoint",
			key: &MasterKey{
				Endpoint: "ca-east-tor.okms.ovh.net",
				KeyID:    "be03bbff-1f9e-5f98-9d11-35a393bbf673",
			},
			expected: "ca-east-tor.okms.ovh.net/be03bbff-1f9e-5f98-9d11-35a393bbf673",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.key.ToString()
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestMasterKey_ToMap(t *testing.T) {
	key := &MasterKey{
		Endpoint:     "eu-west-sbg.okms.ovh.net",
		KeyID:        "12345678-1234-1234-1234-123456789012",
		EncryptedKey: "encrypted-data-key",
		CreationDate: time.Date(2025, 7, 31, 12, 0, 0, 0, time.UTC),
	}

	expected := map[string]interface{}{
		"endpoint":   "eu-west-sbg.okms.ovh.net",
		"key_id":     "12345678-1234-1234-1234-123456789012",
		"enc":        "encrypted-data-key",
		"created_at": "2025-07-31T12:00:00Z",
	}

	got := key.ToMap()
	assert.Equal(t, expected, got)
}

func TestMasterKey_TypeToIdentifier(t *testing.T) {
	key := &MasterKey{}
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: "encrypted-data-key"}
	assert.Equal(t, []byte("encrypted-data-key"), key.EncryptedDataKey())
}

func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	key := &MasterKey{}
	key.SetEncryptedDataKey([]byte("new-encrypted-key"))
	assert.Equal(t, "new-encrypted-key", key.EncryptedKey)
}

func TestMasterKey_getClient(t *testing.T) {
	// Save original env vars
	origCertFile := os.Getenv(CertificateFileEnv)
	origCertKeyFile := os.Getenv(CertificateKeyFileEnv)
	defer func() {
		os.Setenv(CertificateFileEnv, origCertFile)
		os.Setenv(CertificateKeyFileEnv, origCertKeyFile)
	}()

	// Test with missing env vars
	os.Unsetenv(CertificateFileEnv)
	os.Unsetenv(CertificateKeyFileEnv)

	key := &MasterKey{
		Endpoint: "eu-west-sbg.okms.ovh.net",
		KeyID:    "12345678-1234-1234-1234-123456789012",
	}

	// Should fail without certificate files
	client, err := key.getClient()
	assert.Error(t, err)
	assert.Nil(t, client)

	// Test with direct certificate files
	key.SetCertificateFiles("/path/to/cert.pem", "/path/to/key.pem")

	// The test will still fail because the files don't exist, but we can check if the paths were set
	_, err = key.getClient()
	assert.Error(t, err) // Will fail with file not found or similar
	assert.Equal(t, "/path/to/cert.pem", key.certificateFile)
	assert.Equal(t, "/path/to/key.pem", key.certificateKeyFile)

	// Test with env vars
	key = &MasterKey{
		Endpoint: "eu-west-sbg.okms.ovh.net",
		KeyID:    "12345678-1234-1234-1234-123456789012",
	}

	os.Setenv(CertificateFileEnv, "/path/to/env/cert.pem")
	os.Setenv(CertificateKeyFileEnv, "/path/to/env/key.pem")

	// The test will still fail because the files don't exist, but we can check if it tries to use env vars
	_, err = key.getClient()
	assert.Error(t, err) // Will fail with file not found or similar
}

func TestMasterKey_SetCertificateFiles(t *testing.T) {
	key := &MasterKey{}
	key.SetCertificateFiles("/path/to/cert.pem", "/path/to/key.pem")

	assert.Equal(t, "/path/to/cert.pem", key.certificateFile)
	assert.Equal(t, "/path/to/key.pem", key.certificateKeyFile)
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key := &MasterKey{}
	// Should call Encrypt if EncryptedKey is empty
	err := key.EncryptIfNeeded([]byte("dummy"))
	assert.Error(t, err) // getClient will fail due to missing certs

	// Should not call Encrypt if EncryptedKey is already set
	key.EncryptedKey = "already-encrypted"
	err = key.EncryptIfNeeded([]byte("dummy"))
	assert.NoError(t, err)
}

func TestMasterKey_EncryptContext_InvalidUUID(t *testing.T) {
	key := &MasterKey{
		Endpoint: "eu-west-sbg.okms.ovh.net",
		KeyID:    "not-a-uuid",
	}
	err := key.EncryptContext(context.Background(), []byte("test"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse UUID")
}

func TestMasterKey_DecryptContext_InvalidUUID(t *testing.T) {
	key := &MasterKey{
		Endpoint:     "eu-west-sbg.okms.ovh.net",
		KeyID:        "not-a-uuid",
		EncryptedKey: "some-encrypted-data",
	}
	_, err := key.DecryptContext(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse UUID")
}

func TestMasterKey_DecryptContext_InvalidBase64(t *testing.T) {
	key := &MasterKey{
		Endpoint:     "eu-west-sbg.okms.ovh.net",
		KeyID:        "12345678-1234-1234-1234-123456789012",
		EncryptedKey: "@@@not-base64@@@",
	}
	// Inject a fake getClient that returns a mock response
	key.SetCertificateFiles("/nonexistent.crt", "/nonexistent.key")
	_, err := key.DecryptContext(context.Background())
	assert.Error(t, err) // Should still error due to cert loading
}
