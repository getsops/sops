// +build js

package strings_test

import "testing"

func TestBuilderAllocs(t *testing.T) {
	t.Skip("runtime.ReadMemStats, testing.AllocsPerRun not supported in GopherJS")
}

func TestBuilderGrow(t *testing.T) {
	t.Skip("runtime.ReadMemStats, testing.AllocsPerRun not supported in GopherJS")
}
