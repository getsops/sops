// Copyright 2017, OpenCensus Authors
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

package prometheus

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/prometheus/client_golang/prometheus"
)

func newView(measureName string, agg *view.Aggregation) *view.View {
	m := stats.Int64(measureName, "bytes", stats.UnitBytes)
	return &view.View{
		Name:        "foo",
		Description: "bar",
		Measure:     m,
		Aggregation: agg,
	}
}

func TestOnlyCumulativeWindowSupported(t *testing.T) {
	// See Issue https://github.com/census-instrumentation/opencensus-go/issues/214.
	count1 := &view.CountData{Value: 1}
	lastValue1 := &view.LastValueData{Value: 56.7}
	tests := []struct {
		vds  *view.Data
		want int
	}{
		0: {
			vds: &view.Data{
				View: newView("TestOnlyCumulativeWindowSupported/m1", view.Count()),
			},
			want: 0, // no rows present
		},
		1: {
			vds: &view.Data{
				View: newView("TestOnlyCumulativeWindowSupported/m2", view.Count()),
				Rows: []*view.Row{
					{Data: count1},
				},
			},
			want: 1,
		},
		2: {
			vds: &view.Data{
				View: newView("TestOnlyCumulativeWindowSupported/m3", view.LastValue()),
				Rows: []*view.Row{
					{Data: lastValue1},
				},
			},
			want: 1,
		},
	}

	for i, tt := range tests {
		reg := prometheus.NewRegistry()
		collector := newCollector(Options{}, reg)
		collector.addViewData(tt.vds)
		mm, err := reg.Gather()
		if err != nil {
			t.Errorf("#%d: Gather err: %v", i, err)
		}
		reg.Unregister(collector)
		if got, want := len(mm), tt.want; got != want {
			t.Errorf("#%d: got nil %v want nil %v", i, got, want)
		}
	}
}

func TestCollectNonRacy(t *testing.T) {
	// Despite enforcing the singleton, for this case we
	// need an exporter hence won't be using NewExporter.
	exp, err := NewExporter(Options{})
	if err != nil {
		t.Fatalf("NewExporter: %v", err)
	}
	collector := exp.c

	// Synchronize and make sure every goroutine has terminated before we exit
	var waiter sync.WaitGroup
	waiter.Add(3)
	defer waiter.Wait()

	doneCh := make(chan bool)
	// 1. Viewdata write routine at 700ns
	go func() {
		defer waiter.Done()
		tick := time.NewTicker(700 * time.Nanosecond)
		defer tick.Stop()

		defer func() {
			close(doneCh)
		}()

		for i := 0; i < 1e3; i++ {
			count1 := &view.CountData{Value: 1}
			vds := []*view.Data{
				{View: newView(fmt.Sprintf("TestCollectNonRacy/m2-%d", i), view.Count()), Rows: []*view.Row{{Data: count1}}},
			}
			for _, v := range vds {
				exp.ExportView(v)
			}
			<-tick.C
		}
	}()

	inMetricsChan := make(chan prometheus.Metric, 1000)
	// 2. Simulating the Prometheus metrics consumption routine running at 900ns
	go func() {
		defer waiter.Done()
		tick := time.NewTicker(900 * time.Nanosecond)
		defer tick.Stop()

		for {
			select {
			case <-doneCh:
				return
			case <-inMetricsChan:
			}
		}
	}()

	// 3. Collect/Read routine at 800ns
	go func() {
		defer waiter.Done()
		tick := time.NewTicker(800 * time.Nanosecond)
		defer tick.Stop()

		for {
			select {
			case <-doneCh:
				return
			case <-tick.C:
				// Perform some collection here
				collector.Collect(inMetricsChan)
			}
		}
	}()
}

type mSlice []*stats.Int64Measure

func (measures *mSlice) createAndAppend(name, desc, unit string) {
	m := stats.Int64(name, desc, unit)
	*measures = append(*measures, m)
}

type vCreator []*view.View

func (vc *vCreator) createAndAppend(name, description string, keys []tag.Key, measure stats.Measure, agg *view.Aggregation) {
	v := &view.View{
		Name:        name,
		Description: description,
		TagKeys:     keys,
		Measure:     measure,
		Aggregation: agg,
	}
	*vc = append(*vc, v)
}

