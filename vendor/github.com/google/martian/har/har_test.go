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

package har

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/proxyutil"
)

func TestModifyRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/path?query=true", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Add("Request-Header", "first")
	req.Header.Add("Request-Header", "second")

	cookie := &http.Cookie{
		Name:  "request",
		Value: "cookie",
	}
	req.AddCookie(cookie)

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	logger := NewLogger()
	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := log.Version, "1.2"; got != want {
		t.Errorf("log.Version: got %q, want %q", got, want)
	}

	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}

	entry := log.Entries[0]
	if got, want := time.Since(entry.StartedDateTime), time.Second; got > want {
		t.Errorf("entry.StartedDateTime: got %s, want less than %s", got, want)
	}

	hreq := entry.Request
	if got, want := hreq.Method, "GET"; got != want {
		t.Errorf("hreq.Method: got %q, want %q", got, want)
	}

	if got, want := hreq.URL, "http://example.com/path?query=true"; got != want {
		t.Errorf("hreq.URL: got %q, want %q", got, want)
	}

	if got, want := hreq.HTTPVersion, "HTTP/1.1"; got != want {
		t.Errorf("hreq.HTTPVersion: got %q, want %q", got, want)
	}

	if got, want := hreq.BodySize, int64(0); got != want {
		t.Errorf("hreq.BodySize: got %d, want %d", got, want)
	}

	if got, want := hreq.HeadersSize, int64(-1); got != want {
		t.Errorf("hreq.HeadersSize: got %d, want %d", got, want)
	}

	if got, want := len(hreq.QueryString), 1; got != want {
		t.Fatalf("len(hreq.QueryString): got %d, want %q", got, want)
	}

	qs := hreq.QueryString[0]
	if got, want := qs.Name, "query"; got != want {
		t.Errorf("qs.Name: got %q, want %q", got, want)
	}
	if got, want := qs.Value, "true"; got != want {
		t.Errorf("qs.Value: got %q, want %q", got, want)
	}

	wantHeaders := http.Header{
		"Request-Header": {"first", "second"},
		"Cookie":         {cookie.String()},
		"Host":           {"example.com"},
	}
	if got := headersToHTTP(hreq.Headers); !reflect.DeepEqual(got, wantHeaders) {
		t.Errorf("headers:\ngot:\n%+v\nwant:\n%+v", got, wantHeaders)
	}

	if got, want := len(hreq.Cookies), 1; got != want {
		t.Fatalf("len(hreq.Cookies): got %d, want %d", got, want)
	}

	hcookie := hreq.Cookies[0]
	if got, want := hcookie.Name, "request"; got != want {
		t.Errorf("hcookie.Name: got %q, want %q", got, want)
	}
	if got, want := hcookie.Value, "cookie"; got != want {
		t.Errorf("hcookie.Value: got %q, want %q", got, want)
	}
}

func headersToHTTP(hs []Header) http.Header {
	hh := http.Header{}
	for _, h := range hs {
		hh[h.Name] = append(hh[h.Name], h.Value)
	}
	return hh
}

