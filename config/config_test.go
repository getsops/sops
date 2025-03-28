package config

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/getsops/sops/v3/keys"
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

var sampleConfigWithAmbiguousPath = []byte(`
creation_rules:
  - path_regex: foo/*
    kms: "1"
    pgp: "2"
    gcp_kms: "3"
    hc_vault_uris: http://4:8200/v1/4/keys/4
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
        aws_profile: bar
      - arn: foo
        context:
          baz: bam
      - arn: foo
        aws_profile: bar
        context:
          baz: bam
      - arn: foo
        role: '123'
      - arn: foo
        aws_profile: bar
        context:
          baz: bam
        role: '123'
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
        aws_profile: foo
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

var sampleConfigWithMergeType = []byte(`
creation_rules:
  - path_regex: ""
    key_groups:
    # key00
    - hc_vault:
      - 'https://foo.vault:8200/v1/foo/keys/foo-key'
    - merge:
      - merge:
        - pgp:
          # key01
          - foo
          kms:
          # key02
          - arn: foo
            aws_profile: foo
          # key03
          - arn: foo
            aws_profile: bar
            context:
              baz: bam
            role: '123'
          gcp_kms:
          # key04
          - resource_id: foo
          azure_keyvault:
          # key05
          - vaultUrl: https://foo.vault.azure.net
            key: foo-key
            version: fooversion
          hc_vault:
          # key06
          - 'https://bar.vault:8200/v1/bar/keys/bar-key'
        - pgp:
          # key07
          - bar
          kms:
          # key08
          - arn: bar
            aws_profile: bar
          gcp_kms:
          # key09
          - resource_id: bar
          # key10
          - resource_id: baz
          azure_keyvault:
          # key11
          - vaultUrl: https://bar.vault.azure.net
            key: bar-key
            version: barversion
          hc_vault:
          # key12
          - 'https://baz.vault:8200/v1/baz/keys/baz-key'
        pgp:
        # key13
        - baz
        kms:
        # key14
        - arn: baz
          aws_profile: baz
        gcp_kms:
        # duplicate of key09
        - resource_id: bar
        azure_keyvault:
        # duplicate of key05
        - vaultUrl: https://foo.vault.azure.net
          key: foo-key
          version: fooversion
        hc_vault:
        # key15 (duplicate of key00, but that's in a different key_group)
        - 'https://foo.vault:8200/v1/foo/keys/foo-key'
      - pgp:
        # key16
        - qux
        kms:
        # key17
        - arn: qux
          aws_profile: qux
        # key18
        - arn: baz
          aws_profile: bar
        # key19
        - arn: baz
          role: '123'
        gcp_kms:
        # key20
        - resource_id: qux
        # key21
        - resource_id: fnord
        azure_keyvault:
        # key22
        - vaultUrl: https://baz.vault.azure.net
          key: baz-key
          version: bazversion
        hc_vault:
        # key23
        - 'https://qux.vault:8200/v1/qux/keys/qux-key'
      pgp:
      # duplicate of key07
      - bar
      kms:
      # duplicate of key08
      - arn: bar
        aws_profile: bar
      # key24
      - arn: fnord
        aws_profile: fnord
      # duplicate of key03
      - arn: foo
        aws_profile: bar
        context:
          baz: bam
        role: '123'
      gcp_kms:
      # duplicate of key09
      - resource_id: bar
      # duplicate of key21
      - resource_id: fnord
      azure_keyvault:
      # duplicate of key11
      - vaultUrl: https://bar.vault.azure.net
        key: bar-key
        version: barversion
      hc_vault:
      # duplicate of key12
      - 'https://baz.vault:8200/v1/baz/keys/baz-key'
      # key25
      - 'https://fnord.vault:8200/v1/fnord/keys/fnord-key'
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

var sampleConfigWithEncryptedRegexParameters = []byte(`
creation_rules:
  - path_regex: barbar*
    kms: "1"
    pgp: "2"
    encrypted_regex: "^enc:"
    `)

var sampleConfigWithUnencryptedRegexParameters = []byte(`
creation_rules:
  - path_regex: barbar*
    kms: "1"
    pgp: "2"
    unencrypted_regex: "^dec:"
    `)

var sampleConfigWithMACOnlyEncrypted = []byte(`
creation_rules:
  - path_regex: barbar*
    kms: "1"
    pgp: "2"
    mac_only_encrypted: true
    `)

var sampleConfigWithEncryptedCommentRegexParameters = []byte(`
creation_rules:
  - path_regex: barbar*
    kms: "1"
    pgp: "2"
    encrypted_comment_regex: "sops:enc"
    `)

var sampleConfigWithUnencryptedCommentRegexParameters = []byte(`
creation_rules:
  - path_regex: barbar*
    kms: "1"
    pgp: "2"
    unencrypted_comment_regex: "sops:dec"
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

var sampleConfigWithInvalidComplicatedRegexp = []byte(`
creation_rules:
  - path_regex: "[ ]\\K(?<!\\d )(?="
    kms: default
`)

var sampleConfigWithComplicatedRegexp = []byte(`
creation_rules:
  - path_regex: "stage/dev/feature-.*"
    kms: dev-feature
  - path_regex: "stage/dev/.*"
    kms: dev
  - path_regex: "stage/staging/.*"
    kms: staging
  - path_regex: "stage/.*/.*"
    kms: default
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
	bam := "bam"
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
						KMS: []kmsKey{
							{
								Arn:        "foo",
								AwsProfile: "bar",
							},
							{
								Arn: "foo",
								Context: map[string]*string{
									"baz": &bam,
								},
							},
							{
								Arn:        "foo",
								AwsProfile: "bar",
								Context: map[string]*string{
									"baz": &bam,
								},
							},
							{
								Arn:  "foo",
								Role: "123",
							},
							{
								Arn:        "foo",
								AwsProfile: "bar",
								Context: map[string]*string{
									"baz": &bam,
								},
								Role: "123",
							},
						},
						PGP:     []string{"bar"},
						GCPKMS:  []gcpKmsKey{{ResourceID: "foo"}},
						AzureKV: []azureKVKey{{VaultURL: "https://foo.vault.azure.net", Key: "foo-key", Version: "fooversion"}},
						Vault:   []string{"https://foo.vault:8200/v1/foo/keys/foo-key"},
					},
					{
						KMS: []kmsKey{{Arn: "baz", AwsProfile: "foo"}},
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

func id(key keys.MasterKey) string {
	return fmt.Sprintf("%s: %s", key.TypeToIdentifier(), key.ToString())
}

func ids(keys []keys.MasterKey) []string {
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, id(key))
	}
	return result
}

