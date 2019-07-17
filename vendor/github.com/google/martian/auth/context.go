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

package auth

import (
	"sync"

	"github.com/google/martian/v3"
)

const key = "auth.Context"

// Context contains authentication information.
type Context struct {
	mu  sync.RWMutex
	id  string
	err error
}

// FromContext retrieves the auth.Context from the session.
func FromContext(ctx *martian.Context) *Context {
	if v, ok := ctx.Session().Get(key); ok {
		return v.(*Context)
	}

	actx := &Context{}
	ctx.Session().Set(key, actx)

	return actx
}

// ID returns the ID.
func (ctx *Context) ID() string {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	return ctx.id
}

// SetID sets the ID.
func (ctx *Context) SetID(id string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.err = nil

	if id == "" {
		return
	}

	ctx.id = id
}

// SetError sets the error and resets the ID.
func (ctx *Context) SetError(err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	ctx.id = ""
	ctx.err = err
}

// Error returns the error.
func (ctx *Context) Error() error {
	ctx.mu.RLock()
	defer ctx.mu.RUnlock()

	return ctx.err
}
