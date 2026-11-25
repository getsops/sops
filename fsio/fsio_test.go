package fsio

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestReadRegularFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")
	content := []byte("regular file content")

	err := os.WriteFile(filePath, content, 0o600)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	b1, err := Read(filePath)
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}
	if !bytes.Equal(b1, content) {
		t.Errorf("expected %q, got %q", content, b1)
	}

	if _, ok := fileStreamCache.Load(filePath); ok {
		t.Error("expected regular file NOT to be cached, but it was")
	}
}
