package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// CreateAttendanceRequest is the body for POST /attendance.
type CreateAttendanceRequest struct {
	StudentID  uuid.UUID  `json:"student_id" binding:"required"`
	GroupID    uuid.UUID  `json:"group_id" binding:"required"`
	SubjectID  *uuid.UUID `json:"subject_id"` // optional; defaults to group's subject
	LessonDate string     `json:"lesson_date" binding:"required"` // YYYY-MM-DD
	Status     string     `json:"status" binding:"required,oneof=present absent late has nest"`
	Comment    *string    `json:"comment" binding:"omitempty,max=4000"`
}

// UpdateAttendanceRequest is the body for PATCH /attendance/:id.
type UpdateAttendanceRequest struct {
	Status  *string `json:"status" binding:"omitempty,oneof=present absent late has nest"`
	Comment *string `json:"comment" binding:"omitempty,max=4000"`
}

// AttendanceResponse is the API shape for one attendance row.
type AttendanceResponse struct {
	ID                  uuid.UUID `json:"id"`
	StudentID           uuid.UUID `json:"student_id"`
	GroupID             uuid.UUID `json:"group_id"`
	SubjectID           uuid.UUID `json:"subject_id"`
	LessonDate          string    `json:"lesson_date"`
	Status              string    `json:"status"`
	Comment             *string   `json:"comment,omitempty"`
	MarkedByTeacherID   uuid.UUID `json:"marked_by_teacher_id"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// AttendanceListResponse wraps list results.
type AttendanceListResponse struct {
	Items []AttendanceResponse `json:"items"`
}

// AttendanceResponseFrom maps domain to DTO.
func AttendanceResponseFrom(a *domain.Attendance) AttendanceResponse {
	if a == nil {
		return AttendanceResponse{}
	}
	return AttendanceResponse{
		ID:                a.ID,
		StudentID:         a.StudentID,
		GroupID:           a.GroupID,
		SubjectID:         a.SubjectID,
		LessonDate:        a.LessonDate.UTC().Format("2006-01-02"),
		Status:            string(a.Status),
		Comment:           a.Comment,
		MarkedByTeacherID: a.MarkedByTeacherID,
		CreatedAt:         a.CreatedAt,
		UpdatedAt:         a.UpdatedAt,
	}
}

// AttendanceListResponseFrom maps a slice.
func AttendanceListResponseFrom(rows []domain.Attendance) AttendanceListResponse {
	items := make([]AttendanceResponse, 0, len(rows))
	for i := range rows {
		items = append(items, AttendanceResponseFrom(&rows[i]))
	}
	return AttendanceListResponse{Items: items}
}
