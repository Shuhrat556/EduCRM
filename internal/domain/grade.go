package domain

import (
	"time"

	"github.com/google/uuid"
)

// Grade is a weekly rating for a student in a group from a teacher perspective.
type Grade struct {
	ID            uuid.UUID
	StudentID     uuid.UUID
	TeacherID     uuid.UUID // group's assigned teacher (teachers.id)
	GroupID       uuid.UUID
	SubjectID     uuid.UUID
	WeekStartDate time.Time // Monday UTC; defines the rating week
	GradeType     GradeType
	GradeValue    float64
	Comment       *string
	GradedAt      time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
