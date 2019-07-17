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

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

func TestIdModifier(t *testing.T) {
	mod := NewIDModifier()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("X-Martian-ID"), ctx.ID(); got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "X-Martian-ID", got, want)
	}
}

func TestFromJSON(t *testing.T) {
	msg := []byte(`{
			"header.Id": {}
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if _, ok := reqmod.(*idModifier); !ok {
		t.Fatal("reqmod.(*idModifier): got !ok, want ok")
	}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("X-Martian-ID"), ctx.ID(); got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "X-Martian-ID", got, want)
	}
}
