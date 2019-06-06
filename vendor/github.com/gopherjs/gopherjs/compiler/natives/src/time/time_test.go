// +build js

package time_test

import (
	"testing"
)

func TestSleep(t *testing.T) {
	t.Skip("time.Now() is not accurate enough for the test")
}
