package domain

import "errors"

// ErrInvalidTeacherStatus is returned for unknown status values.
var ErrInvalidTeacherStatus = errors.New("invalid teacher status")
