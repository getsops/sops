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

// Package parse constructs martian modifiers from JSON messages.
package parse

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/martian/v3"
)

// ModifierType is the HTTP message type.
type ModifierType string

const (
	// Request modifies an HTTP request.
	Request ModifierType = "request"
	// Response modifies an HTTP response.
	Response ModifierType = "response"
)

// Result holds the parsed modifier and its type.
type Result struct {
	reqmod martian.RequestModifier
	resmod martian.ResponseModifier
}

// NewResult returns a new parse.Result for a given interface{} that implements a modifier
// and a slice of scopes to generate the result for.
//
// Returns nil, error if a given modifier does not support a given scope
func NewResult(mod interface{}, scope []ModifierType) (*Result, error) {
	reqmod, reqOk := mod.(martian.RequestModifier)
	resmod, resOk := mod.(martian.ResponseModifier)
	result := &Result{}
	if scope == nil {
		result.reqmod = reqmod
		result.resmod = resmod
		return result, nil
	}

	for _, s := range scope {
		switch s {
		case Request:
			if !reqOk {
				return nil, fmt.Errorf("parse: invalid scope %q for modifier", "request")
			}

			result.reqmod = reqmod
		case Response:
			if !resOk {
				return nil, fmt.Errorf("parse: invalid scope %q for modifier", "response")
			}

			result.resmod = resmod
		default:
			return nil, fmt.Errorf("parse: invalid scope: %s not in [%q, %q]", s, "request", "response")
		}
	}

	return result, nil
}

// RequestModifier returns the parsed RequestModifier.
//
// Returns nil if the message has no request modifier.
func (r *Result) RequestModifier() martian.RequestModifier {
	return r.reqmod
}

// ResponseModifier returns the parsed ResponseModifier.
//
// Returns nil if the message has no response modifier.
func (r *Result) ResponseModifier() martian.ResponseModifier {
	return r.resmod
}

var (
	parseMu    sync.RWMutex
	parseFuncs = make(map[string]func(b []byte) (*Result, error))
)

// ErrUnknownModifier is the error returned when the message does not
// contain a field representing a known modifier type.
type ErrUnknownModifier struct {
	name string
}

// Error returns a formatted error message for an ErrUnknownModifier.
func (e ErrUnknownModifier) Error() string {
	return fmt.Sprintf("parse: unknown modifier: %s", e.name)
}

// Register registers a parsing function for name that will be used to unmarshal
// a JSON message into the appropriate modifier.
func Register(name string, parseFunc func(b []byte) (*Result, error)) {
	parseMu.Lock()
	defer parseMu.Unlock()

	parseFuncs[name] = parseFunc
}

// FromJSON parses a Modifier JSON message by looking up the named modifier in parseFuncs
// and passing its modifier to the registered parseFunc. Returns a parse.Result containing
// the top-level parsed modifier. If no parser has been registered with the given name
// it returns an error of type ErrUnknownModifier.
func FromJSON(b []byte) (*Result, error) {
	msg := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &msg); err != nil {
		return nil, err
	}

	if len(msg) != 1 {
		ks := ""
		for k := range msg {
			ks += ", " + k
		}

		return nil, fmt.Errorf("parse: expected one modifier, received %d: %s", len(msg), ks)
	}

	parseMu.RLock()
	defer parseMu.RUnlock()
	for k, m := range msg {
		parseFunc, ok := parseFuncs[k]
		if !ok {
			return nil, ErrUnknownModifier{name: k}
		}
		return parseFunc(m)
	}

	return nil, fmt.Errorf("parse: no modifiers found: %v", msg)
}
