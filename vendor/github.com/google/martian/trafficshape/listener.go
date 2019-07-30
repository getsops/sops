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
	"io"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/google/martian/v3/log"
)

// ErrForceClose is an error that communicates the need to close the connection.
type ErrForceClose struct {
	message string
}

func (efc *ErrForceClose) Error() string {
	return efc.message
}

// DefaultBitrate represents the bitrate that will be for all url regexs for which a shape
// has not been specified.
var DefaultBitrate int64 = 500000000000 // 500Gbps (unlimited)

// urlShape contains a rw lock protected shape of a url_regex.
type urlShape struct {
	sync.RWMutex
	Shape *Shape
}

// urlShapes contains a rw lock protected map of url regexs to their URLShapes.
type urlShapes struct {
	sync.RWMutex
	M                map[string]*urlShape
	LastModifiedTime time.Time
}

// Buckets contains the read and write buckets for a url_regex.
type Buckets struct {
	ReadBucket  *Bucket
	WriteBucket *Bucket
}

// NewBuckets returns a *Buckets with the specified up and down bandwidths.
func NewBuckets(up int64, down int64) *Buckets {
	return &Buckets{
		ReadBucket:  NewBucket(up, time.Second),
		WriteBucket: NewBucket(down, time.Second),
	}
}

// ThrottleContext represents whether we are currently in a throttle interval for a particular
// url_regex. If ThrottleNow is true, only then will the current throttle 'Bandwidth' be set
// correctly.
type ThrottleContext struct {
	ThrottleNow bool
	Bandwidth   int64
}

// NextActionInfo represents whether there is an upcoming action. Only if ActionNext is true will the
// Index and ByteOffset be set correctly.
type NextActionInfo struct {
	ActionNext bool
	Index      int64
	ByteOffset int64
}

// Context represents the current information that is needed while writing back to the client.
// Only if Shaping is true, that is we are currently writing back a response that matches a certain
// url_regex will the other values be set correctly. If so, the Buckets represent the buckets
// to be used for the current url_regex. NextActionInfo tells us whether there is an upcoming action
// that needs to be performed, and ThrottleContext tells us whether we are currently in a throttle
// interval (according to the RangeStart). Note, the ThrottleContext is only used once in the start
// to determine the beginning bandwidth. It need not be updated after that. This
// is because the subsequent throttles are captured in the upcoming ChangeBandwidth actions.
// Byte Offset represents the absolute byte offset of response data that we are currently writing back.
// It does not account for the header data.
type Context struct {
	Shaping            bool
	RangeStart         int64
	URLRegex           string
	Buckets            *Buckets
	GlobalBucket       *Bucket
	ThrottleContext    *ThrottleContext
	NextActionInfo     *NextActionInfo
	ByteOffset         int64
	HeaderLen          int64
	HeaderBytesWritten int64
}

// Listener wraps a net.Listener and simulates connection latency and bandwidth
// constraints.
type Listener struct {
	net.Listener

	ReadBucket  *Bucket
	WriteBucket *Bucket

	mu            sync.RWMutex
	latency       time.Duration
	GlobalBuckets map[string]*Bucket
	Shapes        *urlShapes
	defaults      *Default
}

// Conn wraps a net.Conn and simulates connection latency and bandwidth
// constraints. Shapes represents the traffic shape map inherited from the listener.
// Established is the time that this connection was established. LocalBuckets represents a map from
// the url_regexes to their dedicated buckets.
type Conn struct {
	net.Conn
	ReadBucket       *Bucket // Shared by listener.
	WriteBucket      *Bucket // Shared by listener.
	latency          time.Duration
	ronce            sync.Once
	wonce            sync.Once
	Shapes           *urlShapes
	GlobalBuckets    map[string]*Bucket
	LocalBuckets     map[string]*Buckets
	Established      time.Time
	Context          *Context
	DefaultBandwidth Bandwidth
	Listener         *Listener
}

