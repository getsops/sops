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
	"fmt"
	"net/http"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"
)

func TestVerifyRequestsBlankValue(t *testing.T) {
	v := NewVerifier("Martian-Test", "")

	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://www.example.com/%d", i), nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		// Request 1, 3 should fail verification.
		if i%2 == 0 {
			req.Header.Set("Martian-Test", "true")
		}

		if err := v.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}
	}

	merr, ok := v.VerifyRequests().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyRequests(): got no error, want *verify.MultiError")
	}
	if got, want := len(merr.Errors()), 2; got != want {
		t.Fatalf("len(merr.Errors()): got %d, want %d", got, want)
	}

	wants := []string{
		`request(http://www.example.com/1) header verify failure: got no header, want Martian-Test header`,
		`request(http://www.example.com/3) header verify failure: got no header, want Martian-Test header`,
	}
	for i, err := range merr.Errors() {
		if got := err.Error(); got != wants[i] {
			t.Errorf("Errors()[%d]: got %q, want %q", i, got, wants[i])
		}
	}

	v.ResetRequestVerifications()
	if err := v.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestVerifierFromJSON(t *testing.T) {
	msg := []byte(`{
    "header.Verifier": {
      "scope": ["request", "response"],
      "name": "Martian-Test",
      "value": "true"
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
	reqv, ok := reqmod.(verify.RequestVerifier)
	if !ok {
		t.Fatal("reqmod.(verify.RequestVerifier): got !ok, want ok")
	}

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	if err := reqv.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}
	if err := reqv.VerifyRequests(); err == nil {
		t.Error("VerifyRequests(): got nil, want not nil")
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("resmod: got nil, want not nil")
	}
	resv, ok := resmod.(verify.ResponseVerifier)
	if !ok {
		t.Fatal("resmod.(verify.ResponseVerifier): got !ok, want ok")
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resv.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}
	if err := resv.VerifyResponses(); err == nil {
		t.Error("VerifyResponses(): got nil, want not nil")
	}
}

func TestVerifyRequests(t *testing.T) {
	v := NewVerifier("Martian-Test", "testing-even")

	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://www.example.com/%d", i), nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		req.Header.Add("Martian-Test", fmt.Sprintf("test-%d", i))

		// Request 1, 3 should fail verification.
		if i%2 == 0 {
			req.Header.Add("Martian-Test", "testing-even")
		} else {
			req.Header.Add("Martian-Test", "testing-odd")
		}

		if err := v.ModifyRequest(req); err != nil {
			t.Fatalf("%d. ModifyRequest(): got %v, want no error", i, err)
		}
	}

	merr, ok := v.VerifyRequests().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyRequests(): got no error, want *verify.MultiError")
	}
	if got, want := len(merr.Errors()), 2; got != want {
		t.Fatalf("len(merr.Errors()): got %d, want %d", got, want)
	}

	wants := []string{
		`request(http://www.example.com/1) header verify failure: got Martian-Test with value test-1, testing-odd, want value testing-even`,
		`request(http://www.example.com/3) header verify failure: got Martian-Test with value test-3, testing-odd, want value testing-even`,
	}
	for i, err := range merr.Errors() {
		if got := err.Error(); got != wants[i] {
			t.Errorf("Errors()[%d]: got %q, want %q", i, got, wants[i])
		}
	}

	v.ResetRequestVerifications()
	if err := v.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestVerifyResponsesBlankValue(t *testing.T) {
	v := NewVerifier("Martian-Test", "")

	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://www.example.com/%d", i), nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}
		res := proxyutil.NewResponse(200, nil, req)

		// Response 1, 3 should fail verification.
		if i%2 == 0 {
			res.Header.Set("Martian-Test", "true")
		}

		if err := v.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}
	}

	merr, ok := v.VerifyResponses().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyResponses(): got no error, want *verify.MultiError")
	}
	if got, want := len(merr.Errors()), 2; got != want {
		t.Fatalf("len(merr.Errors()): got %d, want %d", got, want)
	}

	wants := []string{
		`response(http://www.example.com/1) header verify failure: got no header, want Martian-Test header`,
		`response(http://www.example.com/3) header verify failure: got no header, want Martian-Test header`,
	}
	for i, err := range merr.Errors() {
		if got := err.Error(); got != wants[i] {
			t.Errorf("Errors()[%d]: got %q, want %q", i, got, wants[i])
		}
	}

	v.ResetResponseVerifications()
	if err := v.VerifyResponses(); err != nil {
		t.Errorf("VerifyResponses(): got %v, want no error", err)
	}
}

func TestVerifyResponses(t *testing.T) {
	v := NewVerifier("Martian-Test", "testing-even")

	for i := 0; i < 4; i++ {
		req, err := http.NewRequest("GET", fmt.Sprintf("http://www.example.com/%d", i), nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}
		res := proxyutil.NewResponse(200, nil, req)

		res.Header.Add("Martian-Test", fmt.Sprintf("test-%d", i))

		// Response 1, 3 should fail verification.
		if i%2 == 0 {
			res.Header.Add("Martian-Test", "testing-even")
		} else {
			res.Header.Add("Martian-Test", "testing-odd")
		}

		if err := v.ModifyResponse(res); err != nil {
			t.Fatalf("%d. ModifyResponse(): got %v, want no error", i, err)
		}
	}

	merr, ok := v.VerifyResponses().(*martian.MultiError)
	if !ok {
		t.Fatal("VerifyResponses(): got no error, want *verify.MultiError")
	}
	if got, want := len(merr.Errors()), 2; got != want {
		t.Fatalf("len(merr.Errors()): got %d, want %d", got, want)
	}

	wants := []string{
		`response(http://www.example.com/1) header verify failure: got Martian-Test with value test-1, testing-odd, want value testing-even`,
		`response(http://www.example.com/3) header verify failure: got Martian-Test with value test-3, testing-odd, want value testing-even`,
	}
	for i, err := range merr.Errors() {
		if got := err.Error(); got != wants[i] {
			t.Errorf("Errors()[%d]: got %q, want %q", i, got, wants[i])
		}
	}

	v.ResetResponseVerifications()
	if err := v.VerifyResponses(); err != nil {
		t.Errorf("VerifyResponses(): got %v, want no error", err)
	}
}
