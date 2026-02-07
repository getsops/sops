package toml

import (
	"reflect"
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
)

func testPlain() []byte {
	return []byte(
		`# 0 comment
0 = 0
1 = 1 # 1 comment

[2]
  21 = [21.1, 21.2] # 21 comment
  22 = 22

[2.1]
  211 = 211

[[2.1.1]]
  2111 = 2111

[[2.1.1]]
  2112 = 2112

[[3]]
  31 = "thirty one"

[[3]]
  32 = "thirty two"

# 4 comment
[4]
  # 41 comment
  41 = 41
`)
}

func testPlainNoComment() []byte {
	return []byte(
		`0 = 0
1 = 1

[2]
  21 = [21.1, 21.2]
  22 = 22

  [2.1]
    211 = 211

    [[2.1.1]]
      2111 = 2111

    [[2.1.1]]
      2112 = 2112

[[3]]
  31 = "thirty one"

[[3]]
  32 = "thirty two"

[4]
  41 = 41
`)
}

func testTreeBranches() sops.TreeBranches {
	return sops.TreeBranches{
		sops.TreeBranch{
			// NOT IMPL on go-toml yet
			// sops.TreeItem{
			// 	Key: sops.Comment{
			// 		Value: " 0 comment",
			// 	},
			// 	Value: interface{}(nil),
			// },
			sops.TreeItem{
				Key:   "0",
				Value: int64(0),
			},
			sops.TreeItem{
				Key:   "1",
				Value: int64(1),
			},
			// NOT IMPL on go-toml yet
			// sops.TreeItem{
			// 	Key: sops.Comment{
			// 		Value: " 1 comment",
			// 	},
			// 	Value: interface{}(nil),
			// },
			sops.TreeItem{
				Key: "2",
				Value: sops.TreeBranch{
					// NOT IMPL on go-toml yet
					// sops.TreeItem{
					// 	Key: sops.Comment{
					// 		Value: " 21 comment",
					// 	},
					// 	Value: interface{}(nil),
					// },
					sops.TreeItem{
						Key: "21",
						Value: []any{
							21.1,
							21.2,
						},
					},
					sops.TreeItem{
						Key:   "22",
						Value: int64(22),
					},
					sops.TreeItem{
						Key: "1",
						Value: sops.TreeBranch{
							sops.TreeItem{
								Key:   "211",
								Value: int64(211),
							},
							sops.TreeItem{
								Key: "1",
								Value: []any{
									sops.TreeBranch{
										sops.TreeItem{
											Key:   "2111",
											Value: int64(2111),
										},
									},
									sops.TreeBranch{
										sops.TreeItem{
											Key:   "2112",
											Value: int64(2112),
										},
									},
								},
							},
						},
					},
				},
			},
			sops.TreeItem{
				Key: "3",
				Value: []any{
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "31",
							Value: "thirty one",
						},
					},
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "32",
							Value: "thirty two",
						},
						// NOT IMPL on go-toml yet
						// sops.TreeItem{
						// 	Key: sops.Comment{
						// 		Value: " 4 comment",
						// 	},
						// 	Value: any(nil),
						// },
					},
				},
			},
			// NOT IMPL on go-toml yet
			// sops.TreeItem{
			// 	Key: sops.Comment{
			// 		Value: " 41 comment",
			// 	},
			// 	Value: any(nil),
			// },
			sops.TreeItem{
				Key: "4",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "41",
						Value: int64(41),
					},
				},
			},
		},
	}
}

func TestLoadPlainFile(t *testing.T) {
	t.Parallel()

	actualBranches, err := (&Store{}).LoadPlainFile(testPlain())
	if err != nil {
		t.Errorf("expected no error, got: %v", err)

		return
	}

	expectedBranches := testTreeBranches()

	if !reflect.DeepEqual(expectedBranches, actualBranches) {
		t.Errorf("expected\n%#v\ngot\n%#v", expectedBranches, actualBranches)

		return
	}
}

