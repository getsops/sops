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

package bttest

import (
	"testing"

	btpb "google.golang.org/genproto/googleapis/bigtable/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMessageOnInvalidKeyRange(t *testing.T) {
	tests := []struct {
		startOpen, startClosed, endOpen, endClosed string
		wantErr                                    string
	}{
		{
			startOpen: "A", startClosed: "A",
			wantErr: "both start_key_closed and start_key_open cannot be set",
		},
		{
			endOpen: "Z", endClosed: "Z",
			wantErr: "both end_key_closed and end_key_open cannot be set",
		},
		{
			startClosed: "Z", endClosed: "A",
			wantErr: "start_key_closed must be less than end_key_closed",
		},
		{
			startClosed: "Z", endOpen: "A",
			wantErr: "start_key_closed must be less than end_key_open",
		},

		{
			startOpen: "Z", endClosed: "A",
			wantErr: "start_key_open must be less than end_key_closed",
		},
		{
			startOpen: "Z", endOpen: "A",
			wantErr: "start_key_open must be less than end_key_open",
		},

		// All values not set.
		{
			startClosed: "", startOpen: "", endClosed: "", endOpen: "", wantErr: "",
		},

		// Indefinite ranges on each side.
		{
			startClosed: "A", startOpen: "", endClosed: "", endOpen: "", wantErr: "",
		},
		{
			startClosed: "", startOpen: "A", endClosed: "", endOpen: "", wantErr: "",
		},
		{
			startClosed: "", startOpen: "", endClosed: "A", endOpen: "", wantErr: "",
		},
		{
			startClosed: "", startOpen: "", endClosed: "", endOpen: "A", wantErr: "",
		},

		// startClosed, endClosed ranges properly set.
		{
			startClosed: "A", startOpen: "", endClosed: "Z", endOpen: "", wantErr: "",
		},
		// startClosed, endOpen range.
		{
			startClosed: "A", startOpen: "", endClosed: "", endOpen: "Z", wantErr: "",
		},
		// startOpen, endClosed ranges properly set.
		{
			startClosed: "", startOpen: "A", endClosed: "Z", endOpen: "", wantErr: "",
		},
		// startOpen, endOpen range.
		{
			startClosed: "", startOpen: "A", endClosed: "", endOpen: "Z", wantErr: "",
		},
	}

	for i, tt := range tests {
		gotErr := messageOnInvalidKeyRanges(
			[]byte(tt.startClosed),
			[]byte(tt.startOpen),
			[]byte(tt.endClosed),
			[]byte(tt.endOpen))

		if gotErr != tt.wantErr {
			t.Errorf("#%d. Error mismatch\nGot:  %q\nGiven:\n%+v", i, gotErr, tt)
		}
	}
}

// This test is just to ensure that the emulator always sends back
// an RPC error with a status code just like Google Bigtable does.
func TestValidateReadRowsRequestSendsRPCError(t *testing.T) {
	tableName := "foo.org/bar"
	// Minimal server to reproduce failures.
	srv := &server{
		tables: map[string]*table{tableName: new(table)},
	}

	badValues := []struct {
		startKeyClosed, startKeyOpen, endKeyClosed, endKeyOpen string
	}{
		{startKeyClosed: "Z", endKeyClosed: "A"},
		{startKeyClosed: "Z", endKeyOpen: "A"},
		{startKeyOpen: "Z", endKeyClosed: "A"},
		{startKeyOpen: "Z", endKeyOpen: "A"},
	}

	for i, tt := range badValues {
		badReq := &btpb.ReadRowsRequest{
			TableName: tableName,
			Rows: &btpb.RowSet{
				RowRanges: []*btpb.RowRange{
					{
						StartKey: &btpb.RowRange_StartKeyClosed{
							StartKeyClosed: []byte(tt.startKeyClosed),
						},
						EndKey: &btpb.RowRange_EndKeyClosed{
							EndKeyClosed: []byte(tt.endKeyClosed),
						},
					},
					{
						StartKey: &btpb.RowRange_StartKeyClosed{
							StartKeyClosed: []byte(tt.startKeyClosed),
						},
						EndKey: &btpb.RowRange_EndKeyOpen{
							EndKeyOpen: []byte(tt.endKeyOpen),
						},
					},
					{
						StartKey: &btpb.RowRange_StartKeyOpen{
							StartKeyOpen: []byte(tt.startKeyOpen),
						},
						EndKey: &btpb.RowRange_EndKeyOpen{
							EndKeyOpen: []byte(tt.endKeyOpen),
						},
					},
					{
						StartKey: &btpb.RowRange_StartKeyOpen{
							StartKeyOpen: []byte(tt.startKeyOpen),
						},
						EndKey: &btpb.RowRange_EndKeyClosed{
							EndKeyClosed: []byte(tt.endKeyClosed),
						},
					},
				},
			},
		}

		err := srv.ReadRows(badReq, nil)
		if err == nil {
			t.Errorf("#%d: unexpectedly returned nil error", i)
			continue
		}

		status, ok := status.FromError(err)
		if !ok {
			t.Errorf("#%d: wrong error type %T, expected status.Error", i, err)
			continue
		}
		if g, w := status.Code(), codes.InvalidArgument; g != w {
			t.Errorf("#%d: wrong error code\nGot  %d %s\nWant %d %s", i, g, g, w, w)
		}
	}
}
