// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bttest

import (
	"context"
	"testing"

	btapb "google.golang.org/genproto/googleapis/bigtable/admin/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteInstance(t *testing.T) {
	srv := &server{
		instances: map[string]*btapb.Instance{
			"projects/test/instances/a1042-instance": {
				Name:        "projects/test/instances/a1042-instance/test1",
				DisplayName: "test1",
				State:       btapb.Instance_READY,
				Labels: map[string]string{
					"component":   "ingest",
					"cost-center": "sales",
				},
			},
		},
	}

	ctx := context.Background()

	// 1. Deletion of a just created instance should succeed.
	delReq := &btapb.DeleteInstanceRequest{Name: "projects/test/instances/a1042-instance"}
	if _, err := srv.DeleteInstance(ctx, delReq); err != nil {
		t.Fatalf("Unexpected failed to delete the newly created instance: %v", err)
	}

	// 2. Deleting a now non-existent instance should fail.
	_, err := srv.DeleteInstance(ctx, delReq)
	if err == nil {
		t.Fatal("Expected an error from deleting a non-existent instance")
	}
	gstatus, _ := status.FromError(err)
	if g, w := gstatus.Code(), codes.NotFound; g != w {
		t.Errorf("Mismatched status code\ngot  %d %s\nwant %d %s", g, g, w, w)
	}
	if g, w := gstatus.Message(), `instance "projects/test/instances/a1042-instance" not found`; g != w {
		t.Errorf("Mismatched status message\ngot  %q\nwant %q", g, w)
	}

	// 3. Now test deletion with invalid names which should always fail.
	invalidNames := []string{
		"", "  ", "//", "/",
		"projects/foo",
		"projects/foo/instances",
		"projects/foo/instances/",
		"projects/foo/instances/a10/bar",
		"projects/foo/instances/a10/bar/fizz",
		"/projects/foo*",
	}

	for _, invalidName := range invalidNames {
		req := &btapb.DeleteInstanceRequest{Name: invalidName}
		_, err := srv.DeleteInstance(ctx, req)
		gstatus, _ := status.FromError(err)
		if g, w := gstatus.Code(), codes.InvalidArgument; g != w {
			t.Errorf("Mismatched status code\ngot  %d %s\nwant %d %s", g, g, w, w)
		}

		wantMsg := "Error in field 'instance_name' : Invalid name for collection instances : " +
			"Should match " + `^projects/[a-z][a-z0-9\\-]+[a-z0-9]/instances/[a-z][a-z0-9\\-]+[a-z0-9]$` +
			" but found '" + invalidName + "'"
		if g, w := gstatus.Message(), wantMsg; g != w {
			t.Errorf("%q: mismatched status message\ngot  %q\nwant %q\n\n", invalidName, g, w)
		}
	}
}
