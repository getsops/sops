package stores

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)


func TestValToString(t *testing.T) {
	assert.Equal(t, "1", ValToString(1))
	assert.Equal(t, "1.0", ValToString(1.0))
	assert.Equal(t, "1.1", ValToString(1.10))
	assert.Equal(t, "1.23", ValToString(1.23))
	assert.Equal(t, "1.2345678901234567", ValToString(1.234567890123456789))
	assert.Equal(t, "200000.0", ValToString(2E5))
	assert.Equal(t, "-2E+10", ValToString(-2E10))
	assert.Equal(t, "2E-10", ValToString(2E-10))
	assert.Equal(t, "1.2345E+100", ValToString(1.2345E100))
	assert.Equal(t, "1.2345E-100", ValToString(1.2345E-100))
	assert.Equal(t, "true", ValToString(true))
	assert.Equal(t, "false", ValToString(false))
	ts, _ := time.Parse(time.RFC3339, "2025-01-02T03:04:05Z")
	assert.Equal(t, "2025-01-02T03:04:05Z", ValToString(ts))
	assert.Equal(t, "a string", ValToString("a string"))
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