// NewListener returns a new bandwidth constrained listener. Defaults to
// DefaultBitrate (uncapped).
func NewListener(l net.Listener) *Listener {
	return &Listener{
		Listener:      l,
		ReadBucket:    NewBucket(DefaultBitrate/8, time.Second),
		WriteBucket:   NewBucket(DefaultBitrate/8, time.Second),
		Shapes:        &urlShapes{M: make(map[string]*urlShape)},
		GlobalBuckets: make(map[string]*Bucket),
		defaults: &Default{
			Bandwidth: Bandwidth{
				Up:   DefaultBitrate / 8,
				Down: DefaultBitrate / 8,
			},
			Latency: 0,
		},
	}
}

// ReadBitrate returns the bitrate in bits per second for reads.
func (l *Listener) ReadBitrate() int64 {
	return l.ReadBucket.Capacity() * 8
}

// SetReadBitrate sets the bitrate in bits per second for reads.
func (l *Listener) SetReadBitrate(bitrate int64) {
	l.ReadBucket.SetCapacity(bitrate / 8)
}

// WriteBitrate returns the bitrate in bits per second for writes.
func (l *Listener) WriteBitrate() int64 {
	return l.WriteBucket.Capacity() * 8
}

// SetWriteBitrate sets the bitrate in bits per second for writes.
func (l *Listener) SetWriteBitrate(bitrate int64) {
	l.WriteBucket.SetCapacity(bitrate / 8)
}

// SetDefaults sets the default traffic shaping parameters for the listener.
func (l *Listener) SetDefaults(defaults *Default) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.defaults = defaults
}

// Defaults returns the default traffic shaping parameters for the listener.
func (l *Listener) Defaults() *Default {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.defaults
}

// Latency returns the latency for connections.
func (l *Listener) Latency() time.Duration {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.latency
}

// SetLatency sets the initial latency for connections.
func (l *Listener) SetLatency(latency time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.latency = latency
}

// GetTrafficShapedConn takes in a normal connection and returns a traffic shaped connection.
func (l *Listener) GetTrafficShapedConn(oc net.Conn) *Conn {
	if tsconn, ok := oc.(*Conn); ok {
		return tsconn
	}
	urlbuckets := make(map[string]*Buckets)
	globalurlbuckets := make(map[string]*Bucket)

	l.Shapes.RLock()
	defaults := l.Defaults()
	latency := l.Latency()
	defaultBandwidth := defaults.Bandwidth
	for regex, shape := range l.Shapes.M {
		// It should be ok to not acquire the read lock on shape, since WriteBucket is never mutated.
		globalurlbuckets[regex] = shape.Shape.WriteBucket
		urlbuckets[regex] = NewBuckets(DefaultBitrate/8, shape.Shape.MaxBandwidth)
	}

	l.Shapes.RUnlock()

	curinfo := &Context{}

	lc := &Conn{
		Conn:             oc,
		latency:          latency,
		ReadBucket:       l.ReadBucket,
		WriteBucket:      l.WriteBucket,
		Shapes:           l.Shapes,
		GlobalBuckets:    globalurlbuckets,
		LocalBuckets:     urlbuckets,
		Context:          curinfo,
		Established:      time.Now(),
		DefaultBandwidth: defaultBandwidth,
		Listener:         l,
	}
	return lc
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	oc, err := l.Listener.Accept()
	if err != nil {
		log.Errorf("trafficshape: failed accepting connection: %v", err)
		return nil, err
	}

	if tconn, ok := oc.(*net.TCPConn); ok {
		log.Debugf("trafficshape: setting keep-alive for TCP connection")
		tconn.SetKeepAlive(true)
		tconn.SetKeepAlivePeriod(3 * time.Minute)
	}
	return l.GetTrafficShapedConn(oc), nil
}

// Close closes the read and write buckets along with the underlying listener.
func (l *Listener) Close() error {
	defer log.Debugf("trafficshape: closed read/write buckets and connection")

	l.ReadBucket.Close()
	l.WriteBucket.Close()

	return l.Listener.Close()
}

