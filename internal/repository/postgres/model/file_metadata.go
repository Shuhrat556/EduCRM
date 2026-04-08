package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileMetadata is the GORM model for upload metadata (storage-agnostic).
type FileMetadata struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerType  string    `gorm:"not null;size:32;index:idx_file_metadata_owner,priority:1"`
	OwnerID    uuid.UUID `gorm:"type:uuid;not null;index:idx_file_metadata_owner,priority:2"`
	FileName   string    `gorm:"not null;size:512"`
	StorageKey string    `gorm:"not null;size:1024;uniqueIndex"`
	FileURL    string    `gorm:"not null;type:text"`
	MimeType   string    `gorm:"not null;size:255"`
	SizeBytes  int64     `gorm:"not null;column:size_bytes"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (FileMetadata) TableName() string {
	return "file_metadata"
}

// BeforeCreate assigns ID when missing.
func (m *FileMetadata) BeforeCreate(_ *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
