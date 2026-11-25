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
	"github.com/getsops/sops/v3/shamir"
)

const threshold = 2

type providerConfig struct {
	ShareEndpoints []string `json:"share_endpoints"`
	Threshold      int      `json:"threshold"`
	Generation     uint64   `json:"generation"`
}

type providerMetadata struct {
	Threshold   int    `json:"threshold"`
	TotalShares int    `json:"total_shares"`
	Generation  uint64 `json:"generation"`
}

type shareReference struct {
	Endpoint string `json:"endpoint"`
	ID       string `json:"id"`
}

type provider struct{}

func (provider) Wrap(ctx context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	var config providerConfig
	if err := sopsplugin.DecodeConfig(request.Config, &config); err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	if config.Threshold != threshold {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("share threshold must be %d", threshold)
	}
	if config.Generation == 0 {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("configuration generation must be non-zero")
	}
	if len(config.ShareEndpoints) < config.Threshold {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("at least %d share endpoints are required", config.Threshold)
	}
	endpoints, err := normalizedEndpoints(config.ShareEndpoints)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	shares, err := shamir.Split(request.DataKey, len(endpoints), config.Threshold)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("split data key: %w", err)
	}
	references := make([]shareReference, 0, len(endpoints))
	for index, endpoint := range endpoints {
		var response struct {
			ID string `json:"id"`
		}
		err := call(ctx, http.MethodPost, endpoint+"/shares", map[string]string{
			"share_b64": base64.StdEncoding.EncodeToString(shares[index]),
		}, &response)
		if err != nil {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("store share %d: %w", index+1, err)
		}
		if response.ID == "" {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("share service %d returned an empty ID", index+1)
		}
		references = append(references, shareReference{Endpoint: endpoint, ID: response.ID})
	}
	encryptedDataKey, err := json.Marshal(references)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("encode share references: %w", err)
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: encryptedDataKey,
		Metadata: map[string]interface{}{
			"threshold":    config.Threshold,
			"total_shares": len(endpoints),
			"generation":   config.Generation,
		},
	}, nil
}

func (provider) Unwrap(ctx context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	var config providerConfig
	if err := sopsplugin.DecodeConfig(request.Config, &config); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	var metadata providerMetadata
	if err := sopsplugin.DecodeMetadata(request.Metadata, &metadata); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	var references []shareReference
	if err := json.Unmarshal(request.EncryptedDataKey, &references); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode share references: %w", err)
	}
	if metadata.Threshold != config.Threshold || metadata.Generation != config.Generation {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("share metadata does not match provider configuration")
	}
	if metadata.TotalShares != len(references) || len(references) < metadata.Threshold {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("stored data key has invalid share references")
	}
	shares := make([][]byte, 0, metadata.Threshold)
	var lastError error
	for _, reference := range references {
		var response struct {
			Share string `json:"share_b64"`
		}
		err := call(ctx, http.MethodGet, reference.Endpoint+"/shares/"+reference.ID, nil, &response)
		if err != nil {
			lastError = err
			continue
		}
		share, err := base64.StdEncoding.DecodeString(response.Share)
		if err != nil {
			lastError = fmt.Errorf("decode share: %w", err)
			continue
		}
		shares = append(shares, share)
		if len(shares) == metadata.Threshold {
			break
		}
	}
	if len(shares) < metadata.Threshold {
		if lastError == nil {
			lastError = fmt.Errorf("no share references were available")
		}
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("need %d shares, retrieved %d: %w", metadata.Threshold, len(shares), lastError)
	}
	dataKey, err := shamir.Combine(shares)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("combine shares: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func normalizedEndpoints(endpoints []string) ([]string, error) {
	normalized := make([]string, len(endpoints))
	for index, endpoint := range endpoints {
		normalized[index] = strings.TrimRight(endpoint, "/")
		if normalized[index] == "" {
			return nil, fmt.Errorf("share endpoint %d must be non-empty", index)
		}
	}
	return normalized, nil
}

func call(ctx context.Context, method, endpoint string, requestBody interface{}, responseBody interface{}) error {
	var body []byte
	var err error
	if requestBody != nil {
		body, err = json.Marshal(requestBody)
		if err != nil {
			return fmt.Errorf("encode share request: %w", err)
		}
	}
	request, err := http.NewRequestWithContext(ctx, method, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create share request: %w", err)
	}
	if requestBody != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{Timeout: 5 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("call share service: %w", err)
	}
	defer response.Body.Close()
	responseData, err := io.ReadAll(io.LimitReader(response.Body, 64*1024))
	if err != nil {
		return fmt.Errorf("read share response: %w", err)
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("share service returned HTTP %s: %s", response.Status, strings.TrimSpace(string(responseData)))
	}
	if err := json.Unmarshal(responseData, responseBody); err != nil {
		return fmt.Errorf("decode share response: %w", err)
	}
	return nil
}

func main() {
	sopsplugin.Serve(provider{})
}
