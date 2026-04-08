package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	teachersvc "github.com/educrm/educrm-backend/internal/usecase/teacher"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// TeacherHandler exposes teacher HTTP endpoints.
type TeacherHandler struct {
	svc *teachersvc.Service
}

// NewTeacherHandler constructs TeacherHandler.
func NewTeacherHandler(svc *teachersvc.Service) *TeacherHandler {
	return &TeacherHandler{svc: svc}
}

// Create godoc
// @Summary Create teacher
// @Description Creates a teacher. Assign groups via the Groups API (each group has one teacher_id).
// @Tags teachers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateTeacherRequest true "Teacher"
// @Success 201 {object} response.Envelope{data=dto.TeacherDetailResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/teachers [post]
func (h *TeacherHandler) Create(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateTeacherRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	status := domain.TeacherStatusActive
	if req.Status != "" {
		s, err := domain.ParseTeacherStatus(req.Status)
		if err != nil {
			response.Error(c, apperror.Validation("status", "Status must be active or inactive"))
			return
		}
		status = s
	}
	out, err := h.svc.Create(c.Request.Context(), teachersvc.CreateInput{
		FullName:          req.FullName,
		Phone:             req.Phone,
		Email:             req.Email,
		Specialization:    req.Specialization,
		PhotoURL:          req.PhotoURL,
		PhotoStorageKey:   req.PhotoStorageKey,
		PhotoContentType:  req.PhotoContentType,
		PhotoOriginalName: req.PhotoOriginalName,
		Status:            status,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.TeacherDetailFrom(out))
}

// List godoc
// @Summary List teachers
// @Description Paginated list with optional status and search.
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page (1-based)" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Param status query string false "Filter: active or inactive"
// @Param q query string false "Search name, email, phone, specialization"
// @Success 200 {object} response.Envelope{data=dto.TeacherListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/teachers [get]
func (h *TeacherHandler) List(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	var q dto.TeacherListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	params := repository.TeacherListParams{
		Search:   q.Q,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	if q.Status != "" {
		st, err := domain.ParseTeacherStatus(q.Status)
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
	response.JSON(c, http.StatusOK, dto.TeacherListResponseFrom(res))
}

// GetByID godoc
// @Summary Get teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.TeacherDetailResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/teachers/{id} [get]
func (h *TeacherHandler) GetByID(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.TeacherDetailFrom(out))
}

// Update godoc
// @Summary Update teacher
// @Tags teachers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID (UUID)"
// @Param body body dto.UpdateTeacherRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.TeacherDetailResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/teachers/{id} [patch]
func (h *TeacherHandler) Update(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateTeacherRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	in := teachersvc.UpdateInput{
		FullName:          req.FullName,
		Phone:             req.Phone,
		Email:             req.Email,
		Specialization:    req.Specialization,
		PhotoURL:          req.PhotoURL,
		PhotoStorageKey:   req.PhotoStorageKey,
		PhotoContentType:  req.PhotoContentType,
		PhotoOriginalName: req.PhotoOriginalName,
	}
	if req.Status != nil {
		st, err := domain.ParseTeacherStatus(*req.Status)
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
	response.JSON(c, http.StatusOK, dto.TeacherDetailFrom(out))
}

// PatchPhoto godoc
// @Summary Update teacher photo metadata
// @Description Stores public URL and optional storage metadata after the client uploads the file to object storage.
// @Tags teachers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Teacher ID (UUID)"
// @Param body body dto.PatchTeacherPhotoRequest true "Photo metadata"
// @Success 200 {object} response.Envelope{data=dto.TeacherDetailResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/teachers/{id}/photo [patch]
func (h *TeacherHandler) PatchPhoto(c *gin.Context) {
	if _, err := RequireStaff(c); err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.PatchTeacherPhotoRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.PatchPhoto(c.Request.Context(), id, teachersvc.PhotoPatchInput{
		PhotoURL:          req.PhotoURL,
		PhotoStorageKey:   req.PhotoStorageKey,
		PhotoContentType:  req.PhotoContentType,
		PhotoOriginalName: req.PhotoOriginalName,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.TeacherDetailFrom(out))
}

// Delete godoc
// @Summary Delete teacher
// @Tags teachers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Teacher ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/teachers/{id} [delete]
func (h *TeacherHandler) Delete(c *gin.Context) {
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
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Teacher deleted successfully"})
}
