package repository

import (
	"context"

	"github.com/google/uuid"
)

// TeacherAssignmentRepository checks teacher↔group↔subject assignments.
type TeacherAssignmentRepository interface {
	// Exists reports whether the teacher is assigned to teach this subject in this group.
	Exists(ctx context.Context, teacherID, groupID, subjectID uuid.UUID) (bool, error)
	// HasAnyAssignmentOnGroup reports whether the teacher has any subject assignment in the group.
	HasAnyAssignmentOnGroup(ctx context.Context, teacherID, groupID uuid.UUID) (bool, error)
	// ListByTeacher returns group and subject pairs.
	ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]TeacherAssignmentRow, error)
}

// TeacherAssignmentRow is a minimal projection for instructor portals.
type TeacherAssignmentRow struct {
	GroupID   uuid.UUID
	SubjectID uuid.UUID
}
