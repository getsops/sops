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

// Package httpspec provides a modifier stack that has been preconfigured to
// provide spec-compliant HTTP proxy behavior.
//
// Related: https://www.mnot.net/blog/2011/07/11/what_proxies_must_do
package httpspec

import (
	"github.com/google/martian/v3/fifo"
	"github.com/google/martian/v3/header"
)

// NewStack returns a martian modifier stack that handles ensuring proper proxy
// behavior, in addition to a fifo.Group that can be used to add additional
// modifiers within the stack.
func NewStack(via string) (outer *fifo.Group, inner *fifo.Group) {
	outer = fifo.NewGroup()

	hbhm := header.NewHopByHopModifier()
	outer.AddRequestModifier(hbhm)
	outer.AddRequestModifier(header.NewForwardedModifier())
	outer.AddRequestModifier(header.NewBadFramingModifier())

	vm := header.NewViaModifier(via)
	outer.AddRequestModifier(vm)

	inner = fifo.NewGroup()
	outer.AddRequestModifier(inner)
	outer.AddResponseModifier(inner)

	outer.AddResponseModifier(vm)
	outer.AddResponseModifier(hbhm)

	return outer, inner
}
