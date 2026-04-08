package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// RoomRepository persists rooms for scheduling.
type RoomRepository interface {
	Create(ctx context.Context, r *domain.Room) error
	Update(ctx context.Context, r *domain.Room) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Room, error)
	List(ctx context.Context, p RoomListParams) ([]domain.Room, int64, error)
}
