package age

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetOutputFromCmd(t *testing.T) {
	const diagnostic = "cannot connect to 1Password app, make sure it is running"
	const secret = "disposable-key-output"
	tests := []struct {
		name       string
		stderr     string
		exitCode   int
		wantSuffix string
	}{
		{"failure with stderr", diagnostic, 7, ": " + diagnostic},
		{"failure without stderr", "", 7, "exit status 7"},
		{"failure with whitespace stderr", " \n\t ", 7, "exit status 7"},
		{"failure with padded stderr", "\n" + diagnostic + "\n", 7, ": " + diagnostic},
		{"success without stderr", "", 0, ""},
		{"success with stderr", diagnostic, 0, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			script := filepath.Join(t.TempDir(), "key-command.sh")
			contents := fmt.Sprintf(
				"#!/bin/bash\nset -euo pipefail\n"+
					"printf '%%s' \"$TEST_KEY\"\nprintf '%%s' \"$TEST_ERROR\" >&2\nexit %d\n",
				tt.exitCode,
			)
			require.NoError(t, os.WriteFile(script, []byte(contents), 0o600))
			out, err := getOutputFromCmd(fmt.Sprintf("bash %q", script), []string{
				"TEST_KEY=" + secret,
				"TEST_ERROR=" + tt.stderr,
			})
			if tt.exitCode == 0 {
				require.NoError(t, err)
				assert.Equal(t, secret, string(out))
				return
			}
			require.Error(t, err)
			assert.Nil(t, out)
			assert.NotContains(t, err.Error(), secret)
			assert.Contains(t, err.Error(), "failed to execute command")
			assert.True(t, strings.HasSuffix(err.Error(), tt.wantSuffix), err.Error())
			var exitError *exec.ExitError
			require.True(t, errors.As(err, &exitError))
			assert.Equal(t, tt.exitCode, exitError.ExitCode())
		})
	}
}
