package domain

import "errors"

// ErrInvalidNotificationType is returned for unknown notification type strings.
var ErrInvalidNotificationType = errors.New("invalid notification type")
