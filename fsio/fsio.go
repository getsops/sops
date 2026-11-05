package fsio

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type cacheEntry struct {
	mu   sync.RWMutex
	data []byte
}

var fileStreamCache sync.Map

// ClearCache wipes the cached stream secrets from memory by overwriting
// the byte slices with zeros before deleting them from the map.
//
// If you are using SOPS as a library, you should call ClearCache after
// completing decryption/encryption operations to ensure no sensitive key
// data remains in memory.
func ClearCache() {
	fileStreamCache.Range(func(key, value any) bool {
		if entry, ok := value.(*cacheEntry); ok {
			entry.mu.Lock()
			for i := range entry.data {
				entry.data[i] = 0
			}
			entry.mu.Unlock()
		}
		fileStreamCache.Delete(key)
		return true
	})
}

// Read reads a file from the given path. If it is a stream (e.g., /dev/fd/* or /proc/*)
// it caches the content in memory to avoid issues with multiple reads from the same stream.
func Read(path string) ([]byte, error) {
	if absPath, err := filepath.Abs(path); err == nil {
		path = absPath
	}
	fileInfo, err := os.Stat(path)
	isStream := err == nil &&
		(fileInfo.Mode()&os.ModeNamedPipe != 0 || fileInfo.Mode()&os.ModeCharDevice != 0 || fileInfo.Mode()&os.ModeSocket != 0)

	if isStream {
		if value, ok := fileStreamCache.Load(path); ok {
			if entry, ok := value.(*cacheEntry); ok {
				entry.mu.RLock()
				defer entry.mu.RUnlock()
				b := make([]byte, len(entry.data))
				copy(b, entry.data)
				return b, nil
			}
		}
	}

	b, err := os.ReadFile(path)
	if err == nil && isStream {
		cachedBytes := make([]byte, len(b))
		copy(cachedBytes, b)
		fileStreamCache.Store(path, &cacheEntry{data: cachedBytes})
	}
	return b, err
}

// Open opens a file from the given path. If it is a stream, it loads the content
// into the cache and returns a reader over the cached bytes.
func Open(path string) (io.ReadCloser, error) {
	if absPath, err := filepath.Abs(path); err == nil {
		path = absPath
	}
	fileInfo, err := os.Stat(path)
	isStream := err == nil &&
		(fileInfo.Mode()&os.ModeNamedPipe != 0 || fileInfo.Mode()&os.ModeCharDevice != 0 || fileInfo.Mode()&os.ModeSocket != 0)

	if isStream {
		b, err := Read(path)
		if err != nil {
			return nil, err
		}
		return io.NopCloser(bytes.NewReader(b)), nil
	}

	return os.Open(path)
}
