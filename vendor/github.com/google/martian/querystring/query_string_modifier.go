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

// Package querystring contains a modifier to rewrite query strings in a request.
package querystring

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

func init() {
	parse.Register("querystring.Modifier", modifierFromJSON)
}

type modifier struct {
	key, value string
}

type modifierJSON struct {
	Name  string               `json:"name"`
	Value string               `json:"value"`
	Scope []parse.ModifierType `json:"scope"`
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *modifier) ModifyRequest(req *http.Request) error {
	query := req.URL.Query()
	query.Set(m.key, m.value)
	req.URL.RawQuery = query.Encode()

	return nil
}

// NewModifier returns a request modifier that will set the query string
// at key with the given value. If the query string key already exists all
// values will be overwritten.
func NewModifier(key, value string) martian.RequestModifier {
	return &modifier{
		key:   key,
		value: value,
	}
}

// modifierFromJSON takes a JSON message as a byte slice and returns
// a querystring.modifier and an error.
//
// Example JSON:
// {
//  "name": "param",
//  "value": "true",
//  "scope": ["request", "response"]
// }
func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewModifier(msg.Name, msg.Value), msg.Scope)
}
