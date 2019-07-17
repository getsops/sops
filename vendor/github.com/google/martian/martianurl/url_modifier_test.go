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
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/martian/v3/parse"
)

func TestNewModifier(t *testing.T) {
	tt := []struct {
		want string
		url  *url.URL
	}{
		{
			want: "https://www.example.com",
			url:  &url.URL{Scheme: "https"},
		},
		{
			want: "http://www.martian.local",
			url:  &url.URL{Host: "www.martian.local"},
		},
		{
			want: "http://www.example.com/test",
			url:  &url.URL{Path: "/test"},
		},
		{
			want: "http://www.example.com?test=true",
			url:  &url.URL{RawQuery: "test=true"},
		},
		{
			want: "http://www.example.com#test",
			url:  &url.URL{Fragment: "test"},
		},
		{
			want: "https://martian.local/test?test=true#test",
			url: &url.URL{
				Scheme:   "https",
				Host:     "martian.local",
				Path:     "/test",
				RawQuery: "test=true",
				Fragment: "test",
			},
		},
	}

	for i, tc := range tt {
		req, err := http.NewRequest("GET", "http://www.example.com", nil)
		if err != nil {
			t.Fatalf("%d. NewRequest(): got %v, want no error", i, err)
		}

		mod := NewModifier(tc.url)

		if err := mod.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %q, want no error", i, err)
		}

		if got := req.URL.String(); got != tc.want {
			t.Errorf("%d. req.URL: got %q, want %q", i, got, tc.want)
		}
	}
}

func TestIntegration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = r.Host
			w.Header().Set("Martian-URL", r.URL.String())
		}))
	defer server.Close()

	u := &url.URL{
		Scheme: "http",
		Host:   server.Listener.Addr().String(),
	}
	m := NewModifier(u)

	req, err := http.NewRequest("GET", "https://example.com/test", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(%q, %q, nil): got %v, want no error", "GET", "http://example.com/test", err)
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("http.DefaultClient.Do(): got %v, want no error", err)
	}

	want := "http://example.com/test"
	if got := res.Header.Get("Martian-URL"); got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-URL", got, want)
	}
}

func TestModifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "url.Modifier": {
      "scope": ["request"],
      "scheme": "https",
      "host": "www.martian.proxy",
      "path": "/testing",
      "query": "test=true"
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

	if got, want := req.URL.Scheme, "https"; got != want {
		t.Errorf("req.URL.Scheme: got %q, want %q", got, want)
	}
	if got, want := req.URL.Host, "www.martian.proxy"; got != want {
		t.Errorf("req.URL.Host: got %q, want %q", got, want)
	}
	if got, want := req.URL.Path, "/testing"; got != want {
		t.Errorf("req.URL.Path: got %q, want %q", got, want)
	}
	if got, want := req.URL.RawQuery, "test=true"; got != want {
		t.Errorf("req.URL.RawQuery: got %q, want %q", got, want)
	}
}
