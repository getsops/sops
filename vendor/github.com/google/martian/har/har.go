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

// Package har collects HTTP requests and responses and stores them in HAR format.
//
// For more information on HAR, see:
// https://w3c.github.io/web-performance/specs/HAR/Overview.html
package har

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/google/martian/v3"
	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/messageview"
	"github.com/google/martian/v3/proxyutil"
)

// Logger maintains request and response log entries.
type Logger struct {
	bodyLogging     func(*http.Response) bool
	postDataLogging func(*http.Request) bool

	creator *Creator

	mu      sync.Mutex
	entries map[string]*Entry
	tail    *Entry
}

// HAR is the top level object of a HAR log.
type HAR struct {
	Log *Log `json:"log"`
}

// Log is the HAR HTTP request and response log.
type Log struct {
	// Version number of the HAR format.
	Version string `json:"version"`
	// Creator holds information about the log creator application.
	Creator *Creator `json:"creator"`
	// Entries is a list containing requests and responses.
	Entries []*Entry `json:"entries"`
}

// Creator is the program responsible for generating the log. Martian, in this case.
type Creator struct {
	// Name of the log creator application.
	Name string `json:"name"`
	// Version of the log creator application.
	Version string `json:"version"`
}

// Entry is a individual log entry for a request or response.
type Entry struct {
	// ID is the unique ID for the entry.
	ID string `json:"_id"`
	// StartedDateTime is the date and time stamp of the request start (ISO 8601).
	StartedDateTime time.Time `json:"startedDateTime"`
	// Time is the total elapsed time of the request in milliseconds.
	Time int64 `json:"time"`
	// Request contains the detailed information about the request.
	Request *Request `json:"request"`
	// Response contains the detailed information about the response.
	Response *Response `json:"response,omitempty"`
	// Cache contains information about a request coming from browser cache.
	Cache *Cache `json:"cache"`
	// Timings describes various phases within request-response round trip. All
	// times are specified in milliseconds.
	Timings *Timings `json:"timings"`
	next    *Entry
}

// Request holds data about an individual HTTP request.
type Request struct {
	// Method is the request method (GET, POST, ...).
	Method string `json:"method"`
	// URL is the absolute URL of the request (fragments are not included).
	URL string `json:"url"`
	// HTTPVersion is the Request HTTP version (HTTP/1.1).
	HTTPVersion string `json:"httpVersion"`
	// Cookies is a list of cookies.
	Cookies []Cookie `json:"cookies"`
	// Headers is a list of headers.
	Headers []Header `json:"headers"`
	// QueryString is a list of query parameters.
	QueryString []QueryString `json:"queryString"`
	// PostData is the posted data information.
	PostData *PostData `json:"postData,omitempty"`
	// HeaderSize is the Total number of bytes from the start of the HTTP request
	// message until (and including) the double CLRF before the body. Set to -1
	// if the info is not available.
	HeadersSize int64 `json:"headersSize"`
	// BodySize is the size of the request body (POST data payload) in bytes. Set
	// to -1 if the info is not available.
	BodySize int64 `json:"bodySize"`
}

// Response holds data about an individual HTTP response.
type Response struct {
	// Status is the response status code.
	Status int `json:"status"`
	// StatusText is the response status description.
	StatusText string `json:"statusText"`
	// HTTPVersion is the Response HTTP version (HTTP/1.1).
	HTTPVersion string `json:"httpVersion"`
	// Cookies is a list of cookies.
	Cookies []Cookie `json:"cookies"`
	// Headers is a list of headers.
	Headers []Header `json:"headers"`
	// Content contains the details of the response body.
	Content *Content `json:"content"`
	// RedirectURL is the target URL from the Location response header.
	RedirectURL string `json:"redirectURL"`
	// HeadersSize is the total number of bytes from the start of the HTTP
	// request message until (and including) the double CLRF before the body.
	// Set to -1 if the info is not available.
	HeadersSize int64 `json:"headersSize"`
	// BodySize is the size of the request body (POST data payload) in bytes. Set
	// to -1 if the info is not available.
	BodySize int64 `json:"bodySize"`
}

