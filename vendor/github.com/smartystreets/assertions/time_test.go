package assertions

import (
	"fmt"
	"time"
)

func (this *AssertionsFixture) TestShouldHappenBefore() {
	this.fail(so(0, ShouldHappenBefore), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenBefore, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenBefore, 1), shouldUseTimes)
	this.fail(so(0, ShouldHappenBefore, time.Now()), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenBefore, 0), shouldUseTimes)

	this.fail(so(january3, ShouldHappenBefore, january1), fmt.Sprintf("Expected '%s' to happen before '%s' (it happened '48h0m0s' after)!", pretty(january3), pretty(january1)))
	this.fail(so(january3, ShouldHappenBefore, january3), fmt.Sprintf("Expected '%s' to happen before '%s' (it happened '0s' after)!", pretty(january3), pretty(january3)))
	this.pass(so(january1, ShouldHappenBefore, january3))
}

func (this *AssertionsFixture) TestShouldHappenOnOrBefore() {
	this.fail(so(0, ShouldHappenOnOrBefore), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenOnOrBefore, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenOnOrBefore, 1), shouldUseTimes)
	this.fail(so(0, ShouldHappenOnOrBefore, time.Now()), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenOnOrBefore, 0), shouldUseTimes)

	this.fail(so(january3, ShouldHappenOnOrBefore, january1), fmt.Sprintf("Expected '%s' to happen before '%s' (it happened '48h0m0s' after)!", pretty(january3), pretty(january1)))
	this.pass(so(january3, ShouldHappenOnOrBefore, january3))
	this.pass(so(january1, ShouldHappenOnOrBefore, january3))
}

func (this *AssertionsFixture) TestShouldHappenAfter() {
	this.fail(so(0, ShouldHappenAfter), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenAfter, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenAfter, 1), shouldUseTimes)
	this.fail(so(0, ShouldHappenAfter, time.Now()), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenAfter, 0), shouldUseTimes)

	this.fail(so(january1, ShouldHappenAfter, january2), fmt.Sprintf("Expected '%s' to happen after '%s' (it happened '24h0m0s' before)!", pretty(january1), pretty(january2)))
	this.fail(so(january1, ShouldHappenAfter, january1), fmt.Sprintf("Expected '%s' to happen after '%s' (it happened '0s' before)!", pretty(january1), pretty(january1)))
	this.pass(so(january3, ShouldHappenAfter, january1))
}

func (this *AssertionsFixture) TestShouldHappenOnOrAfter() {
	this.fail(so(0, ShouldHappenOnOrAfter), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenOnOrAfter, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenOnOrAfter, 1), shouldUseTimes)
	this.fail(so(0, ShouldHappenOnOrAfter, time.Now()), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenOnOrAfter, 0), shouldUseTimes)

	this.fail(so(january1, ShouldHappenOnOrAfter, january2), fmt.Sprintf("Expected '%s' to happen after '%s' (it happened '24h0m0s' before)!", pretty(january1), pretty(january2)))
	this.pass(so(january1, ShouldHappenOnOrAfter, january1))
	this.pass(so(january3, ShouldHappenOnOrAfter, january1))
}

func (this *AssertionsFixture) TestShouldHappenBetween() {
	this.fail(so(0, ShouldHappenBetween), "This assertion requires exactly 2 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenBetween, 1, 2, 3), "This assertion requires exactly 2 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenBetween, 1, 2), shouldUseTimes)
	this.fail(so(0, ShouldHappenBetween, time.Now(), time.Now()), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenBetween, 0, time.Now()), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenBetween, time.Now(), 9), shouldUseTimes)

	this.fail(so(january1, ShouldHappenBetween, january2, january4), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '24h0m0s' outside threshold)!", pretty(january1), pretty(january2), pretty(january4)))
	this.fail(so(january2, ShouldHappenBetween, january2, january4), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '0s' outside threshold)!", pretty(january2), pretty(january2), pretty(january4)))
	this.pass(so(january3, ShouldHappenBetween, january2, january4))
	this.fail(so(january4, ShouldHappenBetween, january2, january4), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '0s' outside threshold)!", pretty(january4), pretty(january2), pretty(january4)))
	this.fail(so(january5, ShouldHappenBetween, january2, january4), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '24h0m0s' outside threshold)!", pretty(january5), pretty(january2), pretty(january4)))
}

func (this *AssertionsFixture) TestShouldHappenOnOrBetween() {
	this.fail(so(0, ShouldHappenOnOrBetween), "This assertion requires exactly 2 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenOnOrBetween, 1, 2, 3), "This assertion requires exactly 2 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenOnOrBetween, 1, time.Now()), shouldUseTimes)
	this.fail(so(0, ShouldHappenOnOrBetween, time.Now(), 1), shouldUseTimes)
	this.fail(so(time.Now(), ShouldHappenOnOrBetween, 0, 1), shouldUseTimes)

	this.fail(so(january1, ShouldHappenOnOrBetween, january2, january4), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '24h0m0s' outside threshold)!", pretty(january1), pretty(january2), pretty(january4)))
	this.pass(so(january2, ShouldHappenOnOrBetween, january2, january4))
	this.pass(so(january3, ShouldHappenOnOrBetween, january2, january4))
	this.pass(so(january4, ShouldHappenOnOrBetween, january2, january4))
	this.fail(so(january5, ShouldHappenOnOrBetween, january2, january4), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '24h0m0s' outside threshold)!", pretty(january5), pretty(january2), pretty(january4)))
}

