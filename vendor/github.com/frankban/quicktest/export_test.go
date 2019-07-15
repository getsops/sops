// Licensed under the MIT license, see LICENCE file for details.

package quicktest

var Prefixf = prefixf

// WithVerbosity returns the given checker with a verbosity level of v.
// A copy of the original checker is made if mutating it is required.
func WithVerbosity(c Checker, v bool) Checker {
	if c, ok := c.(*cmpEqualsChecker); ok {
		c := *c
		c.verbose = func() bool {
			return v
		}
		return &c
	}
	return c
}
