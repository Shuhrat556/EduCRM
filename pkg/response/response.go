package response

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/gin-gonic/gin"
)

// ginCtxRequestID matches middleware.CtxRequestIDKey for correlation in envelopes.
const ginCtxRequestID = "request_id"

func metaFromGin(c *gin.Context) map[string]any {
	if c == nil {
		return nil
	}
	v, ok := c.Get(ginCtxRequestID)
	if !ok {
		return nil
	}
	s, ok := v.(string)
	if !ok || s == "" {
		return nil
	}
	return map[string]any{"request_id": s}
}

// Envelope is the standard JSON envelope for API responses.
type Envelope struct {
	Success bool           `json:"success"`
	Data    any            `json:"data,omitempty"`
	Error   *ErrorBody     `json:"error,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

// ErrorBody is returned to clients for structured errors.
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Kind    string `json:"kind,omitempty"`
}

// JSON writes a success envelope with data.
func JSON(c *gin.Context, status int, data any) {
	env := Envelope{Success: true, Data: data}
	if meta := metaFromGin(c); meta != nil {
		env.Meta = meta
	}
	c.JSON(status, env)
}

// Error writes an error envelope. If err is *apperror.Error, status and body are derived.
func Error(c *gin.Context, err error) {
	if ae, ok := apperror.AsError(err); ok {
		env := Envelope{
			Success: false,
			Error: &ErrorBody{
				Code:    ae.Code,
				Message: ae.Message,
				Kind:    string(ae.Kind),
			},
		}
		if meta := metaFromGin(c); meta != nil {
			env.Meta = meta
		}
		c.JSON(ae.HTTPStatus(), env)
		return
	}
	env := Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    "internal_error",
			Message: "An unexpected error occurred",
			Kind:    string(apperror.KindInternal),
		},
	}
	if meta := metaFromGin(c); meta != nil {
		env.Meta = meta
	}
	c.JSON(http.StatusInternalServerError, env)
}

// NoContent sends 204 without a body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