// Read reads bytes from connection into b, optionally simulating connection
// latency and throttling read throughput based on desired bandwidth
// constraints.
func (c *Conn) Read(b []byte) (int, error) {
	c.ronce.Do(c.sleepLatency)

	n, err := c.ReadBucket.FillThrottle(func(remaining int64) (int64, error) {
		max := remaining
		if l := int64(len(b)); max > l {
			max = l
		}

		n, err := c.Conn.Read(b[:max])
		return int64(n), err
	})
	if err != nil && err != io.EOF {
		log.Errorf("trafficshape: error on throttled read: %v", err)
	}

	return int(n), err
}

// ReadFrom reads data from r until EOF or error, optionally simulating
// connection latency and throttling read throughput based on desired bandwidth
// constraints.
func (c *Conn) ReadFrom(r io.Reader) (int64, error) {
	c.ronce.Do(c.sleepLatency)

	var total int64
	for {
		n, err := c.ReadBucket.FillThrottle(func(remaining int64) (int64, error) {
			return io.CopyN(c.Conn, r, remaining)
		})

		total += n

		if err == io.EOF {
			log.Debugf("trafficshape: exhausted reader successfully")
			return total, nil
		} else if err != nil {
			log.Errorf("trafficshape: failed copying from reader: %v", err)
			return total, err
		}
	}
}

