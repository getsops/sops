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

package martian

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/mitm"
	"github.com/google/martian/v3/proxyutil"
)

type tempError struct{}

func (e *tempError) Error() string   { return "temporary" }
func (e *tempError) Timeout() bool   { return true }
func (e *tempError) Temporary() bool { return true }

type timeoutListener struct {
	net.Listener
	errCount int
	err      error
}

func newTimeoutListener(l net.Listener, errCount int) net.Listener {
	return &timeoutListener{
		Listener: l,
		errCount: errCount,
		err:      &tempError{},
	}
}

func (l *timeoutListener) Accept() (net.Conn, error) {
	if l.errCount > 0 {
		l.errCount--
		return nil, l.err
	}

	return l.Listener.Accept()
}

func TestIntegrationTemporaryTimeout(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(200 * time.Millisecond)

	// Start the proxy with a listener that will return a temporary error on
	// Accept() three times.
	go p.Serve(newTimeoutListener(l, 3))

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Connection", "close")

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
}

func TestIntegrationHTTP(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	p.SetRequestModifier(nil)
	p.SetResponseModifier(nil)

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(200 * time.Millisecond)

	tm := martiantest.NewModifier()

	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.Set("martian.test", "true")
	})

	tm.ResponseFunc(func(res *http.Response) {
		ctx := NewContext(res.Request)
		v, _ := ctx.Get("martian.test")

		res.Header.Set("Martian-Test", v.(string))
	})

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	if got, want := res.Header.Get("Martian-Test"), "true"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Martian-Test", got, want)
	}
}

func TestIntegrationHTTP100Continue(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	p.SetTimeout(2 * time.Second)

	sl, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	go func() {
		conn, err := sl.Accept()
		if err != nil {
			log.Errorf("proxy_test: failed to accept connection: %v", err)
			return
		}
		defer conn.Close()

		log.Infof("proxy_test: accepted connection: %s", conn.RemoteAddr())

		req, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			log.Errorf("proxy_test: failed to read request: %v", err)
			return
		}

		if req.Header.Get("Expect") == "100-continue" {
			log.Infof("proxy_test: received 100-continue request")

			conn.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n"))

			log.Infof("proxy_test: sent 100-continue response")
		} else {
			log.Infof("proxy_test: received non 100-continue request")

			res := proxyutil.NewResponse(417, nil, req)
			res.Header.Set("Connection", "close")
			res.Write(conn)
			return
		}

		res := proxyutil.NewResponse(200, req.Body, req)
		res.Header.Set("Connection", "close")
		res.Write(conn)

		log.Infof("proxy_test: sent 200 response")
	}()

	tm := martiantest.NewModifier()
	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	host := sl.Addr().String()
	raw := fmt.Sprintf("POST http://%s/ HTTP/1.1\r\n"+
		"Host: %s\r\n"+
		"Content-Length: 12\r\n"+
		"Expect: 100-continue\r\n\r\n", host, host)

	if _, err := conn.Write([]byte(raw)); err != nil {
		t.Fatalf("conn.Write(headers): got %v, want no error", err)
	}

	go func() {
		select {
		case <-time.After(time.Second):
			conn.Write([]byte("body content"))
		}
	}()

	res, err := http.ReadResponse(bufio.NewReader(conn), nil)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	got, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll(): got %v, want no error", err)
	}

	if want := []byte("body content"); !bytes.Equal(got, want) {
		t.Errorf("res.Body: got %q, want %q", got, want)
	}

	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
}

