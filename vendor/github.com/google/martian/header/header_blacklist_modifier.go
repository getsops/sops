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
	parse.Register("header.Blacklist", blacklistModifierFromJSON)
}

type blacklistModifier struct {
	names []string
}

type blacklistModifierJSON struct {
	Names []string             `json:"names"`
	Scope []parse.ModifierType `json:"scope"`
}

// ModifyRequest deletes all request headers based on the header name.
func (m *blacklistModifier) ModifyRequest(req *http.Request) error {
	h := proxyutil.RequestHeader(req)

	for _, name := range m.names {
		h.Del(name)
	}

	return nil
}

// ModifyResponse deletes all response headers based on the header name.
func (m *blacklistModifier) ModifyResponse(res *http.Response) error {
	h := proxyutil.ResponseHeader(res)

	for _, name := range m.names {
		h.Del(name)
	}

	return nil
}

// NewBlacklistModifier returns a modifier that will delete any header that
// matches a name contained in the names parameter.
func NewBlacklistModifier(names ...string) martian.RequestResponseModifier {
	return &blacklistModifier{
		names: names,
	}
}

// blacklistModifierFromJSON takes a JSON message as a byte slice and returns
// a blacklistModifier and an error.
//
// Example JSON configuration message:
// {
//   "names": ["X-Header", "Y-Header"],
//   "scope": ["request", "result"]
// }
func blacklistModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &blacklistModifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewBlacklistModifier(msg.Names...), msg.Scope)
}
