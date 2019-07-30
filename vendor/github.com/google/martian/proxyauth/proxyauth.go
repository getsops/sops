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

// Package proxyauth provides authentication support via the
// Proxy-Authorization header.
package proxyauth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/auth"
)

var noop = martian.Noop("proxyauth.Modifier")

// Modifier is the proxy authentication modifier.
type Modifier struct {
	reqmod martian.RequestModifier
	resmod martian.ResponseModifier
}

// NewModifier returns a new proxy authentication modifier.
func NewModifier() *Modifier {
	return &Modifier{
		reqmod: noop,
		resmod: noop,
	}
}

// SetRequestModifier sets the request modifier.
func (m *Modifier) SetRequestModifier(reqmod martian.RequestModifier) {
	if reqmod == nil {
		reqmod = noop
	}

	m.reqmod = reqmod
}

// SetResponseModifier sets the response modifier.
func (m *Modifier) SetResponseModifier(resmod martian.ResponseModifier) {
	if resmod == nil {
		resmod = noop
	}

	m.resmod = resmod
}

// ModifyRequest sets the auth ID in the context from the request iff it has
// not already been set and runs reqmod.ModifyRequest. If the underlying
// modifier has indicated via auth error that no valid auth credentials
// have been found we set ctx.SkipRoundTrip.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	actx := auth.FromContext(ctx)

	actx.SetID(id(req.Header))

	err := m.reqmod.ModifyRequest(req)

	if actx.Error() != nil {
		ctx.SkipRoundTrip()
	}

	return err
}

// ModifyResponse runs resmod.ModifyResponse and modifies the response to
// include the correct status code and headers if auth error is present.
//
// If an error is returned from resmod.ModifyResponse it is returned.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	actx := auth.FromContext(ctx)

	err := m.resmod.ModifyResponse(res)

	if actx.Error() != nil {
		res.StatusCode = http.StatusProxyAuthRequired
		res.Header.Set("Proxy-Authenticate", "Basic")
	}

	return err
}

// id returns an ID derived from the Proxy-Authorization header username and password.
func id(header http.Header) string {
	id := strings.TrimPrefix(header.Get("Proxy-Authorization"), "Basic ")

	data, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return ""
	}

	return string(data)
}
