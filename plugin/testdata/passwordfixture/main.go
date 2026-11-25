package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"

	sopsplugin "github.com/getsops/sops/v3/plugin"
	"golang.org/x/crypto/scrypt"
)

const passwordEnvironmentVariable = "SOPS_TEST_PLUGIN_PASSWORD"

// Keep the test fixture's work factor low; production providers should tune it
// for their supported runtime and security requirements.
const scryptWorkFactor = 1 << 12

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	password, err := passwordFromEnvironment()
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("generate password salt: %w", err)
	}
	key, err := scrypt.Key([]byte(password), salt, scryptWorkFactor, 8, 1, 32)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("derive wrapping key: %w", err)
	}
	gcm, err := newGCM(key)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := append(nonce, gcm.Seal(nil, nonce, request.DataKey, []byte(keyID))...)
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata: map[string]interface{}{
			"key_id": keyID,
			"salt":   base64.StdEncoding.EncodeToString(salt),
			"kdf":    "scrypt",
		},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	password, err := passwordFromEnvironment()
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != keyID || request.Metadata["kdf"] != "scrypt" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("password key metadata does not match configuration")
	}
	saltText, ok := request.Metadata["salt"].(string)
	if !ok || saltText == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing password salt metadata")
	}
	salt, err := base64.StdEncoding.DecodeString(saltText)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode password salt: %w", err)
	}
	key, err := scrypt.Key([]byte(password), salt, scryptWorkFactor, 8, 1, 32)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("derive wrapping key: %w", err)
	}
	gcm, err := newGCM(key)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if len(request.EncryptedDataKey) < gcm.NonceSize() {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("wrapped data key is shorter than the nonce")
	}
	nonce := request.EncryptedDataKey[:gcm.NonceSize()]
	ciphertext := request.EncryptedDataKey[gcm.NonceSize():]
	dataKey, err := gcm.Open(nil, nonce, ciphertext, []byte(keyID))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt with passphrase: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func passwordFromEnvironment() (string, error) {
	password := os.Getenv(passwordEnvironmentVariable)
	if password == "" {
		return "", fmt.Errorf("missing passphrase in %s", passwordEnvironmentVariable)
	}
	return password, nil
}

func configValue(config map[string]interface{}, name string) (string, error) {
	value, ok := config[name].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("missing string plugin config %q", name)
	}
	return value, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}
	return cipher.NewGCM(block)
}

func main() {
	sopsplugin.Serve(provider{})
}
