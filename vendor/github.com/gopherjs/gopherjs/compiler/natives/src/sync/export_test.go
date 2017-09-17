// +build js

package sync

// Referenced by tests, need to have no-op implementations.
var Runtime_procPin = func() int { return 0 }
var Runtime_procUnpin = func() {}