// Cache contains information about a request coming from browser cache.
type Cache struct {
	// Has no fields as they are not supported, but HAR requires the "cache"
	// object to exist.
}

// Timings describes various phases within request-response round trip. All
// times are specified in milliseconds
type Timings struct {
	// Send is the time required to send HTTP request to the server.
	Send int64 `json:"send"`
	// Wait is the time spent waiting for a response from the server.
	Wait int64 `json:"wait"`
	// Receive is the time required to read entire response from server or cache.
	Receive int64 `json:"receive"`
}

// Cookie is the data about a cookie on a request or response.
type Cookie struct {
	// Name is the cookie name.
	Name string `json:"name"`
	// Value is the cookie value.
	Value string `json:"value"`
	// Path is the path pertaining to the cookie.
	Path string `json:"path,omitempty"`
	// Domain is the host of the cookie.
	Domain string `json:"domain,omitempty"`
	// Expires contains cookie expiration time.
	Expires time.Time `json:"-"`
	// Expires8601 contains cookie expiration time in ISO 8601 format.
	Expires8601 string `json:"expires,omitempty"`
	// HTTPOnly is set to true if the cookie is HTTP only, false otherwise.
	HTTPOnly bool `json:"httpOnly,omitempty"`
	// Secure is set to true if the cookie was transmitted over SSL, false
	// otherwise.
	Secure bool `json:"secure,omitempty"`
}

// Header is an HTTP request or response header.
type Header struct {
	// Name is the header name.
	Name string `json:"name"`
	// Value is the header value.
	Value string `json:"value"`
}

// QueryString is a query string parameter on a request.
type QueryString struct {
	// Name is the query parameter name.
	Name string `json:"name"`
	// Value is the query parameter value.
	Value string `json:"value"`
}

// PostData describes posted data on a request.
type PostData struct {
	// MimeType is the MIME type of the posted data.
	MimeType string `json:"mimeType"`
	// Params is a list of posted parameters (in case of URL encoded parameters).
	Params []Param `json:"params"`
	// Text contains the posted data. Although its type is string, it may contain
	// binary data.
	Text string `json:"text"`
}

// pdBinary is the JSON representation of binary PostData.
type pdBinary struct {
	MimeType string `json:"mimeType"`
	// Params is a list of posted parameters (in case of URL encoded parameters).
	Params   []Param `json:"params"`
	Text     []byte  `json:"text"`
	Encoding string  `json:"encoding"`
}

// MarshalJSON returns a JSON representation of binary PostData.
func (p *PostData) MarshalJSON() ([]byte, error) {
	if utf8.ValidString(p.Text) {
		type noMethod PostData // avoid infinite recursion
		return json.Marshal((*noMethod)(p))
	}
	return json.Marshal(pdBinary{
		MimeType: p.MimeType,
		Params:   p.Params,
		Text:     []byte(p.Text),
		Encoding: "base64",
	})
}

// UnmarshalJSON populates PostData based on the []byte representation of
// the binary PostData.
func (p *PostData) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) { // conform to json.Unmarshaler spec
		return nil
	}
	var enc struct {
		Encoding string `json:"encoding"`
	}
	if err := json.Unmarshal(data, &enc); err != nil {
		return err
	}
	if enc.Encoding != "base64" {
		type noMethod PostData // avoid infinite recursion
		return json.Unmarshal(data, (*noMethod)(p))
	}
	var pb pdBinary
	if err := json.Unmarshal(data, &pb); err != nil {
		return err
	}
	p.MimeType = pb.MimeType
	p.Params = pb.Params
	p.Text = string(pb.Text)
	return nil
}

// Param describes an individual posted parameter.
type Param struct {
	// Name of the posted parameter.
	Name string `json:"name"`
	// Value of the posted parameter.
	Value string `json:"value,omitempty"`
	// Filename of a posted file.
	Filename string `json:"fileName,omitempty"`
	// ContentType is the content type of a posted file.
	ContentType string `json:"contentType,omitempty"`
}

