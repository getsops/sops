package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type region struct {
	name     string
	endpoint string
}

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	configuredRegions := regions(request.Config)
	if len(configuredRegions) == 0 {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("no KMS regions configured")
	}
	var lastError error
	for _, region := range configuredRegions {
		var response struct {
			Ciphertext string `json:"ciphertext_b64"`
			Version    string `json:"version"`
		}
		err := call(ctx, region.endpoint+"/wrap", map[string]string{
			"key_id":       keyID,
			"data_key_b64": base64.StdEncoding.EncodeToString(request.DataKey),
		}, &response)
		if err != nil {
			lastError = fmt.Errorf("region %q: %w", region.name, err)
			continue
		}
		if response.Version == "" {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("region %q returned an empty key version", region.name)
		}
		ciphertext, err := base64.StdEncoding.DecodeString(response.Ciphertext)
		if err != nil {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("decode failover ciphertext: %w", err)
		}
		return sopsplugin.WrapKeyResponse{
			EncryptedDataKey: ciphertext,
			Metadata:         map[string]interface{}{"key_id": keyID, "region": region.name, "version": response.Version},
		}, nil
	}
	return sopsplugin.WrapKeyResponse{}, fmt.Errorf("all KMS regions failed: %w", lastError)
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("KMS key metadata does not match configuration")
	}
	regionName, ok := request.Metadata["region"].(string)
	if !ok {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing KMS region metadata")
	}
	endpoint := ""
	for _, region := range regions(request.Config) {
		if region.name == regionName {
			endpoint = region.endpoint
		}
	}
	if endpoint == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("stored KMS region %q is not configured", regionName)
	}
	version, ok := request.Metadata["version"].(string)
	if !ok || version == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing KMS version metadata")
	}
	var response struct {
		DataKey string `json:"data_key_b64"`
	}
	err = call(ctx, endpoint+"/unwrap", map[string]string{
		"key_id":         keyID,
		"version":        version,
		"ciphertext_b64": base64.StdEncoding.EncodeToString(request.EncryptedDataKey),
	}, &response)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("unwrap in region %q: %w", regionName, err)
	}
	dataKey, err := base64.StdEncoding.DecodeString(response.DataKey)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode failover data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func regions(config map[string]interface{}) []region {
	var result []region
	for _, name := range []string{"primary", "secondary"} {
		endpoint, ok := config[name+"_endpoint"].(string)
		if ok && endpoint != "" {
			result = append(result, region{name: name, endpoint: strings.TrimRight(endpoint, "/")})
		}
	}
	return result
}

func configValue(config map[string]interface{}, name string) (string, error) {
	value, ok := config[name].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("missing string plugin config %q", name)
	}
	return value, nil
}

func call(ctx context.Context, endpoint string, requestBody interface{}, responseBody interface{}) error {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encode KMS request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create KMS request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("call KMS: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return fmt.Errorf("read KMS response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("KMS returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	if err := json.Unmarshal(responseData, responseBody); err != nil {
		return fmt.Errorf("decode KMS response: %w", err)
	}
	return nil
}

func main() {
	sopsplugin.Serve(provider{})
}
