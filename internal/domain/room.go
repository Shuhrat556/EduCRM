package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Room is a physical or virtual space used in schedule planning.
type Room struct {
	ID          uuid.UUID
	Name        string
	Capacity    int
	Description *string
	Status      RoomStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NormalizeRoomName trims whitespace for storage.
func NormalizeRoomName(s string) string {
	return strings.TrimSpace(s)
}
