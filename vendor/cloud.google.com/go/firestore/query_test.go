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
	"testing"

	"golang.org/x/net/context"

	"cloud.google.com/go/internal/pretty"
	pb "google.golang.org/genproto/googleapis/firestore/v1beta1"

	"github.com/golang/protobuf/ptypes/wrappers"
)

func TestQueryToProto(t *testing.T) {
	c := &Client{}
	coll := c.Collection("C")
	q := coll.Query
	aFilter, err := filter{[]string{"a"}, ">", 5}.toProto()
	if err != nil {
		t.Fatal(err)
	}
	bFilter, err := filter{[]string{"b"}, "<", "foo"}.toProto()
	if err != nil {
		t.Fatal(err)
	}
	slashStarFilter, err := filter{[]string{"/", "*"}, ">", 5}.toProto()
	if err != nil {
		t.Fatal(err)
	}
	type S struct {
		A int `firestore:"a"`
	}
	for _, test := range []struct {
		in   Query
		want *pb.StructuredQuery
	}{
		{
			in: q.Select(),
			want: &pb.StructuredQuery{
				Select: &pb.StructuredQuery_Projection{
					Fields: []*pb.StructuredQuery_FieldReference{fref1("__name__")},
				},
			},
		},
		{
			in: q.Select("a", "b"),
			want: &pb.StructuredQuery{
				Select: &pb.StructuredQuery_Projection{
					Fields: []*pb.StructuredQuery_FieldReference{fref1("a"), fref1("b")},
				},
			},
		},
		{
			in: q.Select("a", "b").Select("c"), // last wins
			want: &pb.StructuredQuery{
				Select: &pb.StructuredQuery_Projection{
					Fields: []*pb.StructuredQuery_FieldReference{fref1("c")},
				},
			},
		},
		{
			in: q.SelectPaths([]string{"*"}, []string{"/"}),
			want: &pb.StructuredQuery{
				Select: &pb.StructuredQuery_Projection{
					Fields: []*pb.StructuredQuery_FieldReference{fref1("*"), fref1("/")},
				},
			},
		},
		{
			in:   q.Where("a", ">", 5),
			want: &pb.StructuredQuery{Where: aFilter},
		},
		{
			in: q.Where("a", ">", 5).Where("b", "<", "foo"),
			want: &pb.StructuredQuery{
				Where: &pb.StructuredQuery_Filter{
					&pb.StructuredQuery_Filter_CompositeFilter{
						&pb.StructuredQuery_CompositeFilter{
							Op: pb.StructuredQuery_CompositeFilter_AND,
							Filters: []*pb.StructuredQuery_Filter{
								aFilter, bFilter,
							},
						},
					},
				},
			},
		},
		{
			in:   q.WherePath([]string{"/", "*"}, ">", 5),
			want: &pb.StructuredQuery{Where: slashStarFilter},
		},
		{
			in: q.OrderBy("b", Asc).OrderBy("a", Desc).OrderByPath([]string{"~"}, Asc),
			want: &pb.StructuredQuery{
				OrderBy: []*pb.StructuredQuery_Order{
					{fref1("b"), pb.StructuredQuery_ASCENDING},
					{fref1("a"), pb.StructuredQuery_DESCENDING},
					{fref1("~"), pb.StructuredQuery_ASCENDING},
				},
			},
		},
		{
			in: q.Offset(2).Limit(3),
			want: &pb.StructuredQuery{
				Offset: 2,
				Limit:  &wrappers.Int32Value{3},
			},
		},
		{
			in: q.Offset(2).Limit(3).Limit(4).Offset(5), // last wins
			want: &pb.StructuredQuery{
				Offset: 5,
				Limit:  &wrappers.Int32Value{4},
			},
		},
		{
			in: q.OrderBy("a", Asc).StartAt(7).EndBefore(9),
			want: &pb.StructuredQuery{
				OrderBy: []*pb.StructuredQuery_Order{
					{fref1("a"), pb.StructuredQuery_ASCENDING},
				},
				StartAt: &pb.Cursor{
					Values: []*pb.Value{intval(7)},
					Before: true,
				},
				EndAt: &pb.Cursor{
					Values: []*pb.Value{intval(9)},
					Before: true,
				},
			},
		},
		{
			in: q.OrderBy("a", Asc).StartAt(7).EndBefore(9),
			want: &pb.StructuredQuery{
				OrderBy: []*pb.StructuredQuery_Order{
					{fref1("a"), pb.StructuredQuery_ASCENDING},
				},
				StartAt: &pb.Cursor{
					Values: []*pb.Value{intval(7)},
					Before: true,
				},
				EndAt: &pb.Cursor{
					Values: []*pb.Value{intval(9)},
					Before: true,
				},
			},
		},
		{
			in: q.OrderBy("a", Asc).StartAfter(7).EndAt(9),
			want: &pb.StructuredQuery{
				OrderBy: []*pb.StructuredQuery_Order{
					{fref1("a"), pb.StructuredQuery_ASCENDING},
				},
				StartAt: &pb.Cursor{
					Values: []*pb.Value{intval(7)},
					Before: false,
				},
				EndAt: &pb.Cursor{
					Values: []*pb.Value{intval(9)},
					Before: false,
				},
			},
		},
		{
			in: q.OrderBy("a", Asc).OrderBy("b", Desc).StartAfter(7, 8).EndAt(9, 10),
			want: &pb.StructuredQuery{
				OrderBy: []*pb.StructuredQuery_Order{
					{fref1("a"), pb.StructuredQuery_ASCENDING},
					{fref1("b"), pb.StructuredQuery_DESCENDING},
				},
				StartAt: &pb.Cursor{
					Values: []*pb.Value{intval(7), intval(8)},
					Before: false,
				},
				EndAt: &pb.Cursor{
					Values: []*pb.Value{intval(9), intval(10)},
					Before: false,
				},
			},
		},
		{
			// last of StartAt/After wins, same for End
			in: q.OrderBy("a", Asc).
				StartAfter(1).StartAt(2).
				EndAt(3).EndBefore(4),
			want: &pb.StructuredQuery{
				OrderBy: []*pb.StructuredQuery_Order{
					{fref1("a"), pb.StructuredQuery_ASCENDING},
				},
				StartAt: &pb.Cursor{
					Values: []*pb.Value{intval(2)},
					Before: true,
				},
				EndAt: &pb.Cursor{
					Values: []*pb.Value{intval(4)},
					Before: true,
				},
			},
		},
	} {
		got, err := test.in.toProto()
		if err != nil {
			t.Fatalf("%+v: %v", test.in, err)
		}
		test.want.From = []*pb.StructuredQuery_CollectionSelector{{CollectionId: "C"}}
		if !testEqual(got, test.want) {
			t.Errorf("%+v: got\n%v\nwant\n%v", test.in, pretty.Value(got), pretty.Value(test.want))
		}
	}
}

