package sops

import (
	"testing"
	"testing/quick"
)

func TestKMS(t *testing.T) {
	// TODO: make this not terrible and mock KMS with a reverseable operation on the key, or something. Good luck running the tests on a machine that's not mine!
	ks := KMSKeySource{KMS: []KMS{
		KMS{Arn: "arn:aws:kms:us-east-1:927034868273:key/e9fc75db-05e9-44c1-9c35-633922bac347", Role: "", EncryptedKey: ""},
	}}
	f := func(x string) bool {
		ks.EncryptKeys(x)
		v, _ := ks.DecryptKeys()
		if x == "" {
			return true // we can't encrypt an empty string
		}
		return v == x
	}
	config := quick.Config{}
	if testing.Short() {
		config.MaxCount = 10
	}
	if err := quick.Check(f, &config); err != nil {
		t.Error(err)
	}
}
