package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

const tokenEnvironmentVariable = "SOPS_TEST_AGENT_TOKEN"

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	keyID, ok := request.Config["key_id"].(string)
	if !ok || keyID == "" {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("missing string plugin config %q", "key_id")
	}
	publicKeyPEM, ok := request.Config["public_key"].(string)
	if !ok || publicKeyPEM == "" {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("missing string plugin config %q", "public_key")
	}
	publicKey, err := parsePublicKey([]byte(publicKeyPEM))
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, request.DataKey, nil)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("wrap data key: %w", err)
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata:         map[string]interface{}{"key_id": keyID},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	keyID, ok := request.Config["key_id"].(string)
	if !ok || keyID == "" || request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("provider key metadata does not match configuration")
	}
	endpoint, ok := request.Config["agent_endpoint"].(string)
	if !ok || endpoint == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing string plugin config %q", "agent_endpoint")
	}
	body, err := json.Marshal(map[string]string{
		"ciphertext_b64": base64.StdEncoding.EncodeToString(request.EncryptedDataKey),
	})
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("encode agent request: %w", err)
	}
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(endpoint, "/")+"/decrypt", bytes.NewReader(body))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("create agent request: %w", err)
	}
	httpRequest.Header.Set("Content-Type", "application/json")
	httpRequest.Header.Set("Authorization", "Bearer "+os.Getenv(tokenEnvironmentVariable))
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(httpRequest)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("call private-key agent: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("read agent response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("private-key agent returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	var result struct {
		DataKey string `json:"data_key_b64"`
	}
	if err := json.Unmarshal(responseData, &result); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode agent response: %w", err)
	}
	dataKey, err := base64.StdEncoding.DecodeString(result.DataKey)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode unwrapped data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func parsePublicKey(data []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM public key found")
	}
	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := parsed.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA")
	}
	return publicKey, nil
}

func main() {
	sopsplugin.Serve(provider{})
}
