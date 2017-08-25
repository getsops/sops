package yaml

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
  - filename_regex: ""
    kms: foo
    pgp: bar
`)

var sampleConfigWithGroups = []byte(`
creation_rules:
  - filename_regex: foobar*
    kms: "1"
    pgp: "2"
  - filename_regex: ""
    key_groups:
    - kms: foo
      pgp: bar
    - kms: baz
      pgp: qux
`)

func TestLoadConfigFile(t *testing.T) {
	expected := configFile{
		CreationRules: []creationRule{
			creationRule{
				FilenameRegex: "foobar*",
				KMS:           "1",
				PGP:           "2",
			},
			creationRule{
				FilenameRegex: "",
				KMS:           "foo",
				PGP:           "bar",
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
						KMS: "foo",
						PGP: "bar",
					},
					{
						KMS: "baz",
						PGP: "qux",
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
	groups, err := KeyGroupsForFile("foobar2000", sampleConfig, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "2", groups[0][0].ToString())
	assert.Equal(t, "1", groups[0][1].ToString())
	groups, err = KeyGroupsForFile("whatever", sampleConfig, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", groups[0][0].ToString())
	assert.Equal(t, "foo", groups[0][1].ToString())
}

func TestKeyGroupsForFileWithGroups(t *testing.T) {
	groups, err := KeyGroupsForFile("whatever", sampleConfigWithGroups, nil)
	assert.Equal(t, nil, err)
	assert.Equal(t, "bar", groups[0][0].ToString())
	assert.Equal(t, "foo", groups[0][1].ToString())
	assert.Equal(t, "qux", groups[1][0].ToString())
	assert.Equal(t, "baz", groups[1][1].ToString())
}
