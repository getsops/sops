/*
Copyright 2015 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bigtable

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/option"
	btpb "google.golang.org/genproto/googleapis/bigtable/v2"
	"google.golang.org/grpc"
)

func TestPrefix(t *testing.T) {
	for _, test := range []struct {
		prefix, succ string
	}{
		{"", ""},
		{"\xff", ""}, // when used, "" means Infinity
		{"x\xff", "y"},
		{"\xfe", "\xff"},
	} {
		got := prefixSuccessor(test.prefix)
		if got != test.succ {
			t.Errorf("prefixSuccessor(%q) = %q, want %s", test.prefix, got, test.succ)
			continue
		}
		r := PrefixRange(test.prefix)
		if test.succ == "" && r.limit != "" {
			t.Errorf("PrefixRange(%q) got limit %q", test.prefix, r.limit)
		}
		if test.succ != "" && r.limit != test.succ {
			t.Errorf("PrefixRange(%q) got limit %q, want %q", test.prefix, r.limit, test.succ)
		}
	}
}

func TestApplyErrors(t *testing.T) {
	ctx := context.Background()
	table := &Table{
		c: &Client{
			project:  "P",
			instance: "I",
		},
		table: "t",
	}
	f := ColumnFilter("C")
	m := NewMutation()
	m.DeleteRow()
	// Test nested conditional mutations.
	cm := NewCondMutation(f, NewCondMutation(f, m, nil), nil)
	if err := table.Apply(ctx, "x", cm); err == nil {
		t.Error("got nil, want error")
	}
	cm = NewCondMutation(f, nil, NewCondMutation(f, m, nil))
	if err := table.Apply(ctx, "x", cm); err == nil {
		t.Error("got nil, want error")
	}
}

func TestGroupEntries(t *testing.T) {
	for _, test := range []struct {
		desc string
		in   []*entryErr
		size int
		want [][]*entryErr
	}{
		{
			desc: "one entry less than max size is one group",
			in:   []*entryErr{buildEntry(5)},
			size: 10,
			want: [][]*entryErr{{buildEntry(5)}},
		},
		{
			desc: "one entry equal to max size is one group",
			in:   []*entryErr{buildEntry(10)},
			size: 10,
			want: [][]*entryErr{{buildEntry(10)}},
		},
		{
			desc: "one entry greater than max size is one group",
			in:   []*entryErr{buildEntry(15)},
			size: 10,
			want: [][]*entryErr{{buildEntry(15)}},
		},
		{
			desc: "all entries fitting within max size are one group",
			in:   []*entryErr{buildEntry(10), buildEntry(10)},
			size: 20,
			want: [][]*entryErr{{buildEntry(10), buildEntry(10)}},
		},
		{
			desc: "entries each under max size and together over max size are grouped separately",
			in:   []*entryErr{buildEntry(10), buildEntry(10)},
			size: 15,
			want: [][]*entryErr{{buildEntry(10)}, {buildEntry(10)}},
		},
		{
			desc: "entries together over max size are grouped by max size",
			in:   []*entryErr{buildEntry(5), buildEntry(5), buildEntry(5)},
			size: 10,
			want: [][]*entryErr{{buildEntry(5), buildEntry(5)}, {buildEntry(5)}},
		},
		{
			desc: "one entry over max size and one entry under max size are two groups",
			in:   []*entryErr{buildEntry(15), buildEntry(5)},
			size: 10,
			want: [][]*entryErr{{buildEntry(15)}, {buildEntry(5)}},
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			if got, want := groupEntries(test.in, test.size), test.want; !cmp.Equal(mutationCounts(got), mutationCounts(want)) {
				t.Fatalf("[%s] want = %v, got = %v", test.desc, mutationCounts(want), mutationCounts(got))
			}
		})
	}
}

func buildEntry(numMutations int) *entryErr {
	var muts []*btpb.Mutation
	for i := 0; i < numMutations; i++ {
		muts = append(muts, &btpb.Mutation{})
	}
	return &entryErr{Entry: &btpb.MutateRowsRequest_Entry{Mutations: muts}}
}

func mutationCounts(batched [][]*entryErr) []int {
	var res []int
	for _, entries := range batched {
		var count int
		for _, e := range entries {
			count += len(e.Entry.Mutations)
		}
		res = append(res, count)
	}
	return res
}

type requestCountingInterceptor struct {
	grpc.ClientStream
	requestCallback func()
}

func (i *requestCountingInterceptor) SendMsg(m interface{}) error {
	i.requestCallback()
	return i.ClientStream.SendMsg(m)
}

func (i *requestCountingInterceptor) RecvMsg(m interface{}) error {
	return i.ClientStream.RecvMsg(m)
}

func requestCallback(callback func()) func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		return &requestCountingInterceptor{
			ClientStream:    clientStream,
			requestCallback: callback,
		}, err
	}
}

// TestReadRowsInvalidRowSet verifies that the client doesn't send ReadRows() requests with invalid RowSets.
func TestReadRowsInvalidRowSet(t *testing.T) {
	testEnv, err := NewEmulatedEnv(IntegrationTestConfig{})
	if err != nil {
		t.Fatalf("NewEmulatedEnv failed: %v", err)
	}
	var requestCount int
	incrementRequestCount := func() { requestCount++ }
	conn, err := grpc.Dial(testEnv.server.Addr, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(100<<20), grpc.MaxCallRecvMsgSize(100<<20)),
		grpc.WithStreamInterceptor(requestCallback(incrementRequestCount)),
	)
	if err != nil {
		t.Fatalf("grpc.Dial failed: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	adminClient, err := NewAdminClient(ctx, testEnv.config.Project, testEnv.config.Instance, option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer adminClient.Close()
	if err := adminClient.CreateTable(ctx, testEnv.config.Table); err != nil {
		t.Fatalf("CreateTable(%v) failed: %v", testEnv.config.Table, err)
	}
	client, err := NewClient(ctx, testEnv.config.Project, testEnv.config.Instance, option.WithGRPCConn(conn))
	if err != nil {
		t.Fatalf("NewClient failed: %v", err)
	}
	defer client.Close()
	table := client.Open(testEnv.config.Table)
	tests := []struct {
		rr    RowSet
		valid bool
	}{
		{
			rr:    RowRange{},
			valid: true,
		},
		{
			rr:    RowRange{start: "b"},
			valid: true,
		},
		{
			rr:    RowRange{start: "b", limit: "c"},
			valid: true,
		},
		{
			rr:    RowRange{start: "b", limit: "a"},
			valid: false,
		},
		{
			rr:    RowList{"a"},
			valid: true,
		},
		{
			rr:    RowList{},
			valid: false,
		},
	}
	for _, test := range tests {
		requestCount = 0
		err = table.ReadRows(ctx, test.rr, func(r Row) bool { return true })
		if err != nil {
			t.Fatalf("ReadRows(%v) failed: %v", test.rr, err)
		}
		requestValid := requestCount != 0
		if requestValid != test.valid {
			t.Errorf("%s: got %v, want %v", test.rr, requestValid, test.valid)
		}
	}
}
