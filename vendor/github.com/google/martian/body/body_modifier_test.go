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

package body

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/google/martian/v3/messageview"
	"github.com/google/martian/v3/parse"
	"github.com/google/martian/v3/proxyutil"
)

func TestBodyModifier(t *testing.T) {
	mod := NewModifier([]byte("text"), "text/plain")

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Content-Encoding", "gzip")

	if err := mod.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if got, want := req.Header.Get("Content-Type"), "text/plain"; got != want {
		t.Errorf("req.Header.Get(%q): got %v, want %v", "Content-Type", got, want)
	}
	if got, want := req.ContentLength, int64(len([]byte("text"))); got != want {
		t.Errorf("req.ContentLength: got %d, want %d", got, want)
	}
	if got, want := req.Header.Get("Content-Encoding"), ""; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Content-Encoding", got, want)
	}

	got, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	req.Body.Close()

	if want := []byte("text"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	res := proxyutil.NewResponse(200, nil, req)
	res.Header.Set("Content-Encoding", "gzip")

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain"; got != want {
		t.Errorf("res.Header.Get(%q): got %v, want %v", "Content-Type", got, want)
	}
	if got, want := res.ContentLength, int64(len([]byte("text"))); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Content-Encoding"), ""; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Content-Encoding", got, want)
	}

	got, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("text"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}
}
func TestRangeHeaderRequestSingleRange(t *testing.T) {
	mod := NewModifier([]byte("0123456789"), "text/plain")

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Range", "bytes=1-4")

	res := proxyutil.NewResponse(200, nil, req)

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, http.StatusPartialContent; got != want {
		t.Errorf("res.Status: got %v, want %v", got, want)
	}
	if got, want := res.ContentLength, int64(len([]byte("1234"))); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Content-Range"), "bytes 1-4/10"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Content-Encoding", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("1234"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}
}

func TestRangeHeaderRequestSingleRangeHasAllTheBytes(t *testing.T) {
	mod := NewModifier([]byte("0123456789"), "text/plain")

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Range", "bytes=0-")

	res := proxyutil.NewResponse(200, nil, req)

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, http.StatusPartialContent; got != want {
		t.Errorf("res.Status: got %v, want %v", got, want)
	}
	if got, want := res.ContentLength, int64(len([]byte("0123456789"))); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Content-Range"), "bytes 0-9/10"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Content-Encoding", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("0123456789"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}
}

func TestRangeNoEndingIndexSpecified(t *testing.T) {
	mod := NewModifier([]byte("0123456789"), "text/plain")

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Range", "bytes=8-")

	res := proxyutil.NewResponse(200, nil, req)

	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, http.StatusPartialContent; got != want {
		t.Errorf("res.Status: got %v, want %v", got, want)
	}
	if got, want := res.ContentLength, int64(len([]byte("89"))); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Content-Range"), "bytes 8-9/10"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Content-Encoding", got, want)
	}
}

func TestRangeHeaderMultipartRange(t *testing.T) {
	mod := NewModifier([]byte("0123456789"), "text/plain")
	bndry := "3d6b6a416f9b5"
	mod.SetBoundary(bndry)

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Range", "bytes=1-4, 7-9")

	res := proxyutil.NewResponse(200, nil, req)
	if err := mod.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, http.StatusPartialContent; got != want {
		t.Errorf("res.Status: got %v, want %v", got, want)
	}

	if got, want := res.Header.Get("Content-Type"), "multipart/byteranges; boundary=3d6b6a416f9b5"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Content-Type", got, want)
	}

	mv := messageview.New()
	if err := mv.SnapshotResponse(res); err != nil {
		t.Fatalf("mv.SnapshotResponse(res): got %v, want no error", err)
	}

	br, err := mv.BodyReader()
	if err != nil {
		t.Fatalf("mv.BodyReader(): got %v, want no error", err)
	}

	mpr := multipart.NewReader(br, bndry)
	prt1, err := mpr.NextPart()
	if err != nil {
		t.Fatalf("mpr.NextPart(): got %v, want no error", err)
	}
	defer prt1.Close()

	if got, want := prt1.Header.Get("Content-Type"), "text/plain"; got != want {
		t.Errorf("prt1.Header.Get(%q): got %q, want %q", "Content-Type", got, want)
	}

	if got, want := prt1.Header.Get("Content-Range"), "bytes 1-4/10"; got != want {
		t.Errorf("prt1.Header.Get(%q): got %q, want %q", "Content-Range", got, want)
	}

	prt1b, err := ioutil.ReadAll(prt1)
	if err != nil {
		t.Errorf("ioutil.Readall(prt1): got %v, want no error", err)
	}

	if got, want := string(prt1b), "1234"; got != want {
		t.Errorf("prt1 body: got %s, want %s", got, want)
	}

	prt2, err := mpr.NextPart()
	if err != nil {
		t.Fatalf("mpr.NextPart(): got %v, want no error", err)
	}
	defer prt2.Close()

	if got, want := prt2.Header.Get("Content-Type"), "text/plain"; got != want {
		t.Errorf("prt2.Header.Get(%q): got %q, want %q", "Content-Type", got, want)
	}

	if got, want := prt2.Header.Get("Content-Range"), "bytes 7-9/10"; got != want {
		t.Errorf("prt2.Header.Get(%q): got %q, want %q", "Content-Range", got, want)
	}

	prt2b, err := ioutil.ReadAll(prt2)
	if err != io.ErrUnexpectedEOF && err != nil {
		t.Errorf("ioutil.Readall(prt2): got %v, want no error", err)
	}

	if got, want := string(prt2b), "789"; got != want {
		t.Errorf("prt2 body: got %s, want %s", got, want)
	}

	_, err = mpr.NextPart()
	if err == nil {
		t.Errorf("mpr.NextPart: want io.EOF, got no error")
	}
	if err != io.EOF {
		t.Errorf("mpr.NextPart: want io.EOF, got %v", err)
	}
}

func TestModifierFromJSON(t *testing.T) {
	data := base64.StdEncoding.EncodeToString([]byte("data"))
	msg := fmt.Sprintf(`{
	  "body.Modifier":{
		  "scope": ["response"],
  	  "contentType": "text/plain",
	  	"body": %q
    }
	}`, data)

	r, err := parse.FromJSON([]byte(msg))
	if err != nil {
		t.Fatalf("parse.FromJSON(): got %v, want no error", err)
	}

	resmod := r.ResponseModifier()

	if resmod == nil {
		t.Fatalf("resmod: got nil, want not nil")
	}

	req, err := http.NewRequest("GET", "/", strings.NewReader(""))
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := resmod.ModifyResponse(res); err != nil {
		t.Fatalf("resmod.ModifyResponse(): got %v, want no error", err)
	}

	if got, want := res.Header.Get("Content-Type"), "text/plain"; got != want {
		t.Errorf("res.Header.Get(%q): got %v, want %v", "Content-Type", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}
	res.Body.Close()

	if want := []byte("data"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}
}
