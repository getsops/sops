package sops

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUnencryptedSuffix(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "foo_unencrypted",
			Value: "bar",
		},
	}
	tree := Tree{Branch: branch, Metadata: Metadata{UnencryptedSuffix: "_unencrypted"}}
	expected := TreeBranch{
		TreeItem{
			Key:   "foo_unencrypted",
			Value: "bar",
		},
	}
	_, err := tree.Encrypt(bytes.Repeat([]byte("f"), 32))
	if err != nil {
		t.Errorf("Encrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branch, expected) {
		t.Errorf("Trees don't match: \ngot \t\t%+v,\n expected \t\t%+v", tree.Branch, expected)
	}
	_, err = tree.Decrypt(bytes.Repeat([]byte("f"), 32))
	if err != nil {
		t.Errorf("Decrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branch, expected) {
		t.Errorf("Trees don't match: \ngot\t\t\t%+v,\nexpected\t\t%+v", tree.Branch, expected)
	}
}
