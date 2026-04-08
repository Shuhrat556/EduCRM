package repository

import (
	"context"
	"time"
)

// DashboardSnapshot holds aggregate counts for the dashboard (single query).
type DashboardSnapshot struct {
	TotalStudents int64
	TotalTeachers int64
	ActiveGroups  int64
	Debtors       int64
	TodayLessons  int64
	TodayPayments int64
}

// DashboardMonthlyRevenue is paid (non-free) volume recorded in a calendar month (UTC).
type DashboardMonthlyRevenue struct {
	TotalAmountMinor int64
	PaymentCount     int64
}

// DashboardStatsRepository provides read-optimized aggregates for staff dashboards.
type DashboardStatsRepository interface {
	// Snapshot loads headline counts as of ref (UTC calendar date and month).
	Snapshot(ctx context.Context, ref time.Time) (DashboardSnapshot, error)
	// MonthlyPaidRevenue sums amount_minor for paid_full/paid_partial, non-free rows
	// whose created_at falls in [monthStartUTC, monthStartUTC+1 month) in UTC.
	MonthlyPaidRevenue(ctx context.Context, monthStartUTC time.Time) (DashboardMonthlyRevenue, error)
}
