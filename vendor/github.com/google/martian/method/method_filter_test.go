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

package method

import (
	"net/http"
	"testing"

	_ "github.com/google/martian/v3/header"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestFilterModifyRequest(t *testing.T) {
	tt := []struct {
		method string
		want   bool
	}{
		{
			method: "GET",
			want:   true,
		},
		{
			method: "get",
			want:   true,
		},
		{
			method: "POST",
			want:   false,
		},
		{
			method: "DELETE",
			want:   false,
		},
		{
			method: "CONNECT",
			want:   false,
		},
		{
			method: "connect",
			want:   false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}

		mod := NewFilter(tc.method)
		tm := martiantest.NewModifier()
		mod.SetRequestModifier(tm)

		if err := mod.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %q, want no error", i, err)
		}

		if tm.RequestModified() != tc.want {
			t.Errorf("%d. tm.RequestModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestFilterModifyResponse(t *testing.T) {
	tt := []struct {
		method string
		want   bool
	}{
		{
			method: "GET",
			want:   true,
		},
		{
			method: "get",
			want:   true,
		},
		{
			method: "POST",
			want:   false,
		},
		{
			method: "DELETE",
			want:   false,
		},
		{
			method: "CONNECT",
			want:   false,
		},
		{
			method: "connect",
			want:   false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}
		res := proxyutil.NewResponse(200, nil, req)

		mod := NewFilter(tc.method)
		tm := martiantest.NewModifier()
		mod.SetResponseModifier(tm)

		if err := mod.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %q, want no error", i, err)
		}

		if tm.ResponseModified() != tc.want {
			t.Errorf("%d. tm.ResponseModified(): got %t, want %t", i, tm.ResponseModified(), tc.want)
		}
	}

}

func TestFilterFromJSON(t *testing.T) {
	j := `{
		    "method.Filter": {
              "scope": ["request", "response"],
              "method": "GET",
              "modifier": {
                "header.Modifier": {
                  "scope": ["request", "response"],
                  "name": "Mod-Run",
                  "value": "true"
                } 
		      },
		      "else": {
                "header.Modifier": {
                  "scope": ["request", "response"],
                  "name": "Else-Run",
                  "value": "true"
                } 
              }
            }
	      }`

	msg := []byte(j)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("FilterFromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Mod-Run"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Mod-Run", got, want)
	}

	if got, want := req.Header.Get("Else-Run"), ""; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Else-Run", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Mod-Run"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Mod-Run", got, want)
	}

	// test else conditional modifier with POST
	req, err = http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Mod-Run"), ""; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Mod-Run", got, want)
	}

	if got, want := req.Header.Get("Else-Run"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Else-Run", got, want)
	}
}
