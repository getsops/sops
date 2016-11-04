package tests

import (
	"math"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
	"vendored"

	"github.com/gopherjs/gopherjs/tests/otherpkg"
)

func TestSyntax1(t *testing.T) {
	a := 42
	if *&*&a != 42 {
		t.Fail()
	}
}

func TestPointerEquality(t *testing.T) {
	a := 1
	b := 1
	if &a != &a || &a == &b {
		t.Fail()
	}
	m := make(map[*int]int)
	m[&a] = 2
	m[&b] = 3
	if m[&a] != 2 || m[&b] != 3 {
		t.Fail()
	}

	for {
		c := 1
		d := 1
		if &c != &c || &c == &d {
			t.Fail()
		}
		break
	}

	s := struct {
		e int
		f int
	}{1, 1}
	if &s.e != &s.e || &s.e == &s.f {
		t.Fail()
	}

	g := [3]int{1, 2, 3}
	if &g[0] != &g[0] || &g[:][0] != &g[0] || &g[:][0] != &g[:][0] {
		t.Fail()
	}
}

type SingleValue struct {
	Value uint16
}

type OtherSingleValue struct {
	Value uint16
}

func TestStructKey(t *testing.T) {
	m := make(map[SingleValue]int)
	m[SingleValue{Value: 1}] = 42
	m[SingleValue{Value: 2}] = 43
	if m[SingleValue{Value: 1}] != 42 || m[SingleValue{Value: 2}] != 43 || reflect.ValueOf(m).MapIndex(reflect.ValueOf(SingleValue{Value: 1})).Interface() != 42 {
		t.Fail()
	}

	m2 := make(map[interface{}]int)
	m2[SingleValue{Value: 1}] = 42
	m2[SingleValue{Value: 2}] = 43
	m2[OtherSingleValue{Value: 1}] = 44
	if m2[SingleValue{Value: 1}] != 42 || m2[SingleValue{Value: 2}] != 43 || m2[OtherSingleValue{Value: 1}] != 44 || reflect.ValueOf(m2).MapIndex(reflect.ValueOf(SingleValue{Value: 1})).Interface() != 42 {
		t.Fail()
	}
}

func TestSelectOnNilChan(t *testing.T) {
	var c1 chan bool
	c2 := make(chan bool)

	go func() {
		close(c2)
	}()

	select {
	case <-c1:
		t.Fail()
	case <-c2:
		// ok
	}
}

type StructA struct {
	x int
}

type StructB struct {
	StructA
}

func TestEmbeddedStruct(t *testing.T) {
	a := StructA{
		42,
	}
	b := StructB{
		StructA: a,
	}
	b.x = 0
	if a.x != 42 {
		t.Fail()
	}
}

func TestMapStruct(t *testing.T) {
	a := StructA{
		42,
	}
	m := map[int]StructA{
		1: a,
	}
	m[2] = a
	a.x = 0
	if m[1].x != 42 || m[2].x != 42 {
		t.Fail()
	}
}

func TestUnnamedParameters(t *testing.T) {
	ok := false
	defer func() {
		if !ok {
			t.Fail()
		}
	}()
	blockingWithUnnamedParameter(false) // used to cause non-blocking call error, which is ignored by testing
	ok = true
}

func blockingWithUnnamedParameter(bool) {
	c := make(chan int, 1)
	c <- 42
}

func TestGotoLoop(t *testing.T) {
	goto loop
loop:
	for i := 42; ; {
		if i != 42 {
			t.Fail()
		}
		break
	}
}

func TestMaxUint64(t *testing.T) {
	if math.MaxUint64 != 18446744073709551615 {
		t.Fail()
	}
}

func TestCopyBuiltin(t *testing.T) {
	{
		s := []string{"a", "b", "c"}
		copy(s, s[1:])
		if s[0] != "b" || s[1] != "c" || s[2] != "c" {
			t.Fail()
		}
	}
	{
		s := []string{"a", "b", "c"}
		copy(s[1:], s)
		if s[0] != "a" || s[1] != "a" || s[2] != "b" {
			t.Fail()
		}
	}
}

