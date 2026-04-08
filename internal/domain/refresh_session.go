package domain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshSession represents a stored refresh token session (opaque token hashed at rest).
type RefreshSession struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	CreatedAt time.Time
}
