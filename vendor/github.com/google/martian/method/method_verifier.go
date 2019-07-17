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

// Package method provides utilities for working with request methods.
package method

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

type verifier struct {
	method string
	err    *martian.MultiError
}

type verifierJSON struct {
	Method string               `json:"method"`
	Scope  []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("method.Verifier", verifierFromJSON)
}

// NewVerifier returns a new method verifier.
func NewVerifier(method string) (verify.RequestVerifier, error) {
	if method == "" {
		return nil, fmt.Errorf("%s is not a valid HTTP method", method)
	}
	return &verifier{
		method: method,
		err:    martian.NewMultiError(),
	}, nil
}

// ModifyRequest verifies that the request's method matches the given method
// in all modified requests. An error will be added to the contained *MultiError
// if a method is unmatched.
func (v *verifier) ModifyRequest(req *http.Request) error {
	m := req.Method

	if v.method != "" && v.method != m {
		err := fmt.Errorf("request(%v) method verification error: got %v, want %v", req.URL,
			v.method, m)
		v.err.Add(err)
	}

	return nil
}

// VerifyRequests returns an error if verification for any request failed.
// If an error is returned it will be of type *martian.MultiError.
func (v *verifier) VerifyRequests() error {
	if v.err.Empty() {
		return nil
	}

	return v.err
}

// ResetRequestVerifications clears all failed request verifications.
func (v *verifier) ResetRequestVerifications() {
	v.err = martian.NewMultiError()
}

// verifierFromJSON builds a method.Verifier from JSON.
//
// Example JSON:
// {
//   "method.Verifier": {
//     "scope": ["request"],
//     "method": "POST"
//   }
// }
func verifierFromJSON(b []byte) (*parse.Result, error) {

	msg := &verifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}
	v, err := NewVerifier(msg.Method)
	if err != nil {
		return nil, err
	}
	return parse.NewResult(v, msg.Scope)
}
