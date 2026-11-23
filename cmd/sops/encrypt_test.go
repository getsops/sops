package main

import (
	"testing"

	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/stores/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateFileForEncryption(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input string
		code  int
	}{
		{name: "empty", code: codes.NeedAtLeastOneDocument},
		{name: "plain single document", input: "secret: value\n"},
		{name: "plain multiple documents", input: "secret: first\n---\nsecret: second\n"},
		{name: "metadata in first document", input: "sops: {}\n---\nsecret: value\n", code: codes.FileAlreadyEncrypted},
		{name: "metadata in second document", input: "secret: value\n---\nsops: {}\n", code: codes.FileAlreadyEncrypted},
		{name: "metadata in middle document", input: "secret: first\n---\nsops: {}\n---\nsecret: last\n", code: codes.FileAlreadyEncrypted},
		{name: "metadata in last document", input: "secret: first\n---\nsecret: second\n---\nsops: {}\n", code: codes.FileAlreadyEncrypted},
		{name: "nested sops key", input: "secret: first\n---\napp:\n  sops: value\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			store := &yaml.Store{}
			branches, err := store.LoadPlainFile([]byte(tt.input))
			require.NoError(t, err)
			userErr, code := validateFileForEncryption(store, branches)
			assert.Equal(t, tt.code, code)
			switch tt.code {
			case codes.FileAlreadyEncrypted:
				assert.IsType(t, &fileAlreadyEncryptedError{}, userErr)
			case codes.NeedAtLeastOneDocument:
				assert.IsType(t, &needAtLeastOneDocument{}, userErr)
			default:
				assert.Nil(t, userErr)
			}
		})
	}
}
