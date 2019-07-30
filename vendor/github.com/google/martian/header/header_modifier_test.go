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

package header

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestNewHeaderModifier(t *testing.T) {
	mod := NewModifier("testing", "true")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("testing"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "testing", got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("testing"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "testing", got, want)
	}
}

func TestModifyRequestWithHostHeader(t *testing.T) {
	m := NewModifier("Host", "www.google.com")

	req, err := http.NewRequest("GET", "www.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Host, "www.google.com"; got != want {
		t.Errorf("req.Host: got %q, want %q", got, want)
	}
}

func TestModifierFromJSON(t *testing.T) {
	msg := []byte(`{
	  "header.Modifier": {
		  "scope": ["request", "response"],
			"name": "X-Martian",
			"value": "true"
    }
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "http://martian.test", nil)
	req.Header.Add("X-Martian", "false")
	if err != nil {
		t.Fatalf("http.NewRequest(): got %q, want no error", err)
	}

	reqmod := r.RequestModifier()

	if reqmod == nil {
		t.Fatalf("reqmod: got nil, want not nil")
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("X-Martian"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "X-Martian", got, want)
	}

	resmod := r.ResponseModifier()

	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header.Add("X-Martian", "false")
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("X-Martian"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "X-Martian", got, want)
	}
}
