package file

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/sdk/helper/salt"
	"github.com/hashicorp/vault/sdk/logical"
)

func Factory(ctx context.Context, conf *audit.BackendConfig) (audit.Backend, error) {
	if conf.SaltConfig == nil {
		return nil, fmt.Errorf("nil salt config")
	}
	if conf.SaltView == nil {
		return nil, fmt.Errorf("nil salt view")
	}

	path, ok := conf.Config["file_path"]
	if !ok {
		path, ok = conf.Config["path"]
		if !ok {
			return nil, fmt.Errorf("file_path is required")
		}
	}

	// normalize path if configured for stdout
	if strings.EqualFold(path, "stdout") {
		path = "stdout"
	}
	if strings.EqualFold(path, "discard") {
		path = "discard"
	}

	format, ok := conf.Config["format"]
	if !ok {
		format = "json"
	}
	switch format {
	case "json", "jsonx":
	default:
		return nil, fmt.Errorf("unknown format type %q", format)
	}

	// Check if hashing of accessor is disabled
	hmacAccessor := true
	if hmacAccessorRaw, ok := conf.Config["hmac_accessor"]; ok {
		value, err := strconv.ParseBool(hmacAccessorRaw)
		if err != nil {
			return nil, err
		}
		hmacAccessor = value
	}

	// Check if raw logging is enabled
	logRaw := false
	if raw, ok := conf.Config["log_raw"]; ok {
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return nil, err
		}
		logRaw = b
	}

	// Check if mode is provided
	mode := os.FileMode(0600)
	if modeRaw, ok := conf.Config["mode"]; ok {
		m, err := strconv.ParseUint(modeRaw, 8, 32)
		if err != nil {
			return nil, err
		}
		if m != 0 {
			mode = os.FileMode(m)
		}
	}

	b := &Backend{
		path:       path,
		mode:       mode,
		saltConfig: conf.SaltConfig,
		saltView:   conf.SaltView,
		salt:       new(atomic.Value),
		formatConfig: audit.FormatterConfig{
			Raw:          logRaw,
			HMACAccessor: hmacAccessor,
		},
	}

	// Ensure we are working with the right type by explicitly storing a nil of
	// the right type
	b.salt.Store((*salt.Salt)(nil))

	switch format {
	case "json":
		b.formatter.AuditFormatWriter = &audit.JSONFormatWriter{
			Prefix:   conf.Config["prefix"],
			SaltFunc: b.Salt,
		}
	case "jsonx":
		b.formatter.AuditFormatWriter = &audit.JSONxFormatWriter{
			Prefix:   conf.Config["prefix"],
			SaltFunc: b.Salt,
		}
	}

	switch path {
	case "stdout", "discard":
		// no need to test opening file if outputting to stdout or discarding
	default:
		// Ensure that the file can be successfully opened for writing;
		// otherwise it will be too late to catch later without problems
		// (ref: https://github.com/hashicorp/vault/issues/550)
		if err := b.open(); err != nil {
			return nil, errwrap.Wrapf(fmt.Sprintf("sanity check failed; unable to open %q for writing: {{err}}", path), err)
		}
	}

	return b, nil
}

// Backend is the audit backend for the file-based audit store.
//
// NOTE: This audit backend is currently very simple: it appends to a file.
// It doesn't do anything more at the moment to assist with rotation
// or reset the write cursor, this should be done in the future.
type Backend struct {
	path string

	formatter    audit.AuditFormatter
	formatConfig audit.FormatterConfig

	fileLock sync.RWMutex
	f        *os.File
	mode     os.FileMode

	saltMutex  sync.RWMutex
	salt       *atomic.Value
	saltConfig *salt.Config
	saltView   logical.Storage
}

var _ audit.Backend = (*Backend)(nil)

func (b *Backend) Salt(ctx context.Context) (*salt.Salt, error) {
	s := b.salt.Load().(*salt.Salt)
	if s != nil {
		return s, nil
	}

	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()

	s = b.salt.Load().(*salt.Salt)
	if s != nil {
		return s, nil
	}

	newSalt, err := salt.NewSalt(ctx, b.saltView, b.saltConfig)
	if err != nil {
		b.salt.Store((*salt.Salt)(nil))
		return nil, err
	}

	b.salt.Store(newSalt)
	return newSalt, nil
}

func (b *Backend) GetHash(ctx context.Context, data string) (string, error) {
	salt, err := b.Salt(ctx)
	if err != nil {
		return "", err
	}

	return audit.HashString(salt, data), nil
}

func (b *Backend) LogRequest(ctx context.Context, in *logical.LogInput) error {
	var writer io.Writer
	switch b.path {
	case "stdout":
		writer = os.Stdout
	case "discard":
		return nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, 2000))
	err := b.formatter.FormatRequest(ctx, buf, b.formatConfig, in)
	if err != nil {
		return err
	}

	return b.log(ctx, buf, writer)
}

func (b *Backend) log(ctx context.Context, buf *bytes.Buffer, writer io.Writer) error {
	reader := bytes.NewReader(buf.Bytes())

	b.fileLock.Lock()

	if writer == nil {
		if err := b.open(); err != nil {
			b.fileLock.Unlock()
			return err
		}
		writer = b.f
	}

	if _, err := reader.WriteTo(writer); err == nil {
		b.fileLock.Unlock()
		return nil
	} else if b.path == "stdout" {
		b.fileLock.Unlock()
		return err
	}

	// If writing to stdout there's no real reason to think anything would have
	// changed so return above. Otherwise, opportunistically try to re-open the
	// FD, once per call.
	b.f.Close()
	b.f = nil

	if err := b.open(); err != nil {
		b.fileLock.Unlock()
		return err
	}

	reader.Seek(0, io.SeekStart)
	_, err := reader.WriteTo(writer)
	b.fileLock.Unlock()
	return err
}

func (b *Backend) LogResponse(ctx context.Context, in *logical.LogInput) error {
	var writer io.Writer
	switch b.path {
	case "stdout":
		writer = os.Stdout
	case "discard":
		return nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, 6000))
	err := b.formatter.FormatResponse(ctx, buf, b.formatConfig, in)
	if err != nil {
		return err
	}

	return b.log(ctx, buf, writer)
}

// The file lock must be held before calling this
func (b *Backend) open() error {
	if b.f != nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(b.path), b.mode); err != nil {
		return err
	}

	var err error
	b.f, err = os.OpenFile(b.path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, b.mode)
	if err != nil {
		return err
	}

	// Change the file mode in case the log file already existed. We special
	// case /dev/null since we can't chmod it and bypass if the mode is zero
	switch b.path {
	case "/dev/null":
	default:
		if b.mode != 0 {
			err = os.Chmod(b.path, b.mode)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (b *Backend) Reload(_ context.Context) error {
	switch b.path {
	case "stdout", "discard":
		return nil
	}

	b.fileLock.Lock()
	defer b.fileLock.Unlock()

	if b.f == nil {
		return b.open()
	}

	err := b.f.Close()
	// Set to nil here so that even if we error out, on the next access open()
	// will be tried
	b.f = nil
	if err != nil {
		return err
	}

	return b.open()
}

func (b *Backend) Invalidate(_ context.Context) {
	b.saltMutex.Lock()
	defer b.saltMutex.Unlock()
	b.salt.Store((*salt.Salt)(nil))
}
