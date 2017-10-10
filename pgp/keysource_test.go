package pgp

import (
	"bytes"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestPGP(t *testing.T) {
	key := NewMasterKeyFromFingerprint("1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A")
	f := func(x []byte) bool {
		if x == nil || len(x) == 0 {
			return true
		}
		if err := key.Encrypt(x); err != nil {
			t.Errorf("Failed to encrypt: %#v err: %v", x, err)
			return false
		}
		k, err := key.Decrypt()
		if err != nil {
			t.Errorf("Failed to decrypt: %#v err: %v", x, err)
			return false
		}
		return bytes.Equal(x, k)
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestPGPKeySourceFromString(t *testing.T) {
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

func TestRetrievePGPKey(t *testing.T) {
	fingerprint := "1022470DE3F0BC54BC6AB62DE05550BC07FB1A0A"
	_, err := getKeyFromKeyServer("gpg.mozilla.org", fingerprint)
	assert.NoError(t, err)
}
