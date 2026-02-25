//go:build integration

package tencentkms

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testTencentKMSKeyIDEnvVar = "SOPS_TENCENT_KMS_KEY_ID"

// TestEncryptDecryptIntegration tests the full encrypt-decrypt cycle with real KMS operations
// This test requires valid Tencent KMS credentials to be set as environment variables
func TestEncryptDecryptIntegration(t *testing.T) {
	keyID := requireTencentKMSEnv(t)

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

func requireTencentKMSEnv(t *testing.T) string {
	t.Helper()

	requiredVars := []struct {
		name  string
		value string
	}{
		{name: TencentSecretIdEnvVar, value: os.Getenv(TencentSecretIdEnvVar)},
		{name: TencentSecretKeyEnvVar, value: os.Getenv(TencentSecretKeyEnvVar)},
		{name: TencentRegionEnvVar, value: os.Getenv(TencentRegionEnvVar)},
		{name: testTencentKMSKeyIDEnvVar, value: os.Getenv(testTencentKMSKeyIDEnvVar)},
	}

	missing := make([]string, 0)
	var keyID string

	for _, envVar := range requiredVars {
		if envVar.value == "" {
			missing = append(missing, envVar.name)
			continue
		}

		if envVar.name == testTencentKMSKeyIDEnvVar {
			keyID = envVar.value
		}
	}

	if len(missing) > 0 {
		t.Skipf("skip: please configure %s before running integration tests", strings.Join(missing, ", "))
	}

	return keyID
}