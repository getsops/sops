// Copyright 2016 Google Inc. All Rights Reserved.
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

package datastore

import (
	"reflect"
	"testing"

	pb "google.golang.org/genproto/googleapis/datastore/v1"
)

type Simple struct {
	I int64
}

type SimpleWithTag struct {
	I int64 `datastore:"II"`
}

type NestedSimpleWithTag struct {
	A SimpleWithTag `datastore:"AA"`
}

type NestedSliceOfSimple struct {
	A []Simple
}

type SimpleTwoFields struct {
	S  string
	SS string
}

type NestedSimpleAnonymous struct {
	Simple
	X string
}

type NestedSimple struct {
	A Simple
	I int
}

type NestedSimple1 struct {
	A Simple
	X string
}

type NestedSimple2X struct {
	AA NestedSimple
	A  SimpleTwoFields
	S  string
}

type BDotB struct {
	B string `datastore:"B.B"`
}

type ABDotB struct {
	A BDotB
}

type MultiAnonymous struct {
	Simple
	SimpleTwoFields
	X string
}

func TestLoadEntityNestedLegacy(t *testing.T) {
	testCases := []struct {
		desc string
		src  *pb.Entity
		want interface{}
	}{
		{
			desc: "nested",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"X":   {ValueType: &pb.Value_StringValue{StringValue: "two"}},
					"A.I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 2}},
				},
			},
			want: &NestedSimple1{
				A: Simple{I: 2},
				X: "two",
			},
		},
		{
			desc: "nested with tag",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"AA.II": {ValueType: &pb.Value_IntegerValue{IntegerValue: 2}},
				},
			},
			want: &NestedSimpleWithTag{
				A: SimpleWithTag{I: 2},
			},
		},
		{
			desc: "nested with anonymous struct field",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"X": {ValueType: &pb.Value_StringValue{StringValue: "two"}},
					"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 2}},
				},
			},
			want: &NestedSimpleAnonymous{
				Simple: Simple{I: 2},
				X:      "two",
			},
		},
		{
			desc: "nested with dotted field tag",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"A.B.B": {ValueType: &pb.Value_StringValue{StringValue: "bb"}},
				},
			},
			want: &ABDotB{
				A: BDotB{
					B: "bb",
				},
			},
		},
		{
			desc: "nested with multiple anonymous fields",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"I":  {ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
					"S":  {ValueType: &pb.Value_StringValue{StringValue: "S"}},
					"SS": {ValueType: &pb.Value_StringValue{StringValue: "s"}},
					"X":  {ValueType: &pb.Value_StringValue{StringValue: "s"}},
				},
			},
			want: &MultiAnonymous{
				Simple:          Simple{I: 3},
				SimpleTwoFields: SimpleTwoFields{S: "S", SS: "s"},
				X:               "s",
			},
		},
	}

	for _, tc := range testCases {
		dst := reflect.New(reflect.TypeOf(tc.want).Elem()).Interface()
		err := loadEntityProto(dst, tc.src)
		if err != nil {
			t.Errorf("loadEntityProto: %s: %v", tc.desc, err)
			continue
		}

		if !reflect.DeepEqual(tc.want, dst) {
			t.Errorf("%s: compare:\ngot:  %#v\nwant: %#v", tc.desc, dst, tc.want)
		}
	}
}

type WithKey struct {
	X string
	I int
	K *Key `datastore:"__key__"`
}

type NestedWithKey struct {
	Y string
	N WithKey
}

var (
	incompleteKey = newKey("", nil)
	invalidKey    = newKey("s", incompleteKey)
)

