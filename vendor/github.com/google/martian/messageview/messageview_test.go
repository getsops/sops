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

package messageview

import (
	"bufio"
	"bytes"
	"compress/flate"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/google/martian/v3/proxyutil"
)

func TestRequestViewHeadersOnly(t *testing.T) {
	body := strings.NewReader("body content")
	req, err := http.NewRequest("GET", "http://example.com/path?k=v", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.ContentLength = int64(body.Len())
	req.Header.Set("Request-Header", "true")

	mv := New()
	mv.SkipBody(true)
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "GET http://example.com/path?k=v HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Content-Length: 12\r\n" +
		"Request-Header: true\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	if _, err := br.Read(nil); err != io.EOF {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want io.EOF", err)
	}

	r, err := mv.Reader()
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}
}

func TestRequestView(t *testing.T) {
	body := strings.NewReader("body content")
	req, err := http.NewRequest("GET", "http://example.com/path?k=v", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Request-Header", "true")

	// Force Content Length to be unset to simulate lack of Content-Length and
	// Transfer-Encoding which is valid.
	req.ContentLength = -1

	mv := New()
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "GET http://example.com/path?k=v HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Request-Header: true\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "body content"
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	r, err := mv.Reader()
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}

	// Sanity check to ensure it still parses.
	if _, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(got))); err != nil {
		t.Fatalf("http.ReadRequest(): got %v, want no error", err)
	}
}

func TestRequestViewSkipBodyUnlessContentType(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", strings.NewReader("body content"))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.ContentLength = 12
	req.Header.Set("Content-Type", "text/plain; charset=utf-8")

	mv := New()
	mv.SkipBodyUnlessContentType("text/plain")
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "body content"
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	req.Header.Set("Content-Type", "image/png")
	mv = New()
	mv.SkipBodyUnlessContentType("text/plain")
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	br, err = mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	if _, err := br.Read(nil); err != io.EOF {
		t.Fatalf("br.Read(): got %v, want io.EOF", err)
	}
}

func TestRequestViewChunkedTransferEncoding(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/path?k=v", strings.NewReader("body content"))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.TransferEncoding = []string{"chunked"}
	req.Header.Set("Trailer", "Trailer-Header")
	req.Trailer = http.Header{
		"Trailer-Header": []string{"true"},
	}

	mv := New()
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "GET http://example.com/path?k=v HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Trailer: Trailer-Header\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "c\r\nbody content\r\n0\r\n"
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	got, err = ioutil.ReadAll(mv.TrailerReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.TrailerReader()): got %v, want no error", err)
	}

	trailerwant := "Trailer-Header: true\r\n"
	if !bytes.Equal(got, []byte(trailerwant)) {
		t.Fatalf("mv.TrailerReader(): got %q, want %q", got, trailerwant)
	}

	r, err := mv.Reader()
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant + trailerwant); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}

	// Sanity check to ensure it still parses.
	if _, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(got))); err != nil {
		t.Fatalf("http.ReadRequest(): got %v, want no error", err)
	}
}

func TestRequestViewDecodeGzipContentEncoding(t *testing.T) {
	body := new(bytes.Buffer)
	gw := gzip.NewWriter(body)
	gw.Write([]byte("body content"))
	gw.Flush()
	gw.Close()

	req, err := http.NewRequest("GET", "http://example.com/path?k=v", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.TransferEncoding = []string{"chunked"}
	req.Header.Set("Content-Encoding", "gzip")

	mv := New()
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "GET http://example.com/path?k=v HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Encoding: gzip\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader(Decode())
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, wt o error", err)
	}

	bodywant := "body content"

	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	r, err := mv.Reader(Decode())
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant + "\r\n"); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}
}

func TestRequestViewDecodeDeflateContentEncoding(t *testing.T) {
	body := new(bytes.Buffer)
	dw, err := flate.NewWriter(body, -1)
	if err != nil {
		t.Fatalf("flate.NewWriter(): got %v, want no error", err)
	}
	dw.Write([]byte("body content"))
	dw.Flush()
	dw.Close()

	req, err := http.NewRequest("GET", "http://example.com/path?k=v", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.TransferEncoding = []string{"chunked"}
	req.Header.Set("Content-Encoding", "deflate")

	mv := New()
	if err := mv.SnapshotRequest(req); err != nil {
		t.Fatalf("SnapshotRequest(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "GET http://example.com/path?k=v HTTP/1.1\r\n" +
		"Host: example.com\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Encoding: deflate\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader(Decode())
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "body content"

	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	r, err := mv.Reader(Decode())
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant + "\r\n"); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}
}

func TestResponseViewHeadersOnly(t *testing.T) {
	body := strings.NewReader("body content")
	res := proxyutil.NewResponse(200, body, nil)
	res.ContentLength = 12
	res.Header.Set("Response-Header", "true")

	mv := New()
	mv.SkipBody(true)
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 12\r\n" +
		"Response-Header: true\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	if _, err := br.Read(nil); err != io.EOF {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want io.EOF", err)
	}

	r, err := mv.Reader()
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}
}

