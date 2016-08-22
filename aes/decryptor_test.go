package aes

import (
	"fmt"
	"strings"
	"testing"
	"testing/quick"
)

func TestDecrypt(t *testing.T) {
	expected := "foo"
	key := strings.Repeat("f", 32)
	message := `ENC[AES256_GCM,data:oYyi,iv:MyIDYbT718JRr11QtBkcj3Dwm4k1aCGZBVeZf0EyV8o=,tag:t5z2Z023Up0kxwCgw1gNxg==,type:str]`
	decryption, err := Decrypt(message, key, []byte("bar:"))
	if err != nil {
		t.Errorf("%s", err)
	}
	if decryption != expected {
		t.Errorf("Decrypt(\"%s\", \"%s\") == \"%s\", expected %s", message, key, decryption, expected)
	}
}

func TestDecryptInvalidAac(t *testing.T) {
	message := `ENC[AES256_GCM,data:oYyi,iv:MyIDYbT718JRr11QtBkcj3Dwm4k1aCGZBVeZf0EyV8o=,tag:t5z2Z023Up0kxwCgw1gNxg==,type:str]`
	_, err := Decrypt(message, strings.Repeat("f", 32), []byte(""))
	if err == nil {
		t.Errorf("Decrypting with an invalid AAC should fail")
	}
}

func TestRoundtripString(t *testing.T) {
	key := strings.Repeat("f", 32)
	f := func(x string) bool {
		if x == "" {
			return true
		}
		s, err := Encrypt(x, key, []byte(""))
		if err != nil {
			fmt.Println(err)
			return false
		}
		d, err := Decrypt(s, key, []byte(""))
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
	key := strings.Repeat("f", 32)
	f := func(x float64) bool {
		s, err := Encrypt(x, key, []byte(""))
		if err != nil {
			fmt.Println(err)
			return false
		}
		d, err := Decrypt(s, key, []byte(""))
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
	key := strings.Repeat("f", 32)
	f := func(x int) bool {
		s, err := Encrypt(x, key, []byte(""))
		if err != nil {
			fmt.Println(err)
			return false
		}
		d, err := Decrypt(s, key, []byte(""))
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
	key := strings.Repeat("f", 32)
	f := func(x bool) bool {
		s, err := Encrypt(x, key, []byte(""))
		if err != nil {
			fmt.Println(err)
			return false
		}
		d, err := Decrypt(s, key, []byte(""))
		if err != nil {
			return false
		}
		return x == d
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}
