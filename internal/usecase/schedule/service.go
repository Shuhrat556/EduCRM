package schedule

import (
	"context"
	"errors"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service orchestrates schedule use cases.
type Service struct {
	schedules repository.ScheduleRepository
	groups    repository.GroupRepository
	teachers  repository.TeacherRepository
	rooms     repository.RoomRepository
}

// NewService constructs a schedule service.
func NewService(
	schedules repository.ScheduleRepository,
	groups repository.GroupRepository,
	teachers repository.TeacherRepository,
	rooms repository.RoomRepository,
) *Service {
	return &Service{schedules: schedules, groups: groups, teachers: teachers, rooms: rooms}
}

// CreateInput holds validated minutes and weekday.
type CreateInput struct {
	GroupID      uuid.UUID
	TeacherID    uuid.UUID
	RoomID       uuid.UUID
	Weekday      domain.Weekday
	StartMinutes int
	EndMinutes   int
}

// UpdateInput holds optional updates (nil = unchanged).
type UpdateInput struct {
	GroupID         *uuid.UUID
	TeacherID       *uuid.UUID
	RoomID          *uuid.UUID
	Weekday         *domain.Weekday
	StartMinutes    *int
	EndMinutes      *int
}

// ListFilter selects which dimension to list by (exactly one must be set by handler).
type ListFilter struct {
	GroupID   *uuid.UUID
	TeacherID *uuid.UUID
	RoomID    *uuid.UUID
}

// Create validates group/teacher/room, conflicts, and inserts.
func (s *Service) Create(ctx context.Context, in CreateInput) (*domain.Schedule, error) {
	if err := validateTimeRange(in.StartMinutes, in.EndMinutes); err != nil {
		return nil, err
	}
	if err := s.assertGroupExists(ctx, in.GroupID); err != nil {
		return nil, err
	}
	if err := s.assertTeacherExists(ctx, in.TeacherID); err != nil {
		return nil, err
	}
	if err := s.assertRoomExists(ctx, in.RoomID); err != nil {
		return nil, err
	}
	if err := s.assertNoConflicts(ctx, in.RoomID, in.TeacherID, in.Weekday, in.StartMinutes, in.EndMinutes, nil); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	row := &domain.Schedule{
		ID:           uuid.New(),
		GroupID:      in.GroupID,
		TeacherID:    in.TeacherID,
		RoomID:       in.RoomID,
		Weekday:      in.Weekday,
		StartMinutes: in.StartMinutes,
		EndMinutes:   in.EndMinutes,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.schedules.Create(ctx, row); err != nil {
		return nil, apperror.Internal("create schedule").Wrap(err)
	}
	return row, nil
}

// GetByID returns a schedule or not found.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Schedule, error) {
	row, err := s.schedules.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load schedule").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("schedule")
	}
	return row, nil
}

// List returns schedules for one dimension.
func (s *Service) List(ctx context.Context, f ListFilter) ([]domain.Schedule, error) {
	n := 0
	if f.GroupID != nil {
		n++
	}
	if f.TeacherID != nil {
		n++
	}
	if f.RoomID != nil {
		n++
	}
	if n != 1 {
		return nil, apperror.Validation("filter", "provide exactly one of group_id, teacher_id, room_id")
	}
	var (
		rows []domain.Schedule
		err  error
	)
	switch {
	case f.GroupID != nil:
		rows, err = s.schedules.ListByGroup(ctx, *f.GroupID)
	case f.TeacherID != nil:
		rows, err = s.schedules.ListByTeacher(ctx, *f.TeacherID)
	default:
		rows, err = s.schedules.ListByRoom(ctx, *f.RoomID)
	}
	if err != nil {
		return nil, apperror.Internal("list schedules").Wrap(err)
	}
	return rows, nil
}

// Update merges and validates conflicts.
func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*domain.Schedule, error) {
	row, err := s.schedules.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load schedule").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("schedule")
	}
	if in.GroupID != nil {
		if err := s.assertGroupExists(ctx, *in.GroupID); err != nil {
			return nil, err
		}
		row.GroupID = *in.GroupID
	}
	if in.TeacherID != nil {
		if err := s.assertTeacherExists(ctx, *in.TeacherID); err != nil {
			return nil, err
		}
		row.TeacherID = *in.TeacherID
	}
	if in.RoomID != nil {
		if err := s.assertRoomExists(ctx, *in.RoomID); err != nil {
			return nil, err
		}
		row.RoomID = *in.RoomID
	}
	if in.Weekday != nil {
		row.Weekday = *in.Weekday
	}
	if in.StartMinutes != nil {
		row.StartMinutes = *in.StartMinutes
	}
	if in.EndMinutes != nil {
		row.EndMinutes = *in.EndMinutes
	}
	if err := validateTimeRange(row.StartMinutes, row.EndMinutes); err != nil {
		return nil, err
	}
	if err := s.assertNoConflicts(ctx, row.RoomID, row.TeacherID, row.Weekday, row.StartMinutes, row.EndMinutes, &id); err != nil {
		return nil, err
	}
	row.UpdatedAt = time.Now().UTC()
	if err := s.schedules.Update(ctx, row); err != nil {
		return nil, apperror.Internal("update schedule").Wrap(err)
	}
	return row, nil
}

// Delete removes a schedule.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.schedules.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("schedule")
		}
		return apperror.Internal("delete schedule").Wrap(err)
	}
	return nil
}

func validateTimeRange(startMin, endMin int) error {
	if startMin < 0 || startMin >= domain.MinutesPerDay {
		return apperror.Validation("start_time", "must be within the same calendar day (00:00–23:59)")
	}
	if endMin <= startMin || endMin > domain.MinutesPerDay {
		return apperror.Validation("end_time", "must be after start_time and not past 24:00")
	}
	return nil
}

func (s *Service) assertGroupExists(ctx context.Context, id uuid.UUID) error {
	g, err := s.groups.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return apperror.Validation("group_id", "group not found")
	}
	return nil
}

func (s *Service) assertTeacherExists(ctx context.Context, id uuid.UUID) error {
	t, _, err := s.teachers.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load teacher").Wrap(err)
	}
	if t == nil {
		return apperror.Validation("teacher_id", "teacher not found")
	}
	return nil
}

func (s *Service) assertRoomExists(ctx context.Context, id uuid.UUID) error {
	r, err := s.rooms.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load room").Wrap(err)
	}
	if r == nil {
		return apperror.Validation("room_id", "room not found")
	}
	return nil
}

func (s *Service) assertNoConflicts(ctx context.Context, roomID, teacherID uuid.UUID, wd domain.Weekday, startMin, endMin int, excludeID *uuid.UUID) error {
	nRoom, err := s.schedules.CountRoomOverlaps(ctx, roomID, wd, startMin, endMin, excludeID)
	if err != nil {
		return apperror.Internal("check room availability").Wrap(err)
	}
	if nRoom > 0 {
		return apperror.Conflict("room_schedule_conflict", "another session already uses this room at that weekday and time")
	}
	nTeach, err := s.schedules.CountTeacherOverlaps(ctx, teacherID, wd, startMin, endMin, excludeID)
	if err != nil {
		return apperror.Internal("check teacher availability").Wrap(err)
	}
	if nTeach > 0 {
		return apperror.Conflict("teacher_schedule_conflict", "this teacher is already scheduled at that weekday and time")
	}
	return nil
}
