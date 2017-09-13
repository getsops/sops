// +build js

package sync

import "github.com/gopherjs/gopherjs/js"

var semWaiters = make(map[*uint32][]chan bool)

func runtime_Semacquire(s *uint32) {
	if *s == 0 {
		ch := make(chan bool)
		semWaiters[s] = append(semWaiters[s], ch)
		<-ch
	}
	*s--
}

// SemacquireMutex is like Semacquire, but for profiling contended Mutexes.
// Mutex profiling is not supported, so just use the same implementation as runtime_Semacquire.
// TODO: Investigate this. If it's possible to implement, consider doing so, otherwise remove this comment.
func runtime_SemacquireMutex(s *uint32, lifo bool) {
	// TODO: Use lifo if needed/possible.
	runtime_Semacquire(s)
}

func runtime_Semrelease(s *uint32, handoff bool) {
	// TODO: Use handoff if needed/possible.
	*s++

	w := semWaiters[s]
	if len(w) == 0 {
		return
	}

	ch := w[0]
	w = w[1:]
	semWaiters[s] = w
	if len(w) == 0 {
		delete(semWaiters, s)
	}

	ch <- true
}

func runtime_notifyListCheck(size uintptr) {}

func runtime_canSpin(i int) bool {
	return false
}

// Copy of time.runtimeNano.
func runtime_nanotime() int64 {
	const millisecond = 1000000
	return js.Global.Get("Date").New().Call("getTime").Int64() * millisecond
}
