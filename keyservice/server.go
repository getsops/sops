package keyservice

import (
	"go.mozilla.org/sops/kms"
	"go.mozilla.org/sops/pgp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Server struct{}

func (ks *Server) encryptWithPgp(key *PgpKey, plaintext []byte) ([]byte, error) {
	pgpKey := pgp.NewMasterKeyFromFingerprint(key.Fingerprint)
	err := pgpKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(pgpKey.EncryptedKey), nil
}

func (ks *Server) encryptWithKms(key *KmsKey, plaintext []byte) ([]byte, error) {
	var ctx map[string]*string
	for k, v := range key.Context {
		ctx[k] = &v
	}
	kmsKey := kms.MasterKey{
		Arn:               key.Arn,
		Role:              key.Role,
		EncryptionContext: ctx,
	}
	err := kmsKey.Encrypt(plaintext)
	if err != nil {
		return nil, err
	}
	return []byte(kmsKey.EncryptedKey), nil
}

func (ks *Server) decryptWithPgp(key *PgpKey, ciphertext []byte) ([]byte, error) {
	pgpKey := pgp.NewMasterKeyFromFingerprint(key.Fingerprint)
	pgpKey.EncryptedKey = string(ciphertext)
	plaintext, err := pgpKey.Decrypt()
	return []byte(plaintext), err
}

func (ks *Server) decryptWithKms(key *KmsKey, ciphertext []byte) ([]byte, error) {
	var ctx map[string]*string
	for k, v := range key.Context {
		ctx[k] = &v
	}
	kmsKey := kms.MasterKey{
		Arn:               key.Arn,
		Role:              key.Role,
		EncryptionContext: ctx,
	}
	kmsKey.EncryptedKey = string(ciphertext)
	plaintext, err := kmsKey.Decrypt()
	return []byte(plaintext), err
}

func (ks Server) Encrypt(ctx context.Context,
	req *EncryptRequest) (*EncryptResponse, error) {
	key := *req.Key
	switch k := key.KeyType.(type) {
	case *Key_PgpKey:
		ciphertext, err := ks.encryptWithPgp(k.PgpKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		return &EncryptResponse{
			Ciphertext: ciphertext,
		}, nil
	case *Key_KmsKey:
		ciphertext, err := ks.encryptWithKms(k.KmsKey, req.Plaintext)
		if err != nil {
			return nil, err
		}
		return &EncryptResponse{
			Ciphertext: ciphertext,
		}, nil
	case nil:
		return nil, grpc.Errorf(codes.NotFound, "Must provide a key")
	default:
		return nil, grpc.Errorf(codes.NotFound, "Unknown key type")
	}
}

func (ks Server) Decrypt(ctx context.Context,
	req *DecryptRequest) (*DecryptResponse, error) {
	key := *req.Key
	switch k := key.KeyType.(type) {
	case *Key_PgpKey:
		plaintext, err := ks.decryptWithPgp(k.PgpKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		return &DecryptResponse{
			Plaintext: plaintext,
		}, nil
	case *Key_KmsKey:
		plaintext, err := ks.decryptWithKms(k.KmsKey, req.Ciphertext)
		if err != nil {
			return nil, err
		}
		return &DecryptResponse{
			Plaintext: plaintext,
		}, nil
	case nil:
		return nil, grpc.Errorf(codes.NotFound, "Must provide a key")
	default:
		return nil, grpc.Errorf(codes.NotFound, "Unknown key type")
	}
}
