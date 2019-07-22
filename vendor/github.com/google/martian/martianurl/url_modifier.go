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

// Package martianurl provides utilities for modifying, filtering,
// and verifying URLs in martian.Proxy.
package martianurl

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

// Modifier alters the request URL fields to match the fields of
// url and adds a X-Forwarded-Url header that contains the original
// value of the request URL.
type Modifier struct {
	url *url.URL
}

type modifierJSON struct {
	Scheme string               `json:"scheme"`
	Host   string               `json:"host"`
	Path   string               `json:"path"`
	Query  string               `json:"query"`
	Scope  []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("url.Modifier", modifierFromJSON)
}

// ModifyRequest sets the fields of req.URL to m.Url if they are not the zero value.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	if m.url.Scheme != "" {
		req.URL.Scheme = m.url.Scheme
	}
	if m.url.Host != "" {
		req.URL.Host = m.url.Host
	}
	if m.url.Path != "" {
		req.URL.Path = m.url.Path
	}
	if m.url.RawQuery != "" {
		req.URL.RawQuery = m.url.RawQuery
	}
	if m.url.Fragment != "" {
		req.URL.Fragment = m.url.Fragment
	}

	return nil
}

// NewModifier overrides the url of the request.
func NewModifier(url *url.URL) martian.RequestModifier {
	return &Modifier{
		url: url,
	}
}

// modifierFromJSON builds a martianurl.Modifier from JSON.
//
// Example modifier JSON:
// {
//   "martianurl.Modifier": {
//     "scope": ["request"],
//     "scheme": "https",
//     "host": "www.google.com",
//     "path": "/proxy",
//     "query": "testing=true"
//   }
// }
func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	mod := NewModifier(&url.URL{
		Scheme:   msg.Scheme,
		Host:     msg.Host,
		Path:     msg.Path,
		RawQuery: msg.Query,
	})

	return parse.NewResult(mod, msg.Scope)
}
