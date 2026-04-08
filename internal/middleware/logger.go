package middleware

import (
	"log/slog"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs one line per request with duration and status.
// skipPrefixes skips logging when the path starts with any of the given prefixes.
func RequestLogger(log *slog.Logger, skipPrefixes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if shouldSkipPath(path, skipPrefixes) {
			c.Next()
			return
		}
		start := time.Now()
		method := c.Request.Method
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()
		rid := GetRequestID(c)

		log.Info("http_request",
			"method", method,
			"path", path,
			"status", status,
			"duration_ms", duration.Milliseconds(),
			"request_id", rid,
			"client_ip", c.ClientIP(),
		)
	}
}

func shouldSkipPath(path string, prefixes []string) bool {
	for _, p := range prefixes {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if path == p || strings.HasPrefix(path, p+"/") {
			return true
		}
	}
	return false
}
