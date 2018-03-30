// Copyright 2017 Google Inc. All Rights Reserved.
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

package firestore

import (
	"golang.org/x/net/context"
	"testing"

	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

func TestRunTransaction(t *testing.T) {
	ctx := context.Background()
	const db = "projects/projectID/databases/(default)"
	tid := []byte{1}
	c, srv := newMock(t)
	beginReq := &pb.BeginTransactionRequest{Database: db}
	beginRes := &pb.BeginTransactionResponse{Transaction: tid}
	commitReq := &pb.CommitRequest{Database: db, Transaction: tid}
	// Empty transaction.
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(commitReq, &pb.CommitResponse{CommitTime: aTimestamp})
	err := c.RunTransaction(ctx, func(context.Context, *Transaction) error { return nil })
	if err != nil {
		t.Fatal(err)
	}

	// Transaction with read and write.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	aDoc := &pb.Document{
		Name:       db + "/documents/C/a",
		CreateTime: aTimestamp,
		UpdateTime: aTimestamp2,
		Fields:     map[string]*pb.Value{"count": intval(1)},
	}
	srv.addRPC(
		&pb.GetDocumentRequest{
			Name:                db + "/documents/C/a",
			ConsistencySelector: &pb.GetDocumentRequest_Transaction{tid},
		},
		aDoc,
	)
	aDoc2 := &pb.Document{
		Name:   aDoc.Name,
		Fields: map[string]*pb.Value{"count": intval(2)},
	}
	srv.addRPC(
		&pb.CommitRequest{
			Database:    db,
			Transaction: tid,
			Writes: []*pb.Write{{
				Operation:  &pb.Write_Update{aDoc2},
				UpdateMask: &pb.DocumentMask{FieldPaths: []string{"count"}},
				CurrentDocument: &pb.Precondition{
					ConditionType: &pb.Precondition_Exists{true},
				},
			}},
		},
		&pb.CommitResponse{CommitTime: aTimestamp3},
	)
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error {
		docref := c.Collection("C").Doc("a")
		doc, err := tx.Get(docref)
		if err != nil {
			return err
		}
		count, err := doc.DataAt("count")
		if err != nil {
			return err
		}
		tx.UpdateMap(docref, map[string]interface{}{"count": count.(int64) + 1})
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Query
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(
		&pb.RunQueryRequest{
			Parent: db,
			QueryType: &pb.RunQueryRequest_StructuredQuery{
				&pb.StructuredQuery{
					From: []*pb.StructuredQuery_CollectionSelector{{CollectionId: "C"}},
				},
			},
			ConsistencySelector: &pb.RunQueryRequest_Transaction{tid},
		},
		[]interface{}{},
	)
	srv.addRPC(commitReq, &pb.CommitResponse{CommitTime: aTimestamp3})
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error {
		it := tx.Documents(c.Collection("C"))
		_, err := it.Next()
		if err != iterator.Done {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	// Retry entire transaction.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(commitReq, grpc.Errorf(codes.Aborted, ""))
	srv.addRPC(
		&pb.BeginTransactionRequest{
			Database: db,
			Options: &pb.TransactionOptions{
				Mode: &pb.TransactionOptions_ReadWrite_{
					&pb.TransactionOptions_ReadWrite{tid},
				},
			},
		},
		beginRes,
	)
	srv.addRPC(commitReq, &pb.CommitResponse{CommitTime: aTimestamp})
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
}

func TestTransactionErrors(t *testing.T) {
	ctx := context.Background()
	const db = "projects/projectID/databases/(default)"
	c, srv := newMock(t)
	var (
		tid         = []byte{1}
		internalErr = grpc.Errorf(codes.Internal, "so sad")
		beginReq    = &pb.BeginTransactionRequest{
			Database: db,
		}
		beginRes = &pb.BeginTransactionResponse{Transaction: tid}
		getReq   = &pb.GetDocumentRequest{
			Name:                db + "/documents/C/a",
			ConsistencySelector: &pb.GetDocumentRequest_Transaction{tid},
		}
		rollbackReq = &pb.RollbackRequest{Database: db, Transaction: tid}
		commitReq   = &pb.CommitRequest{Database: db, Transaction: tid}
	)

	// BeginTransaction has a permanent error.
	srv.addRPC(beginReq, internalErr)
	err := c.RunTransaction(ctx, func(context.Context, *Transaction) error { return nil })
	if grpc.Code(err) != codes.Internal {
		t.Errorf("got <%v>, want Internal", err)
	}

	// Get has a permanent error.
	get := func(_ context.Context, tx *Transaction) error {
		_, err := tx.Get(c.Doc("C/a"))
		return err
	}
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(getReq, internalErr)
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, get)
	if grpc.Code(err) != codes.Internal {
		t.Errorf("got <%v>, want Internal", err)
	}

	// Get has a permanent error, but the rollback fails. We still
	// return Get's error.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(getReq, internalErr)
	srv.addRPC(rollbackReq, grpc.Errorf(codes.FailedPrecondition, ""))
	err = c.RunTransaction(ctx, get)
	if grpc.Code(err) != codes.Internal {
		t.Errorf("got <%v>, want Internal", err)
	}

	// Commit has a permanent error.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(getReq, &pb.Document{
		Name:       "projects/projectID/databases/(default)/documents/C/a",
		CreateTime: aTimestamp,
		UpdateTime: aTimestamp2,
	})
	srv.addRPC(commitReq, internalErr)
	err = c.RunTransaction(ctx, get)
	if grpc.Code(err) != codes.Internal {
		t.Errorf("got <%v>, want Internal", err)
	}

	// Read after write.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error {
		tx.Delete(c.Doc("C/a"))
		if _, err := tx.Get(c.Doc("C/a")); err != nil {
			return err
		}
		return nil
	})
	if err != errReadAfterWrite {
		t.Errorf("got <%v>, want <%v>", err, errReadAfterWrite)
	}

	// Read after write, with query.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error {
		tx.Delete(c.Doc("C/a"))
		it := tx.Documents(c.Collection("C").Select("x"))
		if _, err := it.Next(); err != iterator.Done {
			return err
		}
		return nil
	})
	if err != errReadAfterWrite {
		t.Errorf("got <%v>, want <%v>", err, errReadAfterWrite)
	}

	// Read after write fails even if the user ignores the read's error.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error {
		tx.Delete(c.Doc("C/a"))
		tx.Get(c.Doc("C/a"))
		return nil
	})
	if err != errReadAfterWrite {
		t.Errorf("got <%v>, want <%v>", err, errReadAfterWrite)
	}

	// Write in read-only transaction.
	srv.reset()
	srv.addRPC(
		&pb.BeginTransactionRequest{
			Database: db,
			Options: &pb.TransactionOptions{
				Mode: &pb.TransactionOptions_ReadOnly_{&pb.TransactionOptions_ReadOnly{}},
			},
		},
		beginRes,
	)
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, func(_ context.Context, tx *Transaction) error {
		return tx.Delete(c.Doc("C/a"))
	}, ReadOnly)
	if err != errWriteReadOnly {
		t.Errorf("got <%v>, want <%v>", err, errWriteReadOnly)
	}

	// Too many retries.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(commitReq, grpc.Errorf(codes.Aborted, ""))
	srv.addRPC(
		&pb.BeginTransactionRequest{
			Database: db,
			Options: &pb.TransactionOptions{
				Mode: &pb.TransactionOptions_ReadWrite_{
					&pb.TransactionOptions_ReadWrite{tid},
				},
			},
		},
		beginRes,
	)
	srv.addRPC(commitReq, grpc.Errorf(codes.Aborted, ""))
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, func(context.Context, *Transaction) error { return nil },
		MaxAttempts(2))
	if grpc.Code(err) != codes.Aborted {
		t.Errorf("got <%v>, want Aborted", err)
	}

	// Nested transaction.
	srv.reset()
	srv.addRPC(beginReq, beginRes)
	srv.addRPC(rollbackReq, &empty.Empty{})
	err = c.RunTransaction(ctx, func(ctx context.Context, tx *Transaction) error {
		return c.RunTransaction(ctx, func(context.Context, *Transaction) error { return nil })
	})
	if got, want := err, errNestedTransaction; got != want {
		t.Errorf("got <%v>, want <%V>", got, want)
	}

	// Non-transactional operation.
	dr := c.Doc("C/d")

	for i, op := range []func(ctx context.Context) error{
		func(ctx context.Context) error { _, err := c.GetAll(ctx, []*DocumentRef{dr}); return err },
		func(ctx context.Context) error { _, _, err := c.Collection("C").Add(ctx, testData); return err },
		func(ctx context.Context) error { _, err := dr.Get(ctx); return err },
		func(ctx context.Context) error { _, err := dr.Create(ctx, testData); return err },
		func(ctx context.Context) error { _, err := dr.Set(ctx, testData); return err },
		func(ctx context.Context) error { _, err := dr.Delete(ctx); return err },
		func(ctx context.Context) error { _, err := dr.UpdateMap(ctx, testData); return err },
		func(ctx context.Context) error {
			_, err := dr.UpdateStruct(ctx, []string{"x"}, struct{}{})
			return err
		},
		func(ctx context.Context) error {
			_, err := dr.UpdatePaths(ctx, []FieldPathUpdate{{Path: []string{"*"}, Value: 1}})
			return err
		},
		func(ctx context.Context) error { it := c.Collections(ctx); _, err := it.Next(); return err },
		func(ctx context.Context) error { it := dr.Collections(ctx); _, err := it.Next(); return err },
		func(ctx context.Context) error { _, err := c.Batch().Commit(ctx); return err },
		func(ctx context.Context) error {
			it := c.Collection("C").Documents(ctx)
			_, err := it.Next()
			return err
		},
	} {
		srv.reset()
		srv.addRPC(beginReq, beginRes)
		srv.addRPC(rollbackReq, &empty.Empty{})
		err = c.RunTransaction(ctx, func(ctx context.Context, _ *Transaction) error {
			return op(ctx)
		})
		if got, want := err, errNonTransactionalOp; got != want {
			t.Errorf("#%d: got <%v>, want <%v>", i, got, want)
		}
	}
}
