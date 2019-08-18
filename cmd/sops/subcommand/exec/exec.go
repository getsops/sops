package exec

import (
	"bytes"
	"log"
	"io/ioutil"
	"syscall"
	"path/filepath"
	"os"
	"os/exec"
	"strings"
)

func init() {
}

type ExecOpts struct {
	Command string
	Plaintext []byte
	Background bool
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

func ExecWithFile(opts ExecOpts) {
	dir, err := ioutil.TempDir("/tmp/", ".sops")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	pipe := GetPipe(dir)
	placeholdered := strings.Replace(opts.Command, "{}", pipe, -1)
	go WritePipe(pipe, []byte(opts.Plaintext))

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

func Exec(opts ExecOpts) {
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
