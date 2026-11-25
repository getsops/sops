package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

const masterSecretEnvironmentVariable = "SOPS_TEST_TENANT_MASTER_SECRET"

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	tenantID, keyID, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	key, err := deriveKey(tenantID, keyID)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	gcm, err := newGCM(key)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("generate nonce: %w", err)
	}
	aad := []byte(tenantID + "\x00" + keyID)
	ciphertext := append(nonce, gcm.Seal(nil, nonce, request.DataKey, aad)...)
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata:         map[string]interface{}{"tenant_id": tenantID, "key_id": keyID},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	tenantID, keyID, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["tenant_id"] != tenantID || request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("tenant key metadata does not match configuration")
	}
	key, err := deriveKey(tenantID, keyID)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
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
	aad := []byte(tenantID + "\x00" + keyID)
	dataKey, err := gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt tenant data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func providerConfig(config map[string]interface{}) (string, string, error) {
	tenantID, err := configValue(config, "tenant_id")
	if err != nil {
		return "", "", err
	}
	keyID, err := configValue(config, "key_id")
	if err != nil {
		return "", "", err
	}
	return tenantID, keyID, nil
}

func deriveKey(tenantID, keyID string) ([]byte, error) {
	secret := os.Getenv(masterSecretEnvironmentVariable)
	if secret == "" {
		return nil, fmt.Errorf("missing tenant master secret in %s", masterSecretEnvironmentVariable)
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(tenantID + "\x00" + keyID))
	return mac.Sum(nil), nil
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
		return nil, fmt.Errorf("derive tenant AES cipher: %w", err)
	}
	return cipher.NewGCM(block)
}

func main() {
	sopsplugin.Serve(provider{})
}