func TestResponseView(t *testing.T) {
	body := strings.NewReader("body content")
	res := proxyutil.NewResponse(200, body, nil)
	res.ContentLength = 12
	res.Header.Set("Response-Header", "true")

	mv := New()
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 12\r\n" +
		"Response-Header: true\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "body content"
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	r, err := mv.Reader()
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}

	// Sanity check to ensure it still parses.
	if _, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(got)), nil); err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
}

func TestResponseViewSkipBodyUnlessContentType(t *testing.T) {
	res := proxyutil.NewResponse(200, strings.NewReader("body content"), nil)
	res.ContentLength = 12
	res.Header.Set("Content-Type", "text/plain; charset=utf-8")

	mv := New()
	mv.SkipBodyUnlessContentType("text/plain")
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "body content"
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	res.Header.Set("Content-Type", "image/png")
	mv = New()
	mv.SkipBodyUnlessContentType("text/plain")
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	br, err = mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	if _, err := br.Read(nil); err != io.EOF {
		t.Fatalf("br.Read(): got %v, want io.EOF", err)
	}
}

func TestResponseViewChunkedTransferEncoding(t *testing.T) {
	body := strings.NewReader("body content")
	res := proxyutil.NewResponse(200, body, nil)
	res.TransferEncoding = []string{"chunked"}
	res.Header.Set("Trailer", "Trailer-Header")
	res.Trailer = http.Header{
		"Trailer-Header": []string{"true"},
	}

	mv := New()
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "HTTP/1.1 200 OK\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Trailer: Trailer-Header\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, want no error", err)
	}

	bodywant := "c\r\nbody content\r\n0\r\n"
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	got, err = ioutil.ReadAll(mv.TrailerReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.TrailerReader()): got %v, want no error", err)
	}

	trailerwant := "Trailer-Header: true\r\n"
	if !bytes.Equal(got, []byte(trailerwant)) {
		t.Fatalf("mv.TrailerReader(): got %q, want %q", got, trailerwant)
	}

	r, err := mv.Reader()
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant + trailerwant); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}

	// Sanity check to ensure it still parses.
	if _, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(got)), nil); err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
}

func TestResponseViewDecodeGzipContentEncoding(t *testing.T) {
	body := new(bytes.Buffer)
	gw := gzip.NewWriter(body)
	gw.Write([]byte("body content"))
	gw.Flush()
	gw.Close()

	res := proxyutil.NewResponse(200, body, nil)
	res.TransferEncoding = []string{"chunked"}
	res.Header.Set("Content-Encoding", "gzip")

	mv := New()
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "HTTP/1.1 200 OK\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Encoding: gzip\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader(Decode())
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, wt o error", err)
	}

	bodywant := "body content"

	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	r, err := mv.Reader(Decode())
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant + "\r\n"); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}
}

func TestResponseViewDecodeGzipContentEncodingPartial(t *testing.T) {
	bodywant := "partial content"
	res := proxyutil.NewResponse(206, strings.NewReader(bodywant), nil)
	res.TransferEncoding = []string{"chunked"}
	res.Header.Set("Content-Encoding", "gzip")

	mv := New()
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}
	br, err := mv.BodyReader(Decode())
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, wt o error", err)
	}
	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}
}

func TestResponseViewDecodeDeflateContentEncoding(t *testing.T) {
	body := new(bytes.Buffer)
	dw, err := flate.NewWriter(body, -1)
	if err != nil {
		t.Fatalf("flate.NewWriter(): got %v, want no error", err)
	}
	dw.Write([]byte("body content"))
	dw.Flush()
	dw.Close()

	res := proxyutil.NewResponse(200, body, nil)
	res.TransferEncoding = []string{"chunked"}
	res.Header.Set("Content-Encoding", "deflate")

	mv := New()
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("SnapshotResponse(): got %v, want no error", err)
	}

	got, err := ioutil.ReadAll(mv.HeaderReader())
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.HeaderReader()): got %v, want no error", err)
	}

	hdrwant := "HTTP/1.1 200 OK\r\n" +
		"Transfer-Encoding: chunked\r\n" +
		"Content-Encoding: deflate\r\n\r\n"

	if !bytes.Equal(got, []byte(hdrwant)) {
		t.Fatalf("mv.HeaderReader(): got %q, want %q", got, hdrwant)
	}

	br, err := mv.BodyReader(Decode())
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	got, err = ioutil.ReadAll(br)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.BodyReader()): got %v, wt o error", err)
	}

	bodywant := "body content"

	if !bytes.Equal(got, []byte(bodywant)) {
		t.Fatalf("mv.BodyReader(): got %q, want %q", got, bodywant)
	}

	r, err := mv.Reader(Decode())
	if err != nil {
		t.Fatalf("mv.Reader(): got %v, want no error", err)
	}
	got, err = ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(mv.Reader()): got %v, want no error", err)
	}

	if want := []byte(hdrwant + bodywant + "\r\n"); !bytes.Equal(got, want) {
		t.Fatalf("mv.Read(): got %q, want %q", got, want)
	}
}
