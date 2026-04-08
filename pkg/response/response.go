package response

import (
	"net/http"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/gin-gonic/gin"
)

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
	c.JSON(status, Envelope{Success: true, Data: data})
}

// Error writes an error envelope. If err is *apperror.Error, status and body are derived.
func Error(c *gin.Context, err error) {
	if ae, ok := apperror.AsError(err); ok {
		c.JSON(ae.HTTPStatus(), Envelope{
			Success: false,
			Error: &ErrorBody{
				Code:    ae.Code,
				Message: ae.Message,
				Kind:    string(ae.Kind),
			},
		})
		return
	}
	c.JSON(http.StatusInternalServerError, Envelope{
		Success: false,
		Error: &ErrorBody{
			Code:    "internal_error",
			Message: "An unexpected error occurred",
			Kind:    string(apperror.KindInternal),
		},
	})
}

// NoContent sends 204 without a body.
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
