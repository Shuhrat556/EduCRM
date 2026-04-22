package payment

import (
	"context"
	"testing"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
)

type payStubUsers struct {
	byID map[uuid.UUID]*domain.User
}

func (s *payStubUsers) FindByLogin(ctx context.Context, login string) (*domain.User, error) {
	return nil, nil
}
func (s *payStubUsers) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}
func (s *payStubUsers) Create(ctx context.Context, u *domain.User) error { return nil }
func (s *payStubUsers) Update(ctx context.Context, u *domain.User) error { return nil }
func (s *payStubUsers) Delete(ctx context.Context, id uuid.UUID) error   { return nil }
func (s *payStubUsers) List(ctx context.Context, p repository.UserListParams) ([]domain.User, int64, error) {
	return nil, 0, nil
}
func (s *payStubUsers) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *payStubUsers) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *payStubUsers) UsernameTaken(ctx context.Context, username string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}

type payStubGroups struct {
	byID map[uuid.UUID]*domain.Group
}

func (s *payStubGroups) Create(ctx context.Context, g *domain.Group) error { return nil }
func (s *payStubGroups) Update(ctx context.Context, g *domain.Group) error { return nil }
func (s *payStubGroups) Delete(ctx context.Context, id uuid.UUID) error     { return nil }
func (s *payStubGroups) FindByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	if s.byID == nil {
		return nil, nil
	}
	return s.byID[id], nil
}
func (s *payStubGroups) List(ctx context.Context, p repository.GroupListParams) ([]domain.Group, int64, error) {
	return nil, 0, nil
}

type payStubMem struct {
	groupByStudent map[uuid.UUID]uuid.UUID
}

func (s *payStubMem) FindGroupIDByStudentUserID(ctx context.Context, studentUserID uuid.UUID) (*uuid.UUID, error) {
	gid, ok := s.groupByStudent[studentUserID]
	if !ok {
		return nil, nil
	}
	return &gid, nil
}
func (s *payStubMem) ListStudentUserIDsByGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	var ids []uuid.UUID
	for sid, gid := range s.groupByStudent {
		if gid == groupID {
			ids = append(ids, sid)
		}
	}
	return ids, nil
}

type payStubRepo struct {
	lastCreate *domain.Payment
	row        *domain.Payment
}

func (s *payStubRepo) Create(ctx context.Context, p *domain.Payment) error {
	s.lastCreate = p
	return nil
}
func (s *payStubRepo) Update(ctx context.Context, p *domain.Payment) error {
	s.row = p
	return nil
}
func (s *payStubRepo) Delete(ctx context.Context, id uuid.UUID) error { return nil }
func (s *payStubRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Payment, error) {
	if s.row == nil {
		return nil, nil
	}
	return s.row, nil
}
func (s *payStubRepo) List(ctx context.Context, p repository.PaymentListParams) ([]domain.Payment, int64, error) {
	return nil, 0, nil
}
func (s *payStubRepo) ListHistoryByStudent(ctx context.Context, studentID uuid.UUID, page, pageSize int) ([]domain.Payment, int64, error) {
	return nil, 0, nil
}

func basePaymentDeps(studentID, groupID uuid.UUID) (*payStubRepo, *payStubGroups, *payStubUsers, *payStubMem) {
	pay := &payStubRepo{}
	grp := &payStubGroups{
		byID: map[uuid.UUID]*domain.Group{
			groupID: {ID: groupID, TeacherID: uuid.New()},
		},
	}
	users := &payStubUsers{
		byID: map[uuid.UUID]*domain.User{
			studentID: {ID: studentID, Role: domain.RoleStudent, IsActive: true},
		},
	}
	mem := &payStubMem{groupByStudent: map[uuid.UUID]uuid.UUID{studentID: groupID}}
	return pay, grp, users, mem
}

func TestCreate_privilegedFieldsAndStaff(t *testing.T) {
	ctx := context.Background()
	studentID, groupID := uuid.New(), uuid.New()
	month := time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC)
	base := CreateInput{
		StudentID:           studentID,
		GroupID:             groupID,
		AmountMinor:         5000,
		Status:              domain.PaymentStatusUnpaid,
		MonthFor:            month,
		PaymentType:         domain.PaymentTypeMonthlyTuition,
		IsFree:              false,
		DiscountAmountMinor: 0,
	}

	tests := []struct {
		name     string
		role     domain.Role
		in       CreateInput
		wantKind apperror.Kind
		wantCode string
		wantOK   bool
	}{
		{
			name:     "teacher rejected",
			role:     domain.RoleTeacher,
			in:       base,
			wantKind: apperror.KindForbidden,
		},
		{
			name: "admin cannot set free",
			role: domain.RoleAdmin,
			in: func() CreateInput {
				x := base
				x.IsFree = true
				return x
			}(),
			wantKind: apperror.KindForbidden,
			wantCode: "free_requires_super_admin",
		},
		{
			name: "admin cannot set discount",
			role: domain.RoleAdmin,
			in: func() CreateInput {
				x := base
				x.DiscountAmountMinor = 100
				return x
			}(),
			wantKind: apperror.KindForbidden,
			wantCode: "discount_requires_super_admin",
		},
		{
			name: "super may set free and discount",
			role: domain.RoleSuperAdmin,
			in: func() CreateInput {
				x := base
				x.IsFree = true
				x.DiscountAmountMinor = 50
				x.AmountMinor = 0
				return x
			}(),
			wantOK: true,
		},
		{
			name: "admin normal row",
			role: domain.RoleAdmin,
			in:   base,
			wantOK: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pay, grp, users, mem := basePaymentDeps(studentID, groupID)
			svc := NewService(pay, grp, users, mem)
			got, err := svc.Create(ctx, tt.role, tt.in)
			if tt.wantOK {
				if err != nil {
					t.Fatal(err)
				}
				if got == nil || pay.lastCreate == nil {
					t.Fatal("expected persisted payment")
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			ae, ok := apperror.AsError(err)
			if !ok {
				t.Fatalf("unexpected %v", err)
			}
			if ae.Kind != tt.wantKind {
				t.Fatalf("kind %s, want %s", ae.Kind, tt.wantKind)
			}
			if tt.wantCode != "" && ae.Code != tt.wantCode {
				t.Fatalf("code %q, want %q", ae.Code, tt.wantCode)
			}
		})
	}
}

