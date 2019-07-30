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

package martianhttp

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/martian/v3/mitm"
)

func TestAuthorityHandler(t *testing.T) {
	ca, _, err := mitm.NewAuthority("martian.proxy", "Martian Authority", time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}

	rw := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/martian/authority.cer", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	h := NewAuthorityHandler(ca)
	h.ServeHTTP(rw, req)

	if got, want := rw.Code, 200; got != want {
		t.Errorf("rw.Code: got %d, want %d", got, want)
	}
	if got, want := rw.Header().Get("Content-Type"), "application/x-x509-ca-cert"; got != want {
		t.Errorf("rw.Header().Get(%q): got %q, want %q", "Content-Type", got, want)
	}

	blk, _ := pem.Decode(rw.Body.Bytes())
	if got, want := blk.Type, "CERTIFICATE"; got != want {
		t.Errorf("rw.Body: got PEM type %q, want %q", got, want)
	}

	cert, err := x509.ParseCertificate(blk.Bytes)
	if err != nil {
		t.Fatalf("x509.ParseCertificate(res.Body): got %v, want no error", err)
	}
	if got, want := cert.Subject.CommonName, "martian.proxy"; got != want {
		t.Errorf("cert.Subject.CommonName: got %q, want %q", got, want)
	}
}
