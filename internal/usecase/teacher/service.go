package teacher

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service orchestrates teacher use cases.
type Service struct {
	teachers repository.TeacherRepository
}

// NewService constructs a teacher service.
func NewService(teachers repository.TeacherRepository) *Service {
	return &Service{teachers: teachers}
}

// Detail is teacher + assigned groups (via groups.teacher_id).
type Detail struct {
	Teacher *domain.Teacher
	Groups  []domain.GroupBrief
}

// ListResult is a paginated teacher list.
type ListResult struct {
	Items    []ListItem `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

// ListItem is one row in list responses.
type ListItem struct {
	Teacher    domain.Teacher `json:"-"`
	GroupCount int            `json:"group_count"`
}

// CreateInput holds create payload.
type CreateInput struct {
	FullName          string
	Phone             *string
	Email             *string
	Specialization    *string
	PhotoURL          *string
	PhotoStorageKey   *string
	PhotoContentType  *string
	PhotoOriginalName *string
	Status            domain.TeacherStatus
}

// UpdateInput holds optional updates.
type UpdateInput struct {
	FullName          *string
	Phone             *string
	Email             *string
	Specialization    *string
	PhotoURL          *string
	PhotoStorageKey   *string
	PhotoContentType  *string
	PhotoOriginalName *string
	Status            *domain.TeacherStatus
}

// PhotoPatchInput updates only photo-related metadata.
type PhotoPatchInput struct {
	PhotoURL          *string
	PhotoStorageKey   *string
	PhotoContentType  *string
	PhotoOriginalName *string
}

// Create validates and stores a teacher.
func (s *Service) Create(ctx context.Context, in CreateInput) (*Detail, error) {
	email := domain.NormalizeTeacherEmail(in.Email)
	phone := domain.NormalizeTeacherPhone(in.Phone)
	if err := s.assertUniqueContact(ctx, email, phone, nil); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	t := &domain.Teacher{
		ID:                uuid.New(),
		FullName:          in.FullName,
		Phone:             phone,
		Email:             email,
		Specialization:    trimPtr(in.Specialization),
		PhotoURL:          trimPtr(in.PhotoURL),
		PhotoStorageKey:   trimPtr(in.PhotoStorageKey),
		PhotoContentType:  trimPtr(in.PhotoContentType),
		PhotoOriginalName: trimPtr(in.PhotoOriginalName),
		Status:            in.Status,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := s.teachers.Create(ctx, t); err != nil {
		return nil, wrapTeacherErr(err)
	}
	return s.GetByID(ctx, t.ID)
}

// GetByID returns teacher detail with groups.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*Detail, error) {
	t, groups, err := s.teachers.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load teacher").Wrap(err)
	}
	if t == nil {
		return nil, apperror.NotFound("teacher")
	}
	return &Detail{Teacher: t, Groups: groups}, nil
}

// List returns a page of teachers.
func (s *Service) List(ctx context.Context, params repository.TeacherListParams) (*ListResult, error) {
	rows, total, err := s.teachers.List(ctx, params)
	if err != nil {
		return nil, apperror.Internal("list teachers").Wrap(err)
	}
	items := make([]ListItem, 0, len(rows))
	for i := range rows {
		items = append(items, ListItem{Teacher: rows[i].Teacher, GroupCount: rows[i].GroupCount})
	}
	page := params.Page
	if page < 1 {
		page = 1
	}
	size := params.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return &ListResult{Items: items, Total: total, Page: page, PageSize: size}, nil
}

// Update applies partial updates.
func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*Detail, error) {
	t, _, err := s.teachers.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load teacher").Wrap(err)
	}
	if t == nil {
		return nil, apperror.NotFound("teacher")
	}
	if in.FullName != nil {
		t.FullName = strings.TrimSpace(*in.FullName)
	}
	if in.Phone != nil {
		t.Phone = domain.NormalizeTeacherPhone(in.Phone)
	}
	if in.Email != nil {
		t.Email = domain.NormalizeTeacherEmail(in.Email)
	}
	if in.Specialization != nil {
		t.Specialization = trimPtr(in.Specialization)
	}
	if in.PhotoURL != nil {
		t.PhotoURL = trimPtr(in.PhotoURL)
	}
	if in.PhotoStorageKey != nil {
		t.PhotoStorageKey = trimPtr(in.PhotoStorageKey)
	}
	if in.PhotoContentType != nil {
		t.PhotoContentType = trimPtr(in.PhotoContentType)
	}
	if in.PhotoOriginalName != nil {
		t.PhotoOriginalName = trimPtr(in.PhotoOriginalName)
	}
	if in.Status != nil {
		t.Status = *in.Status
	}
	t.UpdatedAt = time.Now().UTC()
	if err := s.assertUniqueContact(ctx, t.Email, t.Phone, &id); err != nil {
		return nil, err
	}
	if err := s.teachers.Update(ctx, t); err != nil {
		return nil, wrapTeacherErr(err)
	}
	return s.GetByID(ctx, id)
}

// PatchPhoto updates only photo metadata (after client upload to storage).
func (s *Service) PatchPhoto(ctx context.Context, id uuid.UUID, in PhotoPatchInput) (*Detail, error) {
	if !photoPatchAnyFieldSet(in) {
		return nil, apperror.Validation("photo", "at least one photo field is required")
	}
	t, _, err := s.teachers.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load teacher").Wrap(err)
	}
	if t == nil {
		return nil, apperror.NotFound("teacher")
	}
	if in.PhotoURL != nil {
		t.PhotoURL = trimPtr(in.PhotoURL)
	}
	if in.PhotoStorageKey != nil {
		t.PhotoStorageKey = trimPtr(in.PhotoStorageKey)
	}
	if in.PhotoContentType != nil {
		t.PhotoContentType = trimPtr(in.PhotoContentType)
	}
	if in.PhotoOriginalName != nil {
		t.PhotoOriginalName = trimPtr(in.PhotoOriginalName)
	}
	t.UpdatedAt = time.Now().UTC()
	if err := s.teachers.Update(ctx, t); err != nil {
		return nil, wrapTeacherErr(err)
	}
	return s.GetByID(ctx, id)
}

// Delete removes a teacher. Fails if any group still references this teacher.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.teachers.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("teacher")
		}
		if errors.Is(err, repository.ErrReferenced) {
			return apperror.Conflict("teacher_in_use", "cannot delete teacher assigned to one or more groups")
		}
		return apperror.Internal("delete teacher").Wrap(err)
	}
	return nil
}

func (s *Service) assertUniqueContact(ctx context.Context, email, phone *string, excludeID *uuid.UUID) error {
	if email != nil {
		taken, err := s.teachers.EmailTaken(ctx, *email, excludeID)
		if err != nil {
			return apperror.Internal("check email").Wrap(err)
		}
		if taken {
			return apperror.Conflict("email_taken", "email is already in use")
		}
	}
	if phone != nil {
		taken, err := s.teachers.PhoneTaken(ctx, *phone, excludeID)
		if err != nil {
			return apperror.Internal("check phone").Wrap(err)
		}
		if taken {
			return apperror.Conflict("phone_taken", "phone is already in use")
		}
	}
	return nil
}

func wrapTeacherErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, repository.ErrDuplicate) {
		return apperror.Conflict("unique_violation", "email or phone already in use")
	}
	return apperror.Internal("persist teacher").Wrap(err)
}

func trimPtr(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}

func photoPatchAnyFieldSet(in PhotoPatchInput) bool {
	return in.PhotoURL != nil || in.PhotoStorageKey != nil ||
		in.PhotoContentType != nil || in.PhotoOriginalName != nil
}
