package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	jwtpkg "github.com/educrm/educrm-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Service orchestrates authentication use cases.
type Service struct {
	users      repository.UserRepository
	refresh    repository.RefreshTokenRepository
	jwt        *jwtpkg.Manager
	refreshTTL time.Duration
}

// NewService constructs an auth service.
func NewService(
	users repository.UserRepository,
	refresh repository.RefreshTokenRepository,
	jwtMgr *jwtpkg.Manager,
	refreshTTL time.Duration,
) *Service {
	return &Service{
		users:      users,
		refresh:    refresh,
		jwt:        jwtMgr,
		refreshTTL: refreshTTL,
	}
}

// TokenPair is issued on login and refresh.
type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string
}

// LoginResult is returned from portal-aware login.
type LoginResult struct {
	Tokens              *TokenPair
	ForcePasswordChange bool
	FullName            string
	Username            *string
	Email               *string
	Phone               *string
	Role                domain.Role
}

// UserView is a safe representation of the current user.
type UserView struct {
	ID                  uuid.UUID   `json:"id"`
	FullName            string      `json:"full_name"`
	Username            *string     `json:"username,omitempty"`
	Email               *string     `json:"email,omitempty"`
	Phone               *string     `json:"phone,omitempty"`
	Role                domain.Role `json:"role"`
	IsActive            bool        `json:"is_active"`
	ForcePasswordChange bool        `json:"force_password_change"`
}

// Login validates credentials and returns tokens (legacy; no portal check).
func (s *Service) Login(ctx context.Context, login, password string) (*TokenPair, error) {
	out, err := s.LoginWithRole(ctx, login, password, nil)
	if err != nil {
		return nil, err
	}
	return out.Tokens, nil
}

// LoginWithRole validates credentials. If expectedRole is non-nil, the user role must match (portal guard).
func (s *Service) LoginWithRole(ctx context.Context, login, password string, expectedRole *domain.Role) (*LoginResult, error) {
	u, err := s.users.FindByLogin(ctx, login)
	if err != nil {
		return nil, apperror.Internal("lookup user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.Unauthorized("invalid credentials")
	}
	if !u.IsActive {
		return nil, apperror.Unauthorized("invalid credentials")
	}
	if expectedRole != nil && u.Role != *expectedRole {
		return nil, apperror.New(apperror.KindForbidden, "wrong_portal", "this account cannot sign in on this portal")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, apperror.Unauthorized("invalid credentials")
	}
	pair, err := s.issueTokens(ctx, u)
	if err != nil {
		return nil, err
	}
	return &LoginResult{
		Tokens:              pair,
		ForcePasswordChange: u.ForcePasswordChange,
		FullName:            u.FullName,
		Username:            u.Username,
		Email:               u.Email,
		Phone:               u.Phone,
		Role:                u.Role,
	}, nil
}

// Refresh rotates refresh token and returns a new pair.
func (s *Service) Refresh(ctx context.Context, rawRefresh string) (*TokenPair, error) {
	rawRefresh = strings.TrimSpace(rawRefresh)
	if rawRefresh == "" {
		return nil, apperror.Validation("refresh_token", "refresh_token is required")
	}
	hash := hashRefreshToken(rawRefresh)
	userID, err := s.refresh.FindValidByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Unauthorized("invalid or expired refresh token")
		}
		return nil, apperror.Internal("refresh lookup").Wrap(err)
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.Unauthorized("invalid or expired refresh token")
	}
	if !u.IsActive {
		return nil, apperror.Unauthorized("invalid or expired refresh token")
	}
	newRaw, newHash, exp, err := s.newRefreshValues()
	if err != nil {
		return nil, err
	}
	if err := s.refresh.Replace(ctx, userID, hash, newHash, exp); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperror.Unauthorized("invalid or expired refresh token")
		}
		return nil, apperror.Internal("refresh rotate").Wrap(err)
	}
	access, err := s.jwt.GenerateAccessToken(u.ID.String(), string(u.Role))
	if err != nil {
		return nil, apperror.Internal("sign access token").Wrap(err)
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: newRaw,
		ExpiresIn:    int64(s.jwt.AccessTTL().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// Logout revokes refresh sessions using either a valid access token or a refresh token body.
func (s *Service) Logout(ctx context.Context, authorizationHeader, refreshBody string) error {
	auth := strings.TrimSpace(authorizationHeader)
	const bearer = "bearer "
	if len(auth) >= len(bearer) && strings.EqualFold(auth[:len(bearer)], bearer) {
		raw := strings.TrimSpace(auth[len(bearer):])
		claims, err := s.jwt.ParseAccessToken(raw)
		if err != nil {
			return apperror.Unauthorized("invalid or expired access token")
		}
		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			return apperror.Unauthorized("invalid subject")
		}
		if err := s.refresh.DeleteAllForUser(ctx, uid); err != nil {
			return apperror.Internal("revoke sessions").Wrap(err)
		}
		return nil
	}
	refreshBody = strings.TrimSpace(refreshBody)
	if refreshBody != "" {
		h := hashRefreshToken(refreshBody)
		if err := s.refresh.DeleteByHash(ctx, h); err != nil {
			return apperror.Internal("revoke session").Wrap(err)
		}
		return nil
	}
	return apperror.Validation("auth_logout", "send Authorization: Bearer <access_token> or JSON body with refresh_token")
}

// Me returns the current user without sensitive fields.
func (s *Service) Me(ctx context.Context, userID uuid.UUID) (*UserView, error) {
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return nil, apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return nil, apperror.NotFound("user")
	}
	return &UserView{
		ID:                  u.ID,
		FullName:            u.FullName,
		Username:            u.Username,
		Email:               u.Email,
		Phone:               u.Phone,
		Role:                u.Role,
		IsActive:            u.IsActive,
		ForcePasswordChange: u.ForcePasswordChange,
	}, nil
}

