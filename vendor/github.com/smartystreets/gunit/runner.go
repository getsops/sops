package gunit

import (
	"reflect"
	"testing"
)

// Run receives an instance of a struct that embeds *Fixture.
// The struct definition may include Setup*, Teardown*, and Test*
// methods which will be run as an xUnit-style test fixture.
func Run(fixture interface{}, t *testing.T) {
	t.Parallel()
	run(fixture, t, true)
}

// RunSequential, like Run receives an instance of a struct that embeds *Fixture.
// The fixture is run in much the same way, except that it will not be run in
// parallel with other fixtures in the same package.
func RunSequential(fixture interface{}, t *testing.T) {
	run(fixture, t, false)
}

func run(fixture interface{}, t *testing.T, parallel bool) {
	ensureEmbeddedFixture(fixture, t)
	runner := newFixtureRunner(fixture, t, parallel)
	runner.ScanFixtureForTestCases()
	runner.RunTestCases()
}

func ensureEmbeddedFixture(fixture interface{}, t testingT) {
	fixtureType := reflect.TypeOf(fixture)
	embedded, _ := fixtureType.Elem().FieldByName("Fixture")
	if embedded.Type != embeddedGoodExample.Type {
		t.Fatalf("Type (%v) lacks embedded *gunit.Fixture.", fixtureType)
	}
}

type goodExample struct{ *Fixture }

var embeddedGoodExample, _ = reflect.TypeOf(new(goodExample)).Elem().FieldByName("Fixture")
