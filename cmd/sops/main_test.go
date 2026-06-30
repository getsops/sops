package main

import (
	"flag"
	"os"
	"testing"

	"github.com/getsops/sops/v3/age"
	"github.com/urfave/cli"
)

func TestApplyAgeKeyFileFlag(t *testing.T) {
	t.Run("sets the environment when the flag is provided", func(t *testing.T) {
		// t.Setenv registers restoration of the original value; clearing it
		// afterwards gives the subtest a known, unset starting point.
		t.Setenv(age.SopsAgeKeyFileEnv, "")
		if err := os.Unsetenv(age.SopsAgeKeyFileEnv); err != nil {
			t.Fatalf("failed to unset %s: %v", age.SopsAgeKeyFileEnv, err)
		}

		const want = "/path/to/keys.txt"
		if err := applyAgeKeyFileFlag(want); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := os.Getenv(age.SopsAgeKeyFileEnv); got != want {
			t.Errorf("%s = %q, want %q", age.SopsAgeKeyFileEnv, got, want)
		}
	})

	t.Run("overrides a pre-existing environment value", func(t *testing.T) {
		t.Setenv(age.SopsAgeKeyFileEnv, "/from/env")

		const want = "/from/flag"
		if err := applyAgeKeyFileFlag(want); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := os.Getenv(age.SopsAgeKeyFileEnv); got != want {
			t.Errorf("%s = %q, want %q", age.SopsAgeKeyFileEnv, got, want)
		}
	})

	t.Run("leaves the environment untouched when the flag is empty", func(t *testing.T) {
		const want = "/from/env"
		t.Setenv(age.SopsAgeKeyFileEnv, want)

		if err := applyAgeKeyFileFlag(""); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got := os.Getenv(age.SopsAgeKeyFileEnv); got != want {
			t.Errorf("%s = %q, want %q", age.SopsAgeKeyFileEnv, got, want)
		}
	})
}

// TestAgeKeyFileFlagBefore exercises the actual app.Before hook through a
// cli.Context carrying the --age-key-file flag. Because the hook is wired to
// app.Before (not app.Action), it runs for every command path, including the
// decrypt/exec-env/exec-file subcommands that have their own actions.
func TestAgeKeyFileFlagBefore(t *testing.T) {
	t.Setenv(age.SopsAgeKeyFileEnv, "")
	if err := os.Unsetenv(age.SopsAgeKeyFileEnv); err != nil {
		t.Fatalf("failed to unset %s: %v", age.SopsAgeKeyFileEnv, err)
	}

	const want = "/path/to/keys.txt"
	set := flag.NewFlagSet("test", flag.ContinueOnError)
	set.String("age-key-file", "", "")
	if err := set.Set("age-key-file", want); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	if err := ageKeyFileFlagBefore(cli.NewContext(nil, set, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv(age.SopsAgeKeyFileEnv); got != want {
		t.Errorf("%s = %q, want %q", age.SopsAgeKeyFileEnv, got, want)
	}
}
