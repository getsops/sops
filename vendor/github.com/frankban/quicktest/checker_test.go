// Licensed under the MIT license, see LICENCE file for details.

package quicktest_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	qt "github.com/frankban/quicktest"
)

var (
	goTime = time.Date(2012, 3, 28, 0, 0, 0, 0, time.UTC)
	chInt  = func() chan int {
		ch := make(chan int, 4)
		ch <- 42
		ch <- 47
		return ch
	}()
	sameInts = cmpopts.SortSlices(func(x, y int) bool {
		return x < y
	})
	cmpEqualsGot = struct {
		Strings []interface{}
		Ints    []int
	}{
		Strings: []interface{}{"who", "dalek"},
		Ints:    []int{42, 47},
	}
	cmpEqualsWant = struct {
		Strings []interface{}
		Ints    []int
	}{
		Strings: []interface{}{"who", "dalek"},
		Ints:    []int{42},
	}
)

var checkerTests = []struct {
	about                 string
	checker               qt.Checker
	got                   interface{}
	args                  []interface{}
	verbose               bool
	expectedCheckFailure  string
	expectedNegateFailure string
}{{
	about:   "Equals: same values",
	checker: qt.Equals,
	got:     42,
	args:    []interface{}{42},
	expectedNegateFailure: `
error:
  unexpected success
got:
  int(42)
want:
  <same as "got">
`,
}, {
	about:   "Equals: different values",
	checker: qt.Equals,
	got:     "42",
	args:    []interface{}{"47"},
	expectedCheckFailure: `
error:
  values are not equal
got:
  "42"
want:
  "47"
`,
}, {
	about:   "Equals: different strings with quotes",
	checker: qt.Equals,
	got:     `string "foo"`,
	args:    []interface{}{`string "bar"`},
	expectedCheckFailure: strings.Replace(`
error:
  values are not equal
got:
  ~string "foo"~
want:
  ~string "bar"~
`, "~", "`", -1),
}, {
	about:   "Equals: different types",
	checker: qt.Equals,
	got:     42,
	args:    []interface{}{"42"},
	expectedCheckFailure: `
error:
  values are not equal
got:
  int(42)
want:
  "42"
`,
}, {
	about:   "Equals: nil and nil",
	checker: qt.Equals,
	got:     nil,
	args:    []interface{}{nil},
	expectedNegateFailure: `
error:
  unexpected success
got:
  nil
want:
  <same as "got">
`,
}, {
	about:   "Equals: error is not nil",
	checker: qt.Equals,
	got:     errBadWolf,
	args:    []interface{}{nil},
	expectedCheckFailure: `
error:
  got non-nil error
got:
  bad wolf
    file:line
want:
  nil
`,
}, {
	about:   "Equals: error is not nil: not formatted",
	checker: qt.Equals,
	got: &errTest{
		msg: "bad wolf",
	},
	args: []interface{}{nil},
	expectedCheckFailure: `
error:
  got non-nil error
got:
  e"bad wolf"
want:
  nil
`,
}, {
	about:   "Equals: error is not nil: not formatted and with quotes",
	checker: qt.Equals,
	got: &errTest{
		msg: `failure: "bad wolf"`,
	},
	args: []interface{}{nil},
	expectedCheckFailure: strings.Replace(`
error:
  got non-nil error
got:
  e~failure: "bad wolf"~
want:
  nil
`, "~", "`", -1),
}, {
	about:   "Equals: nil struct",
	checker: qt.Equals,
	got:     (*struct{})(nil),
	args:    []interface{}{nil},
	expectedCheckFailure: `
error:
  values are not equal
got:
  (*struct {})(nil)
want:
  nil
`,
}, {
	about:   "Equals: uncomparable types",
	checker: qt.Equals,
	got: struct {
		Ints []int
	}{
		Ints: []int{42, 47},
	},
	args: []interface{}{struct {
		Ints []int
	}{
		Ints: []int{42, 47},
	}},
	expectedCheckFailure: `
error:
  runtime error: comparing uncomparable type struct { Ints []int }
got:
  struct { Ints []int }{
      Ints: {42, 47},
  }
want:
  <same as "got">
`,
}, {
	about:   "Equals: not enough arguments",
	checker: qt.Equals,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want
`,
}, {
	about:   "Equals: too many arguments",
	checker: qt.Equals,
	args:    []interface{}{nil, 47},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      nil,
      int(47),
  }
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      nil,
      int(47),
  }
want args:
  want
`,
}, {
	about:   "CmpEquals: same values",
	checker: qt.CmpEquals(),
	got:     cmpEqualsGot,
	args:    []interface{}{cmpEqualsGot},
	expectedNegateFailure: `
error:
  unexpected success
got:
  struct { Strings []interface {}; Ints []int }{
      Strings: {
          "who",
          "dalek",
      },
      Ints: {42, 47},
  }
want:
  <same as "got">
`,
}, {
	about:   "CmpEquals: different values",
	checker: qt.CmpEquals(),
	got:     cmpEqualsGot,
	args:    []interface{}{cmpEqualsWant},
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
`, diff(cmpEqualsGot, cmpEqualsWant)),
}, {
	about:   "CmpEquals: different values: verbose",
	checker: qt.CmpEquals(),
	got:     cmpEqualsGot,
	args:    []interface{}{cmpEqualsWant},
	verbose: true,
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
got:
  struct { Strings []interface {}; Ints []int }{
      Strings: {
          "who",
          "dalek",
      },
      Ints: {42, 47},
  }
want:
  struct { Strings []interface {}; Ints []int }{
      Strings: {
          "who",
          "dalek",
      },
      Ints: {42},
  }
`, diff(cmpEqualsGot, cmpEqualsWant)),
}, {
	about:   "CmpEquals: same values with options",
	checker: qt.CmpEquals(sameInts),
	got:     []int{1, 2, 3},
	args: []interface{}{
		[]int{3, 2, 1},
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  []int{1, 2, 3}
want:
  []int{3, 2, 1}
`,
}, {
	about:   "CmpEquals: different values with options",
	checker: qt.CmpEquals(sameInts),
	got:     []int{1, 2, 4},
	args: []interface{}{
		[]int{3, 2, 1},
	},
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
`, diff([]int{1, 2, 4}, []int{3, 2, 1}, sameInts)),
}, {
	about:   "CmpEquals: different values with options: verbose",
	checker: qt.CmpEquals(sameInts),
	got:     []int{1, 2, 4},
	args: []interface{}{
		[]int{3, 2, 1},
	},
	verbose: true,
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
got:
  []int{1, 2, 4}
want:
  []int{3, 2, 1}
`, diff([]int{1, 2, 4}, []int{3, 2, 1}, sameInts)),
}, {
	about:   "CmpEquals: structs with unexported fields not allowed",
	checker: qt.CmpEquals(),
	got: struct{ answer int }{
		answer: 42,
	},
	args: []interface{}{
		struct{ answer int }{
			answer: 42,
		},
	},
	expectedCheckFailure: `
error:
  cannot handle unexported field: root.answer
  consider using a custom Comparer; if you control the implementation of type, you can also consider AllowUnexported or cmpopts.IgnoreUnexported
got:
  struct { answer int }{answer:42}
want:
  <same as "got">
`,
}, {
	about:   "CmpEquals: structs with unexported fields ignored",
	checker: qt.CmpEquals(cmpopts.IgnoreUnexported(struct{ answer int }{})),
	got: struct{ answer int }{
		answer: 42,
	},
	args: []interface{}{
		struct{ answer int }{
			answer: 42,
		},
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  struct { answer int }{answer:42}
want:
  <same as "got">
`,
}, {
	about:   "CmpEquals: same times",
	checker: qt.CmpEquals(),
	got:     goTime,
	args: []interface{}{
		goTime,
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  s"2012-03-28 00:00:00 +0000 UTC"
want:
  <same as "got">
`,
}, {
	about:   "CmpEquals: different times: verbose",
	checker: qt.CmpEquals(),
	got:     goTime.Add(24 * time.Hour),
	args: []interface{}{
		goTime,
	},
	verbose: true,
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
got:
  s"2012-03-29 00:00:00 +0000 UTC"
want:
  s"2012-03-28 00:00:00 +0000 UTC"
`, diff(goTime.Add(24*time.Hour), goTime)),
}, {
	about:   "CmpEquals: not enough arguments",
	checker: qt.CmpEquals(),
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want
`,
}, {
	about:   "CmpEquals: too many arguments",
	checker: qt.CmpEquals(),
	got:     []int{42},
	args:    []interface{}{[]int{42}, "bad wolf"},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      []int{42},
      "bad wolf",
  }
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      []int{42},
      "bad wolf",
  }
want args:
  want
`,
}, {
	about:   "DeepEquals: different values",
	checker: qt.DeepEquals,
	got:     cmpEqualsGot,
	args:    []interface{}{cmpEqualsWant},
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
`, diff(cmpEqualsGot, cmpEqualsWant)),
}, {
	about:   "DeepEquals: different values: verbose",
	checker: qt.DeepEquals,
	got:     cmpEqualsGot,
	args:    []interface{}{cmpEqualsWant},
	verbose: true,
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
got:
  struct { Strings []interface {}; Ints []int }{
      Strings: {
          "who",
          "dalek",
      },
      Ints: {42, 47},
  }
want:
  struct { Strings []interface {}; Ints []int }{
      Strings: {
          "who",
          "dalek",
      },
      Ints: {42},
  }
`, diff(cmpEqualsGot, cmpEqualsWant)),
}, {
	about:   "ContentEquals: same values",
	checker: qt.ContentEquals,
	got:     []string{"these", "are", "the", "voyages"},
	args: []interface{}{
		[]string{"these", "are", "the", "voyages"},
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  []string{"these", "are", "the", "voyages"}
want:
  <same as "got">
`,
}, {
	about:   "ContentEquals: same contents",
	checker: qt.ContentEquals,
	got:     []int{1, 2, 3},
	args: []interface{}{
		[]int{3, 2, 1},
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  []int{1, 2, 3}
want:
  []int{3, 2, 1}
`,
}, {
	about:   "ContentEquals: same contents on complex slice",
	checker: qt.ContentEquals,
	got: []struct {
		Strings []interface{}
		Ints    []int
	}{cmpEqualsGot, cmpEqualsGot, cmpEqualsWant},
	args: []interface{}{
		[]struct {
			Strings []interface{}
			Ints    []int
		}{cmpEqualsWant, cmpEqualsGot, cmpEqualsGot},
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  []struct { Strings []interface {}; Ints []int }{
      {
          Strings: {
              "who",
              "dalek",
          },
          Ints: {42, 47},
      },
      {
          Strings: {
              "who",
              "dalek",
          },
          Ints: {42, 47},
      },
      {
          Strings: {
              "who",
              "dalek",
          },
          Ints: {42},
      },
  }
want:
  []struct { Strings []interface {}; Ints []int }{
      {
          Strings: {
              "who",
              "dalek",
          },
          Ints: {42},
      },
      {
          Strings: {
              "who",
              "dalek",
          },
          Ints: {42, 47},
      },
      {
          Strings: {
              "who",
              "dalek",
          },
          Ints: {42, 47},
      },
  }
`,
}, {
	about:   "ContentEquals: same contents on a nested slice",
	checker: qt.ContentEquals,
	got: struct {
		Nums []int
	}{
		Nums: []int{1, 2, 3, 4},
	},
	args: []interface{}{
		struct {
			Nums []int
		}{
			Nums: []int{4, 3, 2, 1},
		},
	},
	expectedNegateFailure: `
error:
  unexpected success
got:
  struct { Nums []int }{
      Nums: {1, 2, 3, 4},
  }
want:
  struct { Nums []int }{
      Nums: {4, 3, 2, 1},
  }
`,
}, {
	about:   "ContentEquals: slices of different type",
	checker: qt.ContentEquals,
	got:     []string{"bad", "wolf"},
	args: []interface{}{
		[]interface{}{"bad", "wolf"},
	},
	expectedCheckFailure: fmt.Sprintf(`
error:
  values are not deep equal
diff (-got +want):
%s
`, diff([]string{"bad", "wolf"}, []interface{}{"bad", "wolf"})),
}, {
	about:   "ContentEquals: not enough arguments",
	checker: qt.ContentEquals,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want
`,
}, {
	about:   "ContentEquals: too many arguments",
	checker: qt.ContentEquals,
	args:    []interface{}{nil, nil},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      nil,
      nil,
  }
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      nil,
      nil,
  }
want args:
  want
`,
}, {
	about:   "Matches: perfect match",
	checker: qt.Matches,
	got:     "exterminate",
	args:    []interface{}{"exterminate"},
	expectedNegateFailure: `
error:
  unexpected success
got value:
  "exterminate"
regexp:
  <same as "got value">
`,
}, {
	about:   "Matches: match",
	checker: qt.Matches,
	got:     "these are the voyages",
	args:    []interface{}{"these are the .*"},
	expectedNegateFailure: `
error:
  unexpected success
got value:
  "these are the voyages"
regexp:
  "these are the .*"
`,
}, {
	about:   "Matches: match with stringer",
	checker: qt.Matches,
	got:     bytes.NewBufferString("resistance is futile"),
	args:    []interface{}{"resistance is (futile|useful)"},
	expectedNegateFailure: `
error:
  unexpected success
got value:
  s"resistance is futile"
regexp:
  "resistance is (futile|useful)"
`,
}, {
	about:   "Matches: mismatch",
	checker: qt.Matches,
	got:     "voyages",
	args:    []interface{}{"these are the voyages"},
	expectedCheckFailure: `
error:
  value does not match regexp
got value:
  "voyages"
regexp:
  "these are the voyages"
`,
}, {
	about:   "Matches: mismatch with stringer",
	checker: qt.Matches,
	got:     bytes.NewBufferString("voyages"),
	args:    []interface{}{"these are the voyages"},
	expectedCheckFailure: `
error:
  value.String() does not match regexp
got value:
  s"voyages"
regexp:
  "these are the voyages"
`,
}, {
	about:   "Matches: empty pattern",
	checker: qt.Matches,
	got:     "these are the voyages",
	args:    []interface{}{""},
	expectedCheckFailure: `
error:
  value does not match regexp
got value:
  "these are the voyages"
regexp:
  ""
`,
}, {
	about:   "Matches: complex pattern",
	checker: qt.Matches,
	got:     "end of the universe",
	args:    []interface{}{"bad wolf|end of the .*"},
	expectedNegateFailure: `
error:
  unexpected success
got value:
  "end of the universe"
regexp:
  "bad wolf|end of the .*"
`,
}, {
	about:   "Matches: invalid pattern",
	checker: qt.Matches,
	got:     "voyages",
	args:    []interface{}{"("},
	expectedCheckFailure: `
error:
  bad check: cannot compile regexp: error parsing regexp: missing closing ): ` + "`^(()$`" + `
regexp:
  "("
`,
	expectedNegateFailure: `
error:
  bad check: cannot compile regexp: error parsing regexp: missing closing ): ` + "`^(()$`" + `
regexp:
  "("
`,
}, {
	about:   "Matches: pattern not a string",
	checker: qt.Matches,
	got:     "",
	args:    []interface{}{[]int{42}},
	expectedCheckFailure: `
error:
  bad check: regexp is not a string
regexp:
  []int{42}
`,
	expectedNegateFailure: `
error:
  bad check: regexp is not a string
regexp:
  []int{42}
`,
}, {
	about:   "Matches: not a string or as stringer",
	checker: qt.Matches,
	got:     42,
	args:    []interface{}{".*"},
	expectedCheckFailure: `
error:
  bad check: value is not a string or a fmt.Stringer
value:
  int(42)
`,
	expectedNegateFailure: `
error:
  bad check: value is not a string or a fmt.Stringer
value:
  int(42)
`,
}, {
	about:   "Matches: not enough arguments",
	checker: qt.Matches,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
}, {
	about:   "Matches: too many arguments",
	checker: qt.Matches,
	got:     "these are the voyages",
	args:    []interface{}{"these are the .*", nil},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      "these are the .*",
      nil,
  }
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      "these are the .*",
      nil,
  }
want args:
  regexp
`,
}, {
	about:   "ErrorMatches: perfect match",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{"bad wolf"},
	expectedNegateFailure: `
error:
  unexpected success
got error:
  bad wolf
    file:line
regexp:
  "bad wolf"
`,
}, {
	about:   "ErrorMatches: match",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{"bad .*"},
	expectedNegateFailure: `
error:
  unexpected success
got error:
  bad wolf
    file:line
regexp:
  "bad .*"
`,
}, {
	about:   "ErrorMatches: mismatch",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{"exterminate"},
	expectedCheckFailure: `
error:
  error does not match regexp
got error:
  bad wolf
    file:line
regexp:
  "exterminate"
`,
}, {
	about:   "ErrorMatches: empty pattern",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{""},
	expectedCheckFailure: `
error:
  error does not match regexp
got error:
  bad wolf
    file:line
regexp:
  ""
`,
}, {
	about:   "ErrorMatches: complex pattern",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{"bad wolf|end of the universe"},
	expectedNegateFailure: `
error:
  unexpected success
got error:
  bad wolf
    file:line
regexp:
  "bad wolf|end of the universe"
`,
}, {
	about:   "ErrorMatches: invalid pattern",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{"("},
	expectedCheckFailure: `
error:
  bad check: cannot compile regexp: error parsing regexp: missing closing ): ` + "`^(()$`" + `
regexp:
  "("
`,
	expectedNegateFailure: `
error:
  bad check: cannot compile regexp: error parsing regexp: missing closing ): ` + "`^(()$`" + `
regexp:
  "("
`,
}, {
	about:   "ErrorMatches: pattern not a string",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{[]int{42}},
	expectedCheckFailure: `
error:
  bad check: regexp is not a string
regexp:
  []int{42}
`,
	expectedNegateFailure: `
error:
  bad check: regexp is not a string
regexp:
  []int{42}
`,
}, {
	about:   "ErrorMatches: not an error",
	checker: qt.ErrorMatches,
	got:     42,
	args:    []interface{}{".*"},
	expectedCheckFailure: `
error:
  bad check: first argument is not an error
got:
  int(42)
`,
	expectedNegateFailure: `
error:
  bad check: first argument is not an error
got:
  int(42)
`,
}, {
	about:   "ErrorMatches: nil error",
	checker: qt.ErrorMatches,
	got:     nil,
	args:    []interface{}{"some pattern"},
	expectedCheckFailure: `
error:
  got nil error but want non-nil
got error:
  nil
regexp:
  "some pattern"
`,
}, {
	about:   "ErrorMatches: not enough arguments",
	checker: qt.ErrorMatches,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
}, {
	about:   "ErrorMatches: too many arguments",
	checker: qt.ErrorMatches,
	got:     errBadWolf,
	args:    []interface{}{"bad wolf", []string{"bad", "wolf"}},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      "bad wolf",
      []string{"bad", "wolf"},
  }
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      "bad wolf",
      []string{"bad", "wolf"},
  }
