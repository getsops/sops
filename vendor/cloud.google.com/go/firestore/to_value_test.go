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

package firestore

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	ts "github.com/golang/protobuf/ptypes/timestamp"
	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/genproto/googleapis/type/latlng"
)

type testStruct1 struct {
	B  bool
	I  int
	U  uint32
	F  float64
	S  string
	Y  []byte
	T  time.Time
	Ts *ts.Timestamp
	G  *latlng.LatLng
	L  []int
	M  map[string]int
	P  *int
}

var (
	p = new(int)

	testVal1 = testStruct1{
		B:  true,
		I:  1,
		U:  2,
		F:  3.0,
		S:  "four",
		Y:  []byte{5},
		T:  tm,
		Ts: ptm,
		G:  ll,
		L:  []int{6},
		M:  map[string]int{"a": 7},
		P:  p,
	}

	mapVal1 = mapval(map[string]*pb.Value{
		"B":  boolval(true),
		"I":  intval(1),
		"U":  intval(2),
		"F":  floatval(3),
		"S":  {ValueType: &pb.Value_StringValue{"four"}},
		"Y":  bytesval([]byte{5}),
		"T":  tsval(tm),
		"Ts": {ValueType: &pb.Value_TimestampValue{ptm}},
		"G":  geoval(ll),
		"L":  arrayval(intval(6)),
		"M":  mapval(map[string]*pb.Value{"a": intval(7)}),
		"P":  intval(8),
	})
)

