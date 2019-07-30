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

package filter

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"
)

func TestRequestWhenTrueCondition(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.RequestEvaluatesTo(true)
	filter.SetRequestCondition(tmc)

	tmod := martiantest.NewModifier()
	filter.RequestWhenTrue(tmod)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := filter.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := tmod.RequestModified(), true; got != want {
		t.Errorf("tmod.RequestModified(): got %t, want %t", got, want)
	}
}

func TestRequestWithoutSettingCondition(t *testing.T) {
	filter := New()

	// neglect to set a matcher

	tmod := martiantest.NewModifier()
	filter.RequestWhenFalse(tmod)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := filter.ModifyRequest(req); err == nil {
		t.Fatalf("ModifyRequest(): got no error, want error")
	}
}

func TestRequestWhenFalseCondition(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.RequestEvaluatesTo(false)
	filter.SetRequestCondition(tmc)

	tmod := martiantest.NewModifier()
	filter.RequestWhenFalse(tmod)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	if err := filter.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := tmod.RequestModified(), true; got != want {
		t.Errorf("tmod.RequestModified(): got %t, want %t", got, want)
	}
}
func TestResponseWithoutSettingCondition(t *testing.T) {
	filter := New()

	// neglect to set a matcher

	tmod := martiantest.NewModifier()
	filter.ResponseWhenFalse(tmod)

	res := proxyutil.NewResponse(200, nil, nil)

	if err := filter.ModifyResponse(res); err == nil {
		t.Fatalf("ModifyResponse(): got no error, want error")
	}
}

func TestResponseWhenTrueCondition(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.ResponseEvaluatesTo(true)
	filter.SetResponseCondition(tmc)

	tmod := martiantest.NewModifier()
	filter.ResponseWhenTrue(tmod)

	res := proxyutil.NewResponse(200, nil, nil)

	if err := filter.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := tmod.ResponseModified(), true; got != want {
		t.Errorf("tmod.ResponseModified(): got %t, want %t", got, want)
	}
}

func TestResponseWhenFalseCondition(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.ResponseEvaluatesTo(false)
	filter.SetResponseCondition(tmc)

	tmod := martiantest.NewModifier()
	filter.ResponseWhenFalse(tmod)

	res := proxyutil.NewResponse(200, nil, nil)

	if err := filter.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := tmod.ResponseModified(), true; got != want {
		t.Errorf("tmod.ResponseModified(): got %t, want %t", got, want)
	}
}

func TestResetVerifications(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.ResponseEvaluatesTo(true)
	filter.SetResponseCondition(tmc)

	tv := &verify.TestVerifier{
		ResponseError: errors.New("verify response failure"),
	}
	filter.ResponseWhenTrue(tv)

	tv = &verify.TestVerifier{
		RequestError: errors.New("verify request failure"),
	}
	filter.RequestWhenTrue(tv)

	if err := filter.VerifyRequests(); err == nil {
		t.Fatal("VerifyRequests(): got nil, want error")
	}
	if err := filter.VerifyResponses(); err == nil {
		t.Fatal("VerifyResponses(): got nil, want error")
	}

	filter.ResetRequestVerifications()
	filter.ResetResponseVerifications()

	if err := filter.VerifyResponses(); err != nil {
		t.Errorf("VerifyResponses(): got %v, want no error", err)
	}

	if err := filter.VerifyRequests(); err != nil {
		t.Errorf("VerifyRequests(): got %v, want no error", err)
	}
}

func TestPassThroughVerifyRequests(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.RequestEvaluatesTo(true)
	filter.SetRequestCondition(tmc)

	if err := filter.VerifyRequests(); err != nil {
		t.Fatalf("VerifyRequest(): got %v, want no error", err)
	}

	tv := &verify.TestVerifier{
		RequestError: errors.New("verify request failure"),
	}

	filter.RequestWhenTrue(tv)

	if got, want := filter.VerifyRequests().Error(), "verify request failure"; got != want {
		t.Fatalf("VerifyRequests(): got %s, want %s", got, want)
	}
}

func TestPassThroughVerifyResponses(t *testing.T) {
	filter := New()

	tmc := martiantest.NewMatcher()
	tmc.ResponseEvaluatesTo(true)
	filter.SetResponseCondition(tmc)

	if err := filter.VerifyResponses(); err != nil {
		t.Fatalf("VerifyResponses(): got %v, want no error", err)
	}

	tv := &verify.TestVerifier{
		ResponseError: errors.New("verify response failure"),
	}

	filter.ResponseWhenTrue(tv)

	if got, want := filter.VerifyResponses().Error(), "verify response failure"; got != want {
		t.Fatalf("VerifyResponses(): got %s, want %s", got, want)
	}
}