func TestLoadEntityNested(t *testing.T) {
	testCases := []struct {
		desc string
		src  *pb.Entity
		want interface{}
	}{
		{
			desc: "nested basic",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"A": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
							},
						},
					}},
					"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 10}},
				},
			},
			want: &NestedSimple{
				A: Simple{I: 3},
				I: 10,
			},
		},
		{
			desc: "nested with struct tags",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"AA": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"II": {ValueType: &pb.Value_IntegerValue{IntegerValue: 1}},
							},
						},
					}},
				},
			},
			want: &NestedSimpleWithTag{
				A: SimpleWithTag{I: 1},
			},
		},
		{
			desc: "nested 2x",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"AA": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"A": {ValueType: &pb.Value_EntityValue{
									EntityValue: &pb.Entity{
										Properties: map[string]*pb.Value{
											"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
										},
									},
								}},
								"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 1}},
							},
						},
					}},
					"A": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"S":  {ValueType: &pb.Value_StringValue{StringValue: "S"}},
								"SS": {ValueType: &pb.Value_StringValue{StringValue: "s"}},
							},
						},
					}},
					"S": {ValueType: &pb.Value_StringValue{StringValue: "SS"}},
				},
			},
			want: &NestedSimple2X{
				AA: NestedSimple{
					A: Simple{I: 3},
					I: 1,
				},
				A: SimpleTwoFields{S: "S", SS: "s"},
				S: "SS",
			},
		},
		{
			desc: "nested anonymous",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
					"X": {ValueType: &pb.Value_StringValue{StringValue: "SomeX"}},
				},
			},
			want: &NestedSimpleAnonymous{
				Simple: Simple{I: 3},
				X:      "SomeX",
			},
		},
		{
			desc: "nested simple with slice",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"A": {ValueType: &pb.Value_ArrayValue{
						ArrayValue: &pb.ArrayValue{
							Values: []*pb.Value{
								{ValueType: &pb.Value_EntityValue{
									EntityValue: &pb.Entity{
										Properties: map[string]*pb.Value{
											"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
										},
									},
								}},
								{ValueType: &pb.Value_EntityValue{
									EntityValue: &pb.Entity{
										Properties: map[string]*pb.Value{
											"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 4}},
										},
									},
								}},
							},
						},
					}},
				},
			},

			want: &NestedSliceOfSimple{
				A: []Simple{Simple{I: 3}, Simple{I: 4}},
			},
		},
		{
			desc: "nested with multiple anonymous fields",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"I":  {ValueType: &pb.Value_IntegerValue{IntegerValue: 3}},
					"S":  {ValueType: &pb.Value_StringValue{StringValue: "S"}},
					"SS": {ValueType: &pb.Value_StringValue{StringValue: "s"}},
					"X":  {ValueType: &pb.Value_StringValue{StringValue: "ss"}},
				},
			},
			want: &MultiAnonymous{
				Simple:          Simple{I: 3},
				SimpleTwoFields: SimpleTwoFields{S: "S", SS: "s"},
				X:               "ss",
			},
		},
		{
			desc: "nested with dotted field tag",
			src: &pb.Entity{
				Properties: map[string]*pb.Value{
					"A": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"B.B": {ValueType: &pb.Value_StringValue{StringValue: "bb"}},
							},
						},
					}},
				},
			},
			want: &ABDotB{
				A: BDotB{
					B: "bb",
				},
			},
		},
		{
			desc: "nested entity with key",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"Y": {ValueType: &pb.Value_StringValue{StringValue: "yyy"}},
					"N": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Key: keyToProto(testKey1a),
							Properties: map[string]*pb.Value{
								"X": {ValueType: &pb.Value_StringValue{StringValue: "two"}},
								"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 2}},
							},
						},
					}},
				},
			},
			want: &NestedWithKey{
				Y: "yyy",
				N: WithKey{
					X: "two",
					I: 2,
					K: testKey1a,
				},
			},
		},
		{
			desc: "nested entity with invalid key",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"Y": {ValueType: &pb.Value_StringValue{StringValue: "yyy"}},
					"N": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Key: keyToProto(invalidKey),
							Properties: map[string]*pb.Value{
								"X": {ValueType: &pb.Value_StringValue{StringValue: "two"}},
								"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 2}},
							},
						},
					}},
				},
			},
			want: &NestedWithKey{
				Y: "yyy",
				N: WithKey{
					X: "two",
					I: 2,
					K: invalidKey,
				},
			},
		},
	}

	for _, tc := range testCases {
		dst := reflect.New(reflect.TypeOf(tc.want).Elem()).Interface()
		err := loadEntityProto(dst, tc.src)
		if err != nil {
			t.Errorf("loadEntityProto: %s: %v", tc.desc, err)
			continue
		}

		if !reflect.DeepEqual(tc.want, dst) {
			t.Errorf("%s: compare:\ngot:  %#v\nwant: %#v", tc.desc, dst, tc.want)
		}
	}
}

