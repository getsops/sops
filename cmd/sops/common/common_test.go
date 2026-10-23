package common

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getsops/sops/v3"
	. "github.com/getsops/sops/v3/cmd/sops/formats"
	"github.com/getsops/sops/v3/config"
	"github.com/stretchr/testify/assert"
)

func TestLoadEncryptedFileNotSopsFile(t *testing.T) {
	store := StoreForFormat(Json, config.NewStoresConfig())

	t.Run("unparsable input mentions sops", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plain.json")
		assert.NoError(t, os.WriteFile(path, []byte("# this is just a text file\n"), 0o600))
		_, err := LoadEncryptedFile(store, path)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Could not unmarshal input data")
		assert.Contains(t, err.Error(), "sops-encrypted file")
	})

	t.Run("missing metadata is still reported as MetadataNotFound", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "plain.json")
		assert.NoError(t, os.WriteFile(path, []byte(`{"a": "b"}`), 0o600))
		_, err := LoadEncryptedFile(store, path)
		assert.Equal(t, sops.MetadataNotFound, err)
		assert.Contains(t, err.Error(), "does not appear to be encrypted with sops")
	})
}
