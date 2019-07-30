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

package header

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3/filter"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestFilterFromJSON(t *testing.T) {
	msg := []byte(`{
		"header.Filter": {
			"scope": ["request", "response"],
			"name": "Martian-Passthrough",
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

	// Matching condition for request
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Martian-Passthrough", "true")
	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("Martian-Testing"), "true"; got != want {
		t.Fatalf("req.Header.Get(%q): got %q, want %q", "Martian-Testing", got, want)
	}

	// Else condition for request
	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Martian-Passthrough", "false")
	if err := reqmod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("Martian-Testing"), "false"; got != want {
		t.Fatalf("req.Header.Get(%q): got %q, want %q", "Martian-Testing", got, want)
	}

	// Matching condition for response
	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header.Set("Martian-Passthrough", "true")
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.Header.Get("Martian-Testing"), "true"; got != want {
		t.Fatalf("res.Header.Get(%q): got %q, want %q", "Martian-Testing", got, want)
	}

	// Else condition for response
	resmod = r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	res = proxyutil.NewResponse(200, nil, req)
	res.Header.Set("Martian-Passthrough", "false")
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.Header.Get("Martian-Testing"), "false"; got != want {
		t.Fatalf("res.Header.Get(%q): got %q, want %q", "Martian-Testing", got, want)
	}
}

func TestFilterFromJSONWithoutElse(t *testing.T) {
	msg := []byte(`{
		"header.Filter": {
			"scope": ["request", "response"],
			"name": "Martian-Passthrough",
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
	hm := NewMatcher("Martian-Testing", "true")

	tt := []struct {
		name   string
		values []string
		want   bool
	}{
		{
			name:   "Martian-Production",
			values: []string{"true"},
			want:   false,
		},
		{
			name:   "Martian-Testing",
			values: []string{"see-next-value", "true"},
			want:   true,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetRequestCondition(hm)
		f.RequestWhenTrue(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		req.Header[tc.name] = tc.values

		if err := f.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}

		if tm.RequestModified() != tc.want {
			t.Errorf("%d. tm.RequestModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestRequestWhenFalse(t *testing.T) {
	hm := NewMatcher("Martian-Testing", "true")
	tt := []struct {
		name   string
		values []string
		want   bool
	}{
		{
			name:   "Martian-Production",
			values: []string{"true"},
			want:   true,
		},
		{
			name:   "Martian-Testing",
			values: []string{"see-next-value", "true"},
			want:   false,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetRequestCondition(hm)
		f.RequestWhenFalse(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		req.Header[tc.name] = tc.values

		if err := f.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}

		if tm.RequestModified() != tc.want {
			t.Errorf("%d. tm.RequestModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestResponseWhenTrue(t *testing.T) {
	hm := NewMatcher("Martian-Testing", "true")

	tt := []struct {
		name   string
		values []string
		want   bool
	}{
		{
			name:   "Martian-Production",
			values: []string{"true"},
			want:   false,
		},
		{
			name:   "Martian-Testing",
			values: []string{"see-next-value", "true"},
			want:   true,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetResponseCondition(hm)
		f.ResponseWhenTrue(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		res := proxyutil.NewResponse(200, nil, req)

		res.Header[tc.name] = tc.values

		if err := f.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}

		if tm.ResponseModified() != tc.want {
			t.Errorf("%d. tm.ResponseModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}

func TestResponseWhenFalse(t *testing.T) {
	hm := NewMatcher("Martian-Testing", "true")

	tt := []struct {
		name   string
		values []string
		want   bool
	}{
		{
			name:   "Martian-Production",
			values: []string{"true"},
			want:   true,
		},
		{
			name:   "Martian-Testing",
			values: []string{"see-next-value", "true"},
			want:   false,
		},
	}

	for i, tc := range tt {
		tm := martiantest.NewModifier()

		f := filter.New()
		f.SetResponseCondition(hm)
		f.ResponseWhenFalse(tm)

		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}
		res := proxyutil.NewResponse(200, nil, req)

		res.Header[tc.name] = tc.values

		if err := f.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}

		if tm.ResponseModified() != tc.want {
			t.Errorf("%d. tm.ResponseModified(): got %t, want %t", i, tm.RequestModified(), tc.want)
		}
	}
}
