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

// Package stash provides a modifier that stores the request URL in a
// specified header.
package stash

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/martian/v3/parse"
)

func init() {
	parse.Register("stash.Modifier", modifierFromJSON)
}

// Modifier adds a header to the request containing the current state of the URL.
// The header will be named with the value stored in headerName.
// There will be no validation done on this header name.
type Modifier struct {
	headerName string
}

type modifierJSON struct {
	HeaderName string               `json:"headerName"`
	Scope      []parse.ModifierType `json:"scope"`
}

// NewModifier returns a RequestModifier that write the current URL into a header.
func NewModifier(headerName string) *Modifier {
	return &Modifier{headerName: headerName}
}

// ModifyRequest writes the current URL into a header.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	req.Header.Set(m.headerName, req.URL.String())
	return nil
}

// ModifyResponse writes the same header written in the request into the response.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	res.Header.Set(m.headerName, res.Request.Header.Get(m.headerName))
	return nil
}

func modifierFromJSON(b []byte) (*parse.Result, error) {
	// If you would like the saved state of the URL to be written in the response you must specify
	// this modifier's scope as both request and response.
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	mod := NewModifier(msg.HeaderName)
	r, err := parse.NewResult(mod, msg.Scope)
	if err != nil {
		return nil, err
	}

	if r.ResponseModifier() != nil && r.RequestModifier() == nil {
		return nil, fmt.Errorf("to write header on a response, specify scope as both request and response")
	}

	return r, nil
}
