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

package skip

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

func TestRoundTrip(t *testing.T) {
	m := NewRoundTrip()
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if ctx.SkippingRoundTrip() {
		t.Fatal("ctx.SkippingRoundTrip(): got true, want false")
	}

	if err := m.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if !ctx.SkippingRoundTrip() {
		t.Fatal("ctx.SkippingRoundTrip(): got false, want true")
	}
}

func TestFromJSON(t *testing.T) {
	msg := []byte(`{
			"skip.RoundTrip": {}
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if _, ok := reqmod.(*RoundTrip); !ok {
		t.Fatal("reqmod.(*RoundTrip): got !ok, want ok")
	}
}
