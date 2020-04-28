package config

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockFS struct {
	stat func(string) (os.FileInfo, error)
}

func (fs mockFS) Stat(name string) (os.FileInfo, error) {
	return fs.stat(name)
}

func TestFindConfigFileRecursive(t *testing.T) {
	expectedPath := path.Clean("./../../.sops.yaml")
	fs = mockFS{stat: func(name string) (os.FileInfo, error) {
		if name == expectedPath {
			return nil, nil
		}
		return nil, &os.PathError{}
	}}
	filepath, err := FindConfigFile(".")
	assert.Nil(t, err)
	assert.Equal(t, expectedPath, filepath)
}

func TestFindConfigFileCurrentDir(t *testing.T) {
	expectedPath := path.Clean(".sops.yaml")
	fs = mockFS{stat: func(name string) (os.FileInfo, error) {
		if name == expectedPath {
			return nil, nil
		}
		return nil, &os.PathError{}
	}}
	filepath, err := FindConfigFile(".")
	assert.Nil(t, err)
	assert.Equal(t, expectedPath, filepath)
}

var sampleConfig = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "1"
    pgp: "2"
    gcp_kms: "3"
    hc_vault_transit_uri: http://4:8200/v1/4/keys/4
  - path_regex: ""
    kms: foo
    pgp: bar
    gcp_kms: baz
    hc_vault_transit_uri: http://127.0.1.1/v1/baz/keys/baz
`)

var sampleConfigWithPath = []byte(`
creation_rules:
  - path_regex: foo/bar*
    kms: "1"
    pgp: "2"
    gcp_kms: "3"
    hc_vault_uris: http://4:8200/v1/4/keys/4
  - path_regex: somefilename.yml
    kms: bilbo
    pgp: baggins
    gcp_kms: precious
    hc_vault_uris: https://pluto/v1/pluto/keys/pluto
  - path_regex: ""
    kms: foo
    pgp: bar
    gcp_kms: baz
    hc_vault_uris: https://foz:443/v1/foz/keys/foz
`)

var sampleConfigWithGroups = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "1"
    pgp: "2"
  - path_regex: ""
    key_groups:
    - kms:
      - arn: foo
      pgp:
      - bar
      gcp_kms:
      - resource_id: foo
      azure_keyvault:
      - vaultUrl: https://foo.vault.azure.net
        key: foo-key
        version: fooversion
      hc_vault:
      - 'https://foo.vault:8200/v1/foo/keys/foo-key'
    - kms:
      - arn: baz
      pgp:
      - qux
      gcp_kms:
      - resource_id: bar
      - resource_id: baz
      azure_keyvault:
      - vaultUrl: https://bar.vault.azure.net
        key: bar-key
        version: barversion
      hc_vault:
      - 'https://baz.vault:8200/v1/baz/keys/baz-key'
`)

var sampleConfigWithSuffixParameters = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "1"
    pgp: "2"
    unencrypted_suffix: _unencrypted
  - path_regex: bar?foo$
    encrypted_suffix: _enc
    key_groups:
      - kms:
          - arn: baz
        pgp:
          - qux
        gcp_kms:
          - resource_id: bar
          - resource_id: baz
        azure_keyvault:
        - vaultUrl: https://foo.vault.azure.net
          key: foo-key
          version: fooversion
    `)

var sampleConfigWithRegexParameters = []byte(`
creation_rules:
  - path_regex: barbar*
    kms: "1"
    pgp: "2"
    encrypted_regex: "^enc:"
    `)

var sampleConfigWithInvalidParameters = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "1"
    pgp: "2"
    hc_vault_uris: "https://vault.com/v1/bug/keys/pr"
    unencrypted_suffix: _unencrypted
    encrypted_suffix: _enc
    `)

var sampleConfigWithNoMatchingRules = []byte(`
creation_rules:
  - path_regex: notexisting
    pgp: bar
`)

var sampleEmptyConfig = []byte(``)

var sampleConfigWithEmptyCreationRules = []byte(`
creation_rules:
`)

var sampleConfigWithOnlyDestinationRules = []byte(`
destination_rules:
  - path_regex: ""
    s3_bucket: "foobar"
    s3_prefix: "test/"
    recreation_rule:
      pgp: newpgp
`)

var sampleConfigWithDestinationRule = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "1"
    pgp: "2"
    gcp_kms: "3"
  - path_regex: ""
    kms: foo
    pgp: bar
    gcp_kms: baz
destination_rules:
  - path_regex: ""
    s3_bucket: "foobar"
    s3_prefix: "test/"
    recreation_rule:
      pgp: newpgp
`)

var sampleConfigWithVaultDestinationRules = []byte(`
creation_rules:
  - path_regex: foobar*
    kms: "1"
    pgp: "2"
    gcp_kms: "3"
  - path_regex: ""
    kms: foo
    pgp: bar
    gcp_kms: baz
destination_rules:
  - vault_path: "foobar/"
    path_regex: "vault-v2/*"
  - vault_path: "barfoo/"
    vault_kv_mount_name: "kv/"
    vault_kv_version: 1
    path_regex: "vault-v1/*"
