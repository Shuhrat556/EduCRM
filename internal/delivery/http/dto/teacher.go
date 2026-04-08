package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	teachersvc "github.com/educrm/educrm-backend/internal/usecase/teacher"
	"github.com/google/uuid"
)

// CreateTeacherRequest is the body for POST /teachers.
type CreateTeacherRequest struct {
	FullName          string      `json:"full_name" binding:"required,min=2,max=255"`
	Phone             *string     `json:"phone" binding:"omitempty,min=5,max=32"`
	Email             *string     `json:"email" binding:"omitempty,email,max=255"`
	Specialization    *string     `json:"specialization" binding:"omitempty,max=255"`
	PhotoURL          *string     `json:"photo_url" binding:"omitempty,max=2048"`
	PhotoStorageKey   *string     `json:"photo_storage_key" binding:"omitempty,max=512"`
	PhotoContentType  *string     `json:"photo_content_type" binding:"omitempty,max=128"`
	PhotoOriginalName *string     `json:"photo_original_name" binding:"omitempty,max=255"`
	Status            string      `json:"status" binding:"omitempty,oneof=active inactive"`
}

// UpdateTeacherRequest is the body for PATCH /teachers/:id.
type UpdateTeacherRequest struct {
	FullName          *string      `json:"full_name" binding:"omitempty,min=2,max=255"`
	Phone             *string      `json:"phone" binding:"omitempty,min=5,max=32"`
	Email             *string      `json:"email" binding:"omitempty,email,max=255"`
	Specialization    *string      `json:"specialization" binding:"omitempty,max=255"`
	PhotoURL          *string      `json:"photo_url" binding:"omitempty,max=2048"`
	PhotoStorageKey   *string      `json:"photo_storage_key" binding:"omitempty,max=512"`
	PhotoContentType  *string      `json:"photo_content_type" binding:"omitempty,max=128"`
	PhotoOriginalName *string      `json:"photo_original_name" binding:"omitempty,max=255"`
	Status            *string      `json:"status" binding:"omitempty,oneof=active inactive"`
}

// PatchTeacherPhotoRequest updates stored photo metadata (e.g. after S3 upload).
type PatchTeacherPhotoRequest struct {
	PhotoURL          *string `json:"photo_url" binding:"omitempty,max=2048"`
	PhotoStorageKey   *string `json:"photo_storage_key" binding:"omitempty,max=512"`
	PhotoContentType  *string `json:"photo_content_type" binding:"omitempty,max=128"`
	PhotoOriginalName *string `json:"photo_original_name" binding:"omitempty,max=255"`
}

// GroupSummary is a minimal group in teacher responses.
type GroupSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// TeacherSummary is used in list responses.
type TeacherSummary struct {
	ID                uuid.UUID `json:"id"`
	FullName          string    `json:"full_name"`
	Phone             *string   `json:"phone,omitempty"`
	Email             *string   `json:"email,omitempty"`
	Specialization    *string   `json:"specialization,omitempty"`
	PhotoURL          *string   `json:"photo_url,omitempty"`
	PhotoStorageKey   *string   `json:"photo_storage_key,omitempty"`
	PhotoContentType  *string   `json:"photo_content_type,omitempty"`
	PhotoOriginalName *string   `json:"photo_original_name,omitempty"`
	Status            string    `json:"status"`
	GroupCount        int       `json:"group_count"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// TeacherDetailResponse is returned for create/get/update/photo.
type TeacherDetailResponse struct {
	ID                uuid.UUID      `json:"id"`
	FullName          string         `json:"full_name"`
	Phone             *string        `json:"phone,omitempty"`
	Email             *string        `json:"email,omitempty"`
	Specialization    *string        `json:"specialization,omitempty"`
	PhotoURL          *string        `json:"photo_url,omitempty"`
	PhotoStorageKey   *string        `json:"photo_storage_key,omitempty"`
	PhotoContentType  *string        `json:"photo_content_type,omitempty"`
	PhotoOriginalName *string        `json:"photo_original_name,omitempty"`
	Status            string         `json:"status"`
	Groups            []GroupSummary `json:"groups"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
}

// TeacherListResponse is the data envelope for GET /teachers.
type TeacherListResponse struct {
	Items    []TeacherSummary `json:"items"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	PageSize int              `json:"page_size"`
}

// TeacherDetailFrom maps service detail to API DTO.
func TeacherDetailFrom(d *teachersvc.Detail) TeacherDetailResponse {
	if d == nil || d.Teacher == nil {
		return TeacherDetailResponse{}
	}
	t := d.Teacher
	groups := make([]GroupSummary, 0, len(d.Groups))
	for _, g := range d.Groups {
		groups = append(groups, GroupSummary{ID: g.ID, Name: g.Name})
	}
	return TeacherDetailResponse{
		ID:                t.ID,
		FullName:          t.FullName,
		Phone:             t.Phone,
		Email:             t.Email,
		Specialization:    t.Specialization,
		PhotoURL:          t.PhotoURL,
		PhotoStorageKey:   t.PhotoStorageKey,
		PhotoContentType:  t.PhotoContentType,
		PhotoOriginalName: t.PhotoOriginalName,
		Status:            string(t.Status),
		Groups:            groups,
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
	}
}

// TeacherSummaryFrom maps domain teacher + count.
func TeacherSummaryFrom(t *domain.Teacher, groupCount int) TeacherSummary {
	if t == nil {
		return TeacherSummary{}
	}
	return TeacherSummary{
		ID:                t.ID,
		FullName:          t.FullName,
		Phone:             t.Phone,
		Email:             t.Email,
		Specialization:    t.Specialization,
		PhotoURL:          t.PhotoURL,
		PhotoStorageKey:   t.PhotoStorageKey,
		PhotoContentType:  t.PhotoContentType,
		PhotoOriginalName: t.PhotoOriginalName,
		Status:            string(t.Status),
		GroupCount:        groupCount,
		CreatedAt:         t.CreatedAt,
		UpdatedAt:         t.UpdatedAt,
	}
}

// TeacherListResponseFrom maps list result.
func TeacherListResponseFrom(r *teachersvc.ListResult) TeacherListResponse {
	if r == nil {
		return TeacherListResponse{}
	}
	items := make([]TeacherSummary, 0, len(r.Items))
	for i := range r.Items {
		items = append(items, TeacherSummaryFrom(&r.Items[i].Teacher, r.Items[i].GroupCount))
	}
	return TeacherListResponse{
		Items:    items,
		Total:    r.Total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}
}
