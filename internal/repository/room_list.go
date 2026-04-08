package repository

import "github.com/educrm/educrm-backend/internal/domain"

// RoomListParams filters paginated room listing.
type RoomListParams struct {
	Search   string
	Status   *domain.RoomStatus
	Page     int
	PageSize int
}
