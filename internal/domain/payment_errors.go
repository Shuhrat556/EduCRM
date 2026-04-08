package domain

import "errors"

var (
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
	ErrInvalidPaymentType   = errors.New("invalid payment_type")
)
