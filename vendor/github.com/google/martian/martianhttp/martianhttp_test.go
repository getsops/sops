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

package martianhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"

	_ "github.com/google/martian/v3/header"
)

func TestNoModifiers(t *testing.T) {
	m := NewModifier()
	m.SetRequestModifier(nil)
	m.SetResponseModifier(nil)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := m.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}
}

func TestModifyRequest(t *testing.T) {
	m := NewModifier()
	tm := martiantest.NewModifier()

	m.SetRequestModifier(tm)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}

	m.SetRequestModifier(nil)

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
}

func TestModifyResponse(t *testing.T) {
	m := NewModifier()
	tm := martiantest.NewModifier()

	m.SetResponseModifier(tm)

	res := proxyutil.NewResponse(200, nil, nil)
	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}

	m.SetResponseModifier(nil)

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
}

func TestVerifyRequests(t *testing.T) {
	m := NewModifier()

	if err := m.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}

	verr := fmt.Errorf("request verification failure")

	m.SetRequestModifier(&verify.TestVerifier{
		RequestError: verr,
	})

	if err := m.VerifyRequests(); err != verr {
		t.Errorf("VerifyRequests(): got %v, want %v", err, verr)
	}

	m.ResetRequestVerifications()

	if err := m.VerifyRequests(); err != nil {
		t.Errorf("m.VerifyRequests(): got %v, want no error", err)
	}
}

func TestVerifyResponses(t *testing.T) {
	m := NewModifier()

	if err := m.VerifyResponses(); err != nil {
		t.Errorf("VerifyResponses(): got %v, want no error", err)
	}

	verr := fmt.Errorf("response verification failure")
	m.SetResponseModifier(&verify.TestVerifier{
		ResponseError: verr,
	})

	if err := m.VerifyResponses(); err != verr {
		t.Errorf("VerifyResponses(): got %v, want %v", err, verr)
	}

	m.ResetResponseVerifications()

	if err := m.VerifyResponses(); err != nil {
		t.Errorf("m.VerifyResponses(): got %v, want no error", err)
	}
}

func TestServeHTTPInvalidMethod(t *testing.T) {
	m := NewModifier()

	req, err := http.NewRequest("PATCH", "/configure", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, ...): got %v, want no error", "GET", err)
	}
	rw := httptest.NewRecorder()

	m.ServeHTTP(rw, req)
	if got, want := rw.Code, 405; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}
	if got, want := rw.Header().Get("Allow"), "GET, POST"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Allow", got, want)
	}
}

func TestServeHTTPInvalidJSON(t *testing.T) {
	m := NewModifier()

	req, err := http.NewRequest("POST", "/configure", bytes.NewReader([]byte("not-json")))
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, ...): got %v, want no error", "POST", "/configure", err)
	}
	rw := httptest.NewRecorder()

	m.ServeHTTP(rw, req)
	if got, want := rw.Code, 400; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}
}

func TestServeHTTP(t *testing.T) {
	m := NewModifier()

	body := []byte(`{
    "header.Modifier": {
      "scope": ["request", "response"],
			"name": "Martian-Test",
			"value": "true"
		}
	}`)

	req, err := http.NewRequest("POST", "/configure", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rw := httptest.NewRecorder()

	m.ServeHTTP(rw, req)
	if got, want := rw.Code, 200; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}

	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("m.ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("Martian-Test"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Martian-Test", got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("m.ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.Header.Get("Martian-Test"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-Test", got, want)
	}

	req, err = http.NewRequest("GET", "/configure", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	rw = httptest.NewRecorder()

	m.ServeHTTP(rw, req)
	if got, want := rw.Code, 200; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}

	got := new(bytes.Buffer)
	want := new(bytes.Buffer)

	if err := json.Compact(got, body); err != nil {
		t.Fatalf("json.Compact(body): got %v, want no error", err)
	}

	if err := json.Compact(want, rw.Body.Bytes()); err != nil {
		t.Fatalf("json.Compact(rw.Body): got %v, want no error", err)
	}

	if !bytes.Equal(got.Bytes(), want.Bytes()) {
		t.Errorf("rw.Body: got %q, want %q", got.Bytes(), want.Bytes())
	}
}
