//go:build integration

package azkv

import (
	"context"
	"encoding/base64"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/assert"
)

// The following values should be created based on the instructions in:
// https://github.com/mozilla/sops#encrypting-using-azure-key-vault
//
// Additionally required permissions for auto key creation:
// - KeyManagementOperations/Get
// - KeyManagementOperations/Create
var (
	testVaultTenantID     = os.Getenv("SOPS_TEST_AZURE_TENANT_ID")
	testVaultClientID     = os.Getenv("SOPS_TEST_AZURE_CLIENT_ID")
	testVaultClientSecret = os.Getenv("SOPS_TEST_AZURE_CLIENT_SECRET")

	testVaultURL        = os.Getenv("SOPS_TEST_AZURE_VAULT_URL")
	testVaultKeyName    = os.Getenv("SOPS_TEST_AZURE_VAULT_KEY_NAME")
	testVaultKeyVersion = os.Getenv("SOPS_TEST_AZURE_VAULT_KEY_VERSION")
)

func TestMasterKey_Encrypt(t *testing.T) {
	key, err := createTestKMSKeyIfNotExists()
	assert.NoError(t, err)

	data := []byte("to be or not to be static bytes")
	assert.NoError(t, key.Encrypt(data))
	assert.NotEmpty(t, key.EncryptedDataKey())
	assert.NotEqual(t, data, key.EncryptedKey)
}

func TestMasterKey_Decrypt(t *testing.T) {
	key, err := createTestKMSKeyIfNotExists()
	assert.NoError(t, err)

	data := []byte("this is super secret data")

	c, err := azkeys.NewClient(key.VaultURL, key.tokenCredential, nil)
	assert.NoError(t, err)

	resp, err := c.Encrypt(context.Background(), key.Name, key.Version, azkeys.KeyOperationParameters{
		Algorithm: to.Ptr(azkeys.EncryptionAlgorithmRSAOAEP256),
		Value:     data,
	}, nil)
	assert.NoError(t, err)
	key.EncryptedKey = base64.RawURLEncoding.EncodeToString(resp.KeyOperationResult.Result)

	got, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, data, got)
}

func TestMasterKey_EncryptDecrypt_RoundTrip(t *testing.T) {
	key, err := createTestKMSKeyIfNotExists()
	assert.NoError(t, err)

	data := []byte("the earth is round")
	assert.NoError(t, key.Encrypt(data))
	assert.NotNil(t, key.EncryptedDataKey())

	got, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, data, got)
}

func createTestKMSKeyIfNotExists() (*MasterKey, error) {
	token, err := testTokenCredential()
	if err != nil {
		return nil, err
	}

	key := &MasterKey{
		VaultURL: testVaultURL,
		Name:     testVaultKeyName,
		Version:  testVaultKeyVersion,
	}
	NewTokenCredential(token).ApplyToMasterKey(key)

	// If we have been given a version, assume it exists.
	if key.Version == "" {
		c, err := azkeys.NewClient(key.VaultURL, token, nil)
		if err != nil {
			return nil, err
		}

		getResp, err := c.GetKey(context.TODO(), key.Name, key.Version, nil)
		if err == nil {
			key.Version = getResp.KeyBundle.Key.KID.Version()
		}
		if err != nil {
			createResp, err := c.CreateKey(context.TODO(), key.Name, azkeys.CreateKeyParameters{
				Kty:    to.Ptr(azkeys.KeyTypeRSA),
				KeyOps: to.SliceOfPtrs(azkeys.KeyOperationEncrypt, azkeys.KeyOperationDecrypt),
			}, nil)
			if err != nil {
				return nil, err
			}
			key.Version = createResp.Key.KID.Version()
		}
	}

	return key, nil
}

func testTokenCredential() (azcore.TokenCredential, error) {
	return azidentity.NewClientSecretCredential(testVaultTenantID, testVaultClientID, testVaultClientSecret, nil)
}
