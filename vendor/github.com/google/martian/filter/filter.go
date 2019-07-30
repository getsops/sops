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

// Package filter provides a modifier that executes a given set of child
// modifiers based on the evaluated value of the provided conditional.
package filter

import (
	"fmt"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/verify"
)

var noop = martian.Noop("Filter")

// Filter is a modifer that contains conditions to evaluate on request and
// response as well as a set of modifiers to execute based on the value of
// the provided RequestCondition or ResponseCondition.
type Filter struct {
	reqcond RequestCondition
	rescond ResponseCondition

	treqmod martian.RequestModifier
	tresmod martian.ResponseModifier
	freqmod martian.RequestModifier
	fresmod martian.ResponseModifier
}

// New returns a pointer to a Filter with all child modifiers initialized to
// the noop modifier.
func New() *Filter {
	return &Filter{
		treqmod: noop,
		tresmod: noop,
		fresmod: noop,
		freqmod: noop,
	}
}

// SetRequestCondition sets the condition to evaluate on requests.
func (f *Filter) SetRequestCondition(reqcond RequestCondition) {
	f.reqcond = reqcond
}

// SetResponseCondition sets the condition to evaluate on responses.
func (f *Filter) SetResponseCondition(rescond ResponseCondition) {
	f.rescond = rescond
}

// SetRequestModifier sets the martian.RequestModifier that is executed
// when the RequestCondition evaluates to True.  This function is provided
// to maintain backwards compatability with filtering prior to filter.Filter.
func (f *Filter) SetRequestModifier(reqmod martian.RequestModifier) {
	f.RequestWhenTrue(reqmod)
}

// RequestWhenTrue sets the martian.RequestModifier that is executed
// when the RequestCondition evaluates to True.
func (f *Filter) RequestWhenTrue(mod martian.RequestModifier) {
	if mod == nil {
		f.treqmod = noop
		return
	}

	f.treqmod = mod
}

// SetResponseModifier sets the martian.ResponseModifier that is executed
// when the ResponseCondition evaluates to True.  This function is provided
// to maintain backwards compatability with filtering prior to filter.Filter.
func (f *Filter) SetResponseModifier(resmod martian.ResponseModifier) {
	f.ResponseWhenTrue(resmod)
}

// RequestWhenFalse sets the martian.RequestModifier that is executed
// when the RequestCondition evaluates to False.
func (f *Filter) RequestWhenFalse(mod martian.RequestModifier) {
	if mod == nil {
		f.freqmod = noop
		return
	}

	f.freqmod = mod
}

// ResponseWhenTrue sets the martian.ResponseModifier that is executed
// when the ResponseCondition evaluates to True.
func (f *Filter) ResponseWhenTrue(mod martian.ResponseModifier) {
	if mod == nil {
		f.tresmod = noop
		return
	}

	f.tresmod = mod
}

// ResponseWhenFalse sets the martian.ResponseModifier that is executed
// when the ResponseCondition evaluates to False.
func (f *Filter) ResponseWhenFalse(mod martian.ResponseModifier) {
	if mod == nil {
		f.fresmod = noop
		return
	}

	f.fresmod = mod
}

// ModifyRequest evaluates reqcond and executes treqmod iff reqcond evaluates
// to true; otherwise, freqmod is executed.
func (f *Filter) ModifyRequest(req *http.Request) error {
	if f.reqcond == nil {
		return fmt.Errorf("filter.ModifyRequest: no request condition set. Set condition with SetRequestCondition")
	}

	match := f.reqcond.MatchRequest(req)
	if match {
		log.Debugf("filter.ModifyRequest: matched %s", req.URL)
		return f.treqmod.ModifyRequest(req)
	}

	return f.freqmod.ModifyRequest(req)
}

// ModifyResponse evaluates rescond and executes tresmod iff rescond evaluates
// to true; otherwise, fresmod is executed.
func (f *Filter) ModifyResponse(res *http.Response) error {
	if f.rescond == nil {
		return fmt.Errorf("filter.ModifyResponse: no response condition set. Set condition with SetResponseCondition")
	}

	match := f.rescond.MatchResponse(res)
	if match {
		requ := ""
		if res.Request != nil {
			requ = res.Request.URL.String()
		}
		log.Debugf("filter.ModifyResponse: %s", requ)
		return f.tresmod.ModifyResponse(res)
	}

	return f.fresmod.ModifyResponse(res)
}

// VerifyRequests returns an error containing all the verification errors
// returned by request verifiers.
func (f *Filter) VerifyRequests() error {
	merr := martian.NewMultiError()

	freqv, ok := f.freqmod.(verify.RequestVerifier)
	if ok {
		if ve := freqv.VerifyRequests(); ve != nil {
			merr.Add(ve)
		}
	}

	treqv, ok := f.treqmod.(verify.RequestVerifier)
	if ok {
		if ve := treqv.VerifyRequests(); ve != nil {
			merr.Add(ve)
		}
	}

	if merr.Empty() {
		return nil
	}

	return merr
}

// VerifyResponses returns an error containing all the verification errors
// returned by response verifiers.
func (f *Filter) VerifyResponses() error {
	merr := martian.NewMultiError()

	tresv, ok := f.tresmod.(verify.ResponseVerifier)
	if ok {
		if ve := tresv.VerifyResponses(); ve != nil {
			merr.Add(ve)
		}
	}

	fresv, ok := f.fresmod.(verify.ResponseVerifier)
	if ok {
		if ve := fresv.VerifyResponses(); ve != nil {
			merr.Add(ve)
		}
	}

	if merr.Empty() {
		return nil
	}

	return merr
}

// ResetRequestVerifications resets the state of the contained request verifiers.
func (f *Filter) ResetRequestVerifications() {
	if treqv, ok := f.treqmod.(verify.RequestVerifier); ok {
		treqv.ResetRequestVerifications()
	}
	if freqv, ok := f.freqmod.(verify.RequestVerifier); ok {
		freqv.ResetRequestVerifications()
	}
}

// ResetResponseVerifications resets the state of the contained request verifiers.
func (f *Filter) ResetResponseVerifications() {
	if tresv, ok := f.tresmod.(verify.ResponseVerifier); ok {
		tresv.ResetResponseVerifications()
	}
}
