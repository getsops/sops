package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/cmd/sops/codes"
)

func treeWithItems(items ...sops.TreeItem) *sops.Tree {
	return &sops.Tree{Branches: sops.TreeBranches{sops.TreeBranch(items)}}
}

func TestEnvFromTree(t *testing.T) {
	env, err := envFromTree(treeWithItems(
		sops.TreeItem{Key: "A", Value: "1"},
		sops.TreeItem{Key: sops.Comment{Value: "ignored"}, Value: nil},
		sops.TreeItem{Key: "N", Value: 42},
		sops.TreeItem{Key: "B", Value: true},
		sops.TreeItem{Key: "M", Value: "x\ny=z"},
	))
	require.NoError(t, err)
	assert.Equal(t, []string{"A=1", "N=42", "B=true", "M=x\ny=z"}, env)
}

func TestEnvFromTreeErrors(t *testing.T) {
	tests := []struct {
		name     string
		item     sops.TreeItem
		expected string
	}{
		{
			name:     "list value",
			item:     sops.TreeItem{Key: "L", Value: []interface{}{"a", "b"}},
			expected: "cannot use complex value in environment; offending key L",
		},
		{
			name:     "map value",
			item:     sops.TreeItem{Key: "M", Value: sops.TreeBranch{{Key: "a", Value: "b"}}},
			expected: "cannot use complex value in environment; offending key M",
		},
		{
			name:     "non-string key",
			item:     sops.TreeItem{Key: 7, Value: "x"},
			expected: "cannot use non-string keys in environment, got int",
		},
		{
			name:     "'=' in key",
			item:     sops.TreeItem{Key: "A=B", Value: "x"},
			expected: "cannot use keys with '=' in environment: A=B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := envFromTree(treeWithItems(tt.item))
			assert.EqualError(t, err, tt.expected)
		})
	}
}

func TestExitErrorWithFileName(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedMsg  string
		expectedCode int
	}{
		{
			name:         "plain error",
			err:          errors.New("cannot use keys with '=' in environment: A=B"),
			expectedMsg:  "b.yaml: cannot use keys with '=' in environment: A=B",
			expectedCode: codes.ErrorGeneric,
		},
		{
			name:         "exit error keeps its exit code",
			err:          cli.NewExitError("Error reading file: some reason", codes.CouldNotReadInputFile),
			expectedMsg:  "b.yaml: Error reading file: some reason",
			expectedCode: codes.CouldNotReadInputFile,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := exitErrorWithFileName("b.yaml", tt.err)
			assert.EqualError(t, err, tt.expectedMsg)
			var exitCoder cli.ExitCoder
			require.ErrorAs(t, err, &exitCoder)
			assert.Equal(t, tt.expectedCode, exitCoder.ExitCode())
		})
	}
}
