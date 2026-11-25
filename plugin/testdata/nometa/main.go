package main

import (
	"context"
	"fmt"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	wrapped := append([]byte("nometa:"), request.DataKey...)
	return sopsplugin.WrapKeyResponse{EncryptedDataKey: wrapped}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	const prefix = "nometa:"
	if len(request.EncryptedDataKey) < len(prefix) || string(request.EncryptedDataKey[:len(prefix)]) != prefix {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("invalid wrapped data key")
	}
	return sopsplugin.UnwrapKeyResponse{
		DataKey: append([]byte(nil), request.EncryptedDataKey[len(prefix):]...),
	}, nil
}

func main() {
	sopsplugin.Serve(provider{})
}
