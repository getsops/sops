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
	"fmt"
	"net/http"
	"strings"

	"github.com/google/martian/v3"
)

// NewBadFramingModifier makes a best effort to fix inconsistencies in the
// request such as multiple Content-Lengths or the lack of Content-Length and
// improper Transfer-Encoding. If it is unable to determine a proper resolution
// it returns an error.
//
// http://tools.ietf.org/html/draft-ietf-httpbis-p1-messaging-14#section-3.3
func NewBadFramingModifier() martian.RequestModifier {
	return martian.RequestModifierFunc(
		func(req *http.Request) error {
			cls := req.Header["Content-Length"]
			if len(cls) > 0 {
				var length string

				// Iterate over all Content-Length headers, splitting any we find with
				// commas, and check that all Content-Lengths are equal.
				for _, ls := range cls {
					for _, l := range strings.Split(ls, ",") {
						// First length, set it as the canonical Content-Length.
						if length == "" {
							length = strings.TrimSpace(l)
							continue
						}

						// Mismatched Content-Lengths.
						if length != strings.TrimSpace(l) {
							return fmt.Errorf(`bad request framing: multiple mismatched "Content-Length" headers: %v`, cls)
						}
					}
				}

				// All Content-Lengths are equal, remove extras and set it to the
				// canonical value.
				req.Header.Set("Content-Length", length)
			}

			tes := req.Header["Transfer-Encoding"]
			if len(tes) > 0 {
				// Extract the last Transfer-Encoding value, and split on commas.
				last := strings.Split(tes[len(tes)-1], ",")

				// Check that the last, potentially comma-delimited, value is
				// "chunked", else we have no way to determine when the request is
				// finished.
				if strings.TrimSpace(last[len(last)-1]) != "chunked" {
					return fmt.Errorf(`bad request framing: "Transfer-Encoding" header is present, but does not end in "chunked"`)
				}

				// Transfer-Encoding "chunked" takes precedence over
				// Content-Length.
				req.Header.Del("Content-Length")
			}

			return nil
		})
}
