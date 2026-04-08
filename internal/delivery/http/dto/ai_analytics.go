package dto

import "github.com/google/uuid"

// AIAnalyticsFilters optional JSON body for POST /ai/analytics/* endpoints.
type AIAnalyticsFilters struct {
	AsOf      *string    `json:"as_of" example:"2025-04-08T09:00:00Z"`                       // RFC3339 — admin daily, optional debtors anchor
	Month     *string    `json:"month" example:"2025-04"`                                  // YYYY-MM — debtors billed month (default: current UTC month)
	From      *string    `json:"from" example:"2025-03-01T00:00:00Z"`                      // RFC3339 — attendance window start
	To        *string    `json:"to" example:"2025-04-01T00:00:00Z"`                          // RFC3339 — attendance window end
	TeacherID *uuid.UUID `json:"teacher_id" example:"6ba7b810-9dad-11d1-80b4-00c04fd430c8"` // staff: teachers.id for recommendations
	StudentID *uuid.UUID `json:"student_id" example:"550e8400-e29b-41d4-a716-446655440000"` // staff: student user id for warnings
}

// AIAnalyticsResponse is the model output envelope.
type AIAnalyticsResponse struct {
	Output   string `json:"output"`
	Provider string `json:"provider"`
}
