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
	"regexp"

	"github.com/google/martian/v3/filter"
	"github.com/google/martian/v3/parse"
)

func init() {
	parse.Register("url.RegexFilter", regexFilterFromJSON)
}

// URLRegexFilter runs Modifier if the request URL matches the regex, and runs ElseModifier if not.
// This is not to be confused with url.Filter that does string matching on URL segments.
type URLRegexFilter struct {
	*filter.Filter
}

type regexFilterJSON struct {
	Regex        string               `json:"regex"`
	Modifier     json.RawMessage      `json:"modifier"`
	ElseModifier json.RawMessage      `json:"else"`
	Scope        []parse.ModifierType `json:"scope"`
}

// NewRegexFilter constructs a filter that matches on regular expressions.
func NewRegexFilter(r *regexp.Regexp) *URLRegexFilter {
	filter := filter.New()
	matcher := NewRegexMatcher(r)
	filter.SetRequestCondition(matcher)
	filter.SetResponseCondition(matcher)
	return &URLRegexFilter{filter}
}

// regexFilterFromJSON takes a JSON message as a byte slice and returns a
// parse.Result that contains a URLRegexFilter and a scope. The regex syntax is RE2
// as described at https://golang.org/s/re2syntax.
//
// Example JSON configuration message:
// {
//   "scope": ["request", "response"],
//   "regex": ".*www.example.com.*"
// }
func regexFilterFromJSON(b []byte) (*parse.Result, error) {
	msg := &regexFilterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	matcher, err := regexp.Compile(msg.Regex)
	if err != nil {
		return nil, err
	}

	filter := NewRegexFilter(matcher)

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
