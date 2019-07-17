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

package querystring

import "net/http"

// Matcher is a conditonal evalutor of query string parameters
// to be used in structs that take conditions.
type Matcher struct {
	name, value string
}

// NewMatcher builds a new querystring matcher
func NewMatcher(name, value string) *Matcher {
	return &Matcher{name: name, value: value}
}

// MatchRequest evaluates a request and returns whether or not
// the request contains a querystring param that matches the provided name
// and value.
func (m *Matcher) MatchRequest(req *http.Request) bool {
	for n, vs := range req.URL.Query() {
		if m.name == n {
			if m.value == "" {
				return true
			}

			for _, v := range vs {
				if m.value == v {
					return true
				}
			}
		}
	}

	return false
}

// MatchResponse evaluates a response and returns whether or not
// the request that resulted in that response contains a querystring param that matches the provided name
// and value.
func (m *Matcher) MatchResponse(res *http.Response) bool {
	return m.MatchRequest(res.Request)
}
