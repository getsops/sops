package stores

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/azkv"
	"github.com/stretchr/testify/assert"
)

func TestValToString(t *testing.T) {
	assert.Equal(t, "1", ValToString(1))
	assert.Equal(t, "1.0", ValToString(1.0))
	assert.Equal(t, "1.1", ValToString(1.10))
	assert.Equal(t, "1.23", ValToString(1.23))
	assert.Equal(t, "1.2345678901234567", ValToString(1.234567890123456789))
	assert.Equal(t, "200000.0", ValToString(2e5))
	assert.Equal(t, "-2E+10", ValToString(-2e10))
	assert.Equal(t, "2E-10", ValToString(2e-10))
	assert.Equal(t, "1.2345E+100", ValToString(1.2345e100))
	assert.Equal(t, "1.2345E-100", ValToString(1.2345e-100))
	assert.Equal(t, "true", ValToString(true))
	assert.Equal(t, "false", ValToString(false))
	ts, _ := time.Parse(time.RFC3339, "2025-01-02T03:04:05Z")
	assert.Equal(t, "2025-01-02T03:04:05Z", ValToString(ts))
	assert.Equal(t, "a string", ValToString("a string"))
}

func TestAZKVKeyRoundTripPreservesPublicKey(t *testing.T) {
	internal, err := (&azkvkey{
		VaultURL:         "https://test.vault.azure.net",
		Name:             "test-key",
		Version:          "test-version",
		PublicKey:        base64.StdEncoding.EncodeToString([]byte("public-key")),
		CreatedAt:        "2025-01-02T03:04:05Z",
		EncryptedDataKey: "ciphertext",
	}).toInternal()
	assert.NoError(t, err)
	assert.Equal(t, []byte("public-key"), internal.PublicKey)

	keys := azkvKeysFromGroup(sops.KeyGroup{&azkv.MasterKey{
		VaultURL:     "https://test.vault.azure.net",
		Name:         "test-key",
		Version:      "test-version",
		PublicKey:    []byte("public-key"),
		CreationDate: time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC),
		EncryptedKey: "ciphertext",
	}})
	assert.Len(t, keys, 1)
	assert.Equal(t, base64.StdEncoding.EncodeToString([]byte("public-key")), keys[0].PublicKey)
}
