package env

import (
	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops"
	"strings"
	"testing"
)

var PLAIN = []byte(strings.TrimLeft(`
VAR1=val1
VAR2=val2
VAR3_unencrypted=val3
`, "\n"))

var BRANCH = sops.TreeBranch{
	sops.TreeItem{
		Key:   "VAR1",
		Value: "val1",
	},
	sops.TreeItem{
		Key:   "VAR2",
		Value: "val2",
	},
	sops.TreeItem{
		Key:   "VAR3_unencrypted",
		Value: "val3",
	},
}


func TestLoadEncryptedFile(t *testing.T) {
	// FIXME: Implementation?
}

func TestLoadPlainFile(t *testing.T) {
	branch, err := (&Store{}).LoadPlainFile(PLAIN)
	assert.Nil(t, err)
	assert.Equal(t, BRANCH, branch)
}

func TestEmitEncryptedFile(t *testing.T) {
	// FIXME: Implementation?
}

func TestEmitPlainFile(t *testing.T) {
	bytes, err := (&Store{}).EmitPlainFile(BRANCH)
	assert.Nil(t, err)
	assert.Equal(t, PLAIN, bytes)
}

func TestEmitValue(t *testing.T) {
	// FIXME: Implementation?
}
