package kms

import (
	"bytes"
	"testing"
	"testing/quick"
	"time"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mozilla.org/sops/v3/kms/mocks"
)

func TestKMS(t *testing.T) {
	mockKMS := &mocks.KMSAPI{}
	defer mockKMS.AssertExpectations(t)
	kmsSvc = mockKMS
	isMocked = true
	encryptOutput := &kms.EncryptOutput{}
	decryptOutput := &kms.DecryptOutput{}
	mockKMS.On("Encrypt", mock.AnythingOfType("*kms.EncryptInput")).Return(encryptOutput, nil).Run(func(args mock.Arguments) {
		encryptOutput.CiphertextBlob = args.Get(0).(*kms.EncryptInput).Plaintext
	})
	mockKMS.On("Decrypt", mock.AnythingOfType("*kms.DecryptInput")).Return(decryptOutput, nil).Run(func(args mock.Arguments) {
		decryptOutput.Plaintext = args.Get(0).(*kms.DecryptInput).CiphertextBlob
	})
	k := MasterKey{Arn: "arn:aws:kms:us-east-1:927034868273:key/e9fc75db-05e9-44c1-9c35-633922bac347", Role: "", EncryptedKey: ""}
	f := func(x []byte) bool {
		err := k.Encrypt(x)
		if err != nil {
			log.Println(err)
		}
		v, err := k.Decrypt()
		if err != nil {
			log.Println(err)
		}
		return bytes.Equal(v, x)
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
	ks := MasterKeysFromArnString(s, nil, "foo", "SYMMETRIC_DEFAULT")
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

func TestParseEncryptionContext(t *testing.T) {
	value1 := "value1"
	value2 := "value2"
	// map from YAML
	var yamlmap = map[interface{}]interface{}{
		"key1": value1,
		"key2": value2,
	}
	assert.Equal(t, ParseKMSContext(yamlmap), map[string]*string{
		"key1": &value1,
		"key2": &value2,
	})
	assert.Nil(t, ParseKMSContext(map[interface{}]interface{}{}))
	assert.Nil(t, ParseKMSContext(map[interface{}]interface{}{
		"key1": 1,
	}))
	assert.Nil(t, ParseKMSContext(map[interface{}]interface{}{
		1: "value",
	}))
	// map from JSON
	var jsonmap = map[string]interface{}{
		"key1": value1,
		"key2": value2,
	}
	assert.Equal(t, ParseKMSContext(jsonmap), map[string]*string{
		"key1": &value1,
		"key2": &value2,
	})
	assert.Nil(t, ParseKMSContext(map[string]interface{}{}))
	assert.Nil(t, ParseKMSContext(map[string]interface{}{
		"key1": 1,
	}))
	// sops 2.0.x formatted encryption context as a comma-separated list of key:value pairs
	assert.Equal(t, ParseKMSContext("key1:value1,key2:value2"), map[string]*string{
		"key1": &value1,
		"key2": &value2,
	})
	assert.Equal(t, ParseKMSContext("key1:value1"), map[string]*string{
		"key1": &value1,
	})
	assert.Nil(t, ParseKMSContext("key1,key2:value2"))
	assert.Nil(t, ParseKMSContext("key1"))
}

func TestKeyToMap(t *testing.T) {
	value1 := "value1"
	value2 := "value2"
	key := MasterKey{
		CreationDate: time.Date(2016, time.October, 31, 10, 0, 0, 0, time.UTC),
		Arn:          "foo",
		Role:         "bar",
		EncryptedKey: "this is encrypted",
		EncryptionContext: map[string]*string{
			"key1": &value1,
			"key2": &value2,
		},
	}
	assert.Equal(t, map[string]interface{}{
		"arn":        "foo",
		"role":       "bar",
		"enc":        "this is encrypted",
		"created_at": "2016-10-31T10:00:00Z",
		"context": map[string]string{
			"key1": value1,
			"key2": value2,
		},
	}, key.ToMap())
}
