// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pubsub

import (
	"sync"
	"time"

	"golang.org/x/net/context"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
)

// newMessageIterator starts a new streamingMessageIterator.  Stop must be called on the messageIterator
// when it is no longer needed.
// subName is the full name of the subscription to pull messages from.
// ctx is the context to use for acking messages and extending message deadlines.
func newMessageIterator(ctx context.Context, s service, subName string, po *pullOptions) *streamingMessageIterator {
	sp := s.newStreamingPuller(ctx, subName, int32(po.ackDeadline.Seconds()))
	_ = sp.open() // error stored in sp
	return newStreamingMessageIterator(ctx, sp, po)
}

type streamingMessageIterator struct {
	ctx        context.Context
	po         *pullOptions
	sp         *streamingPuller
	kaTicker   *time.Ticker  // keep-alive (deadline extensions)
	ackTicker  *time.Ticker  // message acks
	nackTicker *time.Ticker  // message nacks (more frequent than acks)
	failed     chan struct{} // closed on stream error
	stopped    chan struct{} // closed when Stop is called
	drained    chan struct{} // closed when stopped && no more pending messages
	wg         sync.WaitGroup

	mu                 sync.Mutex
	keepAliveDeadlines map[string]time.Time
	pendingReq         *pb.StreamingPullRequest
	err                error // error from stream failure
}

func newStreamingMessageIterator(ctx context.Context, sp *streamingPuller, po *pullOptions) *streamingMessageIterator {
	// TODO: make kaTicker frequency more configurable. (ackDeadline - 5s) is a
	// reasonable default for now, because the minimum ack period is 10s. This
	// gives us 5s grace.
	keepAlivePeriod := po.ackDeadline - 5*time.Second
	kaTicker := time.NewTicker(keepAlivePeriod)

	// Ack promptly so users don't lose work if client crashes.
	ackTicker := time.NewTicker(100 * time.Millisecond)
	nackTicker := time.NewTicker(100 * time.Millisecond)
	it := &streamingMessageIterator{
		ctx:                ctx,
		sp:                 sp,
		po:                 po,
		kaTicker:           kaTicker,
		ackTicker:          ackTicker,
		nackTicker:         nackTicker,
		failed:             make(chan struct{}),
		stopped:            make(chan struct{}),
		drained:            make(chan struct{}),
		keepAliveDeadlines: map[string]time.Time{},
		pendingReq:         &pb.StreamingPullRequest{},
	}
	it.wg.Add(1)
	go it.sender()
	return it
}

// Subscription.receive will call stop on its messageIterator when finished with it.
// Stop will block until Done has been called on all Messages that have been
// returned by Next, or until the context with which the messageIterator was created
// is cancelled or exceeds its deadline.
func (it *streamingMessageIterator) stop() {
	it.mu.Lock()
	select {
	case <-it.stopped:
	default:
		close(it.stopped)
	}
	it.checkDrained()
	it.mu.Unlock()
	it.wg.Wait()
}

// checkDrained closes the drained channel if the iterator has been stopped and all
// pending messages have either been n/acked or expired.
//
// Called with the lock held.
func (it *streamingMessageIterator) checkDrained() {
	select {
	case <-it.drained:
		return
	default:
	}
	select {
	case <-it.stopped:
		if len(it.keepAliveDeadlines) == 0 {
			close(it.drained)
		}
	default:
	}
}

// Called when a message is acked/nacked.
func (it *streamingMessageIterator) done(ackID string, ack bool) {
	it.mu.Lock()
	defer it.mu.Unlock()
	delete(it.keepAliveDeadlines, ackID)
	if ack {
		it.pendingReq.AckIds = append(it.pendingReq.AckIds, ackID)
	} else {
		it.addDeadlineMod(ackID, 0) // Nack indicated by modifying the deadline to zero.
	}
	it.checkDrained()
}

