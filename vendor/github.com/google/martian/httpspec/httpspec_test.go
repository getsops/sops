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

package httpspec

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
)

func TestNewStack(t *testing.T) {
	stack, fg := NewStack("martian")

	tm := martiantest.NewModifier()
	fg.AddRequestModifier(tm)
	fg.AddResponseModifier(tm)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	// Hop-by-hop header to be removed.
	req.Header.Set("Hop-By-Hop", "true")
	req.Header.Set("Connection", "Hop-By-Hop")

	req.RemoteAddr = "10.0.0.1:5000"

	if err := stack.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Hop-By-Hop"), ""; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Hop-By-Hop", got, want)
	}
	if got, want := req.Header.Get("X-Forwarded-For"), "10.0.0.1"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "X-Forwarded-For", got, want)
	}
	if got, want := req.Header.Get("Via"), "1.1 martian"; !strings.HasPrefix(got, want) {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Via", got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)

	// Hop-by-hop header to be removed.
	res.Header.Set("Hop-By-Hop", "true")
	res.Header.Set("Connection", "Hop-By-Hop")

	if err := stack.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Hop-By-Hop"), ""; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Hop-By-Hop", got, want)
	}
}
