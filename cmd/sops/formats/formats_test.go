package formats

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFromString(t *testing.T) {
	assert.Equal(t, Binary, FormatFromString("foobar"))
	assert.Equal(t, Dotenv, FormatFromString("dotenv"))
	assert.Equal(t, Ini, FormatFromString("ini"))
	assert.Equal(t, Yaml, FormatFromString("yaml"))
	assert.Equal(t, Json, FormatFromString("json"))
}

func TestFormatForPath(t *testing.T) {
	assert.Equal(t, Binary, FormatForPath("/path/to/foobar"))
	assert.Equal(t, Dotenv, FormatForPath("/path/to/foobar.env"))
	assert.Equal(t, Ini, FormatForPath("/path/to/foobar.ini"))
	assert.Equal(t, Json, FormatForPath("/path/to/foobar.json"))
	assert.Equal(t, Yaml, FormatForPath("/path/to/foobar.yml"))
	assert.Equal(t, Yaml, FormatForPath("/path/to/foobar.yaml"))
}

func TestFormatForPathOrString(t *testing.T) {
	assert.Equal(t, Binary, FormatForPathOrString("/path/to/foobar", ""))
	assert.Equal(t, Dotenv, FormatForPathOrString("/path/to/foobar", "dotenv"))
	assert.Equal(t, Dotenv, FormatForPathOrString("/path/to/foobar.env", ""))
	assert.Equal(t, Ini, FormatForPathOrString("/path/to/foobar", "ini"))
	assert.Equal(t, Ini, FormatForPathOrString("/path/to/foobar.ini", ""))
	assert.Equal(t, Json, FormatForPathOrString("/path/to/foobar", "json"))
	assert.Equal(t, Json, FormatForPathOrString("/path/to/foobar.json", ""))
	assert.Equal(t, Yaml, FormatForPathOrString("/path/to/foobar", "yaml"))
	assert.Equal(t, Yaml, FormatForPathOrString("/path/to/foobar.yml", ""))

	assert.Equal(t, Ini, FormatForPathOrString("/path/to/foobar.yml", "ini"))
	assert.Equal(t, Binary, FormatForPathOrString("/path/to/foobar.yml", "binary"))
}
