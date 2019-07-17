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
	"net/http"
	"reflect"
	"testing"
)

func TestRequestHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	h := RequestHeader(req)

	tt := []struct {
		name  string
		value string
	}{
		{
			name:  "Host",
			value: "example.com",
		},
		{
			name:  "Test-Header",
			value: "true",
		},
		{
			name:  "Content-Length",
			value: "100",
		},
		{
			name:  "Transfer-Encoding",
			value: "chunked",
		},
	}

	for i, tc := range tt {
		if err := h.Set(tc.name, tc.value); err != nil {
			t.Errorf("%d. h.Set(%q, %q): got %v, want no error", i, tc.name, tc.value, err)
		}
	}

	if got, want := req.Host, "example.com"; got != want {
		t.Errorf("req.Host: got %q, want %q", got, want)
	}
	if got, want := req.Header.Get("Test-Header"), "true"; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Test-Header", got, want)
	}
	if got, want := req.ContentLength, int64(100); got != want {
		t.Errorf("req.ContentLength: got %d, want %d", got, want)
	}
	if got, want := req.TransferEncoding, []string{"chunked"}; !reflect.DeepEqual(got, want) {
		t.Errorf("req.TransferEncoding: got %v, want %v", got, want)
	}

	if got, want := len(h.Map()), 4; got != want {
		t.Errorf("h.Map(): got %d entries, want %d entries", got, want)
	}

	for n, vs := range h.Map() {
		var want string
		switch n {
		case "Host":
			want = "example.com"
		case "Content-Length":
			want = "100"
		case "Transfer-Encoding":
			want = "chunked"
		case "Test-Header":
			want = "true"
		default:
			t.Errorf("h.Map(): got unexpected %s header", n)
		}

		if got := vs[0]; got != want {
			t.Errorf("h.Map(): got %s header with value %s, want value %s", n, got, want)
		}
	}

	for i, tc := range tt {
		got, ok := h.All(tc.name)
		if !ok {
			t.Errorf("%d. h.All(%q): got false, want true", i, tc.name)
		}

		if want := []string{tc.value}; !reflect.DeepEqual(got, want) {
			t.Errorf("%d. h.All(%q): got %v, want %v", i, tc.name, got, want)
		}

		if got, want := h.Get(tc.name), tc.value; got != want {
			t.Errorf("%d. h.Get(%q): got %q, want %q", i, tc.name, got, want)
		}

		h.Del(tc.name)
	}

	if got, want := req.Host, ""; got != want {
		t.Errorf("req.Host: got %q, want %q", got, want)
	}
	if got, want := req.Header.Get("Test-Header"), ""; got != want {
		t.Errorf("req.Header.Get(%q): got %q, want %q", "Test-Header", got, want)
	}
	if got, want := req.ContentLength, int64(-1); got != want {
		t.Errorf("req.ContentLength: got %d, want %d", got, want)
	}
	if got := req.TransferEncoding; got != nil {
		t.Errorf("req.TransferEncoding: got %v, want nil", got)
	}

	for i, tc := range tt {
		if got, want := h.Get(tc.name), ""; got != want {
			t.Errorf("%d. h.Get(%q): got %q, want %q", i, tc.name, got, want)
		}

		got, ok := h.All(tc.name)
		if ok {
			t.Errorf("%d. h.All(%q): got ok, want !ok", i, tc.name)
		}
		if got != nil {
			t.Errorf("%d. h.All(%q): got %v, want nil", i, tc.name, got)
		}
	}
}

func TestRequestHeaderAdd(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Host = "" // Set to empty so add may overwrite.

	h := RequestHeader(req)

	tt := []struct {
		name             string
		values           []string
		errOnSecondValue bool
	}{
		{
			name:             "Host",
			values:           []string{"example.com", "invalid.com"},
			errOnSecondValue: true,
		},
		{
			name:   "Test-Header",
			values: []string{"first", "second"},
		},
		{
			name:             "Content-Length",
			values:           []string{"100", "101"},
			errOnSecondValue: true,
		},
		{
			name:   "Transfer-Encoding",
			values: []string{"chunked", "gzip"},
		},
	}

	for i, tc := range tt {
		if err := h.Add(tc.name, tc.values[0]); err != nil {
			t.Errorf("%d. h.Add(%q, %q): got %v, want no error", i, tc.name, tc.values[0], err)
		}
		if err := h.Add(tc.name, tc.values[1]); err != nil && !tc.errOnSecondValue {
			t.Errorf("%d. h.Add(%q, %q): got %v, want no error", i, tc.name, tc.values[1], err)
		}
	}

	if got, want := req.Host, "example.com"; got != want {
		t.Errorf("req.Host: got %q, want %q", got, want)
	}
	if got, want := req.Header["Test-Header"], []string{"first", "second"}; !reflect.DeepEqual(got, want) {
		t.Errorf("req.Header[%q]: got %v, want %v", "Test-Header", got, want)
	}
	if got, want := req.ContentLength, int64(100); got != want {
		t.Errorf("req.ContentLength: got %d, want %d", got, want)
	}
	if got, want := req.TransferEncoding, []string{"chunked", "gzip"}; !reflect.DeepEqual(got, want) {
		t.Errorf("req.TransferEncoding: got %v, want %v", got, want)
	}
}

