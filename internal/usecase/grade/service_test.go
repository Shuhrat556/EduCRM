package grade

import (
	"context"
	"math"
	"testing"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
)

type gradeStubRepo struct {
	createErr error
}

func (s *gradeStubRepo) Create(ctx context.Context, g *domain.Grade) error { return s.createErr }
func (s *gradeStubRepo) Update(ctx context.Context, g *domain.Grade) error { return nil }
func (s *gradeStubRepo) Delete(ctx context.Context, id uuid.UUID) error     { return nil }
func (s *gradeStubRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Grade, error) {
	return nil, nil
}
func (s *gradeStubRepo) ListByStudent(ctx context.Context, studentID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Grade, error) {
	return nil, nil
}
func (s *gradeStubRepo) ListByGroup(ctx context.Context, groupID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Grade, error) {
	return nil, nil
}

type gradeStubGroups struct{ g *domain.Group }

func (s *gradeStubGroups) Create(ctx context.Context, g *domain.Group) error { return nil }
func (s *gradeStubGroups) Update(ctx context.Context, g *domain.Group) error { return nil }
func (s *gradeStubGroups) Delete(ctx context.Context, id uuid.UUID) error     { return nil }
func (s *gradeStubGroups) FindByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	return s.g, nil
}
func (s *gradeStubGroups) List(ctx context.Context, p repository.GroupListParams) ([]domain.Group, int64, error) {
	return nil, 0, nil
}

type gradeStubUsers struct{ u *domain.User }

func (s *gradeStubUsers) FindByLogin(ctx context.Context, login string) (*domain.User, error) { return nil, nil }
func (s *gradeStubUsers) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.u, nil
}
func (s *gradeStubUsers) Create(ctx context.Context, u *domain.User) error { return nil }
func (s *gradeStubUsers) Update(ctx context.Context, u *domain.User) error { return nil }
func (s *gradeStubUsers) Delete(ctx context.Context, id uuid.UUID) error   { return nil }
func (s *gradeStubUsers) List(ctx context.Context, p repository.UserListParams) ([]domain.User, int64, error) {
	return nil, 0, nil
}
func (s *gradeStubUsers) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *gradeStubUsers) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *gradeStubUsers) UsernameTaken(ctx context.Context, username string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}

type gradeStubAssign struct{}

func (gradeStubAssign) Exists(ctx context.Context, teacherID, groupID, subjectID uuid.UUID) (bool, error) {
	return true, nil
}
func (gradeStubAssign) HasAnyAssignmentOnGroup(ctx context.Context, teacherID, groupID uuid.UUID) (bool, error) {
	return true, nil
}
func (gradeStubAssign) ListByTeacher(ctx context.Context, teacherID uuid.UUID) ([]repository.TeacherAssignmentRow, error) {
	return nil, nil
}

type gradeStubMem struct{ gid uuid.UUID }

func (s *gradeStubMem) FindGroupIDByStudentUserID(ctx context.Context, studentUserID uuid.UUID) (*uuid.UUID, error) {
	return &s.gid, nil
}
func (s *gradeStubMem) ListStudentUserIDsByGroup(ctx context.Context, groupID uuid.UUID) ([]uuid.UUID, error) {
	return nil, nil
}

type gradeStubLinks struct{}

func (s *gradeStubLinks) FindTeacherIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	return nil, nil
}

func TestCreate_duplicateReturnsConflict(t *testing.T) {
	ctx := context.Background()
	studentID, groupID, teacherID := uuid.New(), uuid.New(), uuid.New()
	svc := NewService(
		&gradeStubRepo{createErr: repository.ErrDuplicate},
		gradeStubAssign{},
		&gradeStubGroups{g: &domain.Group{ID: groupID, TeacherID: teacherID}},
		&gradeStubUsers{u: &domain.User{ID: studentID, Role: domain.RoleStudent}},
		&gradeStubMem{gid: groupID},
		&gradeStubLinks{},
	)
	_, err := svc.Create(ctx, domain.RoleSuperAdmin, uuid.New(), CreateInput{
		StudentID:  studentID,
		GroupID:    groupID,
		GradeType:  domain.GradeTypeTeacherEvaluation,
		GradeValue: 4.5,
	})
	if err == nil {
		t.Fatal("expected conflict")
	}
	ae, ok := apperror.AsError(err)
	if !ok || ae.Kind != apperror.KindConflict || ae.Code != "grade_week_exists" {
		t.Fatalf("got %+v", ae)
	}
}

func TestCreate_invalidGradeValue(t *testing.T) {
	ctx := context.Background()
	studentID, groupID, teacherID := uuid.New(), uuid.New(), uuid.New()
	svc := NewService(
		&gradeStubRepo{},
		gradeStubAssign{},
		&gradeStubGroups{g: &domain.Group{ID: groupID, TeacherID: teacherID}},
		&gradeStubUsers{u: &domain.User{ID: studentID, Role: domain.RoleStudent}},
		&gradeStubMem{gid: groupID},
		&gradeStubLinks{},
	)
	_, err := svc.Create(ctx, domain.RoleAdmin, uuid.New(), CreateInput{
		StudentID:  studentID,
		GroupID:    groupID,
		GradeType:  domain.GradeTypeTeacherEvaluation,
		GradeValue: math.NaN(),
	})
	if err == nil {
		t.Fatal("expected validation")
	}
	ae, _ := apperror.AsError(err)
	if ae.Kind != apperror.KindValidation {
		t.Fatalf("got %v", ae.Kind)
	}
}

func TestList_studentCannotUseGroupFilter(t *testing.T) {
	ctx := context.Background()
	svc := NewService(&gradeStubRepo{}, gradeStubAssign{}, &gradeStubGroups{}, &gradeStubUsers{}, &gradeStubMem{}, &gradeStubLinks{})
	gid := uuid.New()
	_, err := svc.List(ctx, domain.RoleStudent, uuid.New(), ListFilter{GroupID: &gid})
	if err == nil {
		t.Fatal("expected forbidden")
	}
	ae, _ := apperror.AsError(err)
	if ae.Code != "cannot_list_group_grades" {
		t.Fatalf("code %q", ae.Code)
	}
}
