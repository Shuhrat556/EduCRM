package domain

import "time"

// WeekStartUTC returns Monday 00:00:00 UTC for the ISO week containing t.
func WeekStartUTC(t time.Time) time.Time {
	t = t.UTC()
	wd := int(t.Weekday()) // Sunday=0, Monday=1, ...
	daysFromMonday := (wd + 6) % 7
	d := t.AddDate(0, 0, -daysFromMonday)
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, time.UTC)
}
