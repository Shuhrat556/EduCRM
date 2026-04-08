package middleware

import (
	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// RequirePermission aborts with 403 unless the actor's role grants the permission.
func RequirePermission(p rbac.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _, err := ParseActor(c)
		if err != nil {
			c.Abort()
			response.Error(c, err)
			return
		}
		if !rbac.Granted(role, p) {
			c.Abort()
			response.Error(c, apperror.Forbidden("insufficient permissions"))
			return
		}
		c.Next()
	}
}

// RequireAnyPermission aborts unless the actor has at least one permission.
func RequireAnyPermission(perms ...rbac.Permission) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _, err := ParseActor(c)
		if err != nil {
			c.Abort()
			response.Error(c, err)
			return
		}
		if !rbac.GrantedAny(role, perms...) {
			c.Abort()
			response.Error(c, apperror.Forbidden("insufficient permissions"))
			return
		}
		c.Next()
	}
}
