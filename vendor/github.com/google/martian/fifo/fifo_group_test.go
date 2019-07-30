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

package fifo

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"

	_ "github.com/google/martian/v3/header"
)

func TestGroupFromJSON(t *testing.T) {
	msg := []byte(`{
    "fifo.Group": {
      "scope": ["request", "response"],
      "aggregateErrors": true,
      "modifiers": [
        {
          "header.Modifier": {
            "scope": ["request", "response"],
            "name": "X-Testing",
            "value": "true"
          }
        },
        {
          "header.Modifier": {
            "scope": ["request", "response"],
            "name": "Y-Testing",
            "value": "true"
          }
        }
      ]
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
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("X-Testing"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "X-Testing", got, want)
	}
	if got, want := req.Header.Get("Y-Testing"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Y-Testing", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}
	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.Header.Get("X-Testing"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "X-Testing", got, want)
	}
	if got, want := res.Header.Get("Y-Testing"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Y-Testing", got, want)
	}
}

func TestModifyRequest(t *testing.T) {
	fg := NewGroup()
	tm := martiantest.NewModifier()

	fg.AddRequestModifier(tm)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := fg.ModifyRequest(req); err != nil {
		t.Fatalf("fg.ModifyRequest(): got %v, want no error", err)
	}
	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}
}

func TestModifyRequestHaltsOnError(t *testing.T) {
	fg := NewGroup()

	reqerr := errors.New("request error")
	tm := martiantest.NewModifier()
	tm.RequestError(reqerr)
	fg.AddRequestModifier(tm)

	tm2 := martiantest.NewModifier()
	fg.AddRequestModifier(tm2)

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := fg.ModifyRequest(req); err != reqerr {
		t.Fatalf("fg.ModifyRequest(): got %v, want %v", err, reqerr)
	}

	if tm2.RequestModified() {
		t.Error("tm2.RequestModified(): got true, want false")
	}
}

func TestModifyRequestAggregatesErrors(t *testing.T) {
	fg := NewGroup()
	fg.SetAggregateErrors(true)

	reqerr1 := errors.New("1. request error")
	tm := martiantest.NewModifier()
	tm.RequestError(reqerr1)
	fg.AddRequestModifier(tm)

	tm2 := martiantest.NewModifier()
	reqerr2 := errors.New("2. request error")
	tm2.RequestError(reqerr2)
	fg.AddRequestModifier(tm2)

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	merr := martian.NewMultiError()
	merr.Add(reqerr1)
	merr.Add(reqerr2)

	if err := fg.ModifyRequest(req); err == nil {
		t.Fatalf("fg.ModifyRequest(): got %v, want not nil", err)
	}
	if err := fg.ModifyRequest(req); err.Error() != merr.Error() {
		t.Fatalf("fg.ModifyRequest(): got %v, want %v", err, merr)
	}

	if err, want := fg.ModifyRequest(req), "1. request error\n2. request error"; err.Error() != want {
		t.Fatalf("fg.ModifyRequest(): got %v, want %v", err, want)
	}
}

func TestModifyResponse(t *testing.T) {
	fg := NewGroup()
	tm := martiantest.NewModifier()

	fg.AddResponseModifier(tm)

	res := proxyutil.NewResponse(200, nil, nil)
	if err := fg.ModifyResponse(res); err != nil {
		t.Fatalf("fg.ModifyResponse(): got %v, want no error", err)
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
}

func TestModifyResponseHaltsOnError(t *testing.T) {
	fg := NewGroup()

	reserr := errors.New("request error")
	tm := martiantest.NewModifier()
	tm.ResponseError(reserr)
	fg.AddResponseModifier(tm)

	tm2 := martiantest.NewModifier()
	fg.AddResponseModifier(tm2)

	res := proxyutil.NewResponse(200, nil, nil)
	if err := fg.ModifyResponse(res); err != reserr {
		t.Fatalf("fg.ModifyResponse(): got %v, want %v", err, reserr)
	}

	if tm2.ResponseModified() {
		t.Error("tm2.ResponseModified(): got true, want false")
	}
}

func TestModifyResponseAggregatesErrors(t *testing.T) {
	fg := NewGroup()
	fg.SetAggregateErrors(true)

	reserr1 := errors.New("1. response error")
	tm := martiantest.NewModifier()
	tm.ResponseError(reserr1)
	fg.AddResponseModifier(tm)

	tm2 := martiantest.NewModifier()
	reserr2 := errors.New("2. response error")
	tm2.ResponseError(reserr2)
	fg.AddResponseModifier(tm2)

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(200, nil, req)

	merr := martian.NewMultiError()
	merr.Add(reserr1)
	merr.Add(reserr2)

	if err := fg.ModifyResponse(res); err == nil {
		t.Fatalf("fg.ModifyResponse(): got %v, want %v", err, merr)
	}

	if err := fg.ModifyResponse(res); err.Error() != merr.Error() {
		t.Fatalf("fg.ModifyResponse(): got %v, want %v", err, merr)
	}
}

func TestVerifyRequests(t *testing.T) {
	fg := NewGroup()

	if err := fg.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequest(): got %v, want no error", err)
	}

	errs := []error{}
	for i := 0; i < 3; i++ {
		err := fmt.Errorf("%d. verify request failure", i)

		tv := &verify.TestVerifier{
			RequestError: err,
		}
		fg.AddRequestModifier(tv)

		errs = append(errs, err)
	}

	merr, ok := fg.VerifyRequests().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyRequests(): got nil, want *verify.MultiError")
	}

	if !reflect.DeepEqual(merr.Errors(), errs) {
		t.Errorf("merr.Errors(): got %v, want %v", merr.Errors(), errs)
	}
}

func TestVerifyResponses(t *testing.T) {
	fg := NewGroup()

	if err := fg.VerifyResponses(); err != nil {
		t.Fatalf("VerifyResponses(): got %v, want no error", err)
	}

	errs := []error{}
	for i := 0; i < 3; i++ {
		err := fmt.Errorf("%d. verify responses failure", i)

		tv := &verify.TestVerifier{
			ResponseError: err,
		}
		fg.AddResponseModifier(tv)

		errs = append(errs, err)
	}

	merr, ok := fg.VerifyResponses().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyResponses(): got nil, want *verify.MultiError")
	}

	if !reflect.DeepEqual(merr.Errors(), errs) {
		t.Errorf("merr.Errors(): got %v, want %v", merr.Errors(), errs)
	}
}

func TestResets(t *testing.T) {
	fg := NewGroup()

	for i := 0; i < 3; i++ {
		tv := &verify.TestVerifier{
			RequestError:  fmt.Errorf("%d. verify request error", i),
			ResponseError: fmt.Errorf("%d. verify response error", i),
		}
		fg.AddRequestModifier(tv)
		fg.AddResponseModifier(tv)
	}

	if err := fg.VerifyRequests(); err == nil {
		t.Fatal("VerifyRequests(): got nil, want error")
	}
	if err := fg.VerifyResponses(); err == nil {
		t.Fatal("VerifyResponses(): got nil, want error")
	}

	fg.ResetRequestVerifications()
	fg.ResetResponseVerifications()

	if err := fg.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
	if err := fg.VerifyResponses(); err != nil {
		t.Errorf("VerifyResponses(): got %v, want no error", err)
	}
}
