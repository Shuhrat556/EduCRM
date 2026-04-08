package middleware

import (
	"log/slog"
	"runtime/debug"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recovery returns a Gin middleware that logs panics and returns 500 JSON.
func Recovery(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered",
					"error", rec,
					"path", c.Request.URL.Path,
					"stack", string(debug.Stack()),
				)
				c.Abort()
				response.Error(c, apperror.Internal("Internal server error"))
			}
		}()
		c.Next()
	}
}

// ErrorHandler maps handler errors to JSON when handlers use c.Error without writing a body.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Written() || len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		response.Error(c, err)
	}
}
