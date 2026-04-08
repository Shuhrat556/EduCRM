package group

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

// Service orchestrates group (class cohort) use cases.
type Service struct {
	groups   repository.GroupRepository
	subjects repository.SubjectRepository
	teachers repository.TeacherRepository
	rooms    repository.RoomRepository
}

// NewService constructs a group service.
func NewService(
	groups repository.GroupRepository,
	subjects repository.SubjectRepository,
	teachers repository.TeacherRepository,
	rooms repository.RoomRepository,
) *Service {
	return &Service{groups: groups, subjects: subjects, teachers: teachers, rooms: rooms}
}

// CreateInput is validated input for creating a group.
type CreateInput struct {
	Name            string
	SubjectID       uuid.UUID
	TeacherID       uuid.UUID
	RoomID          *uuid.UUID
	StartDate       time.Time
	EndDate         time.Time
	MonthlyFeeMinor int64
	Status          domain.GroupStatus
}

// UpdateInput holds optional field updates (nil = leave unchanged).
type UpdateInput struct {
	Name            *string
	SubjectID       *uuid.UUID
	TeacherID       *uuid.UUID
	RoomID          *uuid.UUID // set room when non-nil and ClearRoom is false
	ClearRoom       bool       // when true, clears room_id
	StartDate       *time.Time
	EndDate         *time.Time
	MonthlyFeeMinor *int64
	Status          *domain.GroupStatus
}

// ListResult is a paginated group list.
type ListResult struct {
	Items    []domain.Group
	Total    int64
	Page     int
	PageSize int
}

// Create validates references and inserts a group.
func (s *Service) Create(ctx context.Context, in CreateInput) (*domain.Group, error) {
	name := domain.NormalizeGroupName(in.Name)
	if name == "" {
		return nil, apperror.Validation("name", "name is required")
	}
	if in.MonthlyFeeMinor < 0 {
		return nil, apperror.Validation("monthly_fee", "must be zero or positive (minor units)")
	}
	start := truncateUTCDate(in.StartDate)
	end := truncateUTCDate(in.EndDate)
	if start.After(end) {
		return nil, apperror.Validation("dates", "start_date must be on or before end_date")
	}
	if err := s.assertActiveSubject(ctx, in.SubjectID); err != nil {
		return nil, err
	}
	if err := s.assertTeacherExists(ctx, in.TeacherID); err != nil {
		return nil, err
	}
	if err := s.assertRoomIfSet(ctx, in.RoomID); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	g := &domain.Group{
		ID:              uuid.New(),
		Name:            name,
		SubjectID:       in.SubjectID,
		TeacherID:       in.TeacherID,
		RoomID:          in.RoomID,
		StartDate:       start,
		EndDate:         end,
		MonthlyFeeMinor: in.MonthlyFeeMinor,
		Status:          in.Status,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.groups.Create(ctx, g); err != nil {
		return nil, apperror.Internal("create group").Wrap(err)
	}
	return g, nil
}

// GetByID returns a group or not found.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	g, err := s.groups.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return nil, apperror.NotFound("group")
	}
	return g, nil
}

// List returns a page of groups.
func (s *Service) List(ctx context.Context, p repository.GroupListParams) (*ListResult, error) {
	rows, total, err := s.groups.List(ctx, p)
	if err != nil {
		return nil, apperror.Internal("list groups").Wrap(err)
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
	return &ListResult{Items: rows, Total: total, Page: page, PageSize: size}, nil
}

// Update merges changes and persists.
func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*domain.Group, error) {
	g, err := s.groups.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return nil, apperror.NotFound("group")
	}
	if in.Name != nil {
		n := domain.NormalizeGroupName(*in.Name)
		if n == "" {
			return nil, apperror.Validation("name", "name cannot be empty")
		}
		g.Name = n
	}
	if in.SubjectID != nil {
		if err := s.assertActiveSubject(ctx, *in.SubjectID); err != nil {
			return nil, err
		}
		g.SubjectID = *in.SubjectID
	}
	if in.TeacherID != nil {
		if err := s.assertTeacherExists(ctx, *in.TeacherID); err != nil {
			return nil, err
		}
		g.TeacherID = *in.TeacherID
	}
	if in.ClearRoom {
		g.RoomID = nil
	} else if in.RoomID != nil {
		if err := s.assertRoomIfSet(ctx, in.RoomID); err != nil {
			return nil, err
		}
		g.RoomID = in.RoomID
	}
	if in.StartDate != nil {
		g.StartDate = truncateUTCDate(*in.StartDate)
	}
	if in.EndDate != nil {
		g.EndDate = truncateUTCDate(*in.EndDate)
	}
	if in.MonthlyFeeMinor != nil {
		if *in.MonthlyFeeMinor < 0 {
			return nil, apperror.Validation("monthly_fee", "must be zero or positive (minor units)")
		}
		g.MonthlyFeeMinor = *in.MonthlyFeeMinor
	}
	if in.Status != nil {
		g.Status = *in.Status
	}
	if g.StartDate.After(g.EndDate) {
		return nil, apperror.Validation("dates", "start_date must be on or before end_date")
	}
	g.UpdatedAt = time.Now().UTC()
	if err := s.groups.Update(ctx, g); err != nil {
		return nil, apperror.Internal("update group").Wrap(err)
	}
	return g, nil
}

// Delete removes a group by id.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.groups.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("group")
		}
		return apperror.Internal("delete group").Wrap(err)
	}
	return nil
}

func (s *Service) assertActiveSubject(ctx context.Context, id uuid.UUID) error {
	sub, err := s.subjects.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load subject").Wrap(err)
	}
	if sub == nil {
		return apperror.Validation("subject_id", "subject not found")
	}
	if sub.Status != domain.SubjectStatusActive {
		return apperror.Validation("subject_id", "subject must be active")
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

func (s *Service) assertRoomIfSet(ctx context.Context, id *uuid.UUID) error {
	if id == nil {
		return nil
	}
	r, err := s.rooms.FindByID(ctx, *id)
	if err != nil {
		return apperror.Internal("load room").Wrap(err)
	}
	if r == nil {
		return apperror.Validation("room_id", "room not found")
	}
	return nil
}

func truncateUTCDate(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}
