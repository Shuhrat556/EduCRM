package dto

import (
	"github.com/educrm/educrm-backend/internal/usecase/auth"
	"github.com/google/uuid"
)

// LoginRequest is the body for POST /auth/login.
type LoginRequest struct {
	Login    string `json:"login" example:"teacher@school.edu" binding:"required,min=3,max=255"`
	Password string `json:"password" example:"your-secure-password" binding:"required,min=8,max=128"`
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

// CurrentUserResponse is returned by GET /auth/me.
type CurrentUserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    *string   `json:"email,omitempty"`
	Phone    *string   `json:"phone,omitempty"`
	Role     string    `json:"role"`
	IsActive bool      `json:"is_active"`
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
		ID:       v.ID,
		Email:    v.Email,
		Phone:    v.Phone,
		Role:     string(v.Role),
		IsActive: v.IsActive,
	}
}
