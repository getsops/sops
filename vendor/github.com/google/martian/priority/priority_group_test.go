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

package priority

import (
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"

	// Import to register header.Modifier with JSON parser.
	_ "github.com/google/martian/v3/header"
)

func TestPriorityGroupModifyRequest(t *testing.T) {
	var order []string

	pg := NewGroup()

	tm50 := martiantest.NewModifier()
	tm50.RequestFunc(func(*http.Request) {
		order = append(order, "tm50")
	})
	pg.AddRequestModifier(tm50, 50)

	tm100a := martiantest.NewModifier()
	tm100a.RequestFunc(func(*http.Request) {
		order = append(order, "tm100a")
	})
	pg.AddRequestModifier(tm100a, 100)

	tm100b := martiantest.NewModifier()
	tm100b.RequestFunc(func(*http.Request) {
		order = append(order, "tm100b")
	})
	pg.AddRequestModifier(tm100b, 100)

	tm75 := martiantest.NewModifier()
	tm75.RequestFunc(func(*http.Request) {
		order = append(order, "tm75")
	})

	if err := pg.RemoveRequestModifier(tm75); err != ErrModifierNotFound {
		t.Fatalf("RemoveRequestModifier(): got %v, want ErrModifierNotFound", err)
	}

	pg.AddRequestModifier(tm75, 100)

	if err := pg.RemoveRequestModifier(tm75); err != nil {
		t.Fatalf("RemoveRequestModifier(): got %v, want no error", err)
	}

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := pg.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := order, []string{"tm100b", "tm100a", "tm50"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("reflect.DeepEqual(%v, %v): got false, want true", got, want)
	}
}

func TestPriorityGroupModifyRequestHaltsOnError(t *testing.T) {
	pg := NewGroup()

	reqerr := errors.New("request error")
	tm := martiantest.NewModifier()
	tm.RequestError(reqerr)

	pg.AddRequestModifier(tm, 100)

	tm2 := martiantest.NewModifier()
	pg.AddRequestModifier(tm2, 75)

	req, err := http.NewRequest("GET", "http://example.com/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := pg.ModifyRequest(req); err != reqerr {
		t.Fatalf("ModifyRequest(): got %v, want %v", err, reqerr)
	}

	if tm2.RequestModified() {
		t.Error("tm2.RequestModified(): got true, want false")
	}
}

func TestPriorityGroupModifyResponse(t *testing.T) {
	var order []string

	pg := NewGroup()

	tm50 := martiantest.NewModifier()
	tm50.ResponseFunc(func(*http.Response) {
		order = append(order, "tm50")
	})
	pg.AddResponseModifier(tm50, 50)

	tm100a := martiantest.NewModifier()
	tm100a.ResponseFunc(func(*http.Response) {
		order = append(order, "tm100a")
	})
	pg.AddResponseModifier(tm100a, 100)

	tm100b := martiantest.NewModifier()
	tm100b.ResponseFunc(func(*http.Response) {
		order = append(order, "tm100b")
	})
	pg.AddResponseModifier(tm100b, 100)

	tm75 := martiantest.NewModifier()
	tm75.ResponseFunc(func(*http.Response) {
		order = append(order, "tm75")
	})

	if err := pg.RemoveResponseModifier(tm75); err != ErrModifierNotFound {
		t.Fatalf("RemoveResponseModifier(): got %v, want ErrModifierNotFound", err)
	}

	pg.AddResponseModifier(tm75, 100)

	if err := pg.RemoveResponseModifier(tm75); err != nil {
		t.Fatalf("RemoveResponseModifier(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, nil)
	if err := pg.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := order, []string{"tm100b", "tm100a", "tm50"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("reflect.DeepEqual(%v, %v): got false, want true", got, want)
	}
}

func TestPriorityGroupModifyResponseHaltsOnError(t *testing.T) {
	pg := NewGroup()

	reserr := errors.New("response error")
	tm := martiantest.NewModifier()
	tm.ResponseError(reserr)

	pg.AddResponseModifier(tm, 100)

	tm2 := martiantest.NewModifier()
	pg.AddResponseModifier(tm2, 75)

	res := proxyutil.NewResponse(200, nil, nil)
	if err := pg.ModifyResponse(res); err != reserr {
		t.Fatalf("ModifyRequest(): got %v, want %v", err, reserr)
	}

	if tm2.ResponseModified() {
		t.Error("tm2.ResponseModified(): got true, want false")
	}
}

func TestGroupFromJSON(t *testing.T) {
	msg := []byte(`{
    "priority.Group": {
    "scope": ["request", "response"],
    "modifiers": [
      {
        "priority": 100,
        "modifier": {
          "header.Modifier": {
            "scope": ["request", "response"],
            "name": "X-Testing",
            "value": "true"
          }
        }
      },
      {
        "priority": 0,
        "modifier": {
          "header.Modifier": {
            "scope": ["request", "response"],
            "name": "Y-Testing",
            "value": "true"
          }
        }
      }
    ]
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
	if got, want := req.Header.Get("X-Testing"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "X-Testing", got, want)
	}
	if got, want := req.Header.Get("Y-Testing"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Y-Testing", got, want)
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if got, want := res.Header.Get("X-Testing"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "X-Testing", got, want)
	}
	if got, want := res.Header.Get("Y-Testing"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Y-Testing", got, want)
	}
}
