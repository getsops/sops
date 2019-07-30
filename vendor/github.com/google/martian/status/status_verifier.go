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

package status

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

const errFormat = "response(%s) status code verify failure: got %d, want %d"

// Verifier verifies the status codes of all responses.
type Verifier struct {
	statusCode int
	err        *martian.MultiError
}

type verifierJSON struct {
	StatusCode int                  `json:"statusCode"`
	Scope      []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("status.Verifier", verifierFromJSON)
}

// NewVerifier returns a new status.Verifier for statusCode.
func NewVerifier(statusCode int) verify.ResponseVerifier {
	return &Verifier{
		statusCode: statusCode,
		err:        martian.NewMultiError(),
	}
}

// ModifyResponse verifies that the status code for all requests
// matches statusCode.
func (v *Verifier) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	if ctx.IsAPIRequest() {
		return nil
	}

	if res.StatusCode != v.statusCode {
		v.err.Add(fmt.Errorf(errFormat, res.Request.URL, res.StatusCode, v.statusCode))
	}

	return nil
}

// VerifyResponses returns an error if verification for any
// request failed.
// If an error is returned it will be of type *martian.MultiError.
func (v *Verifier) VerifyResponses() error {
	if v.err.Empty() {
		return nil
	}

	return v.err
}

// ResetResponseVerifications clears all failed response verifications.
func (v *Verifier) ResetResponseVerifications() {
	v.err = martian.NewMultiError()
}

// verifierFromJSON builds a status.Verifier from JSON.
//
// Example JSON:
// {
//   "status.Verifier": {
//     "scope": ["response"],
//     "statusCode": 401
//   }
// }
func verifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &verifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewVerifier(msg.StatusCode), msg.Scope)
}
