package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Teacher is a teaching staff record (independent from auth users unless linked later).
type Teacher struct {
	ID                uuid.UUID
	FullName          string
	Phone             *string
	Email             *string
	Specialization    *string
	PhotoURL          *string
	PhotoStorageKey   *string
	PhotoContentType  *string
	PhotoOriginalName *string
	Status            TeacherStatus
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NormalizeTeacherEmail lowercases trimmed email or nil.
func NormalizeTeacherEmail(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	t = strings.ToLower(t)
	return &t
}

// NormalizeTeacherPhone trims phone or nil.
func NormalizeTeacherPhone(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}
