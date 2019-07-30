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

package stash

import (
	"net"
	"net/http"
	"testing"

	"github.com/google/martian/v3/fifo"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/port"
	"github.com/google/martian/v3/proxyutil"
)

func TestStashRequest(t *testing.T) {
	fg := fifo.NewGroup()
	fg.AddRequestModifier(NewModifier("stashed-url"))
	pmod := port.NewModifier()
	pmod.UsePort(8080)
	fg.AddRequestModifier(pmod)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := fg.ModifyRequest(req); err != nil {
		t.Fatalf("smod.ModifyRequest(): got %v, want no error", err)
	}

	_, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		t.Fatalf("net.SplitHostPort(%q): got %v, want no error", req.URL.Host, err)
	}

	if got, want := port, "8080"; got != want {
		t.Errorf("port: got %v, want %v", got, want)
	}

	if got, want := req.Header.Get("stashed-url"), "http://example.com"; got != want {
		t.Errorf("stashed-url header: got %v, want %v", got, want)
	}

}

func TestStashRequestResponse(t *testing.T) {
	headerName := "stashed-url"
	originalURL := "http://example.com"
	fg := fifo.NewGroup()
	fg.AddRequestModifier(NewModifier(headerName))
	fg.AddResponseModifier(NewModifier(headerName))
	pmod := port.NewModifier()
	pmod.UsePort(8080)
	fg.AddRequestModifier(pmod)

	req, err := http.NewRequest("GET", originalURL, nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := fg.ModifyRequest(req); err != nil {
		t.Fatalf("smod.ModifyRequest(): got %v, want no error", err)
	}

	_, port, err := net.SplitHostPort(req.URL.Host)
	if err != nil {
		t.Fatalf("net.SplitHostPort(%q): got %v, want no error", req.URL.Host, err)
	}

	if got, want := port, "8080"; got != want {
		t.Errorf("port: got %v, want %v", got, want)
	}

	if got, want := req.Header.Get(headerName), originalURL; got != want {
		t.Errorf("res.Header.Get(%q): got %v, want %v", headerName, got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := fg.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get(headerName), originalURL; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", headerName, got, want)
	}
}

func TestStashInvalidHeaderName(t *testing.T) {
	mod := NewModifier("invalid-chars-actually-work-;><@")

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("smod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("invalid-chars-actually-work-;><@"), "http://example.com"; got != want {
		t.Errorf("stashed-url header: got %v, want %v", got, want)
	}
}

func TestModiferFromJSON(t *testing.T) {
	headerName := "stashed-url"
	originalURL := "http://example.com"
	msg := []byte(`{
    "fifo.Group": {
      "scope": ["request", "response"],
      "modifiers": [
        {
          "stash.Modifier": {
            "scope": ["request", "response"],
            "headerName": "stashed-url"
          }
        },
        {
          "port.Modifier": {
            "scope": ["request"],
            "port": 8080
          }
        }
      ]
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

	req, err := http.NewRequest("GET", originalURL, nil)
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

	if got, want := req.Header.Get(headerName), originalURL; got != want {
		t.Errorf("req.Header.Get(%q) header: got %v, want %v", headerName, got, want)
	}

	resmod := r.ResponseModifier()
	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get(headerName), originalURL; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", headerName, got, want)
	}
}

func TestModiferFromJSONInvalidConfigurations(t *testing.T) {
	msg := []byte(`{
      "stash.Modifier": {
        "scope": ["response"],
        "headerName": "stash-header"
      }
    }`)

	_, err := parse.FromJSON(msg)
	if err == nil {
		t.Fatalf("parseFromJSON(msg): Got no error, but should have gotten one.")
	}
}
