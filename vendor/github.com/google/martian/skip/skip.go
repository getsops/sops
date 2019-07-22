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

// Package skip provides a request modifier to skip the HTTP round-trip.
package skip

import (
	"encoding/json"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

// RoundTrip is a modifier that skips the request round-trip.
type RoundTrip struct{}

type roundTripJSON struct {
	Scope []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("skip.RoundTrip", roundTripFromJSON)
}

// NewRoundTrip returns a new modifier that skips round-trip.
func NewRoundTrip() *RoundTrip {
	return &RoundTrip{}
}

// ModifyRequest skips the request round-trip.
func (r *RoundTrip) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	ctx.SkipRoundTrip()

	return nil
}

// roundTripFromJSON builds a skip.RoundTrip from JSON.

// Example JSON:
// {
//   "skip.RoundTrip": { }
// }
func roundTripFromJSON(b []byte) (*parse.Result, error) {
	msg := &roundTripJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(NewRoundTrip(), msg.Scope)
}
