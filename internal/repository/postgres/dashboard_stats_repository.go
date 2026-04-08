package postgres

import (
	"context"
	"time"

	"github.com/educrm/educrm-backend/internal/repository"
	"gorm.io/gorm"
)

// DashboardStatsRepository implements repository.DashboardStatsRepository.
type DashboardStatsRepository struct {
	db *gorm.DB
}

var _ repository.DashboardStatsRepository = (*DashboardStatsRepository)(nil)

// NewDashboardStatsRepository constructs DashboardStatsRepository.
func NewDashboardStatsRepository(db *gorm.DB) *DashboardStatsRepository {
	return &DashboardStatsRepository{db: db}
}

// Snapshot runs one round-trip with scalar subqueries.
func (r *DashboardStatsRepository) Snapshot(ctx context.Context, ref time.Time) (repository.DashboardSnapshot, error) {
	ref = ref.UTC()
	today := time.Date(ref.Year(), ref.Month(), ref.Day(), 0, 0, 0, 0, time.UTC)
	monthFor := time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, time.UTC)
	weekday := int16(ref.Weekday())

	const q = `
SELECT
  (SELECT COUNT(*) FROM users WHERE role = 'student' AND is_active = true) AS total_students,
  (SELECT COUNT(*) FROM teachers WHERE status = 'active') AS total_teachers,
  (SELECT COUNT(*) FROM groups WHERE status = 'active') AS active_groups,
  (SELECT COUNT(DISTINCT p.student_id) FROM payments p
     WHERE p.deleted_at IS NULL
       AND p.status IN ('unpaid', 'overdue')
       AND p.month_for = ?::date
  ) AS debtors,
  (SELECT COUNT(*) FROM schedules s
     INNER JOIN groups g ON g.id = s.group_id
     WHERE g.status = 'active'
       AND ?::date BETWEEN g.start_date AND g.end_date
       AND s.weekday = ?
  ) AS today_lessons,
  (SELECT COUNT(*) FROM payments p
     WHERE p.deleted_at IS NULL
       AND (
         (p.payment_date IS NOT NULL AND p.payment_date = ?::date)
         OR (p.payment_date IS NULL AND (p.created_at AT TIME ZONE 'UTC')::date = ?::date)
       )
  ) AS today_payments
`
	var row struct {
		TotalStudents int64 `gorm:"column:total_students"`
		TotalTeachers int64 `gorm:"column:total_teachers"`
		ActiveGroups  int64 `gorm:"column:active_groups"`
		Debtors       int64 `gorm:"column:debtors"`
		TodayLessons  int64 `gorm:"column:today_lessons"`
		TodayPayments int64 `gorm:"column:today_payments"`
	}
	err := r.db.WithContext(ctx).Raw(q,
		monthFor,
		today,
		weekday,
		today,
		today,
	).Scan(&row).Error
	if err != nil {
		return repository.DashboardSnapshot{}, err
	}
	return repository.DashboardSnapshot{
		TotalStudents: row.TotalStudents,
		TotalTeachers: row.TotalTeachers,
		ActiveGroups:  row.ActiveGroups,
		Debtors:       row.Debtors,
		TodayLessons:  row.TodayLessons,
		TodayPayments: row.TodayPayments,
	}, nil
}

// MonthlyPaidRevenue aggregates paid rows in the UTC month of monthStartUTC.
func (r *DashboardStatsRepository) MonthlyPaidRevenue(ctx context.Context, monthStartUTC time.Time) (repository.DashboardMonthlyRevenue, error) {
	start := time.Date(monthStartUTC.Year(), monthStartUTC.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	const q = `
SELECT
  COALESCE(SUM(p.amount_minor), 0)::bigint AS total_amount_minor,
  COUNT(*)::bigint AS payment_count
FROM payments p
WHERE p.deleted_at IS NULL
  AND p.status IN ('paid_full', 'paid_partial')
  AND p.is_free = false
  AND p.created_at >= ? AND p.created_at < ?
`
	var row struct {
		TotalAmountMinor int64 `gorm:"column:total_amount_minor"`
		PaymentCount     int64 `gorm:"column:payment_count"`
	}
	err := r.db.WithContext(ctx).Raw(q, start, end).Scan(&row).Error
	if err != nil {
		return repository.DashboardMonthlyRevenue{}, err
	}
	return repository.DashboardMonthlyRevenue{
		TotalAmountMinor: row.TotalAmountMinor,
		PaymentCount:     row.PaymentCount,
	}, nil
}
