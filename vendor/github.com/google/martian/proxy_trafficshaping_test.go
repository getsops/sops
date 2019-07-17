package martian

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/martian/v3/log"
	"github.com/google/martian/v3/martiantest"
	"github.com/google/martian/v3/trafficshape"
)

// Tests that sending data of length 600 bytes with max bandwidth of 100 bytes/s takes
// atleast 4.9s. Uses the Close Connection action to immediately close the connection
// upon the proxy writing 600 bytes. (4.9s ~ 5s = 600b /100b/s - 1s)
func TestConstantThrottleAndClose(t *testing.T) {
	log.SetLevel(log.Info)

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tsl := trafficshape.NewListener(l)
	tsh := trafficshape.NewHandler(tsl)

	// This is the data to be sent.
	testString := strings.Repeat("0", 600)

	// Traffic shaping config request.
	jsonString :=
		`{
				"trafficshape": {
						"shapes": [
							{
								"url_regex": "http://example/example",
								"throttles": [
									{
										"bytes": "0-",
										"bandwidth": 100
									}
								],
								"close_connections": [
									{
										"byte": 600,
										"count": 1
									}
								]
						}
						]
				}
		}`

	tsReq, err := http.NewRequest("POST", "test", bytes.NewBufferString(jsonString))
	rw := httptest.NewRecorder()
	tsh.ServeHTTP(rw, tsReq)
	res := rw.Result()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	p := NewProxy()
	defer p.Close()

	p.SetRequestModifier(nil)
	p.SetResponseModifier(nil)

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(15 * time.Second)

	tm := martiantest.NewModifier()

	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()
	})

	tm.ResponseFunc(func(res *http.Response) {
		res.StatusCode = http.StatusOK
		res.Body = ioutil.NopCloser(bytes.NewBufferString(testString))
	})

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(tsl)

	c1 := make(chan string)
	conn, err := net.Dial("tcp", l.Addr().String())
	defer conn.Close()
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}

	go func() {
		req, err := http.NewRequest("GET", "http://example/example", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		if err := req.WriteProxy(conn); err != nil {
			t.Fatalf("req.WriteProxy(): got %v, want no error", err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			t.Fatalf("http.ReadResponse(): got %v, want no error", err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		bodystr := string(body)
		c1 <- bodystr
	}()

	var bodystr string
	select {
	case bodystringc := <-c1:
		t.Errorf("took < 4.9s, should take at least 4.9s")
		bodystr = bodystringc
	case <-time.After(4900 * time.Millisecond):
		bodystringc := <-c1
		bodystr = bodystringc
	}

	if bodystr != testString {
		t.Errorf("res.Body: got %s, want %s", bodystr, testString)
	}
}

// Tests that sleeping for 5s and then closing the connection
// upon reading 200 bytes, with a bandwidth of 5000 bytes/s
// takes at least 4.9s, and results in a correctly trimmed
// response body. (200 0s instead of 500 0s)
func TestSleepAndClose(t *testing.T) {
	log.SetLevel(log.Info)

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tsl := trafficshape.NewListener(l)
	tsh := trafficshape.NewHandler(tsl)

	// This is the data to be sent.
	testString := strings.Repeat("0", 500)

	// Traffic shaping config request.
	jsonString :=
		`{
				"trafficshape": {
						"shapes": [
							{
								"url_regex": "http://example/example",
								"throttles": [
									{
										"bytes": "0-",
										"bandwidth": 5000
									}
								],
								"halts": [
									{
										"byte": 100,
										"duration": 5000,
										"count": 1
									}
								],
								"close_connections": [
									{
										"byte": 200,
										"count": 1
									}
								]
						}
						]
				}
		}`

	tsReq, err := http.NewRequest("POST", "test", bytes.NewBufferString(jsonString))
	rw := httptest.NewRecorder()
	tsh.ServeHTTP(rw, tsReq)
	res := rw.Result()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	p := NewProxy()
	defer p.Close()

	p.SetRequestModifier(nil)
	p.SetResponseModifier(nil)

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(15 * time.Second)

	tm := martiantest.NewModifier()

	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()
	})

	tm.ResponseFunc(func(res *http.Response) {
		res.StatusCode = http.StatusOK
		res.Body = ioutil.NopCloser(bytes.NewBufferString(testString))
	})

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(tsl)

	c1 := make(chan string)
	conn, err := net.Dial("tcp", l.Addr().String())
	defer conn.Close()
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}

	go func() {
		req, err := http.NewRequest("GET", "http://example/example", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		if err := req.WriteProxy(conn); err != nil {
			t.Fatalf("req.WriteProxy(): got %v, want no error", err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			t.Fatalf("http.ReadResponse(): got %v, want no error", err)
		}
		body, _ := ioutil.ReadAll(res.Body)
		bodystr := string(body)
		c1 <- bodystr
	}()

	var bodystr string
	select {
	case bodystringc := <-c1:
		t.Errorf("took < 4.9s, should take at least 4.9s")
		bodystr = bodystringc
	case <-time.After(4900 * time.Millisecond):
		bodystringc := <-c1
		bodystr = bodystringc
	}

	if want := strings.Repeat("0", 200); bodystr != want {
		t.Errorf("res.Body: got %s, want %s", bodystr, want)
	}
}

// Similar to TestConstantThrottleAndClose, except that it applies
// the throttle only in a specific byte range, and modifies the
// the response to lie in the byte range.
func TestConstantThrottleAndCloseByteRange(t *testing.T) {
	log.SetLevel(log.Info)

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tsl := trafficshape.NewListener(l)
	tsh := trafficshape.NewHandler(tsl)

	// This is the data to be sent.
	testString := strings.Repeat("0", 600)

	// Traffic shaping config request.
	jsonString :=
		`{
				"trafficshape": {
						"shapes": [
							{
								"url_regex": "http://example/example",
								"throttles": [
									{
										"bytes": "500-",
										"bandwidth": 100
									}
								],
								"close_connections": [
									{
										"byte": 1100,
										"count": 1
									}
								]
						}
						]
				}
		}`

	tsReq, err := http.NewRequest("POST", "test", bytes.NewBufferString(jsonString))
	rw := httptest.NewRecorder()
	tsh.ServeHTTP(rw, tsReq)
	res := rw.Result()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	p := NewProxy()
	defer p.Close()

	p.SetRequestModifier(nil)
	p.SetResponseModifier(nil)

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(15 * time.Second)

	tm := martiantest.NewModifier()

	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()
	})

	tm.ResponseFunc(func(res *http.Response) {
		res.StatusCode = http.StatusPartialContent
		res.Body = ioutil.NopCloser(bytes.NewBufferString(testString))
		res.Header.Set("Content-Range", "bytes 500-1100/1100")
	})

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(tsl)

	c1 := make(chan string)
	conn, err := net.Dial("tcp", l.Addr().String())
	defer conn.Close()
	if err != nil {
		t.Fatalf("net.Dial(): got %v, want no error", err)
	}

	go func() {
		req, err := http.NewRequest("GET", "http://example/example", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		if err := req.WriteProxy(conn); err != nil {
			t.Fatalf("req.WriteProxy(): got %v, want no error", err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			t.Fatalf("http.ReadResponse(): got %v, want no error", err)
		}

		body, _ := ioutil.ReadAll(res.Body)
		bodystr := string(body)
		c1 <- bodystr
	}()

	var bodystr string
	select {
	case bodystringc := <-c1:
		t.Errorf("took < 4.9s, should take at least 4.9s")
		bodystr = bodystringc
	case <-time.After(4900 * time.Millisecond):
		bodystringc := <-c1
		bodystr = bodystringc
	}

	if bodystr != testString {
		t.Errorf("res.Body: got %s, want %s", bodystr, testString)
	}
}

// Opens up 5 concurrent connections, and sets the
// max global bandwidth for the url regex to be 250b/s.
// Every connection tries to read 500b of data, but since
// the global bandwidth for the particular regex is 250,
// it should take at least 5 * 500b / 250b/s -1s = 9s to read
// everything.
func TestMaxBandwidth(t *testing.T) {
	log.SetLevel(log.Info)

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tsl := trafficshape.NewListener(l)
	tsh := trafficshape.NewHandler(tsl)

	// This is the data to be sent.
	testString := strings.Repeat("0", 500)

	// Traffic shaping config request.
	jsonString :=
		`{
				"trafficshape": {
						"shapes": [
							{
								"url_regex": "http://example/example",
								"max_global_bandwidth": 250,
								"close_connections": [
									{
										"byte": 500,
										"count": 5
									}
								]
						}
						]
				}
		}`

	tsReq, err := http.NewRequest("POST", "test", bytes.NewBufferString(jsonString))
	rw := httptest.NewRecorder()
	tsh.ServeHTTP(rw, tsReq)
	res := rw.Result()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	p := NewProxy()
	defer p.Close()

	p.SetRequestModifier(nil)
	p.SetResponseModifier(nil)

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(20 * time.Second)

	tm := martiantest.NewModifier()

	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()
	})

	tm.ResponseFunc(func(res *http.Response) {
		res.StatusCode = http.StatusOK
		res.Body = ioutil.NopCloser(bytes.NewBufferString(testString))
	})

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(tsl)

	numChannels := 5

	channels := make([]chan string, numChannels)

	for i := 0; i < numChannels; i++ {
		channels[i] = make(chan string)
	}

	for i := 0; i < numChannels; i++ {
		go func(i int) {
			conn, err := net.Dial("tcp", l.Addr().String())
			defer conn.Close()
			if err != nil {
				t.Fatalf("net.Dial(): got %v, want no error", err)
			}
			req, err := http.NewRequest("GET", "http://example/example", nil)
			if err != nil {
				t.Fatalf("http.NewRequest(): got %v, want no error", err)
			}

			if err := req.WriteProxy(conn); err != nil {
				t.Fatalf("req.WriteProxy(): got %v, want no error", err)
			}

			res, err := http.ReadResponse(bufio.NewReader(conn), req)
			if err != nil {
				t.Fatalf("http.ReadResponse(): got %v, want no error", err)
			}

			body, _ := ioutil.ReadAll(res.Body)
			bodystr := string(body)

			if i != 0 {
				<-channels[i-1]
			}

			channels[i] <- bodystr
		}(i)
	}

	var bodystr string
	select {
	case bodystringc := <-channels[numChannels-1]:
		t.Errorf("took < 8.9s, should take at least 8.9s")
		bodystr = bodystringc
	case <-time.After(8900 * time.Millisecond):
		bodystringc := <-channels[numChannels-1]
		bodystr = bodystringc
	}

	if bodystr != testString {
		t.Errorf("res.Body: got %s, want %s", bodystr, testString)
	}
}