func TestLoadConfigFileWithMerge(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithMergeType, t), "/conf/path", "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(conf.KeyGroups))
	assert.Equal(t, []string{
		"hc_vault: https://foo.vault:8200/v1/foo/keys/foo-key",
	}, ids(conf.KeyGroups[0]))
	assert.Equal(t, []string{
		"pgp: foo",                 // key01
		"kms: foo||foo",            //key02
		"kms: foo+123|baz:bam|bar", //key03
		"gcp_kms: foo",             //key04
		"azure_kv: https://foo.vault.azure.net/keys/foo-key/fooversion", //key05
		"hc_vault: https://bar.vault:8200/v1/bar/keys/bar-key",          //key06
		"pgp: bar",      //key07
		"kms: bar||bar", //key08
		"gcp_kms: bar",  //key09
		"gcp_kms: baz",  //key10
		"azure_kv: https://bar.vault.azure.net/keys/bar-key/barversion", //key11
		"hc_vault: https://baz.vault:8200/v1/baz/keys/baz-key",          //key12
		"pgp: baz",      //key13
		"kms: baz||baz", //key14
		"hc_vault: https://foo.vault:8200/v1/foo/keys/foo-key", //key15
		"pgp: qux",       //key16
		"kms: qux||qux",  //key17
		"kms: baz||bar",  //key18
		"kms: baz+123",   //key19
		"gcp_kms: qux",   //key20
		"gcp_kms: fnord", //key21
		"azure_kv: https://baz.vault.azure.net/keys/baz-key/bazversion", //key22
		"hc_vault: https://qux.vault:8200/v1/qux/keys/qux-key",          //key23
		"kms: fnord||fnord", //key24
		"hc_vault: https://fnord.vault:8200/v1/fnord/keys/fnord-key", //key25
	}, ids(conf.KeyGroups[1]))
}

func TestLoadConfigFileWithNoMatchingRules(t *testing.T) {
	_, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithNoMatchingRules, t), "/conf/path", "foobar2000", nil)
	assert.NotNil(t, err)
}

func TestLoadConfigFileWithInvalidComplicatedRegexp(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithInvalidComplicatedRegexp, t), "/conf/path", "stage/prod/api.yml", nil)
	assert.Equal(t, "can not compile regexp: error parsing regexp: invalid escape sequence: `\\K`", err.Error())
	assert.Nil(t, conf)
}

func TestLoadConfigFileWithComplicatedRegexp(t *testing.T) {
	for filePath, k := range map[string]string{
		"stage/prod/api.yml":        "default",
		"stage/dev/feature-foo.yml": "dev-feature",
		"stage/dev/api.yml":         "dev",
	} {
		conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithComplicatedRegexp, t), "/conf/path", filePath, nil)
		assert.Nil(t, err)
		assert.Equal(t, k, conf.KeyGroups[0][0].ToString())
	}
}

func TestLoadEmptyConfigFile(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleEmptyConfig, t), "/conf/path", "foobar2000", nil)
	assert.Nil(t, conf)
	assert.Nil(t, err)
}

