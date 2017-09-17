package gunit

import (
	"reflect"
	"testing"
)

type testCase struct {
	methodIndex int
	description string
	skipped     bool
	long        bool
	parallel    bool

	setup            int
	teardown         int
	innerFixture     *Fixture
	outerFixtureType reflect.Type
	outerFixture     reflect.Value
}

func newTestCase(methodIndex int, method fixtureMethodInfo, parallel bool) *testCase {
	return &testCase{
		parallel:    parallel,
		methodIndex: methodIndex,
		description: method.name,
		skipped:     method.isSkippedTest,
		long:        method.isLongTest,
	}
}

func (this *testCase) Prepare(setup, teardown int, outerFixtureType reflect.Type) {
	this.setup = setup
	this.teardown = teardown
	this.outerFixtureType = outerFixtureType
}

func (this *testCase) Run(t *testing.T) {
	if this.skipped {
		t.Run(this.description, skip)
	} else if this.long && testing.Short() {
		t.Run(this.description, skipLong)
	} else {
		t.Run(this.description, this.run)
	}
}

func skip(innerT *testing.T)     { innerT.Skip("Skipped test") }
func skipLong(innerT *testing.T) { innerT.Skip("Skipped long-running test") }

func (this *testCase) run(innerT *testing.T) {
	if this.parallel {
		innerT.Parallel()
	}
	this.initializeFixture(innerT)
	defer this.innerFixture.finalize()
	this.runWithSetupAndTeardown()
}
func (this *testCase) initializeFixture(innerT *testing.T) {
	this.innerFixture = newFixture(innerT, testing.Verbose())
	this.outerFixture = reflect.New(this.outerFixtureType.Elem())
	this.outerFixture.Elem().FieldByName("Fixture").Set(reflect.ValueOf(this.innerFixture))
}

func (this *testCase) runWithSetupAndTeardown() {
	this.runSetup()
	defer this.runTeardown()
	this.runTest()
}

func (this *testCase) runSetup() {
	if this.setup < 0 {
		return
	}

	this.outerFixture.Method(this.setup).Call(nil)
}

func (this *testCase) runTest() {
	this.outerFixture.Method(this.methodIndex).Call(nil)
}

func (this *testCase) runTeardown() {
	if this.teardown < 0 {
		return
	}

	this.outerFixture.Method(this.teardown).Call(nil)
}
