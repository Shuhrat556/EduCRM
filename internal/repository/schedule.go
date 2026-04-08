package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// ScheduleRepository persists class schedule slots.
type ScheduleRepository interface {
	Create(ctx context.Context, s *domain.Schedule) error
	Update(ctx context.Context, s *domain.Schedule) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error)
	ListByGroup(ctx context.Context, groupID uuid.UUID) ([]domain.Schedule, error)
	ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]domain.Schedule, error)
	ListByRoom(ctx context.Context, roomID uuid.UUID) ([]domain.Schedule, error)
	// CountRoomOverlaps counts other rows (excluding excludeID) with same room, weekday, overlapping [start,end).
	CountRoomOverlaps(ctx context.Context, roomID uuid.UUID, weekday domain.Weekday, startMin, endMin int, excludeID *uuid.UUID) (int64, error)
	// CountTeacherOverlaps counts other rows with same teacher, weekday, overlapping interval.
	CountTeacherOverlaps(ctx context.Context, teacherID uuid.UUID, weekday domain.Weekday, startMin, endMin int, excludeID *uuid.UUID) (int64, error)
}
