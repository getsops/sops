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

// Package api contains a forwarder to route system HTTP requests to the
// local API server.
package api

import (
	"fmt"
	"net/http"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
)

// Forwarder is a request modifier that routes the request to the API server and
// marks the request for skipped logging.
type Forwarder struct {
	host string
	port int
}

// NewForwarder returns a Forwarder that rewrites requests to host.
func NewForwarder(host string, port int) *Forwarder {
	return &Forwarder{
		host: host,
		port: port,
	}
}

// ModifyRequest forwards the request to the local API server running at f.port,
// downgrades the scheme to http and marks the request context for skipped logging.
// API requests are marked for skipping the roundtrip.
func (f *Forwarder) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	ctx.APIRequest()
	ctx.SkipLogging()

	in := req.URL.String()
	req.URL.Scheme = "http"
	req.URL.Host = fmt.Sprintf("%s:%d", "localhost", f.port)
	out := req.URL.String()
	log.Infof("api.Forwarder: forwarding %s to %s", in, out)

	return nil
}
