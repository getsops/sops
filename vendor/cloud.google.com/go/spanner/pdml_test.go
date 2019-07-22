// Copyright 2018 Google LLC
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

package spanner

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
)

func TestMockPartitionedUpdate(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)

	stmt := NewStatement(updateBarSetFoo)
	rowCount, err := client.PartitionedUpdate(ctx, stmt)
	if err != nil {
		t.Fatal(err)
	}
	want := int64(updateBarSetFooRowCount)
	if rowCount != want {
		t.Errorf("got %d, want %d", rowCount, want)
	}
}

func TestMockPartitionedUpdateWithQuery(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)

	stmt := NewStatement(selectFooFromBar)
	_, err := client.PartitionedUpdate(ctx, stmt)
	wantCode := codes.InvalidArgument
	if serr, ok := err.(*Error); !ok || serr.Code != wantCode {
		t.Errorf("got error %v, want code %s", err, wantCode)
	}
}
