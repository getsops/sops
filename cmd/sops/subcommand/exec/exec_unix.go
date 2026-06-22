//go:build !windows
// +build !windows

package exec

import (
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

var forwardedSignals = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
}

func ExecSyscall(command string, env []string) error {
	return syscall.Exec("/bin/sh", []string{"/bin/sh", "-c", command}, env)
}

func BuildCommand(command string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", command)
}

func RunCommand(cmd *exec.Cmd) error {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	cmd.SysProcAttr.Setpgid = true

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, forwardedSignals...)

	if err := cmd.Start(); err != nil {
		signal.Stop(signals)
		return err
	}

	stopForwarding := make(chan struct{})
	forwardingDone := make(chan struct{})
	go func() {
		defer close(forwardingDone)
		for {
			select {
			case sig := <-signals:
				if err := syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal)); err != nil {
					log.WithError(err).Warn("Failed to forward signal to child process group")
				}
			case <-stopForwarding:
				return
			}
		}
	}()

	err := cmd.Wait()
	signal.Stop(signals)
	close(stopForwarding)
	<-forwardingDone
	return err
}

func WritePipe(pipe string, contents []byte) {
	handle, err := os.OpenFile(pipe, os.O_WRONLY, 0600)

	if err != nil {
		os.Remove(pipe)
		log.Fatal(err)
	}

	handle.Write(contents)
	handle.Close()
}

func GetPipe(dir, filename string) (string, error) {
	tmpfn := filepath.Join(dir, filename)
	err := syscall.Mkfifo(tmpfn, 0600)
	if err != nil {
		return "", err
	}

	return tmpfn, nil
}

func SwitchUser(username string) {
	user, err := user.Lookup(username)
	if err != nil {
		log.Fatal(err)
	}

	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)

	groupIds, err := user.GroupIds()
	var intGroupIds []int
	if err != nil {
		log.Fatal(err)
		intGroupIds = []int{gid}
	} else {
		intGroupIds = make([]int, len(groupIds))
		for i, gid := range groupIds {
			intGroupIds[i], _ = strconv.Atoi(gid)
		}
	}

	err = syscall.Setgroups(intGroupIds)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Setgid(gid)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Setuid(uid)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Setreuid(uid, uid)
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.Setregid(gid, gid)
	if err != nil {
		log.Fatal(err)
	}
}
