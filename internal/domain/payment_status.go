package domain

import "fmt"

// PaymentStatus reflects collection state for monthly billing.
type PaymentStatus string

const (
	PaymentStatusPaidFull    PaymentStatus = "paid_full"
	PaymentStatusPaidPartial PaymentStatus = "paid_partial"
	PaymentStatusUnpaid      PaymentStatus = "unpaid"
	PaymentStatusOverdue     PaymentStatus = "overdue"
)

var validPaymentStatuses = map[PaymentStatus]struct{}{
	PaymentStatusPaidFull:    {},
	PaymentStatusPaidPartial: {},
	PaymentStatusUnpaid:      {},
	PaymentStatusOverdue:     {},
}

// ParsePaymentStatus validates s.
func ParsePaymentStatus(s string) (PaymentStatus, error) {
	t := PaymentStatus(s)
	if _, ok := validPaymentStatuses[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidPaymentStatus, s)
	}
	return t, nil
}
