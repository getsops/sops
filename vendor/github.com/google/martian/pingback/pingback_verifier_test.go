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

package pingback

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

func TestVerifyRequests(t *testing.T) {
	v := NewVerifier(&url.URL{
		Scheme:   "https",
		Host:     "example.com",
		Path:     "/test",
		RawQuery: "testing=true",
	})

	// Initial error state is failure. No pingback has been seen.
	err := v.VerifyRequests()
	if err == nil {
		t.Fatal("v.VerifyRequests(): got nil, want error")
	}

	want := "request(https://example.com/test?testing=true): pingback never occurred"
	if got := err.Error(); got != want {
		t.Errorf("err.Error(): got %q, want %q", got, want)
	}

	// Send non-matching request, error persists.
	req, err := http.NewRequest("GET", "http://www.google.com", nil)
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
	if err := v.VerifyRequests(); err == nil {
		t.Fatal("v.VerifyRequests(): got nil, want error")
	}

	// Send matching requests, clear error.
	req, err = http.NewRequest("GET", "https://example.com/test?testing=true", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, rmv, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer rmv()

	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := v.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequests(): got %v, want no error", err)
	}

	// Send non-matching request again, error is still nil after
	// pingback.
	req, err = http.NewRequest("GET", "http://www.google.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, rm, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer rm()

	if err := v.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := v.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequests(): got %v, want no error", err)
	}

	v.ResetRequestVerifications()
	if err := v.VerifyRequests(); err == nil {
		t.Error("VerifyRequests(): got nil, want error")
	}
}

func TestVerifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "pingback.Verifier": {
      "scope": ["request"],
      "scheme": "https",
      "host": "example.com",
      "path": "/testing",
      "query": "test=true"
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

	if err := reqv.VerifyRequests(); err == nil {
		t.Fatal("VerifyRequests(): got nil, want error")
	}

	req, err := http.NewRequest("GET", "https://example.com/testing?test=true", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, rm, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer rm()

	if err := reqv.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := reqv.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
}
