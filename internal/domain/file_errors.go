package domain

import "errors"

var (
	// ErrInvalidFileOwnerType is returned for unknown owner_type strings.
	ErrInvalidFileOwnerType = errors.New("invalid file owner type")
)