func TestIntegrationHTTPDownstreamProxy(t *testing.T) {
	t.Parallel()

	// Start first proxy to use as downstream.
	dl, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	downstream := NewProxy()
	defer downstream.Close()

	dtr := martiantest.NewTransport()
	dtr.Respond(299)
	downstream.SetRoundTripper(dtr)
	downstream.SetTimeout(600 * time.Millisecond)

	go downstream.Serve(dl)

	// Start second proxy as upstream proxy, will write to downstream proxy.
	ul, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	upstream := NewProxy()
	defer upstream.Close()

	// Set upstream proxy's downstream proxy to the host:port of the first proxy.
	upstream.SetDownstreamProxy(&url.URL{
		Host: dl.Addr().String(),
	})
	upstream.SetTimeout(600 * time.Millisecond)

	go upstream.Serve(ul)

	// Open connection to upstream proxy.
	conn, err := net.Dial("tcp", ul.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	// Response from downstream proxy.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, 299; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}
}

func TestIntegrationHTTPDownstreamProxyError(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	// Set proxy's downstream proxy to invalid host:port to force failure.
	p.SetDownstreamProxy(&url.URL{
		Host: "[::]:0",
	})
	p.SetTimeout(600 * time.Millisecond)

	tm := martiantest.NewModifier()
	reserr := errors.New("response error")
	tm.ResponseError(reserr)

	p.SetResponseModifier(tm)

	go p.Serve(l)

	// Open connection to upstream proxy.
	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", "//example.com:443", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:443 HTTP/1.1
	// Host: example.com
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// Response from upstream proxy, assuming downstream proxy failed to CONNECT.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, 502; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header["Warning"][1], reserr.Error(); !strings.Contains(got, want) {
		t.Errorf("res.Header.get(%q): got %q, want to contain %q", "Warning", got, want)
	}
}

func TestIntegrationTLSHandshakeErrorCallback(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	// Test TLS server.
	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}
	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}

	var herr error
	mc.SetHandshakeErrorCallback(func(_ *http.Request, err error) { herr = fmt.Errorf("handshake error") })
	p.SetMITM(mc)

	tl, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("tls.Listen(): got %v, want no error", err)
	}
	tl = tls.NewListener(tl, mc.TLS())

	go http.Serve(tl, http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(200)
		}))

	tm := martiantest.NewModifier()

	// Force the CONNECT request to dial the local TLS server.
	tm.RequestFunc(func(req *http.Request) {
		req.URL.Host = tl.Addr().String()
	})

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", "//example.com:443", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:443 HTTP/1.1
	// Host: example.com
	//
	// Rewritten to CONNECT to host:port in CONNECT request modifier.
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// CONNECT response after establishing tunnel.
	if _, err := http.ReadResponse(bufio.NewReader(conn), req); err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	tlsconn := tls.Client(conn, &tls.Config{
		ServerName: "example.com",
		// Client has no cert so it will get "x509: certificate signed by unknown authority" from the
		// handshake and send "remote error: bad certificate" to the server.
		RootCAs: x509.NewCertPool(),
	})
	defer tlsconn.Close()

	req, err = http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Connection", "close")

	if got, want := req.Write(tlsconn), "x509: certificate signed by unknown authority"; !strings.Contains(got.Error(), want) {
		t.Fatalf("Got incorrect error from Client Handshake(), got: %v, want: %v", got, want)
	}

	// TODO: herr is not being asserted against. It should be pushed on to a channel
	// of err, and the assertion should pull off of it and assert. That design resulted in the test
	// hanging for unknown reasons.
	t.Skip("skipping assertion of handshake error callback error due to mysterious deadlock")
	if got, want := herr, "remote error: bad certificate"; !strings.Contains(got.Error(), want) {
		t.Fatalf("Got incorrect error from Server Handshake(), got: %v, want: %v", got, want)
	}
}

