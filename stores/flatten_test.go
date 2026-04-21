package stores

import (
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/stretchr/testify/assert"
)

func TestTokenizeFlat(t *testing.T) {
	input := "bar"
	expected := []token{mapToken{"bar"}}
	tokenized := tokenize(input)
	assert.Equal(t, expected, tokenized)
}

func TestTokenizeMap(t *testing.T) {
	input := "bar" + mapSeparator + "foo"
	expected := []token{mapToken{"bar"}, mapToken{"foo"}}
	tokenized := tokenize(input)
	assert.Equal(t, expected, tokenized)
}

func TestTokenizeList(t *testing.T) {
	input := "bar" + listSeparator + "10"
	expected := []token{mapToken{"bar"}, listToken{10}}
	tokenized := tokenize(input)
	assert.Equal(t, expected, tokenized)
}

func TestTokenizeNested(t *testing.T) {
	input := "bar" + listSeparator + "10" + mapSeparator + "baz"
	expected := []token{mapToken{"bar"}, listToken{10}, mapToken{"baz"}}
	tokenized := tokenize(input)
	assert.Equal(t, expected, tokenized)
}

func TestUnflattenTreeBranch(t *testing.T) {
	var (
		input = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key2__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key2__list_1__map_foo",
				Value: "bar",
			},
			sops.TreeItem{
				Key:   "key2__list_1__map_baz",
				Value: "bam",
			},
			sops.TreeItem{
				Key:   "key3__map_foo",
				Value: "bar",
			},
			sops.TreeItem{
				Key:   "key3__map_baz",
				Value: "bam",
			},
			sops.TreeItem{
				Key:   "key4__map_foo",
				Value: "bar",
			},
			sops.TreeItem{
				Key:   "key4__map_baz",
				Value: "bam",
			},
		}
		expectedOutput = sops.TreeBranch{
			sops.TreeItem{
				Key: "key1",
				Value: []interface{}{
					"foo",
				},
			},
			sops.TreeItem{
				Key: "key2",
				Value: []interface{}{
					"foo",
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "baz",
							Value: "bam",
						},
						sops.TreeItem{
							Key:   "foo",
							Value: "bar",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "key3",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "baz",
						Value: "bam",
					},
					sops.TreeItem{
						Key:   "foo",
						Value: "bar",
					},
				},
			},
			sops.TreeItem{
				Key: "key4",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "baz",
						Value: "bam",
					},
					sops.TreeItem{
						Key:   "foo",
						Value: "bar",
					},
				},
			},
		}

		inputDupe = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key1__list_0",
				Value: "bar",
			},
		}

		inputSkip1 = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key1__list_999999999999",
				Value: "bar",
			},
		}

		inputSkip2 = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1__list_1",
				Value: "foo",
			},
		}

		inputCollision1 = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key1__map_foo",
				Value: "bar",
			},
		}

		inputCollision2 = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key1__list_0",
				Value: "bar",
			},
		}
	)

	output, err := unflattenTreeBranch(input)
	assert.Nil(t, err)
	assert.Equal(t, output, expectedOutput)

	output, err = unflattenTreeBranch(inputDupe)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening \"key1__list_0\": Duplicate value", err.Error())
	assert.Nil(t, output)

	output, err = unflattenTreeBranch(inputSkip1)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening: Incomplete list", err.Error())
	assert.Nil(t, output)

	output, err = unflattenTreeBranch(inputSkip2)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening: Incomplete list", err.Error())
	assert.Nil(t, output)

	output, err = unflattenTreeBranch(inputCollision1)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening: Type mismatch", err.Error())
	assert.Nil(t, output)

	output, err = unflattenTreeBranch(inputCollision2)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening: Type mismatch", err.Error())
	assert.Nil(t, output)
}

func TestFlattenTreeBranch(t *testing.T) {
	var (
		input = sops.TreeBranch{
			sops.TreeItem{
				Key: "key1",
				Value: []interface{}{
					"foo",
				},
			},
			sops.TreeItem{
				Key: "key2",
				Value: []interface{}{
					"foo",
					sops.TreeBranch{
						sops.TreeItem{
							Key:   "foo",
							Value: "bar",
						},
						sops.TreeItem{
							Key:   "baz",
							Value: "bam",
						},
					},
				},
			},
			sops.TreeItem{
				Key: "key3",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "foo",
						Value: "bar",
					},
					sops.TreeItem{
						Key:   "baz",
						Value: "bam",
					},
				},
			},
			sops.TreeItem{
				Key: "key4",
				Value: sops.TreeBranch{
					sops.TreeItem{
						Key:   "foo",
						Value: "bar",
					},
					sops.TreeItem{
						Key:   "baz",
						Value: "bam",
					},
				},
			},
		}
		expectedOutput = sops.TreeBranch{
			sops.TreeItem{
				Key:   "prefixkey1__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "prefixkey2__list_0",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "prefixkey2__list_1__map_foo",
				Value: "bar",
			},
			sops.TreeItem{
				Key:   "prefixkey2__list_1__map_baz",
				Value: "bam",
			},
			sops.TreeItem{
				Key:   "prefixkey3__map_foo",
				Value: "bar",
			},
			sops.TreeItem{
				Key:   "prefixkey3__map_baz",
				Value: "bam",
			},
			sops.TreeItem{
				Key:   "prefixkey4__map_foo",
				Value: "bar",
			},
			sops.TreeItem{
				Key:   "prefixkey4__map_baz",
				Value: "bam",
			},
		}

		inputDupe = sops.TreeBranch{
			sops.TreeItem{
				Key:   "key1",
				Value: "foo",
			},
			sops.TreeItem{
				Key:   "key1",
				Value: "bar",
			},
		}
	)

	output, err := flattenTreeBranch(input, "prefix")
	assert.Nil(t, err)
	assert.Equal(t, output, expectedOutput)

	output, err = flattenTreeBranch(inputDupe, "prefix")
	assert.NotNil(t, err)
	assert.Equal(t, "Found key collision \"prefixkey1\" while flattening", err.Error())
	assert.Nil(t, output)
}
