package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

type providerConfig struct {
	KeyringDir string `json:"keyring_dir"`
}

type providerMetadata struct {
	KeyVersion string `json:"key_version"`
}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	var config providerConfig
	if err := sopsplugin.DecodeConfig(request.Config, &config); err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	if config.KeyringDir == "" {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("missing string plugin config %q", "keyring_dir")
	}
	versionBytes, err := os.ReadFile(filepath.Join(config.KeyringDir, "current"))
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("read current key version: %w", err)
	}
	version := strings.TrimSpace(string(versionBytes))
	if err := validateVersion(version); err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	key, err := os.ReadFile(filepath.Join(config.KeyringDir, version+".key"))
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("read key version %q: %w", version, err)
	}
	gcm, err := newGCM(key)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := append(nonce, gcm.Seal(nil, nonce, request.DataKey, []byte(version))...)
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata:         map[string]interface{}{"key_version": version},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	var config providerConfig
	if err := sopsplugin.DecodeConfig(request.Config, &config); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if config.KeyringDir == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing string plugin config %q", "keyring_dir")
	}
	var metadata providerMetadata
	if err := sopsplugin.DecodeMetadata(request.Metadata, &metadata); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if metadata.KeyVersion == "" {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("missing key version metadata")
	}
	version := metadata.KeyVersion
	if err := validateVersion(version); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	key, err := os.ReadFile(filepath.Join(config.KeyringDir, version+".key"))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("read key version %q: %w", version, err)
	}
	gcm, err := newGCM(key)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	if len(request.EncryptedDataKey) < gcm.NonceSize() {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("wrapped data key is shorter than the nonce")
	}
	nonce := request.EncryptedDataKey[:gcm.NonceSize()]
	ciphertext := request.EncryptedDataKey[gcm.NonceSize():]
	dataKey, err := gcm.Open(nil, nonce, ciphertext, []byte(version))
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decrypt with key version %q: %w", version, err)
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func validateVersion(version string) error {
	if version == "" || filepath.Base(version) != version || version == "." || version == ".." {
		return fmt.Errorf("invalid key version %q", version)
	}
	return nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("load keyring AES key: %w", err)
	}
	return cipher.NewGCM(block)
}

func main() {
	sopsplugin.Serve(provider{})
}
