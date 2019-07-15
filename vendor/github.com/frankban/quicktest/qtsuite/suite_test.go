// Licensed under the MIT license, see LICENCE file for details.

package qtsuite_test

import (
	"bytes"
	"fmt"
	"testing"

	qt "github.com/frankban/quicktest"
	"github.com/frankban/quicktest/qtsuite"
)

func TestRunSuite(t *testing.T) {
	c := qt.New(t)
	var calls []call
	tt := &testingT{}
	qtsuite.Run(qt.New(tt), testSuite{calls: &calls})
	c.Assert(calls, qt.DeepEquals, []call{
		{"Test1", 0},
		{"Test4", 0},
	})
}

func TestRunSuiteEmbedded(t *testing.T) {
	c := qt.New(t)
	var calls []call
	tt := &testingT{}
	suite := struct {
		testSuite
	}{testSuite: testSuite{calls: &calls}}
	qtsuite.Run(qt.New(tt), suite)
	c.Assert(calls, qt.DeepEquals, []call{
		{"Test1", 0},
		{"Test4", 0},
	})
}

func TestRunSuitePtr(t *testing.T) {
	c := qt.New(t)
	var calls []call
	tt := &testingT{}
	qtsuite.Run(qt.New(tt), &testSuite{calls: &calls})
	c.Assert(calls, qt.DeepEquals, []call{
		{"Init", 0},
		{"Test1", 1},
		{"Init", 0},
		{"Test4", 1},
	})
}

type testSuite struct {
	init  int
	calls *[]call
}

func (s testSuite) addCall(name string) {
	*s.calls = append(*s.calls, call{Name: name, Init: s.init})
}

func (s *testSuite) Init(*qt.C) {
	s.addCall("Init")
	s.init++
}

func (s testSuite) Test1(*qt.C) {
	s.addCall("Test1")
}

func (s testSuite) Test2() {
	s.addCall("Test2")
}

func (s testSuite) Test3(*testing.T) {
	s.addCall("Test3")
}

func (s testSuite) Test4(*qt.C) {
	s.addCall("Test4")
}

func (s testSuite) Test5(*qt.C) bool {
	s.addCall("Test5")
	return false
}

func (s testSuite) Testa(*qt.C) {
	s.addCall("Testa")
}

type call struct {
	Name string
	Init int
}

func TestInvalidInit(t *testing.T) {
	c := qt.New(t)
	tt := &testingT{}
	tc := qt.New(tt)
	qtsuite.Run(tc, invalidTestSuite{})
	c.Assert(tt.fatalString(), qt.Equals, "wrong signature for Init, must be Init(*quicktest.C)")
}

type invalidTestSuite struct{}

func (invalidTestSuite) Init() {}

// testingT can be passed to qt.New for testing purposes.
type testingT struct {
	testing.TB

	errorBuf bytes.Buffer
	fatalBuf bytes.Buffer

	subTestResult bool
	subTestName   string
	subTestT      *testing.T
}

// Error overrides *testing.T.Error so that messages are collected.
func (t *testingT) Error(a ...interface{}) {
	fmt.Fprint(&t.errorBuf, a...)
}

// Fatal overrides *testing.T.Fatal so that messages are collected and the
// goroutine is not killed.
func (t *testingT) Fatal(a ...interface{}) {
	fmt.Fprint(&t.fatalBuf, a...)
}

// Run overrides *testing.T.Run.
func (t *testingT) Run(name string, f func(t *testing.T)) bool {
	t.subTestName, t.subTestT = name, &testing.T{}
	ch := make(chan struct{})
	// Run the subtest in its own goroutine so that if it calls runtime.GoExit,
	// we can still return appropriately.
	go func() {
		defer close(ch)
		f(t.subTestT)
	}()
	<-ch
	return t.subTestResult
}

// errorString returns the error message.
func (t *testingT) errorString() string {
	return t.errorBuf.String()
}

// fatalString returns the fatal error message.
func (t *testingT) fatalString() string {
	return t.fatalBuf.String()
}
