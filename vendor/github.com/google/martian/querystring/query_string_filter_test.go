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

package querystring

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"

	// Import to register header.Modifier with JSON parser.
	_ "github.com/google/martian/v3/header"
)

func TestNoModifiers(t *testing.T) {
	f := NewFilter("", "")
	f.SetRequestModifier(nil)
	f.SetResponseModifier(nil)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}
}

func TestQueryStringFilterWithQuery(t *testing.T) {
	// Name only, no value.
	f := NewFilter("match", "")

	tm := martiantest.NewModifier()
	f.SetRequestModifier(tm)
	f.SetResponseModifier(tm)

	req, err := http.NewRequest("GET", "http://martian.local?match=any", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
	tm.Reset()

	req, err = http.NewRequest("GET", "http://martian.local?nomatch", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res = proxyutil.NewResponse(200, nil, req)
	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if tm.RequestModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}
	if tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got true, want false")
	}
	tm.Reset()

	// Name and value.
	f = NewFilter("match", "value")
	f.SetRequestModifier(tm)
	f.SetResponseModifier(tm)

	req, err = http.NewRequest("GET", "http://martian.local?match=value", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res = proxyutil.NewResponse(200, nil, req)
	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
	tm.Reset()

	req, err = http.NewRequest("GET", "http://martian.local?match=notvalue", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res = proxyutil.NewResponse(200, nil, req)
	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if tm.RequestModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}
	if tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got true, want false")
	}
	tm.Reset()

	// Explicitly do not match POST data.
	req, err = http.NewRequest("GET", "http://martian.local", strings.NewReader("match=value"))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	res = proxyutil.NewResponse(200, nil, req)
	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if tm.RequestModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}
	if tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got true, want false")
	}
	tm.Reset()
}

func TestFilterFromJSON(t *testing.T) {
	msg := []byte(`{
		"querystring.Filter": {
      "scope": ["request", "response"],
      "name": "param",
      "value": "true",
      "modifier": {
        "header.Modifier": {
          "scope": ["request", "response"],
          "name": "Martian-Modified",
          "value": "true"
        }
      }
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

	req, err := http.NewRequest("GET", "https://martian.test?param=true", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Martian-Modified"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Martian-Modified", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Martian-Modified"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-Modified", got, want)
	}
}

func TestElseCondition(t *testing.T) {
	msg := []byte(`{
		"querystring.Filter": {
      "scope": ["request", "response"],
      "name": "param",
      "value": "true",
      "modifier": {
        "header.Modifier": {
          "scope": ["request", "response"],
          "name": "Martian-Modified",
          "value": "true"
        }
      },
      "else": {
        "header.Modifier": {
          "scope": ["request", "response"],
          "name": "Martian-Modified",
          "value": "false"
        }
      }
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

	req, err := http.NewRequest("GET", "https://martian.test?param=false", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Martian-Modified"), "false"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Martian-Modified", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Martian-Modified"), "false"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-Modified", got, want)
	}
}

func TestVerifyRequests(t *testing.T) {
	f := NewFilter("", "")

	if err := f.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequest(): got %v, want no error", err)
	}

	tv := &verify.TestVerifier{
		RequestError: errors.New("verify request failure"),
	}

	f.SetRequestModifier(tv)

	want := martian.NewMultiError()
	want.Add(tv.RequestError)
	if got := f.VerifyRequests(); got.Error() != want.Error() {
		t.Fatalf("VerifyRequests(): got %v, want %v", got, want)
	}

	f.ResetRequestVerifications()

	if err := f.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequest(): got %v, want no error", err)
	}
}

func TestVerifyResponses(t *testing.T) {
	f := NewFilter("", "")

	if err := f.VerifyResponses(); err != nil {
		t.Fatalf("VerifyResponses(): got %v, want no error", err)
	}

	tv := &verify.TestVerifier{
		ResponseError: errors.New("verify response failure"),
	}

	f.SetResponseModifier(tv)

	want := martian.NewMultiError()
	want.Add(tv.ResponseError)
	if got := f.VerifyResponses(); got.Error() != want.Error() {
		t.Fatalf("VerifyResponses(): got %v, want %v", got, want)
	}

	f.ResetResponseVerifications()

	if err := f.VerifyResponses(); err != nil {
		t.Fatalf("VerifyResponses(): got %v, want no error", err)
	}
}
