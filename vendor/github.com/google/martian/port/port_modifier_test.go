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
	"net"
	"net/http"
	"testing"

	"github.com/google/martian/v3/parse"
)

func TestPortModifierOnPort(t *testing.T) {
	mod := NewModifier()
	mod.UsePort(8080)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	_, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		t.Fatalf("net.SplitHostPort(%q): got %v, want no error", req.URL.Host, err)
	}

	if got, want := port, "8080"; got != want {
		t.Errorf("port: got %v, want %v", got, want)
	}
}

func TestPortModifierWithNoConfiguration(t *testing.T) {
	mod := NewModifier()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.URL.Host, "example.com"; got != want {
		t.Errorf("req.URL.Host: got %v, want %v", got, want)
	}

	req, err = http.NewRequest("GET", "http://example.com:80", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.URL.Host, "example.com:80"; got != want {
		t.Errorf("req.URL.Host: got %v, want %v", got, want)
	}
}

func TestPortModifierDefaultForScheme(t *testing.T) {
	mod := NewModifier()
	mod.DefaultPortForScheme()

	req, err := http.NewRequest("GET", "HtTp://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.URL.Host, "example.com:80"; got != want {
		t.Errorf("req.URL.Host: got %v, want %v", got, want)
	}
}

func TestPortModifierRemove(t *testing.T) {
	mod := NewModifier()
	mod.RemovePort()

	req, err := http.NewRequest("GET", "http://example.com:8080", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.URL.Host, "example.com"; got != want {
		t.Errorf("req.URL.Host: got %v, want %v", got, want)
	}
}

func TestPortModifierAllFields(t *testing.T) {
	mod := NewModifier()
	mod.UsePort(8081)
	mod.DefaultPortForScheme()
	mod.RemovePort()

	req, err := http.NewRequest("GET", "http://example.com:8080", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	// Last configuration was to remove.
	if got, want := req.URL.Host, "example.com"; got != want {
		t.Errorf("req.URL.Host: got %v, want %v", got, want)
	}
}

func TestModiferFromJSON(t *testing.T) {
	msg := []byte(`{
      "port.Modifier": {
        "scope": ["request"],
        "port": 8080
      }
    }`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	_, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		t.Fatalf("net.SplitHostPort(%q): got %v, want no error", req.URL.Host, err)
	}

	if got, want := port, "8080"; got != want {
		t.Errorf("port: got %v, want %v", got, want)
	}

}

func TestModiferFromJSONInvalidConfigurations(t *testing.T) {
	for _, msg := range [][]byte{
		[]byte(`{
				"port.Modifier": {
					"scope": ["request"],
					"port": 8080,
					"defaultForScheme": true,
					"remove": true
				}
			}`),
		[]byte(`{
				"port.Modifier": {
					"scope": ["request"],
					"port": 8080
					"remove": true
				}
			}`),
		[]byte(`{
				"port.Modifier": {
					"scope": ["request"],
					"port": 8080
					"defaultForScheme": true,
				}
			}`),
		[]byte(`{
				"port.Modifier": {
					"scope": ["request"],
					"defaultForScheme": true,
					"remove": true
				}
			}`),
		[]byte(`{
				"port.Modifier": {
					"scope": ["request"],
				}
			}`),
		[]byte(`{
				"port.Modifier": {
					"scope": ["response"],
					"remove": true
				}
			}`),
	} {
		_, err := parse.FromJSON(msg)
		if err == nil {
			t.Fatalf("parseFromJSON(msg): Got no error, but should have gotten one.")
		}
	}
}
