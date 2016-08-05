package decryptor

import (
	"strings"
	"testing"
)

func TestDecrypt(t *testing.T) {
	expected := "foo"
	key := strings.Repeat("f", 32)
	message := `ENC[AES256_GCM,data:xPxG,iv:kMAhrJOMitZZP3C71cA1wnp543hHYFd8+Tv01hdEOqc=,tag:PDVgtlbfBU7A33NKugzNBg==,type:bytes]`
	decryption, err := Decrypt(message, key)
	if err != nil {
		t.Errorf("%s", err)
	}
	if decryption != expected {
		t.Errorf("Decrypt(\"%s\", \"%s\") == \"%s\", expected %s", message, key, decryption, expected)
	}
}
