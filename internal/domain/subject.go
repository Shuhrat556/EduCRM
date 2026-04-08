package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Subject is a reusable course/subject catalog entry for groups and scheduling.
type Subject struct {
	ID          uuid.UUID
	Name        string
	Description *string
	Code        string
	Status      SubjectStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NormalizeSubjectCode returns trimmed uppercase code for storage and uniqueness.
func NormalizeSubjectCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