func TestResponseHeader(t *testing.T) {
	res := NewResponse(200, nil, nil)

	h := ResponseHeader(res)

	tt := []struct {
		name  string
		value string
	}{
		{
			name:  "Test-Header",
			value: "true",
		},
		{
			name:  "Content-Length",
			value: "100",
		},
		{
			name:  "Transfer-Encoding",
			value: "chunked",
		},
	}

	for i, tc := range tt {
		if err := h.Set(tc.name, tc.value); err != nil {
			t.Errorf("%d. h.Set(%q, %q): got %v, want no error", i, tc.name, tc.value, err)
		}
	}

	if got, want := res.Header.Get("Test-Header"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Test-Header", got, want)
	}
	if got, want := res.ContentLength, int64(100); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got, want := res.TransferEncoding, []string{"chunked"}; !reflect.DeepEqual(got, want) {
		t.Errorf("res.TransferEncoding: got %v, want %v", got, want)
	}

	if got, want := len(h.Map()), 3; got != want {
		t.Errorf("h.Map(): got %d entries, want %d entries", got, want)
	}

	for n, vs := range h.Map() {
		var want string
		switch n {
		case "Content-Length":
			want = "100"
		case "Transfer-Encoding":
			want = "chunked"
		case "Test-Header":
			want = "true"
		default:
			t.Errorf("h.Map(): got unexpected %s header", n)
		}

		if got := vs[0]; got != want {
			t.Errorf("h.Map(): got %s header with value %s, want value %s", n, got, want)
		}
	}

	for i, tc := range tt {
		got, ok := h.All(tc.name)
		if !ok {
			t.Errorf("%d. h.All(%q): got false, want true", i, tc.name)
		}

		if want := []string{tc.value}; !reflect.DeepEqual(got, want) {
			t.Errorf("%d. h.All(%q): got %v, want %v", i, tc.name, got, want)
		}

		if got, want := h.Get(tc.name), tc.value; got != want {
			t.Errorf("%d. h.Get(%q): got %q, want %q", i, tc.name, got, want)
		}

		h.Del(tc.name)
	}

	if got, want := res.Header.Get("Test-Header"), ""; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Test-Header", got, want)
	}
	if got, want := res.ContentLength, int64(-1); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got := res.TransferEncoding; got != nil {
		t.Errorf("res.TransferEncoding: got %v, want nil", got)
	}

	for i, tc := range tt {
		if got, want := h.Get(tc.name), ""; got != want {
			t.Errorf("%d. h.Get(%q): got %q, want %q", i, tc.name, got, want)
		}

		got, ok := h.All(tc.name)
		if ok {
			t.Errorf("%d. h.All(%q): got ok, want !ok", i, tc.name)
		}
		if got != nil {
			t.Errorf("%d. h.All(%q): got %v, want nil", i, tc.name, got)
		}
	}
}

func TestResponseHeaderAdd(t *testing.T) {
	res := NewResponse(200, nil, nil)

	h := ResponseHeader(res)

	tt := []struct {
		name             string
		values           []string
		errOnSecondValue bool
	}{
		{
			name:   "Test-Header",
			values: []string{"first", "second"},
		},
		{
			name:             "Content-Length",
			values:           []string{"100", "101"},
			errOnSecondValue: true,
		},
		{
			name:   "Transfer-Encoding",
			values: []string{"chunked", "gzip"},
		},
	}

	for i, tc := range tt {
		if err := h.Add(tc.name, tc.values[0]); err != nil {
			t.Errorf("%d. h.Add(%q, %q): got %v, want no error", i, tc.name, tc.values[0], err)
		}
		if err := h.Add(tc.name, tc.values[1]); err != nil && !tc.errOnSecondValue {
			t.Errorf("%d. h.Add(%q, %q): got %v, want no error", i, tc.name, tc.values[1], err)
		}
	}

	if got, want := res.Header["Test-Header"], []string{"first", "second"}; !reflect.DeepEqual(got, want) {
		t.Errorf("res.Header[%q]: got %v, want %v", "Test-Header", got, want)
	}
	if got, want := res.ContentLength, int64(100); got != want {
		t.Errorf("res.ContentLength: got %d, want %d", got, want)
	}
	if got, want := res.TransferEncoding, []string{"chunked", "gzip"}; !reflect.DeepEqual(got, want) {
		t.Errorf("res.TransferEncoding: got %v, want %v", got, want)
	}
}
