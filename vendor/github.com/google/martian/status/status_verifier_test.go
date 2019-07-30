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

package status

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"
)

func TestVerifyResponses(t *testing.T) {
	v := NewVerifier(301)

	tt := []struct {
		got  int
		want string
	}{
		{200, "response(http://www.example.com) status code verify failure: got 200, want 301"},
		{302, "response(http://www.example.com) status code verify failure: got 302, want 301"},
		{400, "response(http://www.example.com) status code verify failure: got 400, want 301"},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", "http://www.example.com", nil)
		if err != nil {
			t.Fatalf("%d. http.NewRequest(): got %v, want no error", i, err)
		}
		_, remove, err := martian.TestContext(req, nil, nil)
		if err != nil {
			t.Fatalf("TestContext(): got %v, want no error", err)
		}
		defer remove()

		res := proxyutil.NewResponse(tc.got, nil, req)

		if err := v.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}
	}

	merr, ok := v.VerifyResponses().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyResponses(): got nil, want *verify.MultiError")
	}
	errs := merr.Errors()
	if got, want := len(errs), len(tt); got != want {
		t.Fatalf("len(merr.Errors(): got %d, want %d", got, want)
	}

	for i, tc := range tt {
		if got, want := errs[i].Error(), tc.want; got != want {
			t.Errorf("%d. merr.Errors(): got %q, want %q", i, got, want)
		}
	}

	v.ResetResponseVerifications()

	if err := v.VerifyResponses(); err != nil {
		t.Errorf("v.VerifyResponses(): got %v, want no error", err)
	}
}

func TestVerifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "status.Verifier": {
      "scope": ["response"],
      "statusCode": 400
    }
  }`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}
	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}
	resv, ok := resmod.(verify.ResponseVerifier)
	if !ok {
		t.Fatal("reqmod.(verify.RequestVerifier): got !ok, want ok")
	}

	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(200, nil, req)
	if err := resv.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if err := resv.VerifyResponses(); err == nil {
		t.Error("VerifyResponses(): got nil, want not nil")
	}
}