`)

func parseConfigFile(confBytes []byte, t *testing.T) *configFile {
	conf := &configFile{}
	err := conf.load(confBytes)
	assert.Nil(t, err)
	return conf
}

func TestLoadConfigFile(t *testing.T) {
	expected := configFile{
		CreationRules: []creationRule{
			{
				PathRegex: "foobar*",
				KMS:       "1",
				PGP:       "2",
				GCPKMS:    "3",
				VaultURI:  "http://4:8200/v1/4/keys/4",
			},
			{
				PathRegex: "",
				KMS:       "foo",
				PGP:       "bar",
				GCPKMS:    "baz",
				VaultURI:  "http://127.0.1.1/v1/baz/keys/baz",
			},
		},
	}

	conf := configFile{}
	err := conf.load(sampleConfig)
	assert.Nil(t, err)
	assert.Equal(t, expected, conf)
}

func TestLoadConfigFileWithGroups(t *testing.T) {
	expected := configFile{
		CreationRules: []creationRule{
			{
				PathRegex: "foobar*",
				KMS:       "1",
				PGP:       "2",
			},
			{
				PathRegex: "",
				KeyGroups: []keyGroup{
					{
						KMS:     []kmsKey{{Arn: "foo"}},
						PGP:     []string{"bar"},
						GCPKMS:  []gcpKmsKey{{ResourceID: "foo"}},
						AzureKV: []azureKVKey{{VaultURL: "https://foo.vault.azure.net", Key: "foo-key", Version: "fooversion"}},
						Vault:   []string{"https://foo.vault:8200/v1/foo/keys/foo-key"},
					},
					{
						KMS: []kmsKey{{Arn: "baz"}},
						PGP: []string{"qux"},
						GCPKMS: []gcpKmsKey{
							{ResourceID: "bar"},
							{ResourceID: "baz"},
						},
						AzureKV: []azureKVKey{{VaultURL: "https://bar.vault.azure.net", Key: "bar-key", Version: "barversion"}},
						Vault:   []string{"https://baz.vault:8200/v1/baz/keys/baz-key"},
					},
				},
			},
		},
	}

	conf := configFile{}
	err := conf.load(sampleConfigWithGroups)
	assert.Nil(t, err)
	assert.Equal(t, expected, conf)
}

func TestLoadConfigFileWithNoMatchingRules(t *testing.T) {
	_, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithNoMatchingRules, t), "foobar2000", nil)
	assert.NotNil(t, err)
}

func TestLoadEmptyConfigFile(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleEmptyConfig, t), "foobar2000", nil)
	assert.Nil(t, conf)
	assert.Nil(t, err)
}

func TestLoadConfigFileWithEmptyCreationRules(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithEmptyCreationRules, t), "foobar2000", nil)
	assert.Nil(t, conf)
	assert.Nil(t, err)
}

func TestLoadConfigFileWithOnlyDestinationRules(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithOnlyDestinationRules, t), "foobar2000", nil)
	assert.Nil(t, conf)
	assert.Nil(t, err)
}

func TestKeyGroupsForFile(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfig, t), "foobar2000", nil)
	assert.Nil(t, err)
	assert.Equal(t, "2", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "1", conf.KeyGroups[0][1].ToString())
	conf, err = parseCreationRuleForFile(parseConfigFile(sampleConfig, t), "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
}

func TestKeyGroupsForFileWithPath(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithPath, t), "foo/bar2000", nil)
	assert.Nil(t, err)
	assert.Equal(t, "2", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "1", conf.KeyGroups[0][1].ToString())
	conf, err = parseCreationRuleForFile(parseConfigFile(sampleConfigWithPath, t), "somefilename.yml", nil)
	assert.Nil(t, err)
	assert.Equal(t, "baggins", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "bilbo", conf.KeyGroups[0][1].ToString())
	conf, err = parseCreationRuleForFile(parseConfigFile(sampleConfig, t), "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
}

func TestKeyGroupsForFileWithGroups(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithGroups, t), "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
	assert.Equal(t, "qux", conf.KeyGroups[1][0].ToString())
	assert.Equal(t, "baz", conf.KeyGroups[1][1].ToString())
}

func TestLoadConfigFileWithUnencryptedSuffix(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithSuffixParameters, t), "foobar", nil)
	assert.Nil(t, err)
	assert.Equal(t, "_unencrypted", conf.UnencryptedSuffix)
}

func TestLoadConfigFileWithEncryptedSuffix(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithSuffixParameters, t), "barfoo", nil)
	assert.Nil(t, err)
	assert.Equal(t, "_enc", conf.EncryptedSuffix)
}

func TestLoadConfigFileWithEncryptedRegex(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithRegexParameters, t), "barbar", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "^enc:", conf.EncryptedRegex)
}

func TestLoadConfigFileWithInvalidParameters(t *testing.T) {
	_, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithInvalidParameters, t), "foobar", nil)
	assert.NotNil(t, err)
}

func TestLoadConfigFileWithDestinationRule(t *testing.T) {
	conf, err := parseDestinationRuleForFile(parseConfigFile(sampleConfigWithDestinationRule, t), "barfoo", nil)
	assert.Nil(t, err)
	assert.Equal(t, "newpgp", conf.KeyGroups[0][0].ToString())
	assert.NotNil(t, conf.Destination)
	assert.Equal(t, "s3://foobar/test/barfoo", conf.Destination.Path("barfoo"))
}

func TestLoadConfigFileWithVaultDestinationRules(t *testing.T) {
	conf, err := parseDestinationRuleForFile(parseConfigFile(sampleConfigWithVaultDestinationRules, t), "vault-v2/barfoo", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Equal(t, "http://127.0.0.1:8200/v1/secret/data/foobar/barfoo", conf.Destination.Path("barfoo"))
	conf, err = parseDestinationRuleForFile(parseConfigFile(sampleConfigWithVaultDestinationRules, t), "vault-v1/barfoo", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Equal(t, "http://127.0.0.1:8200/v1/kv/barfoo/barfoo", conf.Destination.Path("barfoo"))
}
