package payment

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service orchestrates payment use cases.
type Service struct {
	payments    repository.PaymentRepository
	groups      repository.GroupRepository
	users       repository.UserRepository
	memberships repository.StudentMembershipRepository
}

// NewService constructs a payment service.
func NewService(
	payments repository.PaymentRepository,
	groups repository.GroupRepository,
	users repository.UserRepository,
	memberships repository.StudentMembershipRepository,
) *Service {
	return &Service{
		payments:    payments,
		groups:      groups,
		users:       users,
		memberships: memberships,
	}
}

// CreateInput holds create payload.
type CreateInput struct {
	StudentID           uuid.UUID
	GroupID             uuid.UUID
	AmountMinor         int64
	Status              domain.PaymentStatus
	PaymentDate         *time.Time
	MonthFor            time.Time
	PaymentType         domain.PaymentType
	Comment             *string
	IsFree              bool
	DiscountAmountMinor int64
}

// UpdateInput holds optional updates.
type UpdateInput struct {
	AmountMinor         *int64
	Status              *domain.PaymentStatus
	PaymentDate         *time.Time
	PaymentType         *domain.PaymentType
	Comment             *string
	IsFree              *bool
	DiscountAmountMinor *int64
}

// ListResult is paginated list output.
type ListResult struct {
	Items    []domain.Payment
	Total    int64
	Page     int
	PageSize int
}

// Create validates and inserts a payment.
func (s *Service) Create(ctx context.Context, actorRole domain.Role, in CreateInput) (*domain.Payment, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	if err := s.assertPrivilegedFields(actorRole, in.IsFree, in.DiscountAmountMinor); err != nil {
		return nil, err
	}
	if in.AmountMinor < 0 || in.DiscountAmountMinor < 0 {
		return nil, apperror.Validation("amount", "amount_minor and discount_amount_minor must be >= 0")
	}
	if err := s.assertStudentInGroup(ctx, in.StudentID, in.GroupID); err != nil {
		return nil, err
	}
	monthFor := domain.MonthStartUTC(in.MonthFor)
	now := time.Now().UTC()
	row := &domain.Payment{
		ID:                  uuid.New(),
		StudentID:           in.StudentID,
		GroupID:             in.GroupID,
		AmountMinor:         in.AmountMinor,
		Status:              in.Status,
		PaymentDate:         cloneDatePtr(in.PaymentDate),
		MonthFor:            monthFor,
		PaymentType:         in.PaymentType,
		Comment:             trimComment(in.Comment),
		IsFree:              in.IsFree,
		DiscountAmountMinor: in.DiscountAmountMinor,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if err := s.payments.Create(ctx, row); err != nil {
		return nil, apperror.Internal("create payment").Wrap(err)
	}
	return row, nil
}

// GetByID returns a payment visible to staff or owning student.
func (s *Service) GetByID(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, id uuid.UUID) (*domain.Payment, error) {
	row, err := s.payments.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load payment").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("payment")
	}
	if err := s.ensureCanView(actorRole, actorUserID, row); err != nil {
		return nil, err
	}
	return row, nil
}

