package keyservice

import (
	"fmt"

	"go.mozilla.org/sops/v3/azkv"
	"go.mozilla.org/sops/v3/gcpkms"
	"go.mozilla.org/sops/v3/hcvault"
	"go.mozilla.org/sops/v3/kms"
	"go.mozilla.org/sops/v3/pgp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is a key service server that uses SOPS MasterKeys to fulfill requests
type Server struct {
	// Prompt indicates whether the server should prompt before decrypting or encrypting data
	Prompt bool
}

func (ks *Server) encryptWithPgp(key *PgpKey, plaintext []byte) ([]byte, error) {
	pgpKey := pgp.NewMasterKeyFromFingerprint(key.Fingerprint)
	err := pgpKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(pgpKey.EncryptedKey), nil
}

func (ks *Server) encryptWithKms(key *KmsKey, plaintext []byte) ([]byte, error) {
	kmsKey := kmsKeyToMasterKey(key)
	err := kmsKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(kmsKey.EncryptedKey), nil
}

func (ks *Server) encryptWithGcpKms(key *GcpKmsKey, plaintext []byte) ([]byte, error) {
	gcpKmsKey := gcpkms.MasterKey{
		ResourceID: key.ResourceId,
	}
	err := gcpKmsKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(gcpKmsKey.EncryptedKey), nil
}

func (ks *Server) encryptWithAzureKeyVault(key *AzureKeyVaultKey, plaintext []byte) ([]byte, error) {
	azkvKey := azkv.MasterKey{
		VaultURL: key.VaultUrl,
		Name:     key.Name,
		Version:  key.Version,
	}
	err := azkvKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(azkvKey.EncryptedKey), nil
}

func (ks *Server) encryptWithVault(key *VaultKey, plaintext []byte) ([]byte, error) {
	vaultKey := hcvault.MasterKey{
		VaultAddress: key.VaultAddress,
		EnginePath:   key.EnginePath,
		KeyName:      key.KeyName,
	}
	err := vaultKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(vaultKey.EncryptedKey), nil
}

func (ks *Server) decryptWithPgp(key *PgpKey, ciphertext []byte) ([]byte, error) {
	pgpKey := pgp.NewMasterKeyFromFingerprint(key.Fingerprint)
	pgpKey.EncryptedKey = string(ciphertext)
	plaintext, err := pgpKey.Decrypt()
	return []byte(plaintext), err
}

func (ks *Server) decryptWithKms(key *KmsKey, ciphertext []byte) ([]byte, error) {
	kmsKey := kmsKeyToMasterKey(key)
	kmsKey.EncryptedKey = string(ciphertext)
	plaintext, err := kmsKey.Decrypt()
	return []byte(plaintext), err
}

func (ks *Server) decryptWithGcpKms(key *GcpKmsKey, ciphertext []byte) ([]byte, error) {
	gcpKmsKey := gcpkms.MasterKey{
		ResourceID: key.ResourceId,
	}
	gcpKmsKey.EncryptedKey = string(ciphertext)
	plaintext, err := gcpKmsKey.Decrypt()
	return []byte(plaintext), err
}

func (ks *Server) decryptWithAzureKeyVault(key *AzureKeyVaultKey, ciphertext []byte) ([]byte, error) {
	azkvKey := azkv.MasterKey{
		VaultURL: key.VaultUrl,
		Name:     key.Name,
		Version:  key.Version,
	}
	azkvKey.EncryptedKey = string(ciphertext)
	plaintext, err := azkvKey.Decrypt()
	return []byte(plaintext), err
}

func (ks *Server) decryptWithVault(key *VaultKey, ciphertext []byte) ([]byte, error) {
	vaultKey := hcvault.MasterKey{
		VaultAddress: key.VaultAddress,
		EnginePath:   key.EnginePath,
		KeyName:      key.KeyName,
	}
	vaultKey.EncryptedKey = string(ciphertext)
	plaintext, err := vaultKey.Decrypt()
	return []byte(plaintext), err
}

// Encrypt takes an encrypt request and encrypts the provided plaintext with the provided key, returning the encrypted
// result
func (ks Server) Encrypt(ctx context.Context,
	req *EncryptRequest) (*EncryptResponse, error) {
	key := *req.Key
	var response *EncryptResponse
	switch k := key.KeyType.(type) {
	case *Key_PgpKey:
		ciphertext, err := ks.encryptWithPgp(k.PgpKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		response = &EncryptResponse{
			Ciphertext: ciphertext,
		}
	case *Key_KmsKey:
		ciphertext, err := ks.encryptWithKms(k.KmsKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		response = &EncryptResponse{
			Ciphertext: ciphertext,
		}
	case *Key_GcpKmsKey:
		ciphertext, err := ks.encryptWithGcpKms(k.GcpKmsKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		response = &EncryptResponse{
			Ciphertext: ciphertext,
		}
	case *Key_AzureKeyvaultKey:
		ciphertext, err := ks.encryptWithAzureKeyVault(k.AzureKeyvaultKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		response = &EncryptResponse{
			Ciphertext: ciphertext,
		}
	case *Key_VaultKey:
		ciphertext, err := ks.encryptWithVault(k.VaultKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		response = &EncryptResponse{
			Ciphertext: ciphertext,
		}
	case nil:
		return nil, status.Errorf(codes.NotFound, "Must provide a key")
	default:
		return nil, status.Errorf(codes.NotFound, "Unknown key type")
	}
	if ks.Prompt {
		err := ks.prompt(key, "encrypt")
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}

func keyToString(key Key) string {
	switch k := key.KeyType.(type) {
	case *Key_PgpKey:
		return fmt.Sprintf("PGP key with fingerprint %s", k.PgpKey.Fingerprint)
	case *Key_KmsKey:
		return fmt.Sprintf("AWS KMS key with ARN %s", k.KmsKey.Arn)
	case *Key_GcpKmsKey:
		return fmt.Sprintf("GCP KMS key with resource ID %s", k.GcpKmsKey.ResourceId)
	case *Key_AzureKeyvaultKey:
		return fmt.Sprintf("Azure Key Vault key with URL %s/keys/%s/%s", k.AzureKeyvaultKey.VaultUrl, k.AzureKeyvaultKey.Name, k.AzureKeyvaultKey.Version)
	case *Key_VaultKey:
		return fmt.Sprintf("Hashicorp Vault key with URI %s/v1/%s/keys/%s", k.VaultKey.VaultAddress, k.VaultKey.EnginePath, k.VaultKey.KeyName)
	default:
		return fmt.Sprintf("Unknown key type")
	}
}

func (ks Server) prompt(key Key, requestType string) error {
	keyString := keyToString(key)
	var response string
	for response != "y" && response != "n" {
		fmt.Printf("\nReceived %s request using %s. Respond to request? (y/n): ", requestType, keyString)
		_, err := fmt.Scanln(&response)
		if err != nil {
			return err
		}
	}
	if response == "n" {
		return grpc.Errorf(codes.PermissionDenied, "Request rejected by user")
	}
	return nil
}

// Decrypt takes a decrypt request and decrypts the provided ciphertext with the provided key, returning the decrypted
// result
func (ks Server) Decrypt(ctx context.Context,
	req *DecryptRequest) (*DecryptResponse, error) {
	key := *req.Key
	var response *DecryptResponse
	switch k := key.KeyType.(type) {
	case *Key_PgpKey:
		plaintext, err := ks.decryptWithPgp(k.PgpKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		response = &DecryptResponse{
			Plaintext: plaintext,
		}
	case *Key_KmsKey:
		plaintext, err := ks.decryptWithKms(k.KmsKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		response = &DecryptResponse{
			Plaintext: plaintext,
		}
	case *Key_GcpKmsKey:
		plaintext, err := ks.decryptWithGcpKms(k.GcpKmsKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		response = &DecryptResponse{
			Plaintext: plaintext,
		}
	case *Key_AzureKeyvaultKey:
		plaintext, err := ks.decryptWithAzureKeyVault(k.AzureKeyvaultKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		response = &DecryptResponse{
			Plaintext: plaintext,
		}
	case *Key_VaultKey:
		plaintext, err := ks.decryptWithVault(k.VaultKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		response = &DecryptResponse{
			Plaintext: plaintext,
		}
	case nil:
		return nil, grpc.Errorf(codes.NotFound, "Must provide a key")
	default:
		return nil, grpc.Errorf(codes.NotFound, "Unknown key type")
	}
	if ks.Prompt {
		err := ks.prompt(key, "decrypt")
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}

func kmsKeyToMasterKey(key *KmsKey) kms.MasterKey {
	ctx := make(map[string]*string)
	for k, v := range key.Context {
		value := v // Allocate a new string to prevent the pointer below from referring to only the last iteration value
		ctx[k] = &value
	}
	return kms.MasterKey{
		Arn:               key.Arn,
		Role:              key.Role,
		EncryptionContext: ctx,
		AwsProfile:        key.AwsProfile,
	}
}
