// Licensed under the MIT license, see LICENCE file for details.

package quicktest_test

import (
	"runtime"
	"strconv"
	"strings"
	"testing"

	qt "github.com/frankban/quicktest"
)

// The tests in this file rely on their own source code lines.

func TestReportOutput(t *testing.T) {
	tt := &testingT{}
	c := qt.New(tt)
	c.Assert(42, qt.Equals, 47)
	want := `
error:
  values are not equal
got:
  int(42)
want:
  int(47)
stack:
  $file:19
    c.Assert(42, qt.Equals, 47)
`
	assertReport(t, tt, want)
}

func f1(c *qt.C) {
	f2(c)
}

func f2(c *qt.C) {
	c.Assert(42, qt.IsNil) // Real assertion here!
}

func TestIndirectReportOutput(t *testing.T) {
	tt := &testingT{}
	c := qt.New(tt)
	f1(c)
	want := `
error:
  42 is not nil
got:
  int(42)
stack:
  $file:39
    c.Assert(42, qt.IsNil)
  $file:35
    f2(c)
  $file:45
    f1(c)
`
	assertReport(t, tt, want)
}

func TestMultilineReportOutput(t *testing.T) {
	tt := &testingT{}
	c := qt.New(tt)
	c.Assert(
		"this string", // Comment 1.
		qt.Equals,
		"another string",
		qt.Commentf("a comment"), // Comment 2.
	) // Comment 3.
	want := `
error:
  values are not equal
comment:
  a comment
got:
  "this string"
want:
  "another string"
stack:
  $file:$line
    c.Assert(
        "this string", // Comment 1.
        qt.Equals,
        "another string",
        qt.Commentf("a comment"), // Comment 2.
    )
`
	assertReport(t, tt, want)
}

func TestCmpReportOutput(t *testing.T) {
	tt := &testingT{}
	c := qt.New(tt)
	gotExamples := []*reportExample{{
		AnInt:  42,
		ASlice: []string{},
	}, {
		AnInt:  47,
		ASlice: []string{"these", "are", "the", "voyages"},
	}, {
		AnInt: 1,
	}, {
		AnInt: 2,
	}, {
		ASlice: []string{"foo", "bar"},
	}}
	wantExamples := []*reportExample{{
		AnInt: 42,
	}, {
		AnInt:  47,
		ASlice: []string{"these", "are", "the", "voyages"},
	}, {
		AnInt: 2,
	}, {
		AnInt: 1,
	}, {
		ASlice: []string{"foo"},
	}, {}}
	checker := qt.WithVerbosity(qt.DeepEquals, false)
	c.Assert(gotExamples, checker, wantExamples)
	want := `
error:
  values are not deep equal
diff (-got +want):
    []*quicktest_test.reportExample{
            &{
                    AnInt:  42,
  -                 ASlice: []string{},
  +                 ASlice: nil,
            },
            &{AnInt: 47, ASlice: []string{"these", "are", "the", "voyages"}},
  +         &{AnInt: 2},
            &{AnInt: 1},
  -         &{AnInt: 2},
            &{
                    AnInt: 0,
                    ASlice: []string{
                            "foo",
  -                         "bar",
                    },
            },
  +         &{},
    }
stack:
  $file:121
    c.Assert(gotExamples, checker, wantExamples)
`
	assertReport(t, tt, want)
}

func assertReport(t *testing.T, tt *testingT, want string) {
	got := strings.Replace(tt.fatalString(), "\t", "        ", -1)
	// go-cmp can include non-breaking spaces in its output.
	got = strings.Replace(got, "\u00a0", " ", -1)
	// Adjust for file names in different systems.
	_, file, _, ok := runtime.Caller(0)
	assertBool(t, ok, true)
	want = strings.Replace(want, "$file", file, -1)
	// Adjust for line number based on Go < v1.9 reporting the line where the
	// statement ends.
	line := 65
	vers := runtime.Version()
	if vers == "go1.7" || vers == "go1.8" {
		line = 70
	}
	want = strings.Replace(want, "$line", strconv.Itoa(line), 1)
	if got != want {
		t.Fatalf(`failure:
%q
%q
------------------------------ got ------------------------------
%s------------------------------ want -----------------------------
%s-----------------------------------------------------------------`,
			got, want, got, want)
	}
}

type reportExample struct {
	AnInt  int
	ASlice []string
}
