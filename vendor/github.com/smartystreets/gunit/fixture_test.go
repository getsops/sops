package gunit

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/smartystreets/assertions/should"
)

func TestFinalizeAfterNoActions(t *testing.T) {
	test := Setup(false)

	test.fixture.finalize()

	if test.fakeT.failed {
		t.Error("Fake should not have been marked as failed.")
	}
	if test.out.Len() > 0 {
		t.Errorf("Output was not blank: '%s'", test.out.String())
	}
}

func TestFinalizeAfterFailure(t *testing.T) {
	test := Setup(false)

	test.fakeT.Fail()

	test.fixture.finalize()

	if output := strings.TrimSpace(test.out.String()); strings.Contains(output, "Failure") {
		t.Errorf("Unexpected output: '%s'", output)
	}
}

func TestSoPasses(t *testing.T) {
	test := Setup(false)

	result := test.fixture.So(true, should.BeTrue)
	test.fixture.finalize()

	if !result {
		t.Error("Expected true result, got false")
	}
	if test.out.Len() > 0 {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if test.fakeT.failed {
		t.Error("Test was erroneously marked as failed.")
	}
}

func TestSoFailsAndLogs(t *testing.T) {
	test := Setup(false)

	result := test.fixture.So(true, should.BeFalse)
	test.fixture.finalize()

	if result {
		t.Error("Expected false result, got true")
	}
	if output := test.out.String(); !strings.Contains(output, "Expected:") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
}

func TestAssertPasses(t *testing.T) {
	test := Setup(false)

	test.fixture.Assert(true)
	test.fixture.finalize()

	if test.out.Len() > 0 {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if test.fakeT.failed {
		t.Error("Test was erroneously marked as failed.")
	}
}

func TestAssertFailsAndLogs(t *testing.T) {
	test := Setup(false)

	returned := test.fixture.Assert(false)
	test.fixture.finalize()

	if output := test.out.String(); !strings.Contains(output, "Expected condition to be true, was false instead.") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
	if returned != false {
		t.Error("The same condition should be returned form Assert.")
	}
}

func TestAssertWithCustomMessageFailsAndLogs(t *testing.T) {
	test := Setup(false)

	test.fixture.Assert(false, "gophers!")
	test.fixture.finalize()

	if output := test.out.String(); !strings.Contains(output, "gophers!") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
}

func TestAssertEqualPasses(t *testing.T) {
	test := Setup(false)

	test.fixture.AssertEqual(1, 1)
	test.fixture.finalize()

	if test.out.Len() > 0 {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if test.fakeT.failed {
		t.Error("Test was erroneously marked as failed.")
	}
}

func TestAssertEqualFails(t *testing.T) {
	test := Setup(false)

	returned := test.fixture.AssertEqual(1, 2)
	test.fixture.finalize()

	if output := test.out.String(); !strings.Contains(output, "Expected: [1]\nActual:   [2]") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
	if returned != false {
		t.Error("Should have returned the result of the assertion (false in this case).")
	}
}

func TestAssertSprintEqualPasses(t *testing.T) {
	test := Setup(false)

	returned := test.fixture.AssertSprintEqual(1, 1.0)
	test.fixture.finalize()

	if test.out.Len() > 0 {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if test.fakeT.failed {
		t.Error("Test was erroneously marked as failed.")
	}
	if returned != true {
		t.Error("Should have returned the result of the assertion (true in the case).")
	}
}

func TestAssertSprintEqualFails(t *testing.T) {
	test := Setup(false)

	test.fixture.AssertSprintEqual(1, 2)
	test.fixture.finalize()

	if output := test.out.String(); !strings.Contains(output, "Expected: [1]\nActual:   [2]") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
}

func TestAssertSprintfEqualPasses(t *testing.T) {
	test := Setup(false)

	test.fixture.AssertSprintfEqual(1, uint(1), "%d")
	test.fixture.finalize()

	if test.out.Len() > 0 {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if test.fakeT.failed {
		t.Error("Test was erroneously marked as failed.")
	}
}

func TestAssertSprintfEqualFails(t *testing.T) {
	test := Setup(false)

	test.fixture.AssertSprintfEqual(1, 2, "%d")
	test.fixture.finalize()

	if output := test.out.String(); !strings.Contains(output, "Expected: [1]\nActual:   [2]") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
}

func TestAssertDeepEqualPasses(t *testing.T) {
	test := Setup(false)

	test.fixture.AssertDeepEqual(1, 1)
	test.fixture.finalize()

	if test.out.Len() > 0 {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if test.fakeT.failed {
		t.Error("Test was erroneously marked as failed.")
	}
}

func TestAssertDeepEqualFails(t *testing.T) {
	test := Setup(false)

	test.fixture.AssertDeepEqual(1, 2)
	test.fixture.finalize()

	if output := test.out.String(); !strings.Contains(output, "Expected: [1]\nActual:   [2]") {
		t.Errorf("Unexpected ouput: '%s'", test.out.String())
	}
	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
}

func TestErrorFailsAndLogs(t *testing.T) {
	test := Setup(false)

	test.fixture.Error("1", "2", "3")
	test.fixture.finalize()

	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
	if output := test.out.String(); !strings.Contains(output, "123") {
		t.Errorf("Expected string containing: '123' Got: '%s'", output)
	}
}

func TestErrorfFailsAndLogs(t *testing.T) {
	test := Setup(false)

	test.fixture.Errorf("%s%s%s", "1", "2", "3")
	test.fixture.finalize()

	if !test.fakeT.failed {
		t.Error("Test should have been marked as failed.")
	}
	if output := test.out.String(); !strings.Contains(output, "123") {
		t.Errorf("Expected string containing: '123' Got: '%s'", output)
	}
}

func TestFixturePrinting(t *testing.T) {
	test := Setup(true)

	test.fixture.Print("Print")
	test.fixture.Println("Println")
	test.fixture.Printf("Printf")
	test.fixture.finalize()

	output := test.out.String()
	if !strings.Contains(output, "Print") {
		t.Error("Expected to see 'Print' in the output.")
	}
	if !strings.Contains(output, "Println") {
		t.Error("Expected to see 'Println' in the output.")
	}
	if !strings.Contains(output, "Printf") {
		t.Error("Expected to see 'Printf' in the output.")
	}
	if t.Failed() {
		t.Logf("Actual output: \n%s\n", output)
	}
}

func TestPanicIsRecoveredAndPrintedByFinalize(t *testing.T) {
	test := Setup(false)

	var freakOut = func() {
		defer test.fixture.finalize()
		panic("GOPHERS!")
	}

	freakOut()

	output := test.out.String()
	if !strings.Contains(output, "PANIC: GOPHERS!") {
		t.Errorf("Expected string containing: 'PANIC: GOPHERS!' Got: '%s'", output)
	}
	if !strings.Contains(output, "github.com/smartystreets/gunit.(*Fixture).finalize") {
		t.Error("Expected string containing stack trace information...")
	}
	if !strings.Contains(output, "* (Additional tests may have been skipped as a result of the panic shown above.)") {
		t.Error("Expected string containing warning about additional tests not being run.")
	}
}

func TestFailed(t *testing.T) {
	test := Setup(false)

	if test.fixture.Failed() {
		t.Error("Expected Failed() to return false, got true instead.")
	}

	test.fixture.Error("HI")

	if !test.fixture.Failed() {
		t.Error("Expected Failed() to return true, got false instead.")
	}
}

//////////////////////////////////////////////////////////////////////////////

type FixtureTestState struct {
	fixture *Fixture
	fakeT   *FakeTestingT
	out     *bytes.Buffer
	verbose bool
}

func Setup(verbose bool) *FixtureTestState {
	this := &FixtureTestState{}
	this.out = &bytes.Buffer{}
	this.fakeT = &FakeTestingT{log: this.out}
	this.fixture = newFixture(this.fakeT, verbose)
	return this
}

//////////////////////////////////////////////////////////////////////////////

type FakeTestingT struct {
	log    *bytes.Buffer
	failed bool
}

func (self *FakeTestingT) Log(args ...interface{}) { fmt.Fprint(self.log, args...) }
func (self *FakeTestingT) Fail()                   { self.failed = true }
func (self *FakeTestingT) Failed() bool            { return self.failed }
func (this *FakeTestingT) Fatalf(format string, args ...interface{}) {
	this.Fail()
	this.Log(fmt.Sprintf(format, args))
}

//////////////////////////////////////////////////////////////////////////////
