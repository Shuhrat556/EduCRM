package dto

import (
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/usecase/payment"
	"github.com/google/uuid"
)

// CreatePaymentRequest is the body for POST /payments.
type CreatePaymentRequest struct {
	StudentID           uuid.UUID `json:"student_id" example:"550e8400-e29b-41d4-a716-446655440000" binding:"required"`
	GroupID             uuid.UUID `json:"group_id" example:"6ba7b810-9dad-11d1-80b4-00c04fd430c8" binding:"required"`
	AmountMinor         int64     `json:"amount_minor" example:"150000" binding:"min=0"`
	Status              string    `json:"status" example:"paid_full" binding:"required,oneof=paid_full paid_partial unpaid overdue"`
	PaymentDate         *string   `json:"payment_date" example:"2025-04-01" binding:"omitempty"` // YYYY-MM-DD
	MonthFor            string    `json:"month_for" example:"2025-04" binding:"required"`        // YYYY-MM or YYYY-MM-DD
	PaymentType         string    `json:"payment_type" example:"monthly_tuition" binding:"required,oneof=monthly_tuition partial_payment adjustment other"`
	Comment             *string   `json:"comment" binding:"omitempty,max=4000"`
	IsFree              bool      `json:"is_free"`
	DiscountAmountMinor int64     `json:"discount_amount_minor" binding:"min=0"`
}

// UpdatePaymentRequest is the body for PATCH /payments/:id.
type UpdatePaymentRequest struct {
	AmountMinor         *int64  `json:"amount_minor" binding:"omitempty,min=0"`
	Status              *string `json:"status" binding:"omitempty,oneof=paid_full paid_partial unpaid overdue"`
	PaymentDate         *string `json:"payment_date" binding:"omitempty"`
	PaymentType         *string `json:"payment_type" binding:"omitempty,oneof=monthly_tuition partial_payment adjustment other"`
	Comment             *string `json:"comment" binding:"omitempty,max=4000"`
	IsFree              *bool   `json:"is_free"`
	DiscountAmountMinor *int64  `json:"discount_amount_minor" binding:"omitempty,min=0"`
}

// PaymentResponse is the API shape for a payment.
type PaymentResponse struct {
	ID                  uuid.UUID `json:"id"`
	StudentID           uuid.UUID `json:"student_id"`
	GroupID             uuid.UUID `json:"group_id"`
	AmountMinor         int64     `json:"amount_minor"`
	Status              string    `json:"status"`
	PaymentDate         *string   `json:"payment_date,omitempty"`
	MonthFor            string    `json:"month_for"`
	PaymentType         string    `json:"payment_type"`
	Comment             *string   `json:"comment,omitempty"`
	IsFree              bool      `json:"is_free"`
	DiscountAmountMinor int64     `json:"discount_amount_minor"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// PaymentListResponse paginated list.
type PaymentListResponse struct {
	Items    []PaymentResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

// PaymentResponseFrom maps domain to DTO.
func PaymentResponseFrom(p *domain.Payment) PaymentResponse {
	if p == nil {
		return PaymentResponse{}
	}
	var pay *string
	if p.PaymentDate != nil {
		s := p.PaymentDate.UTC().Format("2006-01-02")
		pay = &s
	}
	return PaymentResponse{
		ID:                  p.ID,
		StudentID:           p.StudentID,
		GroupID:             p.GroupID,
		AmountMinor:         p.AmountMinor,
		Status:              string(p.Status),
		PaymentDate:         pay,
		MonthFor:            p.MonthFor.UTC().Format("2006-01-02"),
		PaymentType:         string(p.PaymentType),
		Comment:             p.Comment,
		IsFree:              p.IsFree,
		DiscountAmountMinor: p.DiscountAmountMinor,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
}

// PaymentListResponseFrom maps list result.
func PaymentListResponseFrom(r *payment.ListResult) PaymentListResponse {
	if r == nil {
		return PaymentListResponse{}
	}
	items := make([]PaymentResponse, 0, len(r.Items))
	for i := range r.Items {
		items = append(items, PaymentResponseFrom(&r.Items[i]))
	}
	return PaymentListResponse{
		Items:    items,
		Total:    r.Total,
		Page:     r.Page,
		PageSize: r.PageSize,
	}
}
