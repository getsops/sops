package keyservice

import (
	"testing"

	"github.com/getsops/sops/v3/kms"
	"github.com/stretchr/testify/assert"
)

func TestKeyFromMasterKeyAwsProfileOverride(t *testing.T) {
	mk := kms.NewMasterKeyWithProfile("arn:aws:kms:us-west-2:111122223333:key/abc", "", nil, "stored")
	key := KeyFromMasterKey(mk)
	assert.Equal(t, "stored", key.GetKmsKey().AwsProfile)
	mk.AwsProfileOverride = "override"
	key = KeyFromMasterKey(mk)
	assert.Equal(t, "override", key.GetKmsKey().AwsProfile)
}
