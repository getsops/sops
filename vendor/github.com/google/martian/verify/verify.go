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

// Package verify provides support for using martian modifiers for request and
// response verifications.
package verify

import (
	"net/http"

	"github.com/google/martian/v3"
)

// RequestVerifier is a RequestModifier that maintains a verification state.
// RequestVerifiers should only return an error from ModifyRequest for errors
// unrelated to the expectation.
type RequestVerifier interface {
	martian.RequestModifier
	VerifyRequests() error
	ResetRequestVerifications()
}

// ResponseVerifier is a ResponseModifier that maintains a verification state.
// ResponseVerifiers should only return an error from ModifyResponse for errors
// unrelated to the expectation.
type ResponseVerifier interface {
	martian.ResponseModifier
	VerifyResponses() error
	ResetResponseVerifications()
}

// RequestResponseVerifier is a RequestVerifier and a ResponseVerifier.
type RequestResponseVerifier interface {
	RequestVerifier
	ResponseVerifier
}

// TestVerifier is a request and response verifier with overridable errors for
// verification.
type TestVerifier struct {
	RequestError  error
	ResponseError error
}

// ModifyRequest is a no-op.
func (tv *TestVerifier) ModifyRequest(*http.Request) error {
	return nil
}

// ModifyResponse is a no-op.
func (tv *TestVerifier) ModifyResponse(*http.Response) error {
	return nil
}

// VerifyRequests returns the set request error.
func (tv *TestVerifier) VerifyRequests() error {
	return tv.RequestError
}

// VerifyResponses returns the set response error.
func (tv *TestVerifier) VerifyResponses() error {
	return tv.ResponseError
}

// ResetRequestVerifications clears out the set request error.
func (tv *TestVerifier) ResetRequestVerifications() {
	tv.RequestError = nil
}

// ResetResponseVerifications clears out the set response error.
func (tv *TestVerifier) ResetResponseVerifications() {
	tv.ResponseError = nil
}
