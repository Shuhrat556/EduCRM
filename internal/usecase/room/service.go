package room

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

const maxRoomCapacity = 100_000

// Service orchestrates room use cases.
type Service struct {
	repo repository.RoomRepository
}

// NewService constructs a room service.
func NewService(repo repository.RoomRepository) *Service {
	return &Service{repo: repo}
}

// ListResult is a paginated room list.
type ListResult struct {
	Items    []domain.Room `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// CreateInput holds create payload.
type CreateInput struct {
	Name        string
	Capacity    int
	Description *string
	Status      domain.RoomStatus
}

// UpdateInput holds optional updates.
type UpdateInput struct {
	Name        *string
	Capacity    *int
	Description *string
	Status      *domain.RoomStatus
}

// Create validates and stores a room.
func (s *Service) Create(ctx context.Context, in CreateInput) (*domain.Room, error) {
	if err := validateCapacity(in.Capacity); err != nil {
		return nil, err
	}
	name := domain.NormalizeRoomName(in.Name)
	if name == "" {
		return nil, apperror.Validation("name", "name is required")
	}
	now := time.Now().UTC()
	r := &domain.Room{
		ID:          uuid.New(),
		Name:        name,
		Capacity:    in.Capacity,
		Description: trimDesc(in.Description),
		Status:      in.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.repo.Create(ctx, r); err != nil {
		return nil, apperror.Internal("create room").Wrap(err)
	}
	return r, nil
}

// GetByID returns a room by id.
func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*domain.Room, error) {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load room").Wrap(err)
	}
	if r == nil {
		return nil, apperror.NotFound("room")
	}
	return r, nil
}

// List returns a page of rooms.
func (s *Service) List(ctx context.Context, params repository.RoomListParams) (*ListResult, error) {
	rows, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, apperror.Internal("list rooms").Wrap(err)
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
	return &ListResult{Items: rows, Total: total, Page: page, PageSize: size}, nil
}

// Update applies partial updates.
func (s *Service) Update(ctx context.Context, id uuid.UUID, in UpdateInput) (*domain.Room, error) {
	r, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load room").Wrap(err)
	}
	if r == nil {
		return nil, apperror.NotFound("room")
	}
	if in.Name != nil {
		n := domain.NormalizeRoomName(*in.Name)
		if n == "" {
			return nil, apperror.Validation("name", "name cannot be empty")
		}
		r.Name = n
	}
	if in.Capacity != nil {
		if err := validateCapacity(*in.Capacity); err != nil {
			return nil, err
		}
		r.Capacity = *in.Capacity
	}
	if in.Description != nil {
		r.Description = trimDesc(in.Description)
	}
	if in.Status != nil {
		r.Status = *in.Status
	}
	r.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, r); err != nil {
		return nil, apperror.Internal("update room").Wrap(err)
	}
	return r, nil
}

// Delete removes a room.
func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("room")
		}
		return apperror.Internal("delete room").Wrap(err)
	}
	return nil
}

func validateCapacity(n int) error {
	if n < 1 {
		return apperror.Validation("capacity", "capacity must be at least 1")
	}
	if n > maxRoomCapacity {
		return apperror.Validation("capacity", "capacity exceeds maximum allowed")
	}
	return nil
}

func trimDesc(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}
