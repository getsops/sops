package main

import (
	"context"
	"errors"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	switch mode(request.Config) {
	case "error":
		return sopsplugin.WrapKeyResponse{}, errors.New("fixture wrap failure")
	case "empty-wrap":
		return sopsplugin.WrapKeyResponse{}, nil
	case "invalid-metadata":
		return sopsplugin.WrapKeyResponse{
			EncryptedDataKey: []byte("wrapped"),
			Metadata:         map[string]interface{}{"unsupported": make(chan int)},
		}, nil
	default:
		return sopsplugin.WrapKeyResponse{
			EncryptedDataKey: []byte("wrapped"),
			Metadata:         map[string]interface{}{},
		}, nil
	}
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	switch mode(request.Config) {
	case "error":
		return sopsplugin.UnwrapKeyResponse{}, errors.New("fixture unwrap failure")
	case "empty-unwrap":
		return sopsplugin.UnwrapKeyResponse{}, nil
	case "short-unwrap":
		return sopsplugin.UnwrapKeyResponse{DataKey: []byte("short")}, nil
	default:
		return sopsplugin.UnwrapKeyResponse{DataKey: []byte("01234567890123456789012345678901")}, nil
	}
}

func mode(config map[string]interface{}) string {
	value, ok := config["mode"].(string)
	if !ok || value == "" {
		return "valid"
	}
	return value
}

func main() {
	sopsplugin.Serve(provider{})
}
