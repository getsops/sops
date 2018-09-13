package assertions

import (
	"fmt"
	"reflect"
	"time"
)

func (this *AssertionsFixture) TestShouldEqual() {
	this.fail(so(1, ShouldEqual), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(1, ShouldEqual, 1, 2), "This assertion requires exactly 1 comparison values (you provided 2).")
	this.fail(so(1, ShouldEqual, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.pass(so(1, ShouldEqual, 1))
	this.fail(so(1, ShouldEqual, 2), "2|1|Expected: '2' Actual: '1' (Should be equal)")
	this.fail(so(1, ShouldEqual, "1"), "1|1|Expected: '1' (string) Actual: '1' (int) (Should be equal, type mismatch)")

	this.pass(so(nil, ShouldEqual, nil))

	this.pass(so(true, ShouldEqual, true))
	this.fail(so(true, ShouldEqual, false), "false|true|Expected: 'false' Actual: 'true' (Should be equal)")

	this.pass(so("hi", ShouldEqual, "hi"))
	this.fail(so("hi", ShouldEqual, "bye"), "bye|hi|Expected: 'bye' Actual: 'hi' (Should be equal)")

	this.pass(so(42, ShouldEqual, uint(42)))

	this.fail(so(Thing1{"hi"}, ShouldEqual, Thing1{}), "{}|{hi}|Expected: '{}' Actual: '{hi}' (Should be equal)")
	this.fail(so(Thing1{"hi"}, ShouldEqual, Thing1{"hi"}), "{hi}|{hi}|Expected: '{hi}' Actual: '{hi}' (Should be equal)")
	this.fail(so(&Thing1{"hi"}, ShouldEqual, &Thing1{"hi"}), "&{hi}|&{hi}|Expected: '&{hi}' Actual: '&{hi}' (Should be equal)")

	this.fail(so(Thing1{}, ShouldEqual, Thing2{}), "{}|{}|Expected: '{}' Actual: '{}' (Should be equal)")

	this.pass(so(ThingWithEqualMethod{"hi"}, ShouldEqual, ThingWithEqualMethod{"hi"}))
	this.fail(so(ThingWithEqualMethod{"hi"}, ShouldEqual, ThingWithEqualMethod{"bye"}),
		"{bye}|{hi}|Expected: '{bye}' Actual: '{hi}' (Should be equal)")
}

func (this *AssertionsFixture) TestTimeEqual() {
	var (
		gopherCon, _ = time.LoadLocation("America/Denver")
		elsewhere, _ = time.LoadLocation("America/New_York")

		timeNow          = time.Now().In(gopherCon)
		timeNowElsewhere = timeNow.In(elsewhere)
		timeLater        = timeNow.Add(time.Nanosecond)
	)

	this.pass(so(timeNow, ShouldNotResemble, timeNowElsewhere)) // Differing *Location field prevents ShouldResemble!
	this.pass(so(timeNow, ShouldEqual, timeNowElsewhere))       // Time.Equal method used to determine exact instant.
	this.pass(so(timeNow, ShouldNotEqual, timeLater))
}

func (this *AssertionsFixture) TestShouldNotEqual() {
	this.fail(so(1, ShouldNotEqual), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(1, ShouldNotEqual, 1, 2), "This assertion requires exactly 1 comparison values (you provided 2).")
	this.fail(so(1, ShouldNotEqual, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.pass(so(1, ShouldNotEqual, 2))
	this.pass(so(1, ShouldNotEqual, "1"))
	this.fail(so(1, ShouldNotEqual, 1), "Expected '1' to NOT equal '1' (but it did)!")

	this.pass(so(true, ShouldNotEqual, false))
	this.fail(so(true, ShouldNotEqual, true), "Expected 'true' to NOT equal 'true' (but it did)!")

	this.pass(so("hi", ShouldNotEqual, "bye"))
	this.fail(so("hi", ShouldNotEqual, "hi"), "Expected 'hi' to NOT equal 'hi' (but it did)!")

	this.pass(so(&Thing1{"hi"}, ShouldNotEqual, &Thing1{"hi"}))
	this.pass(so(Thing1{"hi"}, ShouldNotEqual, Thing1{"hi"}))
	this.pass(so(Thing1{}, ShouldNotEqual, Thing1{}))
	this.pass(so(Thing1{}, ShouldNotEqual, Thing2{}))
}

func (this *AssertionsFixture) TestShouldAlmostEqual() {
	this.fail(so(1, ShouldAlmostEqual), "This assertion requires exactly one comparison value and an optional delta (you provided neither)")
	this.fail(so(1, ShouldAlmostEqual, 1, 2, 3), "This assertion requires exactly one comparison value and an optional delta (you provided more values)")
	this.fail(so(1, ShouldAlmostEqual, "1"), "The comparison value must be a numerical type, but was: string")
	this.fail(so(1, ShouldAlmostEqual, 1, "1"), "The delta value must be a numerical type, but was: string")
	this.fail(so("1", ShouldAlmostEqual, 1), "The actual value must be a numerical type, but was: string")

	// with the default delta
	this.pass(so(0.99999999999999, ShouldAlmostEqual, uint(1)))
	this.pass(so(1, ShouldAlmostEqual, 0.99999999999999))
	this.pass(so(1.3612499999999996, ShouldAlmostEqual, 1.36125))
	this.pass(so(0.7285312499999999, ShouldAlmostEqual, 0.72853125))
	this.fail(so(1, ShouldAlmostEqual, .99), "Expected '1' to almost equal '0.99' (but it didn't)!")

	// with a different delta
	this.pass(so(100.0, ShouldAlmostEqual, 110.0, 10.0))
	this.fail(so(100.0, ShouldAlmostEqual, 111.0, 10.5), "Expected '100' to almost equal '111' (but it didn't)!")

	// various ints should work
	this.pass(so(100, ShouldAlmostEqual, 100.0))
	this.pass(so(int(100), ShouldAlmostEqual, 100.0))
	this.pass(so(int8(100), ShouldAlmostEqual, 100.0))
	this.pass(so(int16(100), ShouldAlmostEqual, 100.0))
	this.pass(so(int32(100), ShouldAlmostEqual, 100.0))
	this.pass(so(int64(100), ShouldAlmostEqual, 100.0))
	this.pass(so(uint(100), ShouldAlmostEqual, 100.0))
	this.pass(so(uint8(100), ShouldAlmostEqual, 100.0))
	this.pass(so(uint16(100), ShouldAlmostEqual, 100.0))
	this.pass(so(uint32(100), ShouldAlmostEqual, 100.0))
	this.pass(so(uint64(100), ShouldAlmostEqual, 100.0))
	this.pass(so(100, ShouldAlmostEqual, 100.0))
	this.fail(so(100, ShouldAlmostEqual, 99.0), "Expected '100' to almost equal '99' (but it didn't)!")

	// floats should work
	this.pass(so(float64(100.0), ShouldAlmostEqual, float32(100.0)))
	this.fail(so(float32(100.0), ShouldAlmostEqual, 99.0, float32(0.1)), "Expected '100' to almost equal '99' (but it didn't)!")
}

func (this *AssertionsFixture) TestShouldNotAlmostEqual() {
	this.fail(so(1, ShouldNotAlmostEqual), "This assertion requires exactly one comparison value and an optional delta (you provided neither)")
	this.fail(so(1, ShouldNotAlmostEqual, 1, 2, 3), "This assertion requires exactly one comparison value and an optional delta (you provided more values)")

	// with the default delta
	this.fail(so(1, ShouldNotAlmostEqual, .99999999999999), "Expected '1' to NOT almost equal '0.99999999999999' (but it did)!")
	this.fail(so(1.3612499999999996, ShouldNotAlmostEqual, 1.36125), "Expected '1.3612499999999996' to NOT almost equal '1.36125' (but it did)!")
	this.pass(so(1, ShouldNotAlmostEqual, .99))

	// with a different delta
	this.fail(so(100.0, ShouldNotAlmostEqual, 110.0, 10.0), "Expected '100' to NOT almost equal '110' (but it did)!")
	this.pass(so(100.0, ShouldNotAlmostEqual, 111.0, 10.5))

	// ints should work
	this.fail(so(100, ShouldNotAlmostEqual, 100.0), "Expected '100' to NOT almost equal '100' (but it did)!")
	this.pass(so(100, ShouldNotAlmostEqual, 99.0))

	// float32 should work
	this.fail(so(float64(100.0), ShouldNotAlmostEqual, float32(100.0)), "Expected '100' to NOT almost equal '100' (but it did)!")
	this.pass(so(float32(100.0), ShouldNotAlmostEqual, 99.0, float32(0.1)))
}

func (this *AssertionsFixture) TestShouldResemble() {
	this.fail(so(Thing1{"hi"}, ShouldResemble), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(Thing1{"hi"}, ShouldResemble, Thing1{"hi"}, Thing1{"hi"}), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.pass(so(Thing1{"hi"}, ShouldResemble, Thing1{"hi"}))
	this.fail(so(Thing1{"hi"}, ShouldResemble, Thing1{"bye"}), `{bye}|{hi}|Expected: 'assertions.Thing1{a:"bye"}' Actual: 'assertions.Thing1{a:"hi"}' (Should resemble)!`)

	var (
		a []int
		b []int = []int{}
	)

	this.fail(so(a, ShouldResemble, b), `[]|[]|Expected: '[]int{}' Actual: '[]int(nil)' (Should resemble)!`)
	this.fail(so(2, ShouldResemble, 1), `1|2|Expected: '1' Actual: '2' (Should resemble)!`)

	this.fail(so(StringStringMapAlias{"hi": "bye"}, ShouldResemble, map[string]string{"hi": "bye"}),
		`map[hi:bye]|map[hi:bye]|Expected: 'map[string]string{"hi":"bye"}' Actual: 'assertions.StringStringMapAlias{"hi":"bye"}' (Should resemble)!`)
	this.fail(so(StringSliceAlias{"hi", "bye"}, ShouldResemble, []string{"hi", "bye"}),
		`[hi bye]|[hi bye]|Expected: '[]string{"hi", "bye"}' Actual: 'assertions.StringSliceAlias{"hi", "bye"}' (Should resemble)!`)

	// some types come out looking the same when represented with "%#v" so we show type mismatch info:
	this.fail(so(StringAlias("hi"), ShouldResemble, "hi"), `hi|hi|Expected: '"hi"' Actual: 'assertions.StringAlias("hi")' (Should resemble)!`)
	this.fail(so(IntAlias(42), ShouldResemble, 42), `42|42|Expected: '42' Actual: 'assertions.IntAlias(42)' (Should resemble)!`)
}

func (this *AssertionsFixture) TestShouldEqualJSON() {
	this.fail(so("hi", ShouldEqualJSON), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("hi", ShouldEqualJSON, "hi", "hi"), "This assertion requires exactly 1 comparison values (you provided 2).")

	// basic identity of keys/values
	this.pass(so(`{"my":"val"}`, ShouldEqualJSON, `{"my":"val"}`))
	this.fail(so(`{"my":"val"}`, ShouldEqualJSON, `{"your":"val"}`),
		`{"your":"val"}|{"my":"val"}|Expected: '{"your":"val"}' Actual: '{"my":"val"}' (Should be equal)`)

	// out of order values causes comparison failure:
	this.pass(so(`{"key0":"val0","key1":"val1"}`, ShouldEqualJSON, `{"key1":"val1","key0":"val0"}`))
	this.fail(so(`{"key0":"val0","key1":"val1"}`, ShouldEqualJSON, `{"key1":"val0","key0":"val0"}`),
		`{"key0":"val0","key1":"val0"}|{"key0":"val0","key1":"val1"}|Expected: '{"key0":"val0","key1":"val0"}' Actual: '{"key0":"val0","key1":"val1"}' (Should be equal)`)

	// missing values causes comparison failure:
	this.fail(so(
		`{"key0":"val0","key1":"val1"}`,
		ShouldEqualJSON,
		`{"key1":"val0"}`),
		`{"key1":"val0"}|{"key0":"val0","key1":"val1"}|Expected: '{"key1":"val0"}' Actual: '{"key0":"val0","key1":"val1"}' (Should be equal)`)

	// whitespace shouldn't interfere with comparison:
	this.pass(so("\n{ \"my\"  :   \"val\"\n}", ShouldEqualJSON, `{"my":"val"}`))

	// Invalid JSON for either actual or expected value is invalid:
	this.fail(so("{}", ShouldEqualJSON, ""), "Expected value not valid JSON: unexpected end of JSON input")
	this.fail(so("", ShouldEqualJSON, "{}"), "Actual value not valid JSON: unexpected end of JSON input")
}

func (this *AssertionsFixture) TestShouldNotResemble() {
	this.fail(so(Thing1{"hi"}, ShouldNotResemble), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(Thing1{"hi"}, ShouldNotResemble, Thing1{"hi"}, Thing1{"hi"}), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.pass(so(Thing1{"hi"}, ShouldNotResemble, Thing1{"bye"}))
	this.fail(so(Thing1{"hi"}, ShouldNotResemble, Thing1{"hi"}),
		`Expected '"assertions.Thing1{a:\"hi\"}"' to NOT resemble '"assertions.Thing1{a:\"hi\"}"' (but it did)!`)

	this.pass(so(map[string]string{"hi": "bye"}, ShouldResemble, map[string]string{"hi": "bye"}))
	this.pass(so(IntAlias(42), ShouldNotResemble, 42))

	this.pass(so(StringSliceAlias{"hi", "bye"}, ShouldNotResemble, []string{"hi", "bye"}))
}

func (this *AssertionsFixture) TestShouldPointTo() {
	t1 := &Thing1{}
	t2 := t1
	t3 := &Thing1{}

	pointer1 := reflect.ValueOf(t1).Pointer()
	pointer3 := reflect.ValueOf(t3).Pointer()

	this.fail(so(t1, ShouldPointTo), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(t1, ShouldPointTo, t2, t3), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.pass(so(t1, ShouldPointTo, t2))
	this.fail(so(t1, ShouldPointTo, t3), fmt.Sprintf(
		"%v|%v|Expected '&{a:}' (address: '%v') and '&{a:}' (address: '%v') to be the same address (but their weren't)!",
		pointer3, pointer1, pointer1, pointer3))

	t4 := Thing1{}
	t5 := t4

	this.fail(so(t4, ShouldPointTo, t5), "Both arguments should be pointers (the first was not)!")
	this.fail(so(&t4, ShouldPointTo, t5), "Both arguments should be pointers (the second was not)!")
	this.fail(so(nil, ShouldPointTo, nil), "Both arguments should be pointers (the first was nil)!")
	this.fail(so(&t4, ShouldPointTo, nil), "Both arguments should be pointers (the second was nil)!")
}

func (this *AssertionsFixture) TestShouldNotPointTo() {
	t1 := &Thing1{}
	t2 := t1
	t3 := &Thing1{}

	pointer1 := reflect.ValueOf(t1).Pointer()

	this.fail(so(t1, ShouldNotPointTo), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(t1, ShouldNotPointTo, t2, t3), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.pass(so(t1, ShouldNotPointTo, t3))
	this.fail(so(t1, ShouldNotPointTo, t2), fmt.Sprintf("Expected '&{a:}' and '&{a:}' to be different references (but they matched: '%v')!", pointer1))

	t4 := Thing1{}
	t5 := t4

	this.fail(so(t4, ShouldNotPointTo, t5), "Both arguments should be pointers (the first was not)!")
	this.fail(so(&t4, ShouldNotPointTo, t5), "Both arguments should be pointers (the second was not)!")
	this.fail(so(nil, ShouldNotPointTo, nil), "Both arguments should be pointers (the first was nil)!")
	this.fail(so(&t4, ShouldNotPointTo, nil), "Both arguments should be pointers (the second was nil)!")
}

func (this *AssertionsFixture) TestShouldBeNil() {
	this.fail(so(nil, ShouldBeNil, nil, nil, nil), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(nil, ShouldBeNil, nil), "This assertion requires exactly 0 comparison values (you provided 1).")

	this.pass(so(nil, ShouldBeNil))
	this.fail(so(1, ShouldBeNil), "Expected: nil Actual: '1'")

	var thing ThingInterface
	this.pass(so(thing, ShouldBeNil))
	thing = &ThingImplementation{}
	this.fail(so(thing, ShouldBeNil), "Expected: nil Actual: '&{}'")

	var thingOne *Thing1
	this.pass(so(thingOne, ShouldBeNil))

	var nilSlice []int = nil
	this.pass(so(nilSlice, ShouldBeNil))

	var nilMap map[string]string = nil
	this.pass(so(nilMap, ShouldBeNil))

	var nilChannel chan int = nil
	this.pass(so(nilChannel, ShouldBeNil))

	var nilFunc func() = nil
	this.pass(so(nilFunc, ShouldBeNil))

	var nilInterface interface{} = nil
	this.pass(so(nilInterface, ShouldBeNil))
}

func (this *AssertionsFixture) TestShouldNotBeNil() {
	this.fail(so(nil, ShouldNotBeNil, nil, nil, nil), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(nil, ShouldNotBeNil, nil), "This assertion requires exactly 0 comparison values (you provided 1).")

	this.fail(so(nil, ShouldNotBeNil), "Expected '<nil>' to NOT be nil (but it was)!")
	this.pass(so(1, ShouldNotBeNil))

	var thing ThingInterface
	this.fail(so(thing, ShouldNotBeNil), "Expected '<nil>' to NOT be nil (but it was)!")
	thing = &ThingImplementation{}
	this.pass(so(thing, ShouldNotBeNil))
}

func (this *AssertionsFixture) TestShouldBeTrue() {
	this.fail(so(true, ShouldBeTrue, 1, 2, 3), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(true, ShouldBeTrue, 1), "This assertion requires exactly 0 comparison values (you provided 1).")

	this.fail(so(false, ShouldBeTrue), "Expected: true Actual: false")
	this.fail(so(1, ShouldBeTrue), "Expected: true Actual: 1")
	this.pass(so(true, ShouldBeTrue))
}

func (this *AssertionsFixture) TestShouldBeFalse() {
	this.fail(so(false, ShouldBeFalse, 1, 2, 3), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(false, ShouldBeFalse, 1), "This assertion requires exactly 0 comparison values (you provided 1).")

	this.fail(so(true, ShouldBeFalse), "Expected: false Actual: true")
	this.fail(so(1, ShouldBeFalse), "Expected: false Actual: 1")
	this.pass(so(false, ShouldBeFalse))
}

func (this *AssertionsFixture) TestShouldBeZeroValue() {
	this.fail(so(0, ShouldBeZeroValue, 1, 2, 3), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(false, ShouldBeZeroValue, true), "This assertion requires exactly 0 comparison values (you provided 1).")

	this.fail(so(1, ShouldBeZeroValue), "0|1|'1' should have been the zero value")                                       //"Expected: (zero value) Actual: 1")
	this.fail(so(true, ShouldBeZeroValue), "false|true|'true' should have been the zero value")                          //"Expected: (zero value) Actual: true")
	this.fail(so("123", ShouldBeZeroValue), "|123|'123' should have been the zero value")                                //"Expected: (zero value) Actual: 123")
	this.fail(so(" ", ShouldBeZeroValue), "| |' ' should have been the zero value")                                      //"Expected: (zero value) Actual:  ")
	this.fail(so([]string{"Nonempty"}, ShouldBeZeroValue), "[]|[Nonempty]|'[Nonempty]' should have been the zero value") //"Expected: (zero value) Actual: [Nonempty]")
	this.fail(so(struct{ a string }{a: "asdf"}, ShouldBeZeroValue), "{}|{asdf}|'{a:asdf}' should have been the zero value")
	this.pass(so(0, ShouldBeZeroValue))
	this.pass(so(false, ShouldBeZeroValue))
	this.pass(so("", ShouldBeZeroValue))
	this.pass(so(struct{}{}, ShouldBeZeroValue))
}

func (this *AssertionsFixture) TestShouldNotBeZeroValue() {
	this.fail(so(0, ShouldNotBeZeroValue, 1, 2, 3), "This assertion requires exactly 0 comparison values (you provided 3).")
	this.fail(so(false, ShouldNotBeZeroValue, true), "This assertion requires exactly 0 comparison values (you provided 1).")

	this.fail(so(0, ShouldNotBeZeroValue), "0|0|'0' should NOT have been the zero value")
	this.fail(so(false, ShouldNotBeZeroValue), "false|false|'false' should NOT have been the zero value")
	this.fail(so("", ShouldNotBeZeroValue), "||'' should NOT have been the zero value")
	this.fail(so(struct{}{}, ShouldNotBeZeroValue), "{}|{}|'{}' should NOT have been the zero value")

	this.pass(so(1, ShouldNotBeZeroValue))
	this.pass(so(true, ShouldNotBeZeroValue))
	this.pass(so("123", ShouldNotBeZeroValue))
	this.pass(so(" ", ShouldNotBeZeroValue))
	this.pass(so([]string{"Nonempty"}, ShouldNotBeZeroValue))
	this.pass(so(struct{ a string }{a: "asdf"}, ShouldNotBeZeroValue))
}
