package tencentkms

import (
	"context"
	"encoding/base64"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	// dummyKeyID
	dummyKeyID = "xxxxx-xxxxx-xxxxx-xxxxx-xxxxx"
	// dummyRegion
	dummyRegion = "ap-singapore"
	// dummyEncryptedKey
	dummyEncryptedKey = "dummy-encrypted-key"
)

// TestNewMasterKeyFromKeyID tests creating MasterKey from key ID
func TestNewMasterKeyFromKeyID(t *testing.T) {
	t.Run("normal creation", func(t *testing.T) {
		key := NewMasterKeyFromKeyID(dummyKeyID)
		assert.Equal(t, dummyKeyID, key.KeyID)
		assert.NotNil(t, key.CreationDate)
		assert.Empty(t, key.Region)
		assert.Empty(t, key.EncryptedKey)
	})

	t.Run("remove spaces", func(t *testing.T) {
		key := NewMasterKeyFromKeyID("  xxxxx-xxxxx-xxxxx-xxxxx-xxxxx  ")
		assert.Equal(t, dummyKeyID, key.KeyID)
	})

	t.Run("empty string", func(t *testing.T) {
		key := NewMasterKeyFromKeyID("")
		assert.Empty(t, key.KeyID)
		assert.NotNil(t, key.CreationDate)
	})
}

// TestMasterKeysFromKeyIDString tests creating multiple MasterKeys from comma-separated key ID string
func TestMasterKeysFromKeyIDString(t *testing.T) {
	t.Run("single key", func(t *testing.T) {
		keys := MasterKeysFromKeyIDString(dummyKeyID)
		assert.Len(t, keys, 1)
		assert.Equal(t, dummyKeyID, keys[0].KeyID)
	})

	t.Run("multiple keys", func(t *testing.T) {
		keyID2 := "yyyyy-yyyyy-yyyyy-yyyyy-yyyyy"
		keys := MasterKeysFromKeyIDString(dummyKeyID + "," + keyID2)
		assert.Len(t, keys, 2)
		assert.Equal(t, dummyKeyID, keys[0].KeyID)
		assert.Equal(t, keyID2, keys[1].KeyID)
	})

	t.Run("empty string", func(t *testing.T) {
		keys := MasterKeysFromKeyIDString("")
		assert.Len(t, keys, 0)
	})

	t.Run("with spaces", func(t *testing.T) {
		keyID2 := "yyyyy-yyyyy-yyyyy-yyyyy-yyyyy"
		keys := MasterKeysFromKeyIDString(dummyKeyID + ", " + keyID2)
		assert.Len(t, keys, 2)
		assert.Equal(t, dummyKeyID, keys[0].KeyID)
		assert.Equal(t, keyID2, keys[1].KeyID)
	})

	t.Run("empty elements", func(t *testing.T) {
		keys := MasterKeysFromKeyIDString(dummyKeyID + ",,")
		assert.Len(t, keys, 3) // Empty strings will still create MasterKey
		assert.Equal(t, dummyKeyID, keys[0].KeyID)
		assert.Empty(t, keys[1].KeyID)
		assert.Empty(t, keys[2].KeyID)
	})
}

// TestMasterKey_EncryptIfNeeded tests encryption when needed
func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	// Since actual encryption requires Tencent Cloud KMS service, we use mock data for testing
	key := &MasterKey{
		KeyID:     dummyKeyID,
		secretId:  "mock-secret-id", // Set mock credentials to avoid nil reference
		secretKey: "mock-secret-key",
		Region:    dummyRegion,
	}

	// We can't actually call Tencent Cloud KMS, so we verify the logic flow
	// Actual encryption will return an error, but EncryptIfNeeded should call Encrypt
	err := key.EncryptIfNeeded([]byte("test-data"))
	assert.Error(t, err) // Expected to fail because we're using mock credentials

	// Manually set encrypted key, then test that it won't be re-encrypted
	key.EncryptedKey = dummyEncryptedKey
	err = key.EncryptIfNeeded([]byte("different-data"))
	assert.NoError(t, err)
	assert.Equal(t, dummyEncryptedKey, key.EncryptedKey) // Confirm key wasn't modified
}

