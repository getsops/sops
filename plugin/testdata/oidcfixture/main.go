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

const clientSecretEnvironmentVariable = "SOPS_TEST_OIDC_CLIENT_SECRET"

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	config, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	token, err := exchangeToken(ctx, config.issuerEndpoint, config.subject)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	var response struct {
		Ciphertext string `json:"ciphertext_b64"`
		Version    string `json:"version"`
	}
	err = callJSON(ctx, http.MethodPost, config.kmsEndpoint+"/wrap", token, map[string]string{
		"key_id":       config.keyID,
		"data_key_b64": base64.StdEncoding.EncodeToString(request.DataKey),
	}, &response)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	if response.Version == "" {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("OIDC KMS returned an empty key version")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(response.Ciphertext)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("decode OIDC KMS ciphertext: %w", err)
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata:         map[string]interface{}{"key_id": config.keyID, "version": response.Version},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	config, err := providerConfig(request.Config)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != config.keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("OIDC KMS metadata does not match configuration")
	}
	version, ok := request.Metadata["version"].(string)
	if !ok || version == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing OIDC KMS version metadata")
	}
	token, err := exchangeToken(ctx, config.issuerEndpoint, config.subject)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	var response struct {
		DataKey string `json:"data_key_b64"`
	}
	err = callJSON(ctx, http.MethodPost, config.kmsEndpoint+"/unwrap", token, map[string]string{
		"key_id":         config.keyID,
		"version":        version,
		"ciphertext_b64": base64.StdEncoding.EncodeToString(request.EncryptedDataKey),
	}, &response)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	dataKey, err := base64.StdEncoding.DecodeString(response.DataKey)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode OIDC KMS data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

type providerSettings struct {
	issuerEndpoint string
	kmsEndpoint    string
	subject        string
	keyID          string
}

func providerConfig(config map[string]interface{}) (providerSettings, error) {
	var settings providerSettings
	var err error
	if settings.issuerEndpoint, err = configValue(config, "issuer_endpoint"); err != nil {
		return providerSettings{}, err
	}
	if settings.kmsEndpoint, err = configValue(config, "kms_endpoint"); err != nil {
		return providerSettings{}, err
	}
	if settings.subject, err = configValue(config, "subject"); err != nil {
		return providerSettings{}, err
	}
	if settings.keyID, err = configValue(config, "key_id"); err != nil {
		return providerSettings{}, err
	}
	settings.issuerEndpoint = strings.TrimRight(settings.issuerEndpoint, "/")
	settings.kmsEndpoint = strings.TrimRight(settings.kmsEndpoint, "/")
	return settings, nil
}

func exchangeToken(ctx context.Context, endpoint, subject string) (string, error) {
	var response struct {
		Token string `json:"access_token"`
	}
	err := callJSON(ctx, http.MethodPost, endpoint+"/token", "", map[string]string{
		"subject":       subject,
		"client_secret": os.Getenv(clientSecretEnvironmentVariable),
	}, &response)
	if err != nil {
		return "", fmt.Errorf("exchange OIDC credentials: %w", err)
	}
	if response.Token == "" {
		return "", fmt.Errorf("OIDC issuer returned an empty access token")
	}
	return response.Token, nil
}

func callJSON(ctx context.Context, method, endpoint, bearerToken string, requestBody interface{}, responseBody interface{}) error {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("encode OIDC KMS request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create OIDC KMS request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	if bearerToken != "" {
		request.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("call OIDC-backed KMS: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return fmt.Errorf("read OIDC KMS response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("OIDC KMS returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	if err := json.Unmarshal(responseData, responseBody); err != nil {
		return fmt.Errorf("decode OIDC KMS response: %w", err)
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

func main() {
	sopsplugin.Serve(provider{})
}
