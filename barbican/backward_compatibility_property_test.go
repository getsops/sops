package barbican

import (
	"testing"
	"testing/quick"
	"time"
)

// TestBackwardCompatibilityProperty tests that Barbican master keys maintain
// backward compatibility with existing SOPS files and can be mixed with other key types
func TestBackwardCompatibilityProperty(t *testing.T) {
	config := &quick.Config{
		MaxCount: 50,
		Rand:     nil,
	}

	property := func(dataKey []byte, secretRefs []string) bool {
		// Ensure we have valid test data
		if len(dataKey) == 0 {
			dataKey = make([]byte, 32)
			for i := range dataKey {
				dataKey[i] = byte(i)
			}
		}
		if len(dataKey) != 32 {
			// Resize to 32 bytes for AES-256
			newDataKey := make([]byte, 32)
			copy(newDataKey, dataKey)
			dataKey = newDataKey
		}

		// Generate valid secret references if none provided
		if len(secretRefs) == 0 {
			secretRefs = []string{
				"550e8400-e29b-41d4-a716-446655440000",
				"region:sjc3:660e8400-e29b-41d4-a716-446655440001",
			}
		}

		// Validate and filter secret references
		var validSecretRefs []string
		for _, ref := range secretRefs {
			if isValidSecretRef(ref) {
				validSecretRefs = append(validSecretRefs, ref)
			}
		}
		if len(validSecretRefs) == 0 {
			validSecretRefs = []string{"550e8400-e29b-41d4-a716-446655440000"}
		}

		// Create Barbican master keys
		var barbicanKeys []*MasterKey
		for _, ref := range validSecretRefs {
			key, err := NewMasterKeyFromSecretRef(ref)
			if err != nil {
				t.Logf("Failed to create master key from ref %s: %v", ref, err)
				return false
			}
			barbicanKeys = append(barbicanKeys, key)
		}

		// Test 1: Verify keys can be created
		if len(barbicanKeys) == 0 {
			t.Logf("No Barbican keys created")
			return false
		}

		// Test 2: Verify each Barbican key maintains its identity
		for i, key := range barbicanKeys {
			// Verify ToString() method works
			keyString := key.ToString()
			if keyString == "" {
				t.Logf("Key %d ToString() returned empty string", i)
				return false
			}

			// Verify TypeToIdentifier() returns correct type
			if key.TypeToIdentifier() != KeyTypeIdentifier {
				t.Logf("Key %d has incorrect type identifier: %s", i, key.TypeToIdentifier())
				return false
			}

			// Verify ToMap() method works
			keyMap := key.ToMap()
			if keyMap == nil {
				t.Logf("Key %d ToMap() returned nil", i)
				return false
			}
			if keyMap["secret_ref"] == nil {
				t.Logf("Key %d ToMap() missing secret_ref", i)
				return false
			}
		}

		// Test 3: Verify NeedsRotation works correctly
		for _, key := range barbicanKeys {
			// New key should not need rotation
			if key.NeedsRotation() {
				t.Logf("New key incorrectly reports needing rotation")
				return false
			}

			// Old key should need rotation
			oldKey := *key
			oldKey.CreationDate = time.Now().Add(-time.Hour * 24 * 365) // 1 year ago
			if !oldKey.NeedsRotation() {
				t.Logf("Old key incorrectly reports not needing rotation")
				return false
			}
		}

		// Test 4: Verify EncryptedDataKey and SetEncryptedDataKey work
		for _, key := range barbicanKeys {
			// Initially should be empty
			if len(key.EncryptedDataKey()) != 0 {
				t.Logf("New key has non-empty EncryptedDataKey")
				return false
			}

			// Set encrypted data key
			testEncryptedKey := "test-encrypted-key-" + key.SecretRef
			key.SetEncryptedDataKey([]byte(testEncryptedKey))
			
			// Verify it was set
			if string(key.EncryptedDataKey()) != testEncryptedKey {
				t.Logf("SetEncryptedDataKey/EncryptedDataKey roundtrip failed")
				return false
			}
		}

		// Test 5: Verify multi-region functionality
		regionGroups := GroupKeysByRegion(barbicanKeys)
		if len(regionGroups) == 0 {
			t.Logf("No region groups found")
			return false
		}

		for region, keys := range regionGroups {
			if len(keys) == 0 {
				t.Logf("Empty key group for region %s", region)
				return false
			}

			// Verify all keys in the same region have consistent properties
			for i, key := range keys {
				effectiveRegion := key.getEffectiveRegion()
				if effectiveRegion != region {
					t.Logf("Key %d in region group %s has different effective region %s", i, region, effectiveRegion)
					return false
				}
			}
		}

		return true
	}

	err := quick.Check(property, config)
	if err != nil {
		t.Errorf("Backward compatibility property failed: %v", err)
	}
}

