package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/stores/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const kmsProfileFile = `{
	"data": "ENC[AES256_GCM,data:abcd,iv:abcd,tag:abcd,type:str]",
	"sops": {
		"kms": [
			{
				"arn": "arn:aws:kms:us-west-2:927034868273:key/e9fc75db-05e9-44c1-9c35-633922bac347",
				"created_at": "2020-01-01T00:00:00Z",
				"enc": "AQIDAHhfoo",
				"aws_profile": "stored"
			}
		],
		"lastmodified": "2020-01-01T00:00:00Z",
		"mac": "ENC[AES256_GCM,data:abcd,iv:abcd,tag:abcd,type:str]",
		"version": "3.9.0"
	}
}`

func TestLoadEncryptedFileWithBugFixesAwsProfile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.json")
	require.NoError(t, os.WriteFile(path, []byte(kmsProfileFile), 0o600))

	tree, err := LoadEncryptedFileWithBugFixes(GenericDecryptOpts{
		InputStore: &json.Store{},
		InputPath:  path,
		AwsProfile: "override",
	})
	require.NoError(t, err)
	key, ok := tree.Metadata.KeyGroups[0][0].(*kms.MasterKey)
	require.True(t, ok)
	assert.Equal(t, "override", key.EffectiveAwsProfile())
	assert.Equal(t, "stored", key.AwsProfile)

	// A tree reloaded from text (as in the editor loop of edit) needs the
	// override applied again.
	reloaded, err := (&json.Store{}).LoadEncryptedFile([]byte(kmsProfileFile))
	require.NoError(t, err)
	SetAWSProfile(reloaded.Metadata.KeyGroups, "override")
	assert.Equal(t, "override", reloaded.Metadata.KeyGroups[0][0].(*kms.MasterKey).EffectiveAwsProfile())
}

func TestLoadEncryptedFileWithBugFixesAwsProfileNotEmitted(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.json")
	require.NoError(t, os.WriteFile(path, []byte(kmsProfileFile), 0o600))
	store := &json.Store{}
	tree, err := LoadEncryptedFileWithBugFixes(GenericDecryptOpts{
		InputStore: store,
		InputPath:  path,
		AwsProfile: "override",
	})
	require.NoError(t, err)
	out, err := store.EmitEncryptedFile(*tree)
	require.NoError(t, err)
	assert.Contains(t, string(out), `"aws_profile": "stored"`)
	assert.NotContains(t, string(out), "override")
}

func TestSetAWSProfile(t *testing.T) {
	newGroups := func() []sops.KeyGroup {
		return []sops.KeyGroup{{&kms.MasterKey{Arn: "a", AwsProfile: "stored"}, &age.MasterKey{Recipient: "r"}}}
	}
	for name, tc := range map[string]struct{ profile, want string }{
		"empty is a no-op":      {"", ""},
		"sets on KMS keys only": {"override", "override"},
	} {
		t.Run(name, func(t *testing.T) {
			groups := newGroups()
			SetAWSProfile(groups, tc.profile)
			key := groups[0][0].(*kms.MasterKey)
			assert.Equal(t, tc.want, key.AwsProfileOverride)
			assert.Equal(t, "stored", key.AwsProfile)
			assert.Equal(t, newGroups()[0][1], groups[0][1])
		})
	}
}

func TestPersistAWSProfile(t *testing.T) {
	groups := []sops.KeyGroup{{&kms.MasterKey{Arn: "a", AwsProfile: "stored"}}}
	PersistAWSProfile(groups, "")
	assert.Equal(t, "stored", groups[0][0].(*kms.MasterKey).AwsProfile)
	PersistAWSProfile(groups, "new")
	key := groups[0][0].(*kms.MasterKey)
	assert.Equal(t, "new", key.AwsProfile)
	assert.Empty(t, key.AwsProfileOverride)
}
