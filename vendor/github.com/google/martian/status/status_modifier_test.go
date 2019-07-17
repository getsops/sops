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

package status

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestFromJSON(t *testing.T) {
	msg := []byte(`{
    "status.Modifier": {
      "scope": ["response"],
      "statusCode": 400
    }
  }`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}
	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, nil)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.StatusCode, 400; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
}

func TestStatusModifierOnResponse(t *testing.T) {
	for i, status := range []int{
		http.StatusForbidden,
		http.StatusOK,
		http.StatusTemporaryRedirect,
	} {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatalf("NewRequest(): got %v, want no error", err)
		}

		res := proxyutil.NewResponse(200, nil, req)

		mod := NewModifier(status)

		if err := mod.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}

		if got, want := res.StatusCode, status; got != want {
			t.Errorf("%d. res.StatusCode: got %v, want %v", i, got, want)
		}
		if got, want := res.Status, fmt.Sprintf("%d %s", res.StatusCode, http.StatusText(status)); got != want {
			t.Errorf("%d. res.Status: got %q, want %q", i, got, want)
		}
	}
}
