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

package martianurl

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

func TestVerifyRequests(t *testing.T) {
	u := &url.URL{
		Scheme:   "https",
		Host:     "*.example.com",
		Path:     "/test",
		RawQuery: "testing=true",
		Fragment: "test",
	}
	v := NewVerifier(u)

	tt := []struct {
		got, want string
	}{
		{
			got: "http://www.example.com/test?testing=true#test",
			want: `request(http://www.example.com/test?testing=true#test) url verify failure:
	Scheme: got "http", want "https"`,
		},
		{
			got: "http://www.martian.test/test?testing=true#test",
			want: `request(http://www.martian.test/test?testing=true#test) url verify failure:
	Scheme: got "http", want "https"
	Host: got "www.martian.test", want "*.example.com"`,
		},
		{
			got: "http://www.martian.test/prod?testing=true#test",
			want: `request(http://www.martian.test/prod?testing=true#test) url verify failure:
	Scheme: got "http", want "https"
	Host: got "www.martian.test", want "*.example.com"
	Path: got "/prod", want "/test"`,
		},
		{
			got: "http://www.martian.test/prod#test",
			want: `request(http://www.martian.test/prod#test) url verify failure:
	Scheme: got "http", want "https"
	Host: got "www.martian.test", want "*.example.com"
	Path: got "/prod", want "/test"
	Query: got "", want "testing=true"`,
		},
		{
			got: "http://www.martian.test/prod#fake",
			want: `request(http://www.martian.test/prod#fake) url verify failure:
	Scheme: got "http", want "https"
	Host: got "www.martian.test", want "*.example.com"
	Path: got "/prod", want "/test"
	Query: got "", want "testing=true"
	Fragment: got "fake", want "test"`,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", tc.got, nil)
		if err != nil {
			t.Fatalf("%d. http.NewRequest(..., %s, ...): got %v, want no error", i, tc.got, err)
		}

		_, remove, err := martian.TestContext(req, nil, nil)
		if err != nil {
			t.Fatalf("TestContext(): got %v, want no error", err)
		}
		defer remove()

		if err := v.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}
	}

	merr, ok := v.VerifyRequests().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyRequests(): got nil, want *verify.MultiError")
	}

	errs := merr.Errors()
	if got, want := len(errs), len(tt); got != want {
		t.Fatalf("len(merr.Errors()): got %d, want %d", got, want)
	}

	for i, tc := range tt {
		if got, want := errs[i].Error(), tc.want; got != want {
			t.Errorf("%d. err.Error(): mismatched error output\ngot: %s\nwant: %s", i, got, want)
		}
	}

	v.ResetRequestVerifications()

	// A valid request.
	req, err := http.NewRequest("GET", "https://www.example.com/test?testing=true#test", nil)
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
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestVerifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "url.Verifier": {
      "scope": ["request"],
      "scheme": "https",
      "host": "www.martian.proxy",
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

	req, err := http.NewRequest("GET", "https://www.martian.proxy/testing?test=false", nil)
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
