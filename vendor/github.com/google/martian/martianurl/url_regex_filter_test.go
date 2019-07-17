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
	"net/http"
	"regexp"
	"testing"

	"github.com/google/martian/v3"
	_ "github.com/google/martian/v3/header"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestRegexFilterModifyRequest(t *testing.T) {
	tt := []struct {
		want   bool
		match  string
		regstr string
	}{
		{
			match:  "https://www.example.com",
			regstr: "https://.*",
			want:   true,
		},
		{
			match:  "http://www.example.com/subpath",
			regstr: ".*www.example.com.*",
			want:   true,
		},
		{
			match:  "https://www.example.com/subpath",
			regstr: ".*www.example.com.*",
			want:   true,
		},
		{
			match:  "http://www.example.com/test",
			regstr: ".*/test",
			want:   true,
		},
		{
			match:  "http://www.example.com?test=true",
			regstr: ".*test=true.*",
			want:   true,
		},
		{
			match:  "http://www.example.com#test",
			regstr: ".*test.*",
			want:   true,
		},
		{
			match:  "https://martian.local/test?test=true#test",
			regstr: "https://martian.local/test\\?test=true#test",
			want:   true,
		},
		{
			match:  "http://www.youtube.com/get_tags?tagone=yes",
			regstr: ".*www.youtube.com/get_tags\\?.*",
			want:   true,
		},
		{
			match:  "https://www.example.com",
			regstr: "http://.*",
			want:   false,
		},
		{
			match:  "http://www.martian.external",
			regstr: ".*www.martian.local.*",
			want:   false,
		},
		{
			match:  "http://www.example.com/testing",
			regstr: ".*/test$",
			want:   false,
		},
		{
			match:  "http://www.example.com?test=false",
			regstr: ".*test=true.*",
			want:   false,
		},
		{
			match:  "http://www.example.com#test",
			regstr: ".*#testing.*",
			want:   false,
		},
		{
			match: "https://martian.local/test?test=true#test",
			// "\\\\" was the old way of adding a backslash in SAVR
			regstr: "https://martian.local/test\\\\?test=true#test",
			want:   false,
		},
		{
			match:  "http://www.youtube.com/get_tags/nope",
			regstr: ".*www.youtube.com/get_ad_tags\\?.*",
			want:   false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", tc.match, nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}
		regex, err := regexp.Compile(tc.regstr)
		if err != nil {
			t.Fatalf("%d. regexp.Compile(): got %v, want no error", i, err)
		}

		var modRun bool
		mod := NewRegexFilter(regex)
		mod.SetRequestModifier(martian.RequestModifierFunc(
			func(*http.Request) error {
				modRun = true
				return nil
			}))

		if err := mod.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %q, want no error", i, err)
		}

		if modRun != tc.want {
			t.Errorf("%d. modRun: got %t, want %t", i, modRun, tc.want)
		}
	}
}

// The matching functionality is already tested above, so this just tests response setting.
func TestRegexFilterModifyResponse(t *testing.T) {
	tt := []struct {
		want   bool
		match  string
		regstr string
	}{
		{
			match:  "https://www.example.com",
			regstr: ".*www.example.com.*",
			want:   true,
		},
		{
			match:  "http://www.martian.external",
			regstr: ".*www.martian.local.*",
			want:   false,
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", tc.match, nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}
		res := proxyutil.NewResponse(200, nil, req)
		regex, err := regexp.Compile(tc.regstr)
		if err != nil {
			t.Fatalf("%d. regexp.Compile(): got %v, want no error", i, err)
		}

		var modRun bool
		mod := NewRegexFilter(regex)
		mod.SetResponseModifier(martian.ResponseModifierFunc(
			func(*http.Response) error {
				modRun = true
				return nil
			}))

		if err := mod.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %q, want no error", i, err)
		}

		if modRun != tc.want {
			t.Errorf("%d. modRun: got %t, want %t", i, modRun, tc.want)
		}
	}
}

func TestRegexFilterFromJSON(t *testing.T) {
	rawMsg := `
	{
		"url.RegexFilter": {
      "scope": ["request", "response"],
			"regex": ".*martian.test.*",
			"modifier": {
				"header.Modifier": {
					"name": "Martian-Test",
					"value": "true",
          "scope": ["request", "response"]
				}
			}
		}
	}`

	r, err := parse.FromJSON([]byte(rawMsg))
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Errorf("reqmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "https://martian.test", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("reqmod.ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Martian-Test"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Martian-Test", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Martian-Test"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-Test", got, want)
	}
}