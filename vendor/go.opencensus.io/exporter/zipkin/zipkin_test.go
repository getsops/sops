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

package zipkin

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/openzipkin/zipkin-go/model"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	"go.opencensus.io/trace"
)

type roundTripper func(*http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

func TestExport(t *testing.T) {
	// Since Zipkin reports in microsecond resolution let's round our Timestamp,
	// so when deserializing Zipkin data in this test we can properly compare.
	now := time.Now().Round(time.Microsecond)
	tests := []struct {
		span *trace.SpanData
		want model.SpanModel
	}{
		{
			span: &trace.SpanData{
				SpanContext: trace.SpanContext{
					TraceID:      trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
					SpanID:       trace.SpanID{17, 18, 19, 20, 21, 22, 23, 24},
					TraceOptions: 1,
				},
				Name:      "name",
				SpanKind:  trace.SpanKindClient,
				StartTime: now,
				EndTime:   now.Add(24 * time.Hour),
				Attributes: map[string]interface{}{
					"stringkey": "value",
					"intkey":    int64(42),
					"boolkey1":  true,
					"boolkey2":  false,
				},
				MessageEvents: []trace.MessageEvent{
					{
						Time:                 now,
						EventType:            trace.MessageEventTypeSent,
						MessageID:            12,
						UncompressedByteSize: 99,
						CompressedByteSize:   98,
					},
				},
				Annotations: []trace.Annotation{
					{
						Time:    now,
						Message: "Annotation",
						Attributes: map[string]interface{}{
							"stringkey": "value",
							"intkey":    int64(42),
							"boolkey1":  true,
							"boolkey2":  false,
						},
					},
				},
				Status: trace.Status{
					Code:    3,
					Message: "error",
				},
			},
			want: model.SpanModel{
				SpanContext: model.SpanContext{
					TraceID: model.TraceID{
						High: 0x0102030405060708,
						Low:  0x090a0b0c0d0e0f10,
					},
					ID:      0x1112131415161718,
					Sampled: &sampledTrue,
				},
				Name:      "name",
				Kind:      model.Client,
				Timestamp: now,
				Duration:  24 * time.Hour,
				Shared:    false,
				Annotations: []model.Annotation{
					{
						Timestamp: now,
						Value:     "Annotation",
					},
					{
						Timestamp: now,
						Value:     "SENT",
					},
				},
				Tags: map[string]string{
					"stringkey":                     "value",
					"intkey":                        "42",
					"boolkey1":                      "true",
					"boolkey2":                      "false",
					"error":                         "INVALID_ARGUMENT",
					"opencensus.status_description": "error",
				},
			},
		},
		{
			span: &trace.SpanData{
				SpanContext: trace.SpanContext{
					TraceID:      trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
					SpanID:       trace.SpanID{17, 18, 19, 20, 21, 22, 23, 24},
					TraceOptions: 1,
				},
				Name:      "name",
				StartTime: now,
				EndTime:   now.Add(24 * time.Hour),
			},
			want: model.SpanModel{
				SpanContext: model.SpanContext{
					TraceID: model.TraceID{
						High: 0x0102030405060708,
						Low:  0x090a0b0c0d0e0f10,
					},
					ID:      0x1112131415161718,
					Sampled: &sampledTrue,
				},
				Name:      "name",
				Timestamp: now,
				Duration:  24 * time.Hour,
				Shared:    false,
			},
		},
		{
			span: &trace.SpanData{
				SpanContext: trace.SpanContext{
					TraceID:      trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
					SpanID:       trace.SpanID{17, 18, 19, 20, 21, 22, 23, 24},
					TraceOptions: 1,
				},
				Name:      "name",
				StartTime: now,
				EndTime:   now.Add(24 * time.Hour),
				Status: trace.Status{
					Code:    0,
					Message: "there is no cause for alarm",
				},
			},
			want: model.SpanModel{
				SpanContext: model.SpanContext{
					TraceID: model.TraceID{
						High: 0x0102030405060708,
						Low:  0x090a0b0c0d0e0f10,
					},
					ID:      0x1112131415161718,
					Sampled: &sampledTrue,
				},
				Name:      "name",
				Timestamp: now,
				Duration:  24 * time.Hour,
				Shared:    false,
				Tags: map[string]string{
					"opencensus.status_description": "there is no cause for alarm",
				},
			},
		},
		{
			span: &trace.SpanData{
				SpanContext: trace.SpanContext{
					TraceID:      trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
					SpanID:       trace.SpanID{17, 18, 19, 20, 21, 22, 23, 24},
					TraceOptions: 1,
				},
				Name:      "name",
				StartTime: now,
				EndTime:   now.Add(24 * time.Hour),
				Status: trace.Status{
					Code: 1234,
				},
			},
			want: model.SpanModel{
				SpanContext: model.SpanContext{
					TraceID: model.TraceID{
						High: 0x0102030405060708,
						Low:  0x090a0b0c0d0e0f10,
					},
					ID:      0x1112131415161718,
					Sampled: &sampledTrue,
				},
				Name:      "name",
				Timestamp: now,
				Duration:  24 * time.Hour,
				Shared:    false,
				Tags: map[string]string{
					"error": "error code 1234",
				},
			},
		},
	}
	for _, tt := range tests {
		got := zipkinSpan(tt.span, nil)
		if len(got.Annotations) != len(tt.want.Annotations) {
			t.Fatalf("zipkinSpan: got %d annotations in span, want %d", len(got.Annotations), len(tt.want.Annotations))
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("zipkinSpan:\n\tgot  %#v\n\twant %#v", got, tt.want)
		}
	}
	for _, tt := range tests {
		ch := make(chan []byte)
		client := http.Client{
			Transport: roundTripper(func(req *http.Request) (*http.Response, error) {
				body, _ := ioutil.ReadAll(req.Body)
				ch <- body
				return &http.Response{StatusCode: 200, Body: ioutil.NopCloser(strings.NewReader(""))}, nil
			}),
		}
		reporter := httpreporter.NewReporter("foo", httpreporter.Client(&client), httpreporter.BatchInterval(time.Millisecond))
		exporter := NewExporter(reporter, nil)
		exporter.ExportSpan(tt.span)
		var data []byte
		select {
		case data = <-ch:
		case <-time.After(2 * time.Second):
			t.Fatalf("span was not exported")
		}
		var spans []model.SpanModel
		json.Unmarshal(data, &spans)
		if len(spans) != 1 {
			t.Fatalf("Export: got %d spans, want 1", len(spans))
		}
		got := spans[0]
		got.SpanContext.Sampled = &sampledTrue // Sampled is not set when the span is reported.
		if len(got.Annotations) != len(tt.want.Annotations) {
			t.Fatalf("Export: got %d annotations in span, want %d", len(got.Annotations), len(tt.want.Annotations))
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Export:\n\tgot  %#v\n\twant %#v", got, tt.want)
		}
	}
}
