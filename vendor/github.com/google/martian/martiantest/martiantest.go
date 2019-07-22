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

// Package martiantest provides helper utilities for testing
// modifiers.
package martiantest

import (
	"net/http"
	"sync/atomic"
)

// Modifier keeps track of the number of requests and responses it has modified
// and can be configured to return errors or run custom functions.
type Modifier struct {
	reqcount int32 // atomic
	rescount int32 // atomic
	reqerr   error
	reserr   error
	reqfunc  func(*http.Request)
	resfunc  func(*http.Response)
}

// NewModifier returns a new test modifier.
func NewModifier() *Modifier {
	return &Modifier{}
}

// RequestCount returns the number of requests modified.
func (m *Modifier) RequestCount() int32 {
	return atomic.LoadInt32(&m.reqcount)
}

// ResponseCount returns the number of responses modified.
func (m *Modifier) ResponseCount() int32 {
	return atomic.LoadInt32(&m.rescount)
}

// RequestModified returns whether a request has been modified.
func (m *Modifier) RequestModified() bool {
	return m.RequestCount() != 0
}

// ResponseModified returns whether a response has been modified.
func (m *Modifier) ResponseModified() bool {
	return m.ResponseCount() != 0
}

// RequestError overrides the error returned by ModifyRequest.
func (m *Modifier) RequestError(err error) {
	m.reqerr = err
}

// ResponseError overrides the error returned by ModifyResponse.
func (m *Modifier) ResponseError(err error) {
	m.reserr = err
}

// RequestFunc is a function to run during ModifyRequest.
func (m *Modifier) RequestFunc(reqfunc func(req *http.Request)) {
	m.reqfunc = reqfunc
}

// ResponseFunc is a function to run during ModifyResponse.
func (m *Modifier) ResponseFunc(resfunc func(res *http.Response)) {
	m.resfunc = resfunc
}

// ModifyRequest increases the count of requests seen and runs reqfunc if configured.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	atomic.AddInt32(&m.reqcount, 1)

	if m.reqfunc != nil {
		m.reqfunc(req)
	}

	return m.reqerr
}

// ModifyResponse increases the count of responses seen and runs resfunc if configured.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	atomic.AddInt32(&m.rescount, 1)

	if m.resfunc != nil {
		m.resfunc(res)
	}

	return m.reserr
}

// Reset resets the request and response counts, the custom
// functions, and the modifier errors.
func (m *Modifier) Reset() {
	atomic.StoreInt32(&m.reqcount, 0)
	atomic.StoreInt32(&m.rescount, 0)

	m.reqfunc = nil
	m.resfunc = nil

	m.reqerr = nil
	m.reserr = nil
}

// Matcher is a stubbed matcher used in tests.
type Matcher struct {
	resval bool
	reqval bool
}

// NewMatcher returns a pointer to martiantest.Matcher with the return values
// for MatchRequest and MatchResponse intiailized to true.
func NewMatcher() *Matcher {
	return &Matcher{resval: true, reqval: true}
}

// ResponseEvaluatesTo sets the value returned by MatchResponse.
func (tm *Matcher) ResponseEvaluatesTo(value bool) {
	tm.resval = value
}

// RequestEvaluatesTo sets the value returned by MatchRequest.
func (tm *Matcher) RequestEvaluatesTo(value bool) {
	tm.reqval = value
}

// MatchRequest returns the stubbed value in tm.reqval.
func (tm *Matcher) MatchRequest(*http.Request) bool {
	return tm.reqval
}

// MatchResponse returns the stubbed value in tm.resval.
func (tm *Matcher) MatchResponse(*http.Response) bool {
	return tm.resval
}
