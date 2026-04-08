package repository

import (
	"context"
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// AttendanceRepository persists lesson attendance rows.
type AttendanceRepository interface {
	Create(ctx context.Context, a *domain.Attendance) error
	Update(ctx context.Context, a *domain.Attendance) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Attendance, error)
	// ListByStudent returns attendance for a student; if viewerTeacherID is set, only rows for groups taught by that teacher.
	ListByStudent(ctx context.Context, studentID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error)
	// ListByGroup returns attendance for a group; viewerTeacherID must match group teacher when non-nil.
	ListByGroup(ctx context.Context, groupID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error)
	// ListByDateRange returns rows with lesson_date in [from, to] inclusive; optional teacher filter via groups.teacher_id.
	ListByDateRange(ctx context.Context, from, to time.Time, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error)
}
