package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTreePath(t *testing.T) {
	path, err := parseTreePath(`["a"][0]['b'][12]`)
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{"a", 0, "b", 12}, path)
}

func TestParseTreePathInvalid(t *testing.T) {
	// These used to make sops panic instead of reporting an error.
	for _, arg := range []string{`[]`, `["]`, `[']`, `["a"][-1]`, `[-1]`, `["a"][]`} {
		path, err := parseTreePath(arg)
		assert.Error(t, err, arg)
		assert.Nil(t, path, arg)
	}
}
