package martiantest

import (
	"net/http"

	"github.com/google/martian/v3/proxyutil"
)

// Transport is an http.RoundTripper for testing.
type Transport struct {
	rtfunc func(*http.Request) (*http.Response, error)
}

// NewTransport builds a new transport that will respond with a 200 OK
// response.
func NewTransport() *Transport {
	tr := &Transport{}
	tr.Respond(200)

	return tr
}

// Respond sets the transport to respond with response with statusCode.
func (tr *Transport) Respond(statusCode int) {
	tr.rtfunc = func(req *http.Request) (*http.Response, error) {
		// Force CONNECT requests to 200 to test CONNECT with downstream proxy.
		if req.Method == "CONNECT" {
			statusCode = 200
		}

		res := proxyutil.NewResponse(statusCode, nil, req)

		return res, nil
	}
}

// RespondError sets the transport to respond with an error on round trip.
func (tr *Transport) RespondError(err error) {
	tr.rtfunc = func(*http.Request) (*http.Response, error) {
		return nil, err
	}
}

// CopyHeaders sets the transport to respond with a 200 OK response with
// headers copied from the request to the response verbatim.
func (tr *Transport) CopyHeaders(names ...string) {
	tr.rtfunc = func(req *http.Request) (*http.Response, error) {
		res := proxyutil.NewResponse(200, nil, req)

		for _, n := range names {
			res.Header.Set(n, req.Header.Get(n))
		}

		return res, nil
	}
}

// Func sets the transport to use the rtfunc.
func (tr *Transport) Func(rtfunc func(*http.Request) (*http.Response, error)) {
	tr.rtfunc = rtfunc
}

// RoundTrip runs the stored round trip func and returns the response.
func (tr *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	return tr.rtfunc(req)
}
