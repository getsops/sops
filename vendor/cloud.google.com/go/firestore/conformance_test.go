// Copyright 2017 Google LLC
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

// A runner for the conformance tests.

package firestore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	pb "cloud.google.com/go/firestore/internal/conformance"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	ts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/api/iterator"
	fspb "google.golang.org/genproto/googleapis/firestore/v1"
)

func TestConformance(t *testing.T) {
	dir := "internal/conformance/testdata"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	wtid := watchTargetID
	watchTargetID = 1
	defer func() { watchTargetID = wtid }()

	for _, f := range files {
		if !strings.Contains(f.Name(), ".json") {
			continue
		}

		inBytes, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			t.Fatalf("%s: %v", f.Name(), err)
		}

		var tf pb.TestFile
		if err := jsonpb.Unmarshal(bytes.NewReader(inBytes), &tf); err != nil {
			t.Fatalf("unmarshalling %s: %v", f.Name(), err)
		}

		for _, tc := range tf.Tests {
			t.Run(tc.Description, func(t *testing.T) {
				c, srv, cleanup := newMock(t)
				defer cleanup()

				if err := runTest(tc, c, srv); err != nil {
					t.Fatal(err)
				}
			})
		}
	}
}

func runTest(test *pb.Test, c *Client, srv *mockServer) error {
	check := func(gotErr error, wantErr bool) error {
		if wantErr && gotErr == nil {
			return errors.New("got nil, want error")
		}
		if !wantErr && gotErr != nil {
			return gotErr
		}
		return nil
	}

	ctx := context.Background()
	switch typedTestcase := test.Test.(type) {
	case *pb.Test_Get:
		req := &fspb.BatchGetDocumentsRequest{
			Database:  c.path(),
			Documents: []string{typedTestcase.Get.DocRefPath},
		}
		srv.addRPC(req, []interface{}{
			&fspb.BatchGetDocumentsResponse{
				Result: &fspb.BatchGetDocumentsResponse_Found{&fspb.Document{
					Name:       typedTestcase.Get.DocRefPath,
					CreateTime: &ts.Timestamp{},
					UpdateTime: &ts.Timestamp{},
				}},
				ReadTime: &ts.Timestamp{},
			},
		})
		ref := docRefFromPath(typedTestcase.Get.DocRefPath, c)
		_, err := ref.Get(ctx)
		if err != nil {
			return err
		}
		// Checking response would just be testing the function converting a Document
		// proto to a DocumentSnapshot, hence uninteresting.

	case *pb.Test_Create:
		srv.addRPC(typedTestcase.Create.Request, commitResponseForSet)
		ref := docRefFromPath(typedTestcase.Create.DocRefPath, c)
		data, err := convertData(typedTestcase.Create.JsonData)
		if err != nil {
			return err
		}
		_, checkErr := ref.Create(ctx, data)
		if err := check(checkErr, typedTestcase.Create.IsError); err != nil {
			return err
		}

	case *pb.Test_Set:
		srv.addRPC(typedTestcase.Set.Request, commitResponseForSet)
		ref := docRefFromPath(typedTestcase.Set.DocRefPath, c)
		data, err := convertData(typedTestcase.Set.JsonData)
		if err != nil {
			return err
		}
		var opts []SetOption
		if typedTestcase.Set.Option != nil {
			opts = []SetOption{convertSetOption(typedTestcase.Set.Option)}
		}
		_, checkErr := ref.Set(ctx, data, opts...)
		if err := check(checkErr, typedTestcase.Set.IsError); err != nil {
			return err
		}

	case *pb.Test_Update:
		// Ignore Update test because we only support UpdatePaths.
		// Not to worry, every Update test has a corresponding UpdatePaths test.

	case *pb.Test_UpdatePaths:
		srv.addRPC(typedTestcase.UpdatePaths.Request, commitResponseForSet)
		ref := docRefFromPath(typedTestcase.UpdatePaths.DocRefPath, c)
		preconds, err := convertPrecondition(typedTestcase.UpdatePaths.Precondition)
		if err != nil {
			return err
		}
		paths := convertFieldPaths(typedTestcase.UpdatePaths.FieldPaths)
		var ups []Update
		for i, p := range paths {
			val, err := convertJSONValue(typedTestcase.UpdatePaths.JsonValues[i])
			if err != nil {
				return err
			}
			ups = append(ups, Update{
				FieldPath: p,
				Value:     val,
			})
		}
		_, checkErr := ref.Update(ctx, ups, preconds...)
		if err := check(checkErr, typedTestcase.UpdatePaths.IsError); err != nil {
			return err
		}

	case *pb.Test_Delete:
		srv.addRPC(typedTestcase.Delete.Request, commitResponseForSet)
		ref := docRefFromPath(typedTestcase.Delete.DocRefPath, c)
		preconds, err := convertPrecondition(typedTestcase.Delete.Precondition)
		if err != nil {
			return err
		}
		_, checkErr := ref.Delete(ctx, preconds...)
		if err := check(checkErr, typedTestcase.Delete.IsError); err != nil {
			return err
		}

	case *pb.Test_Query:
		q, err := convertQuery(typedTestcase.Query)
		if err != nil {
			return err
		}
		got, checkErr := q.toProto()
		if err := check(checkErr, typedTestcase.Query.IsError); err == nil && checkErr == nil {
			if want := typedTestcase.Query.Query; !proto.Equal(got, want) {
				return fmt.Errorf("got:  %s\nwant: %s", proto.MarshalTextString(got), proto.MarshalTextString(want))
			}
		}

	case *pb.Test_Listen:
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		iter := c.Collection("C").OrderBy("a", Asc).Snapshots(ctx)
		var rs []interface{}
		for _, r := range typedTestcase.Listen.Responses {
			rs = append(rs, r)
		}
		srv.addRPC(&fspb.ListenRequest{
			Database:     "projects/projectID/databases/(default)",
			TargetChange: &fspb.ListenRequest_AddTarget{iter.ws.target},
		}, rs)
		got, err := nSnapshots(iter, len(typedTestcase.Listen.Snapshots))
		if err != nil {
			return err
		} else if diff := cmp.Diff(got, typedTestcase.Listen.Snapshots); diff != "" {
			return errors.New(diff)
		}
		if typedTestcase.Listen.IsError {
			_, err := iter.Next()
			if err == nil {
				return fmt.Errorf("got nil, want error")
			}
		}

	default:
		return fmt.Errorf("unknown test type %T", typedTestcase)
	}

	return nil
}

