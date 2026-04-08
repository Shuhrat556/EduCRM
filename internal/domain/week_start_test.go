package domain

import (
	"testing"
	"time"
)

func TestWeekStartUTC(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want string // YYYY-MM-DD Monday UTC
	}{
		{
			name: "Wednesday to same week Monday",
			in:   time.Date(2025, 4, 9, 15, 0, 0, 0, time.UTC), // Wed
			want: "2025-04-07",
		},
		{
			name: "Sunday maps to previous Monday",
			in:   time.Date(2025, 4, 6, 0, 0, 0, 0, time.UTC), // Sun
			want: "2025-03-31",
		},
		{
			name: "Monday unchanged",
			in:   time.Date(2025, 4, 7, 0, 0, 0, 0, time.UTC),
			want: "2025-04-07",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeekStartUTC(tt.in)
			if g := got.UTC().Format("2006-01-02"); g != tt.want {
				t.Fatalf("WeekStartUTC(%v) = %s, want %s", tt.in, g, tt.want)
			}
		})
	}
}
