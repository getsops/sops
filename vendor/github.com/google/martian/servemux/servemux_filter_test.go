// Copyright 2016 Google Inc. All rights reserved.
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
package servemux

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
)

func TestModifyRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("example.com/test", func(rw http.ResponseWriter, req *http.Request) {
		return
	})

	f := NewFilter(mux)
	tm := martiantest.NewModifier()
	f.RequestWhenTrue(tm)
	fm := martiantest.NewModifier()
	f.RequestWhenFalse(fm)

	req, err := http.NewRequest("GET", "http://example.com/test", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := tm.RequestModified(), true; got != want {
		t.Errorf("tm.RequestModified(): got %v, want %v", got, want)
	}

	if got, want := fm.RequestModified(), false; got != want {
		t.Errorf("fm.RequestModified(): got %v, want %v", got, want)
	}

	tm.Reset()
	fm.Reset()

	req, err = http.NewRequest("GET", "http://example.com/nomatch", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := f.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := tm.RequestModified(), false; got != want {
		t.Errorf("tm.RequestModified(): got %v, want %v", got, want)
	}

	if got, want := fm.RequestModified(), true; got != want {
		t.Errorf("fm.RequestModified(): got %v, want %v", got, want)
	}
}

func TestModifyResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("example.com/restest", func(rw http.ResponseWriter, req *http.Request) {
		return
	})

	f := NewFilter(mux)
	tm := martiantest.NewModifier()
	f.ResponseWhenTrue(tm)
	fm := martiantest.NewModifier()
	f.ResponseWhenFalse(fm)

	req, err := http.NewRequest("GET", "http://example.com/restest", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	res := proxyutil.NewResponse(200, nil, req)

	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := tm.ResponseModified(), true; got != want {
		t.Errorf("tm.ResponseModified(): got %v, want %v", got, want)
	}

	if got, want := fm.ResponseModified(), false; got != want {
		t.Errorf("fm.ResponseModified(): got %v, want %v", got, want)
	}

	tm.Reset()
	fm.Reset()

	req, err = http.NewRequest("GET", "http://example.com/nomatch", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	res = proxyutil.NewResponse(200, nil, req)

	if err := f.ModifyResponse(res); err != nil {
		t.Errorf("ModifyResponse(): got %v, want no error", err)
	}

	if tm.ResponseModified() != false {
		t.Errorf("tm.ResponseModified(): got %t, want %t", tm.ResponseModified(), false)
	}

	if got, want := tm.ResponseModified(), false; got != want {
		t.Errorf("tm.ResponseModified(): got %v, want %v", got, want)
	}

	if got, want := fm.ResponseModified(), true; got != want {
		t.Errorf("fm.ResponseModified(): got %v, want %v", got, want)
	}
}
