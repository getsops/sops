// +build js

package reflect_test

import (
	"reflect"
	"testing"
)

func TestAlignment(t *testing.T) {
	t.Skip()
}

func TestSliceOverflow(t *testing.T) {
	t.Skip()
}

func TestFuncLayout(t *testing.T) {
	t.Skip()
}

func TestArrayOfDirectIface(t *testing.T) {
	t.Skip()
}

func TestTypelinksSorted(t *testing.T) {
	t.Skip()
}

func TestGCBits(t *testing.T) {
	t.Skip()
}

func TestChanAlloc(t *testing.T) {
	t.Skip()
}

func TestNameBytesAreAligned(t *testing.T) {
	t.Skip()
}

func TestOffsetLock(t *testing.T) {
	t.Skip()
}

func TestSelectOnInvalid(t *testing.T) {
	reflect.Select([]reflect.SelectCase{
		{
			Dir:  reflect.SelectRecv,
			Chan: reflect.Value{},
		}, {
			Dir:  reflect.SelectSend,
			Chan: reflect.Value{},
			Send: reflect.ValueOf(1),
		}, {
			Dir: reflect.SelectDefault,
		},
	})
}

func TestStructOf(t *testing.T) {
	t.Skip("StructOf")
}

func TestStructOfExportRules(t *testing.T) {
	t.Skip("StructOf")
}

func TestStructOfGC(t *testing.T) {
	t.Skip("StructOf")
}

func TestStructOfAlg(t *testing.T) {
	t.Skip("StructOf")
}

func TestStructOfGenericAlg(t *testing.T) {
	t.Skip("StructOf")
}

func TestStructOfDirectIface(t *testing.T) {
	t.Skip("StructOf")
}

func TestStructOfWithInterface(t *testing.T) {
	t.Skip("StructOf")
}