// Makes 2 requests, with the first one having a byte range starting
// at  byte 250, and adds a close connection action at byte 450.
// The first request should hit the action sooner,
// and delete it. The second request should read the whole
// data (500b)
func TestConcurrentResponseActions(t *testing.T) {
	log.SetLevel(log.Info)

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tsl := trafficshape.NewListener(l)
	tsh := trafficshape.NewHandler(tsl)

	// This is the data to be sent.
	testString := strings.Repeat("0", 500)

	// Traffic shaping config request.
	jsonString :=
		`{
				"trafficshape": {
						"shapes": [
							{
								"url_regex": "http://example/example",
								"throttles": [
									{
										"bytes": "-",
										"bandwidth": 250
									}
								],
								"close_connections": [
									{
										"byte": 450,
										"count": 1
									},
									{
										"byte": 500,
										"count": 1
									}
								]
						}
						]
				}
		}`

	tsReq, err := http.NewRequest("POST", "test", bytes.NewBufferString(jsonString))
	rw := httptest.NewRecorder()
	tsh.ServeHTTP(rw, tsReq)
	res := rw.Result()

	if got, want := res.StatusCode, 200; got != want {
		t.Fatalf("res.StatusCode: got %d, want %d", got, want)
	}

	p := NewProxy()
	defer p.Close()

	p.SetRequestModifier(nil)
	p.SetResponseModifier(nil)

	tr := martiantest.NewTransport()
	p.SetRoundTripper(tr)
	p.SetTimeout(20 * time.Second)

	tm := martiantest.NewModifier()

	tm.RequestFunc(func(req *http.Request) {
		ctx := NewContext(req)
		ctx.SkipRoundTrip()
	})

	tm.ResponseFunc(func(res *http.Response) {
		cr := res.Request.Header.Get("ContentRange")
		res.StatusCode = http.StatusOK
		res.Body = ioutil.NopCloser(bytes.NewBufferString(testString))
		if cr != "" {
			res.StatusCode = http.StatusPartialContent
			res.Header.Set("Content-Range", cr)
		}
	})

	p.SetRequestModifier(tm)
	p.SetResponseModifier(tm)

	go p.Serve(tsl)

	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		defer conn.Close()
		if err != nil {
			t.Fatalf("net.Dial(): got %v, want no error", err)
		}
		req, err := http.NewRequest("GET", "http://example/example", nil)
		req.Header.Set("ContentRange", "bytes 250-1000/1000")
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		if err := req.WriteProxy(conn); err != nil {
			t.Fatalf("req.WriteProxy(): got %v, want no error", err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			t.Fatalf("http.ReadResponse(): got %v, want no error", err)
		}

		body, _ := ioutil.ReadAll(res.Body)
		bodystr := string(body)
		c1 <- bodystr
	}()

	go func() {
		conn, err := net.Dial("tcp", l.Addr().String())
		defer conn.Close()
		if err != nil {
			t.Fatalf("net.Dial(): got %v, want no error", err)
		}
		req, err := http.NewRequest("GET", "http://example/example", nil)
		if err != nil {
			t.Fatalf("http.NewRequest(): got %v, want no error", err)
		}

		if err := req.WriteProxy(conn); err != nil {
			t.Fatalf("req.WriteProxy(): got %v, want no error", err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			t.Fatalf("http.ReadResponse(): got %v, want no error", err)
		}

		body, _ := ioutil.ReadAll(res.Body)
		bodystr := string(body)
		c2 <- bodystr
	}()

	bodystr1 := <-c1
	bodystr2 := <-c2

	if want1 := strings.Repeat("0", 200); bodystr1 != want1 {
		t.Errorf("res.Body: got %s, want %s", bodystr1, want1)
	}
	if want2 := strings.Repeat("0", 500); bodystr2 != want2 {
		t.Errorf("res.Body: got %s, want %s", bodystr2, want2)
	}
}
