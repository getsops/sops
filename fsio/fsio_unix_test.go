//go:build !windows

package fsio

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"testing"
)

func TestReadStreamFile(t *testing.T) {
	tempDir := t.TempDir()
	fifoPath := filepath.Join(tempDir, "test_fifo")

	err := syscall.Mkfifo(fifoPath, 0o600)
	if err != nil {
		t.Fatalf("failed to create named pipe: %v", err)
	}

	content := []byte("super secret key stream")

	go func() {
		f, err := os.OpenFile(fifoPath, os.O_WRONLY, 0)
		if err != nil {
			return
		}
		defer f.Close()
		_, _ = f.Write(content)
	}()

	b1, err := Read(fifoPath)
	if err != nil {
		t.Fatalf("first Read failed: %v", err)
	}
	if !bytes.Equal(b1, content) {
		t.Errorf("expected %q, got %q", content, b1)
	}

	if _, ok := fileStreamCache.Load(fifoPath); !ok {
		t.Error("expected stream file to be cached, but it was not")
	}

	b2, err := Read(fifoPath)
	if err != nil {
		t.Fatalf("second Read failed: %v", err)
	}
	if !bytes.Equal(b2, content) {
		t.Errorf("expected cached %q, got %q", content, b2)
	}

	r, err := Open(fifoPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer r.Close()
	b3, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll failed: %v", err)
	}
	if !bytes.Equal(b3, content) {
		t.Errorf("expected read cached %q, got %q", content, b3)
	}

	// Get the cache entry before clearing to verify it gets zeroed
	entryVal, ok := fileStreamCache.Load(fifoPath)
	if !ok {
		t.Fatal("expected entry to be in cache")
	}
	entry := entryVal.(*cacheEntry)

	// Clear cache and check that bytes are zeroed
	ClearCache()

	// The cached slice should be zeroed
	entry.mu.RLock()
	for i, b := range entry.data {
		if b != 0 {
			t.Errorf("expected cached byte at index %d to be zeroed, got %d", i, b)
		}
	}
	entry.mu.RUnlock()

	if _, ok := fileStreamCache.Load(fifoPath); ok {
		t.Error("expected stream file cache to be cleared, but it was still present")
	}
}

func TestConcurrentReadAndClear(t *testing.T) {
	tempDir := t.TempDir()
	fifoPath := filepath.Join(tempDir, "concurrent_fifo")

	err := syscall.Mkfifo(fifoPath, 0o600)
	if err != nil {
		t.Fatalf("failed to create named pipe: %v", err)
	}

	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}
			f, err := os.OpenFile(fifoPath, os.O_WRONLY, 0)
			if err != nil {
				return
			}
			_, _ = f.Write([]byte("secret_value"))
			f.Close()
		}
	}()

	_, err = Read(fifoPath)
	if err != nil {
		t.Fatalf("initial read failed: %v", err)
	}

	done := make(chan bool)
	go func() {
		for range 100 {
			ClearCache()
			_, _ = Read(fifoPath)
		}
		done <- true
	}()

	go func() {
		for range 100 {
			b, err := Read(fifoPath)
			if err == nil {
				_ = len(b)
				if len(b) > 0 {
					_ = b[0]
				}
			}
		}
		done <- true
	}()

	<-done
	<-done

	close(stop)
	dummy, err := os.OpenFile(fifoPath, os.O_RDONLY|syscall.O_NONBLOCK, 0)
	if err == nil {
		dummy.Close()
	}
}

func TestAnonymousPipe(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer r.Close()

	path := filepath.Join("/proc/self/fd", fmt.Sprintf("%d", r.Fd()))

	content := []byte("anonymous pipe secret key")

	go func() {
		defer w.Close()
		_, _ = w.Write(content)
	}()

	rCloser1, err := Open(path)
	if err != nil {
		t.Fatalf("first Open failed: %v", err)
	}
	b1, err := io.ReadAll(rCloser1)
	rCloser1.Close()
	if err != nil {
		t.Fatalf("first ReadAll failed: %v", err)
	}
	if !bytes.Equal(b1, content) {
		t.Errorf("expected %q, got %q", content, b1)
	}

	if _, ok := fileStreamCache.Load(path); !ok {
		t.Error("expected anonymous pipe to be cached, but it was not")
	}

	rCloser2, err := Open(path)
	if err != nil {
		t.Fatalf("second Open failed: %v", err)
	}
	b2, err := io.ReadAll(rCloser2)
	rCloser2.Close()
	if err != nil {
		t.Fatalf("second ReadAll failed: %v", err)
	}
	if !bytes.Equal(b2, content) {
		t.Errorf("expected cached %q, got %q", content, b2)
	}

	ClearCache()
}

func TestConcurrentReads(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	defer r.Close()

	path := filepath.Join("/proc/self/fd", fmt.Sprintf("%d", r.Fd()))
	content := []byte("concurrent shared secret key")

	go func() {
		defer w.Close()
		_, _ = w.Write(content)
	}()

	_, err = Read(path)
	if err != nil {
		t.Fatalf("initial read failed: %v", err)
	}

	const numGoroutines = 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	for range numGoroutines {
		go func() {
			defer wg.Done()
			b, err := Read(path)
			if err != nil {
				t.Errorf("Read failed: %v", err)
				return
			}
			if !bytes.Equal(b, content) {
				t.Errorf("expected cached %q, got %q", content, b)
			}
		}()
	}

	wg.Wait()
	ClearCache()
}

func TestReadCanonicalization(t *testing.T) {
	tempDir := t.TempDir()
	fifoPath := filepath.Join(tempDir, "canonical_fifo")

	err := syscall.Mkfifo(fifoPath, 0o600)
	if err != nil {
		t.Fatalf("failed to create named pipe: %v", err)
	}

	content := []byte("canonicalized stream content")

	go func() {
		f, err := os.OpenFile(fifoPath, os.O_WRONLY, 0)
		if err != nil {
			return
		}
		defer f.Close()
		_, _ = f.Write(content)
	}()

	absPath := filepath.Join(tempDir, "canonical_fifo")
	altPath := filepath.Join(tempDir, "..", filepath.Base(tempDir), "canonical_fifo")

	b1, err := Read(altPath)
	if err != nil {
		t.Fatalf("first Read with altPath failed: %v", err)
	}
	if !bytes.Equal(b1, content) {
		t.Errorf("expected %q, got %q", content, b1)
	}

	b2, err := Read(absPath)
	if err != nil {
		t.Fatalf("second Read with absPath failed: %v", err)
	}
	if !bytes.Equal(b2, content) {
		t.Errorf("expected cached %q, got %q", content, b2)
	}

	if _, ok := fileStreamCache.Load(absPath); !ok {
		t.Errorf("expected stream to be cached under absolute path %q, but it was not", absPath)
	}

	ClearCache()
}
