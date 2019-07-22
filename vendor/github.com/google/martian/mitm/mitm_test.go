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

package mitm

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestMITM(t *testing.T) {
	ca, priv, err := NewAuthority("martian.proxy", "Martian Authority", 24*time.Hour)
	if err != nil {
		t.Fatalf("NewAuthority(): got %v, want no error", err)
	}

	c, err := NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("NewConfig(): got %v, want no error", err)
	}

	c.SetValidity(20 * time.Hour)
	c.SetOrganization("Test Organization")

	protos := []string{"http/1.1"}

	conf := c.TLS()
	if got := conf.NextProtos; !reflect.DeepEqual(got, protos) {
		t.Errorf("conf.NextProtos: got %v, want %v", got, protos)
	}
	if conf.InsecureSkipVerify {
		t.Error("conf.InsecureSkipVerify: got true, want false")
	}

	// Simulate a TLS connection without SNI.
	clientHello := &tls.ClientHelloInfo{
		ServerName: "",
	}

	if _, err := conf.GetCertificate(clientHello); err == nil {
		t.Fatal("conf.GetCertificate(): got nil, want error")
	}

	// Simulate a TLS connection with SNI.
	clientHello.ServerName = "example.com"

	tlsc, err := conf.GetCertificate(clientHello)
	if err != nil {
		t.Fatalf("conf.GetCertificate(): got %v, want no error", err)
	}

	x509c := tlsc.Leaf
	if got, want := x509c.Subject.CommonName, "example.com"; got != want {
		t.Errorf("x509c.Subject.CommonName: got %q, want %q", got, want)
	}

	c.SkipTLSVerify(true)

	conf = c.TLSForHost("example.com")
	if got := conf.NextProtos; !reflect.DeepEqual(got, protos) {
		t.Errorf("conf.NextProtos: got %v, want %v", got, protos)
	}
	if !conf.InsecureSkipVerify {
		t.Error("conf.InsecureSkipVerify: got false, want true")
	}

	// Set SNI, takes precedence over host.
	clientHello.ServerName = "google.com"
	tlsc, err = conf.GetCertificate(clientHello)
	if err != nil {
		t.Fatalf("conf.GetCertificate(): got %v, want no error", err)
	}

	x509c = tlsc.Leaf
	if got, want := x509c.Subject.CommonName, "google.com"; got != want {
		t.Errorf("x509c.Subject.CommonName: got %q, want %q", got, want)
	}

	// Reset SNI to fallback to hostname.
	clientHello.ServerName = ""
	tlsc, err = conf.GetCertificate(clientHello)
	if err != nil {
		t.Fatalf("conf.GetCertificate(): got %v, want no error", err)
	}

	x509c = tlsc.Leaf
	if got, want := x509c.Subject.CommonName, "example.com"; got != want {
		t.Errorf("x509c.Subject.CommonName: got %q, want %q", got, want)
	}
}

func TestCert(t *testing.T) {
	ca, priv, err := NewAuthority("martian.proxy", "Martian Authority", 24*time.Hour)
	if err != nil {
		t.Fatalf("NewAuthority(): got %v, want no error", err)
	}

	c, err := NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("NewConfig(): got %v, want no error", err)
	}

	tlsc, err := c.cert("example.com")
	if err != nil {
		t.Fatalf("c.cert(%q): got %v, want no error", "example.com:8080", err)
	}

	if tlsc.Certificate == nil {
		t.Error("tlsc.Certificate: got nil, want certificate bytes")
	}
	if tlsc.PrivateKey == nil {
		t.Error("tlsc.PrivateKey: got nil, want private key")
	}

	x509c := tlsc.Leaf
	if x509c == nil {
		t.Fatal("x509c: got nil, want *x509.Certificate")
	}

	if got := x509c.SerialNumber; got.Cmp(MaxSerialNumber) >= 0 {
		t.Errorf("x509c.SerialNumber: got %v, want <= MaxSerialNumber", got)
	}
	if got, want := x509c.Subject.CommonName, "example.com"; got != want {
		t.Errorf("X509c.Subject.CommonName: got %q, want %q", got, want)
	}
	if err := x509c.VerifyHostname("example.com"); err != nil {
		t.Errorf("x509c.VerifyHostname(%q): got %v, want no error", "example.com", err)
	}

	if got, want := x509c.Subject.Organization, []string{"Martian Proxy"}; !reflect.DeepEqual(got, want) {
		t.Errorf("x509c.Subject.Organization: got %v, want %v", got, want)
	}

	if got := x509c.SubjectKeyId; got == nil {
		t.Error("x509c.SubjectKeyId: got nothing, want key ID")
	}
	if !x509c.BasicConstraintsValid {
		t.Error("x509c.BasicConstraintsValid: got false, want true")
	}

	if got, want := x509c.KeyUsage, x509.KeyUsageKeyEncipherment; got&want == 0 {
		t.Error("x509c.KeyUsage: got nothing, want to include x509.KeyUsageKeyEncipherment")
	}
	if got, want := x509c.KeyUsage, x509.KeyUsageDigitalSignature; got&want == 0 {
		t.Error("x509c.KeyUsage: got nothing, want to include x509.KeyUsageDigitalSignature")
	}

	want := []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
	if got := x509c.ExtKeyUsage; !reflect.DeepEqual(got, want) {
		t.Errorf("x509c.ExtKeyUsage: got %v, want %v", got, want)
	}

	if got, want := x509c.DNSNames, []string{"example.com"}; !reflect.DeepEqual(got, want) {
		t.Errorf("x509c.DNSNames: got %v, want %v", got, want)
	}

	before := time.Now().Add(-2 * time.Hour)
	if got := x509c.NotBefore; before.After(got) {
		t.Errorf("x509c.NotBefore: got %v, want after %v", got, before)
	}

	after := time.Now().Add(2 * time.Hour)
	if got := x509c.NotAfter; !after.After(got) {
		t.Errorf("x509c.NotAfter: got %v, want before %v", got, want)
	}

	// Retrieve cached certificate.
	tlsc2, err := c.cert("example.com")
	if err != nil {
		t.Fatalf("c.cert(%q): got %v, want no error", "example.com", err)
	}
	if tlsc != tlsc2 {
		t.Error("tlsc2: got new certificate, want cached certificate")
	}

	// TLS certificate for IP.
	tlsc, err = c.cert("10.0.0.1:8227")
	if err != nil {
		t.Fatalf("c.cert(%q): got %v, want no error", "10.0.0.1:8227", err)
	}
	x509c = tlsc.Leaf

	if got, want := len(x509c.IPAddresses), 1; got != want {
		t.Fatalf("len(x509c.IPAddresses): got %d, want %d", got, want)
	}

	if got, want := x509c.IPAddresses[0], net.ParseIP("10.0.0.1"); !got.Equal(want) {
		t.Fatalf("x509c.IPAddresses: got %v, want %v", got, want)
	}
}
