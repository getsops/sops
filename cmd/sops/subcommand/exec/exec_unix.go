//go:build !windows
// +build !windows

package exec

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"
)

func ExecSyscall(command string, env []string) error {
	return syscall.Exec("/bin/sh", []string{"/bin/sh", "-c", command}, env)
}

func BuildCommand(command string) *exec.Cmd {
	return exec.Command("/bin/sh", "-c", command)
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

func GetPipe(dir, filename string) string {
	tmpfn := filepath.Join(dir, filename)
	err := syscall.Mkfifo(tmpfn, 0600)
	if err != nil {
		log.Fatal(err)
	}

	return tmpfn
}

func SwitchUser(username string) {
	user, err := user.Lookup(username)
	if err != nil {
		log.Fatal(err)
	}

	uid, _ := strconv.Atoi(user.Uid)

	err = syscall.Setgid(uid)
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

	err = syscall.Setregid(uid, uid)
	if err != nil {
		log.Fatal(err)
	}
}
