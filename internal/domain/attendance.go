package domain

import (
	"time"

	"github.com/google/uuid"
)

// Attendance records one student's presence for a group lesson on a calendar date.
type Attendance struct {
	ID                  uuid.UUID
	StudentID           uuid.UUID
	GroupID             uuid.UUID
	LessonDate          time.Time // calendar date (UTC)
	Status              AttendanceStatus
	Comment             *string
	MarkedByTeacherID   uuid.UUID // always the group's assigned teacher at time of marking
	CreatedAt           time.Time
	UpdatedAt           time.Time
}