// TODO descriptions
// TODO cause the array failure
func TestToProtoValue_Conversions(t *testing.T) {
	*p = 8
	for _, test := range []struct {
		desc string
		in   interface{}
		want *pb.Value
	}{
		{
			desc: "nil",
			in:   nil,
			want: nullValue,
		},
		{
			desc: "nil slice",
			in:   []int(nil),
			want: nullValue,
		},
		{
			desc: "nil map",
			in:   map[string]int(nil),
			want: nullValue,
		},
		{
			desc: "nil struct",
			in:   (*testStruct1)(nil),
			want: nullValue,
		},
		{
			desc: "nil timestamp",
			in:   (*ts.Timestamp)(nil),
			want: nullValue,
		},
		{
			desc: "nil latlng",
			in:   (*latlng.LatLng)(nil),
			want: nullValue,
		},
		{
			desc: "nil docref",
			in:   (*DocumentRef)(nil),
			want: nullValue,
		},
		{
			desc: "bool",
			in:   true,
			want: boolval(true),
		},
		{
			desc: "int",
			in:   3,
			want: intval(3),
		},
		{
			desc: "uint32",
			in:   uint32(3),
			want: intval(3),
		},
		{
			desc: "float",
			in:   1.5,
			want: floatval(1.5),
		},
		{
			desc: "string",
			in:   "str",
			want: strval("str"),
		},
		{
			desc: "byte slice",
			in:   []byte{1, 2},
			want: bytesval([]byte{1, 2}),
		},
		{
			desc: "date time",
			in:   tm,
			want: tsval(tm),
		},
		{
			desc: "pointer to timestamp",
			in:   ptm,
			want: &pb.Value{ValueType: &pb.Value_TimestampValue{ptm}},
		},
		{
			desc: "pointer to latlng",
			in:   ll,
			want: geoval(ll),
		},
		{
			desc: "populated slice",
			in:   []int{1, 2},
			want: arrayval(intval(1), intval(2)),
		},
		{
			desc: "pointer to populated slice",
			in:   &[]int{1, 2},
			want: arrayval(intval(1), intval(2)),
		},
		{
			desc: "empty slice",
			in:   []int{},
			want: arrayval(),
		},
		{
			desc: "populated map",
			in:   map[string]int{"a": 1, "b": 2},
			want: mapval(map[string]*pb.Value{"a": intval(1), "b": intval(2)}),
		},
		{
			desc: "empty map",
			in:   map[string]int{},
			want: mapval(map[string]*pb.Value{}),
		},
		{
			desc: "int",
			in:   p,
			want: intval(8),
		},
		{
			desc: "pointer to int",
			in:   &p,
			want: intval(8),
		},
		{
			desc: "populated map",
			in:   map[string]interface{}{"a": 1, "p": p, "s": "str"},
			want: mapval(map[string]*pb.Value{"a": intval(1), "p": intval(8), "s": strval("str")}),
		},
		{
			desc: "map with timestamp",
			in:   map[string]fmt.Stringer{"a": tm},
			want: mapval(map[string]*pb.Value{"a": tsval(tm)}),
		},
		{
			desc: "struct",
			in:   testVal1,
			want: mapVal1,
		},
		{
			desc: "array",
			in:   [1]int{7},
			want: arrayval(intval(7)),
		},
		{
			desc: "pointer to docref",
			in: &DocumentRef{
				ID:   "d",
				Path: "projects/P/databases/D/documents/c/d",
				Parent: &CollectionRef{
					ID:         "c",
					parentPath: "projects/P/databases/D",
					Path:       "projects/P/databases/D/documents/c",
					Query:      Query{collectionID: "c", parentPath: "projects/P/databases/D"},
				},
			},
			want: refval("projects/P/databases/D/documents/c/d"),
		},
		{
			desc: "Transforms are removed, which can lead to leaving nil",
			in:   map[string]interface{}{"a": ServerTimestamp},
			want: nil,
		},
		{
			desc: "Transform nested in map is ignored",
			in: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": ServerTimestamp,
					},
				},
			},
			want: nil,
		},
		{
			desc: "Transforms nested in map are ignored",
			in: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": ServerTimestamp,
						"d": ServerTimestamp,
					},
				},
			},
			want: nil,
		},
		{
			desc: "int nested in map is kept whilst Transforms are ignored",
			in: map[string]interface{}{
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": ServerTimestamp,
						"d": ServerTimestamp,
						"e": 1,
					},
				},
			},
			want: mapval(map[string]*pb.Value{
				"a": mapval(map[string]*pb.Value{
					"b": mapval(map[string]*pb.Value{"e": intval(1)}),
				}),
			}),
		},

		// Transforms are allowed in maps, but won't show up in the returned proto. Instead, we rely
		// on seeing sawTransforms=true and a call to extractTransforms.
		{
			desc: "Transforms in map are ignored, other values are kept (ServerTimestamp)",
			in:   map[string]interface{}{"a": ServerTimestamp, "b": 5},
			want: mapval(map[string]*pb.Value{"b": intval(5)}),
		},
		{
			desc: "Transforms in map are ignored, other values are kept (ArrayUnion)",
			in:   map[string]interface{}{"a": ArrayUnion(1, 2, 3), "b": 5},
			want: mapval(map[string]*pb.Value{"b": intval(5)}),
		},
		{
			desc: "Transforms in map are ignored, other values are kept (ArrayRemove)",
			in:   map[string]interface{}{"a": ArrayRemove(1, 2, 3), "b": 5},
			want: mapval(map[string]*pb.Value{"b": intval(5)}),
		},
	} {
		t.Run(test.desc, func(t *testing.T) {
			got, _, err := toProtoValue(reflect.ValueOf(test.in))
			if err != nil {
				t.Fatalf("%v (%T): %v", test.in, test.in, err)
			}
			if !testEqual(got, test.want) {
				t.Fatalf("%+v (%T):\ngot\n%+v\nwant\n%+v", test.in, test.in, got, test.want)
			}
		})
	}
}

type stringy struct{}

func (stringy) String() string { return "stringy" }

