package postgres

import (
	"context"
	"errors"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileMetadataRepository implements repository.FileMetadataRepository.
type FileMetadataRepository struct {
	db *gorm.DB
}

var _ repository.FileMetadataRepository = (*FileMetadataRepository)(nil)

// NewFileMetadataRepository constructs FileMetadataRepository.
func NewFileMetadataRepository(db *gorm.DB) *FileMetadataRepository {
	return &FileMetadataRepository{db: db}
}

// Create inserts a row.
func (r *FileMetadataRepository) Create(ctx context.Context, m *domain.FileMetadata) error {
	row, err := fileMetadataToModel(m)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(row).Error
}

// FindByID loads by primary key.
func (r *FileMetadataRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.FileMetadata, error) {
	var row model.FileMetadata
	if err := r.db.WithContext(ctx).First(&row, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return fileMetadataToDomain(&row)
}

// List paginates by owner.
func (r *FileMetadataRepository) List(ctx context.Context, p repository.FileMetadataListParams) ([]domain.FileMetadata, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.FileMetadata{}).
			Where("owner_type = ? AND owner_id = ?", p.OwnerType, p.OwnerID)
		return q
	}
	var total int64
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := p.Page
	if page < 1 {
		page = 1
	}
	size := p.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size
	var rows []model.FileMetadata
	if err := build().Order("created_at DESC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]domain.FileMetadata, 0, len(rows))
	for i := range rows {
		d, err := fileMetadataToDomain(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *d)
	}
	return out, total, nil
}

// Delete removes a row.
func (r *FileMetadataRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.FileMetadata{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func fileMetadataToDomain(m *model.FileMetadata) (*domain.FileMetadata, error) {
	ot, err := domain.ParseFileOwnerType(m.OwnerType)
	if err != nil {
		return nil, err
	}
	return &domain.FileMetadata{
		ID:         m.ID,
		OwnerType:  ot,
		OwnerID:    m.OwnerID,
		FileName:   m.FileName,
		StorageKey: m.StorageKey,
		FileURL:    m.FileURL,
		MimeType:   m.MimeType,
		Size:       m.SizeBytes,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}

func fileMetadataToModel(m *domain.FileMetadata) (*model.FileMetadata, error) {
	return &model.FileMetadata{
		ID:         m.ID,
		OwnerType:  string(m.OwnerType),
		OwnerID:    m.OwnerID,
		FileName:   m.FileName,
		StorageKey: m.StorageKey,
		FileURL:    m.FileURL,
		MimeType:   m.MimeType,
		SizeBytes:  m.Size,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}
