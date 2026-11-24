package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
)

func TestCheckNoClobberFlags(t *testing.T) {
	t.Parallel()

	err := CheckNoClobberFlags(false, false, false)
	assert.NoError(t, err)

	err = CheckNoClobberFlags(false, true, false)
	assert.NoError(t, err)

	err = CheckNoClobberFlags(true, true, false)
	assert.NoError(t, err)

	err = CheckNoClobberFlags(true, false, false)
	require.Error(t, err)
	exitErr, ok := err.(*cli.ExitError)
	require.True(t, ok)
	assert.Equal(t, codes.ErrorConflictingParameters, exitErr.ExitCode())
	assert.Contains(t, exitErr.Error(), "--no-clobber requires --output")

	err = CheckNoClobberFlags(true, false, true)
	require.Error(t, err)
	exitErr, ok = err.(*cli.ExitError)
	require.True(t, ok)
	assert.Equal(t, codes.ErrorConflictingParameters, exitErr.ExitCode())
	assert.Contains(t, exitErr.Error(), "cannot use --no-clobber with --in-place")

	err = CheckNoClobberFlags(true, true, true)
	require.Error(t, err)
	exitErr, ok = err.(*cli.ExitError)
	require.True(t, ok)
	assert.Equal(t, codes.ErrorConflictingParameters, exitErr.ExitCode())
	assert.Contains(t, exitErr.Error(), "cannot use --no-clobber with --in-place")
}

func TestOpenOutputFileCreatesNewFile(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")

	file, err := OpenOutputFile(path, true)
	require.NoError(t, err)
	defer file.Close()

	_, err = file.Write([]byte("fresh"))
	require.NoError(t, err)
	require.NoError(t, file.Close())

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, []byte("fresh"), got)
}

func TestOpenOutputFileNoClobberLeavesExistingBytes(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	original := []byte("do-not-clobber-me")
	require.NoError(t, os.WriteFile(path, original, 0o600))

	file, err := OpenOutputFile(path, true)
	require.Error(t, err)
	assert.Nil(t, file)
	assert.True(t, os.IsExist(err), "expected existence error, got %v", err)

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, original, got)
}

func TestOpenOutputFileNoClobberLeavesSymlinkTarget(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	link := filepath.Join(dir, "link.txt")
	original := []byte("symlink-target-bytes")
	require.NoError(t, os.WriteFile(target, original, 0o600))
	require.NoError(t, os.Symlink(target, link))

	file, err := OpenOutputFile(link, true)
	require.Error(t, err)
	assert.Nil(t, file)
	assert.True(t, os.IsExist(err), "expected existence error, got %v", err)

	got, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, original, got)

	info, err := os.Lstat(link)
	require.NoError(t, err)
	assert.NotEqual(t, 0, info.Mode()&os.ModeSymlink)
}

func TestOpenOutputFileWithoutNoClobberOverwrites(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	require.NoError(t, os.WriteFile(path, []byte("old"), 0o600))

	file, err := OpenOutputFile(path, false)
	require.NoError(t, err)
	defer file.Close()
	_, err = file.Write([]byte("new"))
	require.NoError(t, err)
	require.NoError(t, file.Close())

	got, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, []byte("new"), got)
}
