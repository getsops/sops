package exec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeEnv(t *testing.T) {
	tests := []struct {
		name      string
		base      []string
		overrides []string
		expected  []string
	}{
		{
			name:      "empty base",
			overrides: []string{"A=1", "B=2"},
			expected:  []string{"A=1", "B=2"},
		},
		{
			name:     "no overrides",
			base:     []string{"A=1", "B=2"},
			expected: []string{"A=1", "B=2"},
		},
		{
			name:      "disjoint keys keep their order",
			base:      []string{"A=1", "B=2"},
			overrides: []string{"C=3"},
			expected:  []string{"A=1", "B=2", "C=3"},
		},
		{
			name:      "override replaces earlier entry",
			base:      []string{"A=1", "B=1"},
			overrides: []string{"A=2"},
			expected:  []string{"B=1", "A=2"},
		},
		{
			name:      "only the key decides, values may contain '=' and newlines",
			base:      []string{"A=x=y", "B=1"},
			overrides: []string{"A=1=2\n3"},
			expected:  []string{"B=1", "A=1=2\n3"},
		},
		{
			name:      "keys are case-sensitive",
			base:      []string{"a=1"},
			overrides: []string{"A=2"},
			expected:  []string{"a=1", "A=2"},
		},
		{
			name:      "empty value overrides",
			base:      []string{"A=1"},
			overrides: []string{"A="},
			expected:  []string{"A="},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, MergeEnv(tt.base, tt.overrides))
		})
	}
}

func TestMergeEnvDoesNotModifyArguments(t *testing.T) {
	base := []string{"A=1", "B=1"}
	overrides := []string{"A=2"}

	MergeEnv(base, overrides)

	assert.Equal(t, []string{"A=1", "B=1"}, base)
	assert.Equal(t, []string{"A=2"}, overrides)
}
