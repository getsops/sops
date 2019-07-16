// Licensed under the MIT license, see LICENCE file for details.

/*
Package quicktest provides a collection of Go helpers for writing tests.

Quicktest helpers can be easily integrated inside regular Go tests, for
instance:

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

The library provides some base checkers like Equals, DeepEquals, Matches,
ErrorMatches, IsNil and others. More can be added by implementing the Checker
interface.
*/
package quicktest
