// +build js

package bytes_test

import (
	"testing"
)

func dangerousSlice(t *testing.T) []byte {
	t.Skip("dangerousSlice relies on syscall.Getpagesize, which GopherJS doesn't implement")

	panic("unreachable")
}
