package acskms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewMasterKey(t *testing.T) {
	cases := []struct {
		desc           string
		input          string
		expectedRegion string
		expectErr      bool
	}{
		{
			desc:           "valid key ARN",
			input:          "acs:kms:cn-shanghai:1234567890:key/00000000-0000-0000-0000-000000000000",
			expectedRegion: "cn-shanghai",
			expectErr:      false,
		},
		{
			desc:      "alias ARN not supported",
			input:     "acs:kms:cn-hangzhou:1234567890:alias/my-alias",
			expectErr: true,
		},
		{
			desc:      "invalid ARN format",
			input:     "invalid:arn",
			expectErr: true,
		},
		{
			desc:      "missing region",
			input:     "acs:kms::1234567890:key/id",
			expectErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			k, err := NewMasterKey(c.input)
			if c.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, c.input, k.Arn)
				assert.Equal(t, c.expectedRegion, k.Region)
				assert.WithinDuration(t, time.Now().UTC(), k.CreationDate, 5*time.Second)
			}
		})
	}
}

func TestNewMasterKeyFromKeyIDString(t *testing.T) {
	cases := []struct {
		desc      string
		input     string
		count     int
		expectErr bool
	}{
		{
			desc:      "single key",
			input:     "acs:kms:cn-shanghai:1234567890:key/key1",
			count:     1,
			expectErr: false,
		},
		{
			desc:      "multiple keys",
			input:     "acs:kms:cn-shanghai:1234567890:key/key1,acs:kms:cn-hangzhou:1234567890:key/key2",
			count:     2,
			expectErr: false,
		},
		{
			desc:      "empty string",
			input:     "",
			count:     0,
			expectErr: false,
		},
		{
			desc:      "whitespace handling",
			input:     " acs:kms:cn-shanghai:1234567890:key/key1 , acs:kms:cn-hangzhou:1234567890:key/key2 ",
			count:     2,
			expectErr: false,
		},
		{
			desc:      "invalid key in list",
			input:     "acs:kms:cn-shanghai:1234567890:key/key1,invalid",
			count:     0,
			expectErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.desc, func(t *testing.T) {
			keys, err := NewMasterKeyFromKeyIDString(c.input)
			if c.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, keys, c.count)
			}
		})
	}
}

func TestMasterKey_ToMap(t *testing.T) {
	arn := "acs:kms:cn-shanghai:1234567890:key/key1"
	k, err := NewMasterKey(arn)
	assert.NoError(t, err)

	k.EncryptedKey = "encrypted_data"
	m := k.ToMap()

	assert.Equal(t, arn, m["arn"])
	assert.Equal(t, "encrypted_data", m["enc"])
	assert.NotEmpty(t, m["created_at"])
}

func TestMasterKey_MethodProxies(t *testing.T) {
	arn := "acs:kms:cn-shanghai:1234567890:key/key1"
	k, err := NewMasterKey(arn)
	assert.NoError(t, err)

	// Test EncryptedDataKey and SetEncryptedDataKey
	k.SetEncryptedDataKey([]byte("test"))
	assert.Equal(t, []byte("test"), k.EncryptedDataKey())

	// Test ToString
	assert.Equal(t, arn, k.ToString())

	// Test NeedsRotation (should be false as per implementation)
	assert.False(t, k.NeedsRotation())

	// Test EncryptIfNeeded (noop if already encrypted)
	k.EncryptedKey = "already_encrypted"
	err = k.EncryptIfNeeded([]byte("data"))
	assert.NoError(t, err)
	assert.Equal(t, "already_encrypted", k.EncryptedKey)
}