// addDeadlineMod adds the ack ID to the pending request with the given deadline.
//
// Called with the lock held.
func (it *streamingMessageIterator) addDeadlineMod(ackID string, deadlineSecs int32) {
	pr := it.pendingReq
	pr.ModifyDeadlineAckIds = append(pr.ModifyDeadlineAckIds, ackID)
	pr.ModifyDeadlineSeconds = append(pr.ModifyDeadlineSeconds, deadlineSecs)
}

// fail is called when a stream method returns a permanent error.
func (it *streamingMessageIterator) fail(err error) {
	it.mu.Lock()
	if it.err == nil {
		it.err = err
		close(it.failed)
	}
	it.mu.Unlock()
}

// receive makes a call to the stream's Recv method and returns
// its messages.
func (it *streamingMessageIterator) receive() ([]*Message, error) {
	// Stop retrieving messages if the context is done, the stream
	// failed, or the iterator's Stop method was called.
	select {
	case <-it.ctx.Done():
		return nil, it.ctx.Err()
	default:
	}
	it.mu.Lock()
	err := it.err
	it.mu.Unlock()
	if err != nil {
		return nil, err
	}
	// Receive messages from stream. This may block indefinitely.
	msgs, err := it.sp.fetchMessages()
	// The streamingPuller handles retries, so any error here
	// is fatal.
	if err != nil {
		it.fail(err)
		return nil, err
	}
	// We received some messages. Remember them so we can
	// keep them alive.
	deadline := time.Now().Add(it.po.maxExtension)
	it.mu.Lock()
	for _, m := range msgs {
		m.doneFunc = it.done
		it.keepAliveDeadlines[m.ackID] = deadline
	}
	it.mu.Unlock()
	return msgs, nil
}

// sender runs in a goroutine and handles all sends to the stream.
func (it *streamingMessageIterator) sender() {
	defer it.wg.Done()
	defer it.kaTicker.Stop()
	defer it.ackTicker.Stop()
	defer it.nackTicker.Stop()
	defer it.sp.closeSend()

	done := false
	for !done {
		send := false
		select {
		case <-it.ctx.Done():
			// Context canceled or timed out: stop immediately, without
			// another RPC.
			return

		case <-it.failed:
			// Stream failed: nothing to do, so stop immediately.
			return

		case <-it.drained:
			// All outstanding messages have been marked done:
			// nothing left to do except send the final request.
			it.mu.Lock()
			send = (len(it.pendingReq.AckIds) > 0 || len(it.pendingReq.ModifyDeadlineAckIds) > 0)
			done = true

		case <-it.kaTicker.C:
			it.mu.Lock()
			send = it.handleKeepAlives()

		case <-it.nackTicker.C:
			it.mu.Lock()
			send = (len(it.pendingReq.ModifyDeadlineAckIds) > 0)

		case <-it.ackTicker.C:
			it.mu.Lock()
			send = (len(it.pendingReq.AckIds) > 0)

		}
		// Lock is held here.
		if send {
			req := it.pendingReq
			it.pendingReq = &pb.StreamingPullRequest{}
			it.mu.Unlock()
			err := it.sp.send(req)
			if err != nil {
				// The streamingPuller handles retries, so any error here
				// is fatal to the iterator.
				it.fail(err)
				return
			}
		} else {
			it.mu.Unlock()
		}
	}
}

// handleKeepAlives modifies the pending request to include deadline extensions
// for live messages. It also purges expired messages. It reports whether
// there were any live messages.
//
// Called with the lock held.
func (it *streamingMessageIterator) handleKeepAlives() bool {
	live, expired := getKeepAliveAckIDs(it.keepAliveDeadlines)
	for _, e := range expired {
		delete(it.keepAliveDeadlines, e)
	}
	dl := trunc32(int64(it.po.ackDeadline.Seconds()))
	for _, m := range live {
		it.addDeadlineMod(m, dl)
	}
	it.checkDrained()
	return len(live) > 0
}

func getKeepAliveAckIDs(items map[string]time.Time) (live, expired []string) {
	now := time.Now()
	for id, expiry := range items {
		if expiry.Before(now) {
			expired = append(expired, id)
		} else {
			live = append(live, id)
		}
	}
	return live, expired
}
