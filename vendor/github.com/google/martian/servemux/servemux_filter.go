// Copyright 2016 Google Inc. All rights reserved.
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

// Package servemux contains a filter that executes modifiers when there is a
// pattern match in a mux.
package servemux

import (
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/filter"
)

var noop = martian.Noop("mux.Filter")

// Filter is a modifier that executes mod if a pattern is matched in mux.
type Filter struct {
	*filter.Filter
}

// NewFilter constructs a filter that applies the modifier when the request
// url matches a pattern in mux. If no mux is provided, the request is evaluated
// against patterns in http.DefaultServeMux.
func NewFilter(mux *http.ServeMux) *Filter {
	if mux == nil {
		mux = http.DefaultServeMux
	}

	m := NewMatcher(mux)
	f := filter.New()
	f.SetRequestCondition(m)
	f.SetResponseCondition(m)
	return &Filter{Filter: f}
}
