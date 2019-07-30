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
)

func TestSetForwardHeaders(t *testing.T) {
	xfp := "X-Forwarded-Proto"
	xff := "X-Forwarded-For"
	xfh := "X-Forwarded-Host"
	xfu := "X-Forwarded-Url"

	m := NewForwardedModifier()
	req, err := http.NewRequest("GET", "http://martian.local?key=value", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.RemoteAddr = "10.0.0.1:8112"

	if m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get(xfp), "http"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xfp, got, want)
	}
	if got, want := req.Header.Get(xff), "10.0.0.1"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xff, got, want)
	}
	if got, want := req.Header.Get(xfh), "martian.local"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xfh, got, want)
	}
	if got, want := req.Header.Get(xfu), "http://martian.local?key=value"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xfh, got, want)
	}

	// Test with existing X-Forwarded-For.
	req.RemoteAddr = "12.12.12.12"

	if m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get(xff), "10.0.0.1, 12.12.12.12"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xff, got, want)
	}

	// Test that proto, host, and URL headers are preserved if already present.
	req, err = http.NewRequest("GET", "http://example.com/path?k=v", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set(xfp, "https")
	req.Header.Set(xfh, "preserved.host.com")
	req.Header.Set(xfu, "https://preserved.host.com/foo?x=y")

	if m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get(xfp), "https"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xfh, got, want)
	}
	if got, want := req.Header.Get(xfh), "preserved.host.com"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xfh, got, want)
	}
	if got, want := req.Header.Get(xfu), "https://preserved.host.com/foo?x=y"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", xfh, got, want)
	}
}
