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
	"net/http"
	"strings"

	"github.com/google/martian/v3"
)

// Hop-by-hop headers as defined by RFC2616.
//
// http://tools.ietf.org/html/draft-ietf-httpbis-p1-messaging-14#section-7.1.3.1
var hopByHopHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Proxy-Connection", // Non-standard, but required for HTTP/2.
	"Te",
	"Trailer",
	"Transfer-Encoding",
	"Upgrade",
}

type hopByHopModifier struct{}

// NewHopByHopModifier removes Hop-By-Hop headers from requests and
// responses.
func NewHopByHopModifier() martian.RequestResponseModifier {
	return &hopByHopModifier{}
}

// ModifyRequest removes all hop-by-hop headers defined by RFC2616 as
// well as any additional hop-by-hop headers specified in the
// Connection header.
func (m *hopByHopModifier) ModifyRequest(req *http.Request) error {
	removeHopByHopHeaders(req.Header)
	return nil
}

// ModifyResponse removes all hop-by-hop headers defined by RFC2616 as
// well as any additional hop-by-hop headers specified in the
// Connection header.
func (m *hopByHopModifier) ModifyResponse(res *http.Response) error {
	removeHopByHopHeaders(res.Header)
	return nil
}

func removeHopByHopHeaders(header http.Header) {
	// Additional hop-by-hop headers may be specified in `Connection` headers.
	// http://tools.ietf.org/html/draft-ietf-httpbis-p1-messaging-14#section-9.1
	for _, vs := range header["Connection"] {
		for _, v := range strings.Split(vs, ",") {
			k := http.CanonicalHeaderKey(strings.TrimSpace(v))
			header.Del(k)
		}
	}

	for _, k := range hopByHopHeaders {
		header.Del(k)
	}
}
