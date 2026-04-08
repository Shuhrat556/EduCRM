package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	groupsvc "github.com/educrm/educrm-backend/internal/usecase/group"
	"github.com/google/uuid"
)

// CreateGroupRequest is the body for POST /groups.
type CreateGroupRequest struct {
	Name            string      `json:"name" binding:"required,min=1,max=255"`
	SubjectID       uuid.UUID   `json:"subject_id" binding:"required"`
	TeacherID       uuid.UUID   `json:"teacher_id" binding:"required"`
	RoomID          *uuid.UUID  `json:"room_id"`
	StartDate       string      `json:"start_date" binding:"required"`
	EndDate         string      `json:"end_date" binding:"required"`
	MonthlyFee      int64       `json:"monthly_fee" binding:"min=0"` // minor units (e.g. cents)
	Status          string      `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateGroupRequest is the body for PATCH /groups/:id.
type UpdateGroupRequest struct {
	Name            *string     `json:"name" binding:"omitempty,min=1,max=255"`
	SubjectID       *uuid.UUID  `json:"subject_id"`
	TeacherID       *uuid.UUID  `json:"teacher_id"`
	RoomID          *uuid.UUID  `json:"room_id"`
	ClearRoom       *bool       `json:"clear_room"`
	StartDate       *string     `json:"start_date"`
	EndDate         *string     `json:"end_date"`
	MonthlyFee      *int64      `json:"monthly_fee" binding:"omitempty,min=0"`
	Status          *string     `json:"status" binding:"omitempty,oneof=active inactive"`
}

// GroupResponse is a single group in API responses.
type GroupResponse struct {
	ID         uuid.UUID  `json:"id"`
	Name       string     `json:"name"`
	SubjectID  uuid.UUID  `json:"subject_id"`
	TeacherID  uuid.UUID  `json:"teacher_id"`
	RoomID     *uuid.UUID `json:"room_id,omitempty"`
	StartDate  string     `json:"start_date"`
	EndDate    string     `json:"end_date"`
	MonthlyFee int64      `json:"monthly_fee"`
	Status     string     `json:"status"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// GroupListResponse is the data envelope for GET /groups.
type GroupListResponse struct {
	Items    []GroupResponse `json:"items"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// GroupResponseFrom maps domain group to API DTO.
func GroupResponseFrom(g *domain.Group) GroupResponse {
	if g == nil {
		return GroupResponse{}
	}
	return GroupResponse{
		ID:         g.ID,
		Name:       g.Name,
		SubjectID:  g.SubjectID,
		TeacherID:  g.TeacherID,
		RoomID:     g.RoomID,
		StartDate:  g.StartDate.UTC().Format("2006-01-02"),
		EndDate:    g.EndDate.UTC().Format("2006-01-02"),
		MonthlyFee: g.MonthlyFeeMinor,
		Status:     string(g.Status),
		CreatedAt:  g.CreatedAt,
		UpdatedAt:  g.UpdatedAt,
	}
}

// GroupListResponseFrom maps list result.
func GroupListResponseFrom(r *groupsvc.ListResult) GroupListResponse {
	if r == nil {
		return GroupListResponse{}
	}
	items := make([]GroupResponse, 0, len(r.Items))
	for i := range r.Items {
		items = append(items, GroupResponseFrom(&r.Items[i]))
	}
	return GroupListResponse{
		Items:    items,
		Total:    r.Total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}
}
