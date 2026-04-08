package repository

import "github.com/educrm/educrm-backend/internal/domain"

// TeacherListParams filters paginated teacher listing.
type TeacherListParams struct {
	Search   string
	Status   *domain.TeacherStatus
	Page     int
	PageSize int
}

// TeacherListEntry is one row in a teacher list (with group count).
type TeacherListEntry struct {
	Teacher    domain.Teacher
	GroupCount int
}
