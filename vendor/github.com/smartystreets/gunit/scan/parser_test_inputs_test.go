package scan

const comprehensiveTestCode = `package parse

import (
	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

type BowlingGameScoringTests struct {
	*gunit.Fixture

	game *Game
}

func (self *BowlingGameScoringTests) SetupTheGame() {
	self.game = NewGame()
}

func (self *BowlingGameScoringTests) TeardownTheGame() {
	self.game = nil
}

func (self *BowlingGameScoringTests) TestAfterAllGutterBallsTheScoreShouldBeZero() {
	self.rollMany(20, 0)
	self.So(self.game.Score(), should.Equal, 0)
}

func (self *BowlingGameScoringTests) TestAfterAllOnesTheScoreShouldBeTwenty() {
	self.rollMany(20, 1)
	self.So(self.game.Score(), should.Equal, 20)
}

func (self *BowlingGameScoringTests) SkipTestASpareDeservesABonus()      {}

func (self *BowlingGameScoringTests) LongTestPerfectGame() {
	self.rollMany(12, 10)
	self.So(self.game.Score(), should.Equal, 300)
}

func (self *BowlingGameScoringTests) SkipLongTestPerfectGame() {
	self.rollMany(12, 10)
	self.So(self.game.Score(), should.Equal, 300)
}

func (self *BowlingGameScoringTests) rollMany(times, pins int) {
	for x := 0; x < times; x++ {
		self.game.Roll(pins)
	}
}
func (self *BowlingGameScoringTests) rollSpare() {
	self.game.Roll(5)
	self.game.Roll(5)
}
func (self *BowlingGameScoringTests) rollStrike() {
	self.game.Roll(10)
}

func (self *BowlingGameScoringTests) TestNotNiladic_ShouldNotBeCollected(a int) {
	// This should not be collected (it's not niladic)
}
func (self *BowlingGameScoringTests) TestNotVoid_ShouldNOTBeCollected() int {
	return -1
	// This should not be collected (it's not void)
}

//////////////////////////////////////////////////////////////////////////////

// Game contains the state of a bowling game.
type Game struct {
	rolls   []int
	current int
}

// NewGame allocates and starts a new game of bowling.
func NewGame() *Game {
	game := new(Game)
	game.rolls = make([]int, maxThrowsPerGame)
	return game
}

// Roll rolls the ball and knocks down the number of pins specified by pins.
func (self *Game) Roll(pins int) {
	self.rolls[self.current] = pins
	self.current++
}

// Score calculates and returns the player's current score.
func (self *Game) Score() (sum int) {
	for throw, frame := 0, 0; frame < framesPerGame; frame++ {
		if self.isStrike(throw) {
			sum += self.strikeBonusFor(throw)
			throw += 1
		} else if self.isSpare(throw) {
			sum += self.spareBonusFor(throw)
			throw += 2
		} else {
			sum += self.framePointsAt(throw)
			throw += 2
		}
	}
	return sum
}

// isStrike determines if a given throw is a strike or not. A strike is knocking
// down all pins in one throw.
func (self *Game) isStrike(throw int) bool {
	return self.rolls[throw] == allPins
}

// strikeBonusFor calculates and returns the strike bonus for a throw.
func (self *Game) strikeBonusFor(throw int) int {
	return allPins + self.framePointsAt(throw+1)
}

// isSpare determines if a given frame is a spare or not. A spare is knocking
// down all pins in one frame with two throws.
func (self *Game) isSpare(throw int) bool {
	return self.framePointsAt(throw) == allPins
}

// spareBonusFor calculates and returns the spare bonus for a throw.
func (self *Game) spareBonusFor(throw int) int {
	return allPins + self.rolls[throw+2]
}

// framePointsAt computes and returns the score in a frame specified by throw.
func (self *Game) framePointsAt(throw int) int {
	return self.rolls[throw] + self.rolls[throw+1]
}

const (
	// allPins is the number of pins allocated per fresh throw.
	allPins = 10

	// framesPerGame is the number of frames per bowling game.
	framesPerGame = 10

	// maxThrowsPerGame is the maximum number of throws possible in a single game.
	maxThrowsPerGame = 21
)

//////////////////////////////////////////////////////////////////////////////
// These types shouldn't be parsed as fixtures:

type TestFixtureWrongTestCase struct {
	*blah.Fixture
}
type TestFixtureWrongPackage struct {
	*gunit.Fixture2
}

type Hah interface {
	Hi() string
}

type BlahFixture struct {
	blah int
}

//////////////////////////////////////////////////////////////////////////////
`