want args:
  regexp
`,
}, {
	about:   "PanicMatches: perfect match",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{"error: bad wolf"},
	expectedNegateFailure: `
error:
  unexpected success
panic value:
  "error: bad wolf"
function:
  func() {...}
regexp:
  <same as "panic value">
`,
}, {
	about:   "PanicMatches: match",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{"error: .*"},
	expectedNegateFailure: `
error:
  unexpected success
panic value:
  "error: bad wolf"
function:
  func() {...}
regexp:
  "error: .*"
`,
}, {
	about:   "PanicMatches: mismatch",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{"error: exterminate"},
	expectedCheckFailure: `
error:
  panic value does not match regexp
panic value:
  "error: bad wolf"
function:
  func() {...}
regexp:
  "error: exterminate"
`,
}, {
	about:   "PanicMatches: empty pattern",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{""},
	expectedCheckFailure: `
error:
  panic value does not match regexp
panic value:
  "error: bad wolf"
function:
  func() {...}
regexp:
  ""
`,
}, {
	about:   "PanicMatches: complex pattern",
	checker: qt.PanicMatches,
	got:     func() { panic("bad wolf") },
	args:    []interface{}{"bad wolf|end of the universe"},
	expectedNegateFailure: `
error:
  unexpected success
panic value:
  "bad wolf"
