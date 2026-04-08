package handler

import (
	"net/http"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/usecase/aianalytics"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AIAnalyticsHandler exposes AI-powered analytics endpoints.
type AIAnalyticsHandler struct {
	svc *aianalytics.Service
}

// NewAIAnalyticsHandler constructs AIAnalyticsHandler.
func NewAIAnalyticsHandler(svc *aianalytics.Service) *AIAnalyticsHandler {
	return &AIAnalyticsHandler{svc: svc}
}

func bindAIFilters(c *gin.Context) (dto.AIAnalyticsFilters, error) {
	var body dto.AIAnalyticsFilters
	if c.Request.ContentLength <= 0 {
		return body, nil
	}
	if err := BindJSON(c, &body); err != nil {
		return body, err
	}
	return body, nil
}

func parseRFC3339Ptr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseMonthPtr(s *string) (*time.Time, error) {
	if s == nil || *s == "" {
		return nil, nil
	}
	t, err := time.ParseInLocation("2006-01", *s, time.UTC)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// DebtorsSummary godoc
// @Summary AI debtors summary
// @Tags ai-analytics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.AIAnalyticsFilters false "Optional month YYYY-MM"
// @Success 200 {object} response.Envelope{data=dto.AIAnalyticsResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/ai/analytics/debtors-summary [post]
func (h *AIAnalyticsHandler) DebtorsSummary(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	body, err := bindAIFilters(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	month, err := parseMonthPtr(body.Month)
	if err != nil {
		response.Error(c, apperror.Validation("month", "Use YYYY-MM for month"))
		return
	}
	out, err := h.svc.DebtorsSummary(c.Request.Context(), role, month)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AIAnalyticsResponse{Output: out.Output, Provider: out.Provider})
}

// LowAttendance godoc
// @Summary AI low attendance summary
// @Tags ai-analytics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.AIAnalyticsFilters false "Optional from/to RFC3339"
// @Success 200 {object} response.Envelope{data=dto.AIAnalyticsResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/ai/analytics/low-attendance [post]
func (h *AIAnalyticsHandler) LowAttendance(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	body, err := bindAIFilters(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	from, err := parseRFC3339Ptr(body.From)
	if err != nil {
		response.Error(c, apperror.Validation("from", "Use RFC3339 datetime"))
		return
	}
	to, err := parseRFC3339Ptr(body.To)
	if err != nil {
		response.Error(c, apperror.Validation("to", "Use RFC3339 datetime"))
		return
	}
	out, err := h.svc.LowAttendanceSummary(c.Request.Context(), role, from, to)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AIAnalyticsResponse{Output: out.Output, Provider: out.Provider})
}

// AdminDailySummary godoc
// @Summary AI admin daily summary
// @Tags ai-analytics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.AIAnalyticsFilters false "Optional as_of RFC3339"
// @Success 200 {object} response.Envelope{data=dto.AIAnalyticsResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/ai/analytics/admin-daily-summary [post]
func (h *AIAnalyticsHandler) AdminDailySummary(c *gin.Context) {
	role, err := RequireStaff(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	body, err := bindAIFilters(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	asOf, err := parseRFC3339Ptr(body.AsOf)
	if err != nil {
		response.Error(c, apperror.Validation("as_of", "Use RFC3339 datetime"))
		return
	}
	out, err := h.svc.AdminDailySummary(c.Request.Context(), role, asOf)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AIAnalyticsResponse{Output: out.Output, Provider: out.Provider})
}

// TeacherRecommendations godoc
// @Summary AI teacher recommendations
// @Description Teacher uses linked profile; staff must pass teacher_id (teachers table id).
// @Tags ai-analytics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.AIAnalyticsFilters false "teacher_id for staff"
// @Success 200 {object} response.Envelope{data=dto.AIAnalyticsResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/ai/analytics/teacher-recommendations [post]
func (h *AIAnalyticsHandler) TeacherRecommendations(c *gin.Context) {
	role, uid, err := RequireAITeacherRecommendationsActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	body, err := bindAIFilters(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.TeacherRecommendations(c.Request.Context(), role, uid, body.TeacherID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AIAnalyticsResponse{Output: out.Output, Provider: out.Provider})
}

// StudentWarnings godoc
// @Summary AI student warning suggestions
// @Description Student uses self; staff must pass student_id.
// @Tags ai-analytics
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.AIAnalyticsFilters false "student_id for staff"
// @Success 200 {object} response.Envelope{data=dto.AIAnalyticsResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/ai/analytics/student-warnings [post]
func (h *AIAnalyticsHandler) StudentWarnings(c *gin.Context) {
	role, uid, err := RequireAIStudentWarningsActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	body, err := bindAIFilters(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.StudentWarningSuggestions(c.Request.Context(), role, uid, body.StudentID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AIAnalyticsResponse{Output: out.Output, Provider: out.Provider})
}

