package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// CreateScheduleRequest is the body for POST /schedules.
type CreateScheduleRequest struct {
	GroupID   uuid.UUID `json:"group_id" binding:"required"`
	TeacherID uuid.UUID `json:"teacher_id" binding:"required"`
	RoomID    uuid.UUID `json:"room_id" binding:"required"`
	Weekday   int       `json:"weekday" binding:"required,min=0,max=6"`
	StartTime string    `json:"start_time" binding:"required"` // HH:MM 24h
	EndTime   string    `json:"end_time" binding:"required"`   // HH:MM; may use 24:00 for end of day
}

// UpdateScheduleRequest is the body for PATCH /schedules/:id.
type UpdateScheduleRequest struct {
	GroupID   *uuid.UUID `json:"group_id"`
	TeacherID *uuid.UUID `json:"teacher_id"`
	RoomID    *uuid.UUID `json:"room_id"`
	Weekday   *int       `json:"weekday" binding:"omitempty,min=0,max=6"`
	StartTime *string    `json:"start_time"`
	EndTime   *string    `json:"end_time"`
}

// ScheduleResponse is one schedule row in API responses.
type ScheduleResponse struct {
	ID        uuid.UUID `json:"id"`
	GroupID   uuid.UUID `json:"group_id"`
	TeacherID uuid.UUID `json:"teacher_id"`
	RoomID    uuid.UUID `json:"room_id"`
	Weekday   int       `json:"weekday"`
	StartTime string    `json:"start_time"`
	EndTime   string    `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ScheduleListResponse wraps a list of schedules.
type ScheduleListResponse struct {
	Items []ScheduleResponse `json:"items"`
}

// ScheduleResponseFrom maps domain to DTO.
func ScheduleResponseFrom(s *domain.Schedule) ScheduleResponse {
	if s == nil {
		return ScheduleResponse{}
	}
	return ScheduleResponse{
		ID:        s.ID,
		GroupID:   s.GroupID,
		TeacherID: s.TeacherID,
		RoomID:    s.RoomID,
		Weekday:   int(s.Weekday),
		StartTime: FormatMinutesAsClock(s.StartMinutes),
		EndTime:   FormatMinutesAsClock(s.EndMinutes),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// ScheduleListResponseFrom maps a slice.
func ScheduleListResponseFrom(rows []domain.Schedule) ScheduleListResponse {
	items := make([]ScheduleResponse, 0, len(rows))
	for i := range rows {
		items = append(items, ScheduleResponseFrom(&rows[i]))
	}
	return ScheduleListResponse{Items: items}
}
