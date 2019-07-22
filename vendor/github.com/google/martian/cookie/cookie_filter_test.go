// Copyright 2017 Google Inc. All rights reserved.
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

	"github.com/google/martian/v3/filter"
	_ "github.com/google/martian/v3/header"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestFilterFromJSON(t *testing.T) {
	msg := []byte(`{
		"cookie.Filter": {
			"scope": ["request", "response"],
			"name": "martian-cookie",
			"value": "true",
			"modifier": {
				"header.Modifier" : {
					"scope": ["request", "response"],
					"name": "Martian-Testing",
					"value": "true"
				}
			},
			"else": {
				"header.Modifier" : {
					"scope": ["request", "response"],
					"name": "Martian-Testing",
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

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	for _, tc := range []struct {
		name      string
		wantMatch bool
		cookie    *http.Cookie
	}{
		{
			name:      "matching name and value",
			wantMatch: true,
			cookie: &http.Cookie{
				Name:  "martian-cookie",
				Value: "true",
			},
		},
		{
			name:      "matching name with mismatched value",
			wantMatch: false,
			cookie: &http.Cookie{
				Name:  "martian-cookie",
				Value: "false",
			},
		},
		{
			name:      "missing cookie",
			wantMatch: false,
		},
	} {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Errorf("%s: http.NewRequest(): got %v, want no error", tc.name, err)
			continue
		}
		if tc.cookie != nil {
			req.AddCookie(tc.cookie)
		}

		if err := reqmod.ModifyRequest(req); err != nil {
			t.Errorf("%s: ModifyRequest(): got %v, want no error", tc.name, err)
			continue
		}

		want := "false"
		if tc.wantMatch {
			want = "true"
		}
		if got := req.Header.Get("Martian-Testing"); got != want {
			t.Errorf("%s: req.Header.Get(%q): got %q, want %q", "Martian-Testing", tc.name, got, want)
			continue
		}

		res := proxyutil.NewResponse(200, nil, req)
		if tc.cookie != nil {
			c := &http.Cookie{Name: tc.cookie.Name, Value: tc.cookie.Value}
			res.Header.Add("Set-Cookie", c.String())
		}

		if err := resmod.ModifyResponse(res); err != nil {
			t.Fatalf("ModifyResponse(): got %v, want no error", err)
		}

		if got := res.Header.Get("Martian-Testing"); got != want {
			t.Fatalf("res.Header.Get(%q): got %q, want %q", "Martian-Testing", got, want)
		}

	}
}

func TestFilterFromJSONWithoutElse(t *testing.T) {
	msg := []byte(`{
		"cookie.Filter": {
			"scope": ["request", "response"],
			"name": "martian-cookie",
			"value": "true",
			"modifier": {
				"header.Modifier" : {
					"scope": ["request", "response"],
					"name": "Martian-Testing",
					"value": "true"
				}
			}
		}
	}`)
	_, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}
}

func TestRequestWhenTrueCondition(t *testing.T) {
	cm := NewMatcher(&http.Cookie{Name: "Martian-Testing", Value: "true"})

	tt := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "Martian-Production",
			value: "true",
			want:  false,
		},
		{
			name:  "Martian-Testing",
			value: "true",
			want:  true,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetRequestCondition(cm)
		f.RequestWhenTrue(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		req.AddCookie(&http.Cookie{Name: tc.name, Value: tc.value})

		if err := f.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}

		if tm.RequestModified() != tc.want {
			t.Errorf("%d. tm.RequestModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestRequestWhenFalse(t *testing.T) {
	cm := NewMatcher(&http.Cookie{Name: "Martian-Testing", Value: "true"})
	tt := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "Martian-Production",
			value: "true",
			want:  true,
		},
		{
			name:  "Martian-Testing",
			value: "true",
			want:  false,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetRequestCondition(cm)
		f.RequestWhenFalse(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		req.AddCookie(&http.Cookie{Name: tc.name, Value: tc.value})

		if err := f.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}

		if tm.RequestModified() != tc.want {
			t.Errorf("%d. tm.RequestModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestResponseWhenTrue(t *testing.T) {
	cm := NewMatcher(&http.Cookie{Name: "Martian-Testing", Value: "true"})

	tt := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "Martian-Production",
			value: "true",
			want:  false,
		},
		{
			name:  "Martian-Testing",
			value: "true",
			want:  true,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetResponseCondition(cm)
		f.ResponseWhenTrue(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		res := proxyutil.NewResponse(200, nil, req)

		c := &http.Cookie{Name: tc.name, Value: tc.value}
		res.Header.Add("Set-Cookie", c.String())

		if err := f.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}

		if tm.ResponseModified() != tc.want {
			t.Errorf("%d. tm.ResponseModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestResponseWhenFalse(t *testing.T) {
	cm := NewMatcher(&http.Cookie{Name: "Martian-Testing", Value: "true"})

	tt := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "Martian-Production",
			value: "true",
			want:  true,
		},
		{
			name:  "Martian-Testing",
			value: "true",
			want:  false,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetResponseCondition(cm)
		f.ResponseWhenFalse(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		res := proxyutil.NewResponse(200, nil, req)

		c := &http.Cookie{Name: tc.name, Value: tc.value}
		res.Header.Add("Set-Cookie", c.String())

		if err := f.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}

		if tm.ResponseModified() != tc.want {
			t.Errorf("%d. tm.ResponseModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}
