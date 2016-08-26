package yaml

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
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

func TestMasterKeyStringsForFile(t *testing.T) {
	kms, pgp, err := MasterKeyStringsForFile("foobar2000", sampleConfig)
	assert.Equal(t, nil, err)
	assert.Equal(t, "1", kms)
	assert.Equal(t, "2", pgp)
	kms, pgp, err = MasterKeyStringsForFile("whatever", sampleConfig)
	assert.Equal(t, nil, err)
	assert.Equal(t, "foo", kms)
	assert.Equal(t, "bar", pgp)
}
