package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerRequestID = "X-Request-ID"

// CtxRequestIDKey is the Gin context key storing the correlation ID (also used by response envelopes).
const CtxRequestIDKey = "request_id"

// RequestID ensures each request has a correlation ID for logs and tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Writer.Header().Set(headerRequestID, rid)
		c.Set(CtxRequestIDKey, rid)
		c.Next()
	}
}

// GetRequestID returns the request ID from context, if set.
func GetRequestID(c *gin.Context) string {
	v, ok := c.Get(CtxRequestIDKey)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