func TestModifyResponse(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	res := proxyutil.NewResponse(301, strings.NewReader("response body"), req)
	res.ContentLength = 13
	res.Header.Add("Response-Header", "first")
	res.Header.Add("Response-Header", "second")
	res.Header.Set("Location", "google.com")

	expires := time.Now()
	cookie := &http.Cookie{
		Name:     "response",
		Value:    "cookie",
		Path:     "/",
		Domain:   "example.com",
		Expires:  expires,
		Secure:   true,
		HttpOnly: true,
	}
	res.Header.Set("Set-Cookie", cookie.String())

	logger := NewLogger()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}

	hres := log.Entries[0].Response
	if got, want := hres.Status, 301; got != want {
		t.Errorf("hres.Status: got %d, want %d", got, want)
	}

	if got, want := hres.StatusText, "Moved Permanently"; got != want {
		t.Errorf("hres.StatusText: got %q, want %q", got, want)
	}

	if got, want := hres.HTTPVersion, "HTTP/1.1"; got != want {
		t.Errorf("hres.HTTPVersion: got %q, want %q", got, want)
	}

	if got, want := hres.Content.Text, []byte("response body"); !bytes.Equal(got, want) {
		t.Errorf("hres.Content.Text: got %q, want %q", got, want)
	}

	wantHeaders := http.Header{
		"Response-Header": {"first", "second"},
		"Set-Cookie":      {cookie.String()},
		"Location":        {"google.com"},
		"Content-Length":  {"13"},
	}
	if got := headersToHTTP(hres.Headers); !reflect.DeepEqual(got, wantHeaders) {
		t.Errorf("headers:\ngot:\n%+v\nwant:\n%+v", got, wantHeaders)
	}

	if got, want := len(hres.Cookies), 1; got != want {
		t.Fatalf("len(hres.Cookies): got %d, want %d", got, want)
	}

	hcookie := hres.Cookies[0]
	if got, want := hcookie.Name, "response"; got != want {
		t.Errorf("hcookie.Name: got %q, want %q", got, want)
	}
	if got, want := hcookie.Value, "cookie"; got != want {
		t.Errorf("hcookie.Value: got %q, want %q", got, want)
	}
	if got, want := hcookie.Path, "/"; got != want {
		t.Errorf("hcookie.Path: got %q, want %q", got, want)
	}
	if got, want := hcookie.Domain, "example.com"; got != want {
		t.Errorf("hcookie.Domain: got %q, want %q", got, want)
	}
	if got, want := hcookie.Expires, expires; got.Equal(want) {
		t.Errorf("hcookie.Expires: got %s, want %s", got, want)
	}
	if !hcookie.HTTPOnly {
		t.Error("hcookie.HTTPOnly: got false, want true")
	}
	if !hcookie.Secure {
		t.Error("hcookie.Secure: got false, want true")
	}
}

func TestModifyRequestBodyURLEncoded(t *testing.T) {
	logger := NewLogger()

	body := strings.NewReader("first=true&second=false")
	req, err := http.NewRequest("POST", "http://example.com", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Errorf("len(log.Entries): got %v, want %v", got, want)
	}

	pd := log.Entries[0].Request.PostData
	if got, want := pd.MimeType, "application/x-www-form-urlencoded"; got != want {
		t.Errorf("PostData.MimeType: got %v, want %v", got, want)
	}

	if got, want := len(pd.Params), 2; got != want {
		t.Fatalf("len(PostData.Params): got %d, want %d", got, want)
	}

	for _, p := range pd.Params {
		var want string
		switch p.Name {
		case "first":
			want = "true"
		case "second":
			want = "false"
		default:
			t.Errorf("PostData.Params: got %q, want to not be present", p.Name)
			continue
		}

		if got := p.Value; got != want {
			t.Errorf("PostData.Params[%q]: got %q, want %q", p.Name, got, want)
		}
	}
}

func TestModifyRequestBodyArbitraryContentType(t *testing.T) {
	logger := NewLogger()

	body := "arbitrary binary data"
	req, err := http.NewRequest("POST", "http://www.example.com", strings.NewReader(body))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}

	pd := log.Entries[0].Request.PostData
	if got, want := pd.MimeType, ""; got != want {
		t.Errorf("PostData.MimeType: got %q, want %q", got, want)
	}
	if got, want := len(pd.Params), 0; got != want {
		t.Errorf("len(PostData.Params): got %d, want %d", got, want)
	}

	if got, want := pd.Text, body; got != want {
		t.Errorf("PostData.Text: got %q, want %q", got, want)
	}
}

