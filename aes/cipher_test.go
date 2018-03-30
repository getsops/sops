package aes

import (
	"crypto/rand"
	"strings"
	"testing"
	"testing/quick"
)

func TestDecrypt(t *testing.T) {
	expected := "foo"
	key := []byte(strings.Repeat("f", 32))
	message := `ENC[AES256_GCM,data:oYyi,iv:MyIDYbT718JRr11QtBkcj3Dwm4k1aCGZBVeZf0EyV8o=,tag:t5z2Z023Up0kxwCgw1gNxg==,type:str]`
	decryption, err := NewCipher().Decrypt(message, key, "bar:")
	if err != nil {
		t.Errorf("%s", err)
	}
	if decryption != expected {
		t.Errorf("Decrypt(\"%s\", \"%s\") == \"%s\", expected %s", message, key, decryption, expected)
	}
}

func TestDecryptInvalidAad(t *testing.T) {
	message := `ENC[AES256_GCM,data:oYyi,iv:MyIDYbT718JRr11QtBkcj3Dwm4k1aCGZBVeZf0EyV8o=,tag:t5z2Z023Up0kxwCgw1gNxg==,type:str]`
	_, err := NewCipher().Decrypt(message, []byte(strings.Repeat("f", 32)), "")
	if err == nil {
		t.Errorf("Decrypting with an invalid AAC should fail")
	}
}

func TestRoundtripString(t *testing.T) {
	f := func(x, aad string) bool {
		key := make([]byte, 32)
		rand.Read(key)
		s, err := NewCipher().Encrypt(x, key, aad)
		if err != nil {
			log.Println(err)
			return false
		}
		d, err := NewCipher().Decrypt(s, key, aad)
		if err != nil {
			return false
		}
		return x == d
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRoundtripFloat(t *testing.T) {
	key := []byte(strings.Repeat("f", 32))
	f := func(x float64) bool {
		s, err := NewCipher().Encrypt(x, key, "")
		if err != nil {
			log.Println(err)
			return false
		}
		d, err := NewCipher().Decrypt(s, key, "")
		if err != nil {
			return false
		}
		return x == d
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRoundtripInt(t *testing.T) {
	key := []byte(strings.Repeat("f", 32))
	f := func(x int) bool {
		s, err := NewCipher().Encrypt(x, key, "")
		if err != nil {
			log.Println(err)
			return false
		}
		d, err := NewCipher().Decrypt(s, key, "")
		if err != nil {
			return false
		}
		return x == d
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRoundtripBool(t *testing.T) {
	key := []byte(strings.Repeat("f", 32))
	f := func(x bool) bool {
		s, err := NewCipher().Encrypt(x, key, "")
		if err != nil {
			log.Println(err)
			return false
		}
		d, err := NewCipher().Decrypt(s, key, "")
		if err != nil {
			return false
		}
		return x == d
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
