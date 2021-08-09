package exec

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"

	"go.mozilla.org/sops/v3/logging"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	log = logging.NewLogger("EXEC")
}

type ExecOpts struct {
	Command    string
	Plaintext  []byte
	Background bool
	Fifo       bool
	User       string
	Filename   string
}

func GetFile(dir, filename string) *os.File {
	handle, err := ioutil.TempFile(dir, filename)
	if err != nil {
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

	dir, err := ioutil.TempDir("", ".sops")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	var filename string
	if opts.Fifo {
		// fifo handling needs to be async, even opening to write
		// will block if there is no reader present
		filename = GetPipe(dir, opts.Filename)
		go WritePipe(filename, opts.Plaintext)
	} else {
		handle := GetFile(dir, opts.Filename)
		handle.Write(opts.Plaintext)
		handle.Close()
		filename = handle.Name()
	}

	placeholdered := strings.Replace(opts.Command, "{}", filename, -1)
	cmd := BuildCommand(placeholdered)
	cmd.Env = os.Environ()

	if opts.Background {
		return cmd.Start()
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := cmd.Start(); err != nil {
		return err
	}
	redirectSignals(ctx, cmd)

	return cmd.Wait()
}

func ExecWithEnv(opts ExecOpts) error {
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
		return cmd.Start()
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := cmd.Start(); err != nil {
		return err
	}
	redirectSignals(ctx, cmd)

	return cmd.Wait()
}

func redirectSignals(ctx context.Context, cmd *exec.Cmd) {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc)

	go func(c chan os.Signal) {
		for {
			select {
			case sig := <-c:
				if err := cmd.Process.Signal(sig); err != nil {
					log.Warn("Unable to send a signal to the child process")
					continue
				}
			case <-ctx.Done():
				signal.Stop(sigc)
				return
			}
		}
	}(sigc)
}
