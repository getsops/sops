package azkv

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAzureKeySourceFromUrl(t *testing.T) {
	cases := []struct {
		name              string
		input             string
		expectSuccess     bool
		expectedFoundKeys int
		expectedKeys      []MasterKey
	}{
		{
			name:              "Single url",
			input:             "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90",
			expectSuccess:     true,
			expectedFoundKeys: 1,
			expectedKeys: []MasterKey{
				{
					VaultURL: "https://test.vault.azure.net",
					Name:     "test-key",
					Version:  "a2a690a4fcc04166b739da342a912c90",
				},
			},
		},
		{
			name:              "Multiple url",
			input:             "https://test.vault.azure.net/keys/test-key/a2a690a4fcc04166b739da342a912c90,https://test2.vault.azure.net/keys/another-test-key/cf0021e8b743453bae758e7fbf71b60e",
			expectSuccess:     true,
			expectedFoundKeys: 2,
			expectedKeys: []MasterKey{
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
			name:          "Single malformed url",
			input:         "https://test.vault.azure.net/no-keys-here/test-key/a2a690a4fcc04166b739da342a912c90",
			expectSuccess: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			keys, err := MasterKeysFromURLs(c.input)
			if err != nil && c.expectSuccess {
				t.Fatalf("Unexpected error %v", err)
			} else if err == nil && !c.expectSuccess {
				t.Fatal("Expected error, but no error was returned")
			}

			if c.expectedFoundKeys != len(keys) {
				t.Errorf("Unexpected number of keys returned, expected %d, got %d", c.expectedFoundKeys, len(keys))
			}
			for idx := range keys {
				assert.Equal(t, c.expectedKeys[idx].VaultURL, keys[idx].VaultURL)
				assert.Equal(t, c.expectedKeys[idx].Name, keys[idx].Name)
				assert.Equal(t, c.expectedKeys[idx].Version, keys[idx].Version)
			}
		})
	}
}

func TestKeyToMap(t *testing.T) {
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

var azureKeyAcceptanceTestUrl = flag.String("azure-key", "", "URL to Azure Key Vault (note that this can incur real costs!)")

func TestRoundtrip(t *testing.T) {
	if *azureKeyAcceptanceTestUrl == "" {
		t.Skip("Azure URL not provided, skipping acceptance test")
	}

	input := []byte("test-string")

	key, err := NewMasterKeyFromURL(*azureKeyAcceptanceTestUrl)
	if err != nil {
		t.Fatal(err)
	}
	err = key.Encrypt(input)
	if err != nil {
		t.Fatal(err)
	}

	output, err := key.Decrypt()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, input, output)
}
