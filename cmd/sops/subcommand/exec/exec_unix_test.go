//go:build !windows
// +build !windows

package exec

import (
	"os"
	osExec "os/exec"
	"os/user"
	"path/filepath"
	"syscall"
	"testing"
	"time"

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

func TestSignalExitErrorExitCode(t *testing.T) {
	err := &SignalExitError{Signal: syscall.SIGTERM}
	assert.Equal(t, 143, err.ExitCode())
	assert.Contains(t, err.Error(), "terminated by signal")

	err = &SignalExitError{Signal: syscall.SIGINT}
	assert.Equal(t, 130, err.ExitCode())
}

func TestRunForegroundPropagatesChildExitCode(t *testing.T) {
	cmd := BuildCommand("exit 42")
	err := RunForeground(cmd)
	require.Error(t, err)
	exitErr, ok := err.(*osExec.ExitError)
	require.True(t, ok, "expected *exec.ExitError, got %T", err)
	assert.Equal(t, 42, exitErr.ExitCode())
}

func TestRunForegroundReturnsSignalExitErrorForSignaledChild(t *testing.T) {
	// The child kills itself with SIGTERM; RunForeground must surface that as
	// a *SignalExitError so the caller can exit with 128+SIGTERM = 143.
	cmd := BuildCommand("kill -TERM $$")
	err := RunForeground(cmd)
	require.Error(t, err)
	sigErr, ok := err.(*SignalExitError)
	require.True(t, ok, "expected *SignalExitError, got %T (%v)", err, err)
	assert.Equal(t, syscall.SIGTERM, sigErr.Signal)
	assert.Equal(t, 143, sigErr.ExitCode())
}

func TestRunForegroundForwardsSignalsToChildProcessGroup(t *testing.T) {
	dir := t.TempDir()
	ready := filepath.Join(dir, "ready")
	caught := filepath.Join(dir, "caught")

	// The child installs a SIGTERM handler, signals readiness, then loops.
	// When sops (this test process) receives SIGTERM, RunForeground must
	// forward it to the child's process group so the handler runs and the
	// child exits with code 7 instead of being killed outright.
	script := "trap 'echo caught > \"" + caught + "\"; exit 7' TERM; " +
		"echo ready > \"" + ready + "\"; while true; do sleep 0.05; done"
	cmd := BuildCommand(script)

	errCh := make(chan error, 1)
	go func() {
		errCh <- RunForeground(cmd)
	}()

	// Wait until the child has installed its trap.
	require.Eventually(t, func() bool {
		_, err := os.Stat(ready)
		return err == nil
	}, 5*time.Second, 20*time.Millisecond, "child did not become ready")

	// Deliver SIGTERM to ourselves; the forwarding goroutine must relay it.
	require.NoError(t, syscall.Kill(syscall.Getpid(), syscall.SIGTERM))

	select {
	case err := <-errCh:
		require.Error(t, err)
		exitErr, ok := err.(*osExec.ExitError)
		require.True(t, ok, "expected *exec.ExitError, got %T (%v)", err, err)
		assert.Equal(t, 7, exitErr.ExitCode(), "child should have run its TERM handler and exited 7")
	case <-time.After(10 * time.Second):
		t.Fatal("RunForeground did not return after SIGTERM was forwarded")
	}

	data, err := os.ReadFile(caught)
	require.NoError(t, err, "child TERM handler did not run")
	assert.Contains(t, string(data), "caught")
}

func TestRunForegroundStartError(t *testing.T) {
	cmd := BuildCommand("echo hi")
	cmd.Path = "/nonexistent/binary"
	err := RunForeground(cmd)
	require.Error(t, err)
	_, isExitErr := err.(*osExec.ExitError)
	_, isSigErr := err.(*SignalExitError)
	assert.False(t, isExitErr || isSigErr, "start errors must be returned unwrapped, got %T", err)
}
