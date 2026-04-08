package attendance

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service orchestrates attendance use cases.
type Service struct {
	attendance  repository.AttendanceRepository
	groups      repository.GroupRepository
	users       repository.UserRepository
	memberships repository.StudentMembershipRepository
	links       repository.UserTeacherLinkRepository
}

// NewService constructs an attendance service.
func NewService(
	attendance repository.AttendanceRepository,
	groups repository.GroupRepository,
	users repository.UserRepository,
	memberships repository.StudentMembershipRepository,
	links repository.UserTeacherLinkRepository,
) *Service {
	return &Service{
		attendance:  attendance,
		groups:      groups,
		users:       users,
		memberships: memberships,
		links:       links,
	}
}

// CreateInput is validated input for marking attendance.
type CreateInput struct {
	StudentID  uuid.UUID
	GroupID    uuid.UUID
	LessonDate time.Time
	Status     domain.AttendanceStatus
	Comment    *string
}

// UpdateInput holds optional field updates for PATCH.
type UpdateInput struct {
	Status  *domain.AttendanceStatus
	Comment *string
}

// ListFilter selects list mode (exactly one branch used by handler).
type ListFilter struct {
	StudentID *uuid.UUID
	GroupID   *uuid.UUID
	From      *time.Time
	To        *time.Time
}

// Create records attendance for one student on a lesson date.
func (s *Service) Create(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, in CreateInput) (*domain.Attendance, error) {
	filter, err := s.teacherFilter(ctx, actorRole, actorUserID)
	if err != nil {
		return nil, err
	}
	g, err := s.groups.FindByID(ctx, in.GroupID)
	if err != nil {
		return nil, apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return nil, apperror.Validation("group_id", "group not found")
	}
	if err := s.ensureCanModifyGroup(actorRole, filter, g.TeacherID); err != nil {
		return nil, err
	}
	student, err := s.users.FindByID(ctx, in.StudentID)
	if err != nil {
		return nil, apperror.Internal("load student").Wrap(err)
	}
	if student == nil {
		return nil, apperror.Validation("student_id", "user not found")
	}
	if student.Role != domain.RoleStudent {
		return nil, apperror.Validation("student_id", "user must have role student")
	}
	memGroup, err := s.memberships.FindGroupIDByStudentUserID(ctx, in.StudentID)
	if err != nil {
		return nil, apperror.Internal("load enrollment").Wrap(err)
	}
	if memGroup == nil || *memGroup != in.GroupID {
		return nil, apperror.Validation("student_id", "student is not enrolled in this group")
	}
	lessonDate := truncateUTCDate(in.LessonDate)
	now := time.Now().UTC()
	row := &domain.Attendance{
		ID:                uuid.New(),
		StudentID:         in.StudentID,
		GroupID:           in.GroupID,
		LessonDate:        lessonDate,
		Status:            in.Status,
		Comment:           trimComment(in.Comment),
		MarkedByTeacherID: g.TeacherID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
	if err := s.attendance.Create(ctx, row); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, apperror.Conflict("attendance_exists", "attendance already recorded for this student, group, and lesson date")
		}
		return nil, apperror.Internal("create attendance").Wrap(err)
	}
	return row, nil
}

// GetByID returns one row if the actor may view it.
func (s *Service) GetByID(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, id uuid.UUID) (*domain.Attendance, error) {
	filter, err := s.teacherFilter(ctx, actorRole, actorUserID)
	if err != nil {
		return nil, err
	}
	row, err := s.attendance.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load attendance").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("attendance")
	}
	g, err := s.groups.FindByID(ctx, row.GroupID)
	if err != nil {
		return nil, apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return nil, apperror.NotFound("attendance")
	}
	if err := s.ensureCanViewGroup(actorRole, filter, g.TeacherID); err != nil {
		return nil, err
	}
	return row, nil
}