function:
  func() {...}
regexp:
  "bad wolf|end of the universe"
`,
}, {
	about:   "PanicMatches: invalid pattern",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{"("},
	expectedCheckFailure: `
error:
  bad check: cannot compile regexp: error parsing regexp: missing closing ): ` + "`^(()$`" + `
panic value:
  "error: bad wolf"
regexp:
  "("
`,
	expectedNegateFailure: `
error:
  bad check: cannot compile regexp: error parsing regexp: missing closing ): ` + "`^(()$`" + `
panic value:
  "error: bad wolf"
regexp:
  "("
`,
}, {
	about:   "PanicMatches: pattern not a string",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{nil},
	expectedCheckFailure: `
error:
  bad check: regexp is not a string
panic value:
  "error: bad wolf"
regexp:
  nil
`,
	expectedNegateFailure: `
error:
  bad check: regexp is not a string
panic value:
  "error: bad wolf"
regexp:
  nil
`,
}, {
	about:   "PanicMatches: not a function",
	checker: qt.PanicMatches,
	got:     map[string]int{"answer": 42},
	args:    []interface{}{".*"},
	expectedCheckFailure: `
error:
  bad check: first argument is not a function
got:
  map[string]int{"answer":42}
`,
	expectedNegateFailure: `
error:
  bad check: first argument is not a function
got:
  map[string]int{"answer":42}
`,
}, {
	about:   "PanicMatches: not a proper function",
	checker: qt.PanicMatches,
	got:     func(int) { panic("error: bad wolf") },
	args:    []interface{}{".*"},
	expectedCheckFailure: `
error:
  bad check: cannot use a function receiving arguments
function:
  func(int) {...}
`,
	expectedNegateFailure: `
error:
  bad check: cannot use a function receiving arguments
function:
  func(int) {...}
`,
}, {
	about:   "PanicMatches: function returning something",
	checker: qt.PanicMatches,
	got:     func() error { panic("error: bad wolf") },
	args:    []interface{}{".*"},
	expectedNegateFailure: `
error:
  unexpected success
panic value:
  "error: bad wolf"
function:
  func() error {...}
regexp:
  ".*"
`,
}, {
	about:   "PanicMatches: no panic",
	checker: qt.PanicMatches,
	got:     func() {},
	args:    []interface{}{".*"},
	expectedCheckFailure: `
error:
  function did not panic
function:
  func() {...}
regexp:
  ".*"
`,
}, {
	about:   "PanicMatches: not enough arguments",
	checker: qt.PanicMatches,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
}, {
	about:   "PanicMatches: too many arguments",
	checker: qt.PanicMatches,
	got:     func() { panic("error: bad wolf") },
	args:    []interface{}{"error: bad wolf", 42},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      "error: bad wolf",
      int(42),
  }
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      "error: bad wolf",
      int(42),
  }
want args:
  regexp
`,
}, {
	about:   "IsNil: nil",
	checker: qt.IsNil,
	got:     nil,
	expectedNegateFailure: `
error:
  unexpected success
got:
  nil
`,
}, {
	about:   "IsNil: nil struct",
	checker: qt.IsNil,
	got:     (*struct{})(nil),
	expectedNegateFailure: `
error:
  unexpected success
got:
  (*struct {})(nil)
`,
}, {
	about:   "IsNil: nil func",
	checker: qt.IsNil,
	got:     (func())(nil),
	expectedNegateFailure: `
error:
  unexpected success
got:
  func() {...}
`,
}, {
	about:   "IsNil: nil map",
	checker: qt.IsNil,
	got:     (map[string]string)(nil),
	expectedNegateFailure: `
error:
  unexpected success
got:
  map[string]string{}
`,
}, {
	about:   "IsNil: nil slice",
	checker: qt.IsNil,
	got:     ([]int)(nil),
	expectedNegateFailure: `
error:
  unexpected success
got:
  []int(nil)
`,
}, {
	about:   "IsNil: not nil",
	checker: qt.IsNil,
	got:     42,
	expectedCheckFailure: `
error:
  42 is not nil
got:
  int(42)
`,
}, {
	about:   "IsNil: too many arguments",
	checker: qt.IsNil,
	args:    []interface{}{"not nil"},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 1, want 0
got args:
  []interface {}{
      "not nil",
  }
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 1, want 0
got args:
  []interface {}{
      "not nil",
  }
`,
}, {
	about:   "HasLen: arrays with the same length",
	checker: qt.HasLen,
	got:     [4]string{"these", "are", "the", "voyages"},
	args:    []interface{}{4},
	expectedNegateFailure: `
error:
  unexpected success
len(got):
  int(4)
got:
  [4]string{"these", "are", "the", "voyages"}
want length:
  <same as "len(got)">
`,
}, {
	about:   "HasLen: channels with the same length",
	checker: qt.HasLen,
	got:     chInt,
	args:    []interface{}{2},
	expectedNegateFailure: fmt.Sprintf(`
error:
  unexpected success
len(got):
  int(2)
got:
  (chan int)(%v)
want length:
  <same as "len(got)">
`, chInt),
}, {
	about:   "HasLen: maps with the same length",
	checker: qt.HasLen,
	got:     map[string]bool{"true": true},
	args:    []interface{}{1},
	expectedNegateFailure: `
error:
  unexpected success
len(got):
  int(1)
got:
  map[string]bool{"true":true}
want length:
  <same as "len(got)">
`,
}, {
	about:   "HasLen: slices with the same length",
	checker: qt.HasLen,
	got:     []int{},
	args:    []interface{}{0},
	expectedNegateFailure: `
error:
  unexpected success
len(got):
  int(0)
got:
  []int{}
want length:
  <same as "len(got)">
`,
}, {
	about:   "HasLen: strings with the same length",
	checker: qt.HasLen,
	got:     "these are the voyages",
	args:    []interface{}{21},
	expectedNegateFailure: `
error:
  unexpected success
len(got):
  int(21)
got:
  "these are the voyages"
want length:
  <same as "len(got)">
`,
}, {
	about:   "HasLen: arrays with different lengths",
	checker: qt.HasLen,
	got:     [4]string{"these", "are", "the", "voyages"},
	args:    []interface{}{0},
	expectedCheckFailure: `
error:
  unexpected length
len(got):
  int(4)
got:
  [4]string{"these", "are", "the", "voyages"}
want length:
  int(0)
`,
}, {
	about:   "HasLen: channels with different lengths",
	checker: qt.HasLen,
	got:     chInt,
	args:    []interface{}{4},
	expectedCheckFailure: fmt.Sprintf(`
error:
  unexpected length
len(got):
  int(2)
got:
  (chan int)(%v)
want length:
  int(4)
`, chInt),
}, {
	about:   "HasLen: maps with different lengths",
	checker: qt.HasLen,
	got:     map[string]bool{"true": true},
	args:    []interface{}{42},
	expectedCheckFailure: `
error:
  unexpected length
len(got):
  int(1)
got:
  map[string]bool{"true":true}
want length:
  int(42)
`,
}, {
	about:   "HasLen: slices with different lengths",
	checker: qt.HasLen,
	got:     []int{42, 47},
	args:    []interface{}{1},
	expectedCheckFailure: `
error:
  unexpected length
len(got):
  int(2)
got:
  []int{42, 47}
want length:
  int(1)
`,
}, {
	about:   "HasLen: strings with different lengths",
	checker: qt.HasLen,
	got:     "these are the voyages",
	args:    []interface{}{42},
	expectedCheckFailure: `
error:
  unexpected length
len(got):
  int(21)
got:
  "these are the voyages"
want length:
  int(42)
`,
}, {
	about:   "HasLen: value without a length",
	checker: qt.HasLen,
	got:     42,
	args:    []interface{}{42},
	expectedCheckFailure: `
error:
  bad check: first argument has no length
got:
  int(42)
`,
	expectedNegateFailure: `
error:
  bad check: first argument has no length
got:
  int(42)
`,
}, {
	about:   "HasLen: expected value not a number",
	checker: qt.HasLen,
	got:     "these are the voyages",
	args:    []interface{}{"bad wolf"},
	expectedCheckFailure: `
error:
  bad check: length is not an int
length:
  "bad wolf"
`,
	expectedNegateFailure: `
error:
  bad check: length is not an int
length:
  "bad wolf"
`,
}, {
	about:   "HasLen: not enough arguments",
	checker: qt.HasLen,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want length
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  want length
`,
}, {
	about:   "HasLen: too many arguments",
	checker: qt.HasLen,
	got:     []int{42},
	args:    []interface{}{42, 47},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      int(42),
      int(47),
  }
