package dto

import (
	"encoding/json"
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// CreateNotificationRequest is the body for POST /notifications (staff).
type CreateNotificationRequest struct {
	UserID   uuid.UUID       `json:"user_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
	Type     string          `json:"type" example:"grade_posted" binding:"required"`
	Title    string          `json:"title" binding:"required,min=1,max=512"`
	Body     string          `json:"body" binding:"required"`
	Metadata json.RawMessage `json:"metadata" swaggertype:"object"`
}

// NotificationResponse is one in-app notification.
type NotificationResponse struct {
	ID        uuid.UUID       `json:"id"`
	UserID    uuid.UUID       `json:"user_id"`
	Type      string          `json:"type"`
	Title     string          `json:"title"`
	Body      string          `json:"body"`
	ReadAt    *time.Time      `json:"read_at,omitempty"`
	Metadata  json.RawMessage `json:"metadata,omitempty" swaggertype:"object"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// NotificationListResponse is a paginated list.
type NotificationListResponse struct {
	Items    []NotificationResponse `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

// NotificationResponseFrom maps domain to DTO.
func NotificationResponseFrom(n *domain.Notification) NotificationResponse {
	if n == nil {
		return NotificationResponse{}
	}
	var meta json.RawMessage
	if len(n.Metadata) > 0 {
		meta = append(json.RawMessage(nil), n.Metadata...)
	}
	return NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      string(n.Type),
		Title:     n.Title,
		Body:      n.Body,
		ReadAt:    n.ReadAt,
		Metadata:  meta,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

// NotificationListResponseFrom builds list DTO.
func NotificationListResponseFrom(items []domain.Notification, total int64, page, pageSize int) NotificationListResponse {
	out := make([]NotificationResponse, 0, len(items))
	for i := range items {
		out = append(out, NotificationResponseFrom(&items[i]))
	}
	return NotificationListResponse{
		Items:    out,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
