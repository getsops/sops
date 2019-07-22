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

// Package messageview provides no-op snapshots for HTTP requests and
// responses.
package messageview

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
)

// MessageView is a static view of an HTTP request or response.
type MessageView struct {
	message       []byte
	cts           []string
	chunked       bool
	skipBody      bool
	compress      string
	bodyoffset    int64
	traileroffset int64
}

type config struct {
	decode bool
}

// Option is a configuration option for a MessageView.
type Option func(*config)

// Decode sets an option to decode the message body for logging purposes.
func Decode() Option {
	return func(c *config) {
		c.decode = true
	}
}

// New returns a new MessageView.
func New() *MessageView {
	return &MessageView{}
}

// SkipBody will skip reading the body when the view is loaded with a request
// or response.
func (mv *MessageView) SkipBody(skipBody bool) {
	mv.skipBody = skipBody
}

// SkipBodyUnlessContentType will skip reading the body unless the
// Content-Type matches one in cts.
func (mv *MessageView) SkipBodyUnlessContentType(cts ...string) {
	mv.skipBody = true
	mv.cts = cts
}

// SnapshotRequest reads the request into the MessageView. If mv.skipBody is false
// it will also read the body into memory and replace the existing body with
// the in-memory copy. This method is semantically a no-op.
func (mv *MessageView) SnapshotRequest(req *http.Request) error {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "%s %s HTTP/%d.%d\r\n", req.Method,
		req.URL, req.ProtoMajor, req.ProtoMinor)

	if req.Host != "" {
		fmt.Fprintf(buf, "Host: %s\r\n", req.Host)
	}

	if tec := len(req.TransferEncoding); tec > 0 {
		mv.chunked = req.TransferEncoding[tec-1] == "chunked"
		fmt.Fprintf(buf, "Transfer-Encoding: %s\r\n", strings.Join(req.TransferEncoding, ", "))
	}
	if !mv.chunked && req.ContentLength >= 0 {
		fmt.Fprintf(buf, "Content-Length: %d\r\n", req.ContentLength)
	}

	mv.compress = req.Header.Get("Content-Encoding")

	req.Header.WriteSubset(buf, map[string]bool{
		"Host":              true,
		"Content-Length":    true,
		"Transfer-Encoding": true,
	})

	fmt.Fprint(buf, "\r\n")

	mv.bodyoffset = int64(buf.Len())
	mv.traileroffset = int64(buf.Len())

	ct := req.Header.Get("Content-Type")
	if mv.skipBody && !mv.matchContentType(ct) || req.Body == nil {
		mv.message = buf.Bytes()
		return nil
	}

	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	req.Body.Close()

	if mv.chunked {
		cw := httputil.NewChunkedWriter(buf)
		cw.Write(data)
		cw.Close()
	} else {
		buf.Write(data)
	}

	mv.traileroffset = int64(buf.Len())

	req.Body = ioutil.NopCloser(bytes.NewReader(data))

	if req.Trailer != nil {
		req.Trailer.Write(buf)
	} else if mv.chunked {
		fmt.Fprint(buf, "\r\n")
	}

	mv.message = buf.Bytes()

	return nil
}

// SnapshotResponse reads the response into the MessageView. If mv.headersOnly
// is false it will also read the body into memory and replace the existing
// body with the in-memory copy. This method is semantically a no-op.
func (mv *MessageView) SnapshotResponse(res *http.Response) error {
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "HTTP/%d.%d %s\r\n", res.ProtoMajor, res.ProtoMinor, res.Status)

	if tec := len(res.TransferEncoding); tec > 0 {
		mv.chunked = res.TransferEncoding[tec-1] == "chunked"
		fmt.Fprintf(buf, "Transfer-Encoding: %s\r\n", strings.Join(res.TransferEncoding, ", "))
	}
	if !mv.chunked && res.ContentLength >= 0 {
		fmt.Fprintf(buf, "Content-Length: %d\r\n", res.ContentLength)
	}

	mv.compress = res.Header.Get("Content-Encoding")
	// Do not uncompress if we have don't have the full contents.
	if res.StatusCode == http.StatusNoContent || res.StatusCode == http.StatusPartialContent {
		mv.compress = ""
	}

	res.Header.WriteSubset(buf, map[string]bool{
		"Content-Length":    true,
		"Transfer-Encoding": true,
	})

	fmt.Fprint(buf, "\r\n")

	mv.bodyoffset = int64(buf.Len())
	mv.traileroffset = int64(buf.Len())

	ct := res.Header.Get("Content-Type")
	if mv.skipBody && !mv.matchContentType(ct) || res.Body == nil {
		mv.message = buf.Bytes()
		return nil
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	if mv.chunked {
		cw := httputil.NewChunkedWriter(buf)
		cw.Write(data)
		cw.Close()
	} else {
		buf.Write(data)
	}

	mv.traileroffset = int64(buf.Len())

	res.Body = ioutil.NopCloser(bytes.NewReader(data))

	if res.Trailer != nil {
		res.Trailer.Write(buf)
	} else if mv.chunked {
		fmt.Fprint(buf, "\r\n")
	}

	mv.message = buf.Bytes()

	return nil
}

// Reader returns the an io.ReadCloser that reads the full HTTP message.
func (mv *MessageView) Reader(opts ...Option) (io.ReadCloser, error) {
	hr := mv.HeaderReader()
	br, err := mv.BodyReader(opts...)
	if err != nil {
		return nil, err
	}
	tr := mv.TrailerReader()

	return struct {
		io.Reader
		io.Closer
	}{
		Reader: io.MultiReader(hr, br, tr),
		Closer: br,
	}, nil
}

// HeaderReader returns an io.Reader that reads the HTTP Status-Line or
// HTTP Request-Line and headers.
func (mv *MessageView) HeaderReader() io.Reader {
	r := bytes.NewReader(mv.message)
	return io.NewSectionReader(r, 0, mv.bodyoffset)
}

// BodyReader returns an io.ReadCloser that reads the HTTP request or response
// body. If mv.skipBody was set the reader will immediately return io.EOF.
//
// If the Decode option is passed the body will be unchunked if
// Transfer-Encoding is set to "chunked", and will decode the following
// Content-Encodings: gzip, deflate.
func (mv *MessageView) BodyReader(opts ...Option) (io.ReadCloser, error) {
	var r io.Reader

	conf := &config{}
	for _, o := range opts {
		o(conf)
	}

	br := bytes.NewReader(mv.message)
	r = io.NewSectionReader(br, mv.bodyoffset, mv.traileroffset-mv.bodyoffset)

	if !conf.decode {
		return ioutil.NopCloser(r), nil
	}

	if mv.chunked {
		r = httputil.NewChunkedReader(r)
	}
	switch mv.compress {
	case "gzip":
		gr, err := gzip.NewReader(r)
		if err != nil {
			return nil, err
		}
		return gr, nil
	case "deflate":
		return flate.NewReader(r), nil
	default:
		return ioutil.NopCloser(r), nil
	}
}

// TrailerReader returns an io.Reader that reads the HTTP request or response
// trailers, if present.
func (mv *MessageView) TrailerReader() io.Reader {
	r := bytes.NewReader(mv.message)
	end := int64(len(mv.message)) - mv.traileroffset

	return io.NewSectionReader(r, mv.traileroffset, end)
}

func (mv *MessageView) matchContentType(mct string) bool {
	for _, ct := range mv.cts {
		if strings.HasPrefix(mct, ct) {
			return true
		}
	}

	return false
}
