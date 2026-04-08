package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	roomsvc "github.com/educrm/educrm-backend/internal/usecase/room"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// RoomHandler exposes room HTTP endpoints.
type RoomHandler struct {
	svc *roomsvc.Service
}

// NewRoomHandler constructs RoomHandler.
func NewRoomHandler(svc *roomsvc.Service) *RoomHandler {
	return &RoomHandler{svc: svc}
}

// Create godoc
// @Summary Create room
// @Description Creates a room for schedule planning (capacity and active/inactive status).
// @Tags rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateRoomRequest true "Room"
// @Success 201 {object} response.Envelope{data=dto.RoomResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/rooms [post]
func (h *RoomHandler) Create(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateRoomRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	status := domain.RoomStatusActive
	if req.Status != "" {
		s, err := domain.ParseRoomStatus(req.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be active or inactive"))
			return
		}
		status = s
	}
	out, err := h.svc.Create(c.Request.Context(), roomsvc.CreateInput{
		Name:        req.Name,
		Capacity:    req.Capacity,
		Description: req.Description,
		Status:      status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.RoomResponseFrom(out))
}

// List godoc
// @Summary List rooms
// @Description Paginated list with optional status and text search.
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page (1-based)" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Param status query string false "Filter: active or inactive"
// @Param q query string false "Search name or description"
// @Success 200 {object} response.Envelope{data=dto.RoomListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/rooms [get]
func (h *RoomHandler) List(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var q dto.RoomListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	params := repository.RoomListParams{
		Search:   q.Q,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	if q.Status != "" {
		st, err := domain.ParseRoomStatus(q.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be active or inactive"))
			return
		}
		params.Status = &st
	}
	res, err := h.svc.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.RoomListResponseFrom(res))
}

// GetByID godoc
// @Summary Get room by ID
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.RoomResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/rooms/{id} [get]
func (h *RoomHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.RoomResponseFrom(out))
}

// Update godoc
// @Summary Update room
// @Tags rooms
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Room ID (UUID)"
// @Param body body dto.UpdateRoomRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.RoomResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/rooms/{id} [patch]
func (h *RoomHandler) Update(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateRoomRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := roomsvc.UpdateInput{
		Name:        req.Name,
		Capacity:    req.Capacity,
		Description: req.Description,
	}
	if req.Status != nil {
		st, err := domain.ParseRoomStatus(*req.Status)
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
	response.JSON(c, http.StatusOK, dto.RoomResponseFrom(out))
}

// Delete godoc
// @Summary Delete room
// @Tags rooms
// @Security BearerAuth
// @Produce json
// @Param id path string true "Room ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/rooms/{id} [delete]
func (h *RoomHandler) Delete(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Room deleted successfully"})
}
