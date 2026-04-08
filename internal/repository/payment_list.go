package repository

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// PaymentListParams filters paginated payment listing (staff).
type PaymentListParams struct {
	Search      string
	StudentID   *uuid.UUID
	GroupID     *uuid.UUID
	MonthFor    *time.Time // first day of month (UTC)
	Status      *domain.PaymentStatus
	PaymentType *domain.PaymentType
	IsFree      *bool
	Page        int
	PageSize    int
}
