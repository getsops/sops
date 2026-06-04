//go:build !windows
// +build !windows

package exec

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRunCommandForwardsSignalsToChildProcessGroup(t *testing.T) {
	dir := t.TempDir()
	readyFile := filepath.Join(dir, "ready")
	signalFile := filepath.Join(dir, "signal")

	cmd := BuildCommand(`
trap 'echo term > "$SIGNAL_FILE"; exit 42' TERM
echo ready > "$READY_FILE"
while true; do sleep 1; done
`)
	cmd.Env = append(os.Environ(), "READY_FILE="+readyFile, "SIGNAL_FILE="+signalFile)

	errCh := make(chan error, 1)
	go func() {
		errCh <- RunCommand(cmd)
	}()

	require.Eventually(t, func() bool {
		_, err := os.Stat(readyFile)
		return err == nil
	}, 3*time.Second, 25*time.Millisecond)

	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGTERM))

	select {
	case err := <-errCh:
		var exitErr *exec.ExitError
		require.True(t, errors.As(err, &exitErr), "expected command to return an exit error")
		require.Equal(t, 42, exitErr.ExitCode())
	case <-time.After(3 * time.Second):
		if cmd.Process != nil {
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		}
		t.Fatal("command did not receive the forwarded signal")
	}

	contents, err := os.ReadFile(signalFile)
	require.NoError(t, err)
	require.Equal(t, "term\n", string(contents))
}