func TestModifyRequestBodyMultipart(t *testing.T) {
	logger := NewLogger()

	body := new(bytes.Buffer)
	mpw := multipart.NewWriter(body)
	mpw.SetBoundary("boundary")

	if err := mpw.WriteField("key", "value"); err != nil {
		t.Errorf("mpw.WriteField(): got %v, want no error", err)
	}

	w, err := mpw.CreateFormFile("file", "test.txt")
	if _, err = w.Write([]byte("file contents")); err != nil {
		t.Fatalf("Write(): got %v, want no error", err)
	}
	mpw.Close()

	req, err := http.NewRequest("POST", "http://example.com", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Content-Type", mpw.FormDataContentType())

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}

	pd := log.Entries[0].Request.PostData
	if got, want := pd.MimeType, "multipart/form-data"; got != want {
		t.Errorf("PostData.MimeType: got %q, want %q", got, want)
	}
	if got, want := len(pd.Params), 2; got != want {
		t.Errorf("PostData.Params: got %d, want %d", got, want)
	}

	for _, p := range pd.Params {
		var want Param

		switch p.Name {
		case "key":
			want = Param{
				Filename:    "",
				ContentType: "",
				Value:       "value",
			}
		case "file":
			want = Param{
				Filename:    "test.txt",
				ContentType: "application/octet-stream",
				Value:       "file contents",
			}
		default:
			t.Errorf("pd.Params: got %q, want not to be present", p.Name)
			continue
		}

		if got, want := p.Filename, want.Filename; got != want {
			t.Errorf("p.Filename: got %q, want %q", got, want)
		}
		if got, want := p.ContentType, want.ContentType; got != want {
			t.Errorf("p.ContentType: got %q, want %q", got, want)
		}
		if got, want := p.Value, want.Value; got != want {
			t.Errorf("p.Value: got %q, want %q", got, want)
		}
	}
}

func TestModifyRequestErrorsOnDuplicateRequest(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if logger.ModifyRequest(req) == nil {
		t.Fatalf("ModifyRequest(): was supposed to error")
	}
}

func TestHARExportsTime(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	// Simulate fast network round trip.
	time.Sleep(10 * time.Millisecond)

	res := proxyutil.NewResponse(200, nil, req)

	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %v, want %v", got, want)
	}

	entry := log.Entries[0]
	min, max := int64(10), int64(100)
	if got := entry.Time; got < min || got > max {
		t.Errorf("entry.Time: got %dms, want between %dms and %vms", got, min, max)
	}
}

func TestReset(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}

	logger.Reset()

	log = logger.Export().Log
	if got, want := len(log.Entries), 0; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestExportSortsEntries(t *testing.T) {
	logger := NewLogger()
	count := 10

	for i := 0; i < count; i++ {
		req, err := http.NewRequest("GET", "http://example.com", nil)
		if err != nil {
			t.Fatalf("NewRequest(): got %v, want no error", err)
		}

		_, remove, err := martian.TestContext(req, nil, nil)
		if err != nil {
			t.Fatalf("martian.TestContext(): got %v, want no error", err)
		}
		defer remove()

		if err := logger.ModifyRequest(req); err != nil {
			t.Fatalf("ModifyRequest(): got %v, want no error", err)
		}
	}

	log := logger.Export().Log

	for i := 0; i < count-1; i++ {
		first := log.Entries[i]
		second := log.Entries[i+1]

		if got, want := first.StartedDateTime, second.StartedDateTime; got.After(want) {
			t.Errorf("entry.StartedDateTime: got %s, want to be before %s", got, want)
		}
	}
}

