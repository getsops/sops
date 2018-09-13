// Integration tests

package check_test

import (
	. "gopkg.in/check.v1"
)

// -----------------------------------------------------------------------
// Integration test suite.

type integrationS struct{}

var _ = Suite(&integrationS{})

type integrationTestHelper struct{}

func (s *integrationTestHelper) TestMultiLineStringEqualFails(c *C) {
	c.Check("foo\nbar\nbaz\nboom\n", Equals, "foo\nbaar\nbaz\nboom\n")
}

func (s *integrationTestHelper) TestStringEqualFails(c *C) {
	c.Check("foo", Equals, "bar")
}

func (s *integrationTestHelper) TestIntEqualFails(c *C) {
	c.Check(42, Equals, 43)
}

type complexStruct struct {
	r, i int
}

func (s *integrationTestHelper) TestStructEqualFails(c *C) {
	c.Check(complexStruct{1, 2}, Equals, complexStruct{3, 4})
}

func (s *integrationS) TestOutput(c *C) {
	helper := integrationTestHelper{}
	output := String{}
	Run(&helper, &RunConf{Output: &output})
	c.Assert(output.value, Equals, `
----------------------------------------------------------------------
FAIL: integration_test.go:26: integrationTestHelper.TestIntEqualFails

integration_test.go:27:
    c.Check(42, Equals, 43)
... obtained int = 42
... expected int = 43


----------------------------------------------------------------------
FAIL: integration_test.go:18: integrationTestHelper.TestMultiLineStringEqualFails

integration_test.go:19:
    c.Check("foo\nbar\nbaz\nboom\n", Equals, "foo\nbaar\nbaz\nboom\n")
... obtained string = "" +
...     "foo\n" +
...     "bar\n" +
...     "baz\n" +
...     "boom\n"
... expected string = "" +
...     "foo\n" +
...     "baar\n" +
...     "baz\n" +
...     "boom\n"
... String difference:
...     [1]: "bar" != "baar"



----------------------------------------------------------------------
FAIL: integration_test.go:22: integrationTestHelper.TestStringEqualFails

integration_test.go:23:
    c.Check("foo", Equals, "bar")
... obtained string = "foo"
... expected string = "bar"


----------------------------------------------------------------------
FAIL: integration_test.go:34: integrationTestHelper.TestStructEqualFails

integration_test.go:35:
    c.Check(complexStruct{1, 2}, Equals, complexStruct{3, 4})
... obtained check_test.complexStruct = check_test.complexStruct{r:1, i:2}
... expected check_test.complexStruct = check_test.complexStruct{r:3, i:4}
... Difference:
...     r: 1 != 3
...     i: 2 != 4


`)
}