func (this *AssertionsFixture) TestShouldNotHappenOnOrBetween() {
	this.fail(so(0, ShouldNotHappenOnOrBetween), "This assertion requires exactly 2 comparison values (you provided 0).")
	this.fail(so(0, ShouldNotHappenOnOrBetween, 1, 2, 3), "This assertion requires exactly 2 comparison values (you provided 3).")

	this.fail(so(0, ShouldNotHappenOnOrBetween, 1, time.Now()), shouldUseTimes)
	this.fail(so(0, ShouldNotHappenOnOrBetween, time.Now(), 1), shouldUseTimes)
	this.fail(so(time.Now(), ShouldNotHappenOnOrBetween, 0, 1), shouldUseTimes)

	this.pass(so(january1, ShouldNotHappenOnOrBetween, january2, january4))
	this.fail(so(january2, ShouldNotHappenOnOrBetween, january2, january4), fmt.Sprintf("Expected '%s' to NOT happen on or between '%s' and '%s' (but it did)!", pretty(january2), pretty(january2), pretty(january4)))
	this.fail(so(january3, ShouldNotHappenOnOrBetween, january2, january4), fmt.Sprintf("Expected '%s' to NOT happen on or between '%s' and '%s' (but it did)!", pretty(january3), pretty(january2), pretty(january4)))
	this.fail(so(january4, ShouldNotHappenOnOrBetween, january2, january4), fmt.Sprintf("Expected '%s' to NOT happen on or between '%s' and '%s' (but it did)!", pretty(january4), pretty(january2), pretty(january4)))
	this.pass(so(january5, ShouldNotHappenOnOrBetween, january2, january4))
}

func (this *AssertionsFixture) TestShouldHappenWithin() {
	this.fail(so(0, ShouldHappenWithin), "This assertion requires exactly 2 comparison values (you provided 0).")
	this.fail(so(0, ShouldHappenWithin, 1, 2, 3), "This assertion requires exactly 2 comparison values (you provided 3).")

	this.fail(so(0, ShouldHappenWithin, 1, 2), shouldUseDurationAndTime)
	this.fail(so(0, ShouldHappenWithin, oneDay, time.Now()), shouldUseDurationAndTime)
	this.fail(so(time.Now(), ShouldHappenWithin, 0, time.Now()), shouldUseDurationAndTime)

	this.fail(so(january1, ShouldHappenWithin, oneDay, january3), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '24h0m0s' outside threshold)!", pretty(january1), pretty(january2), pretty(january4)))
	this.pass(so(january2, ShouldHappenWithin, oneDay, january3))
	this.pass(so(january3, ShouldHappenWithin, oneDay, january3))
	this.pass(so(january4, ShouldHappenWithin, oneDay, january3))
	this.fail(so(january5, ShouldHappenWithin, oneDay, january3), fmt.Sprintf("Expected '%s' to happen between '%s' and '%s' (it happened '24h0m0s' outside threshold)!", pretty(january5), pretty(january2), pretty(january4)))
}

func (this *AssertionsFixture) TestShouldNotHappenWithin() {
	this.fail(so(0, ShouldNotHappenWithin), "This assertion requires exactly 2 comparison values (you provided 0).")
	this.fail(so(0, ShouldNotHappenWithin, 1, 2, 3), "This assertion requires exactly 2 comparison values (you provided 3).")

	this.fail(so(0, ShouldNotHappenWithin, 1, 2), shouldUseDurationAndTime)
	this.fail(so(0, ShouldNotHappenWithin, oneDay, time.Now()), shouldUseDurationAndTime)
	this.fail(so(time.Now(), ShouldNotHappenWithin, 0, time.Now()), shouldUseDurationAndTime)

	this.pass(so(january1, ShouldNotHappenWithin, oneDay, january3))
	this.fail(so(january2, ShouldNotHappenWithin, oneDay, january3), fmt.Sprintf("Expected '%s' to NOT happen on or between '%s' and '%s' (but it did)!", pretty(january2), pretty(january2), pretty(january4)))
	this.fail(so(january3, ShouldNotHappenWithin, oneDay, january3), fmt.Sprintf("Expected '%s' to NOT happen on or between '%s' and '%s' (but it did)!", pretty(january3), pretty(january2), pretty(january4)))
	this.fail(so(january4, ShouldNotHappenWithin, oneDay, january3), fmt.Sprintf("Expected '%s' to NOT happen on or between '%s' and '%s' (but it did)!", pretty(january4), pretty(january2), pretty(january4)))
	this.pass(so(january5, ShouldNotHappenWithin, oneDay, january3))
}

func (this *AssertionsFixture) TestShouldBeChronological() {
	this.fail(so(0, ShouldBeChronological, 1, 2, 3), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(0, ShouldBeChronological), shouldUseTimeSlice)
	this.fail(so([]time.Time{january5, january1}, ShouldBeChronological),
		"The 'Time' at index [1] should have happened after the previous one (but it didn't!):\n  [0]: 2013-01-05 00:00:00 +0000 UTC\n  [1]: 2013-01-01 00:00:00 +0000 UTC (see, it happened before!)")

	this.pass(so([]time.Time{january1, january2, january3, january4, january5}, ShouldBeChronological))
}

const layout = "2006-01-02 15:04"

var january1, _ = time.Parse(layout, "2013-01-01 00:00")
var january2, _ = time.Parse(layout, "2013-01-02 00:00")
var january3, _ = time.Parse(layout, "2013-01-03 00:00")
var january4, _ = time.Parse(layout, "2013-01-04 00:00")
var january5, _ = time.Parse(layout, "2013-01-05 00:00")

var oneDay, _ = time.ParseDuration("24h0m0s")
var twoDays, _ = time.ParseDuration("48h0m0s")

func pretty(t time.Time) string {
	return fmt.Sprintf("%v", t)
}
