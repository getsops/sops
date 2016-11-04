package tests

import (
	"testing"
)

type S struct {
	x int
}

func (a S) test(b S) {
	a.x = 0
	b.x = 0
}

type A [1]int

func (a A) test(b A) {
	a[0] = 0
	b[0] = 0
}

func TestCopyOnCall(t *testing.T) {
	{
		a := S{1}
		b := S{2}

		a.test(b)
		func() {
			defer a.test(b)
		}()

		if a.x != 1 {
			t.Error("a.x != 1")
		}
		if b.x != 2 {
			t.Error("b.x != 2")
		}
	}
	{
		a := A{1}
		b := A{2}

		a.test(b)
		func() {
			defer a.test(b)
		}()

		if a[0] != 1 {
			t.Error("a[0] != 1")
		}
		if b[0] != 2 {
			t.Error("b[0] != 2")
		}
	}
}

func TestSwap(t *testing.T) {
	{
		a := S{1}
		b := S{2}
		a, b = b, a
		if a.x != 2 || b.x != 1 {
			t.Fail()
		}
	}
	{
		a := A{1}
		b := A{2}
		a, b = b, a
		if a[0] != 2 || b[0] != 1 {
			t.Fail()
		}
	}
}

func TestComposite(t *testing.T) {
	{
		a := S{1}
		s := []S{a}
		s[0].x = 0
		if a.x != 1 {
			t.Fail()
		}
	}
	{
		a := A{1}
		s := []A{a}
		s[0][0] = 0
		if a[0] != 1 {
			t.Fail()
		}
	}
}

func TestAppend(t *testing.T) {
	{
		s := append(make([]S, 3), S{}) // cap(s) == 6
		s = s[:6]
		if s[5].x != 0 {
			t.Fail()
		}
	}

	{
		a := S{1}
		b := []S{{2}}
		s := append([]S{}, b...)
		s[0].x = 0
		if a.x != 1 || b[0].x != 2 {
			t.Fail()
		}
	}
}

type I interface {
	M() int
}

type T S

func (t T) M() int {
	return t.x
}

func TestExplicitConversion(t *testing.T) {
	var coolGuy = S{x: 42}
	var i I
	i = T(coolGuy)
	if i.M() != 42 {
		t.Fail()
	}
}
