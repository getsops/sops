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

// Package martianlog provides a Martian modifier that logs the request and response.
package martianlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/messageview"
	"github.com/google/martian/v3/parse"
)

// Logger is a modifier that logs requests and responses.
type Logger struct {
	log         func(line string)
	headersOnly bool
	decode      bool
}

type loggerJSON struct {
	Scope       []parse.ModifierType `json:"scope"`
	HeadersOnly bool                 `json:"headersOnly"`
	Decode      bool                 `json:"decode"`
}

func init() {
	parse.Register("log.Logger", loggerFromJSON)
}

// NewLogger returns a logger that logs requests and responses, optionally
// logging the body. Log function defaults to martian.Infof.
func NewLogger() *Logger {
	return &Logger{
		log: func(line string) {
			log.Infof(line)
		},
	}
}

// SetHeadersOnly sets whether to log the request/response body in the log.
func (l *Logger) SetHeadersOnly(headersOnly bool) {
	l.headersOnly = headersOnly
}

// SetDecode sets whether to decode the request/response body in the log.
func (l *Logger) SetDecode(decode bool) {
	l.decode = decode
}

// SetLogFunc sets the logging function for the logger.
func (l *Logger) SetLogFunc(logFunc func(line string)) {
	l.log = logFunc
}

// ModifyRequest logs the request, optionally including the body.
//
// The format logged is:
// --------------------------------------------------------------------------------
// Request to http://www.google.com/path?querystring
// --------------------------------------------------------------------------------
// GET /path?querystring HTTP/1.1
// Host: www.google.com
// Connection: close
// Other-Header: values
//
// request content
// --------------------------------------------------------------------------------
func (l *Logger) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	if ctx.SkippingLogging() {
		return nil
	}

	b := &bytes.Buffer{}

	fmt.Fprintln(b, "")
	fmt.Fprintln(b, strings.Repeat("-", 80))
	fmt.Fprintf(b, "Request to %s\n", req.URL)
	fmt.Fprintln(b, strings.Repeat("-", 80))

	mv := messageview.New()
	mv.SkipBody(l.headersOnly)
	if err := mv.SnapshotRequest(req); err != nil {
		return err
	}

	var opts []messageview.Option
	if l.decode {
		opts = append(opts, messageview.Decode())
	}

	r, err := mv.Reader(opts...)
	if err != nil {
		return err
	}

	io.Copy(b, r)

	fmt.Fprintln(b, "")
	fmt.Fprintln(b, strings.Repeat("-", 80))

	l.log(b.String())

	return nil
}

// ModifyResponse logs the response, optionally including the body.
//
// The format logged is:
// --------------------------------------------------------------------------------
// Response from http://www.google.com/path?querystring
// --------------------------------------------------------------------------------
// HTTP/1.1 200 OK
// Date: Tue, 15 Nov 1994 08:12:31 GMT
// Other-Header: values
//
// response content
// --------------------------------------------------------------------------------
func (l *Logger) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	if ctx.SkippingLogging() {
		return nil
	}

	b := &bytes.Buffer{}
	fmt.Fprintln(b, "")
	fmt.Fprintln(b, strings.Repeat("-", 80))
	fmt.Fprintf(b, "Response from %s\n", res.Request.URL)
	fmt.Fprintln(b, strings.Repeat("-", 80))

	mv := messageview.New()
	mv.SkipBody(l.headersOnly)
	if err := mv.SnapshotResponse(res); err != nil {
		return err
	}

	var opts []messageview.Option
	if l.decode {
		opts = append(opts, messageview.Decode())
	}

	r, err := mv.Reader(opts...)
	if err != nil {
		return err
	}

	io.Copy(b, r)

	fmt.Fprintln(b, "")
	fmt.Fprintln(b, strings.Repeat("-", 80))

	l.log(b.String())

	return nil
}

// loggerFromJSON builds a logger from JSON.
//
// Example JSON:
// {
//   "log.Logger": {
//     "scope": ["request", "response"],
//		 "headersOnly": true,
//		 "decode": true
//   }
// }
func loggerFromJSON(b []byte) (*parse.Result, error) {
	msg := &loggerJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	l := NewLogger()
	l.SetHeadersOnly(msg.HeadersOnly)
	l.SetDecode(msg.Decode)

	return parse.NewResult(l, msg.Scope)
}
