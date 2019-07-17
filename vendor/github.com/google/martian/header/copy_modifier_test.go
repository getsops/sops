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

func TestCopyModifier(t *testing.T) {
	m := NewCopyModifier("Original", "Copy")

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Original", "test")

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Copy"), "test"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Copy", got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header.Set("Original", "test")

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Copy"), "test"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Copy", got, want)
	}
}

func TestCopyModifierFromJSON(t *testing.T) {
	msg := []byte(`{
	  "header.Copy": {
			"from": "Original",
			"to": "Copy",
			"scope": ["request", "response"]
    }
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %q, want no error", err)
	}
	req.Header.Set("Original", "test")

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Copy"), "test"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Copy", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header.Set("Original", "test")

	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Copy"), "test"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Copy", got, want)
	}
}
