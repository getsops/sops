// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package port

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/martian/v3/parse"
)

func init() {
	parse.Register("port.Modifier", modifierFromJSON)
}

// Modifier alters the request URL and Host header to use the provided port.
// Only one of port, defaultForScheme, or remove may be specified.  Whichever is set last is the one that will take effect.
// If remove is true, remove the port from the host string ('example.com').
// If defaultForScheme is true, explicitly specify 80 for HTTP or 443 for HTTPS ('http://example.com:80'). Do nothing for a scheme that is not 'http' or 'https'.
// If port is specified, explicitly add it to the host string ('example.com:1234').
// If port is zero and the other fields are false, the request will not be modified.
type Modifier struct {
	port             int
	defaultForScheme bool
	remove           bool
}

type modifierJSON struct {
	Port             int                  `json:"port"`
	DefaultForScheme bool                 `json:"defaultForScheme"`
	Remove           bool                 `json:"remove"`
	Scope            []parse.ModifierType `json:"scope"`
}

// NewModifier returns a RequestModifier that can be configured to alter the request URL and Host header's port.
// One of DefaultPortForScheme, UsePort, or RemovePort should be called to configure this modifier.
func NewModifier() *Modifier {
	return &Modifier{}
}

// DefaultPortForScheme configures the modifier to explicitly specify 80 for HTTP or 443 for HTTPS ('http://example.com:80').
// The modifier will not modify requests with a scheme that is not 'http' or 'https'.
// This overrides any previous configuration for this modifier.
func (m *Modifier) DefaultPortForScheme() {
	m.defaultForScheme = true
	m.remove = false
}

// UsePort configures the modifier to add the specified port to the host string ('example.com:1234').
// This overrides any previous configuration for this modifier.
func (m *Modifier) UsePort(port int) {
	m.port = port
	m.remove = false
	m.defaultForScheme = false
}

// RemovePort configures the modifier to remove the port from the host string ('example.com').
// This overrides any previous configuration for this modifier.
func (m *Modifier) RemovePort() {
	m.remove = true
	m.defaultForScheme = false
}

// ModifyRequest alters the request URL and Host header to modify the port as specified.
// See docs for Modifier for details.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	if m.port == 0 && !m.defaultForScheme && !m.remove {
		return nil
	}

	host := req.URL.Host
	if strings.Contains(host, ":") {
		h, _, err := net.SplitHostPort(host)
		if err != nil {
			return err
		}
		host = h
	}

	if m.remove {
		req.URL.Host = host
		req.Header.Set("Host", host)
		return nil
	}

	if m.defaultForScheme {
		switch req.URL.Scheme {
		case "http":
			hp := net.JoinHostPort(host, "80")
			req.URL.Host = hp
			req.Header.Set("Host", hp)
			return nil
		case "https":
			hp := net.JoinHostPort(host, "443")
			req.URL.Host = hp
			req.Header.Set("Host", hp)
			return nil
		default:
			// Unknown scheme, do nothing.
			return nil
		}
	}

	// Not removing or using default for the scheme, so use the provided port number.
	hp := net.JoinHostPort(host, strconv.Itoa(m.port))
	req.URL.Host = hp
	req.Header.Set("Host", hp)

	return nil
}

func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	errMsg := fmt.Errorf("Must specify only one of port, defaultForScheme or remove")

	mod := NewModifier()
	// Check that exactly one field of port, defaultForScheme, and remove is set.
	switch {
	case msg.Port != 0:
		if msg.DefaultForScheme || msg.Remove {
			return nil, errMsg
		}
		mod.UsePort(msg.Port)
	case msg.DefaultForScheme:
		if msg.Port != 0 || msg.Remove {
			return nil, errMsg
		}
		mod.DefaultPortForScheme()
	case msg.Remove:
		if msg.Port != 0 || msg.DefaultForScheme {
			return nil, errMsg
		}
		mod.RemovePort()
	default:
		return nil, errMsg
	}

	return parse.NewResult(mod, msg.Scope)
}
