package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	groupsvc "github.com/educrm/educrm-backend/internal/usecase/group"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GroupHandler exposes group HTTP endpoints.
type GroupHandler struct {
	svc *groupsvc.Service
}

// NewGroupHandler constructs GroupHandler.
func NewGroupHandler(svc *groupsvc.Service) *GroupHandler {
	return &GroupHandler{svc: svc}
}

// Create godoc
// @Summary Create group
// @Description Creates a class group with one subject, one teacher, optional room, schedule window, monthly fee (minor units), and status.
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateGroupRequest true "Group"
// @Success 201 {object} response.Envelope{data=dto.GroupResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/groups [post]
func (h *GroupHandler) Create(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateGroupRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	start, err := dto.ParseISODate(req.StartDate)
	if err != nil {
		response.Error(c, apperror.Validation("start_date", "Use date format YYYY-MM-DD"))
		return
	}
	end, err := dto.ParseISODate(req.EndDate)
	if err != nil {
		response.Error(c, apperror.Validation("end_date", "Use date format YYYY-MM-DD"))
		return
	}
	status := domain.GroupStatusActive
	if req.Status != "" {
		st, err := domain.ParseGroupStatus(req.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be active or inactive"))
			return
		}
		status = st
	}
	out, err := h.svc.Create(c.Request.Context(), groupsvc.CreateInput{
		Name:            req.Name,
		SubjectID:       req.SubjectID,
		TeacherID:       req.TeacherID,
		RoomID:          req.RoomID,
		StartDate:       start,
		EndDate:         end,
		MonthlyFeeMinor: req.MonthlyFee,
		Status:          status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.GroupResponseFrom(out))
}

// List godoc
// @Summary List groups
// @Description Paginated list with optional filters.
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page (1-based)" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Param q query string false "Search group name"
// @Param status query string false "Filter: active or inactive"
// @Param teacher_id query string false "Filter by teacher UUID"
// @Param subject_id query string false "Filter by subject UUID"
// @Param room_id query string false "Filter by room UUID"
// @Success 200 {object} response.Envelope{data=dto.GroupListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/groups [get]
func (h *GroupHandler) List(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var q dto.GroupListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	params := repository.GroupListParams{
		Search:   q.Q,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	if q.Status != "" {
		st, err := domain.ParseGroupStatus(q.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be active or inactive"))
			return
		}
		params.Status = &st
	}
	if q.TeacherID != "" {
		id := uuid.MustParse(q.TeacherID)
		params.TeacherID = &id
	}
	if q.SubjectID != "" {
		id := uuid.MustParse(q.SubjectID)
		params.SubjectID = &id
	}
	if q.RoomID != "" {
		id := uuid.MustParse(q.RoomID)
		params.RoomID = &id
	}
	res, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.GroupListResponseFrom(res))
}

// GetByID godoc
// @Summary Get group
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.GroupResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/groups/{id} [get]
func (h *GroupHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.GroupResponseFrom(out))
}

// Update godoc
// @Summary Update group
// @Tags groups
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Group ID (UUID)"
// @Param body body dto.UpdateGroupRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.GroupResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/groups/{id} [patch]
func (h *GroupHandler) Update(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateGroupRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := groupsvc.UpdateInput{
		Name:              req.Name,
		SubjectID:         req.SubjectID,
		TeacherID:         req.TeacherID,
		MonthlyFeeMinor:   req.MonthlyFee,
	}
	if req.ClearRoom != nil && *req.ClearRoom {
		in.ClearRoom = true
	} else if req.RoomID != nil {
		in.RoomID = req.RoomID
	}
	if req.StartDate != nil {
		t, err := dto.ParseISODate(*req.StartDate)
		if err != nil {
			response.Error(c, apperror.Validation("start_date", "Use date format YYYY-MM-DD"))
			return
		}
		in.StartDate = &t
	}
	if req.EndDate != nil {
		t, err := dto.ParseISODate(*req.EndDate)
		if err != nil {
			response.Error(c, apperror.Validation("end_date", "Use date format YYYY-MM-DD"))
			return
		}
		in.EndDate = &t
	}
	if req.Status != nil {
		st, err := domain.ParseGroupStatus(*req.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be active or inactive"))
			return
		}
		in.Status = &st
	}
	out, err := h.svc.Update(c.Request.Context(), id, in)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.GroupResponseFrom(out))
}

// Delete godoc
// @Summary Delete group
// @Tags groups
// @Security BearerAuth
// @Produce json
// @Param id path string true "Group ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/groups/{id} [delete]
func (h *GroupHandler) Delete(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Group deleted successfully"})
}
