package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Local stores files on disk under RootDir and builds URLs with PublicBaseURL.
type Local struct {
	RootDir       string
	PublicBaseURL string
}

// NewLocal creates a local disk provider. RootDir is made absolute.
func NewLocal(rootDir, publicBaseURL string) (*Local, error) {
	abs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf("storage local root: %w", err)
	}
	if err := os.MkdirAll(abs, 0o750); err != nil {
		return nil, fmt.Errorf("storage local mkdir: %w", err)
	}
	base := strings.TrimSuffix(strings.TrimSpace(publicBaseURL), "/")
	if base == "" {
		return nil, fmt.Errorf("storage public base URL is empty")
	}
	return &Local{RootDir: abs, PublicBaseURL: base}, nil
}

// Put writes key under RootDir.
func (l *Local) Put(ctx context.Context, key string, r io.Reader, size int64, _ string) error {
	if err := validateStorageKey(key); err != nil {
		return err
	}
	full := filepath.Join(l.RootDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(full), 0o750); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	tmp := full + ".tmp"
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o640)
	if err != nil {
		return fmt.Errorf("open temp: %w", err)
	}
	defer os.Remove(tmp)
	n, err := io.Copy(f, r)
	if closeErr := f.Close(); closeErr != nil && err == nil {
		err = closeErr
	}
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	if size > 0 && n != size {
		return fmt.Errorf("size mismatch: wrote %d of %d", n, size)
	}
	if err := os.Rename(tmp, full); err != nil {
		return fmt.Errorf("rename: %w", err)
	}
	_ = ctx // reserved for cancellation wiring
	return nil
}

// Delete removes the file for key.
func (l *Local) Delete(ctx context.Context, key string) error {
	_ = ctx
	if err := validateStorageKey(key); err != nil {
		return err
	}
	full := filepath.Join(l.RootDir, filepath.FromSlash(key))
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// PublicURL joins the configured base with a slash-encoded key path.
func (l *Local) PublicURL(key string) string {
	return l.PublicBaseURL + "/" + pathEscapeSegments(key)
}

func pathEscapeSegments(key string) string {
	parts := strings.Split(key, "/")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	// simple join; keys are generated without raw spaces
	return strings.Join(parts, "/")
}

func validateStorageKey(key string) error {
	if key == "" || strings.HasPrefix(key, "/") || strings.HasPrefix(key, "\\") {
		return fmt.Errorf("invalid storage key")
	}
	if strings.Contains(key, "..") {
		return fmt.Errorf("invalid storage key")
	}
	return nil
}
