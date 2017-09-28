// +build !appengine

/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Forked from https://camlistore.googlesource.com/camlistore/+/master/pkg/misc/gpgagent/

// Package gpgagent interacts with the local GPG Agent.
package gopgagent /* import "go.mozilla.org/gopgagent" */

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"

	"io"
	"net"
	"net/url"
	"os"
	"os/user"
	"path"
	"strings"
)

// Conn is a connection to the GPG agent.
type Conn struct {
	c  io.ReadWriteCloser
	br *bufio.Reader
}

var (
	ErrNoAgent = errors.New("GPG_AGENT_INFO not set in environment")
	ErrNoData  = errors.New("GPG_ERR_NO_DATA cache miss")
	ErrCancel  = errors.New("gpgagent: Cancel")
)

// NewConn connects to the GPG Agent as described in the
// GPG_AGENT_INFO environment variable.
func NewConn() (*Conn, error) {
	var addr *net.UnixAddr
	if gpgAgentInfo, ok := os.LookupEnv("GPG_AGENT_INFO"); ok {
		sp := strings.SplitN(gpgAgentInfo, ":", 3)
		if len(sp) == 0 || len(sp[0]) == 0 {
			return nil, ErrNoAgent
		}
		addr = &net.UnixAddr{Net: "unix", Name: sp[0]}
	} else {
		// If GPG_AGENT_INFO is not defined, we connect to the default socket,
		// S.gpg-agent, as the gpg-agent documentation recommends.
		// See the --use-standard-socket option in
		// <https://gnupg.org/documentation/manuals/gnupg-2.0/Agent-Options.html>
		currentUser, err := user.Current()
		if err != nil {
			return nil, ErrNoAgent
		}
		sockFile := path.Join(currentUser.HomeDir, ".gnupg", "S.gpg-agent")
		addr = &net.UnixAddr{Net: "unix", Name: sockFile}
	}
	uc, err := net.DialUnix("unix", nil, addr)
	if err != nil {
		return nil, err
	}
	br := bufio.NewReader(uc)
	lineb, err := br.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	line := string(lineb)
	if !strings.HasPrefix(line, "OK") {
		return nil, fmt.Errorf("gpgagent: didn't get OK; got %q", line)
	}
	return &Conn{uc, br}, nil
}

func (c *Conn) Close() error {
	c.br = nil
	return c.c.Close()
}

// PassphraseRequest is a request to get a passphrase from the GPG
// Agent.
type PassphraseRequest struct {
	CacheKey, Error, Prompt, Desc string

	// If the option --no-ask is used and the passphrase is not in
	// the cache the user will not be asked to enter a passphrase
	// but the error code GPG_ERR_NO_DATA is returned.  (ErrNoData)
	NoAsk bool
}

func (c *Conn) RemoveFromCache(cacheKey string) error {
	_, err := fmt.Fprintf(c.c, "CLEAR_PASSPHRASE %s\n", url.QueryEscape(cacheKey))
	if err != nil {
		return err
	}
	lineb, err := c.br.ReadSlice('\n')
	if err != nil {
		return err
	}
	line := string(lineb)
	if !strings.HasPrefix(line, "OK") {
		return fmt.Errorf("gpgagent: CLEAR_PASSPHRASE returned %q", line)
	}
	return nil
}

func (c *Conn) GetPassphrase(pr *PassphraseRequest) (passphrase string, outerr error) {
	defer func() {
		if e, ok := recover().(string); ok {
			passphrase = ""
			outerr = errors.New(e)
		}
	}()
	set := func(cmd string, val string) {
		if val == "" {
			return
		}
		_, err := fmt.Fprintf(c.c, "%s %s\n", cmd, val)
		if err != nil {
			panic("gpgagent: failed to send " + cmd)
		}
		line, _, err := c.br.ReadLine()
		if err != nil {
			panic("gpgagent: failed to read " + cmd)
		}
		if !strings.HasPrefix(string(line), "OK") {
			panic("gpgagent: response to " + cmd + " was " + string(line))
		}
	}
	if d := os.Getenv("DISPLAY"); d != "" {
		set("OPTION", "display="+d)
	}
	tty, err := os.Readlink("/proc/self/fd/0")
	if err == nil {
		set("OPTION", "ttyname="+tty)
	}
	set("OPTION", "ttytype="+os.Getenv("TERM"))
	opts := ""
	if pr.NoAsk {
		opts += "--no-ask "
	}

	encOrX := func(s string) string {
		if s == "" {
			return "X"
		}
		return url.QueryEscape(s)
	}

	_, err = fmt.Fprintf(c.c, "GET_PASSPHRASE %s%s %s %s %s\n",
		opts,
		url.QueryEscape(pr.CacheKey),
		encOrX(pr.Error),
		encOrX(pr.Prompt),
		encOrX(pr.Desc))
	if err != nil {
		return "", err
	}
	lineb, err := c.br.ReadSlice('\n')
	if err != nil {
		return "", err
	}
	line := string(lineb)
	if strings.HasPrefix(line, "OK ") {
		decb, err := hex.DecodeString(line[3 : len(line)-1])
		if err != nil {
			return "", err
		}
		return string(decb), nil
	}
	fields := strings.Split(line, " ")
	if len(fields) >= 2 && fields[0] == "ERR" {
		switch fields[1] {
		case "67108922":
			return "", ErrNoData
		case "83886179":
			return "", ErrCancel
		}
	}
	return "", errors.New(line)
}
