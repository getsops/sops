package keyservice

import (
	"errors"
	"testing"
)

type unsupportedMasterKey struct{}

func (k *unsupportedMasterKey) Encrypt([]byte) error           { return nil }
func (k *unsupportedMasterKey) EncryptIfNeeded([]byte) error   { return nil }
func (k *unsupportedMasterKey) EncryptedDataKey() []byte       { return nil }
func (k *unsupportedMasterKey) SetEncryptedDataKey([]byte)     {}
func (k *unsupportedMasterKey) Decrypt() ([]byte, error)       { return nil, nil }
func (k *unsupportedMasterKey) NeedsRotation() bool            { return false }
func (k *unsupportedMasterKey) ToString() string               { return "unsupported" }
func (k *unsupportedMasterKey) ToMap() map[string]interface{}  { return nil }
func (k *unsupportedMasterKey) TypeToIdentifier() string       { return "unsupported" }

func TestKeyFromMasterKeyOrErrorUnsupportedType(t *testing.T) {
	_, err := KeyFromMasterKeyOrError(&unsupportedMasterKey{})
	if err == nil {
		t.Fatal("expected error for unsupported key type")
	}
	if !errors.Is(err, ErrUnsupportedMasterKeyType) {
		t.Fatalf("expected ErrUnsupportedMasterKeyType, got: %v", err)
	}
}

func TestKeyFromMasterKeyDoesNotPanicOnUnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("expected no panic, got panic: %v", r)
		}
	}()

	_ = KeyFromMasterKey(&unsupportedMasterKey{})
}
