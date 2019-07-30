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

package trafficshape

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func compareActions(testSlice []Action, refSlice []Action) (bool, string) {
	if len(testSlice) != len(refSlice) {
		return false, fmt.Sprintf("length: got %d, want %d", len(testSlice), len(refSlice))
	}
	for i, action := range refSlice {
		failure := false
		switch refAction := action.(type) {
		case *Halt:
			if testAction, ok := testSlice[i].(*Halt); ok {
				if *testAction != *refAction {
					failure = true
				}
			} else {
				failure = true
			}
		case *CloseConnection:
			if testAction, ok := testSlice[i].(*CloseConnection); ok {
				if *testAction != *refAction {
					failure = true
				}
			} else {
				failure = true
			}
		case *ChangeBandwidth:
			if testAction, ok := testSlice[i].(*ChangeBandwidth); ok {
				if *testAction != *refAction {
					failure = true
				}
			} else {
				failure = true
			}
		}
		if failure {
			return false, fmt.Sprintf("Action %d: got %+v, want %+v", i, testSlice[i], action)
		}
	}
	return true, ""
}

func TestHandlerIncorrectInputs(t *testing.T) {
	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tt := []struct {
		testcase string
		body     string
	}{
		{
			testcase: `overlapping throttle`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-1000","bandwidth":100},{"bytes":"700-2000","bandwidth":100}],"close_connections":[{"byte":1078,"count":1}]}]}}`,
		},
		{
			testcase: `negative bandwidth`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-","bandwidth":abc}],"close_connections":[{"byte":1078,"count":1}]}]}}`,
		},
		{
			testcase: `negative close byte`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-1000","bandwidth":100}],"close_connections":[{"byte":-1,"count":1}]}]}}`,
		},
		{
			testcase: `uncompiling regex`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example(","throttles":[{"bytes":"500-1000","bandwidth":100}],"close_connections":[{"byte":100,"count":1}]}]}}`,
		},
		{
			testcase: `missing count`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-1000","bandwidth":100}],"halts":[{"byte":100}]}]}}`,
		},
		{
			testcase: `illformed byte range`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500--1000","bandwidth":100}],"close_connections":[{"byte":10,"count":1}]}]}}`,
		},
		{
			testcase: `throttle end < start`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-255","bandwidth":100}],"close_connections":[{"byte":100,"count":1}]}]}}`,
		},
		{
			testcase: `missing comma`,
			body:     `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-1000","bandwidth":100}]"close_connections":[{"byte":100,"count":1}]}]}}`,
		},
		{
			testcase: `missing regex`,
			body:     `{"trafficshape":{"shapes":[{"throttles":[{"bytes":"500-1000","bandwidth":100}],"close_connections":[{"byte":-1,"count":1}]}]}}`,
		},
		{
			testcase: `negative default bandwidth`,
			body:     `{"trafficshape":{"default":{"bandwidth":{"up":-100000,"down":100000},"latency":1000},"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-1000","bandwidth":100}]",close_connections":[{"byte":100,"count":1}]}]}}`,
		},
		{
			testcase: `negative default latency`,
			body:     `{"trafficshape":{"default":{"bandwidth":{"up":100000,"down":100000},"latency":-1000},"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"500-1000","bandwidth":100}]",close_connections":[{"byte":100,"count":1}]}]}}`,
		},
	}

	for i, tc := range tt {
		t.Logf("case %d: %s", i+1, tc.testcase)

		tsl := NewListener(l)
		defer tsl.Close()

		h := NewHandler(tsl)

		req, err := http.NewRequest("POST", "test", bytes.NewBufferString(tc.body))
		if err != nil {
			t.Fatalf("%d. http.NewRequest(): got %v, want no error", i, err)
		}
		rw := httptest.NewRecorder()

		h.ServeHTTP(rw, req)

		if got := rw.Code; got != 400 {
			t.Errorf("%d. rw.Code: got %d, want %d", i+1, got, 400)
		}
	}
}

func TestHandlerClear(t *testing.T) {

	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tsl := NewListener(l)
	defer tsl.Close()

	h := NewHandler(tsl)
	startTime := time.Now()
	jsonString := `{"trafficshape":{}}`
	req, err := http.NewRequest("POST", "test", bytes.NewBufferString(jsonString))
	if err != nil {
		t.Fatalf("http.NewRequest(): got %v, want no error", err)
	}

	rw := httptest.NewRecorder()

	h.ServeHTTP(rw, req)

	if got, want := rw.Code, 200; got != want {
		t.Errorf(" rw.Code: got %d, want %d", got, want)
	}

	defaults := tsl.Defaults()

	if got, want := defaults.Bandwidth.Down, DefaultBitrate/8; got != want {
		t.Errorf("default downstream bandwidth: got %d, want %d", got, want)
	}
	if got, want := defaults.Latency, int64(0); got != want {
		t.Errorf("default latency: got %d, want %d", got, want)
	}

	if got, want := tsl.WriteBucket.Capacity(), DefaultBitrate/8; got != want {
		t.Errorf("tsl WriteBucket Capacity: got %d, want %d", got, want)
	}

	tsl.Shapes.RLock()
	if got, want := len(tsl.Shapes.M), 0; got != want {
		t.Errorf("length of shape map: got %d, want %d", got, want)
	}
	if modifiedTime := tsl.Shapes.LastModifiedTime; modifiedTime.Before(startTime) {
		t.Errorf("modified time is before start time; should be after")
	}
	tsl.Shapes.RUnlock()
}

func TestHandlerActions(t *testing.T) {
	l, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		t.Fatalf("net.Listen(): got %v, want no error", err)
	}

	tt := []struct {
		jsonString string
		actions    []Action
	}{
		{
			jsonString: `{"trafficshape":{"shapes":[{"url_regex":"http://example/example", "max_global_bandwidth":1000, "throttles":[{"bytes":"500-1000","bandwidth":100},{"bytes":"1000-2000","bandwidth":300},{"bytes":"2001-","bandwidth":400}],
	"halts":[{"byte":530,"duration": 5, "count": 1}],"close_connections":[{"byte":1078,"count":1}]}]}}`,
			actions: []Action{
				&ChangeBandwidth{Byte: 500, Bandwidth: 100},
				&Halt{Byte: 530, Duration: 5, Count: 1},
				&ChangeBandwidth{Byte: 1000, Bandwidth: 300},
				&CloseConnection{Byte: 1078, Count: 1},
				&ChangeBandwidth{Byte: 2000, Bandwidth: 1000},
				&ChangeBandwidth{Byte: 2001, Bandwidth: 400},
			},
		},
		{
			jsonString: `{"trafficshape":{"shapes":[{"url_regex":"http://example/example","throttles":[{"bytes":"-","bandwidth":100}],
			"close_connections":[{"byte":100,"count":1}]}]}}`,
			actions: []Action{
				&ChangeBandwidth{Byte: 0, Bandwidth: 100},
				&CloseConnection{Byte: 100, Count: 1},
			},
		},
	}

	for i, tc := range tt {
		tsl := NewListener(l)
		defer tsl.Close()

		h := NewHandler(tsl)
		startTime := time.Now()
		req, err := http.NewRequest("POST", "test", bytes.NewBufferString(tc.jsonString))
		if err != nil {
			t.Fatalf("%d. http.NewRequest(): got %v, want no error", i, err)
		}

		rw := httptest.NewRecorder()

		h.ServeHTTP(rw, req)

		if got, want := rw.Code, 200; got != want {
			t.Errorf("%d. rw.Code: got %d, want %d", i+1, got, want)
		}

		tsl.Shapes.RLock()
		defer tsl.Shapes.RUnlock()
		if got, want := len(tsl.Shapes.M), 1; got != want {
			t.Errorf("tc.%d length of shape map: got %d, want %d", i+1, got, want)
		}
		tsl.Shapes.M["http://example/example"].RLock()
		defer tsl.Shapes.M["http://example/example"].RUnlock()
		if same, errStr := compareActions(tsl.Shapes.M["http://example/example"].Shape.Actions, tc.actions); !same {
			t.Errorf(errStr)
		}
		if modifiedTime := tsl.Shapes.LastModifiedTime; modifiedTime.Before(startTime) {
			t.Errorf("tc.%d modified time is before start time; should be after", i+1)
		}
	}
}
