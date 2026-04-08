package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AIAnalyticsContextRepository loads structured JSON context for AI analytics prompts.
type AIAnalyticsContextRepository interface {
	DebtorsSummaryData(ctx context.Context, monthForUTC time.Time) (json.RawMessage, error)
	LowAttendanceData(ctx context.Context, fromUTC, toUTC time.Time) (json.RawMessage, error)
	TeacherGroupsData(ctx context.Context, teacherID uuid.UUID) (json.RawMessage, error)
	StudentRiskData(ctx context.Context, studentID uuid.UUID, paymentMonthUTC, attendanceFrom, attendanceTo time.Time) (json.RawMessage, error)
}
