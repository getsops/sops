package main

import (
	"bytes"
	"context"
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

const tokenEnvironmentVariable = "SOPS_TEST_SECRET_REF_TOKEN"

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	endpoint, namespace, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	body, err := json.Marshal(map[string]string{
		"namespace":    namespace,
		"data_key_b64": base64.StdEncoding.EncodeToString(request.DataKey),
	})
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("encode secret request: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/v1/secrets", bytes.NewReader(body))
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("create secret request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+os.Getenv(tokenEnvironmentVariable))
	responseData, err := call(httpRequest)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	var response struct {
		SecretID string `json:"secret_id"`
		Version  string `json:"version"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("decode secret reference response: %w", err)
	}
	if response.SecretID == "" || response.Version == "" {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("secret service returned an empty reference or version")
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: []byte(response.SecretID),
		Metadata: map[string]interface{}{
			"namespace": namespace,
			"secret_id": response.SecretID,
			"version":   response.Version,
		},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	endpoint, namespace, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	secretID, ok := request.Metadata["secret_id"].(string)
	if !ok || secretID == "" || string(request.EncryptedDataKey) != secretID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("encrypted data key does not match stored secret reference")
	}
	if request.Metadata["namespace"] != namespace {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("secret namespace metadata does not match configuration")
	}
	version, ok := request.Metadata["version"].(string)
	if !ok || version == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing string provider metadata %q", "version")
	}

	requestURL, err := url.JoinPath(endpoint, "v1", "secrets", secretID)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("build secret URL: %w", err)
	}
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("parse secret URL: %w", err)
	}
	query := parsedURL.Query()
	query.Set("version", version)
	query.Set("namespace", namespace)
	parsedURL.RawQuery = query.Encode()
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("create secret retrieval request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+os.Getenv(tokenEnvironmentVariable))
	responseData, err := call(httpRequest)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	var response struct {
		DataKey string `json:"data_key_b64"`
	}
	if err := json.Unmarshal(responseData, &response); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode secret response: %w", err)
	}
	dataKey, err := base64.StdEncoding.DecodeString(response.DataKey)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode stored data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func providerConfig(config map[string]interface{}) (string, string, error) {
	endpoint, endpointOK := config["endpoint"].(string)
	namespace, namespaceOK := config["namespace"].(string)
	if !endpointOK || endpoint == "" {
		return "", "", fmt.Errorf("missing string plugin config %q", "endpoint")
	}
	if !namespaceOK || namespace == "" {
		return "", "", fmt.Errorf("missing string plugin config %q", "namespace")
	}
	return strings.TrimRight(endpoint, "/"), namespace, nil
}

func call(request *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call secret service: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return nil, fmt.Errorf("read secret service response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("secret service returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	return responseData, nil
}

func main() {
	sopsplugin.Serve(provider{})
}
