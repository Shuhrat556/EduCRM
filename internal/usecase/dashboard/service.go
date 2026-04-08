package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/educrm/educrm-backend/internal/repository"
)

// Service builds dashboard aggregates for staff.
type Service struct {
	stats repository.DashboardStatsRepository
}

// NewService constructs the dashboard service.
func NewService(stats repository.DashboardStatsRepository) *Service {
	return &Service{stats: stats}
}

// Summary is the full dashboard payload (logic lives here; handler stays thin).
type Summary struct {
	Counts  Counts
	Revenue MonthlyRevenue
}

// Counts mirrors snapshot metrics.
type Counts struct {
	TotalStudents int64
	TotalTeachers int64
	ActiveGroups  int64
	Debtors       int64
	TodayLessons  int64
	TodayPayments int64
}

// MonthlyRevenue is paid (non-free) totals for a calendar month (UTC).
type MonthlyRevenue struct {
	YearMonth        string // YYYY-MM
	TotalAmountMinor int64
	PaymentCount     int64
}

// GetSummary returns snapshot counts and monthly revenue.
// If revenueMonth is nil, the current UTC month is used.
func (s *Service) GetSummary(ctx context.Context, actorRole domain.Role, ref time.Time, revenueMonth *time.Time) (*Summary, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	snap, err := s.stats.Snapshot(ctx, ref)
	if err != nil {
		return nil, apperror.Internal("dashboard snapshot").Wrap(err)
	}
	monthStart := ref.UTC()
	if revenueMonth != nil {
		monthStart = revenueMonth.UTC()
	}
	monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, time.UTC)
	rev, err := s.stats.MonthlyPaidRevenue(ctx, monthStart)
	if err != nil {
		return nil, apperror.Internal("dashboard revenue").Wrap(err)
	}
	return &Summary{
		Counts: Counts{
			TotalStudents: snap.TotalStudents,
			TotalTeachers: snap.TotalTeachers,
			ActiveGroups:  snap.ActiveGroups,
			Debtors:       snap.Debtors,
			TodayLessons:  snap.TodayLessons,
			TodayPayments: snap.TodayPayments,
		},
		Revenue: MonthlyRevenue{
			YearMonth:        monthStart.Format("2006-01"),
			TotalAmountMinor: rev.TotalAmountMinor,
			PaymentCount:     rev.PaymentCount,
		},
	}, nil
}

// ParseYearMonth parses "2006-01" into the first instant of that month UTC, or returns a validation error.
func ParseYearMonth(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, fmt.Errorf("empty year_month")
	}
	t, err := time.ParseInLocation("2006-01", s, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}
