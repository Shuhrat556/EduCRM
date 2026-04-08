package domain

import (
	"time"

	"github.com/google/uuid"
)

// FileMetadata is persisted metadata for a stored blob (any Provider).
type FileMetadata struct {
	ID         uuid.UUID
	OwnerType  FileOwnerType
	OwnerID    uuid.UUID
	FileName   string
	StorageKey string // opaque key for storage.Provider (local path or object key)
	FileURL    string
	MimeType   string
	Size       int64
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
