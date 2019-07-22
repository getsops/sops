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

// Package status contains a modifier to rewrite the status code on a response.
package status

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

type statusModifier struct {
	statusCode int
}

type statusJSON struct {
	StatusCode int                  `json:"statusCode"`
	Scope      []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("status.Modifier", modifierFromJSON)
}

// ModifyResponse overwrites the status text and code on an HTTP response and
// returns nil.
func (s *statusModifier) ModifyResponse(res *http.Response) error {
	res.StatusCode = s.statusCode
	res.Status = fmt.Sprintf("%d %s", s.statusCode, http.StatusText(s.statusCode))

	return nil
}

// NewModifier constructs a statusModifier that overrides response status
// codes with the HTTP status code provided.
func NewModifier(statusCode int) martian.ResponseModifier {
	return &statusModifier{
		statusCode: statusCode,
	}
}

// modifierFromJSON builds a status.Modifier from JSON.
//
// Example JSON:
// {
//   "status.Modifier": {
//     "scope": ["response"],
//     "statusCode": 401
//   }
// }
func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &statusJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewModifier(msg.StatusCode), msg.Scope)
}