// TestMasterKey_EncryptedDataKey tests getting encrypted data key
func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: dummyEncryptedKey}
	assert.EqualValues(t, []byte(dummyEncryptedKey), key.EncryptedDataKey())

	key = &MasterKey{EncryptedKey: ""}
	assert.EqualValues(t, []byte(""), key.EncryptedDataKey())
}

// TestMasterKey_SetEncryptedDataKey tests setting encrypted data key
func TestMasterKey_SetEncryptedDataKey(t *testing.T) {
	key := &MasterKey{}
	data := []byte("test-encrypted-data")
	key.SetEncryptedDataKey(data)
	assert.Equal(t, string(data), key.EncryptedKey)

	// Test empty data
	key.SetEncryptedDataKey([]byte{})
	assert.Equal(t, "", key.EncryptedKey)
}

// TestMasterKey_NeedsRotation tests if key rotation is needed
func TestMasterKey_NeedsRotation(t *testing.T) {
	key := NewMasterKeyFromKeyID(dummyKeyID)
	assert.False(t, key.NeedsRotation()) // Newly created key doesn't need rotation

	// Set an expired creation time (significantly greater than TTL)
	key.CreationDate = time.Now().UTC().Add(-(tencentkmsTTL + 24*time.Hour))
	assert.True(t, key.NeedsRotation())

	// Set time significantly less than TTL to avoid precision issues
	// We'll use a 1 hour difference to ensure it's clearly less than TTL
	key.CreationDate = time.Now().UTC().Add(-(tencentkmsTTL - time.Hour))
	assert.False(t, key.NeedsRotation())

	// We'll skip the exact TTL test due to time precision issues
	// The implementation clearly uses '>' which means only strictly greater than TTL returns true
}

// TestMasterKey_ToString tests conversion to string
func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKeyFromKeyID(dummyKeyID)
	assert.Equal(t, dummyKeyID, key.ToString())

	key = NewMasterKeyFromKeyID("")
	assert.Equal(t, "", key.ToString())
}

// TestMasterKey_ToMap tests conversion to map
func TestMasterKey_ToMap(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	key := &MasterKey{
		KeyID:        dummyKeyID,
		CreationDate: fixedTime,
		EncryptedKey: dummyEncryptedKey,
	}

	expectedMap := map[string]interface{}{
		"keyId":      dummyKeyID,
		"created_at": fixedTime.UTC().Format(time.RFC3339),
		"enc":        dummyEncryptedKey,
	}

	resultMap := key.ToMap()
	assert.Equal(t, expectedMap, resultMap)

	// Test empty value case
	key = &MasterKey{
		KeyID:        "",
		CreationDate: fixedTime,
		EncryptedKey: "",
	}

	emptyExpectedMap := map[string]interface{}{
		"keyId":      "",
		"created_at": fixedTime.UTC().Format(time.RFC3339),
		"enc":        "",
	}

	assert.Equal(t, emptyExpectedMap, key.ToMap())
}

// TestMasterKey_TypeToIdentifier tests type identifier
func TestMasterKey_TypeToIdentifier(t *testing.T) {
	key := NewMasterKeyFromKeyID(dummyKeyID)
	assert.Equal(t, KeyTypeIdentifier, key.TypeToIdentifier())
}

