package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// BindJSON parses JSON into dst and returns a consistent, client-friendly validation error.
func BindJSON(c *gin.Context, dst any) error {
	if err := c.ShouldBindJSON(dst); err != nil {
		return apperrorFromBind(err)
	}
	return nil
}

// BindQuery parses query parameters into dst.
func BindQuery(c *gin.Context, dst any) error {
	if err := c.ShouldBindQuery(dst); err != nil {
		return apperrorFromBind(err)
	}
	return nil
}

func apperrorFromBind(err error) error {
	var syntax *json.SyntaxError
	if errors.As(err, &syntax) {
		return apperror.Validation("invalid_json", "Request body is not valid JSON")
	}
	var unmarshal *json.UnmarshalTypeError
	if errors.As(err, &unmarshal) {
		return apperror.Validation("invalid_type",
			fmt.Sprintf("Field %q must be %s", unmarshal.Field, unmarshal.Type.String()))
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		return apperror.Validation("validation_failed", formatValidationErrors(ve))
	}
	msg := err.Error()
	if strings.Contains(msg, "unexpected end of JSON") {
		return apperror.Validation("invalid_json", "Request body is not valid JSON")
	}
	return apperror.Validation("invalid_request", msg)
}

func formatValidationErrors(ve validator.ValidationErrors) string {
	parts := make([]string, 0, len(ve))
	for _, fe := range ve {
		parts = append(parts, formatFieldError(fe))
	}
	return strings.Join(parts, "; ")
}

func formatFieldError(fe validator.FieldError) string {
	name := fe.Field()
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", name)
	case "required_with":
		return fmt.Sprintf("%s is required when related fields are set", name)
	case "required_with_all":
		return fmt.Sprintf("%s is required", name)
	case "min":
		return fmt.Sprintf("%s must be at least %s", name, fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s", name, fe.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", name)
	case "uuid":
		return fmt.Sprintf("%s must be a UUID", name)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", name, strings.ReplaceAll(fe.Param(), " ", ", "))
	default:
		return fmt.Sprintf("%s failed validation (%s)", name, fe.Tag())
	}
}