// WriteTo writes data to w from the connection, optionally simulating
// connection latency and throttling write throughput based on desired
// bandwidth constraints.
func (c *Conn) WriteTo(w io.Writer) (int64, error) {
	c.wonce.Do(c.sleepLatency)

	var total int64
	for {
		n, err := c.WriteBucket.FillThrottle(func(remaining int64) (int64, error) {
			return io.CopyN(w, c.Conn, remaining)
		})

		total += n

		if err != nil {
			if err != io.EOF {
				log.Errorf("trafficshape: failed copying to writer: %v", err)
			}
			return total, err
		}
	}
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

// CheckExistenceAndValidity checks that the current url regex is present in the map, and that
// the connection was established before the url shape map was last updated. We do not allow the
// updated url shape map to traffic shape older connections.
// Important: Assumes you have acquired the required locks and will release them youself.
func (c *Conn) CheckExistenceAndValidity(URLRegex string) bool {
	shapeStillValid := c.Shapes.LastModifiedTime.Before(c.Established)
	_, p := c.Shapes.M[URLRegex]
	return p && shapeStillValid
}

// GetCurrentThrottle uses binary search to determine if the current byte offset ('start')
// lies within a throttle interval. If so, also returns the bandwidth specified for that interval.
func (c *Conn) GetCurrentThrottle(start int64) *ThrottleContext {
	c.Shapes.RLock()
	defer c.Shapes.RUnlock()

	if !c.CheckExistenceAndValidity(c.Context.URLRegex) {
		log.Debugf("existence check failed")
		return &ThrottleContext{
			ThrottleNow: false,
		}
	}

	c.Shapes.M[c.Context.URLRegex].RLock()
	defer c.Shapes.M[c.Context.URLRegex].RUnlock()

	throttles := c.Shapes.M[c.Context.URLRegex].Shape.Throttles

	if l := len(throttles); l != 0 {
		// ind is the first index in throttles with ByteStart > start.
		// Once we get ind, we can check the previous throttle, if any,
		// to see if its ByteEnd is after 'start'.
		ind := sort.Search(len(throttles),
			func(i int) bool { return throttles[i].ByteStart > start })

		// All throttles have Bytestart > start, hence not in throttle.
		if ind == 0 {
			return &ThrottleContext{
				ThrottleNow: false,
			}
		}

		// No throttle has Bytestart > start, so check the last throttle to
		// see if it ends after 'start'. Note: the last throttle is special
		// since it can have -1 (meaning infinity) as the ByteEnd.
		if ind == l {
			if throttles[l-1].ByteEnd > start || throttles[l-1].ByteEnd == -1 {
				return &ThrottleContext{
					ThrottleNow: true,
					Bandwidth:   throttles[l-1].Bandwidth,
				}
			}
			return &ThrottleContext{
				ThrottleNow: false,
			}
		}

		// Check the previous throttle to see if it ends after 'start'.
		if throttles[ind-1].ByteEnd > start {
			return &ThrottleContext{
				ThrottleNow: true,
				Bandwidth:   throttles[ind-1].Bandwidth,
			}
		}

		return &ThrottleContext{
			ThrottleNow: false,
		}
	}

	return &ThrottleContext{
		ThrottleNow: false,
	}
}

// GetNextActionFromByte takes in a byte offset and uses binary search to determine the upcoming
// action, i.e the first action after the byte that still has a non zero count.
func (c *Conn) GetNextActionFromByte(start int64) *NextActionInfo {
	c.Shapes.RLock()
	defer c.Shapes.RUnlock()

	if !c.CheckExistenceAndValidity(c.Context.URLRegex) {
		log.Debugf("existence check failed")
		return &NextActionInfo{
			ActionNext: false,
		}
	}

	c.Shapes.M[c.Context.URLRegex].RLock()
	defer c.Shapes.M[c.Context.URLRegex].RUnlock()

	actions := c.Shapes.M[c.Context.URLRegex].Shape.Actions

	if l := len(actions); l != 0 {
		ind := sort.Search(len(actions),
			func(i int) bool { return actions[i].getByte() >= start })

		return c.GetNextActionFromIndex(int64(ind))
	}

	return &NextActionInfo{
		ActionNext: false,
	}
}

// GetNextActionFromIndex takes in an index and returns the first action after the index that
// has a non zero count, if there is one.
func (c *Conn) GetNextActionFromIndex(ind int64) *NextActionInfo {
	c.Shapes.RLock()
	defer c.Shapes.RUnlock()

	if !c.CheckExistenceAndValidity(c.Context.URLRegex) {
		return &NextActionInfo{
			ActionNext: false,
		}
	}

	c.Shapes.M[c.Context.URLRegex].RLock()
	defer c.Shapes.M[c.Context.URLRegex].RUnlock()

	actions := c.Shapes.M[c.Context.URLRegex].Shape.Actions

	if l := int64(len(actions)); l != 0 {

		for ind < l && (actions[ind].getCount() == 0) {
			ind++
		}

		if ind >= l {
			return &NextActionInfo{
				ActionNext: false,
			}
		}
		return &NextActionInfo{
			ActionNext: true,
			Index:      ind,
			ByteOffset: actions[ind].getByte(),
		}
	}
	return &NextActionInfo{
		ActionNext: false,
	}
}

// WriteDefaultBuckets writes bytes from b to the connection, optionally simulating
// connection latency and throttling write throughput based on desired
// bandwidth constraints. It uses the WriteBucket inherited from the listener.
func (c *Conn) WriteDefaultBuckets(b []byte) (int, error) {
	c.wonce.Do(c.sleepLatency)

	var total int64
	for len(b) > 0 {
		var max int64

		n, err := c.WriteBucket.FillThrottle(func(remaining int64) (int64, error) {
			max = remaining
			if l := int64(len(b)); remaining >= l {
				max = l
			}

			n, err := c.Conn.Write(b[:max])
			return int64(n), err
		})

		total += n

		if err != nil {
			if err != io.EOF {
				log.Errorf("trafficshape: failed write: %v", err)
			}
			return int(total), err
		}

		b = b[max:]
	}

	return int(total), nil
}

// Write writes bytes from b to the connection, while enforcing throttles and performing actions.
// It uses and updates the Context in the connection.
func (c *Conn) Write(b []byte) (int, error) {
	if !c.Context.Shaping {
		return c.WriteDefaultBuckets(b)
	}
	c.wonce.Do(c.sleepLatency)
	var total int64

	// Write the header if needed, without enforcing any traffic shaping, and without updating
	// ByteOffset.
	if headerToWrite := c.Context.HeaderLen - c.Context.HeaderBytesWritten; headerToWrite > 0 {
		writeAmount := min(int64(len(b)), headerToWrite)

		n, err := c.Conn.Write(b[:writeAmount])

		if err != nil {
			if err != io.EOF {
				log.Errorf("trafficshape: failed write: %v", err)
			}
			return int(n), err
		}
		c.Context.HeaderBytesWritten += writeAmount
		total += writeAmount
		b = b[writeAmount:]
	}

	var amountToWrite int64

	for len(b) > 0 {
		var max int64

		// Determine the amount to be written up till the next action.
		amountToWrite = int64(len(b))
		if c.Context.NextActionInfo.ActionNext {
			amountTillNextAction := c.Context.NextActionInfo.ByteOffset - c.Context.ByteOffset
			if amountTillNextAction <= amountToWrite {
				amountToWrite = amountTillNextAction
			}
		}

		// Write into both the local and global buckets, as well as the underlying connection.
		n, err := c.Context.Buckets.WriteBucket.FillThrottleLocked(func(remaining int64) (int64, error) {
			max = min(remaining, amountToWrite)

			if max == 0 {
				return 0, nil
			}

			return c.Context.GlobalBucket.FillThrottleLocked(func(rem int64) (int64, error) {
				max = min(rem, max)
				n, err := c.Conn.Write(b[:max])

				return int64(n), err
			})
		})

		if err != nil {
			if err != io.EOF {
				log.Errorf("trafficshape: failed write: %v", err)
			}
			return int(total), err
		}

		// Update the current byte offset.
		c.Context.ByteOffset += n
		total += n

		b = b[max:]

		// Check if there was an upcoming action, and that the byte offset matches the action's byte.
		if c.Context.NextActionInfo.ActionNext &&
			c.Context.ByteOffset >= c.Context.NextActionInfo.ByteOffset {
			// Note here, we check again that the url shape map is still valid and that the action still has
			// a non zero count, since that could have been modified since the last time we checked.
			ind := c.Context.NextActionInfo.Index
			c.Shapes.RLock()
			if !c.CheckExistenceAndValidity(c.Context.URLRegex) {
				c.Shapes.RUnlock()
				// Write the remaining b using default buckets, and set Shaping as false
				// so that subsequent calls to Write() also use default buckets
				// without performing any actions.
				c.Context.Shaping = false
				writeTotal, e := c.WriteDefaultBuckets(b)
				return int(total) + writeTotal, e
			}
			c.Shapes.M[c.Context.URLRegex].Lock()
			actions := c.Shapes.M[c.Context.URLRegex].Shape.Actions
			if actions[ind].getCount() != 0 {
				// Update the action count, determine the type of action and perform it.
				actions[ind].decrementCount()
				switch action := actions[ind].(type) {
				case *Halt:
					d := action.Duration
					log.Debugf("trafficshape: Sleeping for time %d ms for urlregex %s at byte offset %d",
						d, c.Context.URLRegex, c.Context.ByteOffset)
					c.Shapes.M[c.Context.URLRegex].Unlock()
					c.Shapes.RUnlock()
					time.Sleep(time.Duration(d) * time.Millisecond)
				case *CloseConnection:
					log.Infof("trafficshape: Closing connection for urlregex %s at byte offset %d",
						c.Context.URLRegex, c.Context.ByteOffset)
					c.Shapes.M[c.Context.URLRegex].Unlock()
					c.Shapes.RUnlock()
					return int(total), &ErrForceClose{message: "Forcing close connection"}
				case *ChangeBandwidth:
					bw := action.Bandwidth
					log.Infof("trafficshape: Changing connection bandwidth to %d for urlregex %s at byte offset %d",
						bw, c.Context.URLRegex, c.Context.ByteOffset)
					c.Shapes.M[c.Context.URLRegex].Unlock()
					c.Shapes.RUnlock()
					c.Context.Buckets.WriteBucket.SetCapacity(bw)
				default:
					c.Shapes.M[c.Context.URLRegex].Unlock()
					c.Shapes.RUnlock()
				}
			} else {
				c.Shapes.M[c.Context.URLRegex].Unlock()
				c.Shapes.RUnlock()
			}
			// Get the next action to be performed, if any.
			c.Context.NextActionInfo = c.GetNextActionFromIndex(ind + 1)
		}
	}
	return int(total), nil
}

func (c *Conn) sleepLatency() {
	log.Debugf("trafficshape: simulating latency: %s", c.latency)
	time.Sleep(c.latency)
}