func TestEmitPlainFile(t *testing.T) {
	t.Parallel()

	branches := testTreeBranches()

	bytes, err := (&Store{}).EmitPlainFile(branches)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)

		return
	}

	if !reflect.DeepEqual(testPlainNoComment(), bytes) {
		t.Errorf("expected\n\n-%s-\n\ngot\n\n-%s-", testPlainNoComment(), bytes)

		return
	}
}

func TestEmitValueString(t *testing.T) {
	t.Parallel()

	bytes, err := (&Store{}).EmitValue("hello")
	assert.Nil(t, err)
	assert.Equal(t, []byte("\"hello\""), bytes)
}

func TestUnmarshalMetadataFromNonSOPSFile(t *testing.T) {
	t.Parallel()

	data := []byte(`hello = 2`)
	_, err := (&Store{}).LoadEncryptedFile(data)
	assert.Equal(t, sops.MetadataNotFound, err)
}

func TestLoadPlainFileRoundTrip(t *testing.T) {
	t.Parallel()

	// Load the plain file
	branches, err := (&Store{}).LoadPlainFile(testPlain())
	assert.Nil(t, err)

	// Emit it back
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	// Load again to verify round-trip works
	branches2, err := (&Store{}).LoadPlainFile(bytes)
	assert.Nil(t, err)

	// Should match the original loaded data
	assert.Equal(t, branches, branches2)
}

func TestEmitValueTreeBranch(t *testing.T) {
	t.Parallel()

	branch := sops.TreeBranch{
		sops.TreeItem{
			Key:   "key1",
			Value: "value1",
		},
		sops.TreeItem{
			Key:   "key2",
			Value: int64(42),
		},
	}

	bytes, err := (&Store{}).EmitValue(branch)
	assert.Nil(t, err)

	// Should be valid TOML
	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, int64(42), result["key2"])
}

func TestEmitValueNumber(t *testing.T) {
	t.Parallel()

	bytes, err := (&Store{}).EmitValue(42)
	assert.Nil(t, err)
	assert.Equal(t, []byte("42"), bytes)
}

func TestEmitValueBool(t *testing.T) {
	t.Parallel()

	bytes, err := (&Store{}).EmitValue(true)
	assert.Nil(t, err)
	assert.Equal(t, []byte("true"), bytes)
}

func TestEmpty(t *testing.T) {
	t.Parallel()

	// Empty TOML file - TOML treats empty input as an empty map (one branch with no items)
	branches, err := (&Store{}).LoadPlainFile([]byte(``))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))
	assert.Equal(t, 0, len(branches[0]))

	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)
	assert.Equal(t, ``, string(bytes))
}

