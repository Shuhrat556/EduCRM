package repository

import (
	"context"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/google/uuid"
)

// PaymentRepository persists student payments (soft-delete preserves history).
type PaymentRepository interface {
	Create(ctx context.Context, p *domain.Payment) error
	Update(ctx context.Context, p *domain.Payment) error
	Delete(ctx context.Context, id uuid.UUID) error // soft delete
	FindByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error)
	List(ctx context.Context, p PaymentListParams) ([]domain.Payment, int64, error)
	// ListHistoryByStudent returns non-deleted payments for one student, newest first.
	ListHistoryByStudent(ctx context.Context, studentID uuid.UUID, page, pageSize int) ([]domain.Payment, int64, error)
}
