// Copyright 2018, OpenCensus Authors
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

package metric

import (
	"math"
	"sync"
	"sync/atomic"
	"time"

	"go.opencensus.io/internal/tagencoding"
	"go.opencensus.io/metric/metricdata"
)

// gauge represents a quantity that can go up an down, for example queue depth
// or number of outstanding requests.
//
// gauge maintains a value for each combination of of label values passed to
// the Set or Add methods.
//
// gauge should not be used directly, use Float64Gauge or Int64Gauge.
type gauge struct {
	vals    sync.Map
	desc    metricdata.Descriptor
	start   time.Time
	keys    []string
	isFloat bool
}

type gaugeEntry interface {
	read(t time.Time) metricdata.Point
}

// Read returns the current values of the gauge as a metric for export.
func (g *gauge) read() *metricdata.Metric {
	now := time.Now()
	m := &metricdata.Metric{
		Descriptor: g.desc,
	}
	g.vals.Range(func(k, v interface{}) bool {
		entry := v.(gaugeEntry)
		key := k.(string)
		labelVals := g.labelValues(key)
		m.TimeSeries = append(m.TimeSeries, &metricdata.TimeSeries{
			StartTime:   now, // Gauge value is instantaneous.
			LabelValues: labelVals,
			Points: []metricdata.Point{
				entry.read(now),
			},
		})
		return true
	})
	return m
}

func (g *gauge) mapKey(labelVals []metricdata.LabelValue) string {
	vb := &tagencoding.Values{}
	for _, v := range labelVals {
		b := make([]byte, 1, len(v.Value)+1)
		if v.Present {
			b[0] = 1
			b = append(b, []byte(v.Value)...)
		}
		vb.WriteValue(b)
	}
	return string(vb.Bytes())
}

func (g *gauge) labelValues(s string) []metricdata.LabelValue {
	vals := make([]metricdata.LabelValue, 0, len(g.keys))
	vb := &tagencoding.Values{Buffer: []byte(s)}
	for range g.keys {
		v := vb.ReadValue()
		if v[0] == 0 {
			vals = append(vals, metricdata.LabelValue{})
		} else {
			vals = append(vals, metricdata.NewLabelValue(string(v[1:])))
		}
	}
	return vals
}

func (g *gauge) entryForValues(labelVals []metricdata.LabelValue, newEntry func() gaugeEntry) interface{} {
	if len(labelVals) != len(g.keys) {
		panic("must supply the same number of label values as keys used to construct this gauge")
	}
	mapKey := g.mapKey(labelVals)
	if entry, ok := g.vals.Load(mapKey); ok {
		return entry
	}
	entry, _ := g.vals.LoadOrStore(mapKey, newEntry())
	return entry
}

// Float64Gauge represents a float64 value that can go up and down.
//
// Float64Gauge maintains a float64 value for each combination of of label values
// passed to the Set or Add methods.
type Float64Gauge struct {
	g gauge
}

// Float64Entry represents a single value of the gauge corresponding to a set
// of label values.
type Float64Entry struct {
	val uint64 // needs to be uint64 for atomic access, interpret with math.Float64frombits
}

func (e *Float64Entry) read(t time.Time) metricdata.Point {
	v := math.Float64frombits(atomic.LoadUint64(&e.val))
	if v < 0 {
		v = 0
	}
	return metricdata.NewFloat64Point(t, v)
}

// GetEntry returns a gauge entry where each key for this gauge has the value
// given.
//
// The number of label values supplied must be exactly the same as the number
// of keys supplied when this gauge was created.
func (g *Float64Gauge) GetEntry(labelVals ...metricdata.LabelValue) *Float64Entry {
	return g.g.entryForValues(labelVals, func() gaugeEntry {
		return &Float64Entry{}
	}).(*Float64Entry)
}

// Set sets the gauge entry value to val.
func (e *Float64Entry) Set(val float64) {
	atomic.StoreUint64(&e.val, math.Float64bits(val))
}

// Add increments the gauge entry value by val.
func (e *Float64Entry) Add(val float64) {
	var swapped bool
	for !swapped {
		oldVal := atomic.LoadUint64(&e.val)
		newVal := math.Float64bits(math.Float64frombits(oldVal) + val)
		swapped = atomic.CompareAndSwapUint64(&e.val, oldVal, newVal)
	}
}

// Int64Gauge represents a int64 gauge value that can go up and down.
//
// Int64Gauge maintains an int64 value for each combination of label values passed to the
// Set or Add methods.
type Int64Gauge struct {
	g gauge
}

// Int64GaugeEntry represents a single value of the gauge corresponding to a set
// of label values.
type Int64GaugeEntry struct {
	val int64
}

func (e *Int64GaugeEntry) read(t time.Time) metricdata.Point {
	v := atomic.LoadInt64(&e.val)
	if v < 0 {
		v = 0.0
	}
	return metricdata.NewInt64Point(t, v)
}

// GetEntry returns a gauge entry where each key for this gauge has the value
// given.
//
// The number of label values supplied must be exactly the same as the number
// of keys supplied when this gauge was created.
func (g *Int64Gauge) GetEntry(labelVals ...metricdata.LabelValue) *Int64GaugeEntry {
	return g.g.entryForValues(labelVals, func() gaugeEntry {
		return &Int64GaugeEntry{}
	}).(*Int64GaugeEntry)
}

// Set sets the value of the gauge entry to the provided value.
func (e *Int64GaugeEntry) Set(val int64) {
	atomic.StoreInt64(&e.val, val)
}

// Add increments the current gauge entry value by val, which may be negative.
func (e *Int64GaugeEntry) Add(val int64) {
	atomic.AddInt64(&e.val, val)
}
