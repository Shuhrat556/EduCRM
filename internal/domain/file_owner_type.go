package domain

import "fmt"

// FileOwnerType classifies what entity owns uploaded content.
type FileOwnerType string

const (
	FileOwnerStudentPhoto FileOwnerType = "student_photo"
	FileOwnerTeacherPhoto FileOwnerType = "teacher_photo"
	FileOwnerDocument     FileOwnerType = "document"
)

var validFileOwnerTypes = map[FileOwnerType]struct{}{
	FileOwnerStudentPhoto: {},
	FileOwnerTeacherPhoto: {},
	FileOwnerDocument:     {},
}

// ParseFileOwnerType validates s.
func ParseFileOwnerType(s string) (FileOwnerType, error) {
	t := FileOwnerType(s)
	if _, ok := validFileOwnerTypes[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidFileOwnerType, s)
	}
	return t, nil
}
