/*
Copyright 2019 Google LLC

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
	"bytes"
	"context"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	pb "cloud.google.com/go/bigtable/internal/conformance"
	"cloud.google.com/go/bigtable/internal/mockserver"
	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/api/option"
	btpb "google.golang.org/genproto/googleapis/bigtable/v2"
	"google.golang.org/grpc"
)

func TestConformance(t *testing.T) {
	ctx := context.Background()

	dir := "internal/conformance/testdata"
	files, err := filepath.Glob(dir + "/*.json")
	if err != nil {
		t.Fatal(err)
	}

	srv, err := mockserver.NewServer("localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer srv.Close()

	conn, err := grpc.Dial(srv.Addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c, err := NewClient(ctx, "some-project", "some-instance", option.WithGRPCConn(conn))
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		inBytes, err := ioutil.ReadFile(f)
		if err != nil {
			t.Fatalf("%s: %v", f, err)
		}

		var tf pb.TestFile
		if err := jsonpb.Unmarshal(bytes.NewReader(inBytes), &tf); err != nil {
			t.Fatalf("unmarshalling %s: %v", f, err)
		}

		for _, tc := range tf.GetReadRowsTests() {
			t.Run(tc.Description, func(t *testing.T) {
				runReadRowsTest(ctx, t, tc, c, srv)
			})
		}
	}
}

func runReadRowsTest(ctx context.Context, t *testing.T, tc *pb.ReadRowsTest, c *Client, srv *mockserver.Server) {
	srv.ReadRowsFn = func(req *btpb.ReadRowsRequest, something btpb.Bigtable_ReadRowsServer) error {
		something.Send(&btpb.ReadRowsResponse{
			Chunks: tc.GetChunks(),
		})

		return nil
	}

	var resIndex int

	// We perform a SingleRow here, but that arg is basically nonsense since
	// the server is hard-coded to return a specific response. As in, we could
	// pass RowRange, ListRows, etc and the result would all be the same.
	err := c.Open("some-table").ReadRows(ctx, SingleRow("some-row"), func(r Row) bool {
		type rowElem struct {
			family    string
			readItems []ReadItem
		}

		// Row comes in as a map, which has undefined iteration order. So, we
		// first stick it in a slice, then sort that slice by family (the
		// results appear ordered as such), then we're ready to use it.
		var byFamily []rowElem
		for family, items := range r {
			byFamily = append(byFamily, rowElem{family: family, readItems: items})
		}
		sort.Slice(byFamily, func(i, j int) bool {
			return strings.Compare(byFamily[i].family, byFamily[j].family) < 0
		})

		for _, row := range byFamily {
			family := row.family
			items := row.readItems
			for _, item := range items {
				want := tc.GetResults()[resIndex]

				if got, want := string(item.Value), want.GetValue(); got != want {
					t.Fatalf("got %s, want %s", got, want)
				}

				if got, want := family, want.GetFamilyName(); got != want {
					t.Fatalf("got %s, want %s", got, want)
				}

				gotMicros := item.Timestamp.Time().UnixNano() / int64(time.Microsecond)
				if got, want := gotMicros, want.GetTimestampMicros(); got != want {
					t.Fatalf("got %d, want %d", got, want)
				}

				if got, want := item.Column, want.GetFamilyName()+":"+want.GetQualifier(); got != want {
					t.Fatalf("got %s, want %s", got, want)
				}

				// TODO: labels do not appear to be accessible. If they ever do become
				// accessible, we should assert on want.GetLabels().

				resIndex++
			}
		}
		return true
	})

	wantNumResults := len(tc.GetResults())

	if wantNumResults == 0 {
		return
	}

	if tc.GetResults()[wantNumResults-1].GetError() {
		// Last expected result is an error, which means we wouldn't
		// count it with gotRowIndex.
		wantNumResults--

		if err == nil {
			t.Fatal("expected err, got nil")
		}
	} else if err != nil {
		t.Fatal(err)
	}

	if got, want := resIndex, wantNumResults; got != want {
		t.Fatalf("got %d results, want %d", got, want)
	}
}
