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

// Package auth provides filtering support for a martian.Proxy based on auth
// ID.
package auth

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/google/martian/v3"
)

// Filter filters RequestModifiers and ResponseModifiers by auth ID.
type Filter struct {
	authRequired bool

	mu      sync.RWMutex
	reqmods map[string]martian.RequestModifier
	resmods map[string]martian.ResponseModifier
}

// NewFilter returns a new auth.Filter.
func NewFilter() *Filter {
	return &Filter{
		reqmods: make(map[string]martian.RequestModifier),
		resmods: make(map[string]martian.ResponseModifier),
	}
}

// SetAuthRequired determines whether the auth ID must have an associated
// RequestModifier or ResponseModifier. If true, it will set auth error.
func (f *Filter) SetAuthRequired(required bool) {
	f.authRequired = required
}

// SetRequestModifier sets the RequestModifier for the given ID. It will
// overwrite any existing modifier with the same ID.
func (f *Filter) SetRequestModifier(id string, reqmod martian.RequestModifier) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if reqmod != nil {
		f.reqmods[id] = reqmod
	} else {
		delete(f.reqmods, id)
	}

	return nil
}

// SetResponseModifier sets the ResponseModifier for the given ID. It will
// overwrite any existing modifier with the same ID.
func (f *Filter) SetResponseModifier(id string, resmod martian.ResponseModifier) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if resmod != nil {
		f.resmods[id] = resmod
	} else {
		delete(f.resmods, id)
	}

	return nil
}

// RequestModifier retrieves the RequestModifier for the given ID. Returns nil
// if no modifier exists for the given ID.
func (f *Filter) RequestModifier(id string) martian.RequestModifier {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.reqmods[id]
}

// ResponseModifier retrieves the ResponseModifier for the given ID. Returns nil
// if no modifier exists for the given ID.
func (f *Filter) ResponseModifier(id string) martian.ResponseModifier {
	f.mu.RLock()
	defer f.mu.RUnlock()

	return f.resmods[id]
}

// ModifyRequest runs the RequestModifier for the associated auth ID. If no
// modifier is found for auth ID then auth error is set.
func (f *Filter) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	actx := FromContext(ctx)

	if reqmod, ok := f.reqmods[actx.ID()]; ok {
		return reqmod.ModifyRequest(req)
	}

	if err := f.requireKnownAuth(actx.ID()); err != nil {
		actx.SetError(err)
	}

	return nil
}

// ModifyResponse runs the ResponseModifier for the associated auth ID. If no
// modifier is found for the auth ID then the auth error is set.
func (f *Filter) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	actx := FromContext(ctx)

	if resmod, ok := f.resmods[actx.ID()]; ok {
		return resmod.ModifyResponse(res)
	}

	if err := f.requireKnownAuth(actx.ID()); err != nil {
		actx.SetError(err)
	}

	return nil
}

func (f *Filter) requireKnownAuth(id string) error {
	_, reqok := f.reqmods[id]
	_, resok := f.resmods[id]

	if !reqok && !resok && f.authRequired {
		return fmt.Errorf("auth: unrecognized credentials: %s", id)
	}

	return nil
}
