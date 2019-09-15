package exec

import (
	"log"
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type ExecOpts struct {
	Command string
	Plaintext []byte
	Background bool
	Fifo bool
	User string
}

func GetFile(dir string) *os.File {
	handle, err := ioutil.TempFile(dir, "tmp-file")
	if err != nil {
		log.Fatal(err)
	}
	return handle
}

func ExecWithFile(opts ExecOpts) {
	if opts.User != "" {
		SwitchUser(opts.User)
	}

	if runtime.GOOS == "windows" && opts.Fifo {
		log.Print("no fifos on windows, use --no-fifo next time")
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
