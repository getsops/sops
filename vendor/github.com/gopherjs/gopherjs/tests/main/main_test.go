package main

import "testing"

func TestNotRunMain(t *testing.T) {
	if mainDidRun {
		t.Error("main function did run")
	}
}
