package domain

import (
	"time"

	"github.com/google/uuid"
)

// Schedule is a recurring weekly slot (same weekday and clock range each week).
type Schedule struct {
	ID           uuid.UUID
	GroupID      uuid.UUID
	TeacherID    uuid.UUID
	RoomID       uuid.UUID
	Weekday      Weekday
	StartMinutes int // minutes from midnight [0, 24*60)
	EndMinutes   int // exclusive upper bound for display often shown as last minute; must be > StartMinutes, max 24*60
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

const MinutesPerDay = 24 * 60
