package main

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	key, keyID, auditEndpoint, err := providerConfig(request.Config)
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
	ciphertext := append(nonce, gcm.Seal(nil, nonce, request.DataKey, []byte(keyID))...)
	if err := writeAuditEvent(ctx, auditEndpoint, keyID, "wrap"); err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata:         map[string]interface{}{"key_id": keyID},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	key, keyID, auditEndpoint, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("audit key metadata does not match configuration")
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
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt data key: %w", err)
	}
	if err := writeAuditEvent(ctx, auditEndpoint, keyID, "unwrap"); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func providerConfig(config map[string]interface{}) ([]byte, string, string, error) {
	keyPath, err := configValue(config, "key_file")
	if err != nil {
		return nil, "", "", err
	}
	key, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, "", "", fmt.Errorf("read audit fixture key: %w", err)
	}
	keyID, err := configValue(config, "key_id")
	if err != nil {
		return nil, "", "", err
	}
	endpoint, err := configValue(config, "audit_endpoint")
	if err != nil {
		return nil, "", "", err
	}
	return key, keyID, strings.TrimRight(endpoint, "/"), nil
}

func writeAuditEvent(ctx context.Context, endpoint, keyID, operation string) error {
	body, err := json.Marshal(map[string]string{"key_id": keyID, "operation": operation})
	if err != nil {
		return fmt.Errorf("encode audit event: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/events", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create audit event request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("send audit event: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 4096))
	if err != nil {
		return fmt.Errorf("read audit response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("audit service returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	return nil
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
		return nil, fmt.Errorf("load audit fixture AES key: %w", err)
	}
	return cipher.NewGCM(block)
}

func main() {
	sopsplugin.Serve(provider{})
}
