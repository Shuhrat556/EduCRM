package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AIAnalyticsContextRepository implements repository.AIAnalyticsContextRepository.
type AIAnalyticsContextRepository struct {
	db *gorm.DB
}

var _ repository.AIAnalyticsContextRepository = (*AIAnalyticsContextRepository)(nil)

// NewAIAnalyticsContextRepository constructs AIAnalyticsContextRepository.
func NewAIAnalyticsContextRepository(db *gorm.DB) *AIAnalyticsContextRepository {
	return &AIAnalyticsContextRepository{db: db}
}

// DebtorsSummaryData returns debtor student rows for a billed month (first day UTC).
func (r *AIAnalyticsContextRepository) DebtorsSummaryData(ctx context.Context, monthForUTC time.Time) (json.RawMessage, error) {
	m := time.Date(monthForUTC.Year(), monthForUTC.Month(), 1, 0, 0, 0, 0, time.UTC)
	const q = `
SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json)::text
FROM (
  SELECT x.student_id, x.open_payment_rows
  FROM (
    SELECT p.student_id::text AS student_id, COUNT(*)::int AS open_payment_rows
    FROM payments p
    WHERE p.deleted_at IS NULL
      AND p.status IN ('unpaid', 'overdue')
      AND p.month_for = ?::date
    GROUP BY p.student_id
  ) x
  ORDER BY x.open_payment_rows DESC, x.student_id
  LIMIT 100
) t
`
	var s string
	if err := r.db.WithContext(ctx).Raw(q, m).Scan(&s).Error; err != nil {
		return nil, err
	}
	return json.RawMessage(s), nil
}

// LowAttendanceData returns students with absence counts in a lesson_date range (inclusive).
func (r *AIAnalyticsContextRepository) LowAttendanceData(ctx context.Context, fromUTC, toUTC time.Time) (json.RawMessage, error) {
	from := time.Date(fromUTC.Year(), fromUTC.Month(), fromUTC.Day(), 0, 0, 0, 0, time.UTC)
	to := time.Date(toUTC.Year(), toUTC.Month(), toUTC.Day(), 0, 0, 0, 0, time.UTC)
	const q = `
SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json)::text
FROM (
  SELECT a.student_id::text AS student_id,
    COUNT(*) FILTER (WHERE a.status = 'absent')::int AS absent_count,
    COUNT(*)::int AS marked_lessons
  FROM attendances a
  WHERE a.lesson_date >= ?::date AND a.lesson_date <= ?::date
  GROUP BY a.student_id
  HAVING COUNT(*) FILTER (WHERE a.status = 'absent') >= 1
  ORDER BY COUNT(*) FILTER (WHERE a.status = 'absent') DESC, a.student_id
  LIMIT 50
) t
`
	var out string
	if err := r.db.WithContext(ctx).Raw(q, from, to).Scan(&out).Error; err != nil {
		return nil, err
	}
	return json.RawMessage(out), nil
}

// TeacherGroupsData returns active groups and enrollment counts for a teacher profile id.
func (r *AIAnalyticsContextRepository) TeacherGroupsData(ctx context.Context, teacherID uuid.UUID) (json.RawMessage, error) {
	const q = `
SELECT COALESCE(json_agg(row_to_json(t)), '[]'::json)::text
FROM (
  SELECT g.id::text AS group_id, g.name AS group_name, g.monthly_fee_minor,
    (SELECT COUNT(*) FROM student_group_memberships sgm WHERE sgm.group_id = g.id)::int AS enrolled_students
  FROM groups g
  WHERE g.teacher_id = ? AND g.status = 'active'
  ORDER BY g.name
) t
`
	var out string
	if err := r.db.WithContext(ctx).Raw(q, teacherID).Scan(&out).Error; err != nil {
		return nil, err
	}
	return json.RawMessage(out), nil
}

// StudentRiskData returns payment and absence signals for one student.
func (r *AIAnalyticsContextRepository) StudentRiskData(ctx context.Context, studentID uuid.UUID, paymentMonthUTC, attendanceFrom, attendanceTo time.Time) (json.RawMessage, error) {
	pm := time.Date(paymentMonthUTC.Year(), paymentMonthUTC.Month(), 1, 0, 0, 0, 0, time.UTC)
	af := time.Date(attendanceFrom.Year(), attendanceFrom.Month(), attendanceFrom.Day(), 0, 0, 0, 0, time.UTC)
	at := time.Date(attendanceTo.Year(), attendanceTo.Month(), attendanceTo.Day(), 0, 0, 0, 0, time.UTC)
	const q = `
SELECT json_build_object(
  'student_id', ?::text,
  'payment_month', ?::text,
  'unpaid_or_overdue_rows_current_month', (
    SELECT COUNT(*)::int FROM payments p
    WHERE p.student_id = ? AND p.deleted_at IS NULL
      AND p.status IN ('unpaid','overdue')
      AND p.month_for = ?::date
  ),
  'absences_in_range', (
    SELECT COUNT(*)::int FROM attendances a
    WHERE a.student_id = ? AND a.status = 'absent'
      AND a.lesson_date >= ?::date AND a.lesson_date <= ?::date
  )
)::text
`
	var out string
	if err := r.db.WithContext(ctx).Raw(q,
		studentID, pm.Format("2006-01-02"),
		studentID, pm,
		studentID, af, at,
	).Scan(&out).Error; err != nil {
		return nil, err
	}
	return json.RawMessage(out), nil
}
