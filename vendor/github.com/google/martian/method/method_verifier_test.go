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

package method

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

func TestVerifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "method.Verifier": {
      "scope": ["request"],
      "method": "POST"
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
	reqv, ok := reqmod.(verify.RequestVerifier)
	if !ok {
		t.Fatal("reqmod.(verify.RequestVerifier): got !ok, want ok")
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := reqv.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := reqv.VerifyRequests(); err == nil {
		t.Error("VerifyRequests(): got nil, want not nil")
	}
}

func TestVerifyRequestPasses(t *testing.T) {
	for _, m := range []string{
		"GET",
		"HEAD",
		"PUT",
		"POST",
		"DELETE",
		"TRACE",
		"OPTIONS",
		"CONNECT",
		"PATCH",
	} {
		v, err := NewVerifier(m)
		if err != nil {
			t.Fatalf("NewVerifier(%q): got %v, want no error", m, err)
		}
		req, err := http.NewRequest(m, "www.google.com", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}
		if err := v.ModifyRequest(req); err != nil {
			t.Fatalf("ModifyRequest(): got %v, want no error", err)
		}
		if err := v.VerifyRequests(); err != nil {
			t.Fatalf("VerifyRequests(): got %v, want no error", err)
		}

		v.ResetRequestVerifications()

		if err := v.VerifyRequests(); err != nil {
			t.Errorf("v.VerifyRequests(): got %v, want no error", err)
		}
	}
}

func TestVerifyPostRequestFailsWithMultiFail(t *testing.T) {
	v, _ := NewVerifier("POST")
	req, err := http.NewRequest("GET", "http://www.google.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest got %v, want no error", err)
	}
	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	merr, ok := v.VerifyRequests().(*martian.MultiError)
	if !ok {
		t.Fatalf("VerifyRequests(): got nil, want *verify.MultiError")
	}

	errs := merr.Errors()
	if len(errs) != 2 {
		t.Fatalf("len(merr.Errors()): got %d, want 2", len(errs))
	}

	expectErr := "request(http://www.google.com) method verification error: got POST, want GET"
	for i := range errs {
		if got, want := errs[i].Error(), expectErr; got != want {
			t.Errorf("%d. err.Error(): mismatched error output\ngot: %s\nwant: %s", i,
				got, want)
		}
	}

	v.ResetRequestVerifications()

	if err := v.VerifyRequests(); err != nil {
		t.Errorf("v.VerifyRequests(): got %v, want no error", err)
	}
}

func TestBadInputToConstructor(t *testing.T) {
	if _, err := NewVerifier(""); err == nil {
		t.Fatalf("NewVerifier(): no error returned for empty")
	}
}
