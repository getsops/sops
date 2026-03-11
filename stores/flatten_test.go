package stores

import (
	"testing"

	"github.com/getsops/sops/v3"
	"github.com/stretchr/testify/assert"
)

func TestFlat(t *testing.T) {
	input := map[string]interface{}{
		"foo": "bar",
	}
	expected := map[string]interface{}{
		"foo": "bar",
	}
	flattened := Flatten(input)
	assert.Equal(t, expected, flattened)
	unflattened := Unflatten(flattened)
	assert.Equal(t, input, unflattened)
}

func TestMap(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": 0,
			"baz": 0,
		},
	}
	expected := map[string]interface{}{
		"foo" + mapSeparator + "bar": 0,
		"foo" + mapSeparator + "baz": 0,
	}
	flattened := Flatten(input)
	assert.Equal(t, expected, flattened)
	unflattened := Unflatten(flattened)
	assert.Equal(t, input, unflattened)
}

func TestFlattenMapMoreNesting(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": map[string]interface{}{
				"baz": 0,
			},
		},
	}
	expected := map[string]interface{}{
		"foo" + mapSeparator + "bar" + mapSeparator + "baz": 0,
	}
	flattened := Flatten(input)
	assert.Equal(t, expected, flattened)
	unflattened := Unflatten(flattened)
	assert.Equal(t, input, unflattened)
}

func TestFlattenList(t *testing.T) {
	input := map[string]interface{}{
		"foo": []interface{}{
			0,
		},
	}
	expected := map[string]interface{}{
		"foo" + listSeparator + "0": 0,
	}
	flattened := Flatten(input)
	assert.Equal(t, expected, flattened)
	unflattened := Unflatten(flattened)
	assert.Equal(t, input, unflattened)
}

func TestFlattenListWithMap(t *testing.T) {
	input := map[string]interface{}{
		"foo": []interface{}{
			map[string]interface{}{
				"bar": 0,
			},
		},
	}
	expected := map[string]interface{}{
		"foo" + listSeparator + "0" + mapSeparator + "bar": 0,
	}
	flattened := Flatten(input)
	assert.Equal(t, expected, flattened)
	unflattened := Unflatten(flattened)
	assert.Equal(t, input, unflattened)
}

func TestFlatten(t *testing.T) {
	input := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]interface{}{
			"foo": 2,
			"bar": map[string]interface{}{
				"foo": 2,
			},
		},
		"qux": []interface{}{
			"hello", 1, 2,
		},
	}
	expected := map[string]interface{}{
		"foo":                        "bar",
		"baz" + mapSeparator + "foo": 2,
		"baz" + mapSeparator + "bar" + mapSeparator + "foo": 2,
		"qux" + listSeparator + "0":                         "hello",
		"qux" + listSeparator + "1":                         1,
		"qux" + listSeparator + "2":                         2,
	}
	flattened := Flatten(input)
	assert.Equal(t, expected, flattened)
	unflattened := Unflatten(flattened)
	assert.Equal(t, input, unflattened)
}

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

func TestFlattenMetadata(t *testing.T) {
	tests := []struct {
		input Metadata
		want  map[string]interface{}
	}{
		{Metadata{MACOnlyEncrypted: false}, map[string]interface{}{"mac_only_encrypted": nil}},
		{Metadata{MACOnlyEncrypted: true}, map[string]interface{}{"mac_only_encrypted": true}},
		{Metadata{MessageAuthenticationCode: "line1\nline2"}, map[string]interface{}{"mac": "line1\nline2"}},
		{Metadata{MessageAuthenticationCode: "line1\n\n\nline2\n\nline3"}, map[string]interface{}{"mac": "line1\n\n\nline2\n\nline3"}},
	}

	for _, tt := range tests {
		got, err := FlattenMetadata(tt.input)
		assert.NoError(t, err)
		for k, v := range tt.want {
			assert.Equal(t, v, got[k])
		}
	}
}

func TestFlattenMetadataToUnflattenMetadata(t *testing.T) {
	tests := []struct {
		input Metadata
	}{
		{Metadata{MACOnlyEncrypted: true}},
		{Metadata{MACOnlyEncrypted: false}},
		{Metadata{ShamirThreshold: 3}},
		{Metadata{MessageAuthenticationCode: "line1\nline2"}},
		{Metadata{MessageAuthenticationCode: "line1\n\n\nline2\n\nline3"}},
	}

	for _, tt := range tests {
		flat, err := FlattenMetadata(tt.input)
		assert.NoError(t, err)
		md, err := UnflattenMetadata(flat)
		assert.NoError(t, err)
		assert.Equal(t, tt.input, md)
	}
}

func TestDecodeNewLines(t *testing.T) {
	tests := []struct {
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{map[string]interface{}{"mac": "line1\\nline2"}, map[string]interface{}{"mac": "line1\nline2"}},
		{map[string]interface{}{"mac": "line1\\n\\n\\nline2\\n\\nline3"}, map[string]interface{}{"mac": "line1\n\n\nline2\n\nline3"}},
	}

	for _, tt := range tests {
		DecodeNewLines(tt.input)
		for k, v := range tt.want {
			assert.Equal(t, v, tt.input[k])
		}
	}
}

func TestEncodeNewLines(t *testing.T) {
	tests := []struct {
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{map[string]interface{}{"mac": "line1\nline2"}, map[string]interface{}{"mac": "line1\\nline2"}},
		{map[string]interface{}{"mac": "line1\n\n\nline2\n\nline3"}, map[string]interface{}{"mac": "line1\\n\\n\\nline2\\n\\nline3"}},
	}

	for _, tt := range tests {
		EncodeNewLines(tt.input)
		for k, v := range tt.want {
			assert.Equal(t, v, tt.input[k])
		}
	}
}

func TestEncodeNonStrings(t *testing.T) {
	tests := []struct {
		input map[string]interface{}
		want  map[string]interface{}
	}{
		{map[string]interface{}{"mac_only_encrypted": false}, map[string]interface{}{"mac_only_encrypted": "false"}},
		{map[string]interface{}{"mac_only_encrypted": true}, map[string]interface{}{"mac_only_encrypted": "true"}},
		{map[string]interface{}{"shamir_threshold": 2}, map[string]interface{}{"shamir_threshold": "2"}},
		{map[string]interface{}{"shamir_threshold": 123}, map[string]interface{}{"shamir_threshold": "123"}},
	}

	for _, tt := range tests {
		EncodeNonStrings(tt.input)
		assert.Equal(t, tt.want, tt.input)
	}
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
		expectedOutputDupe = sops.TreeBranch{
			sops.TreeItem{
				Key: "key1",
				Value: []interface{}{
					"bar",
				},
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
	assert.Nil(t, err)
	assert.Equal(t, output, expectedOutputDupe)

	output, err = unflattenTreeBranch(inputCollision1)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening \"key1__map_foo\": Type mismatch: can only use string key for map", err.Error())
	assert.Nil(t, output)

	output, err = unflattenTreeBranch(inputCollision2)
	assert.NotNil(t, err)
	assert.Equal(t, "Error while unflattening \"key1__list_0\": Type mismatch: can only use integer key for list", err.Error())
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
