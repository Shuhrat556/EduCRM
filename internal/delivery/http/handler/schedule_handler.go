package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/usecase/schedule"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ScheduleHandler exposes schedule HTTP endpoints.
type ScheduleHandler struct {
	svc *schedule.Service
}

// NewScheduleHandler constructs ScheduleHandler.
func NewScheduleHandler(svc *schedule.Service) *ScheduleHandler {
	return &ScheduleHandler{svc: svc}
}

// Create godoc
// @Summary Create schedule entry
// @Description Creates a weekly recurring slot. Weekday is 0=Sunday … 6=Saturday. Times are HH:MM (24h); end may be 24:00. Room and teacher cannot double-book on the same weekday and overlapping time.
// @Tags schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateScheduleRequest true "Schedule"
// @Success 201 {object} response.Envelope{data=dto.ScheduleResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/schedules [post]
func (h *ScheduleHandler) Create(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateScheduleRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	wd, err := domain.ParseWeekday(req.Weekday)
	if err != nil {
		response.Error(c, apperror.Validation("weekday", "Weekday must be 0–6 (Sunday–Saturday)"))
		return
	}
	start, err := dto.ParseClockToMinutes(req.StartTime)
	if err != nil {
		response.Error(c, apperror.Validation("start_time", "Use 24-hour time HH:MM"))
		return
	}
	if start >= domain.MinutesPerDay {
		response.Error(c, apperror.Validation("start_time", "Start time cannot be 24:00"))
		return
	}
	end, err := dto.ParseClockToMinutes(req.EndTime)
	if err != nil {
		response.Error(c, apperror.Validation("end_time", "Use HH:MM or 24:00 for end of day"))
		return
	}
	out, err := h.svc.Create(c.Request.Context(), schedule.CreateInput{
		GroupID:      req.GroupID,
		TeacherID:    req.TeacherID,
		RoomID:       req.RoomID,
		Weekday:      wd,
		StartMinutes: start,
		EndMinutes:   end,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.ScheduleResponseFrom(out))
}

// List godoc
// @Summary List schedules
// @Description Exactly one query param: group_id, teacher_id, or room_id.
// @Tags schedules
// @Security BearerAuth
// @Produce json
// @Param group_id query string false "Group UUID"
// @Param teacher_id query string false "Teacher UUID"
// @Param room_id query string false "Room UUID"
// @Success 200 {object} response.Envelope{data=dto.ScheduleListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/schedules [get]
func (h *ScheduleHandler) List(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var q dto.ScheduleListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	n := 0
	if q.GroupID != "" {
		n++
	}
	if q.TeacherID != "" {
		n++
	}
	if q.RoomID != "" {
		n++
	}
	if n != 1 {
		response.Error(c, apperror.Validation("filter", "Provide exactly one of: group_id, teacher_id, or room_id"))
		return
	}
	var f schedule.ListFilter
	switch {
	case q.GroupID != "":
		id := uuid.MustParse(q.GroupID)
		f.GroupID = &id
	case q.TeacherID != "":
		id := uuid.MustParse(q.TeacherID)
		f.TeacherID = &id
	default:
		id := uuid.MustParse(q.RoomID)
		f.RoomID = &id
	}
	rows, err := h.svc.List(c.Request.Context(), f)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.ScheduleListResponseFrom(rows))
}

// GetByID godoc
// @Summary Get schedule by ID
// @Tags schedules
// @Security BearerAuth
// @Produce json
// @Param id path string true "Schedule ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.ScheduleResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/schedules/{id} [get]
func (h *ScheduleHandler) GetByID(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.ScheduleResponseFrom(out))
}

// Update godoc
// @Summary Update schedule entry
// @Tags schedules
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Schedule ID (UUID)"
// @Param body body dto.UpdateScheduleRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.ScheduleResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/schedules/{id} [patch]
func (h *ScheduleHandler) Update(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateScheduleRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := schedule.UpdateInput{
		GroupID:   req.GroupID,
		TeacherID: req.TeacherID,
		RoomID:    req.RoomID,
	}
	if req.Weekday != nil {
		wd, err := domain.ParseWeekday(*req.Weekday)
		if err != nil {
			response.Error(c, apperror.Validation("weekday", "Weekday must be 0–6 (Sunday–Saturday)"))
			return
		}
		in.Weekday = &wd
	}
	if req.StartTime != nil {
		start, err := dto.ParseClockToMinutes(*req.StartTime)
		if err != nil {
			response.Error(c, apperror.Validation("start_time", "Use 24-hour time HH:MM"))
			return
		}
		if start >= domain.MinutesPerDay {
			response.Error(c, apperror.Validation("start_time", "Start time cannot be 24:00"))
			return
		}
		in.StartMinutes = &start
	}
	if req.EndTime != nil {
		end, err := dto.ParseClockToMinutes(*req.EndTime)
		if err != nil {
			response.Error(c, apperror.Validation("end_time", "Use HH:MM or 24:00 for end of day"))
			return
		}
		in.EndMinutes = &end
	}
	out, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.ScheduleResponseFrom(out))
}

// Delete godoc
// @Summary Delete schedule entry
// @Tags schedules
// @Security BearerAuth
// @Produce json
// @Param id path string true "Schedule ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/schedules/{id} [delete]
func (h *ScheduleHandler) Delete(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Schedule entry deleted successfully"})
}
