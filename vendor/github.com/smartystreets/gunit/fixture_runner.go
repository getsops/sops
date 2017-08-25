package gunit

import (
	"reflect"
	"testing"
)

func newFixtureRunner(fixture interface{}, t *testing.T, parallel bool) *fixtureRunner {
	return &fixtureRunner{
		parallel:    parallel,
		setup:       -1,
		teardown:    -1,
		outerT:      t,
		fixtureType: reflect.ValueOf(fixture).Type(),
	}
}

type fixtureRunner struct {
	outerT      *testing.T
	fixtureType reflect.Type

	parallel bool
	setup    int
	teardown int
	tests    []*testCase
}

func (this *fixtureRunner) ScanFixtureForTestCases() {
	for methodIndex := 0; methodIndex < this.fixtureType.NumMethod(); methodIndex++ {
		methodName := this.fixtureType.Method(methodIndex).Name
		this.scanFixtureMethod(methodIndex, this.newFixtureMethodInfo(methodName))
	}
}

func (this *fixtureRunner) scanFixtureMethod(methodIndex int, method fixtureMethodInfo) {
	switch {
	case method.isSetup:
		this.setup = methodIndex
	case method.isTeardown:
		this.teardown = methodIndex
	case method.isTest:
		this.tests = append(this.tests, newTestCase(methodIndex, method, this.parallel))
	}
}

func (this *fixtureRunner) RunTestCases() {
	if len(this.tests) == 0 {
		this.outerT.Skipf("Fixture (%v) has no test cases.", this.fixtureType)
		return
	}
	for _, test := range this.tests {
		test.Prepare(this.setup, this.teardown, this.fixtureType)
		test.Run(this.outerT)
	}
}
