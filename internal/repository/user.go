package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// UserRepository loads and persists users.
type UserRepository interface {
	FindByLogin(ctx context.Context, login string) (*domain.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Create(ctx context.Context, u *domain.User) error
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, p UserListParams) ([]domain.User, int64, error)
	EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error)
	PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error)
}
