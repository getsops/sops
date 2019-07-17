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

package martianlog

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func ExampleLogger() {
	l := NewLogger()
	l.SetLogFunc(func(line string) {
		// Remove \r to make it easier to test with examples.
		fmt.Print(strings.Replace(line, "\r", "", -1))
	})
	l.SetDecode(true)

	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	gw.Write([]byte("request content"))
	gw.Close()

	req, err := http.NewRequest("GET", "http://example.com/path?querystring", buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.TransferEncoding = []string{"chunked"}
	req.Header.Set("Content-Encoding", "gzip")

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer remove()

	if err := l.ModifyRequest(req); err != nil {
		fmt.Println(err)
		return
	}

	res := proxyutil.NewResponse(200, strings.NewReader("response content"), req)
	res.ContentLength = 16
	res.Header.Set("Date", "Tue, 15 Nov 1994 08:12:31 GMT")
	res.Header.Set("Other-Header", "values")

	if err := l.ModifyResponse(res); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// --------------------------------------------------------------------------------
	// Request to http://example.com/path?querystring
	// --------------------------------------------------------------------------------
	// GET http://example.com/path?querystring HTTP/1.1
	// Host: example.com
	// Transfer-Encoding: chunked
	// Content-Encoding: gzip
	//
	// request content
	//
	// --------------------------------------------------------------------------------
	//
	// --------------------------------------------------------------------------------
	// Response from http://example.com/path?querystring
	// --------------------------------------------------------------------------------
	// HTTP/1.1 200 OK
	// Content-Length: 16
	// Date: Tue, 15 Nov 1994 08:12:31 GMT
	// Other-Header: values
	//
	// response content
	// --------------------------------------------------------------------------------
}

func TestLoggerFromJSON(t *testing.T) {
	msg := []byte(`{
		"log.Logger": {
			"scope": ["request", "response"],
			"headersOnly": true,
			"decode": true
		}
	}`)

	r, err := parse.FromJSON(msg)
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	reqmod := r.RequestModifier()
	if reqmod == nil {
		t.Fatal("r.RequestModifier(): got nil, want not nil")
	}
	if _, ok := reqmod.(*Logger); !ok {
		t.Error("reqmod.(*Logger): got !ok, want ok")
	}

	resmod := r.ResponseModifier()
	if resmod == nil {
		t.Fatal("r.ResponseModifier(): got nil, want not nil")
	}

	l, ok := resmod.(*Logger)
	if !ok {
		t.Error("resmod.(*Logger); got !ok, want ok")
	}

	if !l.headersOnly {
		t.Error("l.headersOnly: got false, want true")
	}

	if !l.decode {
		t.Error("l.decode: got false, want true")
	}
}
