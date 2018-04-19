package gunit

import (
	"reflect"
	"testing"

	"github.com/smartystreets/gunit/scan"
)

func newFixtureRunner(fixture interface{}, t *testing.T, parallel bool, positions scan.TestCasePositions) *fixtureRunner {
	return &fixtureRunner{
		parallel:    parallel,
		setup:       -1,
		teardown:    -1,
		outerT:      t,
		fixtureType: reflect.ValueOf(fixture).Type(),
		positions:   positions,
	}
}

type fixtureRunner struct {
	outerT      *testing.T
	fixtureType reflect.Type

	parallel  bool
	setup     int
	teardown  int
	focus     []*testCase
	tests     []*testCase
	positions scan.TestCasePositions
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
	case method.isFocusTest:
		this.focus = append(this.focus, newTestCase(methodIndex, method, this.parallel, this.positions))
	case method.isTest:
		this.tests = append(this.tests, newTestCase(methodIndex, method, this.parallel, this.positions))
	}
}

func (this *fixtureRunner) RunTestCases() {
	if len(this.focus) > 0 {
		this.tests = append(this.focus, skipped(this.tests)...)
	}
	if len(this.tests) > 0 {
		this.runTestCases(this.tests)
	} else {
		this.outerT.Skipf("Fixture (%v) has no test cases.", this.fixtureType)
	}
}

func (this *fixtureRunner) runTestCases(cases []*testCase) {
	for _, test := range cases {
		test.Prepare(this.setup, this.teardown, this.fixtureType)
		test.Run(this.outerT)
	}
}

func skipped(cases []*testCase) []*testCase {
	for _, test := range cases {
		test.skipped = true
	}
	return cases
}