func TestPointerOfStructConversion(t *testing.T) {
	type A struct {
		Value int
	}

	type B A

	a1 := &A{Value: 1}
	b1 := (*B)(a1)
	b1.Value = 2
	a2 := (*A)(b1)
	a2.Value = 3
	b2 := (*B)(a2)
	b2.Value = 4
	if a1 != a2 || b1 != b2 || a1.Value != 4 || a2.Value != 4 || b1.Value != 4 || b2.Value != 4 {
		t.Fail()
	}
}

func TestCompareStruct(t *testing.T) {
	type A struct {
		Value int
	}

	a := A{42}
	var b interface{} = a
	x := A{0}

	if a != b || a == x || b == x {
		t.Fail()
	}
}

func TestLoopClosure(t *testing.T) {
	type S struct{ fn func() int }
	var fns []*S
	for i := 0; i < 2; i++ {
		z := i
		fns = append(fns, &S{
			fn: func() int {
				return z
			},
		})
	}
	for i, f := range fns {
		if f.fn() != i {
			t.Fail()
		}
	}
}

func TestLoopClosureWithStruct(t *testing.T) {
	type T struct{ A int }
	ts := []T{{0}, {1}, {2}}
	fns := make([]func() T, 3)
	for i, t := range ts {
		t := t
		fns[i] = func() T {
			return t
		}
	}
	for i := range fns {
		if fns[i]().A != i {
			t.Fail()
		}
	}
}

func TestNilInterfaceError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil || !strings.Contains(err.(error).Error(), "nil pointer dereference") {
			t.Fail()
		}
	}()
	var err error
	_ = err.Error()
}

func TestIndexOutOfRangeError(t *testing.T) {
	defer func() {
		if err := recover(); err == nil || !strings.Contains(err.(error).Error(), "index out of range") {
			t.Fail()
		}
	}()
	x := []int{1, 2, 3}[10]
	_ = x
}

func TestNilAtLhs(t *testing.T) {
	type F func(string) string
	var f F
	if nil != f {
		t.Fail()
	}
}

func TestZeroResultByPanic(t *testing.T) {
	if zero() != 0 {
		t.Fail()
	}
}

func zero() int {
	defer func() {
		recover()
	}()
	panic("")
}

func TestNumGoroutine(t *testing.T) {
	n := runtime.NumGoroutine()
	c := make(chan bool)
	go func() {
		<-c
		<-c
		<-c
		<-c
	}()
	c <- true
	c <- true
	c <- true
	if got, want := runtime.NumGoroutine(), n+1; got != want {
		t.Errorf("runtime.NumGoroutine(): Got %d, want %d.", got, want)
	}
	c <- true
}

func TestMapAssign(t *testing.T) {
	x := 0
	m := map[string]string{}
	x, m["foo"] = 5, "bar"
	if x != 5 || m["foo"] != "bar" {
		t.Fail()
	}
}

func TestSwitchStatement(t *testing.T) {
	zero := 0
	var interfaceZero interface{} = zero
	switch {
	case interfaceZero:
		t.Fail()
	default:
		// ok
	}
}

func TestAddAssignOnPackageVar(t *testing.T) {
	otherpkg.Test = 0
	otherpkg.Test += 42
	if otherpkg.Test != 42 {
		t.Fail()
	}
}

func TestPointerOfPackageVar(t *testing.T) {
	otherpkg.Test = 42
	p := &otherpkg.Test
	if *p != 42 {
		t.Fail()
	}
}

func TestFuncInSelect(t *testing.T) {
	f := func(_ func()) chan int {
		return make(chan int, 1)
	}
	select {
	case <-f(func() {}):
	case _ = <-f(func() {}):
	case f(func() {}) <- 42:
	}
}

func TestEscapeAnalysisOnForLoopVariableScope(t *testing.T) {
	for i := 0; ; {
		p := &i
		time.Sleep(0)
		i = 42
		if *p != 42 {
			t.Fail()
		}
		break
	}
}

func TestGoStmtWithStructArg(t *testing.T) {
	type S struct {
		i int
	}

	f := func(s S, c chan int) {
		c <- s.i
		c <- s.i
	}

	c := make(chan int)
	s := S{42}
	go f(s, c)
	s.i = 0
	if <-c != 42 {
		t.Fail()
	}
	if <-c != 42 {
		t.Fail()
	}
}

type methodExprCallType int

