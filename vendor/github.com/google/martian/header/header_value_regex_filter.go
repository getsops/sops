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
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

// ValueRegexFilter executes resmod and reqmod when the header
// value matches regex.
type ValueRegexFilter struct {
	regex  *regexp.Regexp
	header string
	reqmod martian.RequestModifier
	resmod martian.ResponseModifier
}

type headerValueRegexFilterJSON struct {
	Regex      string               `json:"regex"`
	HeaderName string               `json:"header"`
	Modifier   json.RawMessage      `json:"modifier"`
	Scope      []parse.ModifierType `json:"scope"`
}

func init() {
	parse.Register("header.RegexFilter", headerValueRegexFilterFromJSON)
}

// NewValueRegexFilter builds a new header value regex filter.
func NewValueRegexFilter(regex *regexp.Regexp, header string) *ValueRegexFilter {
	return &ValueRegexFilter{
		regex:  regex,
		header: header,
		reqmod: noop,
		resmod: noop,
	}
}

func headerValueRegexFilterFromJSON(b []byte) (*parse.Result, error) {
	msg := &headerValueRegexFilterJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	cr, err := regexp.Compile(msg.Regex)
	if err != nil {
		return nil, err
	}
	filter := NewValueRegexFilter(cr, msg.HeaderName)

	r, err := parse.FromJSON(msg.Modifier)
	if err != nil {
		return nil, err
	}

	reqmod := r.RequestModifier()
	filter.SetRequestModifier(reqmod)

	resmod := r.ResponseModifier()
	filter.SetResponseModifier(resmod)

	return parse.NewResult(filter, msg.Scope)
}

// ModifyRequest runs reqmod iff the value of header matches regex.
func (f *ValueRegexFilter) ModifyRequest(req *http.Request) error {
	hvalue := req.Header.Get(f.header)
	if hvalue == "" {
		return nil
	}

	if f.regex.MatchString(hvalue) {
		return f.reqmod.ModifyRequest(req)
	}

	return nil
}

// ModifyResponse runs resmod iff the value of request header matches regex.
func (f *ValueRegexFilter) ModifyResponse(res *http.Response) error {
	hvalue := res.Request.Header.Get(f.header)
	if hvalue == "" {
		return nil
	}

	if f.regex.MatchString(hvalue) {
		return f.resmod.ModifyResponse(res)
	}

	return nil
}

// SetRequestModifier sets the request modifier of HeaderValueRegexFilter.
func (f *ValueRegexFilter) SetRequestModifier(reqmod martian.RequestModifier) {
	if reqmod == nil {
		f.reqmod = noop
		return
	}

	f.reqmod = reqmod
}

// SetResponseModifier sets the response modifier of HeaderValueRegexFilter.
func (f *ValueRegexFilter) SetResponseModifier(resmod martian.ResponseModifier) {
	if resmod == nil {
		f.resmod = noop
		return
	}

	f.resmod = resmod
}
