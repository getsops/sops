package sops

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type reverseCipher struct{}

// reverse returns its argument string reversed rune-wise left to right.
func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func (c reverseCipher) Encrypt(value interface{}, key []byte, path string) (string, error) {
	b, err := ToBytes(value)
	if err != nil {
		return "", err
	}
	return reverse(string(b)), nil
}
func (c reverseCipher) Decrypt(value string, key []byte, path string) (plaintext interface{}, err error) {
	if value == "error" {
		return nil, fmt.Errorf("Error")
	}
	return reverse(value), nil
}

func TestUnencryptedSuffix(t *testing.T) {
	branches := TreeBranches{
		TreeBranch{
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
		},
	}
	tree := Tree{Branches: branches, Metadata: Metadata{UnencryptedSuffix: "_unencrypted"}}
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
	cipher := reverseCipher{}
	_, err := tree.Encrypt(bytes.Repeat([]byte("f"), 32), cipher)
	if err != nil {
		t.Errorf("Encrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branches[0], expected) {
		t.Errorf("Trees don't match: \ngot \t\t%+v,\n expected \t\t%+v", tree.Branches[0], expected)
	}
	_, err = tree.Decrypt(bytes.Repeat([]byte("f"), 32), cipher)
	if err != nil {
		t.Errorf("Decrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branches[0], expected) {
		t.Errorf("Trees don't match: \ngot\t\t\t%+v,\nexpected\t\t%+v", tree.Branches[0], expected)
	}
}

func TestEncryptedSuffix(t *testing.T) {
	branches := TreeBranches{
		TreeBranch{
			TreeItem{
				Key:   "foo_encrypted",
				Value: "bar",
			},
			TreeItem{
				Key: "bar",
				Value: TreeBranch{
					TreeItem{
						Key:   "foo",
						Value: "bar",
					},
				},
			},
		},
	}
	tree := Tree{Branches: branches, Metadata: Metadata{EncryptedSuffix: "_encrypted"}}
	expected := TreeBranch{
		TreeItem{
			Key:   "foo_encrypted",
			Value: "rab",
		},
		TreeItem{
			Key: "bar",
			Value: TreeBranch{
				TreeItem{
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}
	cipher := reverseCipher{}
	_, err := tree.Encrypt(bytes.Repeat([]byte("f"), 32), cipher)
	if err != nil {
		t.Errorf("Encrypting the tree failed: %s", err)
	}
	if !reflect.DeepEqual(tree.Branches[0], expected) {
		t.Errorf("Trees don't match: \ngot \t\t%+v,\n expected \t\t%+v", tree.Branches[0], expected)
	}
	_, err = tree.Decrypt(bytes.Repeat([]byte("f"), 32), cipher)
	if err != nil {
		t.Errorf("Decrypting the tree failed: %s", err)
	}
	expected[0].Value = "bar"
	if !reflect.DeepEqual(tree.Branches[0], expected) {
		t.Errorf("Trees don't match: \ngot\t\t\t%+v,\nexpected\t\t%+v", tree.Branches[0], expected)
	}
}

type MockCipher struct{}

func (m MockCipher) Encrypt(value interface{}, key []byte, path string) (string, error) {
	return "a", nil
}

func (m MockCipher) Decrypt(value string, key []byte, path string) (interface{}, error) {
	return "a", nil
}

func TestEncrypt(t *testing.T) {
	branches := TreeBranches{
		TreeBranch{
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
			TreeItem{
				Key:   "barfoo",
				Value: nil,
			},
		},
		TreeBranch{
			TreeItem{
				Key:   "foo2",
				Value: "bar",
			},
		},
		TreeBranch{
			TreeItem{
				Key:   "foo3",
				Value: "bar",
			},
		},
	}
	expected := TreeBranches{
		TreeBranch{
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
			TreeItem{
				Key:   "barfoo",
				Value: nil,
			},
		},
		TreeBranch{
			TreeItem{
				Key:   "foo2",
				Value: "a",
			},
		},
		TreeBranch{
			TreeItem{
				Key:   "foo3",
				Value: "a",
			},
		},
	}
	tree := Tree{Branches: branches, Metadata: Metadata{UnencryptedSuffix: DefaultUnencryptedSuffix}}
	tree.Encrypt(bytes.Repeat([]byte{'f'}, 32), MockCipher{})
	if !reflect.DeepEqual(tree.Branches, expected) {
		t.Errorf("%s does not equal expected tree: %s", tree.Branches, expected)
	}
}

func TestDecrypt(t *testing.T) {
	branches := TreeBranches{
		TreeBranch{
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
			TreeItem{
				Key:   "barfoo",
				Value: nil,
			},
		},
		TreeBranch{
			TreeItem{
				Key:   "foo",
				Value: "bar",
			},
			TreeItem{
				Key: "baz",
				Value: TreeBranch{
					TreeItem{
						Key:   "bar",
						Value: "6",
					},
				},
			},
		},
		TreeBranch{
			TreeItem{
				Key:   "foo3",
				Value: "bar",
			},
		},
	}
	expected := TreeBranches{
		TreeBranch{
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
			TreeItem{
				Key:   "barfoo",
				Value: nil,
			},
		},
		TreeBranch{
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
		},
		TreeBranch{
			TreeItem{
				Key:   "foo3",
				Value: "a",
			},
		},
	}
	tree := Tree{Branches: branches, Metadata: Metadata{UnencryptedSuffix: DefaultUnencryptedSuffix}}
	tree.Decrypt(bytes.Repeat([]byte{'f'}, 32), MockCipher{})
	if !reflect.DeepEqual(tree.Branches, expected) {
		t.Errorf("%s does not equal expected tree: %s", tree.Branches[0], expected)
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
	result, err := tree.Truncate([]interface{}{
		"bar",
		"foobar",
		2,
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, result)
}

func TestEncryptComments(t *testing.T) {
	tree := Tree{
		Branches: TreeBranches{
			TreeBranch{
				TreeItem{
					Key:   Comment{"foo"},
					Value: nil,
				},
				TreeItem{
					Key: "list",
					Value: []interface{}{
						"1",
						Comment{"bar"},
						"2",
					},
				},
			},
		},
		Metadata: Metadata{
			UnencryptedSuffix: DefaultUnencryptedSuffix,
		},
	}
	tree.Encrypt(bytes.Repeat([]byte{'f'}, 32), reverseCipher{})
	assert.Equal(t, "oof", tree.Branches[0][0].Key.(Comment).Value)
	assert.Equal(t, "rab", tree.Branches[0][1].Value.([]interface{})[1])
}

func TestDecryptComments(t *testing.T) {
	tree := Tree{
		Branches: TreeBranches{
			TreeBranch{
				TreeItem{
					Key:   Comment{"oof"},
					Value: nil,
				},
				TreeItem{
					Key: "list",
					Value: []interface{}{
						"1",
						Comment{"rab"},
						"2",
					},
				},
				TreeItem{
					Key:   "list",
					Value: nil,
				},
			},
		},
		Metadata: Metadata{
			UnencryptedSuffix: DefaultUnencryptedSuffix,
		},
	}
	tree.Decrypt(bytes.Repeat([]byte{'f'}, 32), reverseCipher{})
	assert.Equal(t, "foo", tree.Branches[0][0].Key.(Comment).Value)
	assert.Equal(t, "bar", tree.Branches[0][1].Value.([]interface{})[1])
}

func TestDecryptUnencryptedComments(t *testing.T) {
	tree := Tree{
		Branches: TreeBranches{
			TreeBranch{
				TreeItem{
					// We use `error` to simulate an error decrypting, the fake cipher will error in this case
					Key:   Comment{"error"},
					Value: nil,
				},
			},
		},
		Metadata: Metadata{},
	}
	tree.Decrypt(bytes.Repeat([]byte{'f'}, 32), reverseCipher{})
	assert.Equal(t, "error", tree.Branches[0][0].Key.(Comment).Value)
}

func TestSetNewKey(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key: "foo",
			Value: TreeBranch{
				TreeItem{
					Key: "bar",
					Value: TreeBranch{
						TreeItem{
							Key:   "baz",
							Value: "foobar",
						},
					},
				},
			},
		},
	}
	set := branch.Set([]interface{}{"foo", "bar", "foo"}, "hello")
	assert.Equal(t, "hello", set[0].Value.(TreeBranch)[0].Value.(TreeBranch)[1].Value)
}

func TestSetArrayDeepNew(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key: "foo",
			Value: []interface{}{
				"one",
				"two",
			},
		},
	}
	set := branch.Set([]interface{}{"foo", 2, "bar"}, "hello")
	assert.Equal(t, "hello", set[0].Value.([]interface{})[2].(TreeBranch)[0].Value)
}

func TestSetNewKeyDeep(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: "bar",
		},
	}
	set := branch.Set([]interface{}{"foo", "bar", "baz"}, "hello")
	assert.Equal(t, "hello", set[0].Value.(TreeBranch)[0].Value.(TreeBranch)[0].Value)
}

func TestSetNewKeyOnEmptyBranch(t *testing.T) {
	branch := TreeBranch{}
	set := branch.Set([]interface{}{"foo", "bar", "baz"}, "hello")
	assert.Equal(t, "hello", set[0].Value.(TreeBranch)[0].Value.(TreeBranch)[0].Value)
}

func TestSetArray(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key: "foo",
			Value: []interface{}{
				"one",
				"two",
				"three",
			},
		},
	}
	set := branch.Set([]interface{}{"foo", 0}, "uno")
	assert.Equal(t, "uno", set[0].Value.([]interface{})[0])
}

