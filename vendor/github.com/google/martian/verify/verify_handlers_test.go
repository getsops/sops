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

package verify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/martian/v3"
)

func TestHandlerServeHTTPUnsupportedMethod(t *testing.T) {
	h := NewHandler()

	for i, m := range []string{"POST", "PUT", "DELETE"} {
		req, err := http.NewRequest(m, "http://example.com", nil)
		if err != nil {
			t.Fatalf("%d. http.NewRequest(): got %v, want no error", i, err)
		}
		rw := httptest.NewRecorder()

		h.ServeHTTP(rw, req)
		if got, want := rw.Code, 405; got != want {
			t.Errorf("%d. rw.Code: got %d, want %d", i, got, want)
		}
		if got, want := rw.Header().Get("Allow"), "GET"; got != want {
			t.Errorf("%d. rw.Header().Get(%q): got %q, want %q", i, "Allow", got, want)
		}
	}
}

func TestHandlerServeHTTPNoVerifiers(t *testing.T) {
	h := NewHandler()
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if got, want := rw.Code, 200; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}
}

func TestHandlerServeHTTP(t *testing.T) {
	merr := martian.NewMultiError()
	merr.Add(fmt.Errorf("first response verification failure"))
	merr.Add(fmt.Errorf("second response verification failure"))

	v := &TestVerifier{
		RequestError:  fmt.Errorf("request verification failure"),
		ResponseError: merr,
	}

	h := NewHandler()
	h.SetRequestVerifier(v)
	h.SetResponseVerifier(v)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if got, want := rw.Code, 200; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}

	buf := new(bytes.Buffer)
	if err := json.Compact(buf, []byte(`{
    "errors": [
      { "message": "request verification failure" },
      { "message": "first response verification failure" },
      { "message": "second response verification failure" }
    ]
  }`)); err != nil {
		t.Fatalf("json.Compact(): got %v, want no error", err)
	}
	// json.(*Encoder).Encode writes a trailing newline, so we will too.
	// see: https://golang.org/src/encoding/json/stream.go
	buf.WriteByte('\n')

	if got, want := rw.Body.Bytes(), buf.Bytes(); !bytes.Equal(got, want) {
		t.Errorf("rw.Body: got %q, want %q", got, want)
	}
}

func TestResetHandlerServeHTTPUnsupportedMethod(t *testing.T) {
	h := NewResetHandler()

	for i, m := range []string{"GET", "PUT", "DELETE"} {
		req, err := http.NewRequest(m, "http://example.com", nil)
		if err != nil {
			t.Fatalf("%d. http.NewRequest(): got %v, want no error", i, err)
		}
		rw := httptest.NewRecorder()

		h.ServeHTTP(rw, req)
		if got, want := rw.Code, 405; got != want {
			t.Errorf("%d. rw.Code: got %d, want %d", i, got, want)
		}
		if got, want := rw.Header().Get("Allow"), "POST"; got != want {
			t.Errorf("%d. rw.Header().Get(%q): got %q, want %q", i, "Allow", got, want)
		}
	}
}

func TestResetHandlerServeHTTPNoVerifiers(t *testing.T) {
	h := NewResetHandler()
	req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if got, want := rw.Code, 204; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}
}

func TestResetHandlerServeHTTP(t *testing.T) {
	v := &TestVerifier{
		RequestError:  fmt.Errorf("request verification failure"),
		ResponseError: fmt.Errorf("response verification failure"),
	}

	h := NewResetHandler()
	h.SetRequestVerifier(v)
	h.SetResponseVerifier(v)

	req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)
	if got, want := rw.Code, 204; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}

	if err := v.VerifyRequests(); err != nil {
		t.Errorf("v.VerifyRequests(): got %v, want no error", err)
	}
	if err := v.VerifyResponses(); err != nil {
		t.Errorf("v.VerifyResponses(): got %v, want no error", err)
	}
}
