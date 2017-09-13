package gunit

import "strings"

type fixtureMethodInfo struct {
	name          string
	isSetup       bool
	isTeardown    bool
	isTest        bool
	isLongTest    bool
	isSkippedTest bool
}

func (this *fixtureRunner) newFixtureMethodInfo(name string) fixtureMethodInfo {
	isTest := strings.HasPrefix(name, "Test")
	isLongTest := strings.HasPrefix(name, "LongTest")
	isSkippedTest := strings.HasPrefix(name, "SkipTest")
	isSkippedLongTest := strings.HasPrefix(name, "SkipLongTest")

	return fixtureMethodInfo{
		name:          name,
		isSetup:       strings.HasPrefix(name, "Setup"),
		isTeardown:    strings.HasPrefix(name, "Teardown"),
		isLongTest:    isLongTest,
		isSkippedTest: isSkippedTest || isSkippedLongTest,
		isTest:        isTest || isLongTest || isSkippedTest || isSkippedLongTest,
	}
}
