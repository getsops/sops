// Copyright 2015 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trafficshape

import (
	"errors"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

func TestBucket(t *testing.T) {
	t.Parallel()

	b := NewBucket(10, 10*time.Millisecond)
	defer b.Close()

	if got, want := b.Capacity(), int64(10); got != want {
		t.Fatalf("b.Capacity(): got %d, want %d", got, want)
	}

	n, err := b.Fill(func(remaining int64) (int64, error) {
		if want := int64(10); remaining != want {
			t.Errorf("remaining: got %d, want %d", remaining, want)
		}
		return 5, nil
	})
	if err != nil {
		t.Fatalf("Fill(): got %v, want no error", err)
	}
	if got, want := n, int64(5); got != want {
		t.Fatalf("n: got %d, want %d", got, want)
	}

	n, err = b.Fill(func(remaining int64) (int64, error) {
		if want := int64(5); remaining != want {
			t.Errorf("remaining: got %d, want %d", remaining, want)
		}
		return 5, nil
	})
	if err != nil {
		t.Fatalf("Fill(): got %v, want no error", err)
	}
	if got, want := n, int64(5); got != want {
		t.Fatalf("n: got %d, want %d", got, want)
	}
	n, err = b.Fill(func(remaining int64) (int64, error) {
		t.Fatal("Fill: executed func when full, want skipped")
		return 0, nil
	})
	if err != nil {
		t.Fatalf("Fill(): got %v, want no error", err)
	}

	// Wait for the bucket to drain.
	for {
		if atomic.LoadInt64(&b.fill) == 0 {
			break
		}
		// Allow for a goroutine switch, required for GOMAXPROCS = 1.
		runtime.Gosched()
	}

	wanterr := errors.New("fill function error")
	n, err = b.Fill(func(remaining int64) (int64, error) {
		if want := int64(10); remaining != want {
			t.Errorf("remaining: got %d, want %d", remaining, want)
		}
		return 0, wanterr
	})
	if err != wanterr {
		t.Fatalf("Fill(): got %v, want %v", err, wanterr)
	}
	if got, want := n, int64(0); got != want {
		t.Fatalf("n: got %d, want %d", got, want)
	}
}

func TestBucketClosed(t *testing.T) {
	t.Parallel()

	b := NewBucket(0, time.Millisecond)
	b.Close()

	if _, err := b.Fill(nil); err != errFillClosedBucket {
		t.Errorf("Fill(): got %v, want errFillClosedBucket", err)
	}
	if _, err := b.FillThrottle(nil); err != errFillClosedBucket {
		t.Errorf("FillThrottle(): got %v, want errFillClosedBucket", err)
	}
}

func TestBucketOverflow(t *testing.T) {
	t.Parallel()

	b := NewBucket(10, 10*time.Millisecond)
	defer b.Close()

	n, err := b.Fill(func(remaining int64) (int64, error) {
		return 11, nil
	})
	if err != nil {
		t.Fatalf("Fill(): got %v, want no error", err)
	}

	n, err = b.Fill(func(int64) (int64, error) {
		t.Fatal("Fill: executed func when full, want skipped")
		return 0, nil
	})
	if err != ErrBucketOverflow {
		t.Fatalf("Fill(): got %v, want ErrBucketOverflow", err)
	}
	if got, want := n, int64(0); got != want {
		t.Fatalf("n: got %d, want %d", got, want)
	}
}

func TestBucketThrottle(t *testing.T) {
	t.Parallel()

	b := NewBucket(50, 50*time.Millisecond)
	defer b.Close()

	closec := make(chan struct{})
	errc := make(chan error, 1)

	fill := func() {
		for {
			select {
			case <-closec:
				return
			default:
				if _, err := b.FillThrottle(func(remaining int64) (int64, error) {
					if remaining < 10 {
						return remaining, nil
					}
					return 10, nil
				}); err != nil {
					select {
					case errc <- err:
					default:
					}
				}
			}
		}
	}

	for i := 0; i < 5; i++ {
		go fill()
	}

	time.Sleep(time.Second)

	close(closec)

	select {
	case err := <-errc:
		t.Fatalf("FillThrottle: got %v, want no error", err)
	default:
	}
}

func TestBucketFillThrottleCloseBeforeTick(t *testing.T) {
	t.Parallel()

	b := NewBucket(0, time.Minute)
	time.AfterFunc(time.Second, func() { b.Close() })

	if _, err := b.FillThrottle(func(int64) (int64, error) {
		t.Fatal("FillThrottle(): executed func after close, want skipped")
		return 0, nil
	}); err != errFillClosedBucket {
		t.Errorf("b.FillThrottle(): got nil, want errFillClosedBucket")
	}
}
