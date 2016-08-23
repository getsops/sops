package kms

import (
	"fmt"
	"testing"
	"testing/quick"
)

func TestKMS(t *testing.T) {
	// TODO: make this not terrible and mock KMS with a reverseable operation on the key, or something. Good luck running the tests on a machine that's not mine!
	k := MasterKey{Arn: "arn:aws:kms:us-east-1:927034868273:key/e9fc75db-05e9-44c1-9c35-633922bac347", Role: "", EncryptedKey: ""}
	f := func(x string) bool {
		err := k.Encrypt(x)
		if err != nil {
			fmt.Println(err)
		}
		v, err := k.Decrypt()
		if err != nil {
			fmt.Println(err)
		}
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

func TestKMSKeySourceFromString(t *testing.T) {
	s := "arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e+arn:aws:iam::927034868273:role/sops-dev, arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
	ks := MasterKeysFromArnString(s)
	k1 := ks[0]
	k2 := ks[1]
	expectedArn1 := "arn:aws:kms:us-east-1:656532927350:key/920aff2e-c5f1-4040-943a-047fa387b27e"
	expectedRole1 := "arn:aws:iam::927034868273:role/sops-dev"
	if k1.Arn != expectedArn1 {
		t.Errorf("ARN mismatch. Expected %s, found %s", expectedArn1, k1.Arn)
	}
	if k1.Role != expectedRole1 {
		t.Errorf("Role mismatch. Expected %s, found %s", expectedRole1, k1.Role)
	}
	expectedArn2 := "arn:aws:kms:ap-southeast-1:656532927350:key/9006a8aa-0fa6-4c14-930e-a2dfb916de1d"
	expectedRole2 := ""
	if k2.Arn != expectedArn2 {
		t.Errorf("ARN mismatch. Expected %s, found %s", expectedArn2, k2.Arn)
	}
	if k2.Role != expectedRole2 {
		t.Errorf("Role mismatch. Expected empty role, found %s.", k2.Role)
	}
}
