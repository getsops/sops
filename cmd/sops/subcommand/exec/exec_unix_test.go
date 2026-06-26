//go:build !windows
// +build !windows

package exec

import (
	"os"
	"os/user"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserEnvReturnsCorrectVars(t *testing.T) {
	currentUser, err := user.Current()
	require.NoError(t, err)

	env := UserEnv(currentUser.Username)

	assert.Contains(t, env, "HOME="+currentUser.HomeDir)
	assert.Contains(t, env, "USER="+currentUser.Username)
	assert.Contains(t, env, "LOGNAME="+currentUser.Username)
	assert.Len(t, env, 3)
}

func TestUserEnvDoesNotModifyProcessEnv(t *testing.T) {
	currentUser, err := user.Current()
	require.NoError(t, err)

	originalHome := os.Getenv("HOME")

	UserEnv(currentUser.Username)

	assert.Equal(t, originalHome, os.Getenv("HOME"),
		"HOME should not be modified in the current process")
}

func TestExecWithFilePassesUserEnvToChild(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("skipping test that requires root privileges")
	}

	currentUser, err := user.Current()
	require.NoError(t, err)

	originalHome := os.Getenv("HOME")

	err = ExecWithFile(ExecOpts{
		Command:   "env | grep ^HOME=",
		Plaintext: []byte("hello"),
		User:      currentUser.Username,
		Fifo:      false,
	})
	require.NoError(t, err)

	assert.Equal(t, originalHome, os.Getenv("HOME"),
		"ExecWithFile should not modify HOME in the current process")
}

func TestExecWithEnvPassesUserEnvToChild(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("skipping test that requires root privileges")
	}

	currentUser, err := user.Current()
	require.NoError(t, err)

	originalHome := os.Getenv("HOME")

	err = ExecWithEnv(ExecOpts{
		Command:   "true",
		Plaintext: []byte{},
		User:      currentUser.Username,
	})
	require.NoError(t, err)

	assert.Equal(t, originalHome, os.Getenv("HOME"),
		"ExecWithEnv should not modify HOME in the current process")
}

func TestExecWithFilePristineIncludesUserEnv(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("skipping test that requires root privileges")
	}

	currentUser, err := user.Current()
	require.NoError(t, err)

	err = ExecWithFile(ExecOpts{
		Command:   "env | grep -q ^HOME=",
		Plaintext: []byte("hello"),
		User:      currentUser.Username,
		Pristine:  true,
		Fifo:      false,
	})
	require.NoError(t, err, "child should have HOME even with --pristine when --user is set")
}

func TestExecWithEnvPristineIncludesUserEnv(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("skipping test that requires root privileges")
	}

	currentUser, err := user.Current()
	require.NoError(t, err)

	err = ExecWithEnv(ExecOpts{
		Command:   "env | grep -q ^HOME=",
		Plaintext: []byte{},
		User:      currentUser.Username,
		Pristine:  true,
	})
	require.NoError(t, err, "child should have HOME even with --pristine when --user is set")
}
