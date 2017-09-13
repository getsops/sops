// +build js

package sync_test

import (
	"testing"
)

func TestPool(t *testing.T) {
	t.Skip()
}

func TestPoolGC(t *testing.T) {
	t.Skip()
}

func TestPoolRelease(t *testing.T) {
	t.Skip()
}

func TestCondCopy(t *testing.T) {
	t.Skip()
}

// TODO: Investigate, fix if possible.
//       It fails with "can't acquire Mutex in 10 seconds"
//       when using Go 1.8 sync.Mutex implementation.
//       It panics with "sync: inconsistent mutex state"
//       with Go 1.9 sync.Mutex implementation.
func TestMutexFairness(t *testing.T) {
	t.Skip("TestMutexFairness fails")
}
