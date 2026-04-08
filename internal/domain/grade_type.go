package domain

import (
	"errors"
	"fmt"
)

// ErrInvalidGradeType is returned for unknown grade_type values.
var ErrInvalidGradeType = errors.New("invalid grade_type")

// GradeType distinguishes who produced the weekly rating.
type GradeType string

const (
	GradeTypeTeacherEvaluation GradeType = "teacher_evaluation"
	GradeTypeStudentEvaluation GradeType = "student_evaluation"
)

var validGradeTypes = map[GradeType]struct{}{
	GradeTypeTeacherEvaluation: {},
	GradeTypeStudentEvaluation: {},
}

// ParseGradeType validates s.
func ParseGradeType(s string) (GradeType, error) {
	t := GradeType(s)
	if _, ok := validGradeTypes[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidGradeType, s)
	}
	return t, nil
}