func TestIntegrationConnect(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	// Test TLS server.
	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}
	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}

	tl, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("tls.Listen(): got %v, want no error", err)
	}
	tl = tls.NewListener(tl, mc.TLS())

	go http.Serve(tl, http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(299)
		}))

	tm := martiantest.NewModifier()
	reqerr := errors.New("request error")
	reserr := errors.New("response error")

	// Force the CONNECT request to dial the local TLS server.
	tm.RequestFunc(func(req *http.Request) {
		req.URL.Host = tl.Addr().String()
	})

	tm.RequestError(reqerr)
	tm.ResponseError(reserr)

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", "//example.com:443", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:443 HTTP/1.1
	// Host: example.com
	//
	// Rewritten to CONNECT to host:port in CONNECT request modifier.
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// CONNECT response after establishing tunnel.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
	if got, want := res.Header.Get("Warning"), reserr.Error(); !strings.Contains(got, want) {
		t.Errorf("res.Header.Get(%q): got %q, want to contain %q", "Warning", got, want)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca)

	tlsconn := tls.Client(conn, &tls.Config{
		ServerName: "example.com",
		RootCAs:    roots,
	})
	defer tlsconn.Close()

	req, err = http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}
	req.Header.Set("Connection", "close")

	// GET / HTTP/1.1
	// Host: example.com
	// Connection: close
	if err := req.Write(tlsconn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	res, err = http.ReadResponse(bufio.NewReader(tlsconn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 299; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Warning"), reserr.Error(); strings.Contains(got, want) {
		t.Errorf("res.Header.Get(%q): got %s, want to not contain %s", "Warning", got, want)
	}
}

func TestIntegrationConnectDownstreamProxy(t *testing.T) {
	t.Parallel()

	// Start first proxy to use as downstream.
	dl, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	downstream := NewProxy()
	defer downstream.Close()

	dtr := martiantest.NewTransport()
	dtr.Respond(299)
	downstream.SetRoundTripper(dtr)

	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", 2*time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}

	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}
	downstream.SetMITM(mc)

	go downstream.Serve(dl)

	// Start second proxy as upstream proxy, will CONNECT to downstream proxy.
	ul, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	upstream := NewProxy()
	defer upstream.Close()

	// Set upstream proxy's downstream proxy to the host:port of the first proxy.
	upstream.SetDownstreamProxy(&url.URL{
		Host: dl.Addr().String(),
	})

	go upstream.Serve(ul)

	// Open connection to upstream proxy.
	conn, err := net.Dial("tcp", ul.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", "//example.com:443", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:443 HTTP/1.1
	// Host: example.com
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// Response from downstream proxy starting MITM.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca)

	tlsconn := tls.Client(conn, &tls.Config{
		// Validate the hostname.
		ServerName: "example.com",
		// The certificate will have been MITM'd, verify using the MITM CA
		// certificate.
		RootCAs: roots,
	})
	defer tlsconn.Close()

	req, err = http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET / HTTP/1.1
	// Host: example.com
	if err := req.Write(tlsconn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// Response from MITM in downstream proxy.
	res, err = http.ReadResponse(bufio.NewReader(tlsconn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 299; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}
}

func TestIntegrationMITM(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	tr := martiantest.NewTransport()
	tr.Func(func(req *http.Request) (*http.Response, error) {
		res := proxyutil.NewResponse(200, nil, req)
		res.Header.Set("Request-Scheme", req.URL.Scheme)

		return res, nil
	})

	p.SetRoundTripper(tr)
	p.SetTimeout(600 * time.Millisecond)

	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", 2*time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}

	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}
	p.SetMITM(mc)

	tm := martiantest.NewModifier()
	reqerr := errors.New("request error")
	reserr := errors.New("response error")
	tm.RequestError(reqerr)
	tm.ResponseError(reserr)

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", "//example.com:443", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:443 HTTP/1.1
	// Host: example.com
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// Response MITM'd from proxy.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	if got, want := res.StatusCode, 200; got != want {

		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Warning"), reserr.Error(); !strings.Contains(got, want) {
		t.Errorf("res.Header.Get(%q): got %q, want to contain %q", "Warning", got, want)
	}

	roots := x509.NewCertPool()
	roots.AddCert(ca)

	tlsconn := tls.Client(conn, &tls.Config{
		ServerName: "example.com",
		RootCAs:    roots,
	})
	defer tlsconn.Close()

	req, err = http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET / HTTP/1.1
	// Host: example.com
	if err := req.Write(tlsconn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// Response from MITM proxy.
	res, err = http.ReadResponse(bufio.NewReader(tlsconn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Request-Scheme"), "https"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Request-Scheme", got, want)
	}
	if got, want := res.Header.Get("Warning"), reserr.Error(); !strings.Contains(got, want) {
		t.Errorf("res.Header.Get(%q): got %q, want to contain %q", "Warning", got, want)
	}
}

func TestIntegrationTransparentHTTP(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(200 * time.Millisecond)

	tm := martiantest.NewModifier()
	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET / HTTP/1.1
	// Host: www.example.com
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	if !tm.RequestModified() {
		t.Error("tm.RequestModified(): got false, want true")
	}
	if !tm.ResponseModified() {
		t.Error("tm.ResponseModified(): got false, want true")
	}
}

func TestIntegrationTransparentMITM(t *testing.T) {
	t.Parallel()

	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", 2*time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}

	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}

	// Start TLS listener with config that will generate certificates based on
	// SNI from connection.
	//
	// BUG: tls.Listen will not accept a tls.Config where Certificates is empty,
	// even though it is supported by tls.Server when GetCertificate is not nil.
	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}
	l = tls.NewListener(l, mc.TLS())

	p := NewProxy()
	defer p.Close()

	tr := martiantest.NewTransport()
	tr.Func(func(req *http.Request) (*http.Response, error) {
		res := proxyutil.NewResponse(200, nil, req)
		res.Header.Set("Request-Scheme", req.URL.Scheme)

		return res, nil
	})

	p.SetRoundTripper(tr)

	tm := martiantest.NewModifier()
	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(l)

	roots := x509.NewCertPool()
	roots.AddCert(ca)

	tlsconn, err := tls.Dial("tcp", l.Addr().String(), &tls.Config{
		// Verify the hostname is example.com.
		ServerName: "example.com",
		// The certificate will have been generated during MITM, so we need to
		// verify it with the generated CA certificate.
		RootCAs: roots,
	})
	if err != nil {
		t.Fatalf("tls.Dial(): got %v, want no error", err)
	}
	defer tlsconn.Close()

	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// Write Encrypted request directly, no CONNECT.
	// GET / HTTP/1.1
	// Host: example.com
	if err := req.Write(tlsconn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	res, err := http.ReadResponse(bufio.NewReader(tlsconn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Request-Scheme"), "https"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Request-Scheme", got, want)
	}

	if !tm.RequestModified() {
		t.Errorf("tm.RequestModified(): got false, want true")
	}
	if !tm.ResponseModified() {
		t.Errorf("tm.ResponseModified(): got false, want true")
	}
}

func TestIntegrationFailedRoundTrip(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	tr := martiantest.NewTransport()
	trerr := errors.New("round trip error")
	tr.RespondError(trerr)
	p.SetRoundTripper(tr)
	p.SetTimeout(200 * time.Millisecond)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	// Response from failed round trip.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 502; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}

	if got, want := res.Header.Get("Warning"), trerr.Error(); !strings.Contains(got, want) {
		t.Errorf("res.Header.Get(%q): got %q, want to contain %q", "Warning", got, want)
	}
}

func TestIntegrationSkipRoundTrip(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	// Transport will be skipped, no 500.
	tr := martiantest.NewTransport()
	tr.Respond(500)
	p.SetRoundTripper(tr)
	p.SetTimeout(200 * time.Millisecond)

	tm := martiantest.NewModifier()
	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()
	})
	p.SetRequestModifier(tm)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	// Response from skipped round trip.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	defer res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
}

