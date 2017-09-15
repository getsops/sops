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
	assert.Equal(t, nil, err)
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
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedPath, filepath)
}

var sampleConfig = []byte(`
creation_rules:
  - filename_regex: foobar*
    kms: "1"
    pgp: "2"
    gcp_kms: "3"
  - filename_regex: ""
    kms: foo
    pgp: bar
    gcp_kms: baz
`)

var sampleConfigWithGroups = []byte(`
creation_rules:
  - filename_regex: foobar*
    kms: "1"
    pgp: "2"
  - filename_regex: ""
    key_groups:
    - kms:
      - arn: foo
      pgp:
      - bar
      gcp_kms:
      - resource_id: foo
    - kms:
      - arn: baz
      pgp:
      - qux
      gcp_kms:
      - resource_id: bar
      - resource_id: baz
`)

func TestLoadConfigFile(t *testing.T) {
	expected := configFile{
		CreationRules: []creationRule{
			creationRule{
				FilenameRegex: "foobar*",
				KMS:           "1",
				PGP:           "2",
				GCPKMS:        "3",
			},
			creationRule{
				FilenameRegex: "",
				KMS:           "foo",
				PGP:           "bar",
				GCPKMS:        "baz",
			},
		},
	}

	conf := configFile{}
	err := conf.load(sampleConfig)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, conf)
}

func TestLoadConfigFileWithGroups(t *testing.T) {
	expected := configFile{
		CreationRules: []creationRule{
			{
				FilenameRegex: "foobar*",
				KMS:           "1",
				PGP:           "2",
			},
			{
				FilenameRegex: "",
				KeyGroups: []keyGroup{
					{
						KMS:    []kmsKey{{Arn: "foo"}},
						PGP:    []string{"bar"},
						GCPKMS: []gcpKmsKey{{ResourceID: "foo"}},
					},
					{
						KMS: []kmsKey{{Arn: "baz"}},
						PGP: []string{"qux"},
						GCPKMS: []gcpKmsKey{
							{ResourceID: "bar"},
							{ResourceID: "baz"},
						},
					},
				},
			},
		},
	}

	conf := configFile{}
	err := conf.load(sampleConfigWithGroups)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, conf)
}

func TestKeyGroupsForFile(t *testing.T) {
	conf, err := loadForFileFromBytes(sampleConfig, "foobar2000", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "2", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "1", conf.KeyGroups[0][1].ToString())
	conf, err = loadForFileFromBytes(sampleConfig, "whatever", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
}

func TestKeyGroupsForFileWithGroups(t *testing.T) {
	conf, err := loadForFileFromBytes(sampleConfigWithGroups, "whatever", nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", conf.KeyGroups[0][0].ToString())
	assert.Equal(t, "foo", conf.KeyGroups[0][1].ToString())
	assert.Equal(t, "qux", conf.KeyGroups[1][0].ToString())
	assert.Equal(t, "baz", conf.KeyGroups[1][1].ToString())
}
