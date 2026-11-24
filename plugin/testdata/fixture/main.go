package main

import (
	"context"
	"fmt"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	label, ok := request.Config["label"].(string)
	if !ok {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("missing label")
	}
	wrapped := make([]byte, len(request.DataKey))
	for i, value := range request.DataKey {
		wrapped[i] = value ^ 0xa7
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: wrapped,
		Metadata:         map[string]interface{}{"label": label},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	label, ok := request.Config["label"].(string)
	if !ok || request.Metadata["label"] != label {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("provider metadata does not match config")
	}
	dataKey := make([]byte, len(request.EncryptedDataKey))
	for i, value := range request.EncryptedDataKey {
		dataKey[i] = value ^ 0xa7
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func main() {
	sopsplugin.Serve(provider{})
}
