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

package martianhttp

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/martiantest"

	_ "github.com/google/martian/v3/header"
)

func TestIntegration(t *testing.T) {
	ptr := martiantest.NewTransport()

	proxy := martian.NewProxy()
	defer proxy.Close()

	proxy.SetRoundTripper(ptr)

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	go proxy.Serve(l)

	m := NewModifier()
	proxy.SetRequestModifier(m)
	proxy.SetResponseModifier(m)

	mux := http.NewServeMux()
	mux.Handle("/", m)

	s := httptest.NewServer(mux)
	defer s.Close()

	body := strings.NewReader(`{
		"header.Modifier": {
      "scope": ["request", "response"],
			"name": "Martian-Test",
			"value": "true"
		}
	}`)

	res, err := http.Post(s.URL, "application/json", body)
	if err != nil {
		t.Fatalf("http.Post(%s): got %v, want no error", s.URL, err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	tr := &http.Transport{
		Proxy: http.ProxyURL(&url.URL{
			Scheme: "http",
			Host:   l.Addr().String(),
		}),
	}
	defer tr.CloseIdleConnections()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Connection", "close")

	res, err = tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("transport.RoundTrip(%q): got %v, want no error", req.URL, err)
	}
	res.Body.Close()

	if got, want := res.Header.Get("Martian-Test"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-Test", got, want)
	}
}
