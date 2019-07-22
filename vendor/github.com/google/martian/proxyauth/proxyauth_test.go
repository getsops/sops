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

package proxyauth

import (
	"encoding/base64"
	"errors"
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/auth"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
)

func encode(v string) string {
	return base64.StdEncoding.EncodeToString([]byte(v))
}

func TestNoModifiers(t *testing.T) {
	m := NewModifier()
	m.SetRequestModifier(nil)
	m.SetResponseModifier(nil)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := m.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := m.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}
}

func TestProxyAuth(t *testing.T) {
	m := NewModifier()
	tm := martiantest.NewModifier()
	m.SetRequestModifier(tm)
	m.SetResponseModifier(tm)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Proxy-Authorization", "Basic "+encode("user:pass"))

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	actx := auth.FromContext(ctx)
	if got, want := actx.ID(), "user:pass"; got != want {
		t.Fatalf("actx.ID(): got %q, want %q", got, want)
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}

	res := proxyutil.NewResponse(200, nil, req)

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Proxy-Authenticate"), ""; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Proxy-Authenticate", got, want)
	}
}

func TestProxyAuthInvalidCredentials(t *testing.T) {
	m := NewModifier()
	authErr := errors.New("auth error")

	tm := martiantest.NewModifier()
	tm.RequestFunc(func(req *http.Request) {
		ctx := martian.NewContext(req)
		actx := auth.FromContext(ctx)

		actx.SetError(authErr)
	})
	tm.ResponseFunc(func(res *http.Response) {
		ctx := martian.NewContext(res.Request)
		actx := auth.FromContext(ctx)

		actx.SetError(authErr)
	})

	m.SetRequestModifier(tm)
	m.SetResponseModifier(tm)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Proxy-Authorization", "Basic "+encode("user:pass"))

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}

	actx := auth.FromContext(ctx)
	if actx.Error() != authErr {
		t.Fatalf("auth.Error(): got %v, want %v", actx.Error(), authErr)
	}
	actx.SetError(nil)

	res := proxyutil.NewResponse(200, nil, req)

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}

	if actx.Error() != authErr {
		t.Fatalf("auth.Error(): got %v, want %v", actx.Error(), authErr)
	}

	if got, want := res.StatusCode, http.StatusProxyAuthRequired; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Proxy-Authenticate"), "Basic"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Proxy-Authenticate", got, want)
	}
}
