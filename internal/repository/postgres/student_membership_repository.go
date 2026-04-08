package postgres

import (
	"context"
	"errors"

	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// StudentMembershipRepository implements repository.StudentMembershipRepository.
type StudentMembershipRepository struct {
	db *gorm.DB
}

var _ repository.StudentMembershipRepository = (*StudentMembershipRepository)(nil)

// NewStudentMembershipRepository constructs StudentMembershipRepository.
func NewStudentMembershipRepository(db *gorm.DB) *StudentMembershipRepository {
	return &StudentMembershipRepository{db: db}
}

// FindGroupIDByStudentUserID returns enrolled group id or nil.
func (r *StudentMembershipRepository) FindGroupIDByStudentUserID(ctx context.Context, studentUserID uuid.UUID) (*uuid.UUID, error) {
	var m model.StudentGroupMembership
	if err := r.db.WithContext(ctx).First(&m, "user_id = ?", studentUserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m.GroupID, nil
}
