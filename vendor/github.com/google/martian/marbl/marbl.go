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

// Package marbl provides HTTP traffic logs streamed over websockets
// that can be added to any point within a Martian modifier tree.
// Marbl transmits HTTP logs that are serialized based on the following
// schema:
//
// Frame Header
// FrameType   uint8
// MessageType uint8
// ID		   [8]byte
// Payload	   HeaderFrame/DataFrame
//
// Header   Frame
// NameLen  uint32
// ValueLen uint32
// Name	    variable
// Value    variable
//
// Data Frame
// Index    uint32
// Terminal uint8
// Len      uint32
// Data     variable
package marbl

import (
	"io"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/proxyutil"
)

// MessageType incicates whether the message represents an HTTP request or response.
type MessageType uint8

const (
	// Unknown type of Message.
	Unknown MessageType = 0x0
	// Request indicates a message that contains an HTTP request.
	Request MessageType = 0x1
	// Response indicates a message that contains an HTTP response.
	Response MessageType = 0x2
)

// FrameType indicates whether the frame contains a Header or Data.
type FrameType uint8

const (
	// UnknownFrame indicates an unknown type of Frame.
	UnknownFrame FrameType = 0x0
	// HeaderFrame indicates a frame that contains a header.
	HeaderFrame FrameType = 0x1
	// DataFrame indicates a frame that contains the payload, usually the body.
	DataFrame FrameType = 0x2
)

// Stream writes logs of requests and responses to a writer.
type Stream struct {
	w      io.Writer
	framec chan []byte
	closec chan struct{}
}

// NewStream initializes a Stream with an io.Writer to log requests and
// responses to. Upon construction, a goroutine is started that listens for frames
// and writes them to w.
func NewStream(w io.Writer) *Stream {
	s := &Stream{
		w:      w,
		framec: make(chan []byte),
		closec: make(chan struct{}),
	}

	go s.loop()

	return s
}

func (s *Stream) loop() {
	for {
		select {
		case f := <-s.framec:
			_, err := s.w.Write(f)
			if err != nil {
				log.Errorf("martian: Error while writing frame")
			}
		case <-s.closec:
			return
		}
	}
}

// Close signals Stream to stop listening for frames in the log loop and stop writing logs.
func (s *Stream) Close() error {
	s.closec <- struct{}{}
	close(s.closec)

	return nil
}

func newFrame(id string, ft FrameType, mt MessageType, plen uint32) []byte {
	f := make([]byte, 0, 10+plen)
	f = append(f, byte(ft), byte(mt))
	f = append(f, id[:8]...)

	return f
}

func (s *Stream) sendHeader(id string, mt MessageType, key, value string) {
	kl := uint32(len(key))
	vl := uint32(len(value))

	f := newFrame(id, HeaderFrame, mt, 64+kl+vl)
	f = append(f, byte(kl>>24), byte(kl>>16), byte(kl>>8), byte(kl))
	f = append(f, byte(vl>>24), byte(vl>>16), byte(vl>>8), byte(vl))
	f = append(f, key[:kl]...)
	f = append(f, value[:vl]...)

	s.framec <- f
}

func (s *Stream) sendData(id string, mt MessageType, i uint32, terminal bool, b []byte, bl int) {
	var ti uint8
	if terminal {
		ti = 1
	}

	f := newFrame(id, DataFrame, mt, 72+uint32(bl))
	f = append(f, byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
	f = append(f, byte(ti))
	f = append(f, byte(bl>>24), byte(bl>>16), byte(bl>>8), byte(bl))
	f = append(f, b[:bl]...)

	s.framec <- f
}

// LogRequest writes an http.Request to Stream with an id unique for the request / response pair.
func (s *Stream) LogRequest(id string, req *http.Request) error {
	s.sendHeader(id, Request, ":method", req.Method)
	s.sendHeader(id, Request, ":scheme", req.URL.Scheme)
	s.sendHeader(id, Request, ":authority", req.URL.Host)
	s.sendHeader(id, Request, ":path", req.URL.EscapedPath())
	s.sendHeader(id, Request, ":query", req.URL.RawQuery)
	s.sendHeader(id, Request, ":proto", req.Proto)
	s.sendHeader(id, Request, ":remote", req.RemoteAddr)
	ts := strconv.FormatInt(time.Now().UnixNano()/1000/1000, 10)
	s.sendHeader(id, Request, ":timestamp", ts)

	ctx := martian.NewContext(req)
	if ctx.IsAPIRequest() {
		s.sendHeader(id, Request, ":api", "true")
	}

	h := proxyutil.RequestHeader(req)

	for k, vs := range h.Map() {
		for _, v := range vs {
			s.sendHeader(id, Request, k, v)
		}
	}

	req.Body = &bodyLogger{
		s:    s,
		id:   id,
		mt:   Request,
		body: req.Body,
	}

	return nil
}

// LogResponse writes an http.Response to Stream with an id unique for the request / response pair.
func (s *Stream) LogResponse(id string, res *http.Response) error {
	s.sendHeader(id, Response, ":proto", res.Proto)
	s.sendHeader(id, Response, ":status", strconv.Itoa(res.StatusCode))
	s.sendHeader(id, Response, ":reason", res.Status)
	ts := strconv.FormatInt(time.Now().UnixNano()/1000/1000, 10)
	s.sendHeader(id, Response, ":timestamp", ts)

	ctx := martian.NewContext(res.Request)
	if ctx.IsAPIRequest() {
		s.sendHeader(id, Response, ":api", "true")
	}

	h := proxyutil.ResponseHeader(res)

	for k, vs := range h.Map() {
		for _, v := range vs {
			s.sendHeader(id, Response, k, v)
		}
	}

	res.Body = &bodyLogger{
		s:    s,
		id:   id,
		mt:   Response,
		body: res.Body,
	}

	return nil
}

type bodyLogger struct {
	index uint32 // atomic
	s     *Stream
	id    string
	mt    MessageType
	body  io.ReadCloser
}

// Read implements the standard Reader interface. Read reads the bytes of the body
// and returns the number of bytes read and an error.
func (bl *bodyLogger) Read(b []byte) (int, error) {
	var terminal bool

	n, err := bl.body.Read(b)
	if err == io.EOF {
		terminal = true
	}

	bl.s.sendData(bl.id, bl.mt, atomic.AddUint32(&bl.index, 1)-1, terminal, b, n)

	return n, err
}

// Close closes the bodyLogger.
func (bl *bodyLogger) Close() error {
	return bl.body.Close()
}
