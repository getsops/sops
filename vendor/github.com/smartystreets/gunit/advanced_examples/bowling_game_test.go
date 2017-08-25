package examples

import (
	"testing"

	"github.com/smartystreets/gunit"
)

func TestBowlingGameScoringFixture(t *testing.T) {
	gunit.Run(new(BowlingGameScoringFixture), t)
}

type BowlingGameScoringFixture struct {
	*gunit.Fixture
	game *Game
}

func (this *BowlingGameScoringFixture) Setup() {
	this.game = NewGame()
}

func (this *BowlingGameScoringFixture) TestAfterAllGutterBallsTheScoreShouldBeZero() {
	this.rollMany(20, 0)
	this.assertScore(0)
}

func (this *BowlingGameScoringFixture) TestAfterAllOnesTheScoreShouldBeTwenty() {
	this.rollMany(20, 1)
	this.assertScore(20)
}

func (this *BowlingGameScoringFixture) TestSpareReceivesSingleRollBonus() {
	this.rollSpare()
	this.game.Roll(4)
	this.game.Roll(3)
	this.rollMany(16, 0)
	this.assertScore(21)
}

func (this *BowlingGameScoringFixture) TestStrikeReceivesDoubleRollBonus() {
	this.rollStrike()
	this.game.Roll(4)
	this.game.Roll(3)
	this.rollMany(16, 0)
	this.assertScore(24)
}

func (this *BowlingGameScoringFixture) TestPerfectGame() {
	this.rollMany(12, 10)
	this.assertScore(300)
}

func (this *BowlingGameScoringFixture) assertScore(expected int) {
	this.AssertEqual(expected, this.game.Score())
}
func (this *BowlingGameScoringFixture) rollMany(times, pins int) {
	for x := 0; x < times; x++ {
		this.game.Roll(pins)
	}
}
func (this *BowlingGameScoringFixture) rollSpare() {
	this.rollMany(2, 5)
}
func (this *BowlingGameScoringFixture) rollStrike() {
	this.game.Roll(10)
}
