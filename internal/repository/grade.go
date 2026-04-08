package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// GradeRepository persists weekly grades.
type GradeRepository interface {
	Create(ctx context.Context, g *domain.Grade) error
	Update(ctx context.Context, g *domain.Grade) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Grade, error)
	// ListByStudent orders by week_start_date DESC. viewerTeacherID limits to groups that teacher teaches; nil means no teacher filter (admin or student listing self).
	ListByStudent(ctx context.Context, studentID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Grade, error)
	// ListByGroup orders by week_start_date DESC, student_id. viewerTeacherID must match group teacher when non-nil.
	ListByGroup(ctx context.Context, groupID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Grade, error)
}
