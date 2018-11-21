package yaml

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops"
)

var PLAIN = []byte(`---
# comment 0
key1: value
key1_a: value
# ^ comment 1
---
key2: value2`)

var BRANCHES = sops.TreeBranches{
	sops.TreeBranch{
		sops.TreeItem{
			Key:   "key1",
			Value: "value",
		},
		sops.TreeItem{
			Key:   "key1_a",
			Value: "value",
		},
		sops.TreeItem{
			Key:   sops.Comment{" ^ comment 1"},
			Value: nil,
		},
	},
	sops.TreeBranch{
		sops.TreeItem{
			Key:   "key2",
			Value: "value2",
		},
	},
}

func TestUnmarshalMetadataFromNonSOPSFile(t *testing.T) {
	data := []byte(`hello: 2`)
	_, err := (&Store{}).LoadEncryptedFile(data)
	assert.Equal(t, sops.MetadataNotFound, err)
}

func TestLoadPlainFile(t *testing.T) {
	branches, err := (&Store{}).LoadPlainFile(PLAIN)
	assert.Nil(t, err)
	assert.Equal(t, BRANCHES, branches)
}