func TestLoadConfigFileWithEmptyCreationRules(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithEmptyCreationRules, t), "/conf/path", "foobar2000", nil)
	assert.Nil(t, conf)
	assert.Nil(t, err)
}

func TestLoadConfigFileWithOnlyDestinationRules(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithOnlyDestinationRules, t), "/conf/path", "foobar2000", nil)
	assert.Nil(t, conf)
	assert.Nil(t, err)
}

func TestKeyGroupsForFile(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfig, t), "/conf/path", "foobar2000", nil)
	assert.Nil(t, err)
	assert.Equal(t, "2", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "1", conf.KeyGroups[0][1].ToString())
	conf, err = parseCreationRuleForFile(parseConfigFile(sampleConfig, t), "/conf/path", "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
}

func TestKeyGroupsForFileWithPath(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithPath, t), "/conf/path", "foo/bar2000", nil)
	assert.Nil(t, err)
	assert.Equal(t, "2", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "1", conf.KeyGroups[0][1].ToString())
	conf, err = parseCreationRuleForFile(parseConfigFile(sampleConfigWithPath, t), "/conf/path", "somefilename.yml", nil)
	assert.Nil(t, err)
	assert.Equal(t, "baggins", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "bilbo", conf.KeyGroups[0][1].ToString())
	conf, err = parseCreationRuleForFile(parseConfigFile(sampleConfig, t), "/conf/path", "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
}

func TestKeyGroupsForFileWithGroups(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithGroups, t), "/conf/path", "whatever", nil)
	assert.Nil(t, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo||bar", conf.KeyGroups[0][1].ToString())
	assert.Equal(t, "foo|baz:bam", conf.KeyGroups[0][2].ToString())
	assert.Equal(t, "foo|baz:bam|bar", conf.KeyGroups[0][3].ToString())
	assert.Equal(t, "foo+123", conf.KeyGroups[0][4].ToString())
	assert.Equal(t, "foo+123|baz:bam|bar", conf.KeyGroups[0][5].ToString())
	assert.Equal(t, "qux", conf.KeyGroups[1][0].ToString())
	assert.Equal(t, "baz||foo", conf.KeyGroups[1][1].ToString())
}

func TestLoadConfigFileWithUnencryptedSuffix(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithSuffixParameters, t), "/conf/path", "foobar", nil)
	assert.Nil(t, err)
	assert.Equal(t, "_unencrypted", conf.UnencryptedSuffix)
}

func TestLoadConfigFileWithEncryptedSuffix(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithSuffixParameters, t), "/conf/path", "barfoo", nil)
	assert.Nil(t, err)
	assert.Equal(t, "_enc", conf.EncryptedSuffix)
}

func TestLoadConfigFileWithUnencryptedRegex(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithUnencryptedRegexParameters, t), "/conf/path", "barbar", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "^dec:", conf.UnencryptedRegex)
}

func TestLoadConfigFileWithEncryptedRegex(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithEncryptedRegexParameters, t), "/conf/path", "barbar", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "^enc:", conf.EncryptedRegex)
}

func TestLoadConfigFileWithMACOnlyEncrypted(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithMACOnlyEncrypted, t), "/conf/path", "barbar", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, conf.MACOnlyEncrypted)
}

func TestLoadConfigFileWithUnencryptedCommentRegex(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithUnencryptedCommentRegexParameters, t), "/conf/path", "barbar", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "sops:dec", conf.UnencryptedCommentRegex)
}

func TestLoadConfigFileWithEncryptedCommentRegex(t *testing.T) {
	conf, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithEncryptedCommentRegexParameters, t), "/conf/path", "barbar", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "sops:enc", conf.EncryptedCommentRegex)
}

func TestLoadConfigFileWithInvalidParameters(t *testing.T) {
	_, err := parseCreationRuleForFile(parseConfigFile(sampleConfigWithInvalidParameters, t), "/conf/path", "foobar", nil)
	assert.NotNil(t, err)
}

func TestLoadConfigFileWithAmbiguousPath(t *testing.T) {
	config := parseConfigFile(sampleConfigWithAmbiguousPath, t)
	_, err := parseCreationRuleForFile(config, "/foo/config", "/foo/foo/bar", nil)
	assert.Nil(t, err)
	_, err = parseCreationRuleForFile(config, "/foo/config", "/foo/fuu/bar", nil)
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
	assert.Contains(t, conf.Destination.Path("barfoo"), "/v1/secret/data/foobar/barfoo")
	conf, err = parseDestinationRuleForFile(parseConfigFile(sampleConfigWithVaultDestinationRules, t), "vault-v1/barfoo", nil)
	assert.Nil(t, err)
	assert.NotNil(t, conf.Destination)
	assert.Contains(t, conf.Destination.Path("barfoo"), "/v1/kv/barfoo/barfoo")
}
