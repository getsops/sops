package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRun_PathNormalization exercises the path-normalization codepaths
// end-to-end through Run. Run takes absolute paths only, so the CLI layer
// (filepath.Abs, cwd lookup) is mocked out by always passing absolute
// paths anchored at t.TempDir().
func TestRun_PathNormalization(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	confPath := filepath.Join(root, ".sops.yaml")
	require.NoError(t, os.WriteFile(confPath, []byte(`
creation_rules:
  - path_regex: '^secrets/.*\.yaml$'
    kms: 'arn:secrets'
  - kms: 'arn:fallback'
`), 0644))

	cases := []struct {
		name         string
		filePath     string
		wantIndex    int
		wantCatchAll bool
	}{
		{
			name:      "absolute path inside the tree → matches relative regex",
			filePath:  filepath.Join(root, "secrets", "db.yaml"),
			wantIndex: 0,
		},
		{
			name:      "absolute path in a nested subdir → matches relative regex",
			filePath:  filepath.Join(root, "secrets", "team", "team.yaml"),
			wantIndex: 0,
		},
		{
			name:         "absolute path outside the tree → falls back, hits catch-all",
			filePath:     filepath.Join(os.TempDir(), "sops-not-in-tree", "x.yaml"),
			wantIndex:    1,
			wantCatchAll: true,
		},
		{
			name:      "non-existent file in tree → matches as if it existed",
			filePath:  filepath.Join(root, "secrets", "future.yaml"),
			wantIndex: 0,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			opts := Opts{
				ConfigPath: confPath,
				FilePath:   tc.filePath,
			}
			output, exitCode, err := Run(opts)
			require.NoError(t, err)
			assert.Equal(t, 0, exitCode)
			require.Len(t, output.CreationRules, 1)
			assert.Equal(t, tc.wantIndex, output.CreationRules[0].RuleIndex)
		})
	}
}

func TestRun_ConfigOverride(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	primaryConf := filepath.Join(dir, ".sops.yaml")
	overrideConf := filepath.Join(dir, "override.yaml")

	require.NoError(t, os.WriteFile(primaryConf, []byte(`creation_rules:
  - kms: 'arn:primary'
`), 0644))
	require.NoError(t, os.WriteFile(overrideConf, []byte(`creation_rules:
  - kms: 'arn:override'
`), 0644))

	opts := Opts{
		ConfigPath: overrideConf,
		FilePath:   filepath.Join(dir, "x.yaml"),
	}
	output, _, err := Run(opts)
	require.NoError(t, err)
	require.Len(t, output.CreationRules, 1)
	require.Len(t, output.CreationRules[0].KMS, 1)
	assert.Equal(t, "arn:override", output.CreationRules[0].KMS[0].Arn)
}

// Confirm the NoRulesMatched code is what we actually use.
func TestRun_NoRulesMatchedCodeIs62(t *testing.T) {
	assert.Equal(t, 62, codes.NoRulesMatched, "spec promises code 62")
}