func fref1(s string) *pb.StructuredQuery_FieldReference {
	return fref([]string{s})
}

func TestQueryToProtoErrors(t *testing.T) {
	q := (&Client{}).Collection("C").Query
	for _, query := range []Query{
		Query{},                                // no collection ID
		q.Where("x", "!=", 1),                  // invalid operator
		q.Where("~", ">", 1),                   // invalid path
		q.WherePath([]string{"*", ""}, ">", 1), // invalid path
		q.StartAt(1),                           // no OrderBy
		q.StartAt(2).OrderBy("x", Asc).OrderBy("y", Desc), // wrong # OrderBy
		q.Select("*"),                                     // invalid path
		q.SelectPaths([]string{"/", "", "~"}),             // invalid path
		q.OrderBy("[", Asc),                               // invalid path
		q.OrderByPath([]string{""}, Desc),                 // invalid path
	} {
		_, err := query.toProto()
		if err == nil {
			t.Errorf("%+v: got nil, want error", query)
		}
	}
}

func TestQueryMethodsDoNotModifyReceiver(t *testing.T) {
	var empty Query

	q := Query{}
	_ = q.Select("a", "b")
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}

	q = Query{}
	q1 := q.Where("a", ">", 3)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}
	// Extra check because Where appends to a slice.
	q1before := q.Where("a", ">", 3) // same as q1
	_ = q1.Where("b", "<", "foo")
	if !testEqual(q1, q1before) {
		t.Errorf("got %+v, want %+v", q1, q1before)
	}

	q = Query{}
	q1 = q.OrderBy("a", Asc)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}
	// Extra check because Where appends to a slice.
	q1before = q.OrderBy("a", Asc) // same as q1
	_ = q1.OrderBy("b", Desc)
	if !testEqual(q1, q1before) {
		t.Errorf("got %+v, want %+v", q1, q1before)
	}

	q = Query{}
	_ = q.Offset(5)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}

	q = Query{}
	_ = q.Limit(5)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}

	q = Query{}
	_ = q.StartAt(5)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}

	q = Query{}
	_ = q.StartAfter(5)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}

	q = Query{}
	_ = q.EndAt(5)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}

	q = Query{}
	_ = q.EndBefore(5)
	if !testEqual(q, empty) {
		t.Errorf("got %+v, want empty", q)
	}
}

func TestQueryFromCollectionRef(t *testing.T) {
	c := &Client{}
	coll := c.Collection("C")
	got := coll.Select("x").Offset(8)
	want := Query{
		c:            c,
		parentPath:   c.path(),
		collectionID: "C",
		selection:    []FieldPath{{"x"}},
		offset:       8,
	}
	if !testEqual(got, want) {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestQueryGetAll(t *testing.T) {
	// This implicitly tests DocumentIterator as well.
	const dbPath = "projects/projectID/databases/(default)"
	ctx := context.Background()
	c, srv := newMock(t)
	docNames := []string{"C/a", "C/b"}
	wantPBDocs := []*pb.Document{
		{
			Name:       dbPath + "/documents/" + docNames[0],
			CreateTime: aTimestamp,
			UpdateTime: aTimestamp,
			Fields:     map[string]*pb.Value{"f": intval(2)},
		},
		{
			Name:       dbPath + "/documents/" + docNames[1],
			CreateTime: aTimestamp2,
			UpdateTime: aTimestamp3,
			Fields:     map[string]*pb.Value{"f": intval(1)},
		},
	}

	srv.addRPC(nil, []interface{}{
		&pb.RunQueryResponse{Document: wantPBDocs[0]},
		&pb.RunQueryResponse{Document: wantPBDocs[1]},
	})
	gotDocs, err := c.Collection("C").Documents(ctx).GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(gotDocs), len(wantPBDocs); got != want {
		t.Errorf("got %d docs, wanted %d", got, want)
	}
	for i, got := range gotDocs {
		want, err := newDocumentSnapshot(c.Doc(docNames[i]), wantPBDocs[i], c)
		if err != nil {
			t.Fatal(err)
		}
		if !testEqual(got, want) {
			// avoid writing a cycle
			got.c = nil
			want.c = nil
			t.Errorf("#%d: got %+v, want %+v", i, pretty.Value(got), pretty.Value(want))
		}
	}
}