type NestedStructPtrs struct {
	*SimpleTwoFields
	Nest      *SimpleTwoFields
	TwiceNest *NestedSimple2
	I         int
}

type NestedSimple2 struct {
	A *Simple
	I int
}

func TestAlreadyPopulatedDst(t *testing.T) {
	testCases := []struct {
		desc string
		src  *pb.Entity
		dst  interface{}
		want interface{}
	}{
		{
			desc: "simple already populated, nil properties",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"I": {ValueType: &pb.Value_NullValue{}},
				},
			},
			dst: &Simple{
				I: 12,
			},
			want: &Simple{},
		},
		{
			desc: "nested structs already populated",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"SS": {ValueType: &pb.Value_StringValue{StringValue: "world"}},
				},
			},
			dst:  &SimpleTwoFields{S: "hello" /* SS: "" */},
			want: &SimpleTwoFields{S: "hello", SS: "world"},
		},
		{
			desc: "nested structs already populated, pValues nil",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"S":    {ValueType: &pb.Value_NullValue{}},
					"SS":   {ValueType: &pb.Value_StringValue{StringValue: "ss hello"}},
					"Nest": {ValueType: &pb.Value_NullValue{}},
					"TwiceNest": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"A": {ValueType: &pb.Value_NullValue{}},
								"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 2}},
							},
						},
					}},
					"I": {ValueType: &pb.Value_IntegerValue{IntegerValue: 5}},
				},
			},
			dst: &NestedStructPtrs{
				&SimpleTwoFields{S: "hello" /* SS: "" */},
				&SimpleTwoFields{ /* S: "" */ SS: "twice hello"},
				&NestedSimple2{
					A: &Simple{I: 2},
					/* I: 0 */
				},
				0,
			},
			want: &NestedStructPtrs{
				&SimpleTwoFields{ /* S: "" */ SS: "ss hello"},
				nil,
				&NestedSimple2{
					/* A: nil, */
					I: 2,
				},
				5,
			},
		},
	}

	for _, tc := range testCases {
		err := loadEntityProto(tc.dst, tc.src)
		if err != nil {
			t.Errorf("loadEntityProto: %s: %v", tc.desc, err)
			continue
		}

		if !reflect.DeepEqual(tc.want, tc.dst) {
			t.Errorf("%s: compare:\ngot:  %#v\nwant: %#v", tc.desc, tc.dst, tc.want)
		}
	}
}

type PLS0 struct {
	A string
}

func (p *PLS0) Load(props []Property) error {
	for _, pp := range props {
		if pp.Name == "A" {
			p.A = pp.Value.(string)
		}
	}
	return nil
}

func (p *PLS0) Save() (props []Property, err error) {
	return []Property{{Name: "A", Value: p.A}}, nil
}

type KeyLoader1 struct {
	A string
	K *Key
}

func (kl *KeyLoader1) Load(props []Property) error {
	for _, pp := range props {
		if pp.Name == "A" {
			kl.A = pp.Value.(string)
		}
	}
	return nil
}

func (kl *KeyLoader1) Save() (props []Property, err error) {
	return []Property{{Name: "A", Value: kl.A}}, nil
}

func (kl *KeyLoader1) LoadKey(k *Key) error {
	kl.K = k
	return nil
}

type KeyLoader2 struct {
	B   int
	Key *Key
}

func (kl *KeyLoader2) Load(props []Property) error {
	for _, pp := range props {
		if pp.Name == "B" {
			kl.B = int(pp.Value.(int64))
		}
	}
	return nil
}

func (kl *KeyLoader2) Save() (props []Property, err error) {
	return []Property{{Name: "B", Value: int64(kl.B)}}, nil
}

func (kl *KeyLoader2) LoadKey(k *Key) error {
	kl.Key = k
	return nil
}

type KeyLoader3 struct {
	C bool
	K *Key
}

func (kl *KeyLoader3) Load(props []Property) error {
	for _, pp := range props {
		if pp.Name == "C" {
			kl.C = pp.Value.(bool)
		}
	}
	return nil
}

