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

	"github.com/google/martian/v3/proxyutil"
)

func TestRemoveHopByHopHeaders(t *testing.T) {
	m := NewHopByHopModifier()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header = http.Header{
		// Additional hop-by-hop headers are listed in the
		// Connection header.
		"Connection": []string{
			"X-Connection",
			"X-Hop-By-Hop, close",
		},

		// RFC hop-by-hop headers.
		"Keep-Alive":          []string{},
		"Proxy-Authenticate":  []string{},
		"Proxy-Authorization": []string{},
		"Te":                []string{},
		"Trailer":           []string{},
		"Transfer-Encoding": []string{},
		"Upgrade":           []string{},
		"Proxy-Connection":  []string{},

		// Hop-by-hop headers listed in the Connection header.
		"X-Connection": []string{},
		"X-Hop-By-Hop": []string{},

		// End-to-end header that should not be removed.
		"X-End-To-End": []string{},
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := len(req.Header), 1; got != want {
		t.Fatalf("len(req.Header): got %d, want %d", got, want)
	}
	if _, ok := req.Header["X-End-To-End"]; !ok {
		t.Errorf("req.Header[%q]: got !ok, want ok", "X-End-To-End")
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header = req.Header
	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := len(res.Header), 1; got != want {
		t.Fatalf("len(res.Header): got %d, want %d", got, want)
	}
	if _, ok := res.Header["X-End-To-End"]; !ok {
		t.Errorf("res.Header[%q]: got !ok, want ok", "X-End-To-End")
	}
}
