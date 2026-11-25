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

const (
	keyEnvironmentVariable   = "SOPS_TEST_APPROVAL_KEY"
	tokenEnvironmentVariable = "SOPS_TEST_APPROVAL_TOKEN"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	gcm, err := newGCM()
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
		Metadata:         map[string]interface{}{"key_id": keyID},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("approval key metadata does not match configuration")
	}
	endpoint, err := configValue(request.Config, "approval_endpoint")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if err := requestApproval(ctx, endpoint, keyID); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	gcm, err := newGCM()
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
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt approved data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func requestApproval(ctx context.Context, endpoint, keyID string) error {
	body, err := json.Marshal(map[string]string{"key_id": keyID, "operation": "unwrap"})
	if err != nil {
		return fmt.Errorf("encode approval request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(endpoint, "/")+"/approve", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create approval request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+os.Getenv(tokenEnvironmentVariable))
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request unwrap approval: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 4096))
	if err != nil {
		return fmt.Errorf("read approval response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("unwrap approval denied: HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
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

func newGCM() (cipher.AEAD, error) {
	key := os.Getenv(keyEnvironmentVariable)
	if key == "" {
		return nil, fmt.Errorf("missing approval fixture key in %s", keyEnvironmentVariable)
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("load approval fixture key: %w", err)
	}
	return cipher.NewGCM(block)
}

func main() {
	sopsplugin.Serve(provider{})
}
