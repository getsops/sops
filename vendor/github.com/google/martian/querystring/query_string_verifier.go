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

package querystring

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

func init() {
	parse.Register("querystring.Verifier", verifierFromJSON)
}

type verifier struct {
	key, value string
	err        *martian.MultiError
}

type verifierJSON struct {
	Name  string               `json:"name"`
	Value string               `json:"value"`
	Scope []parse.ModifierType `json:"scope"`
}

// NewVerifier returns a new param verifier.
func NewVerifier(key, value string) (verify.RequestVerifier, error) {
	if key == "" {
		return nil, fmt.Errorf("no key provided to param verifier")
	}
	return &verifier{
		key:   key,
		value: value,
		err:   martian.NewMultiError(),
	}, nil
}

// ModifyRequest verifies that the request's URL params match the given params
// in all modified requests. If no value is provided, the verifier will only
// check if the given key is present. An error will be added to the contained
// *MultiError if the param is unmatched.
func (v *verifier) ModifyRequest(req *http.Request) error {
	// skip requests to API
	ctx := martian.NewContext(req)
	if ctx.IsAPIRequest() {
		return nil
	}

	if err := req.ParseForm(); err != nil {
		err := fmt.Errorf("request(%v) parsing failed; could not parse query parameters", req.URL)
		v.err.Add(err)
		return nil
	}
	vals, ok := req.Form[v.key]
	if !ok {
		err := fmt.Errorf("request(%v) param verification error: key %v not found", req.URL, v.key)
		v.err.Add(err)
		return nil
	}
	if v.value == "" {
		return nil
	}
	for _, val := range vals {
		if v.value == val {
			return nil
		}
	}
	err := fmt.Errorf("request(%v) param verification error: got %v for key %v, want %v", req.URL, strings.Join(vals, ", "), v.key, v.value)
	v.err.Add(err)
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

// verifierFromJSON builds a querystring.Verifier from JSON.
//
// Example JSON:
// {
//   "querystring.Verifier": {
//     "scope": ["request", "response"],
//     "name": "Martian-Testing",
//     "value": "true"
//   }
// }
func verifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &verifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	v, err := NewVerifier(msg.Name, msg.Value)
	if err != nil {
		return nil, err
	}

	return parse.NewResult(v, msg.Scope)
}
