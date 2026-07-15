package sops_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
)

// TestMACWithCommentInSequence is a regression test for getsops/sops#2243.
//
// A comment inside a sequence is a Comment node in the plaintext tree, but
// Encrypt turns it into an ordinary (non-Comment) encrypted string whose
// metadata records that it is a comment; Decrypt restores it to a Comment.
// Encrypt excludes the comment from the MAC, so Decrypt must exclude it too,
// which requires inspecting the decrypted result's type, not the input's.
// Before the fix the MACs disagreed and decryption of any file with a comment
// inside a sequence failed with "MAC mismatch".
func TestMACWithCommentInSequence(t *testing.T) {
	branches := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{
				Key: "list",
				Value: []interface{}{
					sops.Comment{Value: "commented"},
					"real",
				},
			},
		},
	}
	tree := sops.Tree{Branches: branches, Metadata: sops.Metadata{}}
	cipher := aes.NewCipher()
	key := bytes.Repeat([]byte{'k'}, 32)

	// Encrypt computes the MAC and mutates the tree in place.
	encMAC, err := tree.Encrypt(key, cipher)
	assert.NoError(t, err)

	// The sequence comment is now an encrypted string, not a Comment node.
	seq := tree.Branches[0][0].Value.([]interface{})
	_, stillComment := seq[0].(sops.Comment)
	assert.False(t, stillComment)
	assert.IsType(t, "", seq[0])

	// Decrypting the same tree must reproduce the MAC Encrypt computed...
	decMAC, err := tree.Decrypt(key, cipher)
	assert.NoError(t, err)
	assert.Equal(t, encMAC, decMAC)

	// ...and restore the sequence element to a Comment with its content.
	seq = tree.Branches[0][0].Value.([]interface{})
	comment, ok := seq[0].(sops.Comment)
	assert.True(t, ok)
	assert.Equal(t, "commented", comment.Value)
}
