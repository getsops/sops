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
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"golang.org/x/net/websocket"
)

func TestStreamsInSentOrder(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	handler := NewHandler()
	go http.Serve(l, handler)

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s", l.Addr()), "", "http://localhost/")
	if err != nil {
		t.Fatalf("websocket.Dial(): got %v, want no error", err)
	}
	defer ws.Close()

	// Gives handler time to create the subscription channel.
	time.Sleep(200 * time.Millisecond)

	ws.SetDeadline(time.Now().Add(5 * time.Second))

	// server could still be in the processs of registering the client
	// no easy way to synchronize so we just wait a bit
	time.Sleep(300 * time.Millisecond)

	var iterations int64 = 5000
	go func() {
		for i := int64(0); i < iterations; i++ {
			hex := strconv.FormatInt(int64(i), 16)
			handler.Write([]byte(hex))
		}
	}()

	for i := int64(0); i < iterations; i++ {
		var bytes []byte
		err = websocket.Message.Receive(ws, &bytes)
		if err != nil {
			t.Fatalf("websocket.Conn.Read(): got %v, want no error", err)
		}
		parsed, err := strconv.ParseInt(string(bytes), 16, 64)
		if err != nil {
			t.Fatalf("strconv.ParseInt(): got %v, want no error", err)
		}
		if parsed != i {
			t.Fatalf("Messages arrived out of order, got %d want %d", parsed, i)
		}
	}
}

func TestUnreadsDontBlock(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	handler := NewHandler()
	go http.Serve(l, handler)

	ws, err := websocket.Dial(fmt.Sprintf("ws://%s", l.Addr()), "", "http://localhost/")
	if err != nil {
		t.Fatalf("websocket.Dial(): got %v, want no error", err)
	}
	defer ws.Close()

	// Gives handler time to create the subscription channel.
	time.Sleep(200 * time.Millisecond)

	bytes := make([]byte, 1024)
	_, err = rand.Read(bytes)
	if err != nil {
		t.Fatalf("rand.Read(): got %v, want no error", err)
	}
	// Purposely using more iterations than frame channel size.
	var iterations int64 = 50000
	for i := int64(0); i < iterations; i++ {
		to := doOrTimeout(3*time.Second, func() {
			handler.Write(bytes)
		})
		if to {
			t.Fatalf("handler.Write() Timed out")
		}
	}
}

func doOrTimeout(d time.Duration, f func()) bool {
	done := make(chan interface{})
	go func() {
		f()
		done <- 1
	}()
	select {
	case <-done:
		return false
	case <-time.After(d):
		return true
	}
}
