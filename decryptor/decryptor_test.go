package decryptor

import (
	"strings"
	"testing"
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
