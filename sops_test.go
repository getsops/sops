package sops

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops/aes"
	"go.mozilla.org/sops/kms"
	"reflect"
	"testing"
)

func TestUnencryptedSuffix(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "foo_unencrypted",
			Value: "bar",
		},
		TreeItem{
			Key: "bar_unencrypted",
			Value: TreeBranch{
				TreeItem{
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}
	tree := Tree{Branch: branch, Metadata: Metadata{UnencryptedSuffix: "_unencrypted"}}
	expected := TreeBranch{
		TreeItem{
			Key:   "foo_unencrypted",
			Value: "bar",
		},
		TreeItem{
			Key: "bar_unencrypted",
			Value: TreeBranch{
				TreeItem{
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}
	cipher := aes.Cipher{}
	_, err := tree.Encrypt(bytes.Repeat([]byte("f"), 32), cipher, nil)
	if err != nil {
		t.Errorf("Encrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branch, expected) {
		t.Errorf("Trees don't match: \ngot \t\t%+v,\n expected \t\t%+v", tree.Branch, expected)
	}
	_, err = tree.Decrypt(bytes.Repeat([]byte("f"), 32), cipher, nil)
	if err != nil {
		t.Errorf("Decrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branch, expected) {
		t.Errorf("Trees don't match: \ngot\t\t\t%+v,\nexpected\t\t%+v", tree.Branch, expected)
	}
}

type MockCipher struct{}

func (m MockCipher) Encrypt(value interface{}, key []byte, path string, stashValue interface{}) (string, error) {
	return "a", nil
}

func (m MockCipher) Decrypt(value string, key []byte, path string) (interface{}, interface{}, error) {
	return "a", nil, nil
}

func TestEncrypt(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: "bar",
		},
		TreeItem{
			Key: "baz",
			Value: TreeBranch{
				TreeItem{
					Key:   "bar",
					Value: 5,
				},
			},
		},
		TreeItem{
			Key:   "bar",
			Value: false,
		},
		TreeItem{
			Key:   "foobar",
			Value: 2.12,
		},
	}
	expected := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: "a",
		},
		TreeItem{
			Key: "baz",
			Value: TreeBranch{
				TreeItem{
					Key:   "bar",
					Value: "a",
				},
			},
		},
		TreeItem{
			Key:   "bar",
			Value: "a",
		},
		TreeItem{
			Key:   "foobar",
			Value: "a",
		},
	}
	tree := Tree{Branch: branch, Metadata: Metadata{UnencryptedSuffix: DefaultUnencryptedSuffix}}
	tree.Encrypt(bytes.Repeat([]byte{'f'}, 32), MockCipher{}, make(map[string][]interface{}))
	if !reflect.DeepEqual(tree.Branch, expected) {
		t.Errorf("%s does not equal expected tree: %s", tree.Branch, expected)
	}
}

func TestDecrypt(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: "bar",
		},
		TreeItem{
			Key: "baz",
			Value: TreeBranch{
				TreeItem{
					Key:   "bar",
					Value: "5",
				},
			},
		},
		TreeItem{
			Key:   "bar",
			Value: "false",
		},
		TreeItem{
			Key:   "foobar",
			Value: "2.12",
		},
	}
	expected := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: "a",
		},
		TreeItem{
			Key: "baz",
			Value: TreeBranch{
				TreeItem{
					Key:   "bar",
					Value: "a",
				},
			},
		},
		TreeItem{
			Key:   "bar",
			Value: "a",
		},
		TreeItem{
			Key:   "foobar",
			Value: "a",
		},
	}
	tree := Tree{Branch: branch, Metadata: Metadata{UnencryptedSuffix: DefaultUnencryptedSuffix}}
	tree.Decrypt(bytes.Repeat([]byte{'f'}, 32), MockCipher{}, make(map[string][]interface{}))
	if !reflect.DeepEqual(tree.Branch, expected) {
		t.Errorf("%s does not equal expected tree: %s", tree.Branch, expected)
	}
}

func TestTruncateTree(t *testing.T) {
	tree := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: 2,
		},
		TreeItem{
			Key: "bar",
			Value: TreeBranch{
				TreeItem{
					Key: "foobar",
					Value: []int{
						1,
						2,
						3,
						4,
					},
				},
			},
		},
	}
	expected := 3
	result, err := tree.Truncate(`["bar"]["foobar"][2]`)
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, result)
}

func TestRemoveMasterKeys(t *testing.T) {
	m := Metadata{
		KeySources: []KeySource{
			KeySource{
				Name: "kms",
				Keys: []MasterKey{
					&kms.MasterKey{
						Arn: "foo",
					}, &kms.MasterKey{
						Arn: "bar",
					},
					&kms.MasterKey{
						Arn: "foobar",
					},
				},
			},
		},
	}
	m.RemoveMasterKeys([]MasterKey{
		&kms.MasterKey{
			Arn: "bar",
		},
		&kms.MasterKey{
			Arn: "foobar",
		},
	})
	assert.Equal(t, []MasterKey{
		&kms.MasterKey{
			Arn: "foo",
		},
	}, m.KeySources[0].Keys)
}


func TestInsertOrReplaceValue(t *testing.T) {
	tree := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: 2,
		},
		TreeItem{
			Key: "bar",
			Value: TreeBranch{
				TreeItem{
					Key: "foobar",
					Value: []int{
						1,
						2,
						3,
						4,
					},
				},
			},
		},
	}
	tree = tree.InsertOrReplaceValue("foo", 57)
	assert.Equal(t, tree, TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: 57,
		},
		TreeItem{
			Key: "bar",
			Value: TreeBranch{
				TreeItem{
					Key: "foobar",
					Value: []int{
						1,
						2,
						3,
						4,
					},
				},
			},
		},
	})
	tree = tree.InsertOrReplaceValue("foobar", 100)
	assert.Equal(t, tree, TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: 57,
		},
		TreeItem{
			Key: "bar",
			Value: TreeBranch{
				TreeItem{
					Key: "foobar",
					Value: []int{
						1,
						2,
						3,
						4,
					},
				},
			},
		},
		TreeItem{
			Key: "foobar",
			Value: 100,
		},
	})
}
