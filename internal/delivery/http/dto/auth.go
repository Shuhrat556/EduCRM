package dto

import (
	"github.com/educrm/educrm-backend/internal/usecase/auth"
	"github.com/google/uuid"
)

// LoginRequest is the body for POST /auth/login.
type LoginRequest struct {
	Login    string `json:"login" example:"teacher01" binding:"required,min=3,max=255"`
	Password string `json:"password" example:"your-secure-password" binding:"required,min=8,max=128"`
	// Portal matches the UI entry: student | admin | teacher | super_admin. Required when AUTH_REQUIRE_LOGIN_PORTAL=true (default).
	// Must match the account role or login is rejected. Set AUTH_REQUIRE_LOGIN_PORTAL=false only for legacy clients.
	Portal string `json:"portal" example:"student" binding:"omitempty,oneof=student admin teacher super_admin"`
}

// RefreshRequest is the body for POST /auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"opaque-refresh-token" binding:"required"`
}

// LogoutRequest optionally carries a refresh token when no Bearer header is sent.
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// TokenResponse matches the standard OAuth-style token payload.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// PortalLoginResponse extends tokens with session flags for front-end routing.
type PortalLoginResponse struct {
	TokenResponse
	ForcePasswordChange bool    `json:"force_password_change"`
	FullName            string  `json:"full_name"`
	Username            *string `json:"username,omitempty"`
	Role                string  `json:"role"`
}

// ChangePasswordRequest is the body for POST /auth/change-password (authenticated).
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=1,max=128"`
	NewPassword     string `json:"new_password" binding:"required,min=8,max=128"`
}

// FirstLoginPasswordRequest is the body for POST /auth/first-login/change-password.
type FirstLoginPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8,max=128"`
}

// CurrentUserResponse is returned by GET /auth/me.
type CurrentUserResponse struct {
	ID                  uuid.UUID `json:"id"`
	FullName            string    `json:"full_name"`
	Username            *string   `json:"username,omitempty"`
	Email               *string   `json:"email,omitempty"`
	Phone               *string   `json:"phone,omitempty"`
	Role                string    `json:"role"`
	IsActive            bool      `json:"is_active"`
	ForcePasswordChange bool      `json:"force_password_change"`
}

// TokenResponseFrom maps a service token pair to the API DTO.
func TokenResponseFrom(p *auth.TokenPair) TokenResponse {
	if p == nil {
		return TokenResponse{}
	}
	return TokenResponse{
		AccessToken:  p.AccessToken,
		RefreshToken: p.RefreshToken,
		TokenType:    p.TokenType,
		ExpiresIn:    p.ExpiresIn,
	}
}

// CurrentUserFrom maps a user view to the API DTO.
func CurrentUserFrom(v *auth.UserView) CurrentUserResponse {
	if v == nil {
		return CurrentUserResponse{}
	}
	return CurrentUserResponse{
		ID:                  v.ID,
		FullName:            v.FullName,
		Username:            v.Username,
		Email:               v.Email,
		Phone:               v.Phone,
		Role:                string(v.Role),
		IsActive:            v.IsActive,
		ForcePasswordChange: v.ForcePasswordChange,
	}
}

// PortalLoginResponseFrom maps login result to API DTO.
func PortalLoginResponseFrom(res *auth.LoginResult) PortalLoginResponse {
	if res == nil || res.Tokens == nil {
		return PortalLoginResponse{}
	}
	return PortalLoginResponse{
		TokenResponse:       TokenResponseFrom(res.Tokens),
		ForcePasswordChange: res.ForcePasswordChange,
		FullName:            res.FullName,
		Username:            res.Username,
		Role:                string(res.Role),
	}
}
