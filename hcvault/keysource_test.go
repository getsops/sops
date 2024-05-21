package hcvault

import (
	"fmt"
	logger "log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/mitchellh/go-homedir"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/assert"
)

var (
	// testVaultVersion is the version (image tag) of the Vault server image
	// used to test against.
	testVaultVersion = "1.10.0"
	// testVaultToken is the token of the Vault server.
	testVaultToken = "secret"
	// testEnginePath is the path to mount the Vault Transit on.
	testEnginePath = "sops"
	// testVaultAddress is the HTTP/S address of the Vault server, it is set
	// by TestMain after booting it.
	testVaultAddress string
)

// TestMain initializes a Vault server using Docker, writes the HTTP address to
// testVaultAddress, waits for it to become ready to serve requests, and enables
// Vault Transit on the testEnginePath. It then runs all the tests, which can
// make use of the various `test*` variables.
func TestMain(m *testing.M) {
	// Uses a sensible default on Windows (TCP/HTTP) and Linux/MacOS (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Fatalf("could not connect to docker: %s", err)
	}

	// Pull the image, create a container based on it, and run it
	resource, err := pool.Run("vault", testVaultVersion, []string{"VAULT_DEV_ROOT_TOKEN_ID=" + testVaultToken})
	if err != nil {
		logger.Fatalf("could not start resource: %s", err)
	}

	purgeResource := func() {
		if err := pool.Purge(resource); err != nil {
			logger.Printf("could not purge resource: %s", err)
		}
	}

	testVaultAddress = fmt.Sprintf("http://127.0.0.1:%v", resource.GetPort("8200/tcp"))
	// Wait until Vault is ready to serve requests
	if err := pool.Retry(func() error {
		cfg := api.DefaultConfig()
		cfg.Address = testVaultAddress
		cli, err := api.NewClient(cfg)
		if err != nil {
			return fmt.Errorf("cannot create Vault client: %w", err)
		}
		status, err := cli.Sys().InitStatus()
		if err != nil {
			return err
		}
		if status != true {
			return fmt.Errorf("waiting on Vault server to become ready")
		}
		return nil
	}); err != nil {
		purgeResource()
		logger.Fatalf("could not connect to docker: %s", err)
	}

	if err = enableVaultTransit(testVaultAddress, testVaultToken, testEnginePath); err != nil {
		purgeResource()
		logger.Fatalf("could not enable Vault transit: %s", err)
	}

	// Run the tests, but only if we succeeded in setting up the Vault server
	var code int
	if err == nil {
		code = m.Run()
	}

	// This can't be deferred, as os.Exit simply does not care
	if err := pool.Purge(resource); err != nil {
		logger.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestNewMasterKeysFromURIs(t *testing.T) {
	t.Run("multiple URIs", func(t *testing.T) {
		uris := []string{
			"https://vault.example.com:8200/v1/transit/keys/keyName",
			"", // Empty should be skipped
			"https://vault.me.com/v1/super42/bestmarket/keys/slig",
		}
		keys, err := NewMasterKeysFromURIs(strings.Join(uris, ","))
		assert.NoError(t, err)
		assert.Len(t, keys, 2)
	})

	t.Run("with invalid URI", func(t *testing.T) {
		uris := []string{
			"https://vault.example.com:8200/v1/transit/keys/keyName",
			"vault.me/keys/dev/mykey",
		}
		keys, err := NewMasterKeysFromURIs(strings.Join(uris, ","))
		assert.Error(t, err)
		assert.Nil(t, keys)
	})
}

func TestNewMasterKeyFromURI(t *testing.T) {
	tests := []struct {
		url     string
		want    *MasterKey
		wantErr bool
	}{
		{
			url: "https://vault.example.com:8200/v1/transit/keys/keyName",
			want: &MasterKey{
				VaultAddress: "https://vault.example.com:8200",
				EnginePath:   "transit",
				KeyName:      "keyName",
			},
		},
		{
			url: "https://vault.me.com/v1/super42/bestmarket/keys/slig",
			want: &MasterKey{
				VaultAddress: "https://vault.me.com",
				EnginePath:   "super42/bestmarket",
				KeyName:      "slig",
			},
		},
		{
			url: "http://127.0.0.1:12121/v1/transit/keys/dev",
			want: &MasterKey{
				VaultAddress: "http://127.0.0.1:12121",
				EnginePath:   "transit",
				KeyName:      "dev",
			},
		},
		{
			url:     "vault.me/keys/dev/mykey",
			want:    nil,
			wantErr: true,
		},
		{
			url:     "http://127.0.0.1:12121/v1/keys/dev",
			want:    nil,
			wantErr: true,
		},
		{
			url:     "tcp://127.0.0.1:12121/v1/keys/dev",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			got, err := NewMasterKeyFromURI(tt.url)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if tt.want != nil && got != nil {
				tt.want.CreationDate = got.CreationDate
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestMasterKey_Encrypt(t *testing.T) {
	key := NewMasterKey(testVaultAddress, testEnginePath, "encrypt")
	(Token(testVaultToken)).ApplyToMasterKey(key)
	assert.NoError(t, createVaultKey(key))

	dataKey := []byte("the majority of your brain is fat")
	assert.NoError(t, key.Encrypt(dataKey))
	assert.NotEmpty(t, key.EncryptedKey)

	client, err := vaultClient(key.VaultAddress, key.token)
	assert.NoError(t, err)

	payload := decryptPayload(key.EncryptedKey)
	secret, err := client.Logical().Write(key.decryptPath(), payload)
	assert.NoError(t, err)

	decryptedData, err := dataKeyFromSecret(secret)
	assert.NoError(t, err)
	assert.Equal(t, dataKey, decryptedData)

	key.EnginePath = "invalid"
	assert.Error(t, key.Encrypt(dataKey))

	key.EnginePath = testEnginePath
	key.token = ""
	assert.Error(t, key.Encrypt(dataKey))
}

func TestMasterKey_EncryptIfNeeded(t *testing.T) {
	key := NewMasterKey(testVaultAddress, testEnginePath, "encrypt-if-needed")
	(Token(testVaultToken)).ApplyToMasterKey(key)
	assert.NoError(t, createVaultKey(key))

	assert.NoError(t, key.EncryptIfNeeded([]byte("stingy string")))

	encryptedKey := key.EncryptedKey
	assert.NotEmpty(t, encryptedKey)

	assert.NoError(t, key.EncryptIfNeeded([]byte("stringy sting")))
	assert.Equal(t, encryptedKey, key.EncryptedKey)
}

func TestMasterKey_EncryptedDataKey(t *testing.T) {
	key := &MasterKey{EncryptedKey: "some key"}
	assert.EqualValues(t, key.EncryptedKey, key.EncryptedDataKey())
}

func TestMasterKey_Decrypt(t *testing.T) {
	key := NewMasterKey(testVaultAddress, testEnginePath, "decrypt")
	(Token(testVaultToken)).ApplyToMasterKey(key)
	assert.NoError(t, createVaultKey(key))

	client, err := vaultClient(key.VaultAddress, key.token)
	assert.NoError(t, err)

	dataKey := []byte("the heart of a shrimp is located in its head")
	secret, err := client.Logical().Write(key.encryptPath(), encryptPayload(dataKey))
	assert.NoError(t, err)

	encryptedKey, err := encryptedKeyFromSecret(secret)
	assert.NoError(t, err)

	key.EncryptedKey = encryptedKey
	got, err := key.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, dataKey, got)

	key.EnginePath = "invalid"
	assert.Error(t, key.Encrypt(dataKey))

	key.EnginePath = testEnginePath
	key.token = ""
	assert.Error(t, key.Encrypt(dataKey))
}

func TestMasterKey_EncryptDecrypt_RoundTrip(t *testing.T) {
	token := Token(testVaultToken)

	encryptKey := NewMasterKey(testVaultAddress, testEnginePath, "roundtrip")
	token.ApplyToMasterKey(encryptKey)
	assert.NoError(t, createVaultKey(encryptKey))

	dataKey := []byte("some people have an extra bone in their knee")
	assert.NoError(t, encryptKey.Encrypt(dataKey))
	assert.NotEmpty(t, encryptKey.EncryptedKey)

	decryptKey := NewMasterKey(testVaultAddress, testEnginePath, "roundtrip")
	token.ApplyToMasterKey(decryptKey)
	decryptKey.EncryptedKey = encryptKey.EncryptedKey

	decryptedData, err := decryptKey.Decrypt()
	assert.NoError(t, err)
	assert.Equal(t, dataKey, decryptedData)
}

func TestMasterKey_NeedsRotation(t *testing.T) {
	key := NewMasterKey("", "", "")
	assert.False(t, key.NeedsRotation())

	key.CreationDate = key.CreationDate.Add(-(vaultTTL + time.Second))
	assert.True(t, key.NeedsRotation())
}

func TestMasterKey_ToString(t *testing.T) {
	key := NewMasterKey("https://example.com", "engine", "key-name")
	assert.Equal(t, "https://example.com/v1/engine/keys/key-name", key.ToString())
}

func TestMasterKey_ToMap(t *testing.T) {
	key := &MasterKey{
		KeyName:      "test-key",
		EnginePath:   "engine",
		VaultAddress: testVaultAddress,
		EncryptedKey: "some-encrypted-key",
	}
	assert.Equal(t, map[string]interface{}{
		"vault_address": key.VaultAddress,
		"key_name":      key.KeyName,
		"engine_path":   key.EnginePath,
		"enc":           key.EncryptedKey,
		"created_at":    "0001-01-01T00:00:00Z",
	}, key.ToMap())
}

func Test_encryptedKeyFromSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  *api.Secret
		want    string
		wantErr bool
	}{
		{name: "nil secret", secret: nil, wantErr: true},
		{name: "secret with nil data", secret: &api.Secret{Data: nil}, wantErr: true},
		{name: "secret without ciphertext data", secret: &api.Secret{Data: map[string]interface{}{"other": true}}, wantErr: true},
		{name: "ciphertext non string", secret: &api.Secret{Data: map[string]interface{}{"ciphertext": 123}}, wantErr: true},
		{name: "ciphertext data", secret: &api.Secret{Data: map[string]interface{}{"ciphertext": "secret string"}}, want: "secret string"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encryptedKeyFromSecret(tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_dataKeyFromSecret(t *testing.T) {
	tests := []struct {
		name    string
		secret  *api.Secret
		want    []byte
		wantErr bool
	}{
		{name: "nil secret", secret: nil, wantErr: true},
		{name: "secret with nil data", secret: &api.Secret{Data: nil}, wantErr: true},
		{name: "secret without plaintext data", secret: &api.Secret{Data: map[string]interface{}{"other": true}}, wantErr: true},
		{name: "plaintext non string", secret: &api.Secret{Data: map[string]interface{}{"plaintext": 123}}, wantErr: true},
		{name: "plaintext non base64", secret: &api.Secret{Data: map[string]interface{}{"plaintext": "notbase64"}}, wantErr: true},
		{name: "plaintext base64 data", secret: &api.Secret{Data: map[string]interface{}{"plaintext": "Zm9v"}}, want: []byte("foo")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := dataKeyFromSecret(tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_vaultClient(t *testing.T) {
	t.Run("client", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Reset before and after to make sure the override is taken into
		// account, and restored after the test.
		homedir.Reset()
		t.Cleanup(func() { homedir.Reset() })
		t.Setenv("VAULT_TOKEN", "")
		t.Setenv("HOME", tmpDir)

		got, err := vaultClient(testVaultAddress, "")
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Empty(t, got.Token())
	})

	t.Run("client with VAULT_TOKEN", func(t *testing.T) {
		token := "test-token"
		t.Setenv("VAULT_TOKEN", token)

		got, err := vaultClient(testVaultAddress, "")
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, token, got.Token())
	})

	t.Run("client with token", func(t *testing.T) {
		ignored := "test-token"
		t.Setenv("VAULT_TOKEN", ignored)

		got, err := vaultClient(testVaultAddress, testVaultToken)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, testVaultToken, got.Token())
	})

	t.Run("client with token from file", func(t *testing.T) {
		tmpDir := t.TempDir()

		token := "test-token"
		assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, defaultTokenFile), []byte(token), 0600))

		// Reset before and after to make sure the override is taken into
		// account, and restored after the test.
		homedir.Reset()
		t.Cleanup(func() { homedir.Reset() })
		t.Setenv("VAULT_TOKEN", "")
		t.Setenv("HOME", tmpDir)

		got, err := vaultClient(testVaultAddress, "")
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, token, got.Token())
	})
}

func Test_userVaultToken(t *testing.T) {
	t.Run("reads token from file", func(t *testing.T) {
		tmpDir := t.TempDir()

		token := "test-token"
		assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, defaultTokenFile), []byte(token), 0600))

		// Reset before and after to make sure the override is taken into
		// account, and restored after the test.
		homedir.Reset()
		t.Cleanup(func() { homedir.Reset() })
		t.Setenv("HOME", tmpDir)

		got, err := userVaultToken()
		assert.NoError(t, err)
		assert.Equal(t, token, got)
	})

	t.Run("ignores missing file", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Reset before and after to make sure the override is taken into
		// account, and restored after the test.
		homedir.Reset()
		t.Cleanup(func() { homedir.Reset() })
		t.Setenv("HOME", tmpDir)

		got, err := userVaultToken()
		assert.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("trims spaces", func(t *testing.T) {
		tmpDir := t.TempDir()

		token := "  test-token  "
		assert.NoError(t, os.WriteFile(filepath.Join(tmpDir, defaultTokenFile), []byte(token), 0600))

		// Reset before and after to make sure the override is taken into
		// account, and restored after the test.
		homedir.Reset()
		t.Cleanup(func() { homedir.Reset() })
		t.Setenv("HOME", tmpDir)

		got, err := userVaultToken()
		assert.NoError(t, err)
		assert.Equal(t, "test-token", got)
	})
}

func Test_engineAndKeyFromPath(t *testing.T) {
	t.Run("engine and key", func(t *testing.T) {
		enginePath, key, err := engineAndKeyFromPath("/v1/transit/keys/keyName")
		assert.NoError(t, err)
		assert.Equal(t, "transit", enginePath)
		assert.Equal(t, "keyName", key)
	})

	t.Run("long (nested) path error", func(t *testing.T) {
		_, _, err := engineAndKeyFromPath("/nested/v1/transit/keys/bar")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "running Vault with a prefixed URL is not supported")
	})

	t.Run("invalid format error", func(t *testing.T) {
		_, _, err := engineAndKeyFromPath("/secret/foo/bar")
		assert.Error(t, err)
		assert.ErrorContains(t, err, "vault path does not seem to be formatted correctly")
	})
}

// enableVaultTransit enables the Vault Transit backend on the given enginePath.
func enableVaultTransit(address, token, enginePath string) error {
	client, err := vaultClient(address, token)
	if err != nil {
		return fmt.Errorf("cannot create Vault client: %w", err)
	}

	if err = client.Sys().Mount(enginePath, &api.MountInput{
		Type:        "transit",
		Description: "backend transit used by SOPS",
	}); err != nil {
		return fmt.Errorf("failed to mount transit on engine path '%s': %w", enginePath, err)
	}
	return nil
}

// createVaultKey creates a new RSA-4096 Vault key using the data from the
// provided MasterKey.
func createVaultKey(key *MasterKey) error {
	client, err := vaultClient(key.VaultAddress, key.token)
	if err != nil {
		return fmt.Errorf("cannot create Vault client: %w", err)
	}

	p := path.Join(key.EnginePath, "keys", key.KeyName)
	payload := make(map[string]interface{})
	payload["type"] = "rsa-4096"
	if _, err = client.Logical().Write(p, payload); err != nil {
		return err
	}

	_, err = client.Logical().Read(p)
	return err
}
