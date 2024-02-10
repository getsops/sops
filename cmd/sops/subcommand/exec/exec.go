package exec

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/getsops/sops/v3/logging"

	"github.com/sirupsen/logrus"
)

const (
	FallbackFilename = "tmp-file"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("EXEC")
}

type ExecOpts struct {
	Command    string
	Plaintext  []byte
	Background bool
	Pristine   bool
	Fifo       bool
	User       string
	Filename   string
	Env        []string
}

func GetFile(dir, filename string) *os.File {
	// If no filename is provided, create a random one based on FallbackFilename
	if filename == "" {
		handle, err := os.CreateTemp(dir, FallbackFilename)
		if err != nil {
			log.Fatal(err)
		}
		return handle
	}
	// If a filename is provided, use that one
	handle, err := os.Create(filepath.Join(dir, filename))
	if err != nil {
		log.Fatal(err)
	}
	// read+write for owner only
	if err = handle.Chmod(0600); err != nil {
		log.Fatal(err)
	}
	return handle
}

func ExecWithFile(opts ExecOpts) error {
	if opts.User != "" {
		SwitchUser(opts.User)
	}

	if runtime.GOOS == "windows" && opts.Fifo {
		log.Warn("no fifos on windows, use --no-fifo next time")
		opts.Fifo = false
	}

	dir, err := os.MkdirTemp("", ".sops")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var filename string
	if opts.Fifo {
		// fifo handling needs to be async, even opening to write
		// will block if there is no reader present
		filename = opts.Filename
		if filename == "" {
			filename = FallbackFilename
		}
		filename = GetPipe(dir, filename)
		go WritePipe(filename, opts.Plaintext)
	} else {
		// GetFile handles opts.Filename == "" specially, that's why we have
		// to pass in opts.Filename without handling the fallback here
		handle := GetFile(dir, opts.Filename)
		handle.Write(opts.Plaintext)
		handle.Close()
		filename = handle.Name()
	}

	var env []string
	if !opts.Pristine {
		env = os.Environ()
	}
	env = append(env, opts.Env...)

	placeholdered := strings.Replace(opts.Command, "{}", filename, -1)
	cmd := BuildCommand(placeholdered)
	cmd.Env = env

	if opts.Background {
		return cmd.Start()
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ExecWithEnv(opts ExecOpts) error {
	if opts.User != "" {
		SwitchUser(opts.User)
	}

	var env []string

	if !opts.Pristine {
		env = os.Environ()
	}

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

	env = append(env, opts.Env...)

	cmd := BuildCommand(opts.Command)
	cmd.Env = env

	if opts.Background {
		return cmd.Start()
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
