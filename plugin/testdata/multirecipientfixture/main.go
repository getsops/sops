package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"os"
	"sort"

	sopsplugin "github.com/getsops/sops/v3/plugin"
)

type provider struct{}

func (provider) Wrap(_ context.Context, request sopsplugin.WrapKeyRequest) (sopsplugin.WrapKeyResponse, error) {
	recipients, err := stringMap(request.Config["recipients"], "recipients")
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, err
	}
	recipientIDs := sortedKeys(recipients)
	wrapped := make(map[string]string, len(recipients))
	for _, recipientID := range recipientIDs {
		publicKey, err := parsePublicKey([]byte(recipients[recipientID]))
		if err != nil {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("parse recipient %q public key: %w", recipientID, err)
		}
		ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, request.DataKey, nil)
		if err != nil {
			return sopsplugin.WrapKeyResponse{}, fmt.Errorf("wrap data key for recipient %q: %w", recipientID, err)
		}
		wrapped[recipientID] = base64.StdEncoding.EncodeToString(ciphertext)
	}
	if len(wrapped) == 0 {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("at least one recipient is required")
	}
	encryptedDataKey, err := json.Marshal(wrapped)
	if err != nil {
		return sopsplugin.WrapKeyResponse{}, fmt.Errorf("encode recipient ciphertexts: %w", err)
	}
	return sopsplugin.WrapKeyResponse{
		EncryptedDataKey: encryptedDataKey,
		Metadata:         map[string]interface{}{"recipient_ids": recipientIDs},
	}, nil
}

func (provider) Unwrap(_ context.Context, request sopsplugin.UnwrapKeyRequest) (sopsplugin.UnwrapKeyResponse, error) {
	var wrapped map[string]string
	if err := json.Unmarshal(request.EncryptedDataKey, &wrapped); err != nil {
		return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("decode recipient ciphertexts: %w", err)
	}
	privateKeyFiles, err := stringMap(request.Config["private_key_files"], "private_key_files")
	if err != nil {
		return sopsplugin.UnwrapKeyResponse{}, err
	}
	var failures []string
	for _, recipientID := range sortedKeys(wrapped) {
		privateKeyFile, ok := privateKeyFiles[recipientID]
		if !ok {
			continue
		}
		ciphertext, err := base64.StdEncoding.DecodeString(wrapped[recipientID])
		if err != nil {
			failures = append(failures, fmt.Sprintf("decode recipient %q ciphertext: %v", recipientID, err))
			continue
		}
		privateKeyBytes, err := os.ReadFile(privateKeyFile)
		if err != nil {
			failures = append(failures, fmt.Sprintf("read recipient %q private key: %v", recipientID, err))
			continue
		}
		privateKey, err := parsePrivateKey(privateKeyBytes)
		if err != nil {
			failures = append(failures, fmt.Sprintf("parse recipient %q private key: %v", recipientID, err))
			continue
		}
		dataKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
		if err != nil {
			failures = append(failures, fmt.Sprintf("decrypt for recipient %q: %v", recipientID, err))
			continue
		}
		return sopsplugin.UnwrapKeyResponse{DataKey: dataKey}, nil
	}
	return sopsplugin.UnwrapKeyResponse{}, fmt.Errorf("no configured recipient could unwrap the data key: %v", failures)
}

func stringMap(value interface{}, name string) (map[string]string, error) {
	entries, ok := value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing object plugin config %q", name)
	}
	result := make(map[string]string, len(entries))
	for key, value := range entries {
		text, ok := value.(string)
		if !ok || text == "" {
			return nil, fmt.Errorf("plugin config %q entry %q must be a non-empty string", name, key)
		}
		result[key] = text
	}
	return result, nil
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
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
