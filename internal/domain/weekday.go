package domain

import (
	"errors"
	"fmt"
)

// ErrInvalidWeekday is returned when weekday is not in 0–6.
var ErrInvalidWeekday = errors.New("invalid weekday: expected 0–6 (Sunday–Saturday)")

// Weekday is 0=Sunday through 6=Saturday (same as Go time.Weekday and PostgreSQL EXTRACT(DOW)).
type Weekday int8

const (
	WeekdaySunday    Weekday = 0
	WeekdayMonday    Weekday = 1
	WeekdayTuesday   Weekday = 2
	WeekdayWednesday Weekday = 3
	WeekdayThursday  Weekday = 4
	WeekdayFriday    Weekday = 5
	WeekdaySaturday  Weekday = 6
)

// ParseWeekday validates w is in [0,6].
func ParseWeekday(w int) (Weekday, error) {
	if w < 0 || w > 6 {
		return 0, fmt.Errorf("%w", ErrInvalidWeekday)
	}
	return Weekday(w), nil
}
