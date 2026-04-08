package repository

import (
	"context"

	"github.com/google/uuid"
)

// UserTeacherLinkRepository maps login users (role teacher) to teachers table rows.
type UserTeacherLinkRepository interface {
	// FindTeacherIDByUserID returns linked teacher id or nil if none.
	FindTeacherIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error)
}
