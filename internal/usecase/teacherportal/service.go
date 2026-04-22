package teacherportal

import (
	"context"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/repository"
	schedulesvc "github.com/educrm/educrm-backend/internal/usecase/schedule"
	"github.com/google/uuid"
)

// Service aggregates teacher-scoped portal data (assignments, students).
type Service struct {
	links       repository.UserTeacherLinkRepository
	assignments repository.TeacherAssignmentRepository
	groups      repository.GroupRepository
	subjects    repository.SubjectRepository
	membership  repository.StudentMembershipRepository
	users       repository.UserRepository
	schedule    *schedulesvc.Service
}

// NewService returns a teacher portal service.
func NewService(
	links repository.UserTeacherLinkRepository,
	assignments repository.TeacherAssignmentRepository,
	groups repository.GroupRepository,
	subjects repository.SubjectRepository,
	membership repository.StudentMembershipRepository,
	users repository.UserRepository,
	schedule *schedulesvc.Service,
) *Service {
	return &Service{
		links:       links,
		assignments: assignments,
		groups:      groups,
		subjects:    subjects,
		membership:  membership,
		users:       users,
		schedule:    schedule,
	}
}

func (s *Service) teacherIDForActor(ctx context.Context, actorUserID uuid.UUID) (uuid.UUID, error) {
	tidPtr, err := s.links.FindTeacherIDByUserID(ctx, actorUserID)
	if err != nil {
		return uuid.Nil, err
	}
	if tidPtr == nil || *tidPtr == uuid.Nil {
		return uuid.Nil, apperror.New(apperror.KindForbidden, "teacher_profile_required", "teacher profile not linked")
	}
	return *tidPtr, nil
}

// AssignmentDTO is one row of group/subject assignment for the teacher.
type AssignmentDTO struct {
	GroupID   uuid.UUID `json:"groupId"`
	GroupName string    `json:"groupName"`
	SubjectID uuid.UUID `json:"subjectId"`
	Subject   string    `json:"subject"`
}

// ListAssignments returns group/subject pairs assigned to this teacher.
func (s *Service) ListAssignments(ctx context.Context, actorUserID uuid.UUID) ([]AssignmentDTO, error) {
	tid, err := s.teacherIDForActor(ctx, actorUserID)
	if err != nil {
		return nil, err
	}
	list, err := s.assignments.ListByTeacher(ctx, tid)
	if err != nil {
		return nil, err
	}
	out := make([]AssignmentDTO, 0, len(list))
	for _, a := range list {
		g, e := s.groups.FindByID(ctx, a.GroupID)
		if e != nil || g == nil {
			continue
		}
		var subjName string
		if a.SubjectID != uuid.Nil {
			subj, _ := s.subjects.FindByID(ctx, a.SubjectID)
			if subj != nil {
				subjName = subj.Name
			}
		}
		out = append(out, AssignmentDTO{
			GroupID:   a.GroupID,
			GroupName: g.Name,
			SubjectID: a.SubjectID,
			Subject:   subjName,
		})
	}
	return out, nil
}

// StudentRow is a student visible under the teacher's assignments for a given group.
type StudentRow struct {
	UserID    uuid.UUID `json:"userId"`
	FullName  string    `json:"fullName"`
	Username  string    `json:"username"`
	GroupID   uuid.UUID `json:"groupId"`
	GroupName string    `json:"groupName"`
}

// ListAssignedStudents returns distinct students in groups the teacher is assigned to for the optional group filter.
func (s *Service) ListAssignedStudents(ctx context.Context, actorUserID uuid.UUID, filterGroupID *uuid.UUID) ([]StudentRow, error) {
	tid, err := s.teacherIDForActor(ctx, actorUserID)
	if err != nil {
		return nil, err
	}
	assigns, err := s.assignments.ListByTeacher(ctx, tid)
	if err != nil {
		return nil, err
	}
	groupIDs := make(map[uuid.UUID]struct{})
	for _, a := range assigns {
		if filterGroupID != nil && a.GroupID != *filterGroupID {
			continue
		}
		groupIDs[a.GroupID] = struct{}{}
	}
	if len(groupIDs) == 0 {
		return []StudentRow{}, nil
	}
	seen := make(map[uuid.UUID]struct{})
	var rows []StudentRow
	for gid := range groupIDs {
		g, err := s.groups.FindByID(ctx, gid)
		if err != nil || g == nil {
			continue
		}
		studentUserIDs, err := s.membership.ListStudentUserIDsByGroup(ctx, gid)
		if err != nil {
			return nil, err
		}
		for _, uid := range studentUserIDs {
			if _, ok := seen[uid]; ok {
				continue
			}
			seen[uid] = struct{}{}
			u, err := s.users.FindByID(ctx, uid)
			if err != nil || u == nil {
				continue
			}
			uname := ""
			if u.Username != nil {
				uname = *u.Username
			}
			rows = append(rows, StudentRow{
				UserID:    u.ID,
				FullName:  u.FullName,
				Username:  uname,
				GroupID:   gid,
				GroupName: g.Name,
			})
		}
	}
	return rows, nil
}

// MySchedule returns all schedule slots for the linked teacher record.
func (s *Service) MySchedule(ctx context.Context, actorUserID uuid.UUID) ([]domain.Schedule, error) {
	tid, err := s.teacherIDForActor(ctx, actorUserID)
	if err != nil {
		return nil, err
	}
	return s.schedule.List(ctx, schedulesvc.ListFilter{TeacherID: &tid})
}
