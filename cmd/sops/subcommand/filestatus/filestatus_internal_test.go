package filestatus

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	// this is a generic YAML file **not** sops-encrypted
	testSopsNotEncryptedFile = `
---
hello: world
foo: bar
`
)

func TestFileStatus_checkEncrypted(t *testing.T) {
	exampleFilePath := "../../../../example.yaml"

	encrypted, err := cfs(exampleFilePath)
	require.Nil(t, err, "should not error")
	require.True(t, encrypted, "file should be reported as encrypted")
}

func TestFileStatus_checkPlain(t *testing.T) {
	exampleFilePath := "../../../../functional-tests/res/plainfile.yaml"

	encrypted, err := cfs(exampleFilePath)
	require.Nil(t, err, "should not error")
	require.False(t, encrypted, "file should be reported as encrypted")
}

func TestFileStatus_checkNoMAC(t *testing.T) {
	exampleFilePath := "../../../../functional-tests/res/no_mac.yaml"

	encrypted, err := cfs(exampleFilePath)
	require.Nil(t, err, "should not error")
	require.False(t, encrypted, "file should be reported as encrypted")
}
