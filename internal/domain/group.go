package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Group is a class cohort with one teacher, one subject, optional default room, and schedule window.
type Group struct {
	ID               uuid.UUID
	Name             string
	SubjectID        uuid.UUID
	TeacherID        uuid.UUID
	RoomID           *uuid.UUID
	StartDate        time.Time
	EndDate          time.Time
	MonthlyFeeMinor  int64
	Status           GroupStatus
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// GroupBrief is a minimal projection (e.g. teacher’s group list).
type GroupBrief struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// NormalizeGroupName trims the group name.
func NormalizeGroupName(s string) string {
	return strings.TrimSpace(s)
}
