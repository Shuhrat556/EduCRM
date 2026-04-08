package attendance

import (
	"context"
	"testing"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
)

type attStubRepo struct {
	createErr error
}

func (s *attStubRepo) Create(ctx context.Context, a *domain.Attendance) error { return s.createErr }
func (s *attStubRepo) Update(ctx context.Context, a *domain.Attendance) error { return nil }
func (s *attStubRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.Attendance, error) {
	return nil, nil
}
func (s *attStubRepo) ListByStudent(ctx context.Context, studentID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error) {
	return nil, nil
}
func (s *attStubRepo) ListByGroup(ctx context.Context, groupID uuid.UUID, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error) {
	return nil, nil
}
func (s *attStubRepo) ListByDateRange(ctx context.Context, from, to time.Time, viewerTeacherID *uuid.UUID) ([]domain.Attendance, error) {
	return nil, nil
}

type attStubGroups struct{ g *domain.Group }

func (s *attStubGroups) Create(ctx context.Context, g *domain.Group) error { return nil }
func (s *attStubGroups) Update(ctx context.Context, g *domain.Group) error { return nil }
func (s *attStubGroups) Delete(ctx context.Context, id uuid.UUID) error     { return nil }
func (s *attStubGroups) FindByID(ctx context.Context, id uuid.UUID) (*domain.Group, error) {
	return s.g, nil
}
func (s *attStubGroups) List(ctx context.Context, p repository.GroupListParams) ([]domain.Group, int64, error) {
	return nil, 0, nil
}

type attStubUsers struct{ u *domain.User }

func (s *attStubUsers) FindByLogin(ctx context.Context, login string) (*domain.User, error) { return nil, nil }
func (s *attStubUsers) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.u, nil
}
func (s *attStubUsers) Create(ctx context.Context, u *domain.User) error { return nil }
func (s *attStubUsers) Update(ctx context.Context, u *domain.User) error { return nil }
func (s *attStubUsers) Delete(ctx context.Context, id uuid.UUID) error   { return nil }
func (s *attStubUsers) List(ctx context.Context, p repository.UserListParams) ([]domain.User, int64, error) {
	return nil, 0, nil
}
func (s *attStubUsers) EmailTaken(ctx context.Context, email string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}
func (s *attStubUsers) PhoneTaken(ctx context.Context, phone string, excludeID *uuid.UUID) (bool, error) {
	return false, nil
}

type attStubMem struct{ gid uuid.UUID }

func (s *attStubMem) FindGroupIDByStudentUserID(ctx context.Context, studentUserID uuid.UUID) (*uuid.UUID, error) {
	return &s.gid, nil
}

type attStubLinks struct{}

func (s *attStubLinks) FindTeacherIDByUserID(ctx context.Context, userID uuid.UUID) (*uuid.UUID, error) {
	return nil, nil
}

func TestCreate_duplicateReturnsConflict(t *testing.T) {
	ctx := context.Background()
	studentID, groupID, teacherID := uuid.New(), uuid.New(), uuid.New()
	svc := NewService(
		&attStubRepo{createErr: repository.ErrDuplicate},
		&attStubGroups{g: &domain.Group{ID: groupID, TeacherID: teacherID}},
		&attStubUsers{u: &domain.User{ID: studentID, Role: domain.RoleStudent}},
		&attStubMem{gid: groupID},
		&attStubLinks{},
	)
	_, err := svc.Create(ctx, domain.RoleAdmin, uuid.New(), CreateInput{
		StudentID:  studentID,
		GroupID:    groupID,
		LessonDate: time.Date(2025, 4, 8, 0, 0, 0, 0, time.UTC),
		Status:     domain.AttendancePresent,
	})
	if err == nil {
		t.Fatal("expected conflict")
	}
	ae, ok := apperror.AsError(err)
	if !ok || ae.Kind != apperror.KindConflict || ae.Code != "attendance_exists" {
		t.Fatalf("got %+v", ae)
	}
}

func TestList_filterExactlyOne(t *testing.T) {
	ctx := context.Background()
	svc := NewService(&attStubRepo{}, &attStubGroups{}, &attStubUsers{}, &attStubMem{}, &attStubLinks{})
	sid := uuid.New()
	tests := []struct {
		name    string
		f       ListFilter
		wantErr bool
	}{
		{"none", ListFilter{}, true},
		{"student+group", ListFilter{StudentID: &sid, GroupID: uuidPtr(uuid.New())}, true},
		{"partial range", ListFilter{From: timePtr(time.Now())}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.List(ctx, domain.RoleAdmin, uuid.New(), tt.f)
			if tt.wantErr && err == nil {
				t.Fatal("expected error")
			}
			if !tt.wantErr && err != nil {
				t.Fatal(err)
			}
		})
	}
}

func uuidPtr(u uuid.UUID) *uuid.UUID { return &u }
func timePtr(t time.Time) *time.Time { return &t }
