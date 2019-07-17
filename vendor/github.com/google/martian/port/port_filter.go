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

// Package port provides utilities for modifying and filtering
// based on the port of request URLs.
package port

import (
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

var noop = martian.Noop("port.Filter")

func init() {
	parse.Register("port.Filter", filterFromJSON)
}

// Filter runs modifiers iff the port in the request URL matches port.
type Filter struct {
	reqmod martian.RequestModifier
	resmod martian.ResponseModifier
	port   int
}

type filterJSON struct {
	Port     int                  `json:"port"`
	Modifier json.RawMessage      `json:"modifier"`
	Scope    []parse.ModifierType `json:"scope"`
}

// NewFilter returns a filter that executes modifiers if the port of
// request matches port.
func NewFilter(port int) *Filter {
	return &Filter{
		port:   port,
		reqmod: noop,
		resmod: noop,
	}
}

// SetRequestModifier sets the request modifier.
func (f *Filter) SetRequestModifier(reqmod martian.RequestModifier) {
	if reqmod == nil {
		reqmod = noop
	}

	f.reqmod = reqmod
}

// SetResponseModifier sets the response modifier.
func (f *Filter) SetResponseModifier(resmod martian.ResponseModifier) {
	if resmod == nil {
		resmod = noop
	}

	f.resmod = resmod
}

// ModifyRequest runs the modifier if the port matches the provided port.
func (f *Filter) ModifyRequest(req *http.Request) error {
	var defaultPort int
	if req.URL.Scheme == "http" {
		defaultPort = 80
	}
	if req.URL.Scheme == "https" {
		defaultPort = 443
	}

	hasPort := strings.Contains(req.URL.Host, ":")
	if hasPort {
		_, p, err := net.SplitHostPort(req.URL.Host)
		if err != nil {
			return err
		}

		pt, err := strconv.Atoi(p)
		if err != nil {
			return err
		}
		if pt == f.port {
			return f.reqmod.ModifyRequest(req)
		}
		return nil
	}

	// no port explictly declared - default port
	if f.port == defaultPort {
		return f.reqmod.ModifyRequest(req)
	}

	return nil
}

// ModifyResponse runs the modifier if the request URL matches urlMatcher.
func (f *Filter) ModifyResponse(res *http.Response) error {
	var defaultPort int
	if res.Request.URL.Scheme == "http" {
		defaultPort = 80
	}
	if res.Request.URL.Scheme == "https" {
		defaultPort = 443
	}

	if !strings.Contains(res.Request.URL.Host, ":") && (f.port == defaultPort) {
		return f.resmod.ModifyResponse(res)
	}

	_, p, err := net.SplitHostPort(res.Request.URL.Host)
	if err != nil {
		return err
	}

	pt, err := strconv.Atoi(p)
	if err != nil {
		return err
	}
	if pt == f.port {
		return f.resmod.ModifyResponse(res)
	}

	return nil
}

func filterFromJSON(b []byte) (*parse.Result, error) {
	msg := &filterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	filter := NewFilter(msg.Port)
	r, err := parse.FromJSON(msg.Modifier)
	if err != nil {
		return nil, err
	}

	reqmod := r.RequestModifier()
	if err != nil {
		return nil, err
	}
	if reqmod != nil {
		filter.SetRequestModifier(reqmod)
	}

	resmod := r.ResponseModifier()
	if resmod != nil {
		filter.SetResponseModifier(resmod)
	}

	return parse.NewResult(filter, msg.Scope)
}