// List returns attendance for one of: student, group, or date range.
func (s *Service) List(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, f ListFilter) ([]domain.Attendance, error) {
	filter, err := s.teacherFilter(ctx, actorRole, actorUserID)
	if err != nil {
		return nil, err
	}
	hasStudent := f.StudentID != nil
	hasGroup := f.GroupID != nil
	hasRange := f.From != nil && f.To != nil
	if (f.From != nil) != (f.To != nil) {
		return nil, apperror.Validation("filter", "from and to are both required for date range")
	}
	n := 0
	if hasStudent {
		n++
	}
	if hasGroup {
		n++
	}
	if hasRange {
		n++
	}
	if n != 1 {
		return nil, apperror.Validation("filter", "use exactly one of: student_id, group_id, or from+to date range")
	}
	switch {
	case f.StudentID != nil:
		if filter != nil {
			stuGroup, err := s.memberships.FindGroupIDByStudentUserID(ctx, *f.StudentID)
			if err != nil {
				return nil, apperror.Internal("load enrollment").Wrap(err)
			}
			if stuGroup == nil {
				return nil, apperror.New(apperror.KindForbidden, "cannot_view_student", "student has no group enrollment")
			}
			g, err := s.groups.FindByID(ctx, *stuGroup)
			if err != nil {
				return nil, apperror.Internal("load group").Wrap(err)
			}
			if g == nil || g.TeacherID != *filter {
				return nil, apperror.New(apperror.KindForbidden, "cannot_view_student", "you can only view students in your groups")
			}
		}
		return s.attendance.ListByStudent(ctx, *f.StudentID, filter)
	case f.GroupID != nil:
		g, err := s.groups.FindByID(ctx, *f.GroupID)
		if err != nil {
			return nil, apperror.Internal("load group").Wrap(err)
		}
		if g == nil {
			return nil, apperror.Validation("group_id", "group not found")
		}
		if err := s.ensureCanViewGroup(actorRole, filter, g.TeacherID); err != nil {
			return nil, err
		}
		return s.attendance.ListByGroup(ctx, *f.GroupID, filter)
	default:
		from := truncateUTCDate(*f.From)
		to := truncateUTCDate(*f.To)
		if from.After(to) {
			return nil, apperror.Validation("filter", "from must be on or before to")
		}
		return s.attendance.ListByDateRange(ctx, from, to, filter)
	}
}

// Update changes status and/or comment.
func (s *Service) Update(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, id uuid.UUID, in UpdateInput) (*domain.Attendance, error) {
	filter, err := s.teacherFilter(ctx, actorRole, actorUserID)
	if err != nil {
		return nil, err
	}
	row, err := s.attendance.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load attendance").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("attendance")
	}
	g, err := s.groups.FindByID(ctx, row.GroupID)
	if err != nil {
		return nil, apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return nil, apperror.NotFound("attendance")
	}
	if err := s.ensureCanModifyGroup(actorRole, filter, g.TeacherID); err != nil {
		return nil, err
	}
	if in.Status != nil {
		row.Status = *in.Status
	}
	if in.Comment != nil {
		row.Comment = trimComment(in.Comment)
	}
	row.UpdatedAt = time.Now().UTC()
	if err := s.attendance.Update(ctx, row); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("attendance")
		}
		return nil, apperror.Internal("update attendance").Wrap(err)
	}
	return row, nil
}

// teacherFilter returns nil for admins (no SQL restriction), or linked teacher id for teachers.
func (s *Service) teacherFilter(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID) (*uuid.UUID, error) {
	switch actorRole {
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		return nil, nil
	case domain.RoleTeacher:
		tid, err := s.links.FindTeacherIDByUserID(ctx, actorUserID)
		if err != nil {
			return nil, apperror.Internal("load teacher link").Wrap(err)
		}
		if tid == nil {
			return nil, apperror.New(apperror.KindForbidden, "teacher_not_linked", "link this user to a teacher profile in user_teacher_links to mark or view attendance")
		}
		return tid, nil
	default:
		return nil, apperror.New(apperror.KindForbidden, "insufficient_permissions", "only admin or teacher may access attendance")
	}
}

func (s *Service) ensureCanModifyGroup(actorRole domain.Role, filter *uuid.UUID, groupTeacherID uuid.UUID) error {
	return s.ensureGroupAccess(actorRole, filter, groupTeacherID, true)
}

func (s *Service) ensureCanViewGroup(actorRole domain.Role, filter *uuid.UUID, groupTeacherID uuid.UUID) error {
	return s.ensureGroupAccess(actorRole, filter, groupTeacherID, false)
}

func (s *Service) ensureGroupAccess(actorRole domain.Role, filter *uuid.UUID, groupTeacherID uuid.UUID, modify bool) error {
	switch actorRole {
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		return nil
	case domain.RoleTeacher:
		if filter == nil || *filter != groupTeacherID {
			verb := "view"
			if modify {
				verb = "modify"
			}
			return apperror.New(apperror.KindForbidden, "not_group_teacher", "you can only "+verb+" attendance for your assigned groups")
		}
		return nil
	default:
		return apperror.New(apperror.KindForbidden, "insufficient_permissions", "only admin or teacher may access attendance")
	}
}

func truncateUTCDate(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
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