func (i methodExprCallType) test() int {
	return int(i) + 2
}

func TestMethodExprCall(t *testing.T) {
	if methodExprCallType.test(40) != 42 {
		t.Fail()
	}
}

func TestCopyOnSend(t *testing.T) {
	type S struct{ i int }
	c := make(chan S, 2)
	go func() {
		var s S
		s.i = 42
		c <- s
		select {
		case c <- s:
		}
		s.i = 10
	}()
	if (<-c).i != 42 {
		t.Fail()
	}
	if (<-c).i != 42 {
		t.Fail()
	}
}

func TestEmptySelectCase(t *testing.T) {
	ch := make(chan int, 1)
	ch <- 42

	var v = 0
	select {
	case v = <-ch:
	}
	if v != 42 {
		t.Fail()
	}
}

var a int
var b int
var C int
var D int

var a1 = &a
var a2 = &a
var b1 = &b
var C1 = &C
var C2 = &C
var D1 = &D

func TestPkgVarPointers(t *testing.T) {
	if a1 != a2 || a1 == b1 || C1 != C2 || C1 == D1 {
		t.Fail()
	}
}

func TestStringMap(t *testing.T) {
	m := make(map[string]interface{})
	if m["__proto__"] != nil {
		t.Fail()
	}
	m["__proto__"] = 42
	if m["__proto__"] != 42 {
		t.Fail()
	}
}

type Int int

func (i Int) Value() int {
	return int(i)
}

func (i *Int) ValueByPtr() int {
	return int(*i)
}

func TestWrappedTypeMethod(t *testing.T) {
	i := Int(42)
	p := &i
	if p.Value() != 42 {
		t.Fail()
	}
}

type EmbeddedInt struct {
	Int
}

func TestEmbeddedMethod(t *testing.T) {
	e := EmbeddedInt{42}
	if e.ValueByPtr() != 42 {
		t.Fail()
	}
}

func TestBoolConvert(t *testing.T) {
	if !reflect.ValueOf(true).Convert(reflect.TypeOf(true)).Bool() {
		t.Fail()
	}
}

func TestGoexit(t *testing.T) {
	go func() {
		runtime.Goexit()
	}()
}

func TestVendoring(t *testing.T) {
	if vendored.Answer != 42 {
		t.Fail()
	}
}

func TestShift(t *testing.T) {
	if x := uint(32); uint32(1)<<x != 0 {
		t.Fail()
	}
	if x := uint64(0); uint32(1)<<x != 1 {
		t.Fail()
	}
	if x := uint(4294967295); x>>32 != 0 {
		t.Fail()
	}
	if x := uint(4294967295); x>>35 != 0 {
		t.Fail()
	}
}

func TestTrivialSwitch(t *testing.T) {
	for {
		switch {
		default:
			break
		}
		return
	}
	t.Fail()
}

func TestTupleFnReturnImplicitCast(t *testing.T) {
	var ycalled int = 0
	x := func(fn func() (int, error)) (interface{}, error) {
		return fn()
	}
	y, _ := x(func() (int, error) {
		ycalled++
		return 14, nil
	})
	if y != 14 || ycalled != 1 {
		t.Fail()
	}
}

var tuple2called = 0

func tuple1() (interface{}, error) {
	return tuple2()
}
func tuple2() (int, error) {
	tuple2called++
	return 14, nil
}
func TestTupleReturnImplicitCast(t *testing.T) {
	x, _ := tuple1()
	if x != 14 || tuple2called != 1 {
		t.Fail()
	}
}

func TestDeferNamedTupleReturnImplicitCast(t *testing.T) {
	var ycalled int = 0
	var zcalled int = 0
	z := func() {
		zcalled++
	}
	x := func(fn func() (int, error)) (i interface{}, e error) {
		defer z()
		i, e = fn()
		return
	}
	y, _ := x(func() (int, error) {
		ycalled++
		return 14, nil
	})
	if y != 14 || ycalled != 1 || zcalled != 1 {
		t.Fail()
	}
}

func TestSliceOfString(t *testing.T) {
	defer func() {
		if err := recover(); err == nil || !strings.Contains(err.(error).Error(), "slice bounds out of range") {
			t.Fail()
		}
	}()

	str := "foo"
	print(str[0:10])
}
