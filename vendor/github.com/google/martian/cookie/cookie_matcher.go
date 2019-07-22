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

package cookie

import (
	"net/http"

	"github.com/google/martian/v3/log"
)

// Matcher is a conditonal evalutor of request or
// response cookies to be used in structs that take conditions.
type Matcher struct {
	cookie *http.Cookie
}

// NewMatcher builds a cookie matcher.
func NewMatcher(cookie *http.Cookie) *Matcher {
	return &Matcher{
		cookie: cookie,
	}
}

// MatchRequest evaluates a request and returns whether or not
// the request contains a cookie that matches the provided name, path
// and value.
func (m *Matcher) MatchRequest(req *http.Request) bool {
	for _, c := range req.Cookies() {
		if m.match(c) {
			log.Debugf("cookie.MatchRequest: %s, matched: cookie: %s", req.URL, c)
			return true
		}
	}

	return false
}

// MatchResponse evaluates a response and returns whether or not the response
// contains a cookie that matches the provided name and value.
func (m *Matcher) MatchResponse(res *http.Response) bool {
	for _, c := range res.Cookies() {
		if m.match(c) {
			log.Debugf("cookie.MatchResponse: %s, matched: cookie: %s", res.Request.URL, c)
			return true
		}
	}

	return false
}

func (m *Matcher) match(cs *http.Cookie) bool {
	switch {
	case m.cookie.Name != cs.Name:
		return false
	case m.cookie.Value != "" && m.cookie.Value != cs.Value:
		return false
	}

	return true
}