func (kl *KeyLoader3) Save() (props []Property, err error) {
	return []Property{{Name: "C", Value: kl.C}}, nil
}

func (kl *KeyLoader3) LoadKey(k *Key) error {
	kl.K = k
	return nil
}

type KeyLoader4 struct {
	PLS0
	K *Key
}

func (kl *KeyLoader4) LoadKey(k *Key) error {
	kl.K = k
	return nil
}

type NotKeyLoader struct {
	A string
	K *Key
}

func (p *NotKeyLoader) Load(props []Property) error {
	for _, pp := range props {
		if pp.Name == "A" {
			p.A = pp.Value.(string)
		}
	}
	return nil
}

func (p *NotKeyLoader) Save() (props []Property, err error) {
	return []Property{{Name: "A", Value: p.A}}, nil
}

type NestedKeyLoaders struct {
	Two   *KeyLoader2
	Three []*KeyLoader3
	Four  *KeyLoader4
	PLS   *NotKeyLoader
}

func TestKeyLoader(t *testing.T) {
	testCases := []struct {
		desc string
		src  *pb.Entity
		dst  interface{}
		want interface{}
	}{
		{
			desc: "simple key loader",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"A": {ValueType: &pb.Value_StringValue{StringValue: "hello"}},
				},
			},
			dst: &KeyLoader1{},
			want: &KeyLoader1{
				A: "hello",
				K: testKey0,
			},
		},
		{
			desc: "embedded PLS key loader",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"A": {ValueType: &pb.Value_StringValue{StringValue: "hello"}},
				},
			},
			dst: &KeyLoader4{},
			want: &KeyLoader4{
				PLS0: PLS0{A: "hello"},
				K:    testKey0,
			},
		},
		{
			desc: "nested key loaders",
			src: &pb.Entity{
				Key: keyToProto(testKey0),
				Properties: map[string]*pb.Value{
					"Two": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"B": {ValueType: &pb.Value_IntegerValue{IntegerValue: 12}},
							},
							Key: keyToProto(testKey1a),
						},
					}},
					"Three": {ValueType: &pb.Value_ArrayValue{
						ArrayValue: &pb.ArrayValue{
							Values: []*pb.Value{
								{ValueType: &pb.Value_EntityValue{
									EntityValue: &pb.Entity{
										Properties: map[string]*pb.Value{
											"C": {ValueType: &pb.Value_BooleanValue{BooleanValue: true}},
										},
										Key: keyToProto(testKey1b),
									},
								}},
								{ValueType: &pb.Value_EntityValue{
									EntityValue: &pb.Entity{
										Properties: map[string]*pb.Value{
											"C": {ValueType: &pb.Value_BooleanValue{BooleanValue: false}},
										},
										Key: keyToProto(testKey0),
									},
								}},
							},
						},
					}},
					"Four": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"A": {ValueType: &pb.Value_StringValue{StringValue: "testing"}},
							},
							Key: keyToProto(testKey2a),
						},
					}},
					"PLS": {ValueType: &pb.Value_EntityValue{
						EntityValue: &pb.Entity{
							Properties: map[string]*pb.Value{
								"A": {ValueType: &pb.Value_StringValue{StringValue: "something"}},
							},

							Key: keyToProto(testKey1a),
						},
					}},
				},
			},
			dst: &NestedKeyLoaders{},
			want: &NestedKeyLoaders{
				Two: &KeyLoader2{B: 12, Key: testKey1a},
				Three: []*KeyLoader3{
					{
						C: true,
						K: testKey1b,
					},
					{
						C: false,
						K: testKey0,
					},
				},
				Four: &KeyLoader4{
					PLS0: PLS0{A: "testing"},
					K:    testKey2a,
				},
				PLS: &NotKeyLoader{A: "something"},
			},
		},
	}

	for _, tc := range testCases {
		err := loadEntityProto(tc.dst, tc.src)
		if err != nil {
			t.Errorf("loadEntityProto: %s: %v", tc.desc, err)
			continue
		}

		if !reflect.DeepEqual(tc.want, tc.dst) {
			t.Errorf("%s: compare:\ngot:  %+v\nwant: %+v", tc.desc, tc.dst, tc.want)
		}
	}
}
