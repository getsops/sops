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

package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTPSameOrigin(t *testing.T) {
	var handlerRun bool

	h := NewHandler(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			handlerRun = true
		}))

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if !handlerRun {
		t.Error("handlerRun: got false, want true")
	}
}

func TestServeHTTPPreflight(t *testing.T) {
	var handlerRun bool

	h := NewHandler(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			handlerRun = true
		}))
	h.AllowCredentials(true)

	req, err := http.NewRequest("OPTIONS", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Origin", "http://google.com")
	req.Header.Set("Access-Control-Request-Method", "PUT")
	req.Header.Set("Access-Control-Request-Headers", "Cors-Test")

	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if got, want := rw.Header().Get("Access-Control-Allow-Origin"), "*"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Origin", got, want)
	}
	if got, want := rw.Header().Get("Access-Control-Allow-Methods"), "PUT"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Methods", got, want)
	}
	if got, want := rw.Header().Get("Access-Control-Allow-Headers"), "Cors-Test"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Headers", got, want)
	}
	if got, want := rw.Header().Get("Access-Control-Allow-Credentials"), "true"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Credentials", got, want)
	}

	if handlerRun {
		t.Error("handlerRun: got true, want false")
	}
}

func TestServeHTTPSimple(t *testing.T) {
	var handlerRun bool

	h := NewHandler(http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			handlerRun = true
		}))
	h.SetOrigin("http://martian.local")

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Origin", "http://google.com")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "Cors-Test")

	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if got, want := rw.Header().Get("Access-Control-Allow-Origin"), "http://martian.local"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Origin", got, want)
	}
	if got, want := rw.Header().Get("Access-Control-Allow-Methods"), "GET"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Methods", got, want)
	}
	if got, want := rw.Header().Get("Access-Control-Allow-Headers"), "Cors-Test"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Access-Control-Allow-Headers", got, want)
	}

	if !handlerRun {
		t.Error("handlerRun: got false, want true")
	}
}
