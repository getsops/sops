//go:build integration

package tencentkms

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEncryptDecryptIntegration tests the full encrypt-decrypt cycle with real KMS operations
// This test requires valid Tencent KMS credentials to be set as environment variables
func TestEncryptDecryptIntegration(t *testing.T) {
	// Skip test if credentials are not set
	_ = os.Setenv(TencentSecretIdEnvVar, "")
	_ = os.Setenv(TencentSecretKeyEnvVar, "")
	_ = os.Setenv(TencentRegionEnvVar, "ap-singapore")
	keyID := ""

	// Test cases for encryption and decryption
	testCases := []struct {
		name        string
		plaintext   string
		expectError bool
	}{{
		name:        "Simple string encryption",
		plaintext:   "Hello, Tencent KMS!",
		expectError: false,
	}, {
		name:        "Empty string encryption",
		plaintext:   "",
		expectError: true,
	}, {
		name:        "Long string encryption",
		plaintext:   generateLongString(1000), // 1000 characters
		expectError: false,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new master key
			masterKey := NewMasterKeyFromKeyID(keyID)
			masterKey.CreationDate = time.Now().UTC()

			err := masterKey.Encrypt([]byte(tc.plaintext))
			if tc.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, masterKey.EncryptedKey)

			// Decrypt the data key
			decryptedKey, err := masterKey.Decrypt()
			assert.NoError(t, err)
			assert.NotNil(t, decryptedKey)
			assert.Equal(t, tc.plaintext, string(decryptedKey))

			// Verify key rotation status (should not need rotation if recently created)
			assert.False(t, masterKey.NeedsRotation(), "Newly created key should not need rotation")
		})
	}
}

// generateLongString creates a string of the specified length for testing
func generateLongString(length int) string {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte(65 + (i % 26)) // A-Z characters
	}
	return string(result)
}
