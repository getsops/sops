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

package marbl

import (
	"io"
	"net/http"

	"github.com/google/martian/v3"
)

// Modifier implements the Martian modifier interface so that marbl logs
// can be captured at any point in a Martian modifier tree.
type Modifier struct {
	s *Stream
}

// NewModifier returns a marbl.Modifier initialized with a marbl.Stream.
func NewModifier(w io.Writer) *Modifier {
	return &Modifier{
		s: NewStream(w),
	}
}

// ModifyRequest writes an HTTP request to the log stream.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	return m.s.LogRequest(ctx.ID(), req)
}

// ModifyResponse writes an HTTP response to the log stream.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	return m.s.LogResponse(ctx.ID(), res)
}
