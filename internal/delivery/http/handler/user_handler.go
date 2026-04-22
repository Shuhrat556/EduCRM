package handler

import (
	"net/http"
	"strings"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/repository"
	usersvc "github.com/educrm/educrm-backend/internal/usecase/user"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserHandler exposes user management HTTP endpoints.
type UserHandler struct {
	svc *usersvc.Service
}

// NewUserHandler constructs UserHandler.
func NewUserHandler(svc *usersvc.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

// Create godoc
// @Summary Create user
// @Description Creates a user. Only super_admin may create admin or super_admin. Admin may create teacher/student only. Password is stored hashed.
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body dto.CreateUserRequest true "User"
// @Success 201 {object} response.Envelope{data=dto.UserResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	actor, actorID, err := ParseActor(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.CreateUserRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	role, err := domain.ParseRole(req.Role)
	if err != nil {
		response.Error(c, apperror.Validation("role", "Role must be a valid value (e.g. teacher, student, admin, super_admin)"))
		return
	}
	createdBy := actorID
	out, err := h.svc.Create(c.Request.Context(), actor, usersvc.CreateInput{
		FullName:            strings.TrimSpace(req.FullName),
		Username:            req.Username,
		Email:               req.Email,
		Phone:               req.Phone,
		Password:            req.Password,
		Role:                role,
		ForcePasswordChange: req.ForcePasswordChange,
		CreatedByUserID:     &createdBy,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, dto.UserResponseFrom(out))
}

// List godoc
// @Summary List users
// @Description Paginated list with optional search (email/phone) and role filter. Admins do not see other admins or super_admins.
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page (1-based)" default(1)
// @Param page_size query int false "Page size (max 100)" default(20)
// @Param role query string false "Filter by role"
// @Param q query string false "Search email or phone"
// @Success 200 {object} response.Envelope{data=dto.UserListResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
	actor, err := ActorRole(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	var q dto.UserListQuery
	if err := BindQuery(c, &q); err != nil {
		response.Error(c, err)
		return
	}
	params := repository.UserListParams{
		Search:   q.Q,
		Page:     q.Page,
		PageSize: q.PageSize,
	}
	if q.Role != "" {
		r, err := domain.ParseRole(q.Role)
		if err != nil {
			response.Error(c, apperror.Validation("role", "Invalid role filter"))
			return
		}
		params.Role = &r
	}
	res, err := h.svc.List(c.Request.Context(), actor, params)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.UserListResponseFrom(res))
}

// GetByID godoc
// @Summary Get user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.UserResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	actor, err := ActorRole(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.GetByID(c.Request.Context(), actor, id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.UserResponseFrom(out))
}

// Update godoc
// @Summary Update user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param body body dto.UpdateUserRequest true "Fields to update"
// @Success 200 {object} response.Envelope{data=dto.UserResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 409 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/users/{id} [patch]
func (h *UserHandler) Update(c *gin.Context) {
	actor, err := ActorRole(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	actorID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("Missing authenticated user in request context"))
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.UpdateUserRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.Update(c.Request.Context(), actor, actorID, id, usersvc.UpdateInput{
		FullName: req.FullName,
		Username: req.Username,
		Email:    req.Email,
		Phone:    req.Phone,
		Password: req.Password,
		Role:     req.Role,
		IsActive: req.IsActive,
	})
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.UserResponseFrom(out))
}

// SetStatus godoc
// @Summary Activate or deactivate user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param body body dto.SetUserStatusRequest true "Status"
// @Success 200 {object} response.Envelope{data=dto.UserResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/users/{id}/status [patch]
func (h *UserHandler) SetStatus(c *gin.Context) {
	actor, err := ActorRole(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	actorID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("Missing authenticated user in request context"))
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	var req dto.SetUserStatusRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	out, err := h.svc.SetActive(c.Request.Context(), actor, actorID, id, *req.IsActive)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.UserResponseFrom(out))
}

// Delete godoc
// @Summary Delete user
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 403 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	actor, err := ActorRole(c)
	if err != nil {
		response.Error(c, err)
		return
	}
	actorID, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, apperror.Unauthorized("Missing authenticated user in request context"))
		return
	}
	id, err := PathUUID(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), actor, actorID, id); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "User deleted successfully"})
}
