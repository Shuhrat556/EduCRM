package dto

import (
	"github.com/educrm/educrm-backend/internal/usecase/dashboard"
)

// DashboardCountsResponse is the count block for GET /dashboard/summary.
type DashboardCountsResponse struct {
	TotalStudents int64 `json:"total_students"`
	TotalTeachers int64 `json:"total_teachers"`
	ActiveGroups  int64 `json:"active_groups"`
	Debtors       int64 `json:"debtors"`
	TodayLessons  int64 `json:"today_lessons"`
	TodayPayments int64 `json:"today_payments"`
}

// DashboardMonthlyRevenueResponse is paid revenue for one UTC calendar month.
type DashboardMonthlyRevenueResponse struct {
	YearMonth        string `json:"year_month"`
	TotalAmountMinor int64  `json:"total_amount_minor"`
	PaymentCount     int64  `json:"payment_count"`
}

// DashboardSummaryResponse is the full dashboard payload.
type DashboardSummaryResponse struct {
	Counts  DashboardCountsResponse         `json:"counts"`
	Revenue DashboardMonthlyRevenueResponse `json:"revenue"`
}

// DashboardSummaryFrom maps use case output to API DTOs.
func DashboardSummaryFrom(s *dashboard.Summary) DashboardSummaryResponse {
	if s == nil {
		return DashboardSummaryResponse{}
	}
	return DashboardSummaryResponse{
		Counts: DashboardCountsResponse{
			TotalStudents: s.Counts.TotalStudents,
			TotalTeachers: s.Counts.TotalTeachers,
			ActiveGroups:  s.Counts.ActiveGroups,
			Debtors:       s.Counts.Debtors,
			TodayLessons:  s.Counts.TodayLessons,
			TodayPayments: s.Counts.TodayPayments,
		},
		Revenue: DashboardMonthlyRevenueResponse{
			YearMonth:        s.Revenue.YearMonth,
			TotalAmountMinor: s.Revenue.TotalAmountMinor,
			PaymentCount:     s.Revenue.PaymentCount,
		},
	}
}
