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

// Package martianhttp provides HTTP handlers for managing the state of a martian.Proxy.
package martianhttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

var noop = martian.Noop("martianhttp.Modifier")

// Modifier is a locking modifier that is configured via http.Handler.
type Modifier struct {
	mu     sync.RWMutex
	config []byte
	reqmod martian.RequestModifier
	resmod martian.ResponseModifier
}

// NewModifier returns a new martianhttp.Modifier.
func NewModifier() *Modifier {
	return &Modifier{
		reqmod: noop,
		resmod: noop,
	}
}

// SetRequestModifier sets the request modifier.
func (m *Modifier) SetRequestModifier(reqmod martian.RequestModifier) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.setRequestModifier(reqmod)
}

func (m *Modifier) setRequestModifier(reqmod martian.RequestModifier) {
	if reqmod == nil {
		reqmod = noop
	}

	m.reqmod = reqmod
}

// SetResponseModifier sets the response modifier.
func (m *Modifier) SetResponseModifier(resmod martian.ResponseModifier) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.setResponseModifier(resmod)
}

func (m *Modifier) setResponseModifier(resmod martian.ResponseModifier) {
	if resmod == nil {
		resmod = noop
	}

	m.resmod = resmod
}

// ModifyRequest runs reqmod.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.reqmod.ModifyRequest(req)
}

// ModifyResponse runs resmod.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.resmod.ModifyResponse(res)
}

// VerifyRequests verifies reqmod, iff reqmod is a RequestVerifier.
func (m *Modifier) VerifyRequests() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if reqv, ok := m.reqmod.(verify.RequestVerifier); ok {
		return reqv.VerifyRequests()
	}

	return nil
}

// VerifyResponses verifies resmod, iff resmod is a ResponseVerifier.
func (m *Modifier) VerifyResponses() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if resv, ok := m.resmod.(verify.ResponseVerifier); ok {
		return resv.VerifyResponses()
	}

	return nil
}

// ResetRequestVerifications resets verifications on reqmod, iff reqmod is a
// RequestVerifier.
func (m *Modifier) ResetRequestVerifications() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if reqv, ok := m.reqmod.(verify.RequestVerifier); ok {
		reqv.ResetRequestVerifications()
	}
}

// ResetResponseVerifications resets verifications on resmod, iff resmod is a
// ResponseVerifier.
func (m *Modifier) ResetResponseVerifications() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if resv, ok := m.resmod.(verify.ResponseVerifier); ok {
		resv.ResetResponseVerifications()
	}
}

// ServeHTTP sets or retrieves the JSON-encoded modifier configuration
// depending on request method. POST requests are expected to provide a JSON
// modifier message in the body which will be used to update the contained
// request and response modifiers. GET requests will return the JSON
// (pretty-printed) for the most recent configuration.
func (m *Modifier) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		m.servePOST(rw, req)
		return
	case "GET":
		m.serveGET(rw, req)
		return
	default:
		rw.Header().Set("Allow", "GET, POST")
		rw.WriteHeader(405)
	}
}

func (m *Modifier) servePOST(rw http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		log.Errorf("martianhttp: error reading request body: %v", err)
		return
	}
	req.Body.Close()

	r, err := parse.FromJSON(body)
	if err != nil {
		http.Error(rw, err.Error(), 400)
		log.Errorf("martianhttp: error parsing JSON: %v", err)
		return
	}

	buf := new(bytes.Buffer)
	if err := json.Indent(buf, body, "", "  "); err != nil {
		http.Error(rw, err.Error(), 400)
		log.Errorf("martianhttp: error formatting JSON: %v", err)
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = buf.Bytes()
	m.setRequestModifier(r.RequestModifier())
	m.setResponseModifier(r.ResponseModifier())
}

func (m *Modifier) serveGET(rw http.ResponseWriter, req *http.Request) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	rw.Header().Set("Content-Type", "application/json")
	rw.Write(m.config)
}
