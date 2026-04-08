package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// RegisterFileMetadataRequest records an upload stored outside this API (S3/MinIO).
type RegisterFileMetadataRequest struct {
	OwnerType  string    `json:"owner_type" binding:"required,oneof=student_photo teacher_photo document"`
	OwnerID    uuid.UUID `json:"owner_id" binding:"required"`
	FileName   string    `json:"file_name" binding:"required,min=1,max=512"`
	FileURL    string    `json:"file_url" binding:"required,max=4096"`
	MimeType   string    `json:"mime_type" binding:"required,max=255"`
	Size       int64     `json:"size" binding:"min=0"`
	StorageKey *string   `json:"storage_key" binding:"omitempty,max=1024"`
}

// FileMetadataResponse is the public shape (no storage_key).
type FileMetadataResponse struct {
	ID        uuid.UUID `json:"id"`
	OwnerType string    `json:"owner_type"`
	OwnerID   uuid.UUID `json:"owner_id"`
	FileName  string    `json:"file_name"`
	FileURL   string    `json:"file_url"`
	MimeType  string    `json:"mime_type"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// FileMetadataListResponse paginates file metadata.
type FileMetadataListResponse struct {
	Items    []FileMetadataResponse `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// FileMetadataResponseFrom maps domain to DTO.
func FileMetadataResponseFrom(m *domain.FileMetadata) FileMetadataResponse {
	if m == nil {
		return FileMetadataResponse{}
	}
	return FileMetadataResponse{
		ID:        m.ID,
		OwnerType: string(m.OwnerType),
		OwnerID:   m.OwnerID,
		FileName:  m.FileName,
		FileURL:   m.FileURL,
		MimeType:  m.MimeType,
		Size:      m.Size,
		CreatedAt: m.CreatedAt,
	}
}

// FileMetadataListResponseFrom maps list results.
func FileMetadataListResponseFrom(items []domain.FileMetadata, total int64, page, pageSize int) FileMetadataListResponse {
	out := make([]FileMetadataResponse, 0, len(items))
	for i := range items {
		out = append(out, FileMetadataResponseFrom(&items[i]))
	}
	return FileMetadataListResponse{
		Items:    out,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
