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
	"go.opencensus.io/metric/metricdata"
	"log"
	"time"
)

// Registry creates and manages a set of gauges.
// External synchronization is required if you want to add gauges to the same
// registry from multiple goroutines.
type Registry struct {
	gauges map[string]*gauge
}

// NewRegistry initializes a new Registry.
func NewRegistry() *Registry {
	return &Registry{
		gauges: make(map[string]*gauge),
	}
}

// AddFloat64Gauge creates and adds a new float64-valued gauge to this registry.
func (r *Registry) AddFloat64Gauge(name, description string, unit metricdata.Unit, labelKeys ...string) *Float64Gauge {
	f := &Float64Gauge{
		g: gauge{
			isFloat: true,
		},
	}
	r.initGauge(&f.g, labelKeys, name, description, unit)
	return f
}

// AddInt64Gauge creates and adds a new int64-valued gauge to this registry.
func (r *Registry) AddInt64Gauge(name, description string, unit metricdata.Unit, labelKeys ...string) *Int64Gauge {
	i := &Int64Gauge{}
	r.initGauge(&i.g, labelKeys, name, description, unit)
	return i
}

func (r *Registry) initGauge(g *gauge, labelKeys []string, name string, description string, unit metricdata.Unit) *gauge {
	existing, ok := r.gauges[name]
	if ok {
		if existing.isFloat != g.isFloat {
			log.Panicf("Gauge with name %s already exists with a different type", name)
		}
	}
	g.keys = labelKeys
	g.start = time.Now()
	g.desc = metricdata.Descriptor{
		Name:        name,
		Description: description,
		Unit:        unit,
		LabelKeys:   labelKeys,
	}
	r.gauges[name] = g
	return g
}

// ReadAll reads all gauges in this registry and returns their values as metrics.
func (r *Registry) ReadAll() []*metricdata.Metric {
	ms := make([]*metricdata.Metric, 0, len(r.gauges))
	for _, g := range r.gauges {
		ms = append(ms, g.read())
	}
	return ms
}
