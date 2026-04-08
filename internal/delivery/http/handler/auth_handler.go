package handler

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/delivery/http/dto"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/usecase/auth"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthHandler exposes auth HTTP endpoints.
type AuthHandler struct {
	svc *auth.Service
}

// NewAuthHandler constructs AuthHandler.
func NewAuthHandler(svc *auth.Service) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login godoc
// @Summary Login with email or phone
// @Description Authenticates a user by email (if login contains @) or phone and returns JWT access and refresh tokens.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Credentials"
// @Success 200 {object} response.Envelope{data=dto.TokenResponse}
// @Failure 400 {object} response.Envelope "Invalid JSON or validation error"
// @Failure 401 {object} response.Envelope "Invalid credentials"
// @Failure 500 {object} response.Envelope
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	pair, err := h.svc.Login(c.Request.Context(), req.Login, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.TokenResponseFrom(pair))
}

// Refresh godoc
// @Summary Refresh tokens
// @Description Exchanges a valid refresh token for a new access/refresh pair (rotation).
// @Tags auth
// @Accept json
// @Produce json
// @Param body body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} response.Envelope{data=dto.TokenResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req dto.RefreshRequest
	if err := BindJSON(c, &req); err != nil {
		response.Error(c, err)
		return
	}
	pair, err := h.svc.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.TokenResponseFrom(pair))
}

// Logout godoc
// @Summary Logout
// @Description Revokes refresh tokens using either Authorization Bearer access token (all sessions) or refresh_token in body (single session).
// @Tags auth
// @Accept json
// @Produce json
// @Param Authorization header string false "Bearer access_token"
// @Param body body dto.LogoutRequest false "Optional refresh_token"
// @Success 200 {object} response.Envelope{data=dto.MessageResponse}
// @Failure 400 {object} response.Envelope
// @Failure 401 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req dto.LogoutRequest
	if c.Request.ContentLength > 0 {
		_ = c.ShouldBindJSON(&req)
	}
	if err := h.svc.Logout(c.Request.Context(), c.GetHeader("Authorization"), req.RefreshToken); err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.MessageResponse{Message: "Logged out successfully"})
}

// Me godoc
// @Summary Current user
// @Description Returns the authenticated user profile (requires access token).
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Envelope{data=dto.CurrentUserResponse}
// @Failure 401 {object} response.Envelope
// @Failure 404 {object} response.Envelope
// @Failure 500 {object} response.Envelope
// @Router /api/v1/auth/me [get]
func (h *AuthHandler) Me(c *gin.Context) {
	uidVal, ok := middleware.UserID(c)
	if !ok {
		response.Error(c, errNoUserID())
		return
	}
	u, err := h.svc.Me(c.Request.Context(), uidVal)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.JSON(c, http.StatusOK, dto.CurrentUserFrom(u))
}

func errNoUserID() error {
	return apperror.Unauthorized("missing authentication context")
}
