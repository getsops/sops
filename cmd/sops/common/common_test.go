package common

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteFileAtomic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "secret.yaml")
	require.NoError(t, os.WriteFile(path, []byte("original"), 0o600))

	require.NoError(t, WriteFileAtomic(path, []byte("replacement")))

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "replacement", string(contents))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 1, "temporary file should be removed")
}

func TestWriteFileAtomicCreatesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.yaml")

	require.NoError(t, WriteFileAtomic(path, []byte("data")))

	contents, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "data", string(contents))
}

func TestWriteFileAtomicPreservesPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("file permissions are not supported on Windows")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "secret.yaml")
	require.NoError(t, os.WriteFile(path, []byte("original"), 0o640))

	require.NoError(t, WriteFileAtomic(path, []byte("replacement")))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0o640), info.Mode().Perm())
}

func TestWriteFileAtomicFollowsSymlink(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated privileges on Windows")
	}

	dir := t.TempDir()
	target := filepath.Join(dir, "target.yaml")
	link := filepath.Join(dir, "link.yaml")
	require.NoError(t, os.WriteFile(target, []byte("original"), 0o600))
	require.NoError(t, os.Symlink(target, link))

	require.NoError(t, WriteFileAtomic(link, []byte("replacement")))

	info, err := os.Lstat(link)
	require.NoError(t, err)
	assert.NotZero(t, info.Mode()&os.ModeSymlink, "symlink should be preserved")

	contents, err := os.ReadFile(target)
	require.NoError(t, err)
	assert.Equal(t, "replacement", string(contents))
}

func TestWriteFileAtomicKeepsOriginalOnFailure(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("directory permissions are not enforced for this user")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "secret.yaml")
	require.NoError(t, os.WriteFile(path, []byte("original"), 0o600))
	require.NoError(t, os.Chmod(dir, 0o500))
	t.Cleanup(func() {
		_ = os.Chmod(dir, 0o700)
	})

	err := WriteFileAtomic(path, []byte("replacement"))
	require.Error(t, err)

	contents, readErr := os.ReadFile(path)
	require.NoError(t, readErr)
	assert.Equal(t, "original", string(contents))
}
