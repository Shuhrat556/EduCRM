package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// FileMetadataListParams filters file_metadata rows.
type FileMetadataListParams struct {
	OwnerType string
	OwnerID   uuid.UUID
	Page      int
	PageSize  int
}

// FileMetadataRepository persists upload metadata independent of storage backend.
type FileMetadataRepository interface {
	Create(ctx context.Context, m *domain.FileMetadata) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.FileMetadata, error)
	List(ctx context.Context, p FileMetadataListParams) ([]domain.FileMetadata, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
