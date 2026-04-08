package repository

import (
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// GroupListParams filters paginated group listing.
type GroupListParams struct {
	Search    string
	Status    *domain.GroupStatus
	TeacherID *uuid.UUID
	SubjectID *uuid.UUID
	RoomID    *uuid.UUID
	Page      int
	PageSize  int
}
