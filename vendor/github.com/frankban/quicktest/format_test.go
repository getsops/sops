// Licensed under the MIT license, see LICENCE file for details.

package quicktest_test

import (
	"bytes"
	"testing"

	qt "github.com/frankban/quicktest"
)

var formatTests = []struct {
	about string
	value interface{}
	want  string
}{{
	about: "error value",
	value: errBadWolf,
	want:  "bad wolf\n  file:line",
}, {
	about: "error value: not formatted",
	value: &errTest{
		msg: "exterminate!",
	},
	want: `e"exterminate!"`,
}, {
	about: "error value: with quotes",
	value: &errTest{
		msg: `cannot open "/no/such/file"`,
	},
	want: "e`cannot open \"/no/such/file\"`",
}, {
	about: "error value: multi-line",
	value: &errTest{
		msg: `err:
"these are the voyages"`,
	},
	want: `e"err:\n\"these are the voyages\""`,
}, {
	about: "error value: with backquotes",
	value: &errTest{
		msg: "cannot `open` \"file\"",
	},
	want: `e"cannot ` + "`open`" + ` \"file\""`,
}, {
	about: "stringer",
	value: bytes.NewBufferString("I am a stringer"),
	want:  `s"I am a stringer"`,
}, {
	about: "stringer: with quotes",
	value: bytes.NewBufferString(`I say "hello"`),
	want:  "s`I say \"hello\"`",
}, {
	about: "string",
	value: "these are the voyages",
	want:  `"these are the voyages"`,
}, {
	about: "string: with quotes",
	value: `here is a quote: "`,
	want:  "`here is a quote: \"`",
}, {
	about: "string: multi-line",
	value: `foo
"bar"
`,
	want: `"foo\n\"bar\"\n"`,
}, {
	about: "string: with backquotes",
	value: `"` + "`",
	want:  `"\"` + "`\"",
}, {
	about: "slice",
	value: []int{1, 2, 3},
	want:  "[]int{1, 2, 3}",
}, {
	about: "time",
	value: goTime,
	want:  `s"2012-03-28 00:00:00 +0000 UTC"`,
}}

func TestFormat(t *testing.T) {
	for _, test := range formatTests {
		t.Run(test.about, func(t *testing.T) {
			got := qt.Format(test.value)
			if got != test.want {
				t.Fatalf("format:\ngot  %q\nwant %q", got, test.want)
			}
		})
	}
}