// Content describes details about response content.
type Content struct {
	// Size is the length of the returned content in bytes. Should be equal to
	// response.bodySize if there is no compression and bigger when the content
	// has been compressed.
	Size int64 `json:"size"`
	// MimeType is the MIME type of the response text (value of the Content-Type
	// response header).
	MimeType string `json:"mimeType"`
	// Text contains the response body sent from the server or loaded from the
	// browser cache. This field is populated with textual content only. The text
	// field is either HTTP decoded text or a encoded (e.g. "base64")
	// representation of the response body. Leave out this field if the
	// information is not available.
	Text []byte `json:"text,omitempty"`
	// Encoding used for response text field e.g "base64". Leave out this field
	// if the text field is HTTP decoded (decompressed & unchunked), than
	// trans-coded from its original character set into UTF-8.
	Encoding string `json:"encoding,omitempty"`
}

// Option is a configurable setting for the logger.
type Option func(l *Logger)

// PostDataLogging returns an option that configures request post data logging.
func PostDataLogging(enabled bool) Option {
	return func(l *Logger) {
		l.postDataLogging = func(*http.Request) bool {
			return enabled
		}
	}
}

// PostDataLoggingForContentTypes returns an option that logs request bodies based
// on opting in to the Content-Type of the request.
func PostDataLoggingForContentTypes(cts ...string) Option {
	return func(l *Logger) {
		l.postDataLogging = func(req *http.Request) bool {
			rct := req.Header.Get("Content-Type")

			for _, ct := range cts {
				if strings.HasPrefix(strings.ToLower(rct), strings.ToLower(ct)) {
					return true
				}
			}

			return false
		}
	}
}

// SkipPostDataLoggingForContentTypes returns an option that logs request bodies based
// on opting out of the Content-Type of the request.
func SkipPostDataLoggingForContentTypes(cts ...string) Option {
	return func(l *Logger) {
		l.postDataLogging = func(req *http.Request) bool {
			rct := req.Header.Get("Content-Type")

			for _, ct := range cts {
				if strings.HasPrefix(strings.ToLower(rct), strings.ToLower(ct)) {
					return false
				}
			}

			return true
		}
	}
}

// BodyLogging returns an option that configures response body logging.
func BodyLogging(enabled bool) Option {
	return func(l *Logger) {
		l.bodyLogging = func(*http.Response) bool {
			return enabled
		}
	}
}

// BodyLoggingForContentTypes returns an option that logs response bodies based
// on opting in to the Content-Type of the response.
func BodyLoggingForContentTypes(cts ...string) Option {
	return func(l *Logger) {
		l.bodyLogging = func(res *http.Response) bool {
			rct := res.Header.Get("Content-Type")

			for _, ct := range cts {
				if strings.HasPrefix(strings.ToLower(rct), strings.ToLower(ct)) {
					return true
				}
			}

			return false
		}
	}
}

// SkipBodyLoggingForContentTypes returns an option that logs response bodies based
// on opting out of the Content-Type of the response.
func SkipBodyLoggingForContentTypes(cts ...string) Option {
	return func(l *Logger) {
		l.bodyLogging = func(res *http.Response) bool {
			rct := res.Header.Get("Content-Type")

			for _, ct := range cts {
				if strings.HasPrefix(strings.ToLower(rct), strings.ToLower(ct)) {
					return false
				}
			}

			return true
		}
	}
}

// NewLogger returns a HAR logger. The returned
// logger logs all request post data and response bodies by default.
func NewLogger() *Logger {
	l := &Logger{
		creator: &Creator{
			Name:    "martian proxy",
			Version: "2.0.0",
		},
		entries: make(map[string]*Entry),
	}
	l.SetOption(BodyLogging(true))
	l.SetOption(PostDataLogging(true))
	return l
}

// SetOption sets configurable options on the logger.
func (l *Logger) SetOption(opts ...Option) {
	for _, opt := range opts {
		opt(l)
	}
}

// ModifyRequest logs requests.
func (l *Logger) ModifyRequest(req *http.Request) error {
	ctx := martian.NewContext(req)
	if ctx.SkippingLogging() {
		return nil
	}

	id := ctx.ID()

	return l.RecordRequest(id, req)
}

