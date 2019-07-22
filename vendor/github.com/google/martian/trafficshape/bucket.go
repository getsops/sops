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
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/martian/v3/log"
)

// Bucket is a generic leaky bucket that drains at a configurable interval and
// fills at user defined rate. The bucket may be used concurrently.
type Bucket struct {
	capacity int64 // atomic
	fill     int64 // atomic
	mu       sync.Mutex

	t      *time.Ticker
	closec chan struct{}
}

var (
	// ErrBucketOverflow is an error that indicates the bucket has been overflown
	// by the user. This error is only returned iff fill > capacity.
	ErrBucketOverflow   = errors.New("trafficshape: bucket overflow")
	errFillClosedBucket = errors.New("trafficshape: fill on closed bucket")
)

// NewBucket returns a new leaky bucket with capacity that is drained
// at interval.
func NewBucket(capacity int64, interval time.Duration) *Bucket {
	b := &Bucket{
		capacity: capacity,
		t:        time.NewTicker(interval),
		closec:   make(chan struct{}),
	}

	go b.loop()

	return b
}

// Capacity returns the capacity of the bucket.
func (b *Bucket) Capacity() int64 {
	return atomic.LoadInt64(&b.capacity)
}

// SetCapacity sets the capacity for the bucket and resets the fill to zero.
func (b *Bucket) SetCapacity(capacity int64) {
	log.Infof("trafficshape: set capacity: %d", capacity)

	atomic.StoreInt64(&b.capacity, capacity)
	atomic.StoreInt64(&b.fill, 0)
}

// Close stops the drain loop and marks the bucket as closed.
func (b *Bucket) Close() error {
	log.Debugf("trafficshape: closing bucket")

	// Allow b to be closed multiple times without panicking.
	if b.closed() {
		return nil
	}

	b.t.Stop()
	close(b.closec)

	return nil
}

// FillThrottle calls fn with the available capacity remaining (capacity-fill)
// and fills the bucket with the number of tokens returned by fn. If the
// remaining capacity is <= 0, FillThrottle will wait for the next drain before
// running fn.
//
// If fn returns an error, it will be returned by FillThrottle along with the
// number of tokens processed by fn.
//
// fn is provided the remaining capacity as a soft maximum, fn is allowed to
// use more than the remaining capacity without incurring spillage.
//
// If the bucket is closed when FillThrottle is called, or while waiting for
// the next drain, fn will not be executed and FillThrottle will return with an
// error.
func (b *Bucket) FillThrottle(fn func(int64) (int64, error)) (int64, error) {
	for {
		if b.closed() {
			log.Errorf("trafficshape: fill on closed bucket")
			return 0, errFillClosedBucket
		}

		fill := atomic.LoadInt64(&b.fill)
		capacity := atomic.LoadInt64(&b.capacity)

		if fill < capacity {
			log.Debugf("trafficshape: under capacity (%d/%d)", fill, capacity)

			n, err := fn(capacity - fill)
			fill = atomic.AddInt64(&b.fill, n)

			return n, err
		}

		log.Debugf("trafficshape: bucket full (%d/%d)", fill, capacity)
	}
}

// FillThrottleLocked is like FillThrottle, except that it uses a lock to protect
// the critical section between accessing the fill value and updating it.
func (b *Bucket) FillThrottleLocked(fn func(int64) (int64, error)) (int64, error) {
	for {
		if b.closed() {
			log.Errorf("trafficshape: fill on closed bucket")
			return 0, errFillClosedBucket
		}

		b.mu.Lock()
		fill := atomic.LoadInt64(&b.fill)
		capacity := atomic.LoadInt64(&b.capacity)

		if fill < capacity {

			n, err := fn(capacity - fill)
			fill = atomic.AddInt64(&b.fill, n)
			b.mu.Unlock()
			return n, err
		}
		b.mu.Unlock()

		log.Debugf("trafficshape: bucket full (%d/%d)", fill, capacity)
	}
}

// Fill calls fn with the available capacity remaining (capacity-fill) and
// fills the bucket with the number of tokens returned by fn. If the remaining
// capacity is 0, Fill returns 0, nil. If the remaining capacity is < 0, Fill
// returns 0, ErrBucketOverflow.
//
// If fn returns an error, it will be returned by Fill along with the remaining
// capacity.
//
// fn is provided the remaining capacity as a soft maximum, fn is allowed to
// use more than the remaining capacity without incurring spillage, though this
// will cause subsequent calls to Fill to return ErrBucketOverflow until the
// next drain.
//
// If the bucket is closed when Fill is called, fn will not be executed and
// Fill will return with an error.
func (b *Bucket) Fill(fn func(int64) (int64, error)) (int64, error) {
	if b.closed() {
		log.Errorf("trafficshape: fill on closed bucket")
		return 0, errFillClosedBucket
	}

	fill := atomic.LoadInt64(&b.fill)
	capacity := atomic.LoadInt64(&b.capacity)

	switch {
	case fill < capacity:
		log.Debugf("trafficshape: under capacity (%d/%d)", fill, capacity)

		n, err := fn(capacity - fill)
		fill = atomic.AddInt64(&b.fill, n)

		return n, err
	case fill > capacity:
		log.Debugf("trafficshape: bucket overflow (%d/%d)", fill, capacity)

		return 0, ErrBucketOverflow
	}

	log.Debugf("trafficshape: bucket full (%d/%d)", fill, capacity)
	return 0, nil
}

// loop drains the fill at interval and returns when the bucket is closed.
func (b *Bucket) loop() {
	log.Debugf("trafficshape: started drain loop")
	defer log.Debugf("trafficshape: stopped drain loop")

	for {
		select {
		case t := <-b.t.C:
			atomic.StoreInt64(&b.fill, 0)
			log.Debugf("trafficshape: fill reset @ %s", t)
		case <-b.closec:
			log.Debugf("trafficshape: bucket closed")
			return
		}
	}
}

func (b *Bucket) closed() bool {
	select {
	case <-b.closec:
		return true
	default:
		return false
	}
}