want args:
  want length
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      int(42),
      int(47),
  }
want args:
  want length
`,
}, {
	about:   "Satisfies: success with an error",
	checker: qt.Satisfies,
	got:     qt.BadCheckf("bad wolf"),
	args:    []interface{}{qt.IsBadCheck},
	expectedNegateFailure: `
error:
  unexpected success
arg:
  e"bad check: bad wolf"
predicate function:
  func(error) bool {...}
`,
}, {
	about:   "Satisfies: success with an int",
	checker: qt.Satisfies,
	got:     42,
	args: []interface{}{
		func(v int) bool { return v == 42 },
	},
	expectedNegateFailure: `
error:
  unexpected success
arg:
  int(42)
predicate function:
  func(int) bool {...}
`,
}, {
	about:   "Satisfies: success with nil",
	checker: qt.Satisfies,
	got:     nil,
	args: []interface{}{
		func(v []int) bool { return true },
	},
	expectedNegateFailure: `
error:
  unexpected success
arg:
  nil
predicate function:
  func([]int) bool {...}
`,
}, {
	about:   "Satisfies: failure with an error",
	checker: qt.Satisfies,
	got:     nil,
	args:    []interface{}{qt.IsBadCheck},
	expectedCheckFailure: `
error:
  value does not satisfy predicate function
