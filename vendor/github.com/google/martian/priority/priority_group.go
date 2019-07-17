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

// Package priority allows grouping modifiers and applying them in priority order.
package priority

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
)

var (
	// ErrModifierNotFound is the error returned when attempting to remove a
	// modifier when the modifier does not exist in the group.
	ErrModifierNotFound = errors.New("modifier not found in group")
)

// priorityRequestModifier is a request modifier with a priority.
type priorityRequestModifier struct {
	reqmod   martian.RequestModifier
	priority int64
}

// priorityResponseModifier is a response modifier with a priority.
type priorityResponseModifier struct {
	resmod   martian.ResponseModifier
	priority int64
}

// Group is a group of request and response modifiers ordered by their priority.
type Group struct {
	reqmu   sync.RWMutex
	reqmods []*priorityRequestModifier

	resmu   sync.RWMutex
	resmods []*priorityResponseModifier
}

type groupJSON struct {
	Modifiers []modifierJSON       `json:"modifiers"`
	Scope     []parse.ModifierType `json:"scope"`
}

type modifierJSON struct {
	Priority int64           `json:"priority"`
	Modifier json.RawMessage `json:"modifier"`
}

func init() {
	parse.Register("priority.Group", groupFromJSON)
}

// NewGroup returns a priority group.
func NewGroup() *Group {
	return &Group{}
}

// AddRequestModifier adds a RequestModifier with the given priority.
//
// If a modifier is added with a priority that is equal to an existing priority
// the newer modifier will be added before the existing modifier in the chain.
func (pg *Group) AddRequestModifier(reqmod martian.RequestModifier, priority int64) {
	pg.reqmu.Lock()
	defer pg.reqmu.Unlock()

	preqmod := &priorityRequestModifier{
		reqmod:   reqmod,
		priority: priority,
	}

	for i, m := range pg.reqmods {
		if preqmod.priority >= m.priority {
			pg.reqmods = append(pg.reqmods, nil)
			copy(pg.reqmods[i+1:], pg.reqmods[i:])
			pg.reqmods[i] = preqmod
			return
		}
	}

	// Either this is the first modifier in the list, or the priority is less
	// than all existing modifiers.
	pg.reqmods = append(pg.reqmods, preqmod)
}

// RemoveRequestModifier removes the the highest priority given RequestModifier.
// Returns ErrModifierNotFound if the given modifier does not exist in the group.
func (pg *Group) RemoveRequestModifier(reqmod martian.RequestModifier) error {
	pg.reqmu.Lock()
	defer pg.reqmu.Unlock()

	for i, m := range pg.reqmods {
		if m.reqmod == reqmod {
			copy(pg.reqmods[i:], pg.reqmods[i+1:])
			pg.reqmods[len(pg.reqmods)-1] = nil
			pg.reqmods = pg.reqmods[:len(pg.reqmods)-1]
			return nil
		}
	}

	return ErrModifierNotFound
}

// AddResponseModifier adds a ResponseModifier with the given priority.
//
// If a modifier is added with a priority that is equal to an existing priority
// the newer modifier will be added before the existing modifier in the chain.
func (pg *Group) AddResponseModifier(resmod martian.ResponseModifier, priority int64) {
	pg.resmu.Lock()
	defer pg.resmu.Unlock()

	presmod := &priorityResponseModifier{
		resmod:   resmod,
		priority: priority,
	}

	for i, m := range pg.resmods {
		if presmod.priority >= m.priority {
			pg.resmods = append(pg.resmods, nil)
			copy(pg.resmods[i+1:], pg.resmods[i:])
			pg.resmods[i] = presmod
			return
		}
	}

	// Either this is the first modifier in the list, or the priority is less
	// than all existing modifiers.
	pg.resmods = append(pg.resmods, presmod)
}

// RemoveResponseModifier removes the the highest priority given ResponseModifier.
// Returns ErrModifierNotFound if the given modifier does not exist in the group.
func (pg *Group) RemoveResponseModifier(resmod martian.ResponseModifier) error {
	pg.resmu.Lock()
	defer pg.resmu.Unlock()

	for i, m := range pg.resmods {
		if m.resmod == resmod {
			copy(pg.resmods[i:], pg.resmods[i+1:])
			pg.resmods[len(pg.resmods)-1] = nil
			pg.resmods = pg.resmods[:len(pg.resmods)-1]
			return nil
		}
	}

	return ErrModifierNotFound
}

// ModifyRequest modifies the request. Modifiers are run in descending order of
// their priority. If an error is returned by a RequestModifier the error is
// returned and no further modifiers are run.
func (pg *Group) ModifyRequest(req *http.Request) error {
	pg.reqmu.RLock()
	defer pg.reqmu.RUnlock()

	for _, m := range pg.reqmods {
		if err := m.reqmod.ModifyRequest(req); err != nil {
			return err
		}
	}

	return nil
}

// ModifyResponse modifies the response. Modifiers are run in descending order
// of their priority. If an error is returned by a ResponseModifier the error
// is returned and no further modifiers are run.
func (pg *Group) ModifyResponse(res *http.Response) error {
	pg.resmu.RLock()
	defer pg.resmu.RUnlock()

	for _, m := range pg.resmods {
		if err := m.resmod.ModifyResponse(res); err != nil {
			return err
		}
	}

	return nil
}

// groupFromJSON builds a priority.Group from JSON.
//
// Example JSON:
// {
//   "priority.Group": {
//     "scope": ["request", "response"],
//     "modifiers": [
//       {
//         "priority": 100, // Will run first.
//         "modifier": { ... },
//       },
//       {
//         "priority": 0, // Will run last.
//         "modifier": { ... },
//       }
//     ]
//   }
// }
func groupFromJSON(b []byte) (*parse.Result, error) {
	msg := &groupJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	pg := NewGroup()

	for _, m := range msg.Modifiers {
		r, err := parse.FromJSON(m.Modifier)
		if err != nil {
			return nil, err
		}

		reqmod := r.RequestModifier()
		if reqmod != nil {
			pg.AddRequestModifier(reqmod, m.Priority)
		}

		resmod := r.ResponseModifier()
		if resmod != nil {
			pg.AddResponseModifier(resmod, m.Priority)
		}
	}

	return parse.NewResult(pg, msg.Scope)
}
