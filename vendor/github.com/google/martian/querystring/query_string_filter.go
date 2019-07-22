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

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/filter"
	"github.com/google/martian/v3/parse"
)

var noop = martian.Noop("querystring.Filter")

func init() {
	parse.Register("querystring.Filter", filterFromJSON)
}

// Filter runs modifiers iff the request query parameter for name matches value.
type Filter struct {
	*filter.Filter
}

type filterJSON struct {
	Name         string               `json:"name"`
	Value        string               `json:"value"`
	Modifier     json.RawMessage      `json:"modifier"`
	ElseModifier json.RawMessage      `json:"else"`
	Scope        []parse.ModifierType `json:"scope"`
}

// NewFilter builds a querystring.Filter that filters on name and optionally
// value.
func NewFilter(name, value string) *Filter {
	m := NewMatcher(name, value)
	f := filter.New()
	f.SetRequestCondition(m)
	f.SetResponseCondition(m)
	return &Filter{f}
}

// filterFromJSON takes a JSON message and returns a querystring.Filter.
//
// Example JSON:
// {
//   "name": "param",
//   "value": "example",
//   "scope": ["request", "response"],
//   "modifier": { ... }
// }
func filterFromJSON(b []byte) (*parse.Result, error) {
	msg := &filterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	f := NewFilter(msg.Name, msg.Value)

	r, err := parse.FromJSON(msg.Modifier)
	if err != nil {
		return nil, err
	}

	f.RequestWhenTrue(r.RequestModifier())
	f.ResponseWhenTrue(r.ResponseModifier())

	if len(msg.ElseModifier) > 0 {
		em, err := parse.FromJSON(msg.ElseModifier)
		if err != nil {
			return nil, err
		}

		if em != nil {
			f.RequestWhenFalse(em.RequestModifier())
			f.ResponseWhenFalse(em.ResponseModifier())
		}
	}

	return parse.NewResult(f, msg.Scope)
}
