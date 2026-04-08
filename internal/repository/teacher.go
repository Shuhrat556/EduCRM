package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// TeacherRepository persists teachers. Groups reference teachers via groups.teacher_id.
type TeacherRepository interface {
	Create(ctx context.Context, t *domain.Teacher) error
	Update(ctx context.Context, t *domain.Teacher) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Teacher, []domain.GroupBrief, error)
	List(ctx context.Context, p TeacherListParams) ([]TeacherListEntry, int64, error)
	EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
	PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error)
}
