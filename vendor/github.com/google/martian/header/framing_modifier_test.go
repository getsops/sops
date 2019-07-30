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

package header

import (
	"net/http"
	"reflect"
	"testing"
)

func TestBadFramingMultipleContentLengths(t *testing.T) {
	m := NewBadFramingModifier()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header["Content-Length"] = []string{"42", "42, 42"}

	if err := m.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}
	if got, want := req.Header["Content-Length"], []string{"42"}; !reflect.DeepEqual(got, want) {
		t.Errorf("req.Header[%q]: got %v, want %v", "Content-Length", got, want)
	}

	req.Header["Content-Length"] = []string{"42", "32, 42"}
	if err := m.ModifyRequest(req); err == nil {
		t.Error("ModifyRequest(): got nil, want error")
	}
}

func TestBadFramingTransferEncodingAndContentLength(t *testing.T) {
	m := NewBadFramingModifier()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header["Transfer-Encoding"] = []string{"gzip, chunked"}
	req.Header["Content-Length"] = []string{"42"}

	if err := m.ModifyRequest(req); err != nil {
		t.Errorf("ModifyRequest(): got %v, want no error", err)
	}
	if _, ok := req.Header["Content-Length"]; ok {
		t.Fatalf("req.Header[%q]: got ok, want !ok", "Content-Length")
	}

	req.Header.Set("Transfer-Encoding", "gzip, identity")
	req.Header.Del("Content-Length")
	if err := m.ModifyRequest(req); err == nil {
		t.Error("ModifyRequest(): got nil, want error")
	}
}
