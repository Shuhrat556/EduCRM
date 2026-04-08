package dto

import (
	"fmt"
	"strconv"
	"strings"
)

// ParseClockToMinutes parses "HH:MM" 24-hour to minutes from midnight.
// Allows 24:00 as 1440 (end-of-day) for interval ends only; callers should reject 24:00 for starts.
func ParseClockToMinutes(clock string) (int, error) {
	clock = strings.TrimSpace(clock)
	parts := strings.Split(clock, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("expected HH:MM")
	}
	h, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid hour")
	}
	m, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid minute")
	}
	if h == 24 && m == 0 {
		return 24 * 60, nil
	}
	if h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, fmt.Errorf("time out of range")
	}
	return h*60 + m, nil
}

// FormatMinutesAsClock renders minutes as "HH:MM" (24h). Values >= 1440 clamp to 24:00 for display.
func FormatMinutesAsClock(minutes int) string {
	if minutes < 0 {
		minutes = 0
	}
	if minutes > 24*60 {
		minutes = 24 * 60
	}
	h := minutes / 60
	m := minutes % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}
