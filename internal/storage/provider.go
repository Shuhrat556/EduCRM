package storage

import (
	"context"
	"io"
)

// Provider persists blobs under an opaque storage key and exposes a public URL.
// Local and object stores (S3, MinIO) implement the same contract so the app
// can switch backends without changing file-metadata use cases.
type Provider interface {
	// Put writes the object at key; key must be relative (no ".." or absolute paths).
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	// Delete removes the object at key; idempotent (missing object is OK).
	Delete(ctx context.Context, key string) error
	// PublicURL is the client-facing URL for key (CDN, path, or presigned base).
	PublicURL(key string) string
}
