package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RefreshTokenRepository stores opaque refresh token hashes for rotation and logout.
type RefreshTokenRepository interface {
	Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error
	FindValidByHash(ctx context.Context, tokenHash string) (userID uuid.UUID, err error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteAllForUser(ctx context.Context, userID uuid.UUID) error
	Replace(ctx context.Context, userID uuid.UUID, oldHash, newHash string, expiresAt time.Time) error
}
