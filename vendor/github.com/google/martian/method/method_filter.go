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

package method

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/filter"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
)

var noop = martian.Noop("method.Filter")

func init() {
	parse.Register("method.Filter", filterFromJSON)
}

// Filter runs modifier iff the request method matches the specified method.
type Filter struct {
	*filter.Filter
}

type filterJSON struct {
	Method       string               `json:"method"`
	Modifier     json.RawMessage      `json:"modifier"`
	ElseModifier json.RawMessage      `json:"else"`
	Scope        []parse.ModifierType `json:"scope"`
}

func filterFromJSON(b []byte) (*parse.Result, error) {
	msg := &filterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	filter := NewFilter(msg.Method)

	m, err := parse.FromJSON(msg.Modifier)
	if err != nil {
		return nil, err
	}

	filter.RequestWhenTrue(m.RequestModifier())
	filter.ResponseWhenTrue(m.ResponseModifier())

	if len(msg.ElseModifier) > 0 {
		em, err := parse.FromJSON(msg.ElseModifier)
		if err != nil {
			return nil, err
		}

		if em != nil {
			filter.RequestWhenFalse(em.RequestModifier())
			filter.ResponseWhenFalse(em.ResponseModifier())
		}
	}

	return parse.NewResult(filter, msg.Scope)
}

// NewFilter constructs a filter that applies the modifer when the
// request method matches meth.
func NewFilter(meth string) *Filter {
	log.Debugf("method.NewFilter(%q)", meth)
	m := NewMatcher(meth)
	f := filter.New()
	f.SetRequestCondition(m)
	f.SetResponseCondition(m)
	return &Filter{f}
}

// Matcher is a conditional evaluator of request methods to be used in
// filters that take conditionals.
type Matcher struct {
	method string
}

// NewMatcher builds a new method matcher.
func NewMatcher(method string) *Matcher {
	return &Matcher{
		method: method,
	}
}

// MatchRequest retuns true if m.method matches the request method.
func (m *Matcher) MatchRequest(req *http.Request) bool {
	matched := m.matches(req.Method)
	if matched {
		log.Debugf("method.MatchRequest: matched %s request: %s", req.Method, req.URL)
	}
	return matched
}

// MatchResponse retuns true if m.method matches res.Request.Method.
func (m *Matcher) MatchResponse(res *http.Response) bool {
	matched := m.matches(res.Request.Method)
	if matched {
		log.Debugf("method.MatchResponse: matched %s request: %s", res.Request.Method, res.Request.URL)
	}
	return matched
}

func (m *Matcher) matches(method string) bool {
	return strings.ToUpper(method) == strings.ToUpper(m.method)
}
