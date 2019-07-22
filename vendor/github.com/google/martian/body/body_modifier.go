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

// Package body allows for the replacement of message body on responses.
package body

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/parse"
)

func init() {
	parse.Register("body.Modifier", modifierFromJSON)
}

// Modifier substitutes the body on an HTTP response.
type Modifier struct {
	contentType string
	body        []byte
	boundary    string
}

type modifierJSON struct {
	ContentType string               `json:"contentType"`
	Body        []byte               `json:"body"` // Body is expected to be a Base64 encoded string.
	Scope       []parse.ModifierType `json:"scope"`
}

// NewModifier constructs and returns a body.Modifier.
func NewModifier(b []byte, contentType string) *Modifier {
	log.Debugf("body.NewModifier: len(b): %d, contentType %s", len(b), contentType)
	return &Modifier{
		contentType: contentType,
		body:        b,
		boundary:    randomBoundary(),
	}
}

// modifierFromJSON takes a JSON message as a byte slice and returns a
// body.Modifier and an error.
//
// Example JSON Configuration message:
// {
//   "scope": ["request", "response"],
//   "contentType": "text/plain",
//   "body": "c29tZSBkYXRhIHdpdGggACBhbmQg77u/" // Base64 encoded body
// }
func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &modifierJSON{}
	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	mod := NewModifier(msg.Body, msg.ContentType)
	return parse.NewResult(mod, msg.Scope)
}

// ModifyRequest sets the Content-Type header and overrides the request body.
func (m *Modifier) ModifyRequest(req *http.Request) error {
	log.Debugf("body.ModifyRequest: request: %s", req.URL)
	req.Body.Close()

	req.Header.Set("Content-Type", m.contentType)

	// Reset the Content-Encoding since we know that the new body isn't encoded.
	req.Header.Del("Content-Encoding")

	req.ContentLength = int64(len(m.body))
	req.Body = ioutil.NopCloser(bytes.NewReader(m.body))

	return nil
}

// SetBoundary set the boundary string used for multipart range responses.
func (m *Modifier) SetBoundary(boundary string) {
	m.boundary = boundary
}

// ModifyResponse sets the Content-Type header and overrides the response body.
func (m *Modifier) ModifyResponse(res *http.Response) error {
	log.Debugf("body.ModifyResponse: request: %s", res.Request.URL)
	// Replace the existing body, close it first.
	res.Body.Close()

	res.Header.Set("Content-Type", m.contentType)

	// Reset the Content-Encoding since we know that the new body isn't encoded.
	res.Header.Del("Content-Encoding")

	// If no range request header is present, return the body as the response body.
	if res.Request.Header.Get("Range") == "" {
		res.ContentLength = int64(len(m.body))
		res.Body = ioutil.NopCloser(bytes.NewReader(m.body))

		return nil
	}

	rh := res.Request.Header.Get("Range")
	rh = strings.ToLower(rh)
	sranges := strings.Split(strings.TrimLeft(rh, "bytes="), ",")
	var ranges [][]int
	for _, rng := range sranges {
		if strings.HasSuffix(rng, "-") {
			rng = fmt.Sprintf("%s%d", rng, len(m.body)-1)
		}

		rs := strings.Split(rng, "-")
		if len(rs) != 2 {
			res.StatusCode = http.StatusRequestedRangeNotSatisfiable
			return nil
		}
		start, err := strconv.Atoi(strings.TrimSpace(rs[0]))
		if err != nil {
			return err
		}

		end, err := strconv.Atoi(strings.TrimSpace(rs[1]))
		if err != nil {
			return err
		}

		if start > end {
			res.StatusCode = http.StatusRequestedRangeNotSatisfiable
			return nil
		}

		ranges = append(ranges, []int{start, end})
	}

	// Range request.
	res.StatusCode = http.StatusPartialContent

	// Single range request.
	if len(ranges) == 1 {
		start := ranges[0][0]
		end := ranges[0][1]
		seg := m.body[start : end+1]
		res.ContentLength = int64(len(seg))
		res.Body = ioutil.NopCloser(bytes.NewReader(seg))
		res.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(m.body)))

		return nil
	}

	// Multipart range request.
	var mpbody bytes.Buffer
	mpw := multipart.NewWriter(&mpbody)
	mpw.SetBoundary(m.boundary)

	for _, rng := range ranges {
		start, end := rng[0], rng[1]
		mimeh := make(textproto.MIMEHeader)
		mimeh.Set("Content-Type", m.contentType)
		mimeh.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, len(m.body)))

		seg := m.body[start : end+1]

		pw, err := mpw.CreatePart(mimeh)
		if err != nil {
			return err
		}

		if _, err := pw.Write(seg); err != nil {
			return err
		}
	}
	mpw.Close()

	res.ContentLength = int64(len(mpbody.Bytes()))
	res.Body = ioutil.NopCloser(bytes.NewReader(mpbody.Bytes()))
	res.Header.Set("Content-Type", fmt.Sprintf("multipart/byteranges; boundary=%s", m.boundary))

	return nil
}

// randomBoundary generates a 30 character string for boundaries for mulipart range
// requests. This func panics if io.Readfull fails.
// Borrowed from: https://golang.org/src/mime/multipart/writer.go?#L73
func randomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}
