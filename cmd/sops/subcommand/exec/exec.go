package exec

import (
	"bytes"
	"log"
	"io/ioutil"
	"syscall"
	"path/filepath"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
)

func init() {
}

type ExecOpts struct {
	Command string
	Plaintext []byte
	Background bool
	Fifo bool
	User string
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

func GetPipe(dir string) string {
	tmpfn := filepath.Join(dir, "tmp-file")
	err := syscall.Mkfifo(tmpfn, 0600)
	if err != nil {
		log.Fatal(err)
	}

	return tmpfn
}

func GetFile(dir string) *os.File {
	handle, err := ioutil.TempFile(dir, "tmp-file")
	if err != nil {
		log.Fatal(err)
	}
	return handle
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

func ExecWithFile(opts ExecOpts) {
	if opts.User != "" {
		SwitchUser(opts.User)
	}

	dir, err := ioutil.TempDir("/tmp/", ".sops")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var filename string
	if opts.Fifo {
		// fifo handling needs to be async, even opening to write
		// will block if there is no reader present
		filename = GetPipe(dir)
		go WritePipe(filename, opts.Plaintext)
	} else {
		handle := GetFile(dir)
		handle.Write(opts.Plaintext)
		handle.Close()
		filename = handle.Name()
	}

	placeholdered := strings.Replace(opts.Command, "{}", filename, -1)
	cmd := exec.Command("/bin/sh", "-c", placeholdered)
	cmd.Env = os.Environ()

	if opts.Background {
		cmd.Start()
	} else {
		cmd.Stdin  = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func ExecWithEnv(opts ExecOpts) {
	if opts.User != "" {
		SwitchUser(opts.User)
	}

	env := os.Environ()
	lines := bytes.Split(opts.Plaintext, []byte("\n"))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			continue
		}
		env = append(env, string(line))
	}

	cmd := exec.Command("/bin/sh", "-c", opts.Command)
	cmd.Env = env

	if opts.Background {
		cmd.Start()
	} else {
		cmd.Stdin  = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}