func TestMetricsEndpointOutput(t *testing.T) {
	exporter, err := NewExporter(Options{})
	if err != nil {
		t.Fatalf("failed to create prometheus exporter: %v", err)
	}
	view.RegisterExporter(exporter)

	names := []string{"foo", "bar", "baz"}

	var measures mSlice
	for _, name := range names {
		measures.createAndAppend("tests/"+name, name, "")
	}

	var vc vCreator
	for _, m := range measures {
		vc.createAndAppend(m.Name(), m.Description(), nil, m, view.Count())
	}

	if err := view.Register(vc...); err != nil {
		t.Fatalf("failed to create views: %v", err)
	}
	defer view.Unregister(vc...)

	view.SetReportingPeriod(time.Millisecond)

	for _, m := range measures {
		stats.Record(context.Background(), m.M(1))
	}

	srv := httptest.NewServer(exporter)
	defer srv.Close()

	var i int
	var output string
	for {
		time.Sleep(10 * time.Millisecond)
		if i == 1000 {
			t.Fatal("no output at /metrics (10s wait)")
		}
		i++

		resp, err := http.Get(srv.URL)
		if err != nil {
			t.Fatalf("failed to get /metrics: %v", err)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read body: %v", err)
		}
		resp.Body.Close()

		output = string(body)
		if output != "" {
			break
		}
	}

	if strings.Contains(output, "collected before with the same name and label values") {
		t.Fatal("metric name and labels being duplicated but must be unique")
	}

	if strings.Contains(output, "error(s) occurred") {
		t.Fatal("error reported by prometheus registry")
	}

	for _, name := range names {
		if !strings.Contains(output, "tests_"+name+" 1") {
			t.Fatalf("measurement missing in output: %v", name)
		}
	}
}

func TestCumulativenessFromHistograms(t *testing.T) {
	exporter, err := NewExporter(Options{})
	if err != nil {
		t.Fatalf("failed to create prometheus exporter: %v", err)
	}
	view.RegisterExporter(exporter)
	reportPeriod := time.Millisecond
	view.SetReportingPeriod(reportPeriod)

	m := stats.Float64("tests/bills", "payments by denomination", stats.UnitDimensionless)
	v := &view.View{
		Name:        "cash/register",
		Description: "this is a test",
		Measure:     m,

		// Intentionally used repeated elements in the ascending distribution.
		// to ensure duplicate distribution items are handles.
		Aggregation: view.Distribution(1, 5, 5, 5, 5, 10, 20, 50, 100, 250),
	}

	if err := view.Register(v); err != nil {
		t.Fatalf("Register error: %v", err)
	}
	defer view.Unregister(v)

	// Give the reporter ample time to process registration
	<-time.After(10 * reportPeriod)

	values := []float64{0.25, 245.67, 12, 1.45, 199.9, 7.69, 187.12}
	// We want the results that look like this:
	// 1:   [0.25]      		| 1 + prev(i) = 1 + 0 = 1
	// 5:   [1.45]			| 1 + prev(i) = 1 + 1 = 2
	// 10:	[]			| 1 + prev(i) = 1 + 2 = 3
	// 20:  [12]			| 1 + prev(i) = 1 + 3 = 4
	// 50:  []			| 0 + prev(i) = 0 + 4 = 4
	// 100: []			| 0 + prev(i) = 0 + 4 = 4
	// 250: [187.12, 199.9, 245.67]	| 3 + prev(i) = 3 + 4 = 7
	wantLines := []string{
		`cash_register_bucket{le="1"} 1`,
		`cash_register_bucket{le="5"} 2`,
		`cash_register_bucket{le="10"} 3`,
		`cash_register_bucket{le="20"} 4`,
		`cash_register_bucket{le="50"} 4`,
		`cash_register_bucket{le="100"} 4`,
		`cash_register_bucket{le="250"} 7`,
		`cash_register_bucket{le="+Inf"} 7`,
		`cash_register_sum 654.0799999999999`, // Summation of the input values
		`cash_register_count 7`,
	}

	ctx := context.Background()
	ms := make([]stats.Measurement, 0, len(values))
	for _, value := range values {
		mx := m.M(value)
		ms = append(ms, mx)
	}
	stats.Record(ctx, ms...)

	// Give the recorder ample time to process recording
	<-time.After(10 * reportPeriod)

	cst := httptest.NewServer(exporter)
	defer cst.Close()
	res, err := http.Get(cst.URL)
	if err != nil {
		t.Fatalf("http.Get error: %v", err)
	}
	blob, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Read body error: %v", err)
	}
	str := strings.Trim(string(blob), "\n")
	lines := strings.Split(str, "\n")
	nonComments := make([]string, 0, len(lines))
	for _, line := range lines {
		if !strings.Contains(line, "#") {
			nonComments = append(nonComments, line)
		}
	}

	got := strings.Join(nonComments, "\n")
	want := strings.Join(wantLines, "\n")
	if got != want {
		t.Fatalf("\ngot:\n%s\n\nwant:\n%s\n", got, want)
	}
}
