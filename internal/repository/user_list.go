package repository

import "github.com/educrm/educrm-backend/internal/domain"

// UserListParams filters and paginates user listing.
type UserListParams struct {
	Search       string
	Role         *domain.Role
	Page         int
	PageSize     int
	ExcludeRoles []domain.Role
}
