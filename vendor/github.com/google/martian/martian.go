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

// Package martian provides an HTTP/1.1 proxy with an API for configurable
// request and response modifiers.
package martian

import "net/http"

// RequestModifier is an interface that defines a request modifier that can be
// used by a proxy.
type RequestModifier interface {
	// ModifyRequest modifies the request.
	ModifyRequest(req *http.Request) error
}

// ResponseModifier is an interface that defines a response modifier that can
// be used by a proxy.
type ResponseModifier interface {
	// ModifyResponse modifies the response.
	ModifyResponse(res *http.Response) error
}

// RequestResponseModifier is an interface that is both a ResponseModifier and
// a RequestModifier.
type RequestResponseModifier interface {
	RequestModifier
	ResponseModifier
}

// RequestModifierFunc is an adapter for using a function with the given
// signature as a RequestModifier.
type RequestModifierFunc func(req *http.Request) error

// ResponseModifierFunc is an adapter for using a function with the given
// signature as a ResponseModifier.
type ResponseModifierFunc func(res *http.Response) error

// ModifyRequest modifies the request using the given function.
func (f RequestModifierFunc) ModifyRequest(req *http.Request) error {
	return f(req)
}

// ModifyResponse modifies the response using the given function.
func (f ResponseModifierFunc) ModifyResponse(res *http.Response) error {
	return f(res)
}
