package postgres

import (
	"context"
	"errors"

	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserTeacherLinkRepository implements repository.UserTeacherLinkRepository.
type UserTeacherLinkRepository struct {
	db *gorm.DB
}

var _ repository.UserTeacherLinkRepository = (*UserTeacherLinkRepository)(nil)

// NewUserTeacherLinkRepository constructs UserTeacherLinkRepository.
func NewUserTeacherLinkRepository(db *gorm.DB) *UserTeacherLinkRepository {
	return &UserTeacherLinkRepository{db: db}
}

// FindTeacherIDByUserID returns linked teacher id or nil.
func (r *UserTeacherLinkRepository) FindTeacherIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	var m model.UserTeacherLink
	if err := r.db.WithContext(ctx).First(&m, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m.TeacherID, nil
}
