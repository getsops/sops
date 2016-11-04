// +build js

package sync

var semWaiters = make(map[*uint32][]chan bool)

func runtime_Semacquire(s *uint32) {
	if *s == 0 {
		ch := make(chan bool)
		semWaiters[s] = append(semWaiters[s], ch)
		<-ch
	}
	*s--
}

func runtime_Semrelease(s *uint32) {
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
