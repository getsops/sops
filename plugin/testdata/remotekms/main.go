package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

const tokenEnvironmentVariable = "SOPS_TEST_KMS_TOKEN"

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	endpoint, keyID, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}

	var response struct {
		Ciphertext string `json:"ciphertext"`
		Version    string `json:"version"`
	}
	err = call(ctx, endpoint+"/v1/wrap", map[string]string{
		"key_id":       keyID,
		"data_key_b64": base64.StdEncoding.EncodeToString(request.DataKey),
	}, &response)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(response.Ciphertext)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("decode remote KMS ciphertext: %w", err)
	}
	if response.Version == "" {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("remote KMS returned an empty key version")
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata: map[string]interface{}{
			"key_id":  keyID,
			"version": response.Version,
		},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	endpoint, keyID, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("provider key metadata does not match configuration")
	}
	version, ok := request.Metadata["version"].(string)
	if !ok || version == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing string provider metadata %q", "version")
	}

	var response struct {
		DataKey string `json:"data_key_b64"`
	}
	err = call(ctx, endpoint+"/v1/unwrap", map[string]string{
		"key_id":         keyID,
		"version":        version,
		"ciphertext_b64": base64.StdEncoding.EncodeToString(request.EncryptedDataKey),
	}, &response)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	dataKey, err := base64.StdEncoding.DecodeString(response.DataKey)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode remote KMS data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func providerConfig(config map[string]interface{}) (string, string, error) {
	endpoint, endpointOK := config["endpoint"].(string)
	keyID, keyIDOK := config["key_id"].(string)
	if !endpointOK || endpoint == "" {
		return "", "", fmt.Errorf("missing string plugin config %q", "endpoint")
	}
	if !keyIDOK || keyID == "" {
		return "", "", fmt.Errorf("missing string plugin config %q", "key_id")
	}
	return strings.TrimRight(endpoint, "/"), keyID, nil
}

func call(ctx context.Context, endpoint string, requestBody interface{}, responseBody interface{}) error {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encode remote KMS request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create remote KMS request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+os.Getenv(tokenEnvironmentVariable))

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("call remote KMS: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return fmt.Errorf("read remote KMS response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("remote KMS returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	if err := json.Unmarshal(responseData, responseBody); err != nil {
		return fmt.Errorf("decode remote KMS response: %w", err)
	}
	return nil
}

func main() {
	sopsplugin.Serve(provider{})
}
