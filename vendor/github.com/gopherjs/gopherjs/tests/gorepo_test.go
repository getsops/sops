package tests_test

import (
	"os"
	"os/exec"
	"runtime"
	"testing"
)

// Go repository basic compiler tests, and regression tests for fixed compiler bugs.
func TestGoRepositoryCompilerTests(t *testing.T) {
	if runtime.GOARCH == "js" {
		t.Skip("test meant to be run using normal Go compiler (needs os/exec)")
	}

	args := []string{"go", "run", "run.go", "-summary"}
	if testing.Verbose() {
		args = append(args, "-v")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout
	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

// Test that GopherJS can be vendored into a project, and then used to build Go programs.
// See issue https://github.com/gopherjs/gopherjs/issues/415.
func TestGopherJSCanBeVendored(t *testing.T) {
	if runtime.GOARCH == "js" {
		t.Skip("test meant to be run using normal Go compiler (needs os/exec)")
	}

	cmd := exec.Command("sh", "gopherjsvendored_test.sh")
	cmd.Stderr = os.Stdout
	got, err := cmd.Output()
	if err != nil {
		t.Fatal(err)
	}
	if want := "hello using js pkg\n"; string(got) != want {
		t.Errorf("unexpected stdout from gopherjsvendored_test.sh:\ngot:\n%s\nwant:\n%s", got, want)
	}
}
