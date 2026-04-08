package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	roomsvc "github.com/educrm/educrm-backend/internal/usecase/room"
	"github.com/google/uuid"
)

// CreateRoomRequest is the body for POST /rooms.
type CreateRoomRequest struct {
	Name        string  `json:"name" example:"Lab A" binding:"required,min=1,max=255"`
	Capacity    int     `json:"capacity" example:"24" binding:"required,min=1,max=100000"`
	Description *string `json:"description" example:"Chemistry lab" binding:"omitempty,max=5000"`
	Status      string  `json:"status" example:"active" binding:"omitempty,oneof=active inactive"`
}

// UpdateRoomRequest is the body for PATCH /rooms/:id.
type UpdateRoomRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=255"`
	Capacity    *int    `json:"capacity" binding:"omitempty,min=1,max=100000"`
	Description *string `json:"description" binding:"omitempty,max=5000"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// RoomResponse is the API shape for a room.
type RoomResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Capacity    int       `json:"capacity"`
	Description *string   `json:"description,omitempty"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RoomListResponse is the data payload for GET /rooms.
type RoomListResponse struct {
	Items    []RoomResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// RoomResponseFrom maps domain to DTO.
func RoomResponseFrom(r *domain.Room) RoomResponse {
	if r == nil {
		return RoomResponse{}
	}
	return RoomResponse{
		ID:          r.ID,
		Name:        r.Name,
		Capacity:    r.Capacity,
		Description: r.Description,
		Status:      string(r.Status),
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

// RoomListResponseFrom maps list result.
func RoomListResponseFrom(res *roomsvc.ListResult) RoomListResponse {
	if res == nil {
		return RoomListResponse{}
	}
	items := make([]RoomResponse, 0, len(res.Items))
	for i := range res.Items {
		items = append(items, RoomResponseFrom(&res.Items[i]))
	}
	return RoomListResponse{
		Items:    items,
		Total:    res.Total,
		Page:     res.Page,
		PageSize: res.PageSize,
	}
}
