package scan

import (
	"testing"

	"github.com/smartystreets/assertions"
	"github.com/smartystreets/assertions/should"
)

//////////////////////////////////////////////////////////////////////////////

func TestParseFileWithValidFixturesAndConstructs(t *testing.T) {
	test := &FixtureParsingFixture{t: t, input: comprehensiveTestCode}
	test.ParseFixtures()
	test.AssertFixturesParsedAccuratelyAndCompletely()
}

//////////////////////////////////////////////////////////////////////////////

type FixtureParsingFixture struct {
	t *testing.T

	input      string
	readError  error
	parseError error
	fixtures   []*fixtureInfo
}

func (this *FixtureParsingFixture) ParseFixtures() {
	this.fixtures, this.parseError = scanForFixtures(this.input)
}

func (this *FixtureParsingFixture) AssertFixturesParsedAccuratelyAndCompletely() {
	this.assertFileWasReadWithoutError()
	this.assertFileWasParsedWithoutError()
	this.assertAllFixturesParsed()
	this.assertParsedFixturesAreCorrect()
}
func (this *FixtureParsingFixture) assertFileWasReadWithoutError() {
	if this.readError != nil {
		this.t.Error("Problem: cound't read the input file:", this.readError)
		this.t.FailNow()
	}
}
func (this *FixtureParsingFixture) assertFileWasParsedWithoutError() {
	if this.parseError != nil {
		this.t.Error("Problem: unexpected parsing error: ", this.parseError)
		this.t.FailNow()
	}
}
func (this *FixtureParsingFixture) assertAllFixturesParsed() {
	if len(this.fixtures) != len(expected) {
		this.t.Logf("Problem: Got back the wrong number of fixtures. Expected: %d Got: %d", len(expected), len(this.fixtures))
		this.t.FailNow()
	}
}
func (this *FixtureParsingFixture) assertParsedFixturesAreCorrect() {
	for x := 0; x < len(expected); x++ {
		key := this.fixtures[x].StructName
		if ok, message := assertions.So(this.fixtures[x], should.Resemble, expected[key]); !ok {
			this.t.Errorf("Comparison failure for record: %d\n%s", x, message)
		}
	}
}

func (this *FixtureParsingFixture) AssertErrorWasReturned() {
	if this.parseError == nil {
		this.t.Error("Expected an error, but got nil instead")
	}
}

//////////////////////////////////////////////////////////////////////////////

var expected = map[string]*fixtureInfo{
	"BowlingGameScoringTests": {
		StructName: "BowlingGameScoringTests",
		TestCases: []*testCaseInfo{
			{CharacterPosition: 335, Name: "TestAfterAllGutterBallsTheScoreShouldBeZero",},
			{CharacterPosition: 490, Name: "TestAfterAllOnesTheScoreShouldBeTwenty",},
			{CharacterPosition: 641, Name: "SkipTestASpareDeservesABonus",},
			{CharacterPosition: 718, Name: "LongTestPerfectGame",},
			{CharacterPosition: 852, Name: "SkipLongTestPerfectGame",},
		},
	},
}

//////////////////////////////////////////////////////////////////////////////