func TestSetArrayNew(t *testing.T) {
	branch := TreeBranch{}
	set := branch.Set([]interface{}{"foo", 0, 0}, "uno")
	assert.Equal(t, "uno", set[0].Value.([]interface{})[0].([]interface{})[0])
}

func TestSetExisting(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "foo",
			Value: "foobar",
		},
	}
	set := branch.Set([]interface{}{"foo"}, "bar")
	assert.Equal(t, "bar", set[0].Value)
}

func TestSetArrayLeafNewItem(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key:   "array",
			Value: []interface{}{},
		},
	}
	set := branch.Set([]interface{}{"array", 2}, "hello")
	assert.Equal(t, TreeBranch{
		TreeItem{
			Key: "array",
			Value: []interface{}{
				"hello",
			},
		},
	}, set)
}

func TestSetArrayNonLeaf(t *testing.T) {
	branch := TreeBranch{
		TreeItem{
			Key: "array",
			Value: []interface{}{
				1,
			},
		},
	}
	set := branch.Set([]interface{}{"array", 0, "hello"}, "hello")
	assert.Equal(t, TreeBranch{
		TreeItem{
			Key: "array",
			Value: []interface{}{
				TreeBranch{
					TreeItem{
						Key:   "hello",
						Value: "hello",
					},
				},
			},
		},
	}, set)
}
