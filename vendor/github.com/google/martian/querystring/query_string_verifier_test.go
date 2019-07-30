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

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

func TestVerifyRequestPasses(t *testing.T) {
	v, err := NewVerifier("foo", "bar")
	if err != nil {
		t.Fatalf("NewVerifier(%q, %q): got %v, want no error", "foo", "bar", err)
	}

	req, err := http.NewRequest("GET", "http://www.google.com?foo=baz&foo=bar", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := v.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestVerifyEmptyValue(t *testing.T) {
	v, err := NewVerifier("foo", "")
	if err != nil {
		t.Fatalf("NewVerifier(%q, %q): got %v, want no error", "foo", "", err)
	}
	req, err := http.NewRequest("GET", "http://www.google.com?foo=bar", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()
	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := v.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestFailureWithMissingKey(t *testing.T) {
	v, err := NewVerifier("foo", "bar")
	if err != nil {
		t.Fatalf("NewVerifier(%q, %q): got %v, want no error", "foo", "bar", err)
	}
	req, err := http.NewRequest("GET", "http://www.google.com?fizz=bar", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()
	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	merr, ok := v.VerifyRequests().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyRequests(): got nil, want *verify.MultiError")
	}

	errs := merr.Errors()
	if len(errs) != 1 {
		t.Fatalf("len(merr.Errors()): got %d, want 1", len(errs))
	}

	expectErr := "request(http://www.google.com?fizz=bar) param verification error: key foo not found"
	for i := range errs {
		if got, want := errs[i].Error(), expectErr; got != want {
			t.Errorf("%d. err.Error(): mismatched error output\ngot: %s\nwant: %s", i, got, want)
		}
	}
}

func TestFailureWithMultiFail(t *testing.T) {
	v, err := NewVerifier("foo", "bar")
	if err != nil {
		t.Fatalf("NewVerifier(%q, %q): got %v, want no error", "foo", "bar", err)
	}
	req, err := http.NewRequest("GET", "http://www.google.com?foo=baz", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()
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

	expectErr := "request(http://www.google.com?foo=baz) param verification error: got baz for key foo, want bar"
	for i := range errs {
		if got, want := errs[i].Error(), expectErr; got != want {
			t.Errorf("%d. err.Error(): mismatched error output\ngot: %s\nwant: %s", i,
				got, want)
		}
	}
	v.ResetRequestVerifications()
	if err := v.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestBadInputToConstructor(t *testing.T) {
	if _, err := NewVerifier("", "bar"); err == nil {
		t.Fatalf("NewVerifier(): no error returned for empty key")
	}
}

func TestVerifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "querystring.Verifier": {
      "scope": ["request"],
      "name": "param",
      "value": "true"
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
	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()
	if err := reqv.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := reqv.VerifyRequests(); err == nil {
		t.Error("VerifyRequests(): got nil, want not nil")
	}
}
