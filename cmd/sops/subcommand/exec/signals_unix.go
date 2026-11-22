//go:build !windows
// +build !windows

package exec

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// forwardedSignals are the signals sops relays to a foreground child process.
// When sops runs as PID 1 (e.g. in a container) the kernel delivers stop
// signals only to sops; without forwarding, the child never sees them and
// cannot run its cleanup handlers (see issue #840).
var forwardedSignals = []os.Signal{
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGHUP,
	syscall.SIGQUIT,
}

// SignalExitError is returned by RunForeground when the child process was
// terminated by a signal. ExitCode follows the shell convention 128+signum,
// so callers can propagate a meaningful exit status.
type SignalExitError struct {
	Signal syscall.Signal
}

func (e *SignalExitError) Error() string {
	return fmt.Sprintf("child process terminated by signal %v", e.Signal)
}

// ExitCode returns 128+signum, like a POSIX shell does for a process killed
// by a signal.
func (e *SignalExitError) ExitCode() int {
	return 128 + int(e.Signal)
}

// RunForeground starts cmd in its own process group, forwards the signals in
// forwardedSignals that sops receives to the child's whole process group
// (so grandchildren such as test runners receive them too), waits for the
// child to exit and returns its error: an *exec.ExitError with the child's
// exit code, or a *SignalExitError when the child was killed by a signal.
func RunForeground(cmd *exec.Cmd) error {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}
	stop := forwardSignalsTo(cmd)
	defer stop()
	err := cmd.Wait()
	if exitErr, ok := err.(*exec.ExitError); ok {
		if status, ok := exitErr.Sys().(syscall.WaitStatus); ok && status.Signaled() {
			return &SignalExitError{Signal: status.Signal()}
		}
	}
	return err
}

// forwardSignalsTo relays every signal received on forwardedSignals to the
// process group of cmd. The returned function restores the default signal
// handling.
func forwardSignalsTo(cmd *exec.Cmd) func() {
	ch := make(chan os.Signal, len(forwardedSignals))
	signal.Notify(ch, forwardedSignals...)
	done := make(chan struct{})
	go func() {
		for {
			select {
			case sig := <-ch:
				// A negative pid targets the whole process group. ESRCH
				// means the child already exited and was reaped; there is
				// nothing left to signal.
				_ = syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
			case <-done:
				return
			}
		}
	}()
	return func() {
		signal.Stop(ch)
		close(done)
	}
}
