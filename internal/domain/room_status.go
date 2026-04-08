package domain

import "fmt"

// RoomStatus is persisted as a string.
type RoomStatus string

const (
	RoomStatusActive   RoomStatus = "active"
	RoomStatusInactive RoomStatus = "inactive"
)

var validRoomStatuses = map[RoomStatus]struct{}{
	RoomStatusActive:   {},
	RoomStatusInactive: {},
}

// ParseRoomStatus validates s.
func ParseRoomStatus(s string) (RoomStatus, error) {
	t := RoomStatus(s)
	if _, ok := validRoomStatuses[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidRoomStatus, s)
	}
	return t, nil
}
