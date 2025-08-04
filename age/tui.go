// These functions have been copied from the age project
// https://github.com/FiloSottile/age/blob/v1.0.0/cmd/age/encrypted_keys.go
// https://github.com/FiloSottile/age/blob/3d91014ea095e8d70f7c6c4833f89b53a96e0832/cmd/age/tui.go
//
// Copyright 2021 The age Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in age's LICENSE file at
// https://github.com/FiloSottile/age/blob/v1.0.0/LICENSE
//
// SPDX-License-Identifier: BSD-3-Clause

package age

import (
	"errors"
	"filippo.io/age/plugin"
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"golang.org/x/term"
)

const (
	SopsAgePasswordEnv = "SOPS_AGE_PASSWORD"
)

// readPassphrase reads a passphrase from the terminal. It does not read from a
// non-terminal stdin, so it does not check stdinInUse.
func readPassphrase(prompt string) ([]byte, error) {
	if testing.Testing() {
		password := os.Getenv(SopsAgePasswordEnv)
		if password != "" {
			return []byte(password), nil
		}
	}

	var (
		err        error
		passphrase []byte
	)

	err = withTerminal(func(in, out *os.File) error {
		_, err := fmt.Fprintf(out, "%s ", prompt)
		if err != nil {
			return fmt.Errorf("could not write prompt: %v", err)
		}

		// Use CRLF to work around an apparent bug in WSL2's handling of CONOUT$.
		// Only when running a Windows binary from WSL2, the cursor would not go
		// back to the start of the line with a simple LF. Honestly, it's impressive
		// CONIN$ and CONOUT$ even work at all inside WSL2.
		defer fmt.Fprintf(out, "\r\n")

		if passphrase, err = term.ReadPassword(int(in.Fd())); err != nil {
			return fmt.Errorf("could not read passphrase: %v", err)
		}

		return nil
	})

	return passphrase, err
}

func printf(format string, v ...interface{}) {
	log.Printf("age: "+format, v...)
}

func warningf(format string, v ...interface{}) {
	log.Printf("age: warning: "+format, v...)
}

// If testOnlyPanicInsteadOfExit is true, exit will set testOnlyDidExit and
// panic instead of calling os.Exit. This way, the wrapper in TestMain can
// recover the panic and return the exit code only if it was originated in exit.
var testOnlyPanicInsteadOfExit bool
var testOnlyDidExit bool

func exit(code int) {
	if testOnlyPanicInsteadOfExit {
		testOnlyDidExit = true
		panic(code)
	}
	os.Exit(code)
}

// clearLine clears the current line on the terminal, or opens a new line if
// terminal escape codes don't work.
func clearLine(out io.Writer) {
	const (
		CUI = "\033["   // Control Sequence Introducer
		CPL = CUI + "F" // Cursor Previous Line
		EL  = CUI + "K" // Erase in Line
	)

	// First, open a new line, which is guaranteed to work everywhere. Then, try
	// to erase the line above with escape codes.
	//
	// (We use CRLF instead of LF to work around an apparent bug in WSL2's
	// handling of CONOUT$. Only when running a Windows binary from WSL2, the
	// cursor would not go back to the start of the line with a simple LF.
	// Honestly, it's impressive CONIN$ and CONOUT$ work at all inside WSL2.)
	fmt.Fprintf(out, "\r\n"+CPL+EL)
}

// withTerminal runs f with the terminal input and output files, if available.
// withTerminal does not open a non-terminal stdin, so the caller does not need
// to check stdinInUse.
func withTerminal(f func(in, out *os.File) error) error {
	if runtime.GOOS == "windows" {
		in, err := os.OpenFile("CONIN$", os.O_RDWR, 0)
		if err != nil {
			return err
		}
		defer in.Close()
		out, err := os.OpenFile("CONOUT$", os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		defer out.Close()
		return f(in, out)
	} else if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		defer tty.Close()
		return f(tty, tty)
	} else if term.IsTerminal(int(os.Stdin.Fd())) {
		return f(os.Stdin, os.Stdin)
	} else {
		return fmt.Errorf("standard input is not a terminal, and /dev/tty is not available: %v", err)
	}
}

// readSecret reads a value from the terminal with no echo. The prompt is ephemeral.
func readSecret(prompt string) (s []byte, err error) {
	err = withTerminal(func(in, out *os.File) error {
		fmt.Fprintf(out, "%s ", prompt)
		defer clearLine(out)
		s, err = term.ReadPassword(int(in.Fd()))
		return err
	})
	return
}

// readCharacter reads a single character from the terminal with no echo. The
// prompt is ephemeral.
func readCharacter(prompt string) (c byte, err error) {
	err = withTerminal(func(in, out *os.File) error {
		fmt.Fprintf(out, "%s ", prompt)
		defer clearLine(out)

		oldState, err := term.MakeRaw(int(in.Fd()))
		if err != nil {
			return err
		}
		defer term.Restore(int(in.Fd()), oldState)

		b := make([]byte, 1)
		if _, err := in.Read(b); err != nil {
			return err
		}

		c = b[0]
		return nil
	})
	return
}

var pluginTerminalUI = &plugin.ClientUI{
	DisplayMessage: func(name, message string) error {
		printf("%s plugin: %s", name, message)
		return nil
	},
	RequestValue: func(name, message string, _ bool) (s string, err error) {
		defer func() {
			if err != nil {
				warningf("could not read value for age-plugin-%s: %v", name, err)
			}
		}()
		secret, err := readSecret(message)
		if err != nil {
			return "", err
		}
		return string(secret), nil
	},
	Confirm: func(name, message, yes, no string) (choseYes bool, err error) {
		defer func() {
			if err != nil {
				warningf("could not read value for age-plugin-%s: %v", name, err)
			}
		}()
		if no == "" {
			message += fmt.Sprintf(" (press enter for %q)", yes)
			_, err := readSecret(message)
			if err != nil {
				return false, err
			}
			return true, nil
		}
		message += fmt.Sprintf(" (press [1] for %q or [2] for %q)", yes, no)
		for {
			selection, err := readCharacter(message)
			if err != nil {
				return false, err
			}
			switch selection {
			case '1':
				return true, nil
			case '2':
				return false, nil
			case '\x03': // CTRL-C
				return false, errors.New("user cancelled prompt")
			default:
				warningf("reading value for age-plugin-%s: invalid selection %q", name, selection)
			}
		}
	},
	WaitTimer: func(name string) {
		printf("waiting on %s plugin...", name)
	},
}
