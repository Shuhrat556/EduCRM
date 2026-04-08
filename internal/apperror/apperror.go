package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// Kind classifies errors for HTTP mapping and logging.
type Kind string

const (
	KindInternal        Kind = "internal"
	KindValidation      Kind = "validation"
	KindNotFound        Kind = "not_found"
	KindConflict        Kind = "conflict"
	KindUnauthorized    Kind = "unauthorized"
	KindForbidden       Kind = "forbidden"
	KindTooManyRequests Kind = "too_many_requests"
)

// Error is the application error type with optional wrapped cause.
type Error struct {
	Kind    Kind
	Code    string
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

// HTTPStatus maps the error kind to an HTTP status code.
func (e *Error) HTTPStatus() int {
	switch e.Kind {
	case KindValidation:
		return http.StatusBadRequest
	case KindNotFound:
		return http.StatusNotFound
	case KindConflict:
		return http.StatusConflict
	case KindUnauthorized:
		return http.StatusUnauthorized
	case KindForbidden:
		return http.StatusForbidden
	case KindTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// New constructs an application error.
func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

// Wrap adds a cause while preserving kind/code/message.
func (e *Error) Wrap(err error) *Error {
	out := *e
	out.Err = err
	return &out
}

// AsError extracts *Error from err chain.
func AsError(err error) (*Error, bool) {
	var ae *Error
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

// Common constructors for handlers and services.

func Internal(message string) *Error {
	return New(KindInternal, "internal_error", message)
}

func Validation(code, message string) *Error {
	return New(KindValidation, code, message)
}

func NotFound(resource string) *Error {
	return New(KindNotFound, "not_found", "The requested "+resource+" was not found")
}

func Conflict(code, message string) *Error {
	return New(KindConflict, code, message)
}

func Unauthorized(message string) *Error {
	return New(KindUnauthorized, "unauthorized", message)
}

func Forbidden(message string) *Error {
	return New(KindForbidden, "forbidden", message)
}

// TooManyRequests indicates the client exceeded a rate limit.
func TooManyRequests(message string) *Error {
	return New(KindTooManyRequests, "rate_limited", message)
}
