package tests

import (
	"testing"
	"time"
)

func inner(ch chan struct{}, b bool) ([]byte, error) {
	// ensure gopherjs thinks that this inner function can block
	if b {
		<-ch
	}
	return []byte{}, nil
}

// this function's call to inner never blocks, but the deferred
// statement does.
func outer(ch chan struct{}, b bool) ([]byte, error) {
	defer func() {
		<-ch
	}()

	return inner(ch, b)
}

func TestBlockingInDefer(t *testing.T) {
	defer func() {
		if x := recover(); x != nil {
			t.Error("run time panic: %v", x)
		}
	}()

	ch := make(chan struct{})
	b := false

	go func() {
		time.Sleep(5 * time.Millisecond)
		ch <- struct{}{}
	}()

	outer(ch, b)
}
