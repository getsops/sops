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
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func init() {
	parse.Register("header.Copy", copyModifierFromJSON)
}

type copyModifier struct {
	from, to string
}

type copyModifierJSON struct {
	From  string               `json:"from"`
	To    string               `json:"to"`
	Scope []parse.ModifierType `json:"scope"`
}

// ModifyRequest copies the header in from to the request header for to.
func (m *copyModifier) ModifyRequest(req *http.Request) error {
	log.Debugf("header: copyModifier.ModifyRequest %s, from: %s, to: %s", req.URL, m.from, m.to)
	h := proxyutil.RequestHeader(req)

	return h.Set(m.to, h.Get(m.from))
}

// ModifyResponse copies the header in from to the response header for to.
func (m *copyModifier) ModifyResponse(res *http.Response) error {
	log.Debugf("header: copyModifier.ModifyResponse %s, from: %s, to: %s", res.Request.URL, m.from, m.to)
	h := proxyutil.ResponseHeader(res)

	return h.Set(m.to, h.Get(m.from))
}

// NewCopyModifier returns a modifier that will copy the header in from to the
// header in to.
func NewCopyModifier(from, to string) martian.RequestResponseModifier {
	return &copyModifier{
		from: from,
		to:   to,
	}
}

// copyModifierFromJSON builds a copy modifier from JSON.
//
// Example JSON:
// {
//   "header.Copy": {
//     "scope": ["request", "response"],
//     "from": "Original-Header",
//     "to": "Copy-Header"
//   }
// }
func copyModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &copyModifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewCopyModifier(msg.From, msg.To), msg.Scope)
}