// List paginates for staff.
func (s *Service) List(ctx context.Context, actorRole domain.Role, p repository.PaymentListParams) (*ListResult, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	rows, total, err := s.payments.List(ctx, p)
	if err != nil {
		return nil, apperror.Internal("list payments").Wrap(err)
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
	return &ListResult{Items: rows, Total: total, Page: page, PageSize: size}, nil
}

// History returns paginated payment history for one student (student self or staff with student_id).
func (s *Service) History(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, studentID *uuid.UUID, page, pageSize int) (*ListResult, error) {
	var target uuid.UUID
	switch actorRole {
	case domain.RoleStudent:
		if studentID != nil && *studentID != actorUserID {
			return nil, apperror.New(apperror.KindForbidden, "cannot_view_history", "you may only view your own payment history")
		}
		target = actorUserID
	case domain.RoleAdmin, domain.RoleSuperAdmin:
		if studentID == nil {
			return nil, apperror.Validation("student_id", "required for staff")
		}
		target = *studentID
	default:
		return nil, apperror.New(apperror.KindForbidden, "insufficient_permissions", "only admin or student may view payment history")
	}
	rows, total, err := s.payments.ListHistoryByStudent(ctx, target, page, pageSize)
	if err != nil {
		return nil, apperror.Internal("list payment history").Wrap(err)
	}
	pg := page
	if pg < 1 {
		pg = 1
	}
	sz := pageSize
	if sz < 1 {
		sz = 20
	}
	if sz > 100 {
		sz = 100
	}
	return &ListResult{Items: rows, Total: total, Page: pg, PageSize: sz}, nil
}

// Update applies changes (staff only).
func (s *Service) Update(ctx context.Context, actorRole domain.Role, id uuid.UUID, in UpdateInput) (*domain.Payment, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	row, err := s.payments.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load payment").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("payment")
	}
	if in.IsFree != nil && *in.IsFree {
		if actorRole != domain.RoleSuperAdmin {
			return nil, apperror.New(apperror.KindForbidden, "free_requires_super_admin", "only super_admin may mark a payment as free")
		}
	}
	if in.DiscountAmountMinor != nil && *in.DiscountAmountMinor > 0 {
		if actorRole != domain.RoleSuperAdmin {
			return nil, apperror.New(apperror.KindForbidden, "discount_requires_super_admin", "only super_admin may set a discount amount")
		}
	}
	if in.DiscountAmountMinor != nil && *in.DiscountAmountMinor < 0 {
		return nil, apperror.Validation("discount_amount_minor", "must be >= 0")
	}
	if in.AmountMinor != nil && *in.AmountMinor < 0 {
		return nil, apperror.Validation("amount_minor", "must be >= 0")
	}
	if in.AmountMinor != nil {
		row.AmountMinor = *in.AmountMinor
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if in.PaymentDate != nil {
		row.PaymentDate = cloneDatePtr(in.PaymentDate)
	}
	if in.PaymentType != nil {
		row.PaymentType = *in.PaymentType
	}
	if in.Comment != nil {
		row.Comment = trimComment(in.Comment)
	}
	if in.IsFree != nil {
		row.IsFree = *in.IsFree
	}
	if in.DiscountAmountMinor != nil {
		row.DiscountAmountMinor = *in.DiscountAmountMinor
	}
	row.UpdatedAt = time.Now().UTC()
	if err := s.payments.Update(ctx, row); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("payment")
		}
		return nil, apperror.Internal("update payment").Wrap(err)
	}
	return row, nil
}

// Delete soft-deletes (staff only).
func (s *Service) Delete(ctx context.Context, actorRole domain.Role, id uuid.UUID) error {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return err
	}
	if err := s.payments.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("payment")
		}
		return apperror.Internal("delete payment").Wrap(err)
	}
	return nil
}

func (s *Service) assertPrivilegedFields(actorRole domain.Role, isFree bool, discountMinor int64) error {
	if isFree && actorRole != domain.RoleSuperAdmin {
		return apperror.New(apperror.KindForbidden, "free_requires_super_admin", "only super_admin may mark a payment as free")
	}
	if discountMinor > 0 && actorRole != domain.RoleSuperAdmin {
		return apperror.New(apperror.KindForbidden, "discount_requires_super_admin", "only super_admin may set a discount amount")
	}
	return nil
}

func (s *Service) assertStudentInGroup(ctx context.Context, studentID, groupID uuid.UUID) error {
	student, err := s.users.FindByID(ctx, studentID)
	if err != nil {
		return apperror.Internal("load student").Wrap(err)
	}
	if student == nil {
		return apperror.Validation("student_id", "user not found")
	}
	if student.Role != domain.RoleStudent {
		return apperror.Validation("student_id", "user must have role student")
	}
	g, err := s.groups.FindByID(ctx, groupID)
	if err != nil {
		return apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return apperror.Validation("group_id", "group not found")
	}
	mem, err := s.memberships.FindGroupIDByStudentUserID(ctx, studentID)
	if err != nil {
		return apperror.Internal("load enrollment").Wrap(err)
	}
	if mem == nil || *mem != groupID {
		return apperror.Validation("student_id", "student is not enrolled in this group")
	}
	return nil
}

func (s *Service) ensureCanView(actorRole domain.Role, actorUserID uuid.UUID, row *domain.Payment) error {
	switch actorRole {
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		return nil
	case domain.RoleStudent:
		if row.StudentID != actorUserID {
			return apperror.New(apperror.KindForbidden, "cannot_view_payment", "you may only view your own payments")
		}
		return nil
	default:
		return apperror.New(apperror.KindForbidden, "insufficient_permissions", "cannot view payment")
	}
}

func trimComment(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	if t == "" {
		return nil
	}
	return &t
}

func cloneDatePtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	x := t.UTC()
	d := time.Date(x.Year(), x.Month(), x.Day(), 0, 0, 0, 0, time.UTC)
	return &d
}
