package assert

import (
	"testing"

	"github.com/smartystreets/assertions/internal/unit"
	"github.com/smartystreets/assertions/should"
)

func TestPassedResultFixture(t *testing.T) {
	unit.Run(new(PassedResultFixture), t)
}

type PassedResultFixture struct {
	*unit.Fixture

	result *Result
}

func (this *PassedResultFixture) Setup() {
	this.result = So(1, should.Equal, 1)
	this.result.logger = capture()
	this.result.stdout = this.result.logger.Log
}

func (this *PassedResultFixture) TestQueryFunctions() {
	this.So(this.result.Error(), should.BeNil)
	this.So(this.result.Failed(), should.BeFalse)
	this.So(this.result.Passed(), should.BeTrue)
	this.So(this.result.String(), should.Equal, "✔ So(actual: 1, should.Equal, expected: [1])")
}
func (this *PassedResultFixture) TestPrintln() {
	this.So(this.result.Println(), should.Equal, this.result)
	this.So(this.result.logger.Log.String(), should.BeBlank)
}
func (this *PassedResultFixture) TestLog() {
	this.So(this.result.Log(), should.Equal, this.result)
	this.So(this.result.logger.Log.String(), should.BeBlank)
}
func (this *PassedResultFixture) TestPanic() {
	this.So(this.result.Panic(), should.Equal, this.result)
	this.So(this.result.logger.Log.String(), should.BeBlank)
}
func (this *PassedResultFixture) TestFatal() {
	this.So(this.result.Fatal(), should.Equal, this.result)
	this.So(this.result.logger.Log.String(), should.BeBlank)
}
