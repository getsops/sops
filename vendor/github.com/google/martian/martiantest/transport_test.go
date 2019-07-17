package martiantest

import (
	"errors"
	"net/http"
	"testing"

	"github.com/google/martian/v3/proxyutil"
)

func TestTransport(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	tr := NewTransport()

	res, err := tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("tr.Roundtrip(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}

	// Respond with 301 response.
	tr.Respond(301)
	res, err = tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("tr.Roundtrip(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 301; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}

	// Respond with error.
	trerr := errors.New("transport error")
	tr.RespondError(trerr)

	if _, err := tr.RoundTrip(req); err != trerr {
		t.Fatalf("tr.Roundtrip(): got %v, want %v", err, trerr)
	}

	// Copy headers from request to response.
	req.Header.Set("First-Header", "first")
	req.Header.Set("Second-Header", "second")

	tr.CopyHeaders("First-Header", "Second-Header")

	res, err = tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("tr.Roundtrip(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("First-Header"), "first"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "First-Header", got, want)
	}
	if got, want := res.Header.Get("Second-Header"), "second"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Second-Header", got, want)
	}

	// Custom round trip function.
	tr.Func(func(req *http.Request) (*http.Response, error) {
		res := proxyutil.NewResponse(200, nil, req)
		res.Header.Set("Request-Method", req.Method)

		return res, nil
	})

	res, err = tr.RoundTrip(req)
	if err != nil {
		t.Fatalf("tr.Roundtrip(): got %v, want no error", err)
	}
	res.Body.Close()

	if got, want := res.StatusCode, 200; got != want {
		t.Errorf("res.StatusCode: got %d, want %d", got, want)
	}
	if got, want := res.Header.Get("Request-Method"), "GET"; got != want {
		t.Errorf("res.Header.Get(%q): got %q, want %q", "Request-Method", got, want)
	}
}
