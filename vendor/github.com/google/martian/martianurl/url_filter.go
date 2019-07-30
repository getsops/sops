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
	"net/url"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/filter"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
)

var noop = martian.Noop("url.Filter")

func init() {
	parse.Register("url.Filter", filterFromJSON)
}

// Filter runs modifiers iff the request URL matches all of the segments in url.
type Filter struct {
	*filter.Filter
}

type filterJSON struct {
	Scheme       string               `json:"scheme"`
	Host         string               `json:"host"`
	Path         string               `json:"path"`
	Query        string               `json:"query"`
	Modifier     json.RawMessage      `json:"modifier"`
	ElseModifier json.RawMessage      `json:"else"`
	Scope        []parse.ModifierType `json:"scope"`
}

// NewFilter constructs a filter that applies the modifer when the
// request URL matches all of the provided URL segments.
func NewFilter(u *url.URL) *Filter {
	log.Debugf("martianurl.NewFilter: %s", u)
	m := NewMatcher(u)
	f := filter.New()
	f.SetRequestCondition(m)
	f.SetResponseCondition(m)
	return &Filter{f}
}

// filterFromJSON takes a JSON message as a byte slice and returns a
// parse.Result that contains a URLFilter and a bitmask that represents the
// type of modifier.
//
// Example JSON configuration message:
// {
//   "scheme": "https",
//   "host": "example.com",
//   "path": "/foo/bar",
//   "query": "q=value",
//   "scope": ["request", "response"],
//   "modifier": { ... }
//   "else": { ... }
// }
func filterFromJSON(b []byte) (*parse.Result, error) {
	msg := &filterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	filter := NewFilter(&url.URL{
		Scheme:   msg.Scheme,
		Host:     msg.Host,
		Path:     msg.Path,
		RawQuery: msg.Query,
	})

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
