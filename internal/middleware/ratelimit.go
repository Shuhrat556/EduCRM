package middleware

import (
	"sync"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/pkg/response"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimit applies an in-memory per-IP token bucket. Disabled when rps <= 0.
// skipPrefixes are path prefixes that bypass the limiter (same rules as RequestLogger skips).
func RateLimit(rps float64, burst int, skipPrefixes []string) gin.HandlerFunc {
	if rps <= 0 {
		return func(c *gin.Context) { c.Next() }
	}
	if burst < 1 {
		burst = 1
	}
	lim := rate.Limit(rps)
	var mu sync.Mutex
	byIP := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if shouldSkipPath(path, skipPrefixes) {
			c.Next()
			return
		}

		ip := c.ClientIP()
		mu.Lock()
		limiter, ok := byIP[ip]
		if !ok {
			limiter = rate.NewLimiter(lim, burst)
			byIP[ip] = limiter
		}
		mu.Unlock()

		if !limiter.Allow() {
			response.Error(c, apperror.TooManyRequests("too many requests; try again later"))
			c.Abort()
			return
		}
		c.Next()
	}
}