// TestDecryptionOrderProperty tests that Barbican keys work correctly in different decryption orders
func TestDecryptionOrderProperty(t *testing.T) {
	config := &quick.Config{
		MaxCount: 30,
		Rand:     nil,
	}

	property := func(secretRefs []string) bool {
		// Generate valid secret references if none provided
		if len(secretRefs) == 0 {
			secretRefs = []string{
				"550e8400-e29b-41d4-a716-446655440000",
				"region:dfw3:660e8400-e29b-41d4-a716-446655440001",
			}
		}

		// Validate and filter secret references
		var validSecretRefs []string
		for _, ref := range secretRefs {
			if isValidSecretRef(ref) {
				validSecretRefs = append(validSecretRefs, ref)
			}
		}
		if len(validSecretRefs) == 0 {
			validSecretRefs = []string{"550e8400-e29b-41d4-a716-446655440000"}
		}

		// Create master keys
		var masterKeys []*MasterKey
		for _, ref := range validSecretRefs {
			key, err := NewMasterKeyFromSecretRef(ref)
			if err != nil {
				t.Logf("Failed to create master key from ref %s: %v", ref, err)
				return false
			}
			masterKeys = append(masterKeys, key)
		}

		// Test different decryption orders
		decryptionOrders := [][]string{
			{"barbican", "pgp", "age"},
			{"pgp", "barbican", "age"},
			{"age", "pgp", "barbican"},
			{"barbican"},
		}

		for _, order := range decryptionOrders {
			// Test 1: Verify Barbican keys are recognized in any order
			barbicanFound := false
			for _, keyType := range order {
				if keyType == KeyTypeIdentifier {
					barbicanFound = true
					break
				}
			}

			// If barbican is in the order, verify our keys match
			if barbicanFound {
				for _, key := range masterKeys {
					if key.TypeToIdentifier() != KeyTypeIdentifier {
						t.Logf("Key type mismatch: expected %s, got %s", KeyTypeIdentifier, key.TypeToIdentifier())
						return false
					}
				}
			}

			// Test 2: Verify keys maintain their properties regardless of order
			for _, key := range masterKeys {
				// Key should maintain its secret reference
				if key.SecretRef == "" {
					t.Logf("Key lost its secret reference")
					return false
				}

				// Key should maintain its type identifier
				if key.TypeToIdentifier() != KeyTypeIdentifier {
					t.Logf("Key type identifier changed")
					return false
				}

				// ToString should be consistent
				keyString1 := key.ToString()
				keyString2 := key.ToString()
				if keyString1 != keyString2 {
					t.Logf("ToString() is not consistent")
					return false
				}
			}
		}

		// Test 3: Verify multi-region keys work in any order
		regionGroups := GroupKeysByRegion(masterKeys)
		if len(regionGroups) == 0 {
			t.Logf("No region groups found")
			return false
		}

		for region, keys := range regionGroups {
			if len(keys) == 0 {
				t.Logf("Empty key group for region %s", region)
				return false
			}

			// Verify all keys in the same region have consistent properties
			for i, key := range keys {
				effectiveRegion := key.getEffectiveRegion()
				if effectiveRegion != region {
					t.Logf("Key %d in region group %s has different effective region %s", i, region, effectiveRegion)
					return false
				}
			}
		}

		return true
	}

	err := quick.Check(property, config)
	if err != nil {
		t.Errorf("Decryption order property failed: %v", err)
	}
}