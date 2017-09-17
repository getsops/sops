package gunit

// testingT represents the functional subset from *testing.T needed by Fixture.
type testingT interface {
	Log(args ...interface{})
	Fail()
	Failed() bool
	Fatalf(format string, args ...interface{})
}
