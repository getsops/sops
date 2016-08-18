//Package performance provides gucumber integration tests support.
package performance

import (
	"github.com/lsegal/gucumber"
)

func init() {
	gucumber.Before("@performance", func() {
	})
}
