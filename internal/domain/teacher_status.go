package domain

import "fmt"

// TeacherStatus is persisted as a string.
type TeacherStatus string

const (
	TeacherStatusActive   TeacherStatus = "active"
	TeacherStatusInactive TeacherStatus = "inactive"
)

var validTeacherStatuses = map[TeacherStatus]struct{}{
	TeacherStatusActive:   {},
	TeacherStatusInactive: {},
}

// ParseTeacherStatus validates s.
func ParseTeacherStatus(s string) (TeacherStatus, error) {
	t := TeacherStatus(s)
	if _, ok := validTeacherStatuses[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidTeacherStatus, s)
	}
	return t, nil
}
