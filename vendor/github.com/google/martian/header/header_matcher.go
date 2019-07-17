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

package header

import (
	"net/http"

	"github.com/google/martian/v3/proxyutil"
)

// Matcher is a conditonal evalutor of request or
// response headers to be used in structs that take conditions.
type Matcher struct {
	name, value string
}

// NewMatcher builds a new header matcher.
func NewMatcher(name, value string) *Matcher {
	return &Matcher{
		name:  name,
		value: value,
	}
}

// MatchRequest evaluates a request and returns whether or not
// the request contains a header that matches the provided name
// and value.
func (m *Matcher) MatchRequest(req *http.Request) bool {
	h := proxyutil.RequestHeader(req)

	vs, ok := h.All(m.name)
	if !ok {
		return false
	}

	for _, v := range vs {
		if v == m.value {
			return true
		}
	}

	return false
}

// MatchResponse evaluates a response and returns whether or not
// the response contains a header that matches the provided name
// and value.
func (m *Matcher) MatchResponse(res *http.Response) bool {
	h := proxyutil.ResponseHeader(res)

	vs, ok := h.All(m.name)
	if !ok {
		return false
	}

	for _, v := range vs {
		if v == m.value {
			return true
		}
	}

	return false
}
