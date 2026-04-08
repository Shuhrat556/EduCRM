package repository

import (
	"context"

	"github.com/google/uuid"
)

// StudentMembershipRepository reads student_group_memberships.
type StudentMembershipRepository interface {
	// FindGroupIDByStudentUserID returns the student's group id or nil if not enrolled.
	FindGroupIDByStudentUserID(ctx context.Context, studentUserID uuid.UUID) (*uuid.UUID, error)
}
