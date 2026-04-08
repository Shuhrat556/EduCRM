package grade

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Service orchestrates grade use cases.
type Service struct {
	grades      repository.GradeRepository
	groups      repository.GroupRepository
	users       repository.UserRepository
	memberships repository.StudentMembershipRepository
	links       repository.UserTeacherLinkRepository
}

// NewService constructs a grade service.
func NewService(
	grades repository.GradeRepository,
	groups repository.GroupRepository,
	users repository.UserRepository,
	memberships repository.StudentMembershipRepository,
	links repository.UserTeacherLinkRepository,
) *Service {
	return &Service{
		grades:      grades,
		groups:      groups,
		users:       users,
		memberships: memberships,
		links:       links,
	}
}

// CreateInput holds create payload.
type CreateInput struct {
	StudentID  uuid.UUID
	GroupID    uuid.UUID
	GradeType  domain.GradeType
	GradeValue float64
	Comment    *string
	WeekOf     *time.Time // optional calendar date; week bucket = WeekStartUTC(WeekOf or GradedAt)
	GradedAt   *time.Time // optional; default now
}

// UpdateInput for PATCH.
type UpdateInput struct {
	GradeValue *float64
	Comment    *string
	GradedAt   *time.Time
}

// ListFilter: exactly one of StudentID or GroupID.
type ListFilter struct {
	StudentID *uuid.UUID
	GroupID   *uuid.UUID
}

// Create inserts a weekly grade with duplicate protection per week/type.
func (s *Service) Create(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, in CreateInput) (*domain.Grade, error) {
	if math.IsNaN(in.GradeValue) || math.IsInf(in.GradeValue, 0) {
		return nil, apperror.Validation("grade_value", "must be a finite number")
	}
	g, err := s.groups.FindByID(ctx, in.GroupID)
	if err != nil {
		return nil, apperror.Internal("load group").Wrap(err)
	}
	if g == nil {
		return nil, apperror.Validation("group_id", "group not found")
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
	if err := s.ensureCreateActor(ctx, actorRole, actorUserID, in, g.TeacherID); err != nil {
		return nil, err
	}
	gradedAt := time.Now().UTC()
	if in.GradedAt != nil {
		gradedAt = in.GradedAt.UTC()
	}
	weekRef := gradedAt
	if in.WeekOf != nil {
		weekRef = truncateUTCDate(*in.WeekOf)
	}
	weekStart := domain.WeekStartUTC(weekRef)
	now := time.Now().UTC()
	row := &domain.Grade{
		ID:            uuid.New(),
		StudentID:     in.StudentID,
		TeacherID:     g.TeacherID,
		GroupID:       in.GroupID,
		WeekStartDate: weekStart,
		GradeType:     in.GradeType,
		GradeValue:    in.GradeValue,
		Comment:       trimComment(in.Comment),
		GradedAt:      gradedAt,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.grades.Create(ctx, row); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, apperror.Conflict("grade_week_exists", "a grade for this student, teacher, group, week, and grade_type already exists")
		}
		return nil, apperror.Internal("create grade").Wrap(err)
	}
	return row, nil
}

// GetByID returns a grade if visible to actor.
func (s *Service) GetByID(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, id uuid.UUID) (*domain.Grade, error) {
	row, err := s.grades.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load grade").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("grade")
	}
	if err := s.ensureCanViewRow(ctx, actorRole, actorUserID, row); err != nil {
		return nil, err
	}
	return row, nil
}

// List returns grades for a student or a group.
func (s *Service) List(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, f ListFilter) ([]domain.Grade, error) {
	if (f.StudentID == nil && f.GroupID == nil) || (f.StudentID != nil && f.GroupID != nil) {
		return nil, apperror.Validation("filter", "provide exactly one of student_id or group_id")
	}
	switch actorRole {
	case domain.RoleStudent:
		if f.GroupID != nil {
			return nil, apperror.New(apperror.KindForbidden, "cannot_list_group_grades", "students may only list their own grades")
		}
		if *f.StudentID != actorUserID {
			return nil, apperror.New(apperror.KindForbidden, "cannot_view_other_student", "you may only list your own grades")
		}
		return s.grades.ListByStudent(ctx, *f.StudentID, nil)
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		if f.StudentID != nil {
			return s.grades.ListByStudent(ctx, *f.StudentID, nil)
		}
		return s.grades.ListByGroup(ctx, *f.GroupID, nil)
	case domain.RoleTeacher:
		tid, err := s.linkedTeacherID(ctx, actorUserID)
		if err != nil {
			return nil, err
		}
		if f.StudentID != nil {
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
			if g == nil || g.TeacherID != *tid {
				return nil, apperror.New(apperror.KindForbidden, "cannot_view_student", "you can only view grades for students in your groups")
			}
			return s.grades.ListByStudent(ctx, *f.StudentID, tid)
		}
		grp, err := s.groups.FindByID(ctx, *f.GroupID)
		if err != nil {
			return nil, apperror.Internal("load group").Wrap(err)
		}
		if grp == nil {
			return nil, apperror.Validation("group_id", "group not found")
		}
		if grp.TeacherID != *tid {
			return nil, apperror.New(apperror.KindForbidden, "not_group_teacher", "you can only list grades for your assigned groups")
		}
		return s.grades.ListByGroup(ctx, *f.GroupID, tid)
	default:
		return nil, apperror.New(apperror.KindForbidden, "insufficient_permissions", "unsupported role for listing grades")
	}
}

