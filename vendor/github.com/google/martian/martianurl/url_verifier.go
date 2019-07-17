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

package martianurl

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

const (
	errFormat     = "request(%s) url verify failure:\n%s"
	errPartFormat = "\t%s: got %q, want %q"
)

func init() {
	parse.Register("url.Verifier", verifierFromJSON)
}

// Verifier verifies the structure of URLs.
type Verifier struct {
	url *url.URL
	err *martian.MultiError
}

type verifierJSON struct {
	Scheme string               `json:"scheme"`
	Host   string               `json:"host"`
	Path   string               `json:"path"`
	Query  string               `json:"query"`
	Scope  []parse.ModifierType `json:"scope"`
}

// NewVerifier returns a new URL verifier.
func NewVerifier(url *url.URL) verify.RequestVerifier {
	return &Verifier{
		url: url,
		err: martian.NewMultiError(),
	}
}

// ModifyRequest verifies that the request URL matches all parts of url. If the
// value in url is non-empty it must be an exact match.
func (v *Verifier) ModifyRequest(req *http.Request) error {
	// skip requests to API
	ctx := martian.NewContext(req)
	if ctx.IsAPIRequest() {
		return nil
	}

	var failures []string

	u := req.URL

	if v.url.Scheme != "" && v.url.Scheme != u.Scheme {
		f := fmt.Sprintf(errPartFormat, "Scheme", u.Scheme, v.url.Scheme)
		failures = append(failures, f)
	}
	if v.url.Host != "" && !MatchHost(u.Host, v.url.Host) {
		f := fmt.Sprintf(errPartFormat, "Host", u.Host, v.url.Host)
		failures = append(failures, f)
	}
	if v.url.Path != "" && v.url.Path != u.Path {
		f := fmt.Sprintf(errPartFormat, "Path", u.Path, v.url.Path)
		failures = append(failures, f)
	}
	if v.url.RawQuery != "" && v.url.RawQuery != u.RawQuery {
		f := fmt.Sprintf(errPartFormat, "Query", u.RawQuery, v.url.RawQuery)
		failures = append(failures, f)
	}
	if v.url.Fragment != "" && v.url.Fragment != u.Fragment {
		f := fmt.Sprintf(errPartFormat, "Fragment", u.Fragment, v.url.Fragment)
		failures = append(failures, f)
	}

	if len(failures) > 0 {
		err := fmt.Errorf(errFormat, u, strings.Join(failures, "\n"))
		v.err.Add(err)
	}

	return nil
}

// VerifyRequests returns an error if verification for any request failed.
// If an error is returned it will be of type *martian.MultiError.
func (v *Verifier) VerifyRequests() error {
	if v.err.Empty() {
		return nil
	}

	return v.err
}

// ResetRequestVerifications clears all failed request verifications.
func (v *Verifier) ResetRequestVerifications() {
	v.err = martian.NewMultiError()
}

// verifierFromJSON builds a martianurl.Verifier from JSON.
//
// Example modifier JSON:
// {
//   "martianurl.Verifier": {
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