func TestHTTPThroughConnectWithMITM(t *testing.T) {
	t.Parallel()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	p := NewProxy()
	defer p.Close()

	tm := martiantest.NewModifier()
	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()

		if req.Method != "GET" && req.Method != "CONNECT" {
			t.Errorf("unexpected method on request handler: %v", req.Method)
		}
	})
	p.SetRequestModifier(tm)

	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", 2*time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}

	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}
	p.SetMITM(mc)

	go p.Serve(l)

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", "//example.com:80", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:80 HTTP/1.1
	// Host: example.com
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}

	// Response skipped round trip.
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}

	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	// Response from skipped round trip.
	res, err = http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}

	req, err = http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// GET http://example.com/ HTTP/1.1
	// Host: example.com
	if err := req.WriteProxy(conn); err != nil {
		t.Fatalf("req.WriteProxy(): got %v, want no error", err)
	}

	// Response from skipped round trip.
	res, err = http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
}

func TestServerClosesConnection(t *testing.T) {
	t.Parallel()

	dstl, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("Failed to create http listener: %v", err)
	}
	defer dstl.Close()

	go func() {
		t.Logf("Waiting for server side connection")
		conn, err := dstl.Accept()
		if err != nil {
			t.Fatalf("Got error while accepting connection on destination listener: %v", err)
		}
		t.Logf("Accepted server side connection")

		buf := make([]byte, 16384)
		if _, err := conn.Read(buf); err != nil {
			t.Fatalf("Error reading: %v", err)
		}

		_, err = conn.Write([]byte("HTTP/1.1 301 MOVED PERMANENTLY\r\n" +
			"Server:  \r\n" +
			"Date:  \r\n" +
			"Referer:  \r\n" +
			"Location: http://www.foo.com/\r\n" +
			"Content-type: text/html\r\n" +
			"Connection: close\r\n\r\n"))
		if err != nil {
			t.Fatalf("Got error while writting to connection on destination listener: %v", err)
		}
		conn.Close()
	}()

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	ca, priv, err := mitm.NewAuthority("martian.proxy", "Martian Authority", 2*time.Hour)
	if err != nil {
		t.Fatalf("mitm.NewAuthority(): got %v, want no error", err)
	}

	mc, err := mitm.NewConfig(ca, priv)
	if err != nil {
		t.Fatalf("mitm.NewConfig(): got %v, want no error", err)
	}
	p := NewProxy()
	p.SetMITM(mc)
	defer p.Close()

	// Start the proxy with a listener that will return a temporary error on
	// Accept() three times.
	go p.Serve(newTimeoutListener(l, 3))

	conn, err := net.Dial("tcp", l.Addr().String())
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}
	defer conn.Close()

	req, err := http.NewRequest("CONNECT", fmt.Sprintf("//%s", dstl.Addr().String()), nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	// CONNECT example.com:443 HTTP/1.1
	// Host: example.com
	if err := req.Write(conn); err != nil {
		t.Fatalf("req.Write(): got %v, want no error", err)
	}
	res, err := http.ReadResponse(bufio.NewReader(conn), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	res.Body.Close()

	_, err = conn.Write([]byte("GET / HTTP/1.1\r\n" +
		"User-Agent: curl/7.35.0\r\n" +
		fmt.Sprintf("Host: %s\r\n", dstl.Addr()) +
		"Accept: */*\r\n\r\n"))
	if err != nil {
		t.Fatalf("Error while writing GET request: %v", err)
	}

	res, err = http.ReadResponse(bufio.NewReader(io.TeeReader(conn, os.Stderr)), req)
	if err != nil {
		t.Fatalf("http.ReadResponse(): got %v, want no error", err)
	}
	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("error while ReadAll: %v", err)
	}
	defer res.Body.Close()
}

// TestRacyClose checks that creating a proxy, serving from it, and closing
// it in rapid succession doesn't result in race warnings.
// See https://github.com/google/martian/issues/286.
func TestRacyClose(t *testing.T) {
	t.Parallel()

	log.SetLevel(log.Silent) // avoid "failed to accept" messages because we close l
	openAndConnect := func() {
		l, err := net.Listen("tcp", "[::]:0")
		if err != nil {
			t.Fatalf("net.Listen(): got %v, want no error", err)
		}
		defer l.Close() // to make p.Serve exit

		p := NewProxy()
		go p.Serve(l)
		defer p.Close()

		conn, err := net.Dial("tcp", l.Addr().String())
		if err != nil {
			t.Fatalf("net.Dial(): got %v, want no error", err)
		}
		defer conn.Close()
	}

	// Repeat a bunch of times to make failures more repeatable.
	for i := 0; i < 100; i++ {
		openAndConnect()
	}
}
