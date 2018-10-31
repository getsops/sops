package dotenv

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops"
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

func TestLoadPlainFile(t *testing.T) {
	branch, err := (&Store{}).LoadPlainFile(PLAIN)
	assert.Nil(t, err)
	assert.Equal(t, BRANCH, branch)
}
func TestEmitPlainFile(t *testing.T) {
	bytes, err := (&Store{}).EmitPlainFile(BRANCH)
	assert.Nil(t, err)
	assert.Equal(t, PLAIN, bytes)
}

func TestEmitValueString(t *testing.T) {
	bytes, err := (&Store{}).EmitValue("hello")
	assert.Nil(t, err)
	assert.Equal(t, []byte("hello"), bytes)
}

func TestEmitValueNonstring(t *testing.T) {
	_, err := (&Store{}).EmitValue(BRANCH)
	assert.NotNil(t, err)
}
