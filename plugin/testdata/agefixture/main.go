package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	recipientText, err := configValue(request.Config, "recipient")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	recipient, err := age.ParseX25519Recipient(recipientText)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("parse age recipient: %w", err)
	}
	var ciphertext bytes.Buffer
	writer, err := age.Encrypt(&ciphertext, recipient)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("create age writer: %w", err)
	}
	if _, err := writer.Write(request.DataKey); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("encrypt data key: %w", err)
	}
	if err := writer.Close(); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("close age writer: %w", err)
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext.Bytes(),
		Metadata:         map[string]interface{}{"recipient": recipientText},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	identityPath, err := configValue(request.Config, "identity_file")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	identityText, err := os.ReadFile(identityPath)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("read age identity: %w", err)
	}
	identity, err := age.ParseX25519Identity(string(bytes.TrimSpace(identityText)))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("parse age identity: %w", err)
	}
	if request.Metadata["recipient"] != identity.Recipient().String() {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("age recipient metadata does not match identity")
	}
	reader, err := age.Decrypt(bytes.NewReader(request.EncryptedDataKey), identity)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt age data key: %w", err)
	}
	dataKey, err := io.ReadAll(reader)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("read age data key: %w", err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
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