func TestToProtoValue_Errors(t *testing.T) {
	for _, in := range []interface{}{
		uint64(0),                               // a bad fit for int64
		map[int]bool{},                          // map key type is not string
		make(chan int),                          // can't handle type
		map[string]fmt.Stringer{"a": stringy{}}, // only empty interfaces
		ServerTimestamp,                         // ServerTimestamp can only be a field value
		struct{ A interface{} }{A: ServerTimestamp},
		map[string]interface{}{"a": []interface{}{ServerTimestamp}},
		map[string]interface{}{"a": []interface{}{
			map[string]interface{}{"b": ServerTimestamp},
		}},
		Delete, // Delete should never appear
		[]interface{}{Delete},
		map[string]interface{}{"a": Delete},
		map[string]interface{}{"a": []interface{}{Delete}},

		// Transforms are not allowed to occur in an array.
		[]interface{}{ServerTimestamp},
		[]interface{}{ArrayUnion(1, 2, 3)},
		[]interface{}{ArrayRemove(1, 2, 3)},

		// Transforms are not allowed to occur in a struct.
		struct{ A interface{} }{A: ServerTimestamp},
		struct{ A interface{} }{A: ArrayUnion()},
		struct{ A interface{} }{A: ArrayRemove()},
	} {
		_, _, err := toProtoValue(reflect.ValueOf(in))
		if err == nil {
			t.Errorf("%v: got nil, want error", in)
		}
	}
}

func TestToProtoValue_SawTransform(t *testing.T) {
	for i, in := range []interface{}{
		map[string]interface{}{"a": ServerTimestamp},
		map[string]interface{}{"a": ArrayUnion()},
		map[string]interface{}{"a": ArrayRemove()},
	} {
		_, sawTransform, err := toProtoValue(reflect.ValueOf(in))
		if err != nil {
			t.Fatalf("%d %v: got err %v\nexpected nil", i, in, err)
		}
		if !sawTransform {
			t.Errorf("%d %v: got sawTransform=false, expected sawTransform=true", i, in)
		}
	}
}

type testStruct2 struct {
	Ignore        int       `firestore:"-"`
	Rename        int       `firestore:"a"`
	OmitEmpty     int       `firestore:",omitempty"`
	OmitEmptyTime time.Time `firestore:",omitempty"`
}

func TestToProtoValue_Tags(t *testing.T) {
	in := &testStruct2{
		Ignore:        1,
		Rename:        2,
		OmitEmpty:     3,
		OmitEmptyTime: aTime,
	}
	got, _, err := toProtoValue(reflect.ValueOf(in))
	if err != nil {
		t.Fatal(err)
	}
	want := mapval(map[string]*pb.Value{
		"a":             intval(2),
		"OmitEmpty":     intval(3),
		"OmitEmptyTime": tsval(aTime),
	})
	if !testEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}

	got, _, err = toProtoValue(reflect.ValueOf(testStruct2{}))
	if err != nil {
		t.Fatal(err)
	}
	want = mapval(map[string]*pb.Value{"a": intval(0)})
	if !testEqual(got, want) {
		t.Errorf("got\n%+v\nwant\n%+v", got, want)
	}
}

func TestToProtoValue_Embedded(t *testing.T) {
	// Embedded time.Time, LatLng, or Timestamp should behave like non-embedded.
	type embed struct {
		time.Time
		*latlng.LatLng
		*ts.Timestamp
	}

	got, _, err := toProtoValue(reflect.ValueOf(embed{tm, ll, ptm}))
	if err != nil {
		t.Fatal(err)
	}
	want := mapval(map[string]*pb.Value{
		"Time":      tsval(tm),
		"LatLng":    geoval(ll),
		"Timestamp": {ValueType: &pb.Value_TimestampValue{ptm}},
	})
	if !testEqual(got, want) {
		t.Errorf("got %+v, want %+v", got, want)
	}
}

func TestIsEmpty(t *testing.T) {
	for _, e := range []interface{}{int(0), float32(0), false, "", []int{}, []int(nil), (*int)(nil)} {
		if !isEmptyValue(reflect.ValueOf(e)) {
			t.Errorf("%v (%T): want true, got false", e, e)
		}
	}
	i := 3
	for _, n := range []interface{}{int(1), float32(1), true, "x", []int{1}, &i} {
		if isEmptyValue(reflect.ValueOf(n)) {
			t.Errorf("%v (%T): want false, got true", n, n)
		}
	}
}
