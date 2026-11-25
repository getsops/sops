package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	publicKeyPath, err := configValue(request.Config, "public_key_file")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("read OpenPGP public key: %w", err)
	}
	recipients, err := openpgp.ReadKeyRing(bytes.NewReader(publicKeyData))
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("parse OpenPGP public key: %w", err)
	}
	var ciphertext bytes.Buffer
	armorWriter, err := armor.Encode(&ciphertext, "PGP MESSAGE", nil)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("create OpenPGP armor: %w", err)
	}
	writer, err := openpgp.Encrypt(armorWriter, recipients, nil, &openpgp.FileHints{IsBinary: true}, nil)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("create OpenPGP writer: %w", err)
	}
	if _, err := writer.Write(request.DataKey); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("encrypt OpenPGP data key: %w", err)
	}
	if err := writer.Close(); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("close OpenPGP writer: %w", err)
	}
	if err := armorWriter.Close(); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("close OpenPGP armor: %w", err)
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext.Bytes(),
		Metadata:         map[string]interface{}{"key_id": keyID},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	keyID, err := configValue(request.Config, "key_id")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("OpenPGP key metadata does not match configuration")
	}
	privateKeyPath, err := configValue(request.Config, "private_key_file")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("read OpenPGP private key: %w", err)
	}
	identities, err := openpgp.ReadKeyRing(bytes.NewReader(privateKeyData))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("parse OpenPGP private key: %w", err)
	}
	armored, err := armor.Decode(bytes.NewReader(request.EncryptedDataKey))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode OpenPGP armor: %w", err)
	}
	message, err := openpgp.ReadMessage(armored.Body, identities, nil, nil)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt OpenPGP data key: %w", err)
	}
	dataKey, err := io.ReadAll(message.UnverifiedBody)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("read OpenPGP data key: %w", err)
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
