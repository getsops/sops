package stores

import (
	"testing"

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