func nSnapshots(iter *QuerySnapshotIterator, n int) ([]*pb.Snapshot, error) {
	var snaps []*pb.Snapshot
	for i := 0; i < n; i++ {
		qsnap, err := iter.Next()
		if err != nil {
			return snaps, err
		}
		s := &pb.Snapshot{ReadTime: mustTimestampProto(qsnap.ReadTime)}
		for {
			doc, err := qsnap.Documents.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return snaps, err
			}
			s.Docs = append(s.Docs, doc.proto)
		}
		for _, c := range qsnap.Changes {
			var k pb.DocChange_Kind
			switch c.Kind {
			case DocumentAdded:
				k = pb.DocChange_ADDED
			case DocumentRemoved:
				k = pb.DocChange_REMOVED
			case DocumentModified:
				k = pb.DocChange_MODIFIED
			default:
				panic("bad kind")
			}
			s.Changes = append(s.Changes, &pb.DocChange{
				Kind:     k,
				Doc:      c.Doc.proto,
				OldIndex: int32(c.OldIndex),
				NewIndex: int32(c.NewIndex),
			})
		}
		snaps = append(snaps, s)
	}
	return snaps, nil
}

func docRefFromPath(p string, c *Client) *DocumentRef {
	return &DocumentRef{
		Path:   p,
		ID:     path.Base(p),
		Parent: &CollectionRef{c: c},
	}
}

func convertJSONValue(jv string) (interface{}, error) {
	var val interface{}
	if err := json.Unmarshal([]byte(jv), &val); err != nil {
		return nil, err
	}
	return convertTestValue(val), nil
}

func convertData(jsonData string) (map[string]interface{}, error) {
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &m); err != nil {
		return nil, err
	}
	return convertTestMap(m), nil
}

func convertTestMap(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		m[k] = convertTestValue(v)
	}
	return m
}

func convertTestValue(v interface{}) interface{} {
	switch v := v.(type) {
	case string:
		switch v {
		case "ServerTimestamp":
			return ServerTimestamp
		case "Delete":
			return Delete
		case "NaN":
			return math.NaN()
		default:
			return v
		}
	case float64:
		if v == float64(int(v)) {
			return int(v)
		}
		return v
	case []interface{}:
		if len(v) > 0 {
			if fv, ok := v[0].(string); ok {
				if fv == "ArrayUnion" {
					return ArrayUnion(convertTestValue(v[1:]).([]interface{})...)
				}
				if fv == "ArrayRemove" {
					return ArrayRemove(convertTestValue(v[1:]).([]interface{})...)
				}
			}
		}
		for i, e := range v {
			v[i] = convertTestValue(e)
		}
		return v
	case map[string]interface{}:
		return convertTestMap(v)
	default:
		return v
	}
}

