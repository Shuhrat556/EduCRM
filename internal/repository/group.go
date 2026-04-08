package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// GroupRepository persists class groups.
type GroupRepository interface {
	Create(ctx context.Context, g *domain.Group) error
	Update(ctx context.Context, g *domain.Group) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Group, error)
	List(ctx context.Context, p GroupListParams) ([]domain.Group, int64, error)
}
