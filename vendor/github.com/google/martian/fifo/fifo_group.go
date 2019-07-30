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

// Package fifo provides Group, which is a list of modifiers that are executed
// consecutively. By default, when an error is returned by a modifier, the
// execution of the modifiers is halted, and the error is returned. Optionally,
// when errror aggregation is enabled (by calling SetAggretateErrors(true)), modifier
// execution is not halted, and errors are aggretated and returned after all
// modifiers have been executed.
package fifo

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/verify"
)

// Group is a martian.RequestResponseModifier that maintains lists of
// request and response modifiers executed on a first-in, first-out basis.
type Group struct {
	reqmu   sync.RWMutex
	reqmods []martian.RequestModifier

	resmu   sync.RWMutex
	resmods []martian.ResponseModifier

	aggregateErrors bool
}

type groupJSON struct {
	Modifiers       []json.RawMessage    `json:"modifiers"`
	Scope           []parse.ModifierType `json:"scope"`
	AggregateErrors bool                 `json:"aggregateErrors"`
}

func init() {
	parse.Register("fifo.Group", groupFromJSON)
}

// NewGroup returns a modifier group.
func NewGroup() *Group {
	return &Group{}
}

// SetAggregateErrors sets the error behavior for the Group. When true, the Group will
// continue to execute consecutive modifiers when a modifier in the group encounters an
// error. The Group will then return all errors returned by each modifier after all
// modifiers have been executed.  When false, if an error is returned by a modifier, the
// error is returned by ModifyRequest/Response and no further modifiers are run.
// By default, error aggregation is disabled.
func (g *Group) SetAggregateErrors(aggerr bool) {
	g.aggregateErrors = aggerr
}

// AddRequestModifier adds a RequestModifier to the group's list of request modifiers.
func (g *Group) AddRequestModifier(reqmod martian.RequestModifier) {
	g.reqmu.Lock()
	defer g.reqmu.Unlock()

	g.reqmods = append(g.reqmods, reqmod)
}

// AddResponseModifier adds a ResponseModifier to the group's list of response modifiers.
func (g *Group) AddResponseModifier(resmod martian.ResponseModifier) {
	g.resmu.Lock()
	defer g.resmu.Unlock()

	g.resmods = append(g.resmods, resmod)
}

// ModifyRequest modifies the request. By default, aggregateErrors is false; if an error is
// returned by a RequestModifier the error is returned and no further modifiers are run. When
// aggregateErrors is set to true, the errors returned by each modifier in the group are
// aggregated.
func (g *Group) ModifyRequest(req *http.Request) error {
	log.Debugf("fifo.ModifyRequest: %s", req.URL)
	g.reqmu.RLock()
	defer g.reqmu.RUnlock()

	merr := martian.NewMultiError()

	for _, reqmod := range g.reqmods {
		if err := reqmod.ModifyRequest(req); err != nil {
			if g.aggregateErrors {
				merr.Add(err)
				continue
			}

			return err
		}
	}

	if merr.Empty() {
		return nil
	}

	return merr
}

// ModifyResponse modifies the request. By default, aggregateErrors is false; if an error is
// returned by a RequestModifier the error is returned and no further modifiers are run. When
// aggregateErrors is set to true, the errors returned by each modifier in the group are
// aggregated.
func (g *Group) ModifyResponse(res *http.Response) error {
	requ := ""
	if res.Request != nil {
		requ = res.Request.URL.String()
		log.Debugf("fifo.ModifyResponse: %s", requ)
	}
	g.resmu.RLock()
	defer g.resmu.RUnlock()

	merr := martian.NewMultiError()

	for _, resmod := range g.resmods {
		if err := resmod.ModifyResponse(res); err != nil {
			if g.aggregateErrors {
				merr.Add(err)
				continue
			}

			return err
		}
	}

	if merr.Empty() {
		return nil
	}

	return merr
}

// VerifyRequests returns a MultiError containing all the
// verification errors returned by request verifiers.
func (g *Group) VerifyRequests() error {
	log.Debugf("fifo.VerifyRequests()")
	g.reqmu.Lock()
	defer g.reqmu.Unlock()

	merr := martian.NewMultiError()
	for _, reqmod := range g.reqmods {
		reqv, ok := reqmod.(verify.RequestVerifier)
		if !ok {
			continue
		}

		if err := reqv.VerifyRequests(); err != nil {
			merr.Add(err)
		}
	}

	if merr.Empty() {
		return nil
	}

	return merr
}

// VerifyResponses returns a MultiError containing all the
// verification errors returned by response verifiers.
func (g *Group) VerifyResponses() error {
	log.Debugf("fifo.VerifyResponses()")
	g.resmu.Lock()
	defer g.resmu.Unlock()

	merr := martian.NewMultiError()
	for _, resmod := range g.resmods {
		resv, ok := resmod.(verify.ResponseVerifier)
		if !ok {
			continue
		}

		if err := resv.VerifyResponses(); err != nil {
			merr.Add(err)
		}
	}

	if merr.Empty() {
		return nil
	}

	return merr
}

// ResetRequestVerifications resets the state of the contained request verifiers.
func (g *Group) ResetRequestVerifications() {
	log.Debugf("fifo.ResetRequestVerifications()")
	g.reqmu.Lock()
	defer g.reqmu.Unlock()

	for _, reqmod := range g.reqmods {
		if reqv, ok := reqmod.(verify.RequestVerifier); ok {
			reqv.ResetRequestVerifications()
		}
	}
}

// ResetResponseVerifications resets the state of the contained request verifiers.
func (g *Group) ResetResponseVerifications() {
	log.Debugf("fifo.ResetResponseVerifications()")
	g.resmu.Lock()
	defer g.resmu.Unlock()

	for _, resmod := range g.resmods {
		if resv, ok := resmod.(verify.ResponseVerifier); ok {
			resv.ResetResponseVerifications()
		}
	}
}

// groupFromJSON builds a fifo.Group from JSON.
//
// Example JSON:
// {
//   "fifo.Group" : {
//     "scope": ["request", "result"],
//     "modifiers": [
//       { ... },
//       { ... },
//     ]
//   }
// }
func groupFromJSON(b []byte) (*parse.Result, error) {
	msg := &groupJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	g := NewGroup()
	if msg.AggregateErrors {
		g.SetAggregateErrors(true)
	}

	for _, m := range msg.Modifiers {
		r, err := parse.FromJSON(m)
		if err != nil {
			return nil, err
		}

		reqmod := r.RequestModifier()
		if reqmod != nil {
			g.AddRequestModifier(reqmod)
		}

		resmod := r.ResponseModifier()
		if resmod != nil {
			g.AddResponseModifier(resmod)
		}
	}

	return parse.NewResult(g, msg.Scope)
}
