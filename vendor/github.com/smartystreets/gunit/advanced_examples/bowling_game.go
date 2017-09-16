package examples

type Game struct {
	rolls   [maxThrowsPerGame]int
	current int
}

func NewGame() *Game {
	return new(Game)
}

func (this *Game) Roll(pins int) {
	this.rolls[this.current] = pins
	this.current++
}

func (this *Game) Score() (score int) {
	for frameIndex, frame := 0, 0; frame < framesPerGame; frame++ {
		if this.isStrike(frameIndex) {
			score += allPins + this.nextTwoBallsForStrike(frameIndex)
			frameIndex += 1
		} else if this.isSpare(frameIndex) {
			score += allPins + this.nextBallForSpare(frameIndex)
			frameIndex += 2
		} else {
			score += this.twoBallsInFrame(frameIndex)
			frameIndex += 2
		}
	}
	return score
}

func (this *Game) isSpare(frame int) bool  { return this.twoBallsInFrame(frame) == allPins }
func (this *Game) isStrike(frame int) bool { return this.rolls[frame] == allPins }

func (this *Game) nextTwoBallsForStrike(frame int) int { return this.twoBallsInFrame(frame + 1) }
func (this *Game) nextBallForSpare(frame int) int      { return this.rolls[frame+2] }
func (this *Game) twoBallsInFrame(frame int) int       { return this.rolls[frame] + this.rolls[frame+1] }

const (
	allPins          = 10
	framesPerGame    = 10
	maxThrowsPerGame = 21
)
