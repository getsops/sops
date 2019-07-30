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

// Package header provides utilities for modifying, filtering, and
// verifying headers in martian.Proxy.
package header

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
	"github.com/google/martian/v3/verify"
)

const (
	headerErrFormat = "%s(%s) header verify failure: got no header, want %s header"
	valueErrFormat  = "%s(%s) header verify failure: got %s with value %s, want value %s"
)

type verifier struct {
	name, value string
	reqerr      *martian.MultiError
	reserr      *martian.MultiError
}

type verifierJSON struct {
	Name  string               `json:"name"`
	Value string               `json:"value"`
	Scope []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("header.Verifier", verifierFromJSON)
}

// NewVerifier creates a new header verifier for the given name and value.
func NewVerifier(name, value string) verify.RequestResponseVerifier {
	return &verifier{
		name:   name,
		value:  value,
		reqerr: martian.NewMultiError(),
		reserr: martian.NewMultiError(),
	}
}

// ModifyRequest verifies that the header for name is present in all modified
// requests. If value is non-empty the value must be present in at least one
// header for name. An error will be added to the contained *MultiError for
// every unmatched request.
func (v *verifier) ModifyRequest(req *http.Request) error {
	h := proxyutil.RequestHeader(req)

	vs, ok := h.All(v.name)
	if !ok {
		v.reqerr.Add(fmt.Errorf(headerErrFormat, "request", req.URL, v.name))
		return nil
	}

	for _, value := range vs {
		switch v.value {
		case "", value:
			return nil
		}
	}

	v.reqerr.Add(fmt.Errorf(valueErrFormat, "request", req.URL, v.name,
		strings.Join(vs, ", "), v.value))

	return nil
}

// ModifyResponse verifies that the header for name is present in all modified
// responses. If value is non-empty the value must be present in at least one
// header for name. An error will be added to the contained *MultiError for
// every unmatched response.
func (v *verifier) ModifyResponse(res *http.Response) error {
	h := proxyutil.ResponseHeader(res)

	vs, ok := h.All(v.name)
	if !ok {
		v.reserr.Add(fmt.Errorf(headerErrFormat, "response", res.Request.URL, v.name))
		return nil
	}

	for _, value := range vs {
		switch v.value {
		case "", value:
			return nil
		}
	}

	v.reserr.Add(fmt.Errorf(valueErrFormat, "response", res.Request.URL, v.name,
		strings.Join(vs, ", "), v.value))

	return nil
}

// VerifyRequests returns an error if verification for any request failed.
// If an error is returned it will be of type *martian.MultiError.
func (v *verifier) VerifyRequests() error {
	if v.reqerr.Empty() {
		return nil
	}

	return v.reqerr
}

// VerifyResponses returns an error if verification for any request failed.
// If an error is returned it will be of type *martian.MultiError.
func (v *verifier) VerifyResponses() error {
	if v.reserr.Empty() {
		return nil
	}

	return v.reserr
}

// ResetRequestVerifications clears all failed request verifications.
func (v *verifier) ResetRequestVerifications() {
	v.reqerr = martian.NewMultiError()
}

// ResetResponseVerifications clears all failed response verifications.
func (v *verifier) ResetResponseVerifications() {
	v.reserr = martian.NewMultiError()
}

// verifierFromJSON builds a header.Verifier from JSON.
//
// Example JSON:
// {
//   "name": "header.Verifier",
//   "scope": ["request", "result"],
//   "modifier": {
//     "name": "Martian-Testing",
//     "value": "true"
//   }
// }
func verifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &verifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewVerifier(msg.Name, msg.Value), msg.Scope)
}