func TestCreate_negativeAmounts(t *testing.T) {
	ctx := context.Background()
	studentID, groupID := uuid.New(), uuid.New()
	pay, grp, users, mem := basePaymentDeps(studentID, groupID)
	svc := NewService(pay, grp, users, mem)
	in := CreateInput{
		StudentID: studentID, GroupID: groupID,
		AmountMinor: -1, Status: domain.PaymentStatusUnpaid,
		MonthFor: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		PaymentType: domain.PaymentTypeMonthlyTuition,
	}
	_, err := svc.Create(ctx, domain.RoleAdmin, in)
	if err == nil {
		t.Fatal("expected validation error")
	}
	ae, _ := apperror.AsError(err)
	if ae.Kind != apperror.KindValidation {
		t.Fatalf("got %v", ae.Kind)
	}
}

func TestGetByID_studentCannotViewOthers(t *testing.T) {
	ctx := context.Background()
	studentID, otherStudent, groupID := uuid.New(), uuid.New(), uuid.New()
	pay, grp, users, mem := basePaymentDeps(studentID, groupID)
	users.byID[otherStudent] = &domain.User{ID: otherStudent, Role: domain.RoleStudent}
	pid := uuid.New()
	pay.row = &domain.Payment{ID: pid, StudentID: otherStudent, GroupID: groupID}
	svc := NewService(pay, grp, users, mem)
	_, err := svc.GetByID(ctx, domain.RoleStudent, studentID, pid)
	if err == nil {
		t.Fatal("expected forbidden")
	}
	ae, _ := apperror.AsError(err)
	if ae.Kind != apperror.KindForbidden || ae.Code != "cannot_view_payment" {
		t.Fatalf("got %+v", ae)
	}
}

func TestHistory_rules(t *testing.T) {
	ctx := context.Background()
	selfID, other := uuid.New(), uuid.New()
	pay := &payStubRepo{}
	svc := NewService(pay, &payStubGroups{}, &payStubUsers{}, &payStubMem{})

	tests := []struct {
		name     string
		role     domain.Role
		actor    uuid.UUID
		student  *uuid.UUID
		wantKind apperror.Kind
		wantOK   bool
	}{
		{
			name:     "student cannot pass different student_id",
			role:     domain.RoleStudent,
			actor:    selfID,
			student:  &other,
			wantKind: apperror.KindForbidden,
		},
		{
			name:    "student omits student_id uses self",
			role:    domain.RoleStudent,
			actor:   selfID,
			student: nil,
			wantOK:  true,
		},
		{
			name:     "admin requires student_id",
			role:     domain.RoleAdmin,
			actor:    uuid.New(),
			student:  nil,
			wantKind: apperror.KindValidation,
		},
		{
			name:     "teacher forbidden",
			role:     domain.RoleTeacher,
			actor:    uuid.New(),
			student:  &selfID,
			wantKind: apperror.KindForbidden,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.History(ctx, tt.role, tt.actor, tt.student, 1, 20)
			if tt.wantOK {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil {
				t.Fatal("expected error")
			}
			ae, ok := apperror.AsError(err)
			if !ok || ae.Kind != tt.wantKind {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestUpdate_freeDiscountSuperOnly(t *testing.T) {
	ctx := context.Background()
	pid := uuid.New()
	row := &domain.Payment{
		ID: pid, StudentID: uuid.New(), GroupID: uuid.New(),
		AmountMinor: 1000, Status: domain.PaymentStatusUnpaid,
		MonthFor: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		PaymentType: domain.PaymentTypeMonthlyTuition,
	}
	pay := &payStubRepo{row: row}
	svc := NewService(pay, &payStubGroups{}, &payStubUsers{}, &payStubMem{})

	t.Run("admin cannot set is_free", func(t *testing.T) {
		v := true
		_, err := svc.Update(ctx, domain.RoleAdmin, pid, UpdateInput{IsFree: &v})
		if err == nil {
			t.Fatal("expected error")
		}
		ae, _ := apperror.AsError(err)
		if ae.Code != "free_requires_super_admin" {
			t.Fatalf("code %q", ae.Code)
		}
	})
	t.Run("super can set is_free", func(t *testing.T) {
		v := true
		got, err := svc.Update(ctx, domain.RoleSuperAdmin, pid, UpdateInput{IsFree: &v})
		if err != nil {
			t.Fatal(err)
		}
		if !got.IsFree {
			t.Fatal("expected free")
		}
	})
}
