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
	"encoding/json"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/filter"
	"github.com/google/martian/v3/parse"
)

var noop = martian.Noop("header.Filter")

// Filter filters requests and responses based on header name and value.
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

func init() {
	parse.Register("header.Filter", filterFromJSON)
}

// NewFilter builds a new header filter.
func NewFilter(name, value string) *Filter {
	m := NewMatcher(http.CanonicalHeaderKey(name), value)
	f := filter.New()
	f.SetRequestCondition(m)
	f.SetResponseCondition(m)
	return &Filter{f}
}

// filterFromJSON builds a header.Filter from JSON.
//
// Example JSON:
// {
//   "scope": ["request", "result"],
//   "name": "Martian-Testing",
//   "value": "true",
//   "modifier": { ... },
//   "else": { ... }
// }
func filterFromJSON(b []byte) (*parse.Result, error) {
	msg := &filterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	filter := NewFilter(msg.Name, msg.Value)

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
