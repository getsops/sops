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
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func init() {
	parse.Register("header.Append", appendModifierFromJSON)
}

type appendModifier struct {
	name, value string
}

type appendModifierJSON struct {
	Name  string               `json:"name"`
	Value string               `json:"value"`
	Scope []parse.ModifierType `json:"scope"`
}

// ModifyRequest appends the header at name with value to the request.
func (m *appendModifier) ModifyRequest(req *http.Request) error {
	return proxyutil.RequestHeader(req).Add(m.name, m.value)
}

// ModifyResponse appends the header at name with value to the response.
func (m *appendModifier) ModifyResponse(res *http.Response) error {
	return proxyutil.ResponseHeader(res).Add(m.name, m.value)
}

// NewAppendModifier returns an appendModifier that will append a header with
// with the given name and value for both requests and responses. Existing
// headers with the same name will be left in place.
func NewAppendModifier(name, value string) martian.RequestResponseModifier {
	return &appendModifier{
		name:  http.CanonicalHeaderKey(name),
		value: value,
	}
}

// appendModifierFromJSON takes a JSON message as a byte slice and returns
// an appendModifier and an error.
//
// Example JSON configuration message:
// {
//  "scope": ["request", "result"],
//  "name": "X-Martian",
//  "value": "true"
// }
func appendModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	modifier := NewAppendModifier(msg.Name, msg.Value)

	return parse.NewResult(modifier, msg.Scope)
}
