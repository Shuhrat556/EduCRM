package handler

import (
	"net/http"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/usecase/dashboard"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// DashboardHandler exposes staff dashboard endpoints.
type DashboardHandler struct {
	svc *dashboard.Service
}

// NewDashboardHandler constructs DashboardHandler.
func NewDashboardHandler(svc *dashboard.Service) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

// Summary godoc
// @Summary Dashboard statistics
// @Description Headline counts (UTC date for "today"). Debtors = distinct students with unpaid/overdue rows for the current billing month (month_for). Today lessons = schedule slots on today's weekday for active groups whose date range includes today. Today payments = rows with payment_date = today (UTC) or, if null, created_at date = today. Revenue = sum of amount_minor for paid_full/paid_partial, non-free payments created in the selected UTC month (default current).
// @Tags dashboard
// @Security BearerAuth
// @Produce json
// @Param as_of query string false "Reference instant RFC3339 (default: now UTC)"
// @Param year_month query string false "Revenue month YYYY-MM (UTC, default: current month)"
// @Success 200 {object} response.Envelope{data=dto.DashboardSummaryResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/dashboard/summary [get]
func (h *DashboardHandler) Summary(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	ref := time.Now().UTC()
	if s := c.Query("as_of"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			response.Error(c, apperror.Validation("as_of", "Use RFC3339 instant (e.g. 2025-04-08T12:00:00Z)"))
			return
		}
		ref = t.UTC()
	}
	var revMonth *time.Time
	if ym := c.Query("year_month"); ym != "" {
		m, err := dashboard.ParseYearMonth(ym)
		if err != nil {
			response.Error(c, apperror.Validation("year_month", "Use calendar month YYYY-MM"))
			return
		}
		revMonth = &m
	}
	out, err := h.svc.GetSummary(c.Request.Context(), role, ref, revMonth)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.DashboardSummaryFrom(out))
}
