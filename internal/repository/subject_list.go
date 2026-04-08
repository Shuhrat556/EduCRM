package repository

import "github.com/educrm/educrm-backend/internal/domain"

// SubjectListParams filters paginated subject listing.
type SubjectListParams struct {
	Search   string
	Status   *domain.SubjectStatus
	Page     int
	PageSize int
}
