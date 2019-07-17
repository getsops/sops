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
)

const idHeaderName string = "X-Martian-ID"

func init() {
	parse.Register("header.Id", idModifierFromJSON)
}

type idModifier struct{}

type idModifierJSON struct {
	Scope []parse.ModifierType `json:"scope"`
}

// NewIDModifier returns a request modifier that will set a header with the name
// X-Martian-ID with a value that is a unique identifier for the request. In the case
// that the X-Martian-ID header is already set, the header is unmodified.
func NewIDModifier() martian.RequestModifier {
	return &idModifier{}
}

// ModifyRequest sets the X-Martian-ID header with a unique identifier.  In the case
// that the X-Martian-ID header is already set, the header is unmodified.
func (im *idModifier) ModifyRequest(req *http.Request) error {
	// Do not rewrite an ID if req already has one
	if req.Header.Get(idHeaderName) != "" {
		return nil
	}

	ctx := martian.NewContext(req)
	req.Header.Set(idHeaderName, ctx.ID())

	return nil
}

func idModifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &idModifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	modifier := NewIDModifier()

	return parse.NewResult(modifier, msg.Scope)
}
