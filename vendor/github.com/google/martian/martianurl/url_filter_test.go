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

package martianurl

import (
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"

	_ "github.com/google/martian/v3/header"
)

func TestFilterModifyRequest(t *testing.T) {
	tt := []struct {
		want  bool
		match string
		url   *url.URL
	}{
		{
			match: "https://www.example.com",
			url:   &url.URL{Scheme: "https"},
			want:  true,
		},
		{
			match: "http://www.martian.local",
			url:   &url.URL{Host: "*.martian.local"},
			want:  true,
		},
		{
			match: "http://www.example.com/test",
			url:   &url.URL{Path: "/test"},
			want:  true,
		},
		{
			match: "http://www.example.com?test=true",
			url:   &url.URL{RawQuery: "test=true"},
			want:  true,
		},
		{
			match: "http://www.example.com#test",
			url:   &url.URL{Fragment: "test"},
			want:  true,
		},
		{
			match: "https://martian.local/test?test=true#test",
			url: &url.URL{
				Scheme:   "https",
				Host:     "martian.local",
				Path:     "/test",
				RawQuery: "test=true",
				Fragment: "test",
			},
			want: true,
		},
		{
			match: "https://www.example.com",
			url:   &url.URL{Scheme: "http"},
			want:  false,
		},
		{
			match: "http://www.martian.external",
			url:   &url.URL{Host: "www.martian.local"},
			want:  false,
		},
		{
			match: "http://www.example.com/testing",
			url:   &url.URL{Path: "/test"},
			want:  false,
		},
		{
			match: "http://www.example.com?test=false",
			url:   &url.URL{RawQuery: "test=true"},
			want:  false,
		},
		{
			match: "http://www.example.com#test",
			url:   &url.URL{Fragment: "testing"},
			want:  false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", tc.match, nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}

		mod := NewFilter(tc.url)
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
		want  bool
		match string
		url   *url.URL
	}{
		{
			match: "https://www.example.com",
			url:   &url.URL{Scheme: "https"},
			want:  true,
		},
		{
			match: "http://www.martian.local",
			url:   &url.URL{Host: "www.martian.local"},
			want:  true,
		},
		{
			match: "http://www.example.com/test",
			url:   &url.URL{Path: "/test"},
			want:  true,
		},
		{
			match: "http://www.example.com?test=true",
			url:   &url.URL{RawQuery: "test=true"},
			want:  true,
		},
		{
			match: "http://www.example.com#test",
			url:   &url.URL{Fragment: "test"},
			want:  true,
		},
		{
			match: "https://martian.local/test?test=true#test",
			url: &url.URL{
				Scheme:   "https",
				Host:     "martian.local",
				Path:     "/test",
				RawQuery: "test=true",
				Fragment: "test",
			},
			want: true,
		},
		{
			match: "https://www.example.com",
			url:   &url.URL{Scheme: "http"},
			want:  false,
		},
		{
			match: "http://www.martian.external",
			url:   &url.URL{Host: "www.martian.local"},
			want:  false,
		},
		{
			match: "http://www.example.com/testing",
			url:   &url.URL{Path: "/test"},
			want:  false,
		},
		{
			match: "http://www.example.com?test=false",
			url:   &url.URL{RawQuery: "test=true"},
			want:  false,
		},
		{
			match: "http://www.example.com#test",
			url:   &url.URL{Fragment: "testing"},
			want:  false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", tc.match, nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}
		res := proxyutil.NewResponse(200, nil, req)

		mod := NewFilter(tc.url)
		tm := martiantest.NewModifier()
		mod.SetResponseModifier(tm)

		if err := mod.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %q, want no error", i, err)
		}

		if tm.ResponseModified() != tc.want {
			t.Errorf("tm.ResponseModified(): got %t, want %t", tm.ResponseModified(), tc.want)
		}
	}
}

func TestFilterFromJSON(t *testing.T) {
	msg := []byte(`{
		"url.Filter": {
          "scope": ["request", "response"],
          "scheme": "https",
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
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("FilterFromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("reqmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "https://martian.test", nil)
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

	// test else conditional modifier with scheme of http
	req, err = http.NewRequest("GET", "http://martian.test", nil)
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
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Mod-Run", got, want)
	}

}

func TestPassThroughVerifyRequests(t *testing.T) {
	u := &url.URL{Host: "www.martian.local"}
	f := NewFilter(u)

	if err := f.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequest(): got %v, want no error", err)
	}

	tv := &verify.TestVerifier{
		RequestError: errors.New("verify request failure"),
	}

	f.SetRequestModifier(tv)

	if got, want := f.VerifyRequests().Error(), "verify request failure"; got != want {
		t.Fatalf("VerifyRequests(): got %s, want %s", got, want)
	}
}

func TestPassThroughVerifyResponses(t *testing.T) {
	u := &url.URL{Host: "www.martian.local"}
	f := NewFilter(u)
	if err := f.VerifyResponses(); err != nil {
		t.Fatalf("VerifyResponses(): got %v, want no error", err)
	}

	tv := &verify.TestVerifier{
		ResponseError: errors.New("verify response failure"),
	}

	f.SetResponseModifier(tv)

	if got, want := f.VerifyResponses().Error(), "verify response failure"; got != want {
		t.Fatalf("VerifyResponses(): got %s, want %s", got, want)
	}
}

func TestResets(t *testing.T) {
	u := &url.URL{Host: "www.martian.local"}
	f := NewFilter(u)

	tv := &verify.TestVerifier{
		ResponseError: errors.New("verify response failure"),
	}
	f.SetResponseModifier(tv)

	tv = &verify.TestVerifier{
		RequestError: errors.New("verify request failure"),
	}
	f.SetRequestModifier(tv)

	if err := f.VerifyRequests(); err == nil {
		t.Fatal("VerifyRequests(): got nil, want error")
	}
	if err := f.VerifyResponses(); err == nil {
		t.Fatal("VerifyResponses(): got nil, want error")
	}

	f.ResetRequestVerifications()
	f.ResetResponseVerifications()

	if err := f.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
	if err := f.VerifyResponses(); err != nil {
		t.Errorf("VerifyResponses(): got %v, want no error", err)
	}
}
