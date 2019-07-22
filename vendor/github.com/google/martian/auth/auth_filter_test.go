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

package auth

import (
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
)

func TestFilter(t *testing.T) {
	f := NewFilter()
	if f.RequestModifier("id") != nil {
		t.Fatalf("f.RequestModifier(%q): got reqmod, want nil", "id")
	}
	if f.ResponseModifier("id") != nil {
		t.Fatalf("f.ResponseModifier(%q): got resmod, want nil", "id")
	}

	tm := martiantest.NewModifier()
	f.SetRequestModifier("id", tm)
	f.SetResponseModifier("id", tm)

	if f.RequestModifier("id") != tm {
		t.Errorf("f.RequestModifier(%q): got nil, want martiantest.Modifier", "id")
	}
	if f.ResponseModifier("id") != tm {
		t.Errorf("f.ResponseModifier(%q): got nil, want martiantest.Modifier", "id")
	}
}

func TestModifyRequest(t *testing.T) {
	f := NewFilter()

	tm := martiantest.NewModifier()
	f.SetRequestModifier("id", tm)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	// No ID, auth required.
	f.SetAuthRequired(true)

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := f.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	actx := FromContext(ctx)
	if actx.Error() == nil {
		t.Error("actx.Error(): got nil, want error")
	}
	if tm.RequestModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}
	tm.Reset()

	// No ID, auth not required.
	f.SetAuthRequired(false)
	actx.SetError(nil)

	if err := f.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if actx.Error() != nil {
		t.Errorf("actx.Error(): got %v, want no error", err)
	}
	if tm.RequestModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}

	// Valid ID.
	actx.SetError(nil)
	actx.SetID("id")

	if err := f.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if actx.Error() != nil {
		t.Errorf("actx.Error(): got %v, want no error", actx.Error())
	}
	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}
}

func TestModifyResponse(t *testing.T) {
	f := NewFilter()

	tm := martiantest.NewModifier()
	f.SetResponseModifier("id", tm)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	res := proxyutil.NewResponse(200, nil, req)

	// No ID, auth required.
	f.SetAuthRequired(true)

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := f.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	actx := FromContext(ctx)
	if actx.Error() == nil {
		t.Error("actx.Error(): got nil, want error")
	}
	if tm.ResponseModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}

	// No ID, no auth required.
	f.SetAuthRequired(false)
	actx.SetError(nil)

	if err := f.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got true, want false")
	}

	// Valid ID.
	actx.SetID("id")

	if err := f.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
}
