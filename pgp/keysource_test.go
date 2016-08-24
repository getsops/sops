package pgp

import (
	"bytes"
	"testing"
	"testing/quick"
)

func TestGPG(t *testing.T) {
	key := NewMasterKeyFromFingerprint("64FEF099B0544CF975BCD408A014A073E0848B51")
	f := func(x []byte) bool {
		if x == nil || len(x) == 0 {
			return true
		}
		key.Encrypt(x)
		k, _ := key.Decrypt()
		return bytes.Equal(x, k)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestGPGKeySourceFromString(t *testing.T) {
	s := "C8C5 2C0A B2A4 8174 01E8  12C8 F3CC 3233 3FAD 9F1E, C8C5 2C0A B2A4 8174 01E8  12C8 F3CC 3233 3FAD 9F1E"
	ks := MasterKeysFromFingerprintString(s)
	expected := "C8C52C0AB2A4817401E812C8F3CC32333FAD9F1E"
	if ks[0].Fingerprint != expected {
		t.Errorf("Fingerprint does not match. Got %s, expected %s", ks[0].Fingerprint, expected)
	}

	if ks[1].Fingerprint != expected {
		t.Error("Fingerprint does not match")
	}
}
