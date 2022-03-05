package pgp

import (
	"bytes"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestPGP(t *testing.T) {
	key := NewMasterKeyFromFingerprint("FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4")
	f := func(x []byte) bool {
		if x == nil || len(x) == 0 {
			return true
		}
		if err := key.Encrypt(x); err != nil {
			t.Errorf("Failed to encrypt: %#v err: %w", x, err)
			return false
		}
		k, err := key.Decrypt()
		if err != nil {
			t.Errorf("Failed to decrypt: %#v err: %w", x, err)
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
	// Requires a key available in https://keys.openpgp.org/ *with identity information* (that is, an email address).
	// See https://keys.openpgp.org/about/faq#verify-multiple for details about identity information.
	// We use the key of release@mozilla.com for here.
	fingerprint := "14F26682D0916CDD81E37B6D61B7B526D98F0353"
	_, err := getKeyFromKeyServer(fingerprint)
	assert.NoError(t, err)
}
