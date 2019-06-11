package assert

import (
	"testing"

	"github.com/smartystreets/assertions/internal/unit"
	"github.com/smartystreets/assertions/should"
)

func TestFailedResultFixture(t *testing.T) {
	unit.Run(new(FailedResultFixture), t)
}

type FailedResultFixture struct {
	*unit.Fixture

	result *Result
}

func (this *FailedResultFixture) Setup() {
	this.result = So(1, should.Equal, 2)
	this.result.logger = capture()
	this.result.stdout = this.result.logger.Log
}

func (this *FailedResultFixture) assertLogMessageContents() {
	this.So(this.result.logger.Log.String(), should.ContainSubstring, "✘ So(actual: 1, should.Equal, expected: [2])")
	this.So(this.result.logger.Log.String(), should.ContainSubstring, "Assertion failure at ")
	this.So(this.result.logger.Log.String(), should.EndWith, "Expected: '2'\nActual:   '1'\n(Should be equal)\n")
}

func (this *FailedResultFixture) TestQueryFunctions() {
	this.So(this.result.Failed(), should.BeTrue)
	this.So(this.result.Passed(), should.BeFalse)
	this.So(this.result.logger.Log.Len(), should.Equal, 0)

	this.result.logger.Print(this.result.String())
	this.result.logger.Print(this.result.Error())
	this.assertLogMessageContents()
}

func (this *FailedResultFixture) TestPrintln() {
	this.So(this.result.Println(), should.Equal, this.result)
	this.assertLogMessageContents()
}

func (this *FailedResultFixture) TestLog() {
	this.So(this.result.Log(), should.Equal, this.result)
	this.assertLogMessageContents()
}

func (this *FailedResultFixture) TestPanic() {
	this.So(func() { this.result.Panic() }, should.Panic)
	this.assertLogMessageContents()
}

func (this *FailedResultFixture) TestFatal() {
	this.So(this.result.Fatal(), should.Equal, this.result)
	this.assertLogMessageContents()
}
