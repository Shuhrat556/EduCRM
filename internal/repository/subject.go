package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// SubjectRepository reads subjects for validation and joins.
type SubjectRepository interface {
	Exists(ctx context.Context, id uuid.UUID) (bool, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Subject, error)
}
