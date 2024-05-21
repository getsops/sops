package dotenv

import (
	"strings"
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/stretchr/testify/assert"
)

var PLAIN = []byte(strings.TrimLeft(`
VAR1=val1
VAR2=val2
#comment
VAR3_unencrypted=val3
VAR4=val4\nval4
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
		Key:   sops.Comment{Value: "comment"},
		Value: nil,
	},
	sops.TreeItem{
		Key:   "VAR3_unencrypted",
		Value: "val3",
	},
	sops.TreeItem{
		Key:   "VAR4",
		Value: "val4\nval4",
	},
}

func TestLoadPlainFile(t *testing.T) {
	branches, err := (&Store{}).LoadPlainFile(PLAIN)
	assert.Nil(t, err)
	assert.Equal(t, BRANCH, branches[0])
}
func TestEmitPlainFile(t *testing.T) {
	branches := sops.TreeBranches{
		BRANCH,
	}
	bytes, err := (&Store{}).EmitPlainFile(branches)
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

func TestEmitEncryptedFileStability(t *testing.T) {
	// emit the same tree multiple times to ensure the output is stable
	// i.e. emitting the same tree always yields exactly the same output
	var previous []byte
	for i := 0; i < 10; i += 1 {
		bytes, err := (&Store{}).EmitEncryptedFile(sops.Tree{
			Branches: []sops.TreeBranch{{}},
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, bytes)
		if previous != nil {
			assert.Equal(t, previous, bytes)
		}
		previous = bytes
	}
}

func TestHasSopsTopLevelKey(t *testing.T) {
	ok := (&Store{}).HasSopsTopLevelKey(sops.TreeBranch{
		sops.TreeItem{
			Key:   "sops",
			Value: "value",
		},
	})
	assert.Equal(t, ok, false)
	ok = (&Store{}).HasSopsTopLevelKey(sops.TreeBranch{
		sops.TreeItem{
			Key:   "sops_",
			Value: "value",
		},
	})
	assert.Equal(t, ok, true)
}