// RecordRequest logs the HTTP request with the given ID. The ID should be unique
// per request/response pair.
func (l *Logger) RecordRequest(id string, req *http.Request) error {
	hreq, err := NewRequest(req, l.postDataLogging(req))
	if err != nil {
		return err
	}

	entry := &Entry{
		ID:              id,
		StartedDateTime: time.Now().UTC(),
		Request:         hreq,
		Cache:           &Cache{},
		Timings:         &Timings{},
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.entries[id]; exists {
		return fmt.Errorf("Duplicate request ID: %s", id)
	}
	l.entries[id] = entry
	if l.tail == nil {
		l.tail = entry
	}
	entry.next = l.tail.next
	l.tail.next = entry
	l.tail = entry

	return nil
}

// NewRequest constructs and returns a Request from req. If withBody is true,
// req.Body is read to EOF and replaced with a copy in a bytes.Buffer. An error
// is returned (and req.Body may be in an intermediate state) if an error is
// returned from req.Body.Read.
func NewRequest(req *http.Request, withBody bool) (*Request, error) {
	r := &Request{
		Method:      req.Method,
		URL:         req.URL.String(),
		HTTPVersion: req.Proto,
		HeadersSize: -1,
		BodySize:    req.ContentLength,
		QueryString: []QueryString{},
		Headers:     headers(proxyutil.RequestHeader(req).Map()),
		Cookies:     cookies(req.Cookies()),
	}

	for n, vs := range req.URL.Query() {
		for _, v := range vs {
			r.QueryString = append(r.QueryString, QueryString{
				Name:  n,
				Value: v,
			})
		}
	}

	pd, err := postData(req, withBody)
	if err != nil {
		return nil, err
	}
	r.PostData = pd

	return r, nil
}

// ModifyResponse logs responses.
func (l *Logger) ModifyResponse(res *http.Response) error {
	ctx := martian.NewContext(res.Request)
	if ctx.SkippingLogging() {
		return nil
	}
	id := ctx.ID()

	return l.RecordResponse(id, res)
}

// RecordResponse logs an HTTP response, associating it with the previously-logged
// HTTP request with the same ID.
func (l *Logger) RecordResponse(id string, res *http.Response) error {
	hres, err := NewResponse(res, l.bodyLogging(res))
	if err != nil {
		return err
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if e, ok := l.entries[id]; ok {
		e.Response = hres
		e.Time = time.Since(e.StartedDateTime).Nanoseconds() / 1000000
	}

	return nil
}

// NewResponse constructs and returns a Response from resp. If withBody is true,
// resp.Body is read to EOF and replaced with a copy in a bytes.Buffer. An error
// is returned (and resp.Body may be in an intermediate state) if an error is
// returned from resp.Body.Read.
func NewResponse(res *http.Response, withBody bool) (*Response, error) {
	r := &Response{
		HTTPVersion: res.Proto,
		Status:      res.StatusCode,
		StatusText:  http.StatusText(res.StatusCode),
		HeadersSize: -1,
		BodySize:    res.ContentLength,
		Headers:     headers(proxyutil.ResponseHeader(res).Map()),
		Cookies:     cookies(res.Cookies()),
	}

	if res.StatusCode >= 300 && res.StatusCode < 400 {
		r.RedirectURL = res.Header.Get("Location")
	}

	r.Content = &Content{
		Encoding: "base64",
		MimeType: res.Header.Get("Content-Type"),
	}

	if withBody {
		mv := messageview.New()
		if err := mv.SnapshotResponse(res); err != nil {
			return nil, err
		}

		br, err := mv.BodyReader(messageview.Decode())
		if err != nil {
			return nil, err
		}

		body, err := ioutil.ReadAll(br)
		if err != nil {
			return nil, err
		}

		r.Content.Text = body
		r.Content.Size = int64(len(body))
	}
	return r, nil
}

// Export returns the in-memory log.
func (l *Logger) Export() *HAR {
	l.mu.Lock()
	defer l.mu.Unlock()

	es := make([]*Entry, 0, len(l.entries))
	curr := l.tail
	for curr != nil {
		curr = curr.next
		es = append(es, curr)
		if curr == l.tail {
			break
		}
	}

	return l.makeHAR(es)
}

// ExportAndReset returns the in-memory log for completed requests, clearing them.
func (l *Logger) ExportAndReset() *HAR {
	l.mu.Lock()
	defer l.mu.Unlock()

	es := make([]*Entry, 0, len(l.entries))
	curr := l.tail
	prev := l.tail
	var first *Entry
	for curr != nil {
		curr = curr.next
		if curr.Response != nil {
			es = append(es, curr)
			delete(l.entries, curr.ID)
		} else {
			if first == nil {
				first = curr
			}
			prev.next = curr
			prev = curr
		}
		if curr == l.tail {
			break
		}
	}
	if len(l.entries) == 0 {
		l.tail = nil
	} else {
		l.tail = prev
		l.tail.next = first
	}

	return l.makeHAR(es)
}

func (l *Logger) makeHAR(es []*Entry) *HAR {
	return &HAR{
		Log: &Log{
			Version: "1.2",
			Creator: l.creator,
			Entries: es,
		},
	}
}

// Reset clears the in-memory log of entries.
func (l *Logger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = make(map[string]*Entry)
	l.tail = nil
}

func cookies(cs []*http.Cookie) []Cookie {
	hcs := make([]Cookie, 0, len(cs))

	for _, c := range cs {
		var expires string
		if !c.Expires.IsZero() {
			expires = c.Expires.Format(time.RFC3339)
		}

		hcs = append(hcs, Cookie{
			Name:        c.Name,
			Value:       c.Value,
			Path:        c.Path,
			Domain:      c.Domain,
			HTTPOnly:    c.HttpOnly,
			Secure:      c.Secure,
			Expires:     c.Expires,
			Expires8601: expires,
		})
	}

	return hcs
}

func headers(hs http.Header) []Header {
	hhs := make([]Header, 0, len(hs))

	for n, vs := range hs {
		for _, v := range vs {
			hhs = append(hhs, Header{
				Name:  n,
				Value: v,
			})
		}
	}

	return hhs
}

func postData(req *http.Request, logBody bool) (*PostData, error) {
	// If the request has no body (no Content-Length and Transfer-Encoding isn't
	// chunked), skip the post data.
	if req.ContentLength <= 0 && len(req.TransferEncoding) == 0 {
		return nil, nil
	}

	ct := req.Header.Get("Content-Type")
	mt, ps, err := mime.ParseMediaType(ct)
	if err != nil {
		log.Errorf("har: cannot parse Content-Type header %q: %v", ct, err)
		mt = ct
	}

	pd := &PostData{
		MimeType: mt,
		Params:   []Param{},
	}

	if !logBody {
		return pd, nil
	}

	mv := messageview.New()
	if err := mv.SnapshotRequest(req); err != nil {
		return nil, err
	}

	br, err := mv.BodyReader()
	if err != nil {
		return nil, err
	}

	switch mt {
	case "multipart/form-data":
		mpr := multipart.NewReader(br, ps["boundary"])

		for {
			p, err := mpr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			defer p.Close()

			body, err := ioutil.ReadAll(p)
			if err != nil {
				return nil, err
			}

			pd.Params = append(pd.Params, Param{
				Name:        p.FormName(),
				Filename:    p.FileName(),
				ContentType: p.Header.Get("Content-Type"),
				Value:       string(body),
			})
		}
	case "application/x-www-form-urlencoded":
		body, err := ioutil.ReadAll(br)
		if err != nil {
			return nil, err
		}

		vs, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}

		for n, vs := range vs {
			for _, v := range vs {
				pd.Params = append(pd.Params, Param{
					Name:  n,
					Value: v,
				})
			}
		}
	default:
		body, err := ioutil.ReadAll(br)
		if err != nil {
			return nil, err
		}

		pd.Text = string(body)
	}

	return pd, nil
}
