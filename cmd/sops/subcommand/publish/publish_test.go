package publish

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3/config"
	"github.com/stretchr/testify/assert"
)

func TestRunWithoutMatchingDestination(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".sops.yaml")
	err := os.WriteFile(configPath, []byte(`
destination_rules:
  - vault_path: "foo/"
    path_regex: \.enc\.yml$
`), 0o600)
	assert.NoError(t, err)
	inputPath := filepath.Join(dir, "parameters.yml")
	err = os.WriteFile(inputPath, nil, 0o600)
	assert.NoError(t, err)

	opts := Opts{
		ConfigPath: configPath,
		InputPath:  inputPath,
		RootPath:   dir,
		Recursive:  true,
	}
	assert.NoError(t, Run(opts))

	opts.Recursive = false
	assert.ErrorIs(t, Run(opts), config.ErrNoMatchingDestination)
}

func TestRunRecursiveWithoutDestinationRules(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, ".sops.yaml")
	err := os.WriteFile(configPath, []byte("creation_rules:\n  - pgp: foo\n"), 0o600)
	assert.NoError(t, err)
	inputPath := filepath.Join(dir, "parameters.yml")
	err = os.WriteFile(inputPath, nil, 0o600)
	assert.NoError(t, err)

	err = Run(Opts{
		ConfigPath: configPath,
		InputPath:  inputPath,
		RootPath:   dir,
		Recursive:  true,
	})
	assert.Error(t, err)
	assert.NotErrorIs(t, err, config.ErrNoMatchingDestination)
}
