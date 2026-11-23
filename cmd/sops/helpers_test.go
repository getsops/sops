package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	fage "filippo.io/age"
	"github.com/stretchr/testify/require"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/cmd/sops/formats"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/keyservice"
)

// isolateAgeEnvironment makes sure that the age key source only sees the
// identities passed through SOPS_AGE_KEY by newTestAgeKey, and not keys of
// the user running the tests (like ~/.ssh/id_ed25519).
func isolateAgeEnvironment(t *testing.T) {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(home, ".config"))
	for _, name := range []string{
		age.SopsAgeKeyEnv,
		age.SopsAgeKeyFileEnv,
		age.SopsAgeKeyCmdEnv,
		age.SopsAgeSshPrivateKeyFileEnv,
		age.SopsAgeSshPrivateKeyCmdEnv,
	} {
		// t.Setenv registers the cleanup that restores the original value.
		t.Setenv(name, "")
		require.NoError(t, os.Unsetenv(name))
	}
}

// newTestAgeKey generates a throwaway age identity, makes it available for
// decryption through SOPS_AGE_KEY, and returns the matching master key.
func newTestAgeKey(t *testing.T) *age.MasterKey {
	t.Helper()
	identity, err := fage.GenerateX25519Identity()
	require.NoError(t, err)
	identities := []string{identity.String()}
	if existing, ok := os.LookupEnv(age.SopsAgeKeyEnv); ok && existing != "" {
		identities = append([]string{existing}, identities...)
	}
	t.Setenv(age.SopsAgeKeyEnv, strings.Join(identities, "\n"))
	key, err := age.MasterKeyFromRecipient(identity.Recipient().String())
	require.NoError(t, err)
	return key
}

func testYAMLStore() common.Store {
	return common.StoreForFormat(formats.Yaml, config.NewStoresConfig())
}

func testKeyServices() []keyservice.KeyServiceClient {
	return []keyservice.KeyServiceClient{keyservice.NewLocalClient()}
}

// encryptTestFile writes plaintext YAML to a temporary file, encrypts it in
// place with the given key groups, and returns the file's path.
func encryptTestFile(t *testing.T, plaintext string, groups []sops.KeyGroup) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "secret.yaml")
	require.NoError(t, os.WriteFile(path, []byte(plaintext), 0600))
	store := testYAMLStore()
	encrypted, err := encrypt(encryptOpts{
		Cipher:      aes.NewCipher(),
		InputStore:  store,
		OutputStore: store,
		InputPath:   path,
		KeyServices: testKeyServices(),
		encryptConfig: encryptConfig{
			KeyGroups:         groups,
			UnencryptedSuffix: sops.DefaultUnencryptedSuffix,
		},
	})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(path, encrypted, 0600))
	return path
}

// decryptTestFile decrypts the file at path and returns the plaintext YAML.
func decryptTestFile(t *testing.T, path string) string {
	t.Helper()
	store := testYAMLStore()
	plaintext, err := decrypt(decryptOpts{
		Cipher:      aes.NewCipher(),
		InputStore:  store,
		OutputStore: store,
		InputPath:   path,
		KeyServices: testKeyServices(),
	})
	require.NoError(t, err)
	return string(plaintext)
}