// Update modifies value/comment/graded_at.
func (s *Service) Update(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, id uuid.UUID, in UpdateInput) (*domain.Grade, error) {
	row, err := s.grades.FindByID(ctx, id)
	if err != nil {
		return nil, apperror.Internal("load grade").Wrap(err)
	}
	if row == nil {
		return nil, apperror.NotFound("grade")
	}
	if err := s.ensureCanModifyRow(ctx, actorRole, actorUserID, row); err != nil {
		return nil, err
	}
	if in.GradeValue != nil {
		if math.IsNaN(*in.GradeValue) || math.IsInf(*in.GradeValue, 0) {
			return nil, apperror.Validation("grade_value", "must be a finite number")
		}
		row.GradeValue = *in.GradeValue
	}
	if in.Comment != nil {
		row.Comment = trimComment(in.Comment)
	}
	if in.GradedAt != nil {
		row.GradedAt = in.GradedAt.UTC()
	}
	row.UpdatedAt = time.Now().UTC()
	if err := s.grades.Update(ctx, row); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("grade")
		}
		return nil, apperror.Internal("update grade").Wrap(err)
	}
	return row, nil
}

// Delete removes a grade.
func (s *Service) Delete(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, id uuid.UUID) error {
	row, err := s.grades.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load grade").Wrap(err)
	}
	if row == nil {
		return apperror.NotFound("grade")
	}
	if err := s.ensureCanModifyRow(ctx, actorRole, actorUserID, row); err != nil {
		return err
	}
	if err := s.grades.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("grade")
		}
		return apperror.Internal("delete grade").Wrap(err)
	}
	return nil
}

func (s *Service) linkedTeacherID(ctx context.Context, actorUserID uuid.UUID) (*uuid.UUID, error) {
	tid, err := s.links.FindTeacherIDByUserID(ctx, actorUserID)
	if err != nil {
		return nil, apperror.Internal("load teacher link").Wrap(err)
	}
	if tid == nil {
		return nil, apperror.New(apperror.KindForbidden, "teacher_not_linked", "link this user to a teacher profile in user_teacher_links")
	}
	return tid, nil
}

func (s *Service) ensureCreateActor(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, in CreateInput, groupTeacherID uuid.UUID) error {
	switch in.GradeType {
	case domain.GradeTypeTeacherEvaluation:
		switch actorRole {
		case domain.RoleSuperAdmin, domain.RoleAdmin:
			return nil
		case domain.RoleTeacher:
			tid, err := s.linkedTeacherID(ctx, actorUserID)
			if err != nil {
				return err
			}
			if *tid != groupTeacherID {
				return apperror.New(apperror.KindForbidden, "not_group_teacher", "only the assigned group teacher can create teacher evaluations")
			}
			return nil
		default:
			return apperror.New(apperror.KindForbidden, "insufficient_permissions", "only admin or assigned teacher can create teacher evaluations")
		}
	case domain.GradeTypeStudentEvaluation:
		switch actorRole {
		case domain.RoleSuperAdmin, domain.RoleAdmin:
			return nil
		case domain.RoleStudent:
			if actorUserID != in.StudentID {
				return apperror.New(apperror.KindForbidden, "self_only", "you may only submit your own student evaluation")
			}
			return nil
		case domain.RoleTeacher:
			tid, err := s.linkedTeacherID(ctx, actorUserID)
			if err != nil {
				return err
			}
			if *tid != groupTeacherID {
				return apperror.New(apperror.KindForbidden, "not_group_teacher", "only the assigned group teacher can record student evaluations on behalf of students")
			}
			return nil
		default:
			return apperror.New(apperror.KindForbidden, "insufficient_permissions", "cannot create student evaluation")
		}
	default:
		return apperror.Validation("grade_type", "invalid grade type")
	}
}

func (s *Service) ensureCanViewRow(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, row *domain.Grade) error {
	switch actorRole {
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		return nil
	case domain.RoleStudent:
		if row.StudentID != actorUserID {
			return apperror.New(apperror.KindForbidden, "cannot_view_grade", "you may only view your own grades")
		}
		return nil
	case domain.RoleTeacher:
		tid, err := s.linkedTeacherID(ctx, actorUserID)
		if err != nil {
			return err
		}
		if row.TeacherID != *tid {
			return apperror.New(apperror.KindForbidden, "cannot_view_grade", "you can only view grades for your assigned groups")
		}
		return nil
	default:
		return apperror.New(apperror.KindForbidden, "insufficient_permissions", "cannot view grades")
	}
}

func (s *Service) ensureCanModifyRow(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, row *domain.Grade) error {
	switch actorRole {
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		return nil
	case domain.RoleStudent:
		if row.GradeType != domain.GradeTypeStudentEvaluation {
			return apperror.New(apperror.KindForbidden, "cannot_modify_grade", "students may only change their own student evaluations")
		}
		if row.StudentID != actorUserID {
			return apperror.New(apperror.KindForbidden, "cannot_modify_grade", "you may only modify your own grades")
		}
		return nil
	case domain.RoleTeacher:
		tid, err := s.linkedTeacherID(ctx, actorUserID)
		if err != nil {
			return err
		}
		if row.TeacherID != *tid {
			return apperror.New(apperror.KindForbidden, "cannot_modify_grade", "you can only modify grades for your assigned groups")
		}
		return nil
	default:
		return apperror.New(apperror.KindForbidden, "insufficient_permissions", "cannot modify grades")
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
