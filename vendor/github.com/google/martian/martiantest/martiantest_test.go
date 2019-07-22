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

package martiantest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/martian/v3/proxyutil"
)

func TestModifier(t *testing.T) {
	var reqrun bool
	var resrun bool

	moderr := errors.New("modifier error")
	tm := NewModifier()
	tm.RequestError(moderr)
	tm.RequestFunc(func(*http.Request) {
		reqrun = true
	})

	tm.ResponseError(moderr)
	tm.ResponseFunc(func(*http.Response) {
		resrun = true
	})

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := tm.ModifyRequest(req); err != moderr {
		t.Fatalf("tm.ModifyRequest(): got %v, want %v", err, moderr)
	}
	if !tm.RequestModified() {
		t.Errorf("tm.RequestModified(): got false, want true")
	}
	if tm.RequestCount() != 1 {
		t.Errorf("tm.RequestCount(): got %d, want %d", tm.RequestCount(), 1)
	}
	if !reqrun {
		t.Error("reqrun: got false, want true")
	}

	res := proxyutil.NewResponse(200, nil, req)

	if err := tm.ModifyResponse(res); err != moderr {
		t.Fatalf("tm.ModifyResponse(): got %v, want %v", err, moderr)
	}
	if !tm.ResponseModified() {
		t.Errorf("tm.ResponseModified(): got false, want true")
	}
	if tm.ResponseCount() != 1 {
		t.Errorf("tm.ResponseCount(): got %d, want %d", tm.ResponseCount(), 1)
	}
	if !resrun {
		t.Error("resrun: got false, want true")
	}

	tm.Reset()

	if tm.RequestModified() {
		t.Error("tm.RequestModified(): got true, want false")
	}
	if tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got true, want false")
	}
}