func TestEmptyTable(t *testing.T) {
	t.Parallel()

	data := []byte(`[empty]
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))
	assert.Equal(t, 1, len(branches[0]))

	// Re-emit and verify
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	// Should contain the empty table
	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)
	assert.Contains(t, result, "empty")
}

func TestHasSopsTopLevelKey(t *testing.T) {
	t.Parallel()

	ok := (&Store{}).HasSopsTopLevelKey(sops.TreeBranch{
		sops.TreeItem{
			Key:   "sops",
			Value: "value",
		},
	})
	assert.True(t, ok)

	ok = (&Store{}).HasSopsTopLevelKey(sops.TreeBranch{
		sops.TreeItem{
			Key:   "sops_",
			Value: "value",
		},
	})
	assert.False(t, ok)

	ok = (&Store{}).HasSopsTopLevelKey(sops.TreeBranch{
		sops.TreeItem{
			Key:   "other",
			Value: "value",
		},
	})
	assert.False(t, ok)
}

func TestLoadEncryptedFile(t *testing.T) {
	t.Parallel()

	// Create a sample encrypted TOML with metadata
	data := []byte(`key1 = "value1"
key2 = 42

[sops]
  version = "3.7.0"
  mac = "ENC[AES256_GCM,data:abc123,iv:def456,tag:ghi789,type:str]"
  lastmodified = "2023-01-01T00:00:00Z"

  [[sops.kms]]
    arn = "arn:aws:kms:us-east-1:123456789012:key/12345678-1234-1234-1234-123456789012"
    created_at = "2023-01-01T00:00:00Z"
    enc = "encrypted-data-key"
`)

	tree, err := (&Store{}).LoadEncryptedFile(data)
	assert.Nil(t, err)

	// Verify metadata was loaded
	assert.NotNil(t, tree.Metadata)
	assert.Equal(t, "3.7.0", tree.Metadata.Version)
	assert.Equal(t, "ENC[AES256_GCM,data:abc123,iv:def456,tag:ghi789,type:str]", tree.Metadata.MessageAuthenticationCode)

	// Verify data was loaded
	assert.Equal(t, 1, len(tree.Branches))
	// The branch should contain key1, key2, and sops
	assert.GreaterOrEqual(t, len(tree.Branches[0]), 2)
}

func TestEmitEncryptedFile(t *testing.T) {
	t.Parallel()

	// Create a simple tree with metadata
	tree := sops.Tree{
		Branches: sops.TreeBranches{
			sops.TreeBranch{
				sops.TreeItem{
					Key:   "key1",
					Value: "value1",
				},
				sops.TreeItem{
					Key:   "key2",
					Value: int64(42),
				},
			},
		},
		Metadata: sops.Metadata{
			Version:                   "3.7.0",
			MessageAuthenticationCode: "test-mac",
		},
	}

	bytes, err := (&Store{}).EmitEncryptedFile(tree)
	assert.Nil(t, err)

	// Should be valid TOML
	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	// Should contain data
	assert.Equal(t, "value1", result["key1"])
	assert.Equal(t, int64(42), result["key2"])

	// Should contain sops metadata
	assert.Contains(t, result, "sops")
	sopsMap := result["sops"].(map[string]any)
	assert.Equal(t, "3.7.0", sopsMap["version"])
	assert.Equal(t, "test-mac", sopsMap["mac"])
}

func TestNestedStructures(t *testing.T) {
	t.Parallel()

	data := []byte(`[server]
  host = "localhost"
  port = 8080

  [server.database]
    name = "mydb"
    user = "admin"

[[users]]
  name = "Alice"
  role = "admin"

[[users]]
  name = "Bob"
  role = "user"
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))

	// Re-emit and verify structure is preserved
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	// Verify nested table
	assert.Contains(t, result, "server")
	server := result["server"].(map[string]any)
	assert.Equal(t, "localhost", server["host"])
	assert.Equal(t, int64(8080), server["port"])

	db := server["database"].(map[string]any)
	assert.Equal(t, "mydb", db["name"])
	assert.Equal(t, "admin", db["user"])

	// Verify array of tables
	assert.Contains(t, result, "users")
	users := result["users"].([]any)
	assert.Equal(t, 2, len(users))

	user1 := users[0].(map[string]any)
	assert.Equal(t, "Alice", user1["name"])
	assert.Equal(t, "admin", user1["role"])

	user2 := users[1].(map[string]any)
	assert.Equal(t, "Bob", user2["name"])
	assert.Equal(t, "user", user2["role"])
}

func TestErrorOnMultipleBranches(t *testing.T) {
	t.Parallel()

	branches := sops.TreeBranches{
		sops.TreeBranch{
			sops.TreeItem{Key: "key1", Value: "value1"},
		},
		sops.TreeBranch{
			sops.TreeItem{Key: "key2", Value: "value2"},
		},
	}

	// TOML can only contain one document
	_, err := (&Store{}).EmitPlainFile(branches)
	assert.NotNil(t, err)
	assert.Equal(t, errTOMLUniqueDocument, err)
}

