// +build !js

package tests_test

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"testing"
)

// Test for internalization/externalization of time.Time/Date when time package is imported
// but time.Time is unused, causing it to be DCEed (or time package not imported at all).
//
// See https://github.com/gopherjs/gopherjs/issues/279.
func TestTimeInternalizationExternalization(t *testing.T) {
	got, err := exec.Command("gopherjs", "run", filepath.Join("testdata", "time_inexternalization.go")).Output()
	if err != nil {
		t.Fatalf("%v:\n%s", err, got)
	}

	want, err := ioutil.ReadFile(filepath.Join("testdata", "time_inexternalization.out"))
	if err != nil {
		t.Fatalf("error reading .out file: %v", err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("got != want:\ngot:\n%s\nwant:\n%s", got, want)
	}
}