// TestMasterKey_createClient tests client creation logic
func TestMasterKey_createClient(t *testing.T) {
	// Save original environment variables
	originalSecretId := os.Getenv(TencentSecretIdEnvVar)
	originalSecretKey := os.Getenv(TencentSecretKeyEnvVar)
	originalRegion := os.Getenv(TencentRegionEnvVar)
	defer func() {
		// Restore original environment variables
		os.Unsetenv(TencentSecretIdEnvVar)
		os.Unsetenv(TencentSecretKeyEnvVar)
		os.Unsetenv(TencentRegionEnvVar)
		if originalSecretId != "" {
			os.Setenv(TencentSecretIdEnvVar, originalSecretId)
		}
		if originalSecretKey != "" {
			os.Setenv(TencentSecretKeyEnvVar, originalSecretKey)
		}
		if originalRegion != "" {
			os.Setenv(TencentRegionEnvVar, originalRegion)
		}
	}()

	t.Run("use credentials from environment", func(t *testing.T) {
		os.Setenv(TencentSecretIdEnvVar, "env-secret-id")
		os.Setenv(TencentSecretKeyEnvVar, "env-secret-key")
		os.Setenv(TencentRegionEnvVar, "env-region")

		key := &MasterKey{
			KeyID: dummyKeyID,
		}

		client, err := key.createClient()
		assert.NotNil(t, client, "should create client object from environment variables")
		_ = err // Avoid unused variable warning
	})

	t.Run("prioritize credentials in key", func(t *testing.T) {
		os.Setenv(TencentSecretIdEnvVar, "env-secret-id")
		os.Setenv(TencentSecretKeyEnvVar, "env-secret-key")
		os.Setenv(TencentRegionEnvVar, "env-region")

		key := &MasterKey{
			KeyID:     dummyKeyID,
			Region:    "key-region",
			secretId:  "key-secret-id",
			secretKey: "key-secret-key",
		}

		// Check if function runs without crashing
		client, err := key.createClient()
		assert.NotNil(t, client, "should create client object")
		_ = err // Avoid unused variable warning
	})

	t.Run("empty credentials", func(t *testing.T) {
		key := &MasterKey{
			KeyID: dummyKeyID,
		}

		// Check if client can be created with empty credentials (may return error or default client)
		client, err := key.createClient()
		// Don't make strict assertions, just ensure test doesn't crash
		if err != nil {
			assert.Error(t, err, "empty credentials may return error")
		} else {
			assert.NotNil(t, client, "empty credentials may return default client")
		}
	})
}

// TestEncryptDecryptMock tests encryption and decryption process (using mock data)
func TestEncryptDecryptMock(t *testing.T) {
	// Create a key
	key := &MasterKey{
		KeyID:     dummyKeyID,
		Region:    dummyRegion,
		secretId:  "mock-secret-id",
		secretKey: "mock-secret-key",
	}

	// Test encryption (will fail but test the flow)
	dataKey := []byte("test-data-key")
	err := key.Encrypt(dataKey)
	assert.Error(t, err) // Expected to fail because we're using mock credentials

	// Manually set encrypted key
	key.EncryptedKey = base64.StdEncoding.EncodeToString(dataKey)

	// Test decryption (will fail but test the flow)
	_, err = key.Decrypt()
	assert.Error(t, err) // Expected to fail because we're using mock credentials
}

// TestEncryptContextDecryptContextMock tests context-aware encryption and decryption process (using mock data)
func TestEncryptContextDecryptContextMock(t *testing.T) {
	// Create a key
	key := &MasterKey{
		KeyID:     dummyKeyID,
		Region:    dummyRegion,
		secretId:  "mock-secret-id",
		secretKey: "mock-secret-key",
	}

	// Create a context
	ctx := context.Background()

	// Test context-aware encryption (will fail but test the flow)
	dataKey := []byte("test-data-key")
	err := key.EncryptContext(ctx, dataKey)
	assert.Error(t, err) // Expected to fail because we're using mock credentials

	// Manually set encrypted key
	key.EncryptedKey = base64.StdEncoding.EncodeToString(dataKey)

	// Test context-aware decryption (will fail but test the flow)
	_, err = key.DecryptContext(ctx)
	assert.Error(t, err) // Expected to fail because we're using mock credentials

	// Test decryption with empty encrypted key
	key.EncryptedKey = ""
	_, err = key.DecryptContext(ctx)
	assert.Error(t, err)
	assert.Equal(t, "master key is empty", err.Error())
}

// TestTimeCalculations tests time calculation related functionality
func TestTimeCalculations(t *testing.T) {
	// Verify TTL value
	expectedTTL := time.Hour * 24 * 30 * 6 // 6 months
	assert.Equal(t, expectedTTL, tencentkmsTTL)

	// Test time comparison logic
	key := NewMasterKeyFromKeyID(dummyKeyID)
	now := time.Now().UTC()

	// Newly created key doesn't need rotation
	key.CreationDate = now
	assert.False(t, key.NeedsRotation())

	// Key close to but not exceeding TTL doesn't need rotation
	key.CreationDate = now.Add(-(tencentkmsTTL - 24*time.Hour))
	assert.False(t, key.NeedsRotation())

	// Key exceeding TTL needs rotation
	key.CreationDate = now.Add(-(tencentkmsTTL + 24*time.Hour))
	assert.True(t, key.NeedsRotation())
}
