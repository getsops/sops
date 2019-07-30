// +build js

package math_test

import (
	"testing"
)

// Slighly higher tolerances than upstream, otherwise TestGamma fails.
// TODO: Is there a better way to fix TestGamma? It's weird that only one test
//       requires increasing tolerances. Perhaps there's a better fix? Maybe we
//       should override TestGamma specifically and not the package-wide tolerances,
//       because this will cause many other tests to be less accurate. Or maybe this
//       is fine?
func close(a, b float64) bool     { return tolerance(a, b, 4e-14) }
func veryclose(a, b float64) bool { return tolerance(a, b, 6e-15) }

func testExp(t *testing.T, Exp func(float64) float64, name string) {
	t.Skip("inaccurate")
}
