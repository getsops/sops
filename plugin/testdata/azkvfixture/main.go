package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	publicKeyFile, err := configValue(request.Config, "publicKeyFile")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	vaultURL, err := configValue(request.Config, "vaultUrl")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	key, err := configValue(request.Config, "key")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	version, err := configValue(request.Config, "version")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	publicKeyBytes, err := os.ReadFile(publicKeyFile)
	if err != nil {
		savedPublicKey, ok := request.Metadata["public_key"].(string)
		if !ok {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("public key file unavailable and no saved public key in metadata")
		}
		publicKeyBytes = []byte(savedPublicKey)
	}
	publicKey, err := parsePublicKey(publicKeyBytes)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, request.DataKey, nil)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: ciphertext,
		Metadata: map[string]interface{}{
			"key_id":     fmt.Sprintf("%s/keys/%s/%s", vaultURL, key, version),
			"public_key": string(publicKeyBytes),
		},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	privateKeyFile, err := configValue(request.Config, "privateKeyFile")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	vaultURL, err := configValue(request.Config, "vaultUrl")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	key, err := configValue(request.Config, "key")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	version, err := configValue(request.Config, "version")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	keyID := fmt.Sprintf("%s/keys/%s/%s", vaultURL, key, version)
	if request.Metadata["key_id"] != keyID {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("provider key metadata does not match configuration")
	}
	privateKeyBytes, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	dataKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, request.EncryptedDataKey, nil)
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
}

func configValue(config map[string]interface{}, key string) (string, error) {
	value, ok := config[key].(string)
	if !ok || value == "" {
		return "", fmt.Errorf("missing string plugin config %q", key)
	}
	return value, nil
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

func parsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("no PEM private key found")
	}
	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	privateKey, ok := parsed.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not RSA")
	}
	return privateKey, nil
}

func main() {
	sopsplugin.Serve(provider{})
}
