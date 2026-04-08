package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/educrm/educrm-backend/internal/repository/postgres/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentRepository implements repository.PaymentRepository.
type PaymentRepository struct {
	db *gorm.DB
}

var _ repository.PaymentRepository = (*PaymentRepository)(nil)

// NewPaymentRepository constructs PaymentRepository.
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

// Create inserts a payment.
func (r *PaymentRepository) Create(ctx context.Context, p *domain.Payment) error {
	m, err := paymentToModel(p)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(m).Error
}

// Update updates mutable scalar fields (not student/group/month_for).
func (r *PaymentRepository) Update(ctx context.Context, p *domain.Payment) error {
	m, err := paymentToModel(p)
	if err != nil {
		return err
	}
	res := r.db.WithContext(ctx).Model(&model.Payment{}).Where("id = ?", p.ID).Updates(map[string]any{
		"amount_minor":          m.AmountMinor,
		"status":                m.Status,
		"payment_date":          m.PaymentDate,
		"payment_type":          m.PaymentType,
		"comment":               m.Comment,
		"is_free":               m.IsFree,
		"discount_amount_minor": m.DiscountAmountMinor,
		"updated_at":            m.UpdatedAt,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Delete soft-deletes a payment.
func (r *PaymentRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res := r.db.WithContext(ctx).Delete(&model.Payment{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindByID loads a non-deleted payment.
func (r *PaymentRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	var m model.Payment
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return paymentToDomain(&m)
}

// List paginates filtered payments (excludes soft-deleted).
func (r *PaymentRepository) List(ctx context.Context, p repository.PaymentListParams) ([]domain.Payment, int64, error) {
	build := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&model.Payment{})
		if p.Search != "" {
			term := "%" + escapeLikePattern(p.Search) + "%"
			q = q.Where("payments.comment ILIKE ? ESCAPE '\\'", term)
		}
		if p.StudentID != nil {
			q = q.Where("payments.student_id = ?", *p.StudentID)
		}
		if p.GroupID != nil {
			q = q.Where("payments.group_id = ?", *p.GroupID)
		}
		if p.MonthFor != nil {
			m := truncateUTCDate(*p.MonthFor)
			q = q.Where("payments.month_for = ?", m)
		}
		if p.Status != nil {
			q = q.Where("payments.status = ?", string(*p.Status))
		}
		if p.PaymentType != nil {
			q = q.Where("payments.payment_type = ?", string(*p.PaymentType))
		}
		if p.IsFree != nil {
			q = q.Where("payments.is_free = ?", *p.IsFree)
		}
		return q
	}
	var total int64
	if err := build().Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page := p.Page
	if page < 1 {
		page = 1
	}
	size := p.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size
	var rows []model.Payment
	if err := build().Order("payments.month_for DESC, payments.created_at DESC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out, err := paymentsToDomain(rows)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

// ListHistoryByStudent lists payments for a student (newest first).
func (r *PaymentRepository) ListHistoryByStudent(ctx context.Context, studentID uuid.UUID, page, pageSize int) ([]domain.Payment, int64, error) {
	base := r.db.WithContext(ctx).Model(&model.Payment{}).Where("student_id = ?", studentID)
	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page < 1 {
		page = 1
	}
	size := pageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	offset := (page - 1) * size
	var rows []model.Payment
	if err := base.Order("month_for DESC, created_at DESC").Limit(size).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out, err := paymentsToDomain(rows)
	if err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func paymentToDomain(m *model.Payment) (*domain.Payment, error) {
	st, err := domain.ParsePaymentStatus(m.Status)
	if err != nil {
		return nil, err
	}
	pt, err := domain.ParsePaymentType(m.PaymentType)
	if err != nil {
		return nil, err
	}
	var payDate *time.Time
	if m.PaymentDate != nil {
		d := truncateUTCDate(*m.PaymentDate)
		payDate = &d
	}
	return &domain.Payment{
		ID:                  m.ID,
		StudentID:           m.StudentID,
		GroupID:             m.GroupID,
		AmountMinor:         m.AmountMinor,
		Status:              st,
		PaymentDate:         payDate,
		MonthFor:            domain.MonthStartUTC(m.MonthFor),
		PaymentType:         pt,
		Comment:             m.Comment,
		IsFree:              m.IsFree,
		DiscountAmountMinor: m.DiscountAmountMinor,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
	}, nil
}

func paymentToModel(p *domain.Payment) (*model.Payment, error) {
	m := &model.Payment{
		ID:                  p.ID,
		StudentID:           p.StudentID,
		GroupID:             p.GroupID,
		AmountMinor:         p.AmountMinor,
		Status:              string(p.Status),
		MonthFor:            domain.MonthStartUTC(p.MonthFor),
		PaymentType:         string(p.PaymentType),
		Comment:             p.Comment,
		IsFree:              p.IsFree,
		DiscountAmountMinor: p.DiscountAmountMinor,
		CreatedAt:           p.CreatedAt,
		UpdatedAt:           p.UpdatedAt,
	}
	if p.PaymentDate != nil {
		d := truncateUTCDate(*p.PaymentDate)
		m.PaymentDate = &d
	}
	return m, nil
}

func paymentsToDomain(rows []model.Payment) ([]domain.Payment, error) {
	out := make([]domain.Payment, 0, len(rows))
	for i := range rows {
		p, err := paymentToDomain(&rows[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *p)
	}
	return out, nil
}
