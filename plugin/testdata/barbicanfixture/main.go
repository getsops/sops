package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

const tokenEnvironmentVariable = "SOPS_TEST_BARBICAN_TOKEN"

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	endpoint, secretID, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	secret, version, err := getSecret(ctx, endpoint, secretID, "current")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	gcm, err := newGCM(secret)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("generate local encryption nonce: %w", err)
	}
	ciphertext := append(nonce, gcm.Seal(nil, nonce, request.DataKey, []byte(secretID))...)
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata: map[string]interface{}{
			"secret_id":      secretID,
			"secret_version": version,
		},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	endpoint, secretID, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["secret_id"] != secretID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("provider secret metadata does not match configuration")
	}
	version, ok := request.Metadata["secret_version"].(string)
	if !ok || version == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing string provider metadata %q", "secret_version")
	}
	secret, _, err := getSecret(ctx, endpoint, secretID, version)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	gcm, err := newGCM(secret)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if len(request.EncryptedDataKey) < gcm.NonceSize() {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("wrapped data key is shorter than the nonce")
	}
	nonce := request.EncryptedDataKey[:gcm.NonceSize()]
	ciphertext := request.EncryptedDataKey[gcm.NonceSize():]
	dataKey, err := gcm.Open(nil, nonce, ciphertext, []byte(secretID))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt data key with retrieved secret: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func providerConfig(config map[string]interface{}) (string, string, error) {
	endpoint, endpointOK := config["endpoint"].(string)
	secretID, secretIDOK := config["secret_id"].(string)
	if !endpointOK || endpoint == "" {
		return "", "", fmt.Errorf("missing string plugin config %q", "endpoint")
	}
	if !secretIDOK || secretID == "" {
		return "", "", fmt.Errorf("missing string plugin config %q", "secret_id")
	}
	return strings.TrimRight(endpoint, "/"), secretID, nil
}

func getSecret(ctx context.Context, endpoint, secretID, version string) ([]byte, string, error) {
	requestURL, err := url.JoinPath(endpoint, "v1", "secrets", secretID)
	if err != nil {
		return nil, "", fmt.Errorf("build secret URL: %w", err)
	}
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, "", fmt.Errorf("parse secret URL: %w", err)
	}
	query := parsedURL.Query()
	query.Set("version", version)
	parsedURL.RawQuery = query.Encode()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return nil, "", fmt.Errorf("create secret request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+os.Getenv(tokenEnvironmentVariable))
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, "", fmt.Errorf("retrieve secret from Barbican-like service: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return nil, "", fmt.Errorf("read secret response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, "", fmt.Errorf("secret service returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	var result struct {
		Secret  string `json:"secret_b64"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(responseData, &result); err != nil {
		return nil, "", fmt.Errorf("decode secret response: %w", err)
	}
	secret, err := base64.StdEncoding.DecodeString(result.Secret)
	if err != nil {
		return nil, "", fmt.Errorf("decode retrieved secret: %w", err)
	}
	if result.Version == "" {
		return nil, "", fmt.Errorf("secret service returned an empty version")
	}
	return secret, result.Version, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create local cipher from retrieved secret: %w", err)
	}
	return cipher.NewGCM(block)
}

func main() {
	sopsplugin.Serve(provider{})
}
