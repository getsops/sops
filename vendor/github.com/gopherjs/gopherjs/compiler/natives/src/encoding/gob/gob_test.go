// +build js

package gob

import (
	"bytes"
	"reflect"
	"testing"
)

// TODO: TestEndToEnd override can be removed once the bug with Marr field is fixed.
func TestEndToEnd(t *testing.T) {
	type T2 struct {
		T string
	}
	type T3 struct {
		X float64
		Z *int
	}
	type T1 struct {
		A, B, C  int
		M        map[string]*float64
		M2       map[int]T3
		Mstring  map[string]string
		Mintptr  map[int]*int
		Mcomp    map[complex128]complex128
		Marr     map[[2]string][2]*float64
		EmptyMap map[string]int // to check that we receive a non-nil map.
		N        *[3]float64
		Strs     *[2]string
		Int64s   *[]int64
		RI       complex64
		S        string
		Y        []byte
		T        *T2
	}
	pi := 3.14159
	e := 2.71828
	two := 2.0
	meaning := 42
	fingers := 5
	s1 := "string1"
	s2 := "string2"
	var comp1 complex128 = complex(1.0, 1.0)
	var comp2 complex128 = complex(1.0, 1.0)
	var arr1 [2]string
	arr1[0] = s1
	arr1[1] = s2
	var arr2 [2]string
	arr2[0] = s2
	arr2[1] = s1
	var floatArr1 [2]*float64
	floatArr1[0] = &pi
	floatArr1[1] = &e
	var floatArr2 [2]*float64
	floatArr2[0] = &e
	floatArr2[1] = &two
	t1 := &T1{
		A:       17,
		B:       18,
		C:       -5,
		M:       map[string]*float64{"pi": &pi, "e": &e},
		M2:      map[int]T3{4: T3{X: pi, Z: &meaning}, 10: T3{X: e, Z: &fingers}},
		Mstring: map[string]string{"pi": "3.14", "e": "2.71"},
		Mintptr: map[int]*int{meaning: &fingers, fingers: &meaning},
		Mcomp:   map[complex128]complex128{comp1: comp2, comp2: comp1},
		// TODO: Fix this problem:
		// 	TypeError: dst.$set is not a function
		// 	    at typedmemmove (/github.com/gopherjs/gopherjs/reflect.go:487:3)
		//Marr:     map[[2]string][2]*float64{arr1: floatArr1, arr2: floatArr2},
		EmptyMap: make(map[string]int),
		N:        &[3]float64{1.5, 2.5, 3.5},
		Strs:     &[2]string{s1, s2},
		Int64s:   &[]int64{77, 89, 123412342134},
		RI:       17 - 23i,
		S:        "Now is the time",
		Y:        []byte("hello, sailor"),
		T:        &T2{"this is T2"},
	}
	b := new(bytes.Buffer)
	err := NewEncoder(b).Encode(t1)
	if err != nil {
		t.Error("encode:", err)
	}
	var _t1 T1
	err = NewDecoder(b).Decode(&_t1)
	if err != nil {
		t.Fatal("decode:", err)
	}
	if !reflect.DeepEqual(t1, &_t1) {
		t.Errorf("encode expected %v got %v", *t1, _t1)
	}
	// Be absolutely sure the received map is non-nil.
	if t1.EmptyMap == nil {
		t.Errorf("nil map sent")
	}
	if _t1.EmptyMap == nil {
		t.Errorf("nil map received")
	}
}

func TestTypeRace(t *testing.T) {
	// encoding/gob currently uses nosync. This test uses sync.WaitGroup and
	// cannot succeed when nosync is used.
	t.Skip("using nosync")
}
