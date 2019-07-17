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

// Package cookie allows for the modification of cookies on http requests and responses.
package cookie

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
)

func init() {
	parse.Register("cookie.Modifier", modifierFromJSON)
}

type modifier struct {
	cookie *http.Cookie
}

type modifierJSON struct {
	Name     string               `json:"name"`
	Value    string               `json:"value"`
	Path     string               `json:"path"`
	Domain   string               `json:"domain"`
	Expires  time.Time            `json:"expires"`
	Secure   bool                 `json:"secure"`
	HTTPOnly bool                 `json:"httpOnly"`
	MaxAge   int                  `json:"maxAge"`
	Scope    []parse.ModifierType `json:"scope"`
}

// ModifyRequest adds cookie to the request.
func (m *modifier) ModifyRequest(req *http.Request) error {
	req.AddCookie(m.cookie)
	log.Debugf("cookie.ModifyRequest: %s: cookie: %s", req.URL, m.cookie)

	return nil
}

// ModifyResponse sets cookie on the response.
func (m *modifier) ModifyResponse(res *http.Response) error {
	res.Header.Add("Set-Cookie", m.cookie.String())
	log.Debugf("cookie.ModifyResponse: %s: cookie: %s", res.Request.URL, m.cookie)

	return nil
}

// NewModifier returns a modifier that injects the provided cookie into the
// request or response.
func NewModifier(cookie *http.Cookie) martian.RequestResponseModifier {
	return &modifier{
		cookie: cookie,
	}
}

// modifierFromJSON takes a JSON message as a byte slice and returns a
// CookieModifier and an error.
//
// Example JSON Configuration message:
// {
//   "name": "Martian-Cookie",
//   "value": "some value",
//   "path": "/some/path",
//   "domain": "example.com",
//   "expires": "2025-04-12T23:20:50.52Z", // RFC 3339
//   "secure": true,
//   "httpOnly": false,
//   "maxAge": 86400,
//   "scope": ["request", "result"]
// }
func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	c := &http.Cookie{
		Name:     msg.Name,
		Value:    msg.Value,
		Path:     msg.Path,
		Domain:   msg.Domain,
		Expires:  msg.Expires,
		Secure:   msg.Secure,
		HttpOnly: msg.HTTPOnly,
		MaxAge:   msg.MaxAge,
	}

	return parse.NewResult(NewModifier(c), msg.Scope)
}
