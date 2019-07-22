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

func TestModifyRequestWithMultipleHeaders(t *testing.T) {
	m := NewAppendModifier("X-Repeated", "modifier")

	req, err := http.NewRequest("GET", "www.example.com", nil)
	req.Header.Add("X-Repeated", "original")
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header["X-Repeated"][0], "original"; got != want {
		t.Errorf("req.Header[\"X-Repeated\"][0]: got %q, want %q", got, want)
	}
	if got, want := req.Header["X-Repeated"][1], "modifier"; got != want {
		t.Errorf("req.Header[\"X-Repeated\"][1]: got %q, want %q", got, want)
	}
}

func TestAppendModifierFromJSON(t *testing.T) {
	msg := []byte(`{
		"header.Append": {
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

	if n := len(req.Header["X-Martian"]); n != 2 {
		t.Errorf("res.Header[%q]: got len %d, want 2", "X-Martian", n)
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

	if n := len(res.Header["X-Martian"]); n != 2 {
		t.Errorf("res.Header[%q]: got len %d, want 2", "X-Martian", n)
	}
}
