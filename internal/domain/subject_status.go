package domain

import "fmt"

// SubjectStatus is persisted as a string (course/subject catalog).
type SubjectStatus string

const (
	SubjectStatusActive   SubjectStatus = "active"
	SubjectStatusInactive SubjectStatus = "inactive"
)

var validSubjectStatuses = map[SubjectStatus]struct{}{
	SubjectStatusActive:   {},
	SubjectStatusInactive: {},
}

// ParseSubjectStatus validates s.
func ParseSubjectStatus(s string) (SubjectStatus, error) {
	t := SubjectStatus(s)
	if _, ok := validSubjectStatuses[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidSubjectStatus, s)
	}
	return t, nil
}
