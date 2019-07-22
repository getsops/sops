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
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/martian/v3/log"
)

// Handler configures a trafficshape.Listener.
type Handler struct {
	l *Listener
}

// Throttle represents a byte interval with a specific bandwidth.
type Throttle struct {
	Bytes     string `json:"bytes"`
	Bandwidth int64  `json:"bandwidth"`
	ByteStart int64
	ByteEnd   int64
}

// Action represents an arbitrary event that needs to be executed while writing back to the client.
type Action interface {
	// Byte offset to perform Action at.
	getByte() int64
	// Number of times to perform the action. -1 for infinite times.
	getCount() int64
	// Update the count when performing an action.
	decrementCount()
}

// Halt is the event that represents a period of time to sleep while writing.
// It implements the Action interface.
type Halt struct {
	Byte     int64 `json:"byte"`
	Duration int64 `json:"duration"`
	Count    int64 `json:"count"`
}

func (h *Halt) getByte() int64 {
	return h.Byte
}

func (h *Halt) getCount() int64 {
	return h.Count
}

func (h *Halt) decrementCount() {
	if h.Count > 0 {
		h.Count--
	}
}

// CloseConnection is an event that represents the closing of a connection with a client.
// It implements the Action interface.
type CloseConnection struct {
	Byte  int64 `json:"byte"`
	Count int64 `json:"count"`
}

func (cc *CloseConnection) getByte() int64 {
	return cc.Byte
}

func (cc *CloseConnection) getCount() int64 {
	return cc.Count
}

func (cc *CloseConnection) decrementCount() {
	if cc.Count > 0 {
		cc.Count--
	}
}

// Shape encloses the traffic shape of a particular url regex.
type Shape struct {
	URLRegex         string             `json:"url_regex"`
	MaxBandwidth     int64              `json:"max_global_bandwidth"`
	Throttles        []*Throttle        `json:"throttles"`
	Halts            []*Halt            `json:"halts"`
	CloseConnections []*CloseConnection `json:"close_connections"`
	// Actions are populated after processing Throttles, Halts and CloseConnections.
	// Actions is sorted in the order of byte offset.
	Actions []Action
	// WriteBucket is initialized by us using MaxBandwidth.
	WriteBucket *Bucket
}

// Bandwidth encloses information about the upstream and downstream bandwidths.
type Bandwidth struct {
	Up   int64 `json:"up"`
	Down int64 `json:"down"`
}

// Default encloses information about the default traffic shaping parameters: bandwidth and latency.
type Default struct {
	Bandwidth Bandwidth `json:"bandwidth"`
	Latency   int64     `json:"latency"`
}

// Trafficshape contains global shape of traffic, i.e information about shape of each url specified and
// the default traffic shaping parameters.
type Trafficshape struct {
	Defaults *Default `json:"default"`
	Shapes   []*Shape `json:"shapes"`
}

// ConfigRequest represents a request to configure the global traffic shape.
type ConfigRequest struct {
	Trafficshape *Trafficshape `json:"trafficshape"`
}

// ChangeBandwidth represents the event of changing the current bandwidth. It is used as an
// endpoint of a Throttle. It implements the Action interface.
type ChangeBandwidth struct {
	Byte      int64
	Bandwidth int64
}

func (cb *ChangeBandwidth) getByte() int64 {
	return cb.Byte
}

func (cb *ChangeBandwidth) getCount() int64 {
	return -1
}

// No op. This is because Throttles have infinite count.
func (cb *ChangeBandwidth) decrementCount() {
}

// NewHandler returns an http.Handler to configure traffic shaping.
func NewHandler(l *Listener) *Handler {
	return &Handler{
		l: l,
	}
}

// ServeHTTP configures latency and bandwidth constraints.
//
// The "latency" query string parameter accepts a duration string in any format
// supported by time.ParseDuration.
// The "up" and "down" query string parameters accept integers as bits per
// second to be used for read and write throughput.
func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log.Infof("trafficshape: configuration request")

	receivedConfig := &ConfigRequest{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "Error reading request body", 400)
		return
	}
	bodystr := string(body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := json.NewDecoder(req.Body).Decode(&receivedConfig); err != nil {
		log.Errorf("Error while parsing the received json: %v", err)
		http.Error(rw, err.Error(), 400)
		return
	}

	if receivedConfig.Trafficshape == nil {
		http.Error(rw, "Error: trafficshape property not found", 400)
		return
	}

	defaults := receivedConfig.Trafficshape.Defaults
	if defaults == nil {
		defaults = &Default{}
	}

	if defaults.Bandwidth.Up < 0 || defaults.Bandwidth.Down < 0 || defaults.Latency < 0 {
		http.Error(rw, "Error: Invalid Defaults", 400)
		return
	}

	if defaults.Bandwidth.Up == 0 {
		defaults.Bandwidth.Up = DefaultBitrate / 8
	}
	if defaults.Bandwidth.Down == 0 {
		defaults.Bandwidth.Down = DefaultBitrate / 8
	}

	// Parse and verify the received shapes.
	if err := parseShapes(receivedConfig.Trafficshape); err != nil {
		http.Error(rw, err.Error(), 400)
		return
	}

	// Update the Listener with the new traffic shape.
	h.l.Shapes.Lock()

	h.l.Shapes.LastModifiedTime = time.Now()
	h.l.ReadBucket.SetCapacity(defaults.Bandwidth.Down)
	h.l.WriteBucket.SetCapacity(defaults.Bandwidth.Up)
	h.l.SetLatency(time.Duration(defaults.Latency) * time.Millisecond)
	h.l.SetDefaults(defaults)

	h.l.Shapes.M = make(map[string]*urlShape)
	for _, shape := range receivedConfig.Trafficshape.Shapes {
		h.l.Shapes.M[shape.URLRegex] = &urlShape{Shape: shape}
	}
	// Update the time that the map was last modified to the current time.
	h.l.Shapes.LastModifiedTime = time.Now()
	h.l.Shapes.Unlock()

	rw.WriteHeader(http.StatusOK)
	io.WriteString(rw, bodystr)
}
