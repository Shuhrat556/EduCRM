package domain

import "fmt"

// PaymentType categorizes the payment record.
type PaymentType string

const (
	PaymentTypeMonthlyTuition PaymentType = "monthly_tuition"
	PaymentTypePartialPayment PaymentType = "partial_payment"
	PaymentTypeAdjustment     PaymentType = "adjustment"
	PaymentTypeOther          PaymentType = "other"
)

var validPaymentTypes = map[PaymentType]struct{}{
	PaymentTypeMonthlyTuition:  {},
	PaymentTypePartialPayment: {},
	PaymentTypeAdjustment:     {},
	PaymentTypeOther:          {},
}

// ParsePaymentType validates s.
func ParsePaymentType(s string) (PaymentType, error) {
	t := PaymentType(s)
	if _, ok := validPaymentTypes[t]; !ok {
		return "", fmt.Errorf("%w: %q", ErrInvalidPaymentType, s)
	}
	return t, nil
}
