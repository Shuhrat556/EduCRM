package repository

import "errors"

// ErrNotFound is returned when a lookup yields no row.
var ErrNotFound = errors.New("repository: not found")

// ErrDuplicate is returned when a unique constraint is violated.
var ErrDuplicate = errors.New("repository: duplicate key")

// ErrReferenced is returned when a row cannot be deleted due to foreign key references.
var ErrReferenced = errors.New("repository: referenced by other rows")
