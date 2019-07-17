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

package ipauth

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/auth"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
)

func TestModifyRequest(t *testing.T) {
	m := NewModifier()
	m.SetRequestModifier(nil)

	req, err := http.NewRequest("CONNECT", "https://www.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	actx := auth.FromContext(ctx)

	if got, want := actx.ID(), ""; got != want {
		t.Errorf("actx.ID(): got %q, want %q", got, want)
	}

	// IP with port and modifier with error.
	tm := martiantest.NewModifier()
	reqerr := errors.New("request error")
	tm.RequestError(reqerr)

	req.RemoteAddr = "1.1.1.1:8111"
	m.SetRequestModifier(tm)

	if err := m.ModifyRequest(req); err != reqerr {
		t.Fatalf("ModifyConnectRequest(): got %v, want %v", err, reqerr)
	}

	if got, want := actx.ID(), "1.1.1.1"; got != want {
		t.Errorf("actx.ID(): got %q, want %q", got, want)
	}

	// IP without port and modifier with auth error.
	req.RemoteAddr = "4.4.4.4"

	authErr := errors.New("auth error")
	tm.RequestError(nil)
	tm.RequestFunc(func(req *http.Request) {
		ctx := martian.NewContext(req)
		actx := auth.FromContext(ctx)

		actx.SetError(authErr)
	})

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := actx.ID(), ""; got != want {
		t.Errorf("actx.ID(): got %q, want %q", got, want)
	}
}

func TestModifyResponse(t *testing.T) {
	m := NewModifier()
	m.SetResponseModifier(nil)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(200, nil, req)
	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	// Modifier with error.
	tm := martiantest.NewModifier()
	reserr := errors.New("response error")
	tm.ResponseError(reserr)

	m.SetResponseModifier(tm)
	if err := m.ModifyResponse(res); err != reserr {
		t.Fatalf("ModifyResponse(): got %v, want %v", err, reserr)
	}

	// Modifier with auth error.
	tm.ResponseError(nil)
	authErr := errors.New("auth error")
	tm.ResponseFunc(func(res *http.Response) {
		ctx := martian.NewContext(res.Request)
		actx := auth.FromContext(ctx)

		actx.SetError(authErr)
	})

	actx := auth.FromContext(ctx)
	actx.SetID("bad-auth")

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.StatusCode, 403; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}

	if got, want := actx.Error(), authErr; got != want {
		t.Errorf("actx.Error(): got %v, want %v", got, want)
	}
}
