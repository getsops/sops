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

package martian

import (
	"net/http"

	"github.com/google/martian/v3/log"
)

type noopModifier struct {
	id string
}

// Noop returns a modifier that does not change the request or the response.
func Noop(id string) RequestResponseModifier {
	return &noopModifier{
		id: id,
	}
}

// ModifyRequest logs a debug line.
func (nm *noopModifier) ModifyRequest(*http.Request) error {
	log.Debugf("%s: no request modifier configured", nm.id)
	return nil
}

// ModifyResponse logs a debug line.
func (nm *noopModifier) ModifyResponse(*http.Response) error {
	log.Debugf("%s: no response modifier configured", nm.id)
	return nil
}
