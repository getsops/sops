// Copyright 2016 Google Inc. All rights reserved.
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

package servemux

import "net/http"

// Matcher is a conditional evaluator of request urls against patterns registered
// in mux.
type Matcher struct {
	mux *http.ServeMux
}

// NewMatcher builds a new servemux.Matcher.
func NewMatcher(mux *http.ServeMux) *Matcher {
	return &Matcher{
		mux: mux,
	}
}

// MatchRequest returns true if the request URL matches any pattern in mux. If no
// pattern is matched, false is returned.
func (m *Matcher) MatchRequest(req *http.Request) bool {
	if _, pattern := m.mux.Handler(req); pattern != "" {
		return true
	}

	return false
}

// MatchResponse returns true if the request URL associated with the response matches
// any pattern in mux. If pattern is matched, false is returned.
func (m *Matcher) MatchResponse(res *http.Response) bool {
	if _, pattern := m.mux.Handler(res.Request); pattern != "" {
		return true
	}

	return false
}
