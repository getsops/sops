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
	"strings"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/proxyutil"
)

func TestViaModifier(t *testing.T) {
	m := NewViaModifier("martian")
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	res := proxyutil.NewResponse(200, nil, req)

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("Via"), "1.1 martian"; !strings.HasPrefix(got, want) {
		t.Errorf("req.Header.Get(%q): got %q, want prefixed with %q", "Via", got, want)
	}

	if err := m.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	req.Header.Set("Via", "1.0\talpha\t(martian)")
	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header.Get("Via"), "1.0\talpha\t(martian), 1.1 martian"; !strings.HasPrefix(got, want) {
		t.Errorf("req.Header.Get(%q): got %q, want prefixed with %q", "Via", got, want)
	}

	m.SetBoundary("boundary")
	req.Header.Set("Via", "1.0\talpha\t(martian), 1.1 martian-boundary, 1.1 beta")
	if err := m.ModifyRequest(req); err == nil {
		t.Fatal("ModifyRequest(): got nil, want request loop error")
	}
	if !ctx.SkippingRoundTrip() {
		t.Errorf("ctx.SkippingRoundTrip(): got false, want true")
	}

	if err := m.ModifyResponse(res); err == nil {
		t.Fatal("ModifyResponse(): got nil, want request loop error")
	}
	if got, want := res.StatusCode, 400; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Status, http.StatusText(400); got != want {
		t.Errorf("res.Status: got %q, want %q", got, want)
	}
}