arg:
  nil
predicate function:
  func(error) bool {...}
`,
}, {
	about:   "Satisfies: failure with a string",
	checker: qt.Satisfies,
	got:     "bad wolf",
	args: []interface{}{
		func(string) bool { return false },
	},
	expectedCheckFailure: `
error:
  value does not satisfy predicate function
arg:
  "bad wolf"
predicate function:
  func(string) bool {...}
`,
}, {
	about:   "Satisfies: not a function",
	checker: qt.Satisfies,
	got:     42,
	args:    []interface{}{42},
	expectedCheckFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  int(42)
`,
	expectedNegateFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  int(42)
`,
}, {
	about:   "Satisfies: function accepting no arguments",
	checker: qt.Satisfies,
	got:     42,
	args: []interface{}{
		func() bool { return true },
	},
	expectedCheckFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func() bool {...}
`,
	expectedNegateFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func() bool {...}
`,
}, {
	about:   "Satisfies: function accepting too many arguments",
	checker: qt.Satisfies,
	got:     42,
	args: []interface{}{
		func(int, string) bool { return false },
	},
	expectedCheckFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(int, string) bool {...}
`,
	expectedNegateFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(int, string) bool {...}
`,
}, {
	about:   "Satisfies: function returning no arguments",
	checker: qt.Satisfies,
	got:     42,
	args: []interface{}{
		func(error) {},
	},
	expectedCheckFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(error) {...}
`,
	expectedNegateFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(error) {...}
