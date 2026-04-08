package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Notification is an in-app notification for a single user.
type Notification struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Type      NotificationType
	Title     string
	Body      string
	ReadAt    *time.Time
	Metadata  json.RawMessage // optional JSON (e.g. telegram_chat_id, deep links for push)
	CreatedAt time.Time
	UpdatedAt time.Time
}
