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

	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestBlacklistModifierOnRequest(t *testing.T) {
	mod := NewBlacklistModifier("X-Testing")

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	req.Header.Set("X-Testing", "value")
	req.Header.Set("Y-Testing", "value")

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if _, ok := req.Header["X-Testing"]; ok {
		t.Errorf("req.Header[%q]: got true, want false", "X-Testing")
	}

	if _, ok := req.Header["Y-Testing"]; !ok {
		t.Errorf("req.Header[%q]: got false, want true", "Y-Testing")
	}
}

func TestBlacklistModifierOnResponse(t *testing.T) {
	mod := NewBlacklistModifier("X-Testing")

	res := proxyutil.NewResponse(200, nil, nil)

	res.Header.Set("X-Testing", "value")
	res.Header.Set("Y-Testing", "value")

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if _, ok := res.Header["X-Testing"]; ok {
		t.Errorf("res.Header[%q]: got true, want false", "X-Testing")
	}

	if _, ok := res.Header["Y-Testing"]; !ok {
		t.Errorf("res.Header[%q]: got false, want true", "Y-Testing")
	}
}

func TestBlacklistModifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "header.Blacklist": {
  		"scope": ["request", "response"],
			"names": ["X-Testing", "Y-Testing"]
		}
	}`)
	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "http://martian.test", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %q, want no error", err)
	}

	req.Header.Set("X-Testing", "value")
	req.Header.Set("Y-Testing", "value")
	req.Header.Set("Z-Testing", "value")

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatalf("reqmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header.Set("X-Testing", "value")
	res.Header.Set("Y-Testing", "value")
	res.Header.Set("Z-Testing", "value")

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	tt := []struct {
		header string
		want   string
	}{
		{
			header: "X-Testing",
			want:   "",
		},
		{
			header: "Y-Testing",
			want:   "",
		},
		{
			header: "Z-Testing",
			want:   "value",
		},
	}

	for i, tc := range tt {
		if err := reqmod.ModifyRequest(req); err != nil {
			t.Fatalf("%d. reqmod.ModifyRequest(): got %v, want no error", i, err)
		}

		if got, want := req.Header.Get(tc.header), tc.want; got != want {
			t.Errorf("%d. req.Header.Get(%q): got %q, want %q", i, tc.header, got, want)
		}

		if err := resmod.ModifyResponse(res); err != nil {
			t.Fatalf("%d. resmod.ModifyResponse(): got %v, want no error", i, err)
		}

		if got, want := res.Header.Get(tc.header), tc.want; got != want {
			t.Errorf("%d. res.Header.Get(%q): got %q, want %q", i, tc.header, got, want)
		}
	}
}
