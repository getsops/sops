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

package querystring

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3/parse"
)

func TestNewQueryStringModifier(t *testing.T) {
	mod := NewModifier("testing", "true")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.URL.Query().Get("testing"), "true"; got != want {
		t.Errorf("req.URL.Query().Get(%q): got %q, want %q", "testing", got, want)
	}
}

func TestQueryStringModifierQueryExists(t *testing.T) {
	mod := NewModifier("testing", "true")

	req, err := http.NewRequest("GET", "/?testing=false", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.URL.Query().Get("testing"), "true"; got != want {
		t.Errorf("req.URL.Query().Get(%q): got %q, want %q", "testing", got, want)
	}
}

func TestQueryStringModifierQueryExistsMultipleKeys(t *testing.T) {
	mod := NewModifier("testing", "true")

	req, err := http.NewRequest("GET", "/?testing=false&testing=foo&foo=bar", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.URL.Query().Get("testing"), "true"; got != want {
		t.Errorf("req.URL.Query().Get(%q): got %q, want %q", "testing", got, want)
	}
	if got, want := req.URL.Query().Get("foo"), "bar"; got != want {
		t.Errorf("req.URL.Query().Get(%q): got %q, want %q", "testing", got, want)
	}
}

func TestModifierFromJSON(t *testing.T) {
	msg := []byte(`
	{
		"querystring.Modifier": {
      "scope": ["request"],
			"name": "param",
			"value": "true"
		}
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "http://martian.test", nil)
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

	if got, want := req.URL.Query().Get("param"), "true"; got != want {
		t.Errorf("req.URL.Query().Get(%q): got %q, want %q", "param", got, want)
	}
}
