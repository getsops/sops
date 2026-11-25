package updatekeys

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	agecrypto "filippo.io/age"
	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/pgp"
	"github.com/getsops/sops/v3/stores/yaml"
	"github.com/stretchr/testify/require"
)

type ageTestServer struct {
	keyservice.Server
	identities age.ParsedIdentities
}

func (s ageTestServer) Decrypt(_ context.Context, request *keyservice.DecryptRequest) (*keyservice.DecryptResponse, error) {
	key, err := age.MasterKeyFromRecipient(request.Key.GetAgeKey().Recipient)
	if err != nil {
		return nil, err
	}
	key.SetEncryptedDataKey(request.Ciphertext)
	s.identities.ApplyToMasterKey(key)
	plaintext, err := key.Decrypt()
	return &keyservice.DecryptResponse{Plaintext: plaintext}, err
}

func testKey(recipient string) *age.MasterKey {
	return &age.MasterKey{Recipient: recipient}
}

func TestKeyGroupsChanged(t *testing.T) {
	first := testKey("age1first")
	second := testKey("age1second")

	tests := []struct {
		name    string
		current []sops.KeyGroup
		desired []sops.KeyGroup
		changed bool
	}{
		{
			name:    "same order",
			current: []sops.KeyGroup{{first, second}},
			desired: []sops.KeyGroup{{testKey("age1first"), testKey("age1second")}},
			changed: false,
		},
		{
			name:    "key order changes",
			current: []sops.KeyGroup{{first, second}},
			desired: []sops.KeyGroup{{testKey("age1second"), testKey("age1first")}},
			changed: true,
		},
		{
			name:    "group order changes",
			current: []sops.KeyGroup{{first}, {second}},
			desired: []sops.KeyGroup{{testKey("age1second")}, {testKey("age1first")}},
			changed: true,
		},
		{
			name:    "different type ordering is not persisted",
			current: []sops.KeyGroup{{&pgp.MasterKey{Fingerprint: "fingerprint"}, first}},
			desired: []sops.KeyGroup{{first, &pgp.MasterKey{Fingerprint: "fingerprint"}}},
		},
		{
			name:    "key added",
			current: []sops.KeyGroup{{first}},
			desired: []sops.KeyGroup{{testKey("age1first"), testKey("age1second")}},
			changed: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := keyGroupsChanged(test.current, test.desired); got != test.changed {
				t.Fatalf("keyGroupsChanged() = %v, want %v", got, test.changed)
			}
		})
	}
}

func TestUpdateKeysReordersAgeRecipients(t *testing.T) {
	first, err := agecrypto.GenerateX25519Identity()
	require.NoError(t, err)
	second, err := agecrypto.GenerateX25519Identity()
	require.NoError(t, err)
	services := []keyservice.KeyServiceClient{keyservice.NewCustomLocalClient(ageTestServer{
		identities: age.ParsedIdentities{first, second},
	})}
	store := yaml.NewStore(&config.YAMLStoreConfig{})
	tree := sops.Tree{
		Branches: sops.TreeBranches{{{Key: "secret", Value: "value"}}},
		Metadata: sops.Metadata{KeyGroups: []sops.KeyGroup{{
			testKey(first.Recipient().String()), testKey(second.Recipient().String()),
		}}},
	}
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(services)
	require.Empty(t, errs)
	cipher := aes.NewCipher()
	require.NoError(t, common.EncryptTree(common.EncryptTreeOpts{Tree: &tree, Cipher: cipher, DataKey: dataKey}))
	before, err := store.EmitEncryptedFile(tree)
	require.NoError(t, err)
	dir := t.TempDir()
	input := filepath.Join(dir, "secret.yaml")
	configuration := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(input, before, 0600))
	rule := fmt.Sprintf("creation_rules:\n  - age: %s,%s\n", second.Recipient(), first.Recipient())
	require.NoError(t, os.WriteFile(configuration, []byte(rule), 0600))
	opts := Opts{InputPath: input, ConfigPath: configuration, KeyServices: services}
	require.NoError(t, UpdateKeys(opts))
	after, err := os.ReadFile(input)
	require.NoError(t, err)
	loaded, err := store.LoadEncryptedFile(after)
	require.NoError(t, err)
	require.Equal(t, second.Recipient().String(), loaded.Metadata.KeyGroups[0][0].ToString())
	require.Equal(t, first.Recipient().String(), loaded.Metadata.KeyGroups[0][1].ToString())
	require.Equal(t, tree.Branches, loaded.Branches)
	require.Equal(t, tree.Metadata.MessageAuthenticationCode, loaded.Metadata.MessageAuthenticationCode)
	_, err = common.DecryptTree(common.DecryptTreeOpts{Tree: &loaded, Cipher: cipher, KeyServices: opts.KeyServices})
	require.NoError(t, err)
	require.Equal(t, "value", loaded.Branches[0][0].Value)
	// A second update must leave the file byte-for-byte unchanged.
	require.NoError(t, UpdateKeys(opts))
	again, err := os.ReadFile(input)
	require.NoError(t, err)
	require.Equal(t, after, again)
}

func TestUpdateKeysMixedTypesAlreadyCurrent(t *testing.T) {
	identity, err := agecrypto.GenerateX25519Identity()
	require.NoError(t, err)
	store := yaml.NewStore(&config.YAMLStoreConfig{})
	tree := sops.Tree{Branches: sops.TreeBranches{{{Key: "public", Value: "value"}}}, Metadata: sops.Metadata{
		LastModified: time.Now().UTC(),
		KeyGroups: []sops.KeyGroup{{testKey(identity.Recipient().String()),
			&pgp.MasterKey{Fingerprint: "0123456789012345678901234567890123456789", CreationDate: time.Now().UTC()},
		}},
	}}
	before, err := store.EmitEncryptedFile(tree)
	require.NoError(t, err)
	dir := t.TempDir()
	input := filepath.Join(dir, "secret.yaml")
	configuration := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(input, before, 0600))
	rule := fmt.Sprintf("creation_rules:\n  - key_groups:\n      - age: [%s]\n        pgp: [0123456789012345678901234567890123456789]\n", identity.Recipient())
	require.NoError(t, os.WriteFile(configuration, []byte(rule), 0600))
	// No key service: the unchanged file must not require decryption.
	require.NoError(t, UpdateKeys(Opts{InputPath: input, ConfigPath: configuration}))
	after, err := os.ReadFile(input)
	require.NoError(t, err)
	require.Equal(t, before, after)
}
