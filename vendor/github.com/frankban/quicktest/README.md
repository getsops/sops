[![GoDoc](https://godoc.org/github.com/frankban/quicktest?status.svg)](https://godoc.org/github.com/frankban/quicktest)
[![Build Status](https://travis-ci.org/frankban/quicktest.svg?branch=master)](https://travis-ci.org/frankban/quicktest)

# quicktest

A collection of Go helpers for writing tests.

## Installation

To install the package, run `go get github.com/frankban/quicktest`.

## Usage

Quicktest helpers can be easily integrated inside regular Go tests, for
instance:
```go
import qt "github.com/frankban/quicktest"

func TestFoo(t *testing.T) {
    t.Run("numbers", func(t *testing.T) {
        c := qt.New(t)
        numbers, err := somepackage.Numbers()
        c.Assert(numbers, qt.DeepEquals, []int{42, 47})
        c.Assert(err, qt.ErrorMatches, "bad wolf")
    })
    t.Run("nil", func(t *testing.T) {
        c := qt.New(t)
        got := somepackage.MaybeNil()
        c.Assert(got, qt.IsNil, qt.Commentf("value: %v", somepackage.Value))
    })
}
```
The library provides some base checkers like Equals, DeepEquals, Matches,
ErrorMatches, IsNil and others. More can be added by implementing the Checker
interface.

See the
[go documentation](https://godoc.org/github.com/frankban/quicktest) for this
library.
