package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const headerRequestID = "X-Request-ID"
const ctxRequestIDKey = "request_id"

// RequestID ensures each request has a correlation ID for logs and tracing.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader(headerRequestID)
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Writer.Header().Set(headerRequestID, rid)
		c.Set(ctxRequestIDKey, rid)
		c.Next()
	}
}

// GetRequestID returns the request ID from context, if set.
func GetRequestID(c *gin.Context) string {
	v, ok := c.Get(ctxRequestIDKey)
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
