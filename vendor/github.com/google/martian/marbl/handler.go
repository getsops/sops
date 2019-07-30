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

package marbl

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"

	"github.com/google/martian/v3/log"

	"golang.org/x/net/websocket"
)

// Handler exposes marbl logs over websockets.
type Handler struct {
	mu   sync.RWMutex
	subs map[string]chan<- []byte
}

// NewHandler instantiates a Handler with an empty set of subscriptions.
func NewHandler() *Handler {
	return &Handler{
		subs: make(map[string]chan<- []byte),
	}
}

// Write writes frames to all websocket subscribers and returns the number
// of bytes written and an error.
func (h *Handler) Write(b []byte) (int, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var wg sync.WaitGroup
	for id, framec := range h.subs {
		wg.Add(1)
		go func(id string, fc chan<- []byte) {
			defer wg.Done()
			select {
			case fc <- b:
			default:
				log.Errorf("logstream: buffer full for connection, dropping")
				go h.unsubscribe(id)
			}
		}(id, framec)
	}
	wg.Wait()

	return len(b), nil
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	websocket.Server{Handler: h.streamLogs}.ServeHTTP(rw, req)
}

func (h *Handler) streamLogs(conn *websocket.Conn) {
	defer conn.Close()

	id, err := newID()
	if err != nil {
		log.Errorf("logstream: failed to create ID: %v", err)
		return
	}
	framec := make(chan []byte, 16384)

	h.subscribe(id, framec)
	defer h.unsubscribe(id)

	for b := range framec {
		if err := websocket.Message.Send(conn, b); err != nil {
			log.Errorf("logstream: failed to send message: %v", err)
			return
		}
	}
}

func newID() (string, error) {
	src := make([]byte, 8)
	if _, err := rand.Read(src); err != nil {
		return "", err
	}

	return hex.EncodeToString(src), nil
}

func (h *Handler) unsubscribe(id string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if fc, ok := h.subs[id]; ok {
		close(fc)
		delete(h.subs, id)
	}
}

func (h *Handler) subscribe(id string, framec chan<- []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if fc, ok := h.subs[id]; ok {
		// TODO: Re-pick the id.
		log.Errorf("Resubscribing with ID: %v", id)
		// Close the channel for now so the websocket gets disconnected,
		// instead of silently failing.
		close(fc)
	}
	h.subs[id] = framec
}
