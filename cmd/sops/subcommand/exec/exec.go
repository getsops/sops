package exec

import (
	"bytes"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"go.mozilla.org/sops/logging"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("EXEC")
}

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
		log.Warn("no fifos on windows, use --no-fifo next time")
		opts.Fifo = false
	}

	dir, err := ioutil.TempDir("", ".sops")
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
	cmd := BuildCommand(placeholdered)
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

	cmd := BuildCommand(opts.Command)
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
