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

package cookie

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestCookieModifier(t *testing.T) {
	cookie := &http.Cookie{
		Name:  "name",
		Value: "value",
	}

	mod := NewModifier(cookie)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := len(req.Cookies()), 1; got != want {
		t.Errorf("len(req.Cookies): got %v, want %v", got, want)
	}
	if got, want := req.Cookies()[0].Name, cookie.Name; got != want {
		t.Errorf("req.Cookies()[0].Name: got %v, want %v", got, want)
	}
	if got, want := req.Cookies()[0].Value, cookie.Value; got != want {
		t.Errorf("req.Cookies()[0].Value: got %v, want %v", got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := len(res.Cookies()), 1; got != want {
		t.Errorf("len(res.Cookies): got %v, want %v", got, want)
	}
	if got, want := res.Cookies()[0].Name, cookie.Name; got != want {
		t.Errorf("res.Cookies()[0].Name: got %v, want %v", got, want)
	}
	if got, want := res.Cookies()[0].Value, cookie.Value; got != want {
		t.Errorf("res.Cookies()[0].Value: got %v, want %v", got, want)
	}
}

func TestModifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "cookie.Modifier": {
      "scope": ["request", "response"],
		  "name": "martian",
			"value": "value"
		}
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "http://example.com/path/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()

	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := len(req.Cookies()), 1; got != want {
		t.Fatalf("len(req.Cookies): got %v, want %v", got, want)
	}
	if got, want := req.Cookies()[0].Name, "martian"; got != want {
		t.Errorf("req.Cookies()[0].Name: got %v, want %v", got, want)
	}
	if got, want := req.Cookies()[0].Value, "value"; got != want {
		t.Errorf("req.Cookies()[0].Value: got %v, want %v", got, want)
	}

	resmod := r.ResponseModifier()

	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := len(res.Cookies()), 1; got != want {
		t.Fatalf("len(res.Cookies): got %v, want %v", got, want)
	}
	if got, want := res.Cookies()[0].Name, "martian"; got != want {
		t.Errorf("res.Cookies()[0].Name: got %v, want %v", got, want)
	}
	if got, want := res.Cookies()[0].Value, "value"; got != want {
		t.Errorf("res.Cookies()[0].Value: got %v, want %v", got, want)
	}
}
