package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// matchFromBytes writes the given .sops.yaml contents to a temp dir and
// returns the MatchResult for an absolute file path. The file path can be
// outside the temp dir — in that case the matcher falls back to absolute-
// path matching (which is fine for catch-all creation rules with no
// path_regex). Tests that exercise path normalization should anchor the
// file path inside the temp dir.
func matchFromBytes(t *testing.T, confBytes []byte, absFilePath string) *config.MatchResult {
	t.Helper()
	dir := t.TempDir()
	confPath := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(confPath, confBytes, 0644))
	mr, err := config.MatchRulesForFile(confPath, absFilePath)
	require.NoError(t, err)
	return mr
}

func TestBuildCreationRuleView_StringKMS(t *testing.T) {
	confBytes := []byte(`
creation_rules:
  - kms: 'arn:aws:kms:us-east-1:1:key/short'
`)
	mr := matchFromBytes(t, confBytes, "/conf/x.yaml")
	require.NotNil(t, mr.CreationRule)

	view, err := buildCreationRuleView(mr.CreationRule)
	require.NoError(t, err)
	assert.Equal(t, 0, view.RuleIndex)
	require.Len(t, view.KMS, 1)
	assert.Equal(t, "arn:aws:kms:us-east-1:1:key/short", view.KMS[0].Arn)
	// Other recipient fields stay empty/omitted.
	assert.Empty(t, view.Age)
	assert.Empty(t, view.PGP)
}

func TestBuildCreationRuleView_KeyGroups(t *testing.T) {
	confBytes := []byte(`
creation_rules:
  - shamir_threshold: 2
    key_groups:
      - kms:
          - arn: 'arn:group1'
            role: 'arn:role/r1'
        pgp:
          - 'FP1'
      - age:
          - 'age1g2'
`)
	mr := matchFromBytes(t, confBytes, "/conf/x.yaml")
	require.NotNil(t, mr.CreationRule)

	view, err := buildCreationRuleView(mr.CreationRule)
	require.NoError(t, err)
	assert.Equal(t, 2, view.ShamirThreshold)
	require.Len(t, view.KeyGroups, 2)
	require.Len(t, view.KeyGroups[0].KMS, 1)
	assert.Equal(t, "arn:group1", view.KeyGroups[0].KMS[0].Arn)
	assert.Equal(t, "arn:role/r1", view.KeyGroups[0].KMS[0].Role)
	assert.Equal(t, []string{"FP1"}, view.KeyGroups[0].PGP)
	require.Len(t, view.KeyGroups[1].Age, 1)
	// Flat recipient fields should be omitted (omitempty) when key_groups populated.
	assert.Empty(t, view.KMS)
	assert.Empty(t, view.Age)
	assert.Empty(t, view.PGP)
}

func TestBuildCreationRuleView_KMSObjectForm(t *testing.T) {
	confBytes := []byte(`
creation_rules:
  - kms:
      - arn: 'arn:rich'
        role: 'arn:r'
        aws_profile: 'prod'
        context:
          team: 'payments'
`)
	mr := matchFromBytes(t, confBytes, "/conf/x.yaml")
	view, err := buildCreationRuleView(mr.CreationRule)
	require.NoError(t, err)
	require.Len(t, view.KMS, 1)
	assert.Equal(t, "arn:rich", view.KMS[0].Arn)
	assert.Equal(t, "arn:r", view.KMS[0].Role)
	assert.Equal(t, "prod", view.KMS[0].AwsProfile)
	require.NotNil(t, view.KMS[0].Context)
	assert.Equal(t, "payments", *view.KMS[0].Context["team"])
}

func TestBuildDestinationRuleView_S3(t *testing.T) {
	confBytes := []byte(`
destination_rules:
  - path_regex: 'publish/.*\.json$'
    s3_bucket: 'org-secrets'
    s3_prefix: 'sops/'
    omit_extensions: true
    recreation_rule:
      kms: 'arn:recreation'
`)
	mr := matchFromBytes(t, confBytes, "/conf/publish/app.json")
	require.NotNil(t, mr.DestinationRule)

	view, err := buildDestinationRuleView(mr.DestinationRule)
	require.NoError(t, err)
	assert.Equal(t, 0, view.RuleIndex)
	assert.Equal(t, `publish/.*\.json$`, view.PathRegex)
	assert.True(t, view.OmitExtensions)
	require.NotNil(t, view.Destination)
	assert.Equal(t, "s3", view.Destination.Type)
	assert.Equal(t, "org-secrets", view.Destination.Bucket)
	assert.Equal(t, "sops/", view.Destination.Prefix)
	require.NotNil(t, view.RecreationRule)
	require.Len(t, view.RecreationRule.KMS, 1)
	assert.Equal(t, "arn:recreation", view.RecreationRule.KMS[0].Arn)
}

func TestBuildDestinationRuleView_VaultNoRecreation(t *testing.T) {
	confBytes := []byte(`
destination_rules:
  - path_regex: 'publish/.*'
    vault_path: 'secret/data/app'
    vault_address: 'https://v.example.com'
`)
	mr := matchFromBytes(t, confBytes, "/conf/publish/app.yaml")
	require.NotNil(t, mr.DestinationRule)

	view, err := buildDestinationRuleView(mr.DestinationRule)
	require.NoError(t, err)
	require.NotNil(t, view.Destination)
	assert.Equal(t, "vault", view.Destination.Type)
	assert.Equal(t, "https://v.example.com", view.Destination.Address)
	assert.Equal(t, "secret/data/app", view.Destination.Path)
	// No recreation_rule in source → nil.
	assert.Nil(t, view.RecreationRule)
}

