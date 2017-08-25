package assertions

func (this *AssertionsFixture) TestShouldStartWith() {
	this.fail(so("", ShouldStartWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("", ShouldStartWith, "asdf", "asdf"), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.pass(so("", ShouldStartWith, ""))
	this.fail(so("", ShouldStartWith, "x"), "x||Expected '' to start with 'x' (but it didn't)!")
	this.pass(so("abc", ShouldStartWith, "abc"))
	this.fail(so("abc", ShouldStartWith, "abcd"), "abcd|abc|Expected 'abc' to start with 'abcd' (but it didn't)!")

	this.pass(so("superman", ShouldStartWith, "super"))
	this.fail(so("superman", ShouldStartWith, "bat"), "bat|sup...|Expected 'superman' to start with 'bat' (but it didn't)!")
	this.fail(so("superman", ShouldStartWith, "man"), "man|sup...|Expected 'superman' to start with 'man' (but it didn't)!")

	this.fail(so(1, ShouldStartWith, 2), "Both arguments to this assertion must be strings (you provided int and int).")
}

func (this *AssertionsFixture) TestShouldNotStartWith() {
	this.fail(so("", ShouldNotStartWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("", ShouldNotStartWith, "asdf", "asdf"), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.fail(so("", ShouldNotStartWith, ""), "Expected '<empty>' NOT to start with '<empty>' (but it did)!")
	this.fail(so("superman", ShouldNotStartWith, "super"), "Expected 'superman' NOT to start with 'super' (but it did)!")
	this.pass(so("superman", ShouldNotStartWith, "bat"))
	this.pass(so("superman", ShouldNotStartWith, "man"))

	this.fail(so(1, ShouldNotStartWith, 2), "Both arguments to this assertion must be strings (you provided int and int).")
}

func (this *AssertionsFixture) TestShouldEndWith() {
	this.fail(so("", ShouldEndWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("", ShouldEndWith, "", ""), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.pass(so("", ShouldEndWith, ""))
	this.fail(so("", ShouldEndWith, "z"), "z||Expected '' to end with 'z' (but it didn't)!")
	this.pass(so("xyz", ShouldEndWith, "xyz"))
	this.fail(so("xyz", ShouldEndWith, "wxyz"), "wxyz|xyz|Expected 'xyz' to end with 'wxyz' (but it didn't)!")

	this.pass(so("superman", ShouldEndWith, "man"))
	this.fail(so("superman", ShouldEndWith, "super"), "super|...erman|Expected 'superman' to end with 'super' (but it didn't)!")
	this.fail(so("superman", ShouldEndWith, "blah"), "blah|...rman|Expected 'superman' to end with 'blah' (but it didn't)!")

	this.fail(so(1, ShouldEndWith, 2), "Both arguments to this assertion must be strings (you provided int and int).")
}

func (this *AssertionsFixture) TestShouldNotEndWith() {
	this.fail(so("", ShouldNotEndWith), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("", ShouldNotEndWith, "", ""), "This assertion requires exactly 1 comparison values (you provided 2).")

	this.fail(so("", ShouldNotEndWith, ""), "Expected '<empty>' NOT to end with '<empty>' (but it did)!")
	this.fail(so("superman", ShouldNotEndWith, "man"), "Expected 'superman' NOT to end with 'man' (but it did)!")
	this.pass(so("superman", ShouldNotEndWith, "super"))

	this.fail(so(1, ShouldNotEndWith, 2), "Both arguments to this assertion must be strings (you provided int and int).")
}

func (this *AssertionsFixture) TestShouldContainSubstring() {
	this.fail(so("asdf", ShouldContainSubstring), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("asdf", ShouldContainSubstring, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(123, ShouldContainSubstring, 23), "Both arguments to this assertion must be strings (you provided int and int).")

	this.pass(so("asdf", ShouldContainSubstring, "sd"))
	this.fail(so("qwer", ShouldContainSubstring, "sd"), "sd|qwer|Expected 'qwer' to contain substring 'sd' (but it didn't)!")
}

func (this *AssertionsFixture) TestShouldNotContainSubstring() {
	this.fail(so("asdf", ShouldNotContainSubstring), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so("asdf", ShouldNotContainSubstring, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	this.fail(so(123, ShouldNotContainSubstring, 23), "Both arguments to this assertion must be strings (you provided int and int).")

	this.pass(so("qwer", ShouldNotContainSubstring, "sd"))
	this.fail(so("asdf", ShouldNotContainSubstring, "sd"), "Expected 'asdf' NOT to contain substring 'sd' (but it did)!")
}

func (this *AssertionsFixture) TestShouldBeBlank() {
	this.fail(so("", ShouldBeBlank, "adsf"), "This assertion requires exactly 0 comparison values (you provided 1).")
	this.fail(so(1, ShouldBeBlank), "The argument to this assertion must be a string (you provided int).")

	this.fail(so("asdf", ShouldBeBlank), "|asdf|Expected 'asdf' to be blank (but it wasn't)!")
	this.pass(so("", ShouldBeBlank))
}

func (this *AssertionsFixture) TestShouldNotBeBlank() {
	this.fail(so("", ShouldNotBeBlank, "adsf"), "This assertion requires exactly 0 comparison values (you provided 1).")
	this.fail(so(1, ShouldNotBeBlank), "The argument to this assertion must be a string (you provided int).")

	this.fail(so("", ShouldNotBeBlank), "Expected value to NOT be blank (but it was)!")
	this.pass(so("asdf", ShouldNotBeBlank))
}

func (this *AssertionsFixture) TestShouldEqualWithout() {
	this.fail(so("", ShouldEqualWithout, ""), "This assertion requires exactly 2 comparison values (you provided 1).")
	this.fail(so(1, ShouldEqualWithout, 2, 3), "All arguments to this assertion must be strings (you provided: [int int int]).")

	this.fail(so("asdf", ShouldEqualWithout, "qwer", "q"), "Expected 'asdf' to equal 'qwer' but without any 'q' (but it didn't).")
	this.pass(so("asdf", ShouldEqualWithout, "df", "as"))
}

func (this *AssertionsFixture) TestShouldEqualTrimSpace() {
	this.fail(so(" asdf ", ShouldEqualTrimSpace), "This assertion requires exactly 1 comparison values (you provided 0).")
	this.fail(so(1, ShouldEqualTrimSpace, 2), "Both arguments to this assertion must be strings (you provided int and int).")

	this.fail(so("asdf", ShouldEqualTrimSpace, "qwer"), "qwer|asdf|Expected: 'qwer' Actual: 'asdf' (Should be equal)")
	this.pass(so(" asdf\t\n", ShouldEqualTrimSpace, "asdf"))
}
