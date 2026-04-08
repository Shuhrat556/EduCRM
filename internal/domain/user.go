package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// User is a domain entity (no persistence tags).
type User struct {
	ID           uuid.UUID
	Email        *string
	Phone        *string
	PasswordHash string
	Role         Role
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NormalizeEmail returns a trimmed lower-case email or nil when empty.
func NormalizeEmail(s *string) *string {
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

// NormalizePhone returns trimmed phone or nil when empty.
func NormalizePhone(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}
