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

package proxyutil

import (
	"fmt"
	"net/http"
	"strconv"
)

// Header is a generic representation of a set of HTTP headers for requests and
// responses.
type Header struct {
	h http.Header

	host func() string
	cl   func() int64
	te   func() []string

	setHost func(string)
	setCL   func(int64)
	setTE   func([]string)
}

// RequestHeader returns a new set of headers from a request.
func RequestHeader(req *http.Request) *Header {
	return &Header{
		h:       req.Header,
		host:    func() string { return req.Host },
		cl:      func() int64 { return req.ContentLength },
		te:      func() []string { return req.TransferEncoding },
		setHost: func(host string) { req.Host = host },
		setCL:   func(cl int64) { req.ContentLength = cl },
		setTE:   func(te []string) { req.TransferEncoding = te },
	}
}

// ResponseHeader returns a new set of headers from a request.
func ResponseHeader(res *http.Response) *Header {
	return &Header{
		h:       res.Header,
		host:    func() string { return "" },
		cl:      func() int64 { return res.ContentLength },
		te:      func() []string { return res.TransferEncoding },
		setHost: func(string) {},
		setCL:   func(cl int64) { res.ContentLength = cl },
		setTE:   func(te []string) { res.TransferEncoding = te },
	}
}

// Set sets value at header name for the request or response.
func (h *Header) Set(name, value string) error {
	switch http.CanonicalHeaderKey(name) {
	case "Host":
		h.setHost(value)
	case "Content-Length":
		cl, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		h.setCL(cl)
	case "Transfer-Encoding":
		h.setTE([]string{value})
	default:
		h.h.Set(name, value)
	}

	return nil
}

// Add appends the value to the existing header at name for the request or
// response.
func (h *Header) Add(name, value string) error {
	switch http.CanonicalHeaderKey(name) {
	case "Host":
		if h.host() != "" {
			return fmt.Errorf("proxyutil: illegal header multiple: %s", "Host")
		}

		return h.Set(name, value)
	case "Content-Length":
		if h.cl() > 0 {
			return fmt.Errorf("proxyutil: illegal header multiple: %s", "Content-Length")
		}

		return h.Set(name, value)
	case "Transfer-Encoding":
		h.setTE(append(h.te(), value))
	default:
		h.h.Add(name, value)
	}

	return nil
}

// Get returns the first value at header name for the request or response.
func (h *Header) Get(name string) string {
	switch http.CanonicalHeaderKey(name) {
	case "Host":
		return h.host()
	case "Content-Length":
		if h.cl() < 0 {
			return ""
		}

		return strconv.FormatInt(h.cl(), 10)
	case "Transfer-Encoding":
		if len(h.te()) < 1 {
			return ""
		}

		return h.te()[0]
	default:
		return h.h.Get(name)
	}
}

// All returns all the values for header name. If the header does not exist it
// returns nil, false.
func (h *Header) All(name string) ([]string, bool) {
	switch http.CanonicalHeaderKey(name) {
	case "Host":
		if h.host() == "" {
			return nil, false
		}

		return []string{h.host()}, true
	case "Content-Length":
		if h.cl() <= 0 {
			return nil, false
		}

		return []string{strconv.FormatInt(h.cl(), 10)}, true
	case "Transfer-Encoding":
		if h.te() == nil {
			return nil, false
		}

		return h.te(), true
	default:
		vs, ok := h.h[http.CanonicalHeaderKey(name)]
		return vs, ok
	}
}

// Del deletes the header at name for the request or response.
func (h *Header) Del(name string) {
	switch http.CanonicalHeaderKey(name) {
	case "Host":
		h.setHost("")
	case "Content-Length":
		h.setCL(-1)
	case "Transfer-Encoding":
		h.setTE(nil)
	default:
		h.h.Del(name)
	}
}

// Map returns an http.Header that includes Host, Content-Length, and
// Transfer-Encoding.
func (h *Header) Map() http.Header {
	hm := make(http.Header)

	for k, vs := range h.h {
		hm[k] = vs
	}

	for _, k := range []string{
		"Host",
		"Content-Length",
		"Transfer-Encoding",
	} {
		vs, ok := h.All(k)
		if !ok {
			continue
		}

		hm[k] = vs
	}

	return hm
}
