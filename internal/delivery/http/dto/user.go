package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/usecase/user"
	"github.com/google/uuid"
)

// CreateUserRequest is the body for POST /users.
type CreateUserRequest struct {
	FullName            string  `json:"full_name" example:"Jane Doe" binding:"required,min=1,max=255"`
	Username            *string `json:"username" example:"jstudent01" binding:"omitempty,min=3,max=64"`
	Email               *string `json:"email" example:"student@school.edu" binding:"omitempty,email,max=255"`
	Phone               *string `json:"phone" example:"+15551234567" binding:"omitempty,min=5,max=32"`
	Password            string  `json:"password" example:"min-8-chars" binding:"required,min=8,max=128"`
	Role                string  `json:"role" example:"student" binding:"required"`
	ForcePasswordChange bool    `json:"force_password_change"`
}

// UpdateUserRequest is the body for PATCH /users/:id (all fields optional).
type UpdateUserRequest struct {
	FullName *string `json:"full_name" binding:"omitempty,min=1,max=255"`
	Username *string `json:"username" binding:"omitempty,min=3,max=64"`
	Email    *string `json:"email" binding:"omitempty,email,max=255"`
	Phone    *string `json:"phone" binding:"omitempty,min=5,max=32"`
	Password *string `json:"password" binding:"omitempty,min=8,max=128"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

// SetUserStatusRequest toggles activation.
type SetUserStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

// UserResponse is a public user representation.
type UserResponse struct {
	ID                  uuid.UUID `json:"id"`
	FullName            string    `json:"full_name"`
	Username            *string   `json:"username,omitempty"`
	Email               *string   `json:"email,omitempty"`
	Phone               *string   `json:"phone,omitempty"`
	Role                string    `json:"role"`
	IsActive            bool      `json:"is_active"`
	ForcePasswordChange bool      `json:"force_password_change"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// UserListResponse is the data payload for GET /users.
type UserListResponse struct {
	Items    []UserResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// UserResponseFrom maps a service projection to API DTO.
func UserResponseFrom(p *user.UserPublic) UserResponse {
	if p == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:                  p.ID,
		FullName:            p.FullName,
		Username:            p.Username,
		Email:               p.Email,
		Phone:               p.Phone,
		Role:                string(p.Role),
		IsActive:            p.IsActive,
		ForcePasswordChange: p.ForcePasswordChange,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

// UserListResponseFrom maps a list result.
func UserListResponseFrom(r *user.ListResult) UserListResponse {
	if r == nil {
		return UserListResponse{}
	}
	items := make([]UserResponse, 0, len(r.Items))
	for i := range r.Items {
		items = append(items, UserResponseFrom(&r.Items[i]))
	}
	return UserListResponse{
		Items:    items,
		Total:    r.Total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}
}
