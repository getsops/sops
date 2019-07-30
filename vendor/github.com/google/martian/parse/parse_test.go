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

package parse

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/martiantest"
)

func TestFromJSON(t *testing.T) {
	msg := []byte(`{
		"first.Modifier": { },
		"second.Modifier": { }
	}`)

	if _, err := FromJSON(msg); err == nil {
		t.Error("FromJson(): got nil, want more than one key error")
	}

	Register("martiantest.Modifier", func(b []byte) (*Result, error) {
		type testJSON struct {
			Scope []ModifierType `json:"scope"`
		}

		msg := &testJSON{}
		if err := json.Unmarshal(b, msg); err != nil {
			return nil, err
		}

		tm := martiantest.NewModifier()
		return NewResult(tm, msg.Scope)
	})

	msg = []byte(`{
		"martiantest.Modifier": { }
	}`)

	r, err := FromJSON(msg)
	if err != nil {
		t.Fatalf("FromJSON(): got %v, want no error", err)
	}

	if _, ok := r.RequestModifier().(*martiantest.Modifier); !ok {
		t.Fatal("r.RequestModifier().(*martiantest.Modifier): got !ok, want ok")
	}

	if _, ok := r.ResponseModifier().(*martiantest.Modifier); !ok {
		t.Fatal("r.ResponseModifier().(*martiantest.Modifier): got !ok, want ok")
	}

	msg = []byte(`{
	  "martiantest.Modifier": {
      "scope": ["request"]
    }
	}`)

	r, err = FromJSON(msg)
	if err != nil {
		t.Fatalf("FromJSON(): got %v, want no error", err)
	}

	if _, ok := r.RequestModifier().(*martiantest.Modifier); !ok {
		t.Fatal("r.RequestModifier().(*martiantest.Modifier): got !ok, want ok")
	}

	resmod := r.ResponseModifier()
	if resmod != nil {
		t.Error("r.ResponseModifier(): got not nil, want nil")
	}

	msg = []byte(`{
	  "martiantest.Modifier": {
      "scope": ["response"]
    }
	}`)

	r, err = FromJSON(msg)
	if err != nil {
		t.Fatalf("FromJSON(): got %v, want no error", err)
	}

	if _, ok := r.ResponseModifier().(*martiantest.Modifier); !ok {
		t.Fatal("r.ResponseModifier().(*martiantest.Modifier): got !ok, want ok")
	}

	reqmod := r.RequestModifier()
	if reqmod != nil {
		t.Error("r.RequestModifier(): got not nil, want nil")
	}
}

func TestNewResultMismatchedScopes(t *testing.T) {
	reqmod := martian.RequestModifierFunc(
		func(*http.Request) error {
			return nil
		})
	resmod := martian.ResponseModifierFunc(
		func(*http.Response) error {
			return nil
		})

	if _, err := NewResult(reqmod, []ModifierType{Response}); err == nil {
		t.Error("NewResult(reqmod, RESPONSE): got nil, want error")
	}

	if _, err := NewResult(resmod, []ModifierType{Request}); err == nil {
		t.Error("NewResult(resmod, REQUEST): got nil, want error")
	}

	if _, err := NewResult(reqmod, []ModifierType{ModifierType("unknown")}); err == nil {
		t.Error("NewResult(resmod, REQUEST): got nil, want error")
	}
}

func TestResultModifierAccessors(t *testing.T) {
	tm := martiantest.NewModifier()

	r := &Result{
		reqmod: tm,
		resmod: nil,
	}
	if reqmod := r.RequestModifier(); reqmod == nil {
		t.Error("r.RequestModifier: got nil, want reqmod")
	}

	if resmod := r.ResponseModifier(); resmod != nil {
		t.Error("r.ResponseModifier: got resmod, want nil")
	}

	r = &Result{
		reqmod: nil,
		resmod: tm,
	}
	if reqmod := r.RequestModifier(); reqmod != nil {
		t.Errorf("r.RequestModifier: got reqmod, want nil")
	}

	if resmod := r.ResponseModifier(); resmod == nil {
		t.Error("r.ResponseModifier: got nil, want resmod")
	}
}

func TestParseUnknownModifierReturnsError(t *testing.T) {
	msg := []byte(`{
	  "unknown.Key": {
      "scope": ["request", "response"]
		}
	}`)

	_, err := FromJSON(msg)

	umerr, ok := err.(ErrUnknownModifier)
	if !ok {
		t.Fatalf("FromJSON(): got %v, want ErrUnknownModifier", err)
	}

	if got, want := umerr.Error(), "parse: unknown modifier: unknown.Key"; got != want {
		t.Errorf("Error(): got %q, want %q", got, want)
	}
}
