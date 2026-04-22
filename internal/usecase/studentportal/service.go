package studentportal

import (
	"context"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	attendancesvc "github.com/educrm/educrm-backend/internal/usecase/attendance"
	gradesvc "github.com/educrm/educrm-backend/internal/usecase/grade"
	schedulesvc "github.com/educrm/educrm-backend/internal/usecase/schedule"
	"github.com/google/uuid"
)

// Service exposes student self-service reads (grades, schedule, attendance).
type Service struct {
	membership repository.StudentMembershipRepository
	schedule   *schedulesvc.Service
	grades     *gradesvc.Service
	attendance *attendancesvc.Service
}

// NewService returns a student portal service.
func NewService(
	membership repository.StudentMembershipRepository,
	schedule *schedulesvc.Service,
	grades *gradesvc.Service,
	attendance *attendancesvc.Service,
) *Service {
	return &Service{
		membership: membership,
		schedule:   schedule,
		grades:     grades,
		attendance: attendance,
	}
}

func (s *Service) groupIDForStudent(ctx context.Context, studentUserID uuid.UUID) (uuid.UUID, error) {
	gidPtr, err := s.membership.FindGroupIDByStudentUserID(ctx, studentUserID)
	if err != nil {
		return uuid.Nil, err
	}
	if gidPtr == nil || *gidPtr == uuid.Nil {
		return uuid.Nil, apperror.New(apperror.KindNotFound, "group_not_found", "student has no group assignment")
	}
	return *gidPtr, nil
}

// MyGrades returns grades for the authenticated student.
func (s *Service) MyGrades(ctx context.Context, studentUserID uuid.UUID) ([]domain.Grade, error) {
	return s.grades.List(ctx, domain.RoleStudent, studentUserID, gradesvc.ListFilter{
		StudentID: &studentUserID,
	})
}

// MySchedule returns schedule entries for the student's group.
func (s *Service) MySchedule(ctx context.Context, studentUserID uuid.UUID) ([]domain.Schedule, error) {
	gid, err := s.groupIDForStudent(ctx, studentUserID)
	if err != nil {
		return nil, err
	}
	return s.schedule.List(ctx, schedulesvc.ListFilter{GroupID: &gid})
}

// MyAttendance returns attendance rows for the authenticated student.
func (s *Service) MyAttendance(ctx context.Context, studentUserID uuid.UUID) ([]domain.Attendance, error) {
	return s.attendance.List(ctx, domain.RoleStudent, studentUserID, attendancesvc.ListFilter{
		StudentID: &studentUserID,
	})
}
