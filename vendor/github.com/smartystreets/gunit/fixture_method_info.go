package gunit

import "strings"

type fixtureMethodInfo struct {
	name          string
	isSetup       bool
	isTeardown    bool
	isTest        bool
	isFocusTest   bool
	isLongTest    bool
	isSkippedTest bool
}

func (this *fixtureRunner) newFixtureMethodInfo(name string) fixtureMethodInfo {
	var (
		isSetup           = strings.HasPrefix(name, "Setup")
		isTeardown        = strings.HasPrefix(name, "Teardown")
		isTest            = strings.HasPrefix(name, "Test")
		isLongTest        = strings.HasPrefix(name, "LongTest")
		isFocusTest       = strings.HasPrefix(name, "FocusTest")
		isFocusLongTest   = strings.HasPrefix(name, "FocusLongTest")
		isSkippedTest     = strings.HasPrefix(name, "SkipTest")
		isSkippedLongTest = strings.HasPrefix(name, "SkipLongTest")
	)

	return fixtureMethodInfo{
		name:          name,
		isSetup:       isSetup,
		isTeardown:    isTeardown,
		isLongTest:    isLongTest || isSkippedLongTest || isFocusLongTest,
		isFocusTest:   isFocusTest || isFocusLongTest,
		isSkippedTest: isSkippedTest || isSkippedLongTest,
		isTest:        isTest || isLongTest || isSkippedTest || isSkippedLongTest || isFocusTest || isFocusLongTest,
	}
}
