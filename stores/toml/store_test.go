package toml

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mozilla.org/sops/v3"
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
						Value: []interface{}{
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
								Value: []interface{}{
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
				Value: []interface{}{
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
						// 	Value: interface{}(nil),
						// },
					},
				},
			},
			// NOT IMPL on go-toml yet
			// sops.TreeItem{
			// 	Key: sops.Comment{
			// 		Value: " 41 comment",
			// 	},
			// 	Value: interface{}(nil),
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
