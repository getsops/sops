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

package proxyutil

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "http://www.example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Close = true

	res := NewResponse(200, nil, req)
	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Status, "200 OK"; got != want {
		t.Errorf("res.Status: got %q, want %q", got, want)
	}
	if !res.Close {
		t.Error("res.Close: got false, want true")
	}
	if got, want := res.Proto, "HTTP/1.1"; got != want {
		t.Errorf("res.Proto: got %q, want %q", got, want)
	}
	if got, want := res.ProtoMajor, 1; got != want {
		t.Errorf("res.ProtoMajor: got %d, want %d", got, want)
	}
	if got, want := res.ProtoMinor, 1; got != want {
		t.Errorf("res.ProtoMinor: got %d, want %d", got, want)
	}
	if res.Header == nil {
		t.Error("res.Header: got nil, want header")
	}
	if _, ok := res.Body.(io.ReadCloser); !ok {
		t.Error("res.Body.(io.ReadCloser): got !ok, want ok")
	}
	if got, want := res.Request, req; got != want {
		t.Errorf("res.Request: got %v, want %v", got, want)
	}
}

func TestWarning(t *testing.T) {
	hdr := http.Header{}
	err := fmt.Errorf("modifier error")

	Warning(hdr, err)

	if got, want := len(hdr["Warning"]), 1; got != want {
		t.Fatalf("len(hdr[%q]): got %d, want %d", "Warning", got, want)
	}

	want := `199 "martian" "modifier error"`
	if got := hdr["Warning"][0]; !strings.HasPrefix(got, want) {
		t.Errorf("hdr[%q][0]: got %q, want to have prefix %q", "Warning", got, want)
	}

	hdr.Set("Date", "Mon, 02 Jan 2006 15:04:05 GMT")
	Warning(hdr, err)

	if got, want := len(hdr["Warning"]), 2; got != want {
		t.Fatalf("len(hdr[%q]): got %d, want %d", "Warning", got, want)
	}

	want = `199 "martian" "modifier error" "Mon, 02 Jan 2006 15:04:05 GMT"`
	if got := hdr["Warning"][1]; got != want {
		t.Errorf("hdr[%q][1]: got %q, want %q", "Warning", got, want)
	}
}
