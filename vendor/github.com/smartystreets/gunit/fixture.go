package gunit

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// Fixture keeps track of test status (failed, passed, skipped) and
// handles custom logging for xUnit style tests as an embedded field.
// The Fixture manages an instance of *testing.T. Certain methods
// defined herein merely forward to calls on the *testing.T:
//
//     - Fixture.Error(...) ----> *testing.T.Error
//     - Fixture.Errorf(...) ---> *testing.T.Errorf
//     - Fixture.Print(...) ----> *testing.T.Log or fmt.Print
//     - Fixture.Printf(...) ---> *testing.T.Logf or fmt.Printf
//     - Fixture.Println(...) --> *testing.T.Log or fmt.Println
//     - Fixture.Failed() ------> *testing.T.Failed()
//     - Fixture.fail() --------> *testing.T.Fail()
//
// We don't use these methods much, preferring instead to lean heavily
// on Fixture.So and the rich set of should-style assertions provided at
// github.com/smartystreets/assertions/should
type Fixture struct {
	t       testingT
	log     *bytes.Buffer
	verbose bool
}

func newFixture(t testingT, verbose bool) *Fixture {
	return &Fixture{t: t, verbose: verbose, log: &bytes.Buffer{}}
}

// So is a convenience method for reporting assertion failure messages,
// from the many assertion functions found in github.com/smartystreets/assertions/should.
// Example: this.So(actual, should.Equal, expected)
func (this *Fixture) So(actual interface{}, assert assertion, expected ...interface{}) bool {

	failure := assert(actual, expected...)
	failed := len(failure) > 0
	if failed {
		this.t.Fail()
		this.fail(failure)
	}
	return !failed
}

// Assert tests a boolean which, if not true, marks the current test case as failed and
// prints the provided message.
func (this *Fixture) Assert(condition bool, messages ...string) bool {
	if !condition {
		if len(messages) == 0 {
			messages = append(messages, "Expected condition to be true, was false instead.")
		}
		this.fail(strings.Join(messages, ", "))
	}
	return condition
}
func (this *Fixture) AssertEqual(expected, actual interface{}) bool {
	return this.Assert(expected == actual, fmt.Sprintf(comparisonFormat, fmt.Sprint(expected), fmt.Sprint(actual)))
}
func (this *Fixture) AssertSprintEqual(expected, actual interface{}) bool {
	return this.AssertEqual(fmt.Sprint(expected), fmt.Sprint(actual))
}
func (this *Fixture) AssertSprintfEqual(expected, actual interface{}, format string) bool {
	return this.AssertEqual(fmt.Sprintf(format, expected), fmt.Sprintf(format, actual))
}
func (this *Fixture) AssertDeepEqual(expected, actual interface{}) bool {
	return this.Assert(reflect.DeepEqual(expected, actual),
		fmt.Sprintf(comparisonFormat, fmt.Sprintf("%#v", expected), fmt.Sprintf("%#v", actual)))
}

func (this *Fixture) Error(args ...interface{})            { this.fail(fmt.Sprint(args...)) }
func (this *Fixture) Errorf(f string, args ...interface{}) { this.fail(fmt.Sprintf(f, args...)) }

func (this *Fixture) Print(a ...interface{})                 { fmt.Fprint(this.log, a...) }
func (this *Fixture) Printf(format string, a ...interface{}) { fmt.Fprintf(this.log, format, a...) }
func (this *Fixture) Println(a ...interface{})               { fmt.Fprintln(this.log, a...) }

func (this *Fixture) Failed() bool { return this.t.Failed() }

func (this *Fixture) fail(failure string) {
	this.t.Fail()
	this.Print(newFailureReport(failure))
}

func (this *Fixture) finalize() {
	if r := recover(); r != nil {
		this.recoverPanic(r)
	}

	if this.t.Failed() || (this.verbose && this.log.Len() > 0) {
		this.t.Log("\n" + strings.TrimSpace(this.log.String()) + "\n")
	}
}
func (this *Fixture) recoverPanic(r interface{}) {
	this.Println("PANIC:", r)
	buffer := make([]byte, 1024*16)
	runtime.Stack(buffer, false)
	this.Println(strings.TrimSpace(string(buffer)))
	this.Println("* (Additional tests may have been skipped as a result of the panic shown above.)")
	this.t.Fail()
}

const comparisonFormat = "Expected: [%s]\nActual:   [%s]"

// assertion is a copy of github.com/smartystreets/assertions.assertion.
type assertion func(actual interface{}, expected ...interface{}) string
