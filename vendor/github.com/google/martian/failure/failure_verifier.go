// Copyright 2017 Google Inc. All rights reserved.
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

// Package failure provides a verifier that always fails, adding a given message
// to the multierror log. This can be used to turn any filter in to a defacto
// verifier, by wrapping it in a filter and thus causing a verifier failure whenever
// a request passes the filter.
package failure

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

func init() {
	parse.Register("failure.Verifier", verifierFromJSON)
}

type verifier struct {
	message string
	merr    *martian.MultiError
}

type verifierJSON struct {
	Message string               `json:"message"`
	Scope   []parse.ModifierType `json:"scope"`
}

// NewVerifier returns a new failing verifier.
func NewVerifier(message string) (verify.RequestVerifier, error) {
	return &verifier{
		message: message,
		merr:    martian.NewMultiError(),
	}, nil
}

// ModifyRequest adds an error message containing the message field in the verifier to the verifier errors.
// This means that any time a request hits the verifier it's treated as an error.
func (v *verifier) ModifyRequest(req *http.Request) error {
	err := fmt.Errorf("request(%v) verification error: %s", req.URL, v.message)
	v.merr.Add(err)
	return nil
}

// VerifyRequests returns an error if any requests have hit the verifier.
// If an error is returned it will be of type *martian.MultiError.
func (v *verifier) VerifyRequests() error {
	if v.merr.Empty() {
		return nil
	}

	return v.merr
}

// ResetRequestVerifications clears all failed request verifications.
func (v *verifier) ResetRequestVerifications() {
	v.merr = martian.NewMultiError()
}

// verifierFromJSON builds a failure.Verifier from JSON
//
// Example JSON:
// {
//   "failure.Verifier": {
//     "scope": ["request", "response"],
//     "message": "Request passed a filter it should not have"
//   }
// }
func verifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &verifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	v, err := NewVerifier(msg.Message)
	if err != nil {
		return nil, err
	}

	return parse.NewResult(v, msg.Scope)
}