func TestExportIgnoresOrphanedResponse(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	// Reset before the response comes back.
	logger.Reset()

	res := proxyutil.NewResponse(200, nil, req)
	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 0; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestExportAndResetResetsCompleteRequests(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	logger.ExportAndReset()

	log := logger.Export().Log
	if got, want := len(log.Entries), 0; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestExportAndResetLeavesPendingRequests(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	logger.ExportAndReset()

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestExportAndResetExportsCompleteRequests(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.ExportAndReset().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestExportAndResetExportsCompleteRequestsWithPendingLeft(t *testing.T) {
	logger := NewLogger()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err = martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	_, remove, err = martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.ExportAndReset().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}

	log = logger.Export().Log
	if got, want := len(log.Entries), 2; got != want {
		t.Errorf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestSkippingLogging(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	ctx, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	ctx.SkipLogging()

	logger := NewLogger()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	res := proxyutil.NewResponse(200, nil, req)
	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 0; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}
}

func TestOptionResponseBodyLogging(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("NewRequest(): got %v, want no error", err)
	}

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	bdr := strings.NewReader("{\"response\": \"body\"}")
	res := proxyutil.NewResponse(200, bdr, req)
	res.ContentLength = int64(bdr.Len())
	res.Header.Set("Content-Type", "application/json")

	logger := NewLogger()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log := logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}

	if got, want := string(log.Entries[0].Response.Content.Text), "{\"response\": \"body\"}"; got != want {
		t.Fatalf("log.Entries[0].Response.Content.Text: got %s, want %s", got, want)
	}

	logger = NewLogger()
	logger.SetOption(BodyLogging(false))

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log = logger.Export().Log
	if got, want := len(log.Entries), 1; got != want {
		t.Fatalf("len(log.Entries): got %d, want %d", got, want)
	}

	if got, want := string(log.Entries[0].Response.Content.Text), ""; got != want {
		t.Fatalf("log.Entries[0].Response.Content: got %s, want %s", got, want)
	}

	logger = NewLogger()
	logger.SetOption(BodyLoggingForContentTypes("application/json"))

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log = logger.Export().Log
	if got, want := string(log.Entries[0].Response.Content.Text), "{\"response\": \"body\"}"; got != want {
		t.Fatalf("log.Entries[0].Response.Content: got %s, want %s", got, want)
	}

	logger = NewLogger()
	logger.SetOption(SkipBodyLoggingForContentTypes("application/json"))

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	if err := logger.ModifyResponse(res); err != nil {
		t.Fatalf("ModifyResponse(): got %v, want no error", err)
	}

	log = logger.Export().Log
	if got, want := string(log.Entries[0].Response.Content.Text), ""; got != want {
		t.Fatalf("log.Entries[0].Response.Content: got %v, want %v", got, want)
	}
}

func TestOptionRequestPostDataLogging(t *testing.T) {
	logger := NewLogger()
	logger.SetOption(PostDataLoggingForContentTypes("application/x-www-form-urlencoded"))

	body := strings.NewReader("first=true&second=false")
	req, err := http.NewRequest("POST", "http://example.com", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, remove, err := martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log := logger.Export().Log

	for _, param := range log.Entries[0].Request.PostData.Params {
		if param.Name == "first" {
			if got, want := param.Value, "true"; got != want {
				t.Fatalf("Params[%q].Value: got %s, want %s", param.Name, got, want)
			}
		}

		if param.Name == "second" {
			if got, want := param.Value, "false"; got != want {
				t.Fatalf("Params[%q].Value: got %s, want %s", param.Name, got, want)
			}
		}
	}

	logger = NewLogger()
	logger.SetOption(SkipPostDataLoggingForContentTypes("application/x-www-form-urlencoded"))

	body = strings.NewReader("first=true&second=false")
	req, err = http.NewRequest("POST", "http://example.com", body)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, remove, err = martian.TestContext(req, nil, nil)
	if err != nil {
		t.Fatalf("martian.TestContext(): got %v, want no error", err)
	}
	defer remove()

	if err := logger.ModifyRequest(req); err != nil {
		t.Fatalf("ModifyRequest(): got %v, want no error", err)
	}

	log = logger.Export().Log
	if got, want := len(log.Entries[0].Request.PostData.Params), 0; got != want {
		t.Fatalf("len(log.Entries[0].Request.PostData.Params): got %v, want %v", got, want)
	}
}

func TestJSONMarshalPostData(t *testing.T) {
	// Verify that encoding/json round-trips har.PostData with both text and binary data.
	for _, text := range []string{"hello", string([]byte{150, 151, 152})} {
		want := &PostData{
			MimeType: "m",
			Params:   []Param{{Name: "n", Value: "v"}},
			Text:     text,
		}
		data, err := json.Marshal(want)
		if err != nil {
			t.Fatal(err)
		}
		var got PostData
		if err := json.Unmarshal(data, &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(&got, want) {
			t.Errorf("got %+v, want %+v", &got, want)
		}
	}
}
