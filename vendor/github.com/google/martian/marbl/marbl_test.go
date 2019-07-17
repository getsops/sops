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
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/proxyutil"
)

func TestMarkAPIRequestsWithHeader(t *testing.T) {
	areq, err := http.NewRequest("POST", "http://localhost:8080/configure", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(areq, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	ctx.APIRequest()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, removereq, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer removereq()

	var b bytes.Buffer

	s := NewStream(&b)
	s.LogRequest("00000000", areq)
	s.LogRequest("00000001", req)
	s.Close()

	headers := make(map[string]string)
	reader := NewReader(&b)

	for {
		frame, err := reader.ReadFrame()
		if frame == nil {
			break
		}
		if err != nil && err != io.EOF {
			t.Fatalf("reader.ReadFrame(): got %v, want no error or io.EOF", err)
		}

		headerFrame, ok := frame.(Header)
		if !ok {
			t.Fatalf("frame.(Header): couldn't convert frame '%v' to a headerFrame", frame)
		}
		headers[headerFrame.ID+headerFrame.Name] = headerFrame.Value
	}

	apih, ok := headers["00000000:api"]
	if !ok {
		t.Errorf("headers[00000000:api]: got no such header, want :api (headers were: %v)", headers)
	}

	if got, want := apih, "true"; got != want {
		t.Errorf("headers[%q]: got %v, want %q", "00000000:api", got, want)
	}

	_, ok = headers["00000001:api"]
	if got, want := ok, false; got != want {
		t.Error("headers[00000001:api]: got :api header, want no header for non-api requests")
	}
}

func TestSendTimestampWithLogRequest(t *testing.T) {
	req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	var b bytes.Buffer
	s := NewStream(&b)

	before := time.Now().UnixNano() / 1000 / 1000
	s.LogRequest("Fake_Id0", req)
	s.Close()
	after := time.Now().UnixNano() / 1000 / 1000

	headers := make(map[string]string)
	reader := NewReader(&b)

	for {
		frame, err := reader.ReadFrame()
		if frame == nil {
			break
		}
		if err != nil && err != io.EOF {
			t.Fatalf("reader.ReadFrame(): got %v, want no error or io.EOF", err)
		}

		headerFrame, ok := frame.(Header)
		if !ok {
			t.Fatalf("frame.(Header): couldn't convert frame '%v' to a headerFrame", frame)
		}
		headers[headerFrame.Name] = headerFrame.Value
	}

	timestr, ok := headers[":timestamp"]
	if !ok {
		t.Fatalf("headers[:timestamp]: got no such header, want :timestamp (headers were: %v)", headers)
	}
	ts, err := strconv.ParseInt(timestr, 10, 64)
	if err != nil {
		t.Fatalf("strconv.ParseInt: got %s, want no error. Invalidly formatted timestamp ('%s')", err, timestr)
	}
	if ts < before || ts > after {
		t.Fatalf("headers[:timestamp]: got %d, want timestamp between %d and %d", ts, before, after)
	}
}

func TestSendTimestampWithLogResponse(t *testing.T) {
	req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(200, nil, req)
	var b bytes.Buffer
	s := NewStream(&b)

	before := time.Now().UnixNano() / 1000 / 1000
	s.LogResponse("Fake_Id1", res)
	s.Close()
	after := time.Now().UnixNano() / 1000 / 1000

	headers := make(map[string]string)
	reader := NewReader(&b)

	for {
		frame, err := reader.ReadFrame()
		if frame == nil {
			break
		}
		if err != nil && err != io.EOF {
			t.Fatalf("reader.ReadFrame(): got %v, want no error or io.EOF", err)
		}

		headerFrame, ok := frame.(Header)
		if !ok {
			t.Fatalf("frame.(Header): couldn't convert frame '%v' to a headerFrame", frame)
		}
		headers[headerFrame.Name] = headerFrame.Value
	}

	timestr, ok := headers[":timestamp"]
	if !ok {
		t.Fatalf("headers[:timestamp]: got no such header, want :timestamp (headers were: %v)", headers)
	}
	ts, err := strconv.ParseInt(timestr, 10, 64)
	if err != nil {
		t.Fatalf("strconv.ParseInt: got %s, want no error. Invalidly formatted timestamp ('%s')", err, timestr)
	}
	if ts < before || ts > after {
		t.Fatalf("headers[:timestamp]: got %d, want timestamp between %d and %d (headers were: %v)", ts, before, after, headers)
	}
}

func TestBodyLoggingWithOneRead(t *testing.T) {
	// Test scenario:
	// 1. Prepare HTTP request with body containing a string.
	// 2. Initialize marbl logging on this request.
	// 3. Read body of the request in single Read() and verity that it matches
	//    original string.
	// 4. Parse marbl data, extract DataFrames and verify that they match
	// .  original string.
	body := "hello, world"
	req, err := http.NewRequest("POST", "http://example.com", strings.NewReader(body))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	var b bytes.Buffer
	s := NewStream(&b)
	s.LogRequest("Fake_Id0", req)

	// Read request body into big slice.
	bodybytes := make([]byte, 100)

	// First read. Due to implementation details of strings.Read
	// it reads all bytes but doesn't return EOF.
	n, err := req.Body.Read(bodybytes)
	if n != len(body) {
		t.Fatalf("req.Body.Read(): expected to read %v bytes but read %v", len(body), n)
	}
	if body != string(bodybytes[:n]) {
		t.Fatalf("req.Body.Read(): expected to read %v but read %v", body, string(bodybytes[:n]))
	}
	if err != nil {
		t.Fatalf("req.Body.Read(): first read expected to be successful but got error %v", err)
	}

	// second read. We already consumed the whole string on the first read
	// so now it should be 0 bytes and EOF.
	n, err = req.Body.Read(bodybytes)
	if n != 0 {
		t.Fatalf("req.Body.Read(): expected to read 0 bytes but read %v", n)
	}
	if err != io.EOF {
		t.Fatalf("req.Body.Read(): expected EOF but got %v", err)
	}

	s.Close()
	reader := NewReader(&b)
	bodybytes = readAllDataFrames(reader, "Fake_Id0", t)
	if len(bodybytes) != len(body) {
		t.Fatalf("readAllDataFrames(): expected .marbl data to have %v bytes, but got %v", len(body), len(bodybytes))
	}
	if body != string(bodybytes) {
		t.Fatalf("readAllDataFrames(): expected .marbl data to have string %v but got %v", body, string(bodybytes))
	}
}

func TestBodyLogging_ManyReads(t *testing.T) {
	// Test scenario:
	// 1. Prepare HTTP request with body containing a string.
	// 2. Initialize marbl logging on this request.
	// 3. Read body of the request in many reads, 1 byte per read and
	// .  verify that it matches original string.
	// 4. Parse marbl data, extract DataFrames and verify that they match
	// .  original string.
	body := "hello, world"
	req, err := http.NewRequest("POST", "http://example.com", strings.NewReader(body))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	var b bytes.Buffer
	s := NewStream(&b)
	s.LogRequest("Fake_Id0", req)

	// Read request body into single byte slice.
	bodybytes := make([]byte, 1)

	for i := 0; i < len(body); i++ {
		// first read
		n, err := req.Body.Read(bodybytes)
		if n != 1 {
			t.Fatalf("req.Body.Read(): expected to read 1 byte but read %v", n)
		}
		if body[i] != bodybytes[0] {
			t.Fatalf("req.Body.Read(): expected to read %v but read %v", body[i], bodybytes[0])
		}
		if err != nil {
			t.Fatalf("req.Body.Read(): read expected to be successfully but got error %v", err)
		}
	}

	// last read. We already consumed the whole string on the previous reads
	// so now it should be 0 bytes and EOF.
	n, err := req.Body.Read(bodybytes)
	if n != 0 {
		t.Fatalf("req.Body.Read(): expected to read 0 bytes but read %v", n)
	}
	if err != io.EOF {
		t.Fatalf("req.Body.Read(): expected EOF but got %v", err)
	}

	s.Close()
	reader := NewReader(&b)
	bodybytes = readAllDataFrames(reader, "Fake_Id0", t)
	if len(bodybytes) != len(body) {
		t.Fatalf("readAllDataFrames(): expected .marbl data to have %v bytes, but got %v", len(body), len(bodybytes))
	}
	if body != string(bodybytes) {
		t.Fatalf("readAllDataFrames(): expected .marbl data to have string %v but got %v", body, string(bodybytes))
	}
}

func TestReturnOriginalRequestPathAndQuery(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo%20bar?baz%20qux", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("TestContext(): got %v, want no error", err)
	}
	defer remove()

	var b bytes.Buffer
	s := NewStream(&b)

	s.LogRequest("Fake_Id0", req)
	s.Close()

	headers := make(map[string]string)
	reader := NewReader(&b)

	for {
		frame, err := reader.ReadFrame()
		if frame == nil {
			break
		}
		if err != nil && err != io.EOF {
			t.Fatalf("reader.ReadFrame(): got %v, want no error or io.EOF", err)
		}

		headerFrame, ok := frame.(Header)
		if !ok {
			t.Fatalf("frame.(Header): couldn't convert frame '%v' to a headerFrame", frame)
		}
		headers[headerFrame.Name] = headerFrame.Value
	}

	path := headers[":path"]
	if path != "/foo%20bar" {
		t.Fatalf("headers[:path]: expected /foo%%20bar but got %s", path)
	}
	query := headers[":query"]
	if query != "baz%20qux" {
		t.Fatalf("headers[:query]: expected baz%%20qux but got %s", query)
	}
}

// readAllDataFrames reads all DataFrames with reader, filters the one that match provided
// id and assembles data from all frames into single slice. It expects that
// there is only one slice of DataFrames with provided id.
func readAllDataFrames(reader *Reader, id string, t *testing.T) []byte {
	res := make([]byte, 0)
	term := false
	var i uint32
	for {
		frame, _ := reader.ReadFrame()
		if frame == nil {
			break
		}
		if frame.FrameType() == DataFrame {
			df := frame.(Data)
			if df.ID != id {
				continue
			}
			if term {
				t.Fatal("DataFrame after terminal frame are not allowed.")
			}
			if df.Index != i {
				t.Fatalf("expected DataFrame index %v but got %v", i, df.Index)
			}
			term = df.Terminal
			res = append(res, df.Data...)
			i++
		}
	}
	
	if !term {
		t.Fatal("didn't see terminal DataFrame")
	}
	
	return res
}
