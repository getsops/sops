// package unit implements a light-weight x-Unit style testing framework.
// It is basically a scaled-down version of github.com/smartystreets/gunit.
// See https://smartystreets.com/blog/2018/07/lets-build-xunit-in-go for
// an explanation of the basic moving parts.
package unit

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func Run(fixture interface{}, t *testing.T) {
	fixtureType := reflect.TypeOf(fixture)

	for x := 0; x < fixtureType.NumMethod(); x++ {
		testMethodName := fixtureType.Method(x).Name
		if strings.HasPrefix(testMethodName, "Test") {
			t.Run(testMethodName, func(t *testing.T) {
				instance := reflect.New(fixtureType.Elem())

				innerFixture := newFixture(t, testing.Verbose())
				field := instance.Elem().FieldByName("Fixture")
				field.Set(reflect.ValueOf(innerFixture))

				defer innerFixture.Finalize()

				if setup := instance.MethodByName("Setup"); setup.IsValid() {
					setup.Call(nil)
				}

				instance.MethodByName(testMethodName).Call(nil)

				if teardown := instance.MethodByName("Teardown"); teardown.IsValid() {
					teardown.Call(nil)
				}
			})
		}
	}
}

type Fixture struct {
	t       *testing.T
	log     *bytes.Buffer
	verbose bool
}

func newFixture(t *testing.T, verbose bool) *Fixture {
	return &Fixture{t: t, verbose: verbose, log: &bytes.Buffer{}}
}

func (this *Fixture) So(actual interface{}, assert assertion, expected ...interface{}) bool {
	failure := assert(actual, expected...)
	failed := len(failure) > 0
	if failed {
		this.fail(failure)
	}
	return !failed
}

func (this *Fixture) fail(failure string) {
	this.t.Fail()
	this.Print(failure)
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

const comparisonFormat = "Expected: [%s]\nActual:   [%s]"

func (this *Fixture) Error(args ...interface{})            { this.fail(fmt.Sprint(args...)) }
func (this *Fixture) Errorf(f string, args ...interface{}) { this.fail(fmt.Sprintf(f, args...)) }

func (this *Fixture) Print(a ...interface{})                 { fmt.Fprint(this.log, a...) }
func (this *Fixture) Printf(format string, a ...interface{}) { fmt.Fprintf(this.log, format, a...) }
func (this *Fixture) Println(a ...interface{})               { fmt.Fprintln(this.log, a...) }

func (this *Fixture) Write(p []byte) (int, error) { return this.log.Write(p) }
func (this *Fixture) Failed() bool                { return this.t.Failed() }
func (this *Fixture) Name() string                { return this.t.Name() }

func (this *Fixture) Finalize() {
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
	this.t.Fail()
}

// assertion is a copy of github.com/smartystreets/assertions.assertion.
type assertion func(actual interface{}, expected ...interface{}) string
