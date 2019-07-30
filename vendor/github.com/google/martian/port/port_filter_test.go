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

package port

import (
	"net/http"
	"net/url"
	"testing"

	_ "github.com/google/martian/v3/header"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestFilterModifyRequest(t *testing.T) {
	tt := []struct {
		want bool
		url  *url.URL
		port int
	}{
		{
			url:  &url.URL{Scheme: "http", Host: "example.com"},
			port: 80,
			want: true,
		},
		{
			url:  &url.URL{Scheme: "http", Host: "example.com:80"},
			port: 80,
			want: true,
		},
		{
			url:  &url.URL{Scheme: "http", Host: "example.com"},
			port: 123,
			want: false,
		},
		{
			url:  &url.URL{Scheme: "http", Host: "example.com:8080"},
			port: 123,
			want: false,
		},
		{
			url:  &url.URL{Scheme: "https", Host: "example.com"},
			port: 443,
			want: true,
		},
		{
			url:  &url.URL{Scheme: "https", Host: "example.com:443"},
			port: 443,
			want: true,
		},
		{
			url:  &url.URL{Scheme: "https", Host: "example.com"},
			port: 123,
			want: false,
		},
		{
			url:  &url.URL{Scheme: "https", Host: "example.com:8080"},
			port: 123,
			want: false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", tc.url.String(), nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}

		mod := NewFilter(tc.port)
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

func TestFilterFromJSON(t *testing.T) {
	msg := []byte(`{
		"port.Filter": {
          "scope": ["request", "response"],
          "port": 8080,
          "modifier": {
            "header.Modifier": {
              "scope": ["request", "response"],
              "name": "Mod-Run",
              "value": "true"
          }
        }
      }
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("FilterFromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "https://example.com:8080", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Mod-Run"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Mod-Run", got, want)
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
}
