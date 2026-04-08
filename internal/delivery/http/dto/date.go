package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
)

// ParseISODate parses a calendar date "YYYY-MM-DD" in UTC.
func ParseISODate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty date")
	}
	t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// ParseMonthFor accepts YYYY-MM or YYYY-MM-DD and returns the first day of that month (UTC).
func ParseMonthFor(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty month_for")
	}
	if len(s) == 7 && s[4] == '-' {
		t, err := time.ParseInLocation("2006-01", s, time.UTC)
		if err != nil {
			return time.Time{}, err
		}
		return domain.MonthStartUTC(t), nil
	}
	t, err := ParseISODate(s)
	if err != nil {
		return time.Time{}, err
	}
	return domain.MonthStartUTC(t), nil
}
