package domain

import "errors"

// ErrInvalidRoomStatus is returned for unknown room status values.
var ErrInvalidRoomStatus = errors.New("invalid room status")
