// Licensed under the MIT license, see LICENCE file for details.

package quicktest_test

import (
	"testing"

	qt "github.com/frankban/quicktest"
)

func TestCommentf(t *testing.T) {
	c := qt.Commentf("the answer is %d", 42)
	comment := c.String()
	expectedComment := "the answer is 42"
	if comment != expectedComment {
		t.Fatalf("comment error:\ngot  %q\nwant %q", comment, expectedComment)
	}
}

func TestConstantCommentf(t *testing.T) {
	const expectedComment = "bad wolf"
	c := qt.Commentf(expectedComment)
	comment := c.String()
	if comment != expectedComment {
		t.Fatalf("constant comment error:\ngot  %q\nwant %q", comment, expectedComment)
	}
}
