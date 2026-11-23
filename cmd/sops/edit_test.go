package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
)

func TestEditKeepsChanges(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the test editor is a shell script")
	}
	for _, showMasterKeys := range []bool{false, true} {
		name := "plain"
		if showMasterKeys {
			name = "show-master-keys"
		}
		t.Run(name, func(t *testing.T) {
			isolateAgeEnvironment(t)
			key := newTestAgeKey(t)
			path := encryptTestFile(t, "secret: old\n", []sops.KeyGroup{{key}})

			editor := filepath.Join(t.TempDir(), "editor.sh")
			script := "#!/bin/sh\nsed 's/secret: old/secret: new/' \"$1\" > \"$1.tmp\" && mv \"$1.tmp\" \"$1\"\n"
			require.NoError(t, os.WriteFile(editor, []byte(script), 0700))
			t.Setenv("SOPS_EDITOR", editor)

			store := testYAMLStore()
			output, err := edit(editOpts{
				Cipher:         aes.NewCipher(),
				InputStore:     store,
				OutputStore:    store,
				InputPath:      path,
				KeyServices:    testKeyServices(),
				ShowMasterKeys: showMasterKeys,
			})
			require.NoError(t, err)
			require.NoError(t, os.WriteFile(path, output, 0600))

			assert.Equal(t, "secret: new\n", decryptTestFile(t, path))
		})
	}
}
