package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// NotificationListParams filters notifications for a user.
type NotificationListParams struct {
	UserID     uuid.UUID
	UnreadOnly bool
	Page       int
	PageSize   int
}

// NotificationRepository persists in-app notifications.
type NotificationRepository interface {
	Create(ctx context.Context, n *domain.Notification) error
	Update(ctx context.Context, n *domain.Notification) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Notification, error)
	List(ctx context.Context, p NotificationListParams) ([]domain.Notification, int64, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