// ChangePassword replaces the password when the current password matches.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" || len(newPassword) < 8 {
		return apperror.Validation("new_password", "password must be at least 8 characters")
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return apperror.NotFound("user")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(currentPassword)); err != nil {
		return apperror.New(apperror.KindForbidden, "invalid_password", "current password is incorrect")
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return apperror.Internal("hash password").Wrap(err)
	}
	u.PasswordHash = hash
	u.ForcePasswordChange = false
	u.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, u); err != nil {
		return apperror.Internal("update password").Wrap(err)
	}
	return nil
}

// FirstLoginSetPassword sets a new password when force_password_change is true (no current password).
func (s *Service) FirstLoginSetPassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" || len(newPassword) < 8 {
		return apperror.Validation("new_password", "password must be at least 8 characters")
	}
	u, err := s.users.FindByID(ctx, userID)
	if err != nil {
		return apperror.Internal("load user").Wrap(err)
	}
	if u == nil {
		return apperror.NotFound("user")
	}
	if !u.ForcePasswordChange {
		return apperror.New(apperror.KindForbidden, "password_change_not_required", "use POST /auth/change-password with your current password")
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return apperror.Internal("hash password").Wrap(err)
	}
	u.PasswordHash = hash
	u.ForcePasswordChange = false
	u.UpdatedAt = time.Now().UTC()
	if err := s.users.Update(ctx, u); err != nil {
		return apperror.Internal("update password").Wrap(err)
	}
	return nil
}

func (s *Service) issueTokens(ctx context.Context, u *domain.User) (*TokenPair, error) {
	raw, h, exp, err := s.newRefreshValues()
	if err != nil {
		return nil, err
	}
	if err := s.refresh.Create(ctx, u.ID, h, exp); err != nil {
		return nil, apperror.Internal("store refresh token").Wrap(err)
	}
	access, err := s.jwt.GenerateAccessToken(u.ID.String(), string(u.Role))
	if err != nil {
		return nil, apperror.Internal("sign access token").Wrap(err)
	}
	return &TokenPair{
		AccessToken:  access,
		RefreshToken: raw,
		ExpiresIn:    int64(s.jwt.AccessTTL().Seconds()),
		TokenType:    "Bearer",
	}, nil
}

func (s *Service) newRefreshValues() (raw string, hash string, expiresAt time.Time, err error) {
	expiresAt = time.Now().UTC().Add(s.refreshTTL)
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", time.Time{}, apperror.Internal("entropy").Wrap(err)
	}
	raw = hex.EncodeToString(buf)
	return raw, hashRefreshToken(raw), expiresAt, nil
}

func hashRefreshToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// HashPassword returns a bcrypt hash for bootstrapping users (e.g. seeds or admin tooling).
func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