`,
}, {
	about:   "Satisfies: function returning too many argments",
	checker: qt.Satisfies,
	got:     42,
	args: []interface{}{
		func(int) (bool, error) { return true, nil },
	},
	expectedCheckFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(int) (bool, error) {...}
`,
	expectedNegateFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(int) (bool, error) {...}
`,
}, {
	about:   "Satisfies: function not returning a bool",
	checker: qt.Satisfies,
	got:     42,
	args: []interface{}{
		func(int) error { return nil },
	},
	expectedCheckFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(int) error {...}
`,
	expectedNegateFailure: `
error:
  bad check: predicate function is not a func(T) bool
predicate function:
  func(int) error {...}
`,
}, {
	about:   "Satisfies: type mismatch",
	checker: qt.Satisfies,
	got:     42,
	args:    []interface{}{qt.IsBadCheck},
	expectedCheckFailure: `
error:
  bad check: cannot use value of type int as type error in argument to predicate function
arg:
  int(42)
predicate function:
  func(error) bool {...}
`,
	expectedNegateFailure: `
error:
  bad check: cannot use value of type int as type error in argument to predicate function
arg:
  int(42)
predicate function:
  func(error) bool {...}
`,
}, {
	about:   "Satisfies: nil value that cannot be nil",
	checker: qt.Satisfies,
	got:     nil,
	args: []interface{}{
		func(string) bool { return true },
	},
	expectedCheckFailure: `
error:
  bad check: cannot use nil as type string in argument to predicate function
predicate function:
  func(string) bool {...}
`,
	expectedNegateFailure: `
error:
  bad check: cannot use nil as type string in argument to predicate function
predicate function:
  func(string) bool {...}
`,
}, {
	about:   "Satisfies: not enough arguments",
	checker: qt.Satisfies,
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  predicate function
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  predicate function
`,
}, {
	about:   "Satisfies: too many arguments",
	checker: qt.Satisfies,
	got:     42,
	args:    []interface{}{func() bool { return true }, 1, 2},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 3, want 1
got args:
  []interface {}{
      func() bool {...},
      int(1),
      int(2),
  }
want args:
  predicate function
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 3, want 1
got args:
  []interface {}{
      func() bool {...},
      int(1),
      int(2),
  }
want args:
  predicate function
`,
}, {
	about:   "Not: success",
	checker: qt.Not(qt.IsNil),
	got:     42,
	expectedNegateFailure: `
error:
  42 is not nil
got:
  int(42)
`,
}, {
	about:   "Not: failure",
	checker: qt.Not(qt.IsNil),
	got:     nil,
	expectedCheckFailure: `
error:
  unexpected success
got:
  nil
`,
}, {
	about:   "Not: not enough arguments",
	checker: qt.Not(qt.PanicMatches),
	expectedCheckFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
	expectedNegateFailure: `
error:
  bad check: not enough arguments provided to checker: got 0, want 1
want args:
  regexp
`,
}, {
	about:   "Not: too many arguments",
	checker: qt.Not(qt.Equals),
	args:    []interface{}{42, nil},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      int(42),
      nil,
  }
