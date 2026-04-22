package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// CreateGradeRequest is the body for POST /grades.
type CreateGradeRequest struct {
	StudentID  uuid.UUID  `json:"student_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
	GroupID    uuid.UUID  `json:"group_id" example:"6ba7b810-9dad-11d1-80b4-00c04fd430c8" binding:"required"`
	SubjectID  *uuid.UUID `json:"subject_id"` // optional; defaults to group's subject
	GradeType  string     `json:"grade_type" example:"teacher_evaluation" binding:"required,oneof=teacher_evaluation student_evaluation"`
	GradeValue float64   `json:"grade_value" example:"4.5" binding:"required"`
	Comment    *string   `json:"comment" binding:"omitempty,max=4000"`
	WeekOf     *string   `json:"week_of" binding:"omitempty"` // YYYY-MM-DD; week bucket is Monday of that week in UTC
	GradedAt   *string   `json:"graded_at" binding:"omitempty"` // RFC3339; default now
}

// UpdateGradeRequest is the body for PATCH /grades/:id.
type UpdateGradeRequest struct {
	GradeValue *float64 `json:"grade_value" example:"4.5"`
	Comment    *string  `json:"comment" example:"Great progress" binding:"omitempty,max=4000"`
	GradedAt   *string  `json:"graded_at" example:"2025-04-08T12:00:00Z" binding:"omitempty"` // RFC3339
}

// GradeResponse is one grade in API responses.
type GradeResponse struct {
	ID            uuid.UUID `json:"id"`
	StudentID     uuid.UUID `json:"student_id"`
	TeacherID     uuid.UUID `json:"teacher_id"`
	GroupID       uuid.UUID `json:"group_id"`
	SubjectID     uuid.UUID `json:"subject_id"`
	WeekStartDate string    `json:"week_start_date"`
	GradeType     string    `json:"grade_type"`
	GradeValue    float64   `json:"grade_value"`
	Comment       *string   `json:"comment,omitempty"`
	GradedAt      time.Time `json:"graded_at"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// GradeListResponse wraps a list of grades.
type GradeListResponse struct {
	Items []GradeResponse `json:"items"`
}

// GradeResponseFrom maps domain to DTO.
func GradeResponseFrom(g *domain.Grade) GradeResponse {
	if g == nil {
		return GradeResponse{}
	}
	return GradeResponse{
		ID:            g.ID,
		StudentID:     g.StudentID,
		TeacherID:     g.TeacherID,
		GroupID:       g.GroupID,
		SubjectID:     g.SubjectID,
		WeekStartDate: g.WeekStartDate.UTC().Format("2006-01-02"),
		GradeType:     string(g.GradeType),
		GradeValue:    g.GradeValue,
		Comment:       g.Comment,
		GradedAt:      g.GradedAt,
		CreatedAt:     g.CreatedAt,
		UpdatedAt:     g.UpdatedAt,
	}
}

// GradeListResponseFrom maps a slice.
func GradeListResponseFrom(rows []domain.Grade) GradeListResponse {
	items := make([]GradeResponse, 0, len(rows))
	for i := range rows {
		items = append(items, GradeResponseFrom(&rows[i]))
	}
	return GradeListResponse{Items: items}
}
