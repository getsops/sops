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
//

package ocgrpc

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"go.opencensus.io/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"context"

	"go.opencensus.io/metric/metricdata"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"google.golang.org/grpc/stats"
)

func TestClientDefaultCollections(t *testing.T) {
	k1, _ := tag.NewKey("k1")
	k2, _ := tag.NewKey("k2")

	type tagPair struct {
		k tag.Key
		v string
	}

	type wantData struct {
		v    func() *view.View
		rows []*view.Row
	}
	type rpc struct {
		tags        []tagPair
		tagInfo     *stats.RPCTagInfo
		inPayloads  []*stats.InPayload
		outPayloads []*stats.OutPayload
		end         *stats.End
	}

	type testCase struct {
		label string
		rpcs  []*rpc
		wants []*wantData
	}
	tcs := []testCase{
		{
			label: "1",
			rpcs: []*rpc{
				{
					[]tagPair{{k1, "v1"}},
					&stats.RPCTagInfo{FullMethodName: "/package.service/method"},
					[]*stats.InPayload{
						{Length: 10},
					},
					[]*stats.OutPayload{
						{Length: 10},
					},
					&stats.End{Error: nil},
				},
			},
			wants: []*wantData{
				{
					func() *view.View { return ClientSentMessagesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 1, 1, 1, 0),
						},
					},
				},
				{
					func() *view.View { return ClientReceivedMessagesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 1, 1, 1, 0),
						},
					},
				},
				{
					func() *view.View { return ClientSentBytesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 10, 10, 10, 0),
						},
					},
				},
				{
					func() *view.View { return ClientReceivedBytesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 1, 10, 10, 10, 0),
						},
					},
				},
			},
		},
		{
			label: "2",
			rpcs: []*rpc{
				{
					[]tagPair{{k1, "v1"}},
					&stats.RPCTagInfo{FullMethodName: "/package.service/method"},
					[]*stats.InPayload{
						{Length: 10},
					},
					[]*stats.OutPayload{
						{Length: 10},
						{Length: 10},
						{Length: 10},
					},
					&stats.End{Error: nil},
				},
				{
					[]tagPair{{k1, "v11"}},
					&stats.RPCTagInfo{FullMethodName: "/package.service/method"},
					[]*stats.InPayload{
						{Length: 10},
						{Length: 10},
					},
					[]*stats.OutPayload{
						{Length: 10},
						{Length: 10},
					},
					&stats.End{Error: status.Error(codes.Canceled, "canceled")},
				},
			},
			wants: []*wantData{
				{
					func() *view.View { return ClientSentMessagesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 2, 2, 3, 2.5, 0.5),
						},
					},
				},
				{
					func() *view.View { return ClientReceivedMessagesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 2, 1, 2, 1.5, 0.5),
						},
					},
				},
			},
		},
		{
			label: "3",
			rpcs: []*rpc{
				{
					[]tagPair{{k1, "v1"}},
					&stats.RPCTagInfo{FullMethodName: "/package.service/method"},
					[]*stats.InPayload{
						{Length: 1},
					},
					[]*stats.OutPayload{
						{Length: 1},
						{Length: 1024},
						{Length: 65536},
					},
					&stats.End{Error: nil},
				},
				{
					[]tagPair{{k1, "v1"}, {k2, "v2"}},
					&stats.RPCTagInfo{FullMethodName: "/package.service/method"},
					[]*stats.InPayload{
						{Length: 1024},
					},
					[]*stats.OutPayload{
						{Length: 4096},
						{Length: 16384},
					},
					&stats.End{Error: status.Error(codes.Canceled, "canceled")},
				},
				{
					[]tagPair{{k1, "v11"}, {k2, "v22"}},
					&stats.RPCTagInfo{FullMethodName: "/package.service/method"},
					[]*stats.InPayload{
						{Length: 2048},
						{Length: 16384},
					},
					[]*stats.OutPayload{
						{Length: 2048},
						{Length: 4096},
						{Length: 16384},
					},
					&stats.End{Error: status.Error(codes.Aborted, "aborted")},
				},
			},
			wants: []*wantData{
				{
					func() *view.View { return ClientSentMessagesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 3, 2, 3, 2.666666666, 0.333333333*2),
						},
					},
				},
				{
					func() *view.View { return ClientReceivedMessagesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 2, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 3, 1, 2, 1.333333333, 0.333333333*2),
						},
					},
				},
				{
					func() *view.View { return ClientSentBytesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{0, 0, 0, 0, 2 /*16384*/, 1 /*65536*/, 0, 0, 0, 0, 0, 0, 0, 0}, 3, 20480, 66561, 36523, 1.355519318e+09),
						},
					},
				},
				{
					func() *view.View { return ClientReceivedBytesPerRPCView },
					[]*view.Row{
						{
							Tags: []tag.Tag{
								{Key: KeyClientMethod, Value: "package.service/method"},
							},
							Data: newDistributionData([]int64{1, 1, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0}, 3, 1, 18432, 6485.666667, 2.1459558466666666e+08),
						},
					},
				},
			},
		},
	}

	views := []*view.View{
		ClientSentBytesPerRPCView,
		ClientReceivedBytesPerRPCView,
		ClientRoundtripLatencyView,
		ClientCompletedRPCsView,
		ClientSentMessagesPerRPCView,
		ClientReceivedMessagesPerRPCView,
	}

	for _, tc := range tcs {
		// Register views.
		if err := view.Register(views...); err != nil {
			t.Error(err)
		}

		h := &ClientHandler{}
		h.StartOptions.Sampler = trace.NeverSample()
		for _, rpc := range tc.rpcs {
			var mods []tag.Mutator
			for _, t := range rpc.tags {
				mods = append(mods, tag.Upsert(t.k, t.v))
			}
			ctx, err := tag.New(context.Background(), mods...)
			if err != nil {
				t.Errorf("%q: NewMap = %v", tc.label, err)
			}
			encoded := tag.Encode(tag.FromContext(ctx))
			ctx = stats.SetTags(context.Background(), encoded)
			ctx = h.TagRPC(ctx, rpc.tagInfo)
			for _, out := range rpc.outPayloads {
				out.Client = true
				h.HandleRPC(ctx, out)
			}
			for _, in := range rpc.inPayloads {
				in.Client = true
				h.HandleRPC(ctx, in)
			}
			rpc.end.Client = true
			h.HandleRPC(ctx, rpc.end)
		}

		for _, wantData := range tc.wants {
			gotRows, err := view.RetrieveData(wantData.v().Name)
			if err != nil {
				t.Errorf("%q: RetrieveData(%q) = %v", tc.label, wantData.v().Name, err)
				continue
			}

			for _, gotRow := range gotRows {
				if !containsRow(wantData.rows, gotRow) {
					t.Errorf("%q: unwanted row for view %q = %v", tc.label, wantData.v().Name, gotRow)
					break
				}
			}

			for _, wantRow := range wantData.rows {
				if !containsRow(gotRows, wantRow) {
					t.Errorf("%q: row missing for view %q; want %v", tc.label, wantData.v().Name, wantRow)
					break
				}
			}
		}

		// Unregister views to cleanup.
		view.Unregister(views...)
	}
}