want args:
  want
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 2, want 1
got args:
  []interface {}{
      int(42),
      nil,
  }
want args:
  want
`,
}, {
	about:   "Contains with string",
	checker: qt.Contains,
	got:     "hello, world",
	args:    []interface{}{"world"},
	expectedNegateFailure: `
error:
  unexpected success
got:
  "hello, world"
want:
  "world"
`,
}, {
	about:   "Contains with string no match",
	checker: qt.Contains,
	got:     "hello, world",
	args:    []interface{}{"worlds"},
	expectedCheckFailure: `
error:
  no substring match found
got:
  "hello, world"
want:
  "worlds"
`,
}, {
	about:   "Contains with slice",
	checker: qt.Contains,
	got:     []string{"a", "b", "c"},
	args:    []interface{}{"a"},
	expectedNegateFailure: `
error:
  unexpected success
got:
  []string{"a", "b", "c"}
want:
  "a"
`,
}, {
	about:   "Contains with map",
	checker: qt.Contains,
	// Note: we can't use more than one element here because
	// pretty.Print output is non-deterministic.
	// https://github.com/kr/pretty/issues/47
	got:  map[string]string{"a": "d"},
	args: []interface{}{"d"},
	expectedNegateFailure: `
error:
  unexpected success
got:
  map[string]string{"a":"d"}
want:
  "d"
