package notify

import (
	"context"

	"github.com/google/uuid"
)

// PushProvider delivers mobile or web push notifications (FCM, APNs, Web Push, etc.).
// Add an implementation without changing in-app notification persistence or HTTP handlers.
type PushProvider interface {
	Send(ctx context.Context, payload PushPayload) error
}

// PushPayload is a normalized push envelope for future providers.
type PushPayload struct {
	UserID uuid.UUID
	Title  string
	Body   string
	Data   map[string]string
}

// NoopPushProvider satisfies PushProvider without sending.
type NoopPushProvider struct{}

// Send implements PushProvider as a no-op.
func (NoopPushProvider) Send(ctx context.Context, payload PushPayload) error {
	_ = ctx
	_ = payload
	return nil
}