func TestClientRecordExemplar(t *testing.T) {
	key, _ := tag.NewKey("test_key")
	tagInfo := &stats.RPCTagInfo{FullMethodName: "/package.service/method"}
	out := &stats.OutPayload{Length: 2000}
	end := &stats.End{Error: nil}

	if err := view.Register(ClientSentBytesPerRPCView); err != nil {
		t.Error(err)
	}
	h := &ClientHandler{}
	h.StartOptions.Sampler = trace.AlwaysSample()
	ctx, err := tag.New(context.Background(), tag.Upsert(key, "test_val"))
	if err != nil {
		t.Error(err)
	}
	encoded := tag.Encode(tag.FromContext(ctx))
	ctx = stats.SetTags(context.Background(), encoded)
	ctx = h.TagRPC(ctx, tagInfo)

	out.Client = true
	h.HandleRPC(ctx, out)
	end.Client = true
	h.HandleRPC(ctx, end)

	span := trace.FromContext(ctx)
	if span == nil {
		t.Fatal("expected non-nil span, got nil")
	}
	if !span.IsRecordingEvents() {
		t.Errorf("span should be sampled")
	}
	attachments := map[string]interface{}{metricdata.AttachmentKeySpanContext: span.SpanContext()}
	wantExemplar := &metricdata.Exemplar{Value: 2000, Attachments: attachments}

	rows, err := view.RetrieveData(ClientSentBytesPerRPCView.Name)
	if err != nil {
		t.Fatal("Error RetrieveData ", err)
	}
	if len(rows) == 0 {
		t.Fatal("No data was recorded.")
	}
	data := rows[0].Data
	dis, ok := data.(*view.DistributionData)
	if !ok {
		t.Fatal("want DistributionData, got ", data)
	}
	// Only recorded value is 2000, which falls into the second bucket (1024, 2048].
	wantBuckets := []int64{0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	if !reflect.DeepEqual(dis.CountPerBucket, wantBuckets) {
		t.Errorf("want buckets %v, got %v", wantBuckets, dis.CountPerBucket)
	}
	for i, e := range dis.ExemplarsPerBucket {
		// Only the second bucket should have an exemplar.
		if i == 1 {
			if diff := cmpExemplar(e, wantExemplar); diff != "" {
				t.Fatalf("Unexpected Exemplar -got +want: %s", diff)
			}
		} else if e != nil {
			t.Errorf("want nil exemplar, got %v", e)
		}
	}

	// Unregister views to cleanup.
	view.Unregister(ClientSentBytesPerRPCView)
}

// containsRow returns true if rows contain r.
func containsRow(rows []*view.Row, r *view.Row) bool {
	for _, x := range rows {
		if r.Equal(x) {
			return true
		}
	}
	return false
}

// Compare exemplars while ignoring exemplar timestamp, since timestamp is non-deterministic.
func cmpExemplar(got, want *metricdata.Exemplar) string {
	return cmp.Diff(got, want, cmpopts.IgnoreFields(metricdata.Exemplar{}, "Timestamp"), cmpopts.IgnoreUnexported(metricdata.Exemplar{}))
}
