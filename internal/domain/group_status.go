package domain

import "fmt"

// GroupStatus is persisted as a string.
type GroupStatus string

const (
	GroupStatusActive   GroupStatus = "active"
	GroupStatusInactive GroupStatus = "inactive"
)

var validGroupStatuses = map[GroupStatus]struct{}{
	GroupStatusActive:   {},
	GroupStatusInactive: {},
}

// ParseGroupStatus validates s.
func ParseGroupStatus(s string) (GroupStatus, error) {
	t := GroupStatus(s)
	if _, ok := validGroupStatuses[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidGroupStatus, s)
	}
	return t, nil
}
