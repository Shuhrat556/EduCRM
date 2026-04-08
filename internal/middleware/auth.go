package middleware

import (
	"strings"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/pkg/jwt"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ctxAuthUserID = "auth_user_id"
	ctxAuthRole   = "auth_role"
)

// AuthRequired validates a Bearer access JWT and stores user id and role in the Gin context.
func AuthRequired(jwtMgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		raw, ok := bearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.Abort()
			response.Error(c, apperror.Unauthorized("missing or invalid authorization header"))
			return
		}
		claims, err := jwtMgr.ParseAccessToken(raw)
		if err != nil {
			c.Abort()
			response.Error(c, apperror.Unauthorized("invalid or expired access token"))
			return
		}
		uid, err := uuid.Parse(claims.Subject)
		if err != nil {
			c.Abort()
			response.Error(c, apperror.Unauthorized("invalid subject"))
			return
		}
		if _, err := domain.ParseRole(claims.Role); err != nil {
			c.Abort()
			response.Error(c, apperror.Unauthorized("invalid role in token"))
			return
		}
		c.Set(ctxAuthUserID, uid)
		c.Set(ctxAuthRole, claims.Role)
		c.Next()
	}
}

// UserID returns the authenticated user id set by AuthRequired.
func UserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ctxAuthUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

// UserRole returns the role string from the access token claims.
func UserRole(c *gin.Context) (string, bool) {
	v, ok := c.Get(ctxAuthRole)
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

// ParseActor returns the authenticated role and user id (requires AuthRequired).
func ParseActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	uid, ok := UserID(c)
	if !ok {
		return "", uuid.Nil, apperror.Unauthorized("missing user")
	}
	rs, ok := UserRole(c)
	if !ok {
		return "", uuid.Nil, apperror.Unauthorized("missing role")
	}
	role, err := domain.ParseRole(rs)
	if err != nil {
		return "", uuid.Nil, apperror.Unauthorized("invalid role")
	}
	return role, uid, nil
}

// RequireRoles denies the request when the caller's role is not in the allow list.
// Prefer rbac.RequirePermission / RequireAnyPermission in router for explicit capability checks.
func RequireRoles(allowed ...domain.Role) gin.HandlerFunc {
	set := make(map[string]struct{}, len(allowed))
	for _, r := range allowed {
		set[string(r)] = struct{}{}
	}
	return func(c *gin.Context) {
		role, ok := UserRole(c)
		if !ok {
			c.Abort()
			response.Error(c, apperror.Unauthorized("missing role"))
			return
		}
		if _, ok := set[role]; !ok {
			c.Abort()
			response.Error(c, apperror.Forbidden("insufficient role"))
			return
		}
		c.Next()
	}
}

func bearerToken(header string) (token string, ok bool) {
	h := strings.TrimSpace(header)
	const prefix = "bearer "
	if len(h) < len(prefix) || !strings.EqualFold(h[:len(prefix)], prefix) {
		return "", false
	}
	t := strings.TrimSpace(h[len(prefix):])
	if t == "" {
		return "", false
	}
	return t, true
}
