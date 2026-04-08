package domain

import "errors"

// ErrInvalidSubjectStatus is returned for unknown status values.
var ErrInvalidSubjectStatus = errors.New("invalid subject status")
