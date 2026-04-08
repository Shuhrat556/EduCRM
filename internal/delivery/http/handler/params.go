package handler

import (
	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PathUUID parses a path parameter as UUID (e.g. Param("id") for /resource/:id).
func PathUUID(c *gin.Context, param string) (uuid.UUID, error) {
	raw := c.Param(param)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperror.Validation(param, "Invalid ID: use a UUID in canonical form (e.g. 550e8400-e29b-41d4-a716-446655440000)")
	}
	return id, nil
}
