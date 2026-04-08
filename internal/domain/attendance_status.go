package domain

import (
	"errors"
	"fmt"
)

// ErrInvalidAttendanceStatus is returned for unknown attendance status values.
var ErrInvalidAttendanceStatus = errors.New("invalid attendance status")

// AttendanceStatus is persisted as a string.
type AttendanceStatus string

const (
	AttendancePresent AttendanceStatus = "present"
	AttendanceAbsent  AttendanceStatus = "absent"
	AttendanceLate    AttendanceStatus = "late"
)

var validAttendanceStatuses = map[AttendanceStatus]struct{}{
	AttendancePresent: {},
	AttendanceAbsent:  {},
	AttendanceLate:    {},
}

// ParseAttendanceStatus validates s.
func ParseAttendanceStatus(s string) (AttendanceStatus, error) {
	t := AttendanceStatus(s)
	if _, ok := validAttendanceStatuses[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidAttendanceStatus, s)
	}
	return t, nil
}