`,
}, {
	about:   "Contains with non-string",
	checker: qt.Contains,
	got:     "aa",
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: strings can only contain strings, not int
`,
	expectedNegateFailure: `
error:
  bad check: strings can only contain strings, not int
`,
}, {
	about:   "All slice equals",
	checker: qt.All(qt.Equals),
	got:     []string{"a", "a"},
	args:    []interface{}{"a"},
	expectedNegateFailure: `
error:
  unexpected success
container:
  []string{"a", "a"}
want:
  "a"
`,
}, {
	about:   "All slice match",
	checker: qt.All(qt.Matches),
	got:     []string{"red", "blue", "green"},
	args:    []interface{}{".*e.*"},
	expectedNegateFailure: `
error:
  unexpected success
container:
  []string{"red", "blue", "green"}
regexp:
  ".*e.*"
`,
}, {
	about:   "All nested match",
	checker: qt.All(qt.All(qt.Matches)),
	got:     [][]string{{"hello", "goodbye"}, {"red", "blue"}, {}},
	args:    []interface{}{".*e.*"},
	expectedNegateFailure: `
error:
  unexpected success
container:
  [][]string{
      {"hello", "goodbye"},
      {"red", "blue"},
      {},
  }
regexp:
  ".*e.*"
`,
}, {
	about:   "All nested mismatch",
	checker: qt.All(qt.All(qt.Matches)),
	got:     [][]string{{"hello", "goodbye"}, {"black", "blue"}, {}},
	args:    []interface{}{".*e.*"},
	expectedCheckFailure: `
error:
  mismatch at index 1
error:
  mismatch at index 0
error:
  value does not match regexp
first mismatched element:
  "black"
`,
}, {
	about:   "All slice mismatch",
	checker: qt.All(qt.Matches),
	got:     []string{"red", "black"},
	args:    []interface{}{".*e.*"},
	expectedCheckFailure: `
error:
  mismatch at index 1
error:
  value does not match regexp
first mismatched element:
  "black"
`,
}, {
	about:   "All slice mismatch with DeepEqual",
	checker: qt.All(qt.DeepEquals),
	got:     [][]string{{"a", "b"}, {"a", "c"}},
	args:    []interface{}{[]string{"a", "b"}},
	expectedCheckFailure: `
error:
  mismatch at index 1
error:
  values are not deep equal
diff (-got +want):
` + diff([]string{"a", "c"}, []string{"a", "b"}) + `
`,
}, {
	about:   "All bad checker args count",
	checker: qt.All(qt.IsNil),
	got:     []int{},
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 1, want 0
got args:
  []interface {}{
      int(5),
  }
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 1, want 0
got args:
  []interface {}{
      int(5),
  }
`,
}, {
	about:   "All bad checker args",
	checker: qt.All(qt.Matches),
	got:     []string{"hello"},
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: at index 0: bad check: regexp is not a string
`,
	expectedNegateFailure: `
error:
  bad check: at index 0: bad check: regexp is not a string
`,
}, {
	about:   "All with non-container",
	checker: qt.All(qt.Equals),
	got:     5,
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: map, slice or array required
`,
	expectedNegateFailure: `
error:
  bad check: map, slice or array required
`,
}, {
	about:   "All mismatch with map",
	checker: qt.All(qt.Matches),
	got:     map[string]string{"a": "red", "b": "black"},
	args:    []interface{}{".*e.*"},
	expectedCheckFailure: `
error:
  mismatch at key "b"
error:
  value does not match regexp
first mismatched element:
  "black"
`,
}, {
	about:   "Any with non-container",
	checker: qt.Any(qt.Equals),
	got:     5,
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: map, slice or array required
`,
	expectedNegateFailure: `
error:
  bad check: map, slice or array required
`,
}, {
	about:   "Any no match",
	checker: qt.Any(qt.Equals),
	got:     []int{},
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  no matching element found
container:
  []int{}
want:
  int(5)
`,
}, {
	about:   "Any bad checker arg count",
	checker: qt.Any(qt.IsNil),
	got:     []int{},
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: too many arguments provided to checker: got 1, want 0
got args:
  []interface {}{
      int(5),
  }
`,
	expectedNegateFailure: `
error:
  bad check: too many arguments provided to checker: got 1, want 0
got args:
  []interface {}{
      int(5),
  }
`,
}, {
	about:   "Any bad checker args",
	checker: qt.Any(qt.Matches),
	got:     []string{"hello"},
	args:    []interface{}{5},
	expectedCheckFailure: `
error:
  bad check: at index 0: bad check: regexp is not a string
`,
	expectedNegateFailure: `
error:
  bad check: at index 0: bad check: regexp is not a string
`,
}}

func TestCheckers(t *testing.T) {
	for _, test := range checkerTests {
		checker := qt.WithVerbosity(test.checker, test.verbose)
		t.Run(test.about, func(t *testing.T) {
			tt := &testingT{}
			c := qt.New(tt)
			ok := c.Check(test.got, checker, test.args...)
			checkResult(t, ok, tt.errorString(), test.expectedCheckFailure)
		})
		t.Run("Not "+test.about, func(t *testing.T) {
			tt := &testingT{}
			c := qt.New(tt)
			ok := c.Check(test.got, qt.Not(checker), test.args...)
			checkResult(t, ok, tt.errorString(), test.expectedNegateFailure)
		})
	}
}

func diff(x, y interface{}, opts ...cmp.Option) string {
	d := cmp.Diff(x, y, opts...)
	return strings.TrimSuffix(qt.Prefixf("  ", "%s", d), "\n")
}
