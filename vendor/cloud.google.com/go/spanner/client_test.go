/*
Copyright 2017 Google LLC

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

package spanner

import (
	"context"
	"io"
	"strings"
	"testing"

	"cloud.google.com/go/spanner/internal/testutil"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	gstatus "google.golang.org/grpc/status"
)

// Test validDatabaseName()
func TestValidDatabaseName(t *testing.T) {
	validDbURI := "projects/spanner-cloud-test/instances/foo/databases/foodb"
	invalidDbUris := []string{
		// Completely wrong DB URI.
		"foobarDB",
		// Project ID contains "/".
		"projects/spanner-cloud/test/instances/foo/databases/foodb",
		// No instance ID.
		"projects/spanner-cloud-test/instances//databases/foodb",
	}
	if err := validDatabaseName(validDbURI); err != nil {
		t.Errorf("validateDatabaseName(%q) = %v, want nil", validDbURI, err)
	}
	for _, d := range invalidDbUris {
		if err, wantErr := validDatabaseName(d), "should conform to pattern"; !strings.Contains(err.Error(), wantErr) {
			t.Errorf("validateDatabaseName(%q) = %q, want error pattern %q", validDbURI, err, wantErr)
		}
	}
}

func TestReadOnlyTransactionClose(t *testing.T) {
	// Closing a ReadOnlyTransaction shouldn't panic.
	c := &Client{}
	tx := c.ReadOnlyTransaction()
	tx.Close()
}

func TestClient_Single(t *testing.T) {
	t.Parallel()
	err := testSingleQuery(t, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Single_Unavailable(t *testing.T) {
	t.Parallel()
	err := testSingleQuery(t, gstatus.Error(codes.Unavailable, "Temporary unavailable"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestClient_Single_InvalidArgument(t *testing.T) {
	t.Parallel()
	err := testSingleQuery(t, gstatus.Error(codes.InvalidArgument, "Invalid argument"))
	if err == nil {
		t.Fatalf("missing expected error")
	} else if gstatus.Code(err) != codes.InvalidArgument {
		t.Fatal(err)
	}
}

func testSingleQuery(t *testing.T, serverError error) error {
	config := ClientConfig{}
	server, client := newSpannerInMemTestServerWithConfig(t, config)
	defer server.teardown(client)
	if serverError != nil {
		server.testSpanner.SetError(serverError)
	}
	ctx := context.Background()
	iter := client.Single().Query(ctx, NewStatement(selectSingerIDAlbumIDAlbumTitleFromAlbums))
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
	}
	return nil
}

func createSimulatedExecutionTimeWithTwoUnavailableErrors(method string) map[string]testutil.SimulatedExecutionTime {
	errors := make([]error, 2)
	errors[0] = gstatus.Error(codes.Unavailable, "Temporary unavailable")
	errors[1] = gstatus.Error(codes.Unavailable, "Temporary unavailable")
	executionTimes := make(map[string]testutil.SimulatedExecutionTime)
	executionTimes[method] = testutil.SimulatedExecutionTime{
		Errors: errors,
	}
	return executionTimes
}

func TestClient_ReadOnlyTransaction(t *testing.T) {
	t.Parallel()
	if err := testReadOnlyTransaction(t, make(map[string]testutil.SimulatedExecutionTime)); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadOnlyTransaction_UnavailableOnSessionCreate(t *testing.T) {
	t.Parallel()
	if err := testReadOnlyTransaction(t, createSimulatedExecutionTimeWithTwoUnavailableErrors(testutil.MethodCreateSession)); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadOnlyTransaction_UnavailableOnBeginTransaction(t *testing.T) {
	t.Parallel()
	if err := testReadOnlyTransaction(t, createSimulatedExecutionTimeWithTwoUnavailableErrors(testutil.MethodBeginTransaction)); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadOnlyTransaction_UnavailableOnExecuteStreamingSql(t *testing.T) {
	t.Parallel()
	if err := testReadOnlyTransaction(t, createSimulatedExecutionTimeWithTwoUnavailableErrors(testutil.MethodExecuteStreamingSql)); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadOnlyTransaction_UnavailableOnCreateSessionAndBeginTransaction(t *testing.T) {
	t.Parallel()
	exec := map[string]testutil.SimulatedExecutionTime{
		testutil.MethodCreateSession:    {Errors: []error{gstatus.Error(codes.Unavailable, "Temporary unavailable")}},
		testutil.MethodBeginTransaction: {Errors: []error{gstatus.Error(codes.Unavailable, "Temporary unavailable")}},
	}
	if err := testReadOnlyTransaction(t, exec); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadOnlyTransaction_UnavailableOnCreateSessionAndInvalidArgumentOnBeginTransaction(t *testing.T) {
	t.Parallel()
	exec := map[string]testutil.SimulatedExecutionTime{
		testutil.MethodCreateSession:    {Errors: []error{gstatus.Error(codes.Unavailable, "Temporary unavailable")}},
		testutil.MethodBeginTransaction: {Errors: []error{gstatus.Error(codes.InvalidArgument, "Invalid argument")}},
	}
	if err := testReadOnlyTransaction(t, exec); err == nil {
		t.Fatalf("Missing expected exception")
	} else if gstatus.Code(err) != codes.InvalidArgument {
		t.Fatalf("Got unexpected exception: %v", err)
	}
}

func testReadOnlyTransaction(t *testing.T, executionTimes map[string]testutil.SimulatedExecutionTime) error {
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)
	for method, exec := range executionTimes {
		server.testSpanner.PutExecutionTime(method, exec)
	}
	ctx := context.Background()
	tx := client.ReadOnlyTransaction()
	defer tx.Close()
	iter := tx.Query(ctx, NewStatement(selectSingerIDAlbumIDAlbumTitleFromAlbums))
	defer iter.Stop()
	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		var singerID, albumID int64
		var albumTitle string
		if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
			return err
		}
	}
	return nil
}

func TestClient_ReadWriteTransaction(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, make(map[string]testutil.SimulatedExecutionTime), 1); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransactionCommitAborted(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodCommitTransaction: {Errors: []error{gstatus.Error(codes.Aborted, "Transaction aborted")}},
	}, 2); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransactionExecuteStreamingSqlAborted(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodExecuteStreamingSql: {Errors: []error{gstatus.Error(codes.Aborted, "Transaction aborted")}},
	}, 2); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransaction_UnavailableOnBeginTransaction(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodBeginTransaction: {Errors: []error{gstatus.Error(codes.Unavailable, "Unavailable")}},
	}, 1); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransaction_UnavailableOnBeginAndAbortOnCommit(t *testing.T) {
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodBeginTransaction:  {Errors: []error{gstatus.Error(codes.Unavailable, "Unavailable")}},
		testutil.MethodCommitTransaction: {Errors: []error{gstatus.Error(codes.Aborted, "Aborted")}},
	}, 2); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransaction_UnavailableOnExecuteStreamingSql(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodExecuteStreamingSql: {Errors: []error{gstatus.Error(codes.Unavailable, "Unavailable")}},
	}, 1); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransaction_UnavailableOnBeginAndExecuteStreamingSqlAndTwiceAbortOnCommit(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodBeginTransaction:    {Errors: []error{gstatus.Error(codes.Unavailable, "Unavailable")}},
		testutil.MethodExecuteStreamingSql: {Errors: []error{gstatus.Error(codes.Unavailable, "Unavailable")}},
		testutil.MethodCommitTransaction:   {Errors: []error{gstatus.Error(codes.Aborted, "Aborted"), gstatus.Error(codes.Aborted, "Aborted")}},
	}, 3); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransaction_AbortedOnExecuteStreamingSqlAndCommit(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodExecuteStreamingSql: {Errors: []error{gstatus.Error(codes.Aborted, "Aborted")}},
		testutil.MethodCommitTransaction:   {Errors: []error{gstatus.Error(codes.Aborted, "Aborted"), gstatus.Error(codes.Aborted, "Aborted")}},
	}, 4); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransactionCommitAbortedAndUnavailable(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodCommitTransaction: {
			Errors: []error{
				gstatus.Error(codes.Aborted, "Transaction aborted"),
				gstatus.Error(codes.Unavailable, "Unavailable"),
			},
		},
	}, 2); err != nil {
		t.Fatal(err)
	}
}

func TestClient_ReadWriteTransactionCommitAlreadyExists(t *testing.T) {
	t.Parallel()
	if err := testReadWriteTransaction(t, map[string]testutil.SimulatedExecutionTime{
		testutil.MethodCommitTransaction: {Errors: []error{gstatus.Error(codes.AlreadyExists, "A row with this key already exists")}},
	}, 1); err != nil {
		if gstatus.Code(err) != codes.AlreadyExists {
			t.Fatalf("Got unexpected error %v, expected %v", err, codes.AlreadyExists)
		}
	} else {
		t.Fatalf("Missing expected exception")
	}
}

func testReadWriteTransaction(t *testing.T, executionTimes map[string]testutil.SimulatedExecutionTime, expectedAttempts int) error {
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)
	for method, exec := range executionTimes {
		server.testSpanner.PutExecutionTime(method, exec)
	}
	var attempts int
	ctx := context.Background()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *ReadWriteTransaction) error {
		attempts++
		iter := tx.Query(ctx, NewStatement(selectSingerIDAlbumIDAlbumTitleFromAlbums))
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var singerID, albumID int64
			var albumTitle string
			if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if expectedAttempts != attempts {
		t.Fatalf("unexpected number of attempts: %d, expected %d", attempts, expectedAttempts)
	}
	return nil
}

func TestClient_ApplyAtLeastOnce(t *testing.T) {
	t.Parallel()
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)
	ms := []*Mutation{
		Insert("Accounts", []string{"AccountId", "Nickname", "Balance"}, []interface{}{int64(1), "Foo", int64(50)}),
		Insert("Accounts", []string{"AccountId", "Nickname", "Balance"}, []interface{}{int64(2), "Bar", int64(1)}),
	}
	server.testSpanner.PutExecutionTime(testutil.MethodCommitTransaction,
		testutil.SimulatedExecutionTime{
			Errors: []error{gstatus.Error(codes.Aborted, "Transaction aborted")},
		})
	_, err := client.Apply(context.Background(), ms, ApplyAtLeastOnce())
	if err != nil {
		t.Fatal(err)
	}
}

// PartitionedUpdate should not retry on aborted.
func TestClient_PartitionedUpdate(t *testing.T) {
	t.Parallel()
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)
	// PartitionedDML transactions are not committed.
	server.testSpanner.PutExecutionTime(testutil.MethodExecuteStreamingSql,
		testutil.SimulatedExecutionTime{
			Errors: []error{gstatus.Error(codes.Aborted, "Transaction aborted")},
		})
	_, err := client.PartitionedUpdate(context.Background(), NewStatement(updateBarSetFoo))
	if err == nil {
		t.Fatalf("Missing expected Aborted exception")
	} else {
		if gstatus.Code(err) != codes.Aborted {
			t.Fatalf("Got unexpected error %v, expected Aborted", err)
		}
	}
}

func TestReadWriteTransaction_ErrUnexpectedEOF(t *testing.T) {
	server, client := newSpannerInMemTestServer(t)
	defer server.teardown(client)
	var attempts int
	ctx := context.Background()
	_, err := client.ReadWriteTransaction(ctx, func(ctx context.Context, tx *ReadWriteTransaction) error {
		attempts++
		iter := tx.Query(ctx, NewStatement(selectSingerIDAlbumIDAlbumTitleFromAlbums))
		defer iter.Stop()
		for {
			row, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var singerID, albumID int64
			var albumTitle string
			if err := row.Columns(&singerID, &albumID, &albumTitle); err != nil {
				return err
			}
		}
		return io.ErrUnexpectedEOF
	})
	if err != io.ErrUnexpectedEOF {
		t.Fatalf("Missing expected error %v, got %v", io.ErrUnexpectedEOF, err)
	}
	if attempts != 1 {
		t.Fatalf("unexpected number of attempts: %d, expected %d", attempts, 1)
	}
}
