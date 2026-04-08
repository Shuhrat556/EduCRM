package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS configures cross-origin responses. When allowedOrigins contains exactly "*",
// the wildcard origin is emitted (credentials must be false). Otherwise the request
// Origin header is echoed only when it matches an entry in allowedOrigins (case-sensitive).
func CORS(allowedOrigins []string, allowCredentials bool) gin.HandlerFunc {
	if len(allowedOrigins) == 0 {
		allowedOrigins = []string{"*"}
	}
	wildcardOnly := len(allowedOrigins) == 1 && allowedOrigins[0] == "*"

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		var allowOrigin string
		switch {
		case wildcardOnly:
			allowOrigin = "*"
		case origin != "" && originAllowed(origin, allowedOrigins):
			allowOrigin = origin
		case len(allowedOrigins) == 1:
			allowOrigin = allowedOrigins[0]
		default:
			// Non-browser or same-origin: still set a safe default for tools that send Origin.
			allowOrigin = allowedOrigins[0]
		}

		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Request-ID")
		c.Header("Access-Control-Max-Age", "86400")
		if allowCredentials && allowOrigin != "" && allowOrigin != "*" {
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func originAllowed(origin string, allowed []string) bool {
	for _, o := range allowed {
		if o == origin {
			return true
		}
	}
	return false
}
