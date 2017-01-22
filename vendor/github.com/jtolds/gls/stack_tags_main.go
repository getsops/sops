// +build !js

package gls

// This file is used for standard Go builds, which have the expected runtime
// support

import (
	"runtime"
)

func getStack(offset, amount int) []uintptr {
	stack := make([]uintptr, amount)
	return stack[:runtime.Callers(offset, stack)]
}

func findPtr() uintptr {
	pc, _, _, ok := runtime.Caller(3)
	if !ok {
		panic("failed to find function pointer")
	}
	return pc
}
