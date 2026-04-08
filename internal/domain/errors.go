package domain

import "errors"

var (
	// ErrInvalidRole is returned when a role string is not one of the known roles.
	ErrInvalidRole = errors.New("invalid role")
)