func TestRun_NoMatchEmptyArrays(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(confPath, []byte(`creation_rules: []`), 0644))

	opts := Opts{
		ConfigPath: confPath,
		FilePath:   filepath.Join(dir, "anything.yaml"),
	}
	out, exitCode, err := Run(opts)
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
	require.NotNil(t, out)
	assert.Equal(t, 1, out.SchemaVersion)
	// Critical: not nil, not "null" in JSON.
	assert.NotNil(t, out.CreationRules)
	assert.NotNil(t, out.DestinationRules)
	assert.Empty(t, out.CreationRules)
	assert.Empty(t, out.DestinationRules)

	// Marshal and verify the JSON shape directly.
	b, err := json.Marshal(out)
	require.NoError(t, err)
	assert.Contains(t, string(b), `"creation_rules":[]`)
	assert.Contains(t, string(b), `"destination_rules":[]`)
	assert.NotContains(t, string(b), `"creation_rules":null`)
}

func TestRun_RequireMatchExitsNonZero(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(confPath, []byte(`creation_rules: []`), 0644))

	opts := Opts{
		ConfigPath:   confPath,
		FilePath:     filepath.Join(dir, "anything.yaml"),
		RequireMatch: true,
	}
	out, exitCode, err := Run(opts)
	// JSON is still produced even on the error path.
	require.NotNil(t, out)
	assert.Equal(t, codes.NoRulesMatched, exitCode)
	require.Error(t, err)
}

func TestBuildCreationRuleView_SkipsFlatWhenKeyGroupsPresent(t *testing.T) {
	// sops uses key_groups XOR flat fields. When both are written, the flat
	// fields are ignored by the encryption pipeline. The view must reflect
	// that to avoid misleading users about "which keys will encrypt this".
	confBytes := []byte(`
creation_rules:
  - kms: 'arn:flat-IGNORED-BY-SOPS'
    pgp: 'IGNORED-PGP'
    age: 'age1IGNORED'
    key_groups:
      - kms: [{ arn: 'arn:effective' }]
        pgp: ['REALPGP']
`)
	mr := matchFromBytes(t, confBytes, "/conf/x.yaml")
	view, err := buildCreationRuleView(mr.CreationRule)
	require.NoError(t, err)
	require.Len(t, view.KeyGroups, 1)
	assert.Equal(t, "arn:effective", view.KeyGroups[0].KMS[0].Arn)
	// Flat fields are dead per sops's parser; the view should not emit them.
	assert.Empty(t, view.KMS)
	assert.Empty(t, view.PGP)
	assert.Empty(t, view.Age)
}

func TestBuildCreationRuleView_AzureFlatURLParsed(t *testing.T) {
	// Flat azure_keyvault URLs should split into vaultUrl/key/version,
	// matching how azkv.NewMasterKeyFromURL parses them at runtime.
	confBytes := []byte(`
creation_rules:
  - azure_keyvault: 'https://akv1.vault.azure.net/keys/k1/v123'
`)
	mr := matchFromBytes(t, confBytes, "/conf/x.yaml")
	view, err := buildCreationRuleView(mr.CreationRule)
	require.NoError(t, err)
	require.Len(t, view.AzureKeyVault, 1)
	assert.Equal(t, "https://akv1.vault.azure.net", view.AzureKeyVault[0].VaultURL)
	assert.Equal(t, "k1", view.AzureKeyVault[0].Key)
	assert.Equal(t, "v123", view.AzureKeyVault[0].Version)
}

func TestBuildCreationRuleView_AzureFlatURLWithoutVersion(t *testing.T) {
	confBytes := []byte(`
creation_rules:
  - azure_keyvault: 'https://akv1.vault.azure.net/keys/k1'
`)
	mr := matchFromBytes(t, confBytes, "/conf/x.yaml")
	view, err := buildCreationRuleView(mr.CreationRule)
	require.NoError(t, err)
	require.Len(t, view.AzureKeyVault, 1)
	assert.Equal(t, "https://akv1.vault.azure.net", view.AzureKeyVault[0].VaultURL)
	assert.Equal(t, "k1", view.AzureKeyVault[0].Key)
	assert.Equal(t, "", view.AzureKeyVault[0].Version)
}

func TestRun_RequireMatchSucceedsOnMatch(t *testing.T) {
	dir := t.TempDir()
	confPath := filepath.Join(dir, ".sops.yaml")
	require.NoError(t, os.WriteFile(confPath, []byte(`
creation_rules:
  - kms: 'arn:catch'
`), 0644))

	opts := Opts{
		ConfigPath:   confPath,
		FilePath:     filepath.Join(dir, "x.yaml"),
		RequireMatch: true,
	}
	out, exitCode, err := Run(opts)
	require.NoError(t, err)
	assert.Equal(t, 0, exitCode)
	require.Len(t, out.CreationRules, 1)
}