func convertSetOption(opt *pb.SetOption) SetOption {
	if opt.All {
		return MergeAll
	}
	return Merge(convertFieldPaths(opt.Fields)...)
}

func convertFieldPaths(fps []*pb.FieldPath) []FieldPath {
	var res []FieldPath
	for _, fp := range fps {
		res = append(res, fp.Field)
	}
	return res
}

func convertPrecondition(fp *fspb.Precondition) ([]Precondition, error) {
	if fp == nil {
		return nil, nil
	}
	var pc Precondition
	switch fp := fp.ConditionType.(type) {
	case *fspb.Precondition_Exists:
		pc = exists(fp.Exists)
	case *fspb.Precondition_UpdateTime:
		tm, err := ptypes.Timestamp(fp.UpdateTime)
		if err != nil {
			return nil, err
		}
		pc = LastUpdateTime(tm)
	default:
		return nil, fmt.Errorf("unknown precondition type %T", fp)
	}
	return []Precondition{pc}, nil
}

func convertQuery(qt *pb.QueryTest) (*Query, error) {
	parts := strings.Split(qt.CollPath, "/")
	q := Query{
		parentPath:   strings.Join(parts[:len(parts)-2], "/"),
		collectionID: parts[len(parts)-1],
		path:         qt.CollPath,
	}
	for _, c := range qt.Clauses {
		switch c := c.Clause.(type) {
		case *pb.Clause_Select:
			q = q.SelectPaths(convertFieldPaths(c.Select.Fields)...)
		case *pb.Clause_OrderBy:
			var dir Direction
			switch c.OrderBy.Direction {
			case "asc":
				dir = Asc
			case "desc":
				dir = Desc
			default:
				return nil, fmt.Errorf("bad direction: %q", c.OrderBy.Direction)
			}
			q = q.OrderByPath(FieldPath(c.OrderBy.Path.Field), dir)
		case *pb.Clause_Where:
			val, err := convertJSONValue(c.Where.JsonValue)
			if err != nil {
				return nil, err
			}
			q = q.WherePath(FieldPath(c.Where.Path.Field), c.Where.Op, val)
		case *pb.Clause_Offset:
			q = q.Offset(int(c.Offset))
		case *pb.Clause_Limit:
			q = q.Limit(int(c.Limit))
		case *pb.Clause_StartAt:
			cs, err := convertCursor(c.StartAt)
			if err != nil {
				return nil, err
			}
			q = q.StartAt(cs...)
		case *pb.Clause_StartAfter:
			cs, err := convertCursor(c.StartAfter)
			if err != nil {
				return nil, err
			}
			q = q.StartAfter(cs...)
		case *pb.Clause_EndAt:
			cs, err := convertCursor(c.EndAt)
			if err != nil {
				return nil, err
			}
			q = q.EndAt(cs...)
		case *pb.Clause_EndBefore:
			cs, err := convertCursor(c.EndBefore)
			if err != nil {
				return nil, err
			}
			q = q.EndBefore(cs...)
		default:
			return nil, fmt.Errorf("bad clause type %T", c)
		}
	}
	return &q, nil
}

// Returns args to a cursor method (StartAt, etc.).
func convertCursor(c *pb.Cursor) ([]interface{}, error) {
	if c.DocSnapshot != nil {
		ds, err := convertDocSnapshot(c.DocSnapshot)
		if err != nil {
			return nil, err
		}
		return []interface{}{ds}, nil
	}
	var vals []interface{}
	for _, jv := range c.JsonValues {
		v, err := convertJSONValue(jv)
		if err != nil {
			return nil, err
		}
		vals = append(vals, v)
	}
	return vals, nil
}

func convertDocSnapshot(ds *pb.DocSnapshot) (*DocumentSnapshot, error) {
	data, err := convertData(ds.JsonData)
	if err != nil {
		return nil, err
	}
	doc, transformPaths, err := toProtoDocument(data)
	if err != nil {
		return nil, err
	}
	if len(transformPaths) > 0 {
		return nil, errors.New("saw transform paths in DocSnapshot")
	}
	return &DocumentSnapshot{
		Ref: &DocumentRef{
			Path:   ds.Path,
			Parent: &CollectionRef{Path: path.Dir(ds.Path)},
		},
		proto: doc,
	}, nil
}
