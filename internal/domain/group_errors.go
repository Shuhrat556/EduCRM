package domain

import "errors"

// ErrInvalidGroupStatus is returned for unknown group status values.
var ErrInvalidGroupStatus = errors.New("invalid group status")
