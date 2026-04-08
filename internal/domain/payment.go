package domain

import (
	"time"

	"github.com/google/uuid"
)

// Payment is a billing row for a student in a group (supports partials via multiple rows per month).
type Payment struct {
	ID                  uuid.UUID
	StudentID           uuid.UUID
	GroupID             uuid.UUID
	AmountMinor         int64
	Status              PaymentStatus
	PaymentDate         *time.Time // calendar date when money received; nil for unpaid/overdue if unknown
	MonthFor            time.Time  // first day of billed month (UTC)
	PaymentType         PaymentType
	Comment             *string
	IsFree              bool
	DiscountAmountMinor int64
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

// MonthStartUTC returns the first instant of the calendar month for t (UTC).
func MonthStartUTC(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}
