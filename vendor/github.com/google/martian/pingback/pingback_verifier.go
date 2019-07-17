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

// Package pingback provides verification that specific URLs have been seen by
// the proxy.
package pingback

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

const (
	errFormat = "request(%s): pingback never occurred"
)

func init() {
	parse.Register("pingback.Verifier", verifierFromJSON)
}

// Verifier verifies that the specific URL has been seen.
type Verifier struct {
	url *url.URL
	err error
}

type verifierJSON struct {
	Scheme string               `json:"scheme"`
	Host   string               `json:"host"`
	Path   string               `json:"path"`
	Query  string               `json:"query"`
	Scope  []parse.ModifierType `json:"scope"`
}

// NewVerifier returns a new pingback verifier.
func NewVerifier(url *url.URL) verify.RequestVerifier {
	return &Verifier{
		url: url,
		err: fmt.Errorf(errFormat, url.String()),
	}
}

// ModifyRequest verifies that the request URL matches all parts of url.
//
// If the value in url is non-empty, it must be an exact match. If the URL
// matches the pingback, it is recorded by setting the error to nil. The error
// will continue to be nil until the verifier has been reset, regardless of
// subsequent requests matching.
func (v *Verifier) ModifyRequest(req *http.Request) error {
	// skip requests to API
	ctx := martian.NewContext(req)
	if ctx.IsAPIRequest() {
		return nil
	}

	u := req.URL

	switch {
	case v.url.Scheme != "" && v.url.Scheme != u.Scheme:
	case v.url.Host != "" && v.url.Host != u.Host:
	case v.url.Path != "" && v.url.Path != u.Path:
	case v.url.RawQuery != "" && v.url.RawQuery != u.RawQuery:
	default:
		v.err = nil
	}

	return nil
}

// VerifyRequests returns an error if pingback never occurred.
func (v *Verifier) VerifyRequests() error {
	return v.err
}

// ResetRequestVerifications clears the failed request verification.
func (v *Verifier) ResetRequestVerifications() {
	v.err = fmt.Errorf(errFormat, v.url.String())
}

// verifierFromJSON builds a pingback.Verifier from JSON.
//
// Example JSON:
// {
//   "pingback.Verifier": {
//     "scope": ["request"],
//     "scheme": "https",
//     "host": "www.google.com",
//     "path": "/proxy",
//     "query": "testing=true"
//   }
// }
func verifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &verifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	v := NewVerifier(&url.URL{
		Scheme:   msg.Scheme,
		Host:     msg.Host,
		Path:     msg.Path,
		RawQuery: msg.Query,
	})

	return parse.NewResult(v, msg.Scope)
}