func TestArraysOfArrays(t *testing.T) {
	t.Parallel()

	data := []byte(`arrays = [[["foo", "bar"], ["baz"]], [["qux"]]]
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))

	// Re-emit and verify structure is preserved
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	arrays := result["arrays"].([]any)
	assert.Equal(t, 2, len(arrays))

	// First nested array
	arr0 := arrays[0].([]any)
	assert.Equal(t, 2, len(arr0))
	arr00 := arr0[0].([]any)
	assert.Equal(t, []any{"foo", "bar"}, arr00)
	arr01 := arr0[1].([]any)
	assert.Equal(t, []any{"baz"}, arr01)

	// Second nested array
	arr1 := arrays[1].([]any)
	assert.Equal(t, 1, len(arr1))
	arr10 := arr1[0].([]any)
	assert.Equal(t, []any{"qux"}, arr10)
}

func TestSpecialCharacters(t *testing.T) {
	t.Parallel()

	// Test that TOML handles various special characters
	data := []byte(`simple = "hello world"
with_space = "hello world"
number = 42
`)

	// Load the TOML
	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))

	// Emit and verify round-trip works
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	// Re-load to verify round-trip
	branches2, err := (&Store{}).LoadPlainFile(bytes)
	assert.Nil(t, err)

	// Should be equal after round-trip
	assert.Equal(t, branches, branches2)
}

func TestEmitValueNonString(t *testing.T) {
	t.Parallel()

	// TreeBranch should work
	branch := sops.TreeBranch{
		sops.TreeItem{Key: "key", Value: "value"},
	}
	_, err := (&Store{}).EmitValue(branch)
	assert.Nil(t, err)

	// Other complex types should also work
	_, err = (&Store{}).EmitValue(42)
	assert.Nil(t, err)

	_, err = (&Store{}).EmitValue(true)
	assert.Nil(t, err)
}

func TestMixedArrayTypes(t *testing.T) {
	t.Parallel()

	// TOML requires arrays to have consistent types
	// Test that we can handle arrays with mixed simple types
	data := []byte(`mixed_numbers = [1, 2.5, 3]
strings = ["hello", "world"]
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)

	// Re-emit and verify
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	// Verify mixed_numbers array
	nums := result["mixed_numbers"].([]any)
	assert.Equal(t, 3, len(nums))

	// Verify strings array
	strs := result["strings"].([]any)
	assert.Equal(t, 2, len(strs))
	assert.Equal(t, "hello", strs[0])
	assert.Equal(t, "world", strs[1])
}

func TestDateTimeValues(t *testing.T) {
	t.Parallel()

	// TOML has native datetime support
	data := []byte(`datetime = 1979-05-27T07:32:00Z
date = 1979-05-27
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))

	// Re-emit and verify structure is preserved
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	// Should contain datetime fields
	assert.Contains(t, result, "datetime")
	assert.Contains(t, result, "date")
}

func TestInlineTable(t *testing.T) {
	t.Parallel()

	data := []byte(`name = { first = "Tom", last = "Preston-Werner" }
point = { x = 1, y = 2 }
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(branches))

	// Re-emit and verify structure is preserved
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	// Verify inline tables are loaded correctly
	name := result["name"].(map[string]any)
	assert.Equal(t, "Tom", name["first"])
	assert.Equal(t, "Preston-Werner", name["last"])

	point := result["point"].(map[string]any)
	assert.Equal(t, int64(1), point["x"])
	assert.Equal(t, int64(2), point["y"])
}

func TestBooleanValues(t *testing.T) {
	t.Parallel()

	data := []byte(`enabled = true
disabled = false
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)

	// Verify boolean values are preserved
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	assert.Equal(t, true, result["enabled"])
	assert.Equal(t, false, result["disabled"])
}

func TestFloatValues(t *testing.T) {
	t.Parallel()

	data := []byte(`pi = 3.14159
negative = -0.01
exponent = 5e+22
`)

	branches, err := (&Store{}).LoadPlainFile(data)
	assert.Nil(t, err)

	// Re-emit and verify
	bytes, err := (&Store{}).EmitPlainFile(branches)
	assert.Nil(t, err)

	var result map[string]any
	err = toml.Unmarshal(bytes, &result)
	assert.Nil(t, err)

	assert.Equal(t, 3.14159, result["pi"])
	assert.Equal(t, -0.01, result["negative"])
	assert.Equal(t, 5e+22, result["exponent"])
}
