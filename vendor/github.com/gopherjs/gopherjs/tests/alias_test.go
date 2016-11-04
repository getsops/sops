package tests

import (
	"testing"
)

type foo struct {
	a int
}

func Test1(t *testing.T) {
	calls := 0
	bar := func() *foo {
		calls++
		return &foo{42}
	}
	q := &bar().a
	if calls != 1 {
		t.Error("Should've been a call")
	}
	*q = 40
	if calls != 1 {
		t.Error("Wrong number of calls: ", calls, ", should be 1")
	}
	if *q != 40 {
		t.Error("*q != 40")
	}
}

func Test2(t *testing.T) {
	f := foo{}
	p := &f.a
	f = foo{}
	f.a = 4
	if *p != 4 {
		t.Error("*p != 4")
	}
}

func Test3(t *testing.T) {
	f := foo{}
	p := &f
	f = foo{4}
	if p.a != 4 {
		t.Error("p.a != 4")
	}
}

func Test4(t *testing.T) {
	f := struct {
		a struct {
			b int
		}
	}{}
	p := &f.a
	q := &p.b
	r := &(*p).b
	*r = 4
	p = nil
	if *r != 4 {
		t.Error("*r != 4")
	}
	if *q != 4 {
		t.Error("*q != 4")
	}
}

func Test5(t *testing.T) {
	f := struct {
		a [3]int
	}{[3]int{6, 6, 6}}
	s := f.a[:]
	f.a = [3]int{4, 4, 4}
	if s[1] != 4 {
		t.Error("s[1] != 4")
	}
}
