// Licensed under the MIT license, see LICENCE file for details.

/*
Package qtsuite allows quicktest to run test suites.

A test suite is a value with one or more test methods.
For example, the following code defines a suite of test functions that starts
an HTTP server before running each test, and tears it down afterwards:

	type suite struct {
		url string
	}

	func (s *suite) Init(c *qt.C) {
		hnd := func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(w, "%s %s", req.Method, req.URL.Path)
		}
		srv := httptest.NewServer(http.HandlerFunc(hnd))
		c.Defer(srv.Close)
		s.url = srv.URL
	}

	func (s *suite) TestGet(c *qt.C) {
		c.Parallel()
		resp, err := http.Get(s.url)
		c.Assert(err, qt.Equals, nil)
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		c.Assert(err, qt.Equals, nil)
		c.Assert(string(b), qt.Equals, "GET /")
	}

	func (s *suite) TestHead(c *qt.C) {
		c.Parallel()
		resp, err := http.Head(s.url + "/path")
		c.Assert(err, qt.Equals, nil)
		defer resp.Body.Close()
		b, err := ioutil.ReadAll(resp.Body)
		c.Assert(err, qt.Equals, nil)
		c.Assert(string(b), qt.Equals, "")
		c.Assert(resp.ContentLength, qt.Equals, int64(10))
	}

The above code could be invoked from a test function like this:

	func TestHTTPMethods(t *testing.T) {
		qtsuite.Run(qt.New(t), &suite{"http://example.com"})
	}
*/
package qtsuite

import (
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	qt "github.com/frankban/quicktest"
)

// Run runs each test method defined on the given value as a separate
// subtest. A test is a method of the form
//	func (T) TestXxx(*quicktest.C)
// where Xxx does not start with a lowercase letter.
//
// If suite is a pointer, the value pointed to is copied before any
// methods are invoked on it; a new copy is made for each test. This
// means that it is OK for tests to modify fields in suite concurrently
// if desired - it's OK to call c.Parallel().
//
// If suite has a method of the form
//	func (T) Init(*quicktest.C)
// this method will be invoked before each test run.
func Run(c *qt.C, suite interface{}) {
	sv := reflect.ValueOf(suite)
	st := sv.Type()
	init, hasInit := st.MethodByName("Init")
	if hasInit && !isValidMethod(init) {
		c.Fatal("wrong signature for Init, must be Init(*quicktest.C)")
	}
	for i := 0; i < st.NumMethod(); i++ {
		m := st.Method(i)
		if !isTestMethod(m) {
			continue
		}
		c.Run(m.Name, func(c *qt.C) {
			if !isValidMethod(m) {
				c.Fatalf("wrong signature for %s, must be %s(*quicktest.C)", m.Name, m.Name)
			}

			sv := sv
			if st.Kind() == reflect.Ptr {
				sv1 := reflect.New(st.Elem())
				sv1.Elem().Set(sv.Elem())
				sv = sv1
			}
			args := []reflect.Value{sv, reflect.ValueOf(c)}
			if hasInit {
				init.Func.Call(args)
			}
			m.Func.Call(args)
		})
	}
}

var cType = reflect.TypeOf(&qt.C{})

func isTestMethod(m reflect.Method) bool {
	if !strings.HasPrefix(m.Name, "Test") {
		return false
	}
	r, n := utf8.DecodeRuneInString(m.Name[4:])
	return n == 0 || !unicode.IsLower(r)
}

func isValidMethod(m reflect.Method) bool {
	return m.Type.NumIn() == 2 && m.Type.NumOut() == 0 && m.Type.In(1) == cType
}
