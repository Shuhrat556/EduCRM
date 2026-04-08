package aianalytics

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/educrm/educrm-backend/internal/ai"
	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/educrm/educrm-backend/internal/repository"
	"github.com/google/uuid"
)

// Result is a completed analytics generation.
type Result struct {
	Output   string
	Provider string
}

// Service orchestrates context loading, prompt rendering, and provider calls.
type Service struct {
	provider  ai.Provider
	prompts   *ai.PromptCatalog
	context   repository.AIAnalyticsContextRepository
	dashboard repository.DashboardStatsRepository
	users     repository.UserRepository
	links     repository.UserTeacherLinkRepository
}

// NewService constructs the AI analytics service.
func NewService(
	provider ai.Provider,
	prompts *ai.PromptCatalog,
	ctxRepo repository.AIAnalyticsContextRepository,
	dashboard repository.DashboardStatsRepository,
	users repository.UserRepository,
	links repository.UserTeacherLinkRepository,
) *Service {
	return &Service{
		provider:  provider,
		prompts:   prompts,
		context:   ctxRepo,
		dashboard: dashboard,
		users:     users,
		links:     links,
	}
}

// DebtorsSummary runs the debtors analytics prompt for staff.
func (s *Service) DebtorsSummary(ctx context.Context, actorRole domain.Role, monthFor *time.Time) (*Result, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	ref := time.Now().UTC()
	if monthFor != nil {
		ref = monthFor.UTC()
	}
	monthStart := time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, time.UTC)
	data, err := s.context.DebtorsSummaryData(ctx, monthStart)
	if err != nil {
		return nil, apperror.Internal("load debtors context").Wrap(err)
	}
	return s.runScenario(ctx, "debtors_summary", string(data))
}

// LowAttendanceSummary runs attendance risk analytics for staff.
func (s *Service) LowAttendanceSummary(ctx context.Context, actorRole domain.Role, from, to *time.Time) (*Result, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	toD := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	fromD := toD.AddDate(0, 0, -30)
	if from != nil {
		fromD = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	}
	if to != nil {
		toD = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, time.UTC)
	}
	data, err := s.context.LowAttendanceData(ctx, fromD, toD)
	if err != nil {
		return nil, apperror.Internal("load attendance context").Wrap(err)
	}
	return s.runScenario(ctx, "low_attendance_summary", string(data))
}

// AdminDailySummary composes dashboard snapshot JSON and runs the admin brief prompt.
func (s *Service) AdminDailySummary(ctx context.Context, actorRole domain.Role, asOf *time.Time) (*Result, error) {
	if err := rbac.RequireStaff(actorRole); err != nil {
		return nil, err
	}
	ref := time.Now().UTC()
	if asOf != nil {
		ref = asOf.UTC()
	}
	snap, err := s.dashboard.Snapshot(ctx, ref)
	if err != nil {
		return nil, apperror.Internal("dashboard snapshot").Wrap(err)
	}
	monthStart := time.Date(ref.Year(), ref.Month(), 1, 0, 0, 0, 0, time.UTC)
	rev, err := s.dashboard.MonthlyPaidRevenue(ctx, monthStart)
	if err != nil {
		return nil, apperror.Internal("monthly revenue").Wrap(err)
	}
	payload, err := json.Marshal(map[string]any{
		"snapshot": snap,
		"revenue":  rev,
	})
	if err != nil {
		return nil, apperror.Internal("encode dashboard").Wrap(err)
	}
	return s.runScenario(ctx, "admin_daily_summary", string(payload))
}

// TeacherRecommendations resolves the teacher profile and loads group context.
func (s *Service) TeacherRecommendations(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, teacherID *uuid.UUID) (*Result, error) {
	tid, err := s.resolveTeacherID(ctx, actorRole, actorUserID, teacherID)
	if err != nil {
		return nil, err
	}
	data, err := s.context.TeacherGroupsData(ctx, tid)
	if err != nil {
		return nil, apperror.Internal("load teacher context").Wrap(err)
	}
	return s.runScenario(ctx, "teacher_recommendations", string(data))
}

// StudentWarningSuggestions loads risk signals for a student (self or staff).
func (s *Service) StudentWarningSuggestions(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, studentID *uuid.UUID) (*Result, error) {
	sid, err := s.resolveStudentID(actorRole, actorUserID, studentID)
	if err != nil {
		return nil, err
	}
	if err := s.assertStudentUser(ctx, sid); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	from := now.AddDate(0, 0, -30)
	fromD := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, time.UTC)
	toD := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	data, err := s.context.StudentRiskData(ctx, sid, monthStart, fromD, toD)
	if err != nil {
		return nil, apperror.Internal("load student risk context").Wrap(err)
	}
	return s.runScenario(ctx, "student_warning_suggestions", string(data))
}

func (s *Service) runScenario(ctx context.Context, scenario, dataJSON string) (*Result, error) {
	sys, user, err := s.prompts.Render(scenario, dataJSON)
	if err != nil {
		return nil, apperror.Internal("render prompts").Wrap(err)
	}
	out, err := s.provider.Generate(ctx, ai.GenerateInput{SystemPrompt: sys, UserPrompt: user})
	if err != nil {
		return nil, apperror.Internal("ai generate").Wrap(err)
	}
	name := out.ProviderName
	if strings.TrimSpace(name) == "" {
		name = s.provider.Name()
	}
	return &Result{Output: out.Text, Provider: name}, nil
}

func (s *Service) resolveTeacherID(ctx context.Context, actorRole domain.Role, actorUserID uuid.UUID, staffTeacherID *uuid.UUID) (uuid.UUID, error) {
	if rbac.IsStaff(actorRole) && staffTeacherID != nil {
		return *staffTeacherID, nil
	}
	if actorRole != domain.RoleTeacher && !rbac.IsStaff(actorRole) {
		return uuid.Nil, apperror.New(apperror.KindForbidden, "forbidden", "only teacher or staff may request teacher recommendations")
	}
	if actorRole == domain.RoleTeacher {
		tid, err := s.links.FindTeacherIDByUserID(ctx, actorUserID)
		if err != nil {
			return uuid.Nil, apperror.Internal("load teacher link").Wrap(err)
		}
		if tid == nil {
			return uuid.Nil, apperror.Validation("teacher", "no teacher profile linked to this user")
		}
		return *tid, nil
	}
	return uuid.Nil, apperror.Validation("teacher_id", "required for staff")
}

func (s *Service) resolveStudentID(actorRole domain.Role, actorUserID uuid.UUID, staffStudentID *uuid.UUID) (uuid.UUID, error) {
	if rbac.IsStaff(actorRole) && staffStudentID != nil {
		return *staffStudentID, nil
	}
	if actorRole == domain.RoleStudent {
		return actorUserID, nil
	}
	if !rbac.IsStaff(actorRole) {
		return uuid.Nil, apperror.New(apperror.KindForbidden, "forbidden", "only student or staff may request student warnings")
	}
	return uuid.Nil, apperror.Validation("student_id", "required for staff")
}

func (s *Service) assertStudentUser(ctx context.Context, id uuid.UUID) error {
	u, err := s.users.FindByID(ctx, id)
	if err != nil {
		return apperror.Internal("load user").Wrap(err)
	}
	if u == nil || u.Role != domain.RoleStudent {
		return apperror.Validation("student_id", "must be a student user")
	}
	return nil
}

