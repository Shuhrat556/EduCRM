package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/usecase/attendance"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AttendanceHandler exposes attendance HTTP endpoints.
type AttendanceHandler struct {
	svc *attendance.Service
}

// NewAttendanceHandler constructs AttendanceHandler.
func NewAttendanceHandler(svc *attendance.Service) *AttendanceHandler {
	return &AttendanceHandler{svc: svc}
}

// Create godoc
// @Summary Mark attendance
// @Description Records attendance for a student on a lesson date. Admins may mark for any group; teachers only for groups they teach (user_teacher_links). marked_by_teacher_id is set to the group's assigned teacher.
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateAttendanceRequest true "Attendance"
// @Success 201 {object} response.Envelope{data=dto.AttendanceResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/attendance [post]
func (h *AttendanceHandler) Create(c *gin.Context) {
	role, actorID, err := RequireAttendanceActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateAttendanceRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	ld, err := dto.ParseISODate(req.LessonDate)
	if err != nil {
		response.Error(c, apperror.Validation("lesson_date", "Use date format YYYY-MM-DD"))
		return
	}
	st, err := domain.ParseAttendanceStatus(req.Status)
	if err != nil {
		response.Error(c, apperror.Validation("status", "status must be present, absent, or late"))
		return
	}
	out, err := h.svc.Create(c.Request.Context(), role, actorID, attendance.CreateInput{
		StudentID:  req.StudentID,
		GroupID:    req.GroupID,
		LessonDate: ld,
		Status:     st,
		Comment:    req.Comment,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.AttendanceResponseFrom(out))
}

// List godoc
// @Summary List attendance
// @Description Exactly one mode: student_id, group_id, or from+to (YYYY-MM-DD) date range.
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param student_id query string false "Student user UUID"
// @Param group_id query string false "Group UUID"
// @Param from query string false "Range start (YYYY-MM-DD)"
// @Param to query string false "Range end (YYYY-MM-DD)"
// @Success 200 {object} response.Envelope{data=dto.AttendanceListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/attendance [get]
func (h *AttendanceHandler) List(c *gin.Context) {
	role, actorID, err := RequireAttendanceActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.AttendanceListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	var f attendance.ListFilter
	if q.StudentID != "" {
		id := uuid.MustParse(q.StudentID)
		f.StudentID = &id
	}
	if q.GroupID != "" {
		id := uuid.MustParse(q.GroupID)
		f.GroupID = &id
	}
	if q.From != "" || q.To != "" {
		if q.From == "" || q.To == "" {
			response.Error(c, apperror.Validation("filter", "For a date range, both from and to are required (YYYY-MM-DD)"))
			return
		}
		fromT, err := dto.ParseISODate(q.From)
		if err != nil {
			response.Error(c, apperror.Validation("from", "Use date format YYYY-MM-DD"))
			return
		}
		toT, err := dto.ParseISODate(q.To)
		if err != nil {
			response.Error(c, apperror.Validation("to", "Use date format YYYY-MM-DD"))
			return
		}
		f.From = &fromT
		f.To = &toT
	}
	rows, err := h.svc.List(c.Request.Context(), role, actorID, f)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AttendanceListResponseFrom(rows))
}

// GetByID godoc
// @Summary Get attendance by ID
// @Tags attendance
// @Security BearerAuth
// @Produce json
// @Param id path string true "Attendance ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.AttendanceResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/attendance/{id} [get]
func (h *AttendanceHandler) GetByID(c *gin.Context) {
	role, actorID, err := RequireAttendanceActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.GetByID(c.Request.Context(), role, actorID, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AttendanceResponseFrom(out))
}

// Update godoc
// @Summary Update attendance
// @Tags attendance
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Attendance ID (UUID)"
// @Param body body dto.UpdateAttendanceRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.AttendanceResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/attendance/{id} [patch]
func (h *AttendanceHandler) Update(c *gin.Context) {
	role, actorID, err := RequireAttendanceActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateAttendanceRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := attendance.UpdateInput{Comment: req.Comment}
	if req.Status != nil {
		st, err := domain.ParseAttendanceStatus(*req.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "status must be present, absent, or late"))
			return
		}
		in.Status = &st
	}
	if in.Status == nil && req.Comment == nil {
		response.Error(c, apperror.Validation("body", "Provide at least one of status or comment"))
		return
	}
	out, err := h.svc.Update(c.Request.Context(), role, actorID, id, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.AttendanceResponseFrom(out))
}
