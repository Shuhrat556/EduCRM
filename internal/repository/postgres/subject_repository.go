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

// SubjectRepository implements repository.SubjectRepository.
type SubjectRepository struct {
	db *gorm.DB
}

var _ repository.SubjectRepository = (*SubjectRepository)(nil)

// NewSubjectRepository constructs SubjectRepository.
func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

// Exists reports whether a subject row exists.
func (r *SubjectRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&model.Subject{}).Where("id = ?", id).Count(&n).Error
	return n > 0, err
}

// FindByID loads a subject by id.
func (r *SubjectRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error) {
	var m model.Subject
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return subjectToDomain(&m)
}

func subjectToDomain(m *model.Subject) (*domain.Subject, error) {
	st, err := domain.ParseSubjectStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return &domain.Subject{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Code:        m.Code,
		Status:      st,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}
