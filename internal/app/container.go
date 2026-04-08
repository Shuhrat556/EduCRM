package app

import (
	"log/slog"

	"github.com/educrm/educrm-backend/internal/config"
	"github.com/educrm/educrm-backend/internal/usecase/aianalytics"
	"github.com/educrm/educrm-backend/internal/usecase/attendance"
	"github.com/educrm/educrm-backend/internal/usecase/auth"
	"github.com/educrm/educrm-backend/internal/usecase/dashboard"
	"github.com/educrm/educrm-backend/internal/usecase/file"
	"github.com/educrm/educrm-backend/internal/usecase/grade"
	"github.com/educrm/educrm-backend/internal/usecase/group"
	"github.com/educrm/educrm-backend/internal/usecase/notification"
	"github.com/educrm/educrm-backend/internal/usecase/payment"
	"github.com/educrm/educrm-backend/internal/usecase/room"
	"github.com/educrm/educrm-backend/internal/usecase/schedule"
	"github.com/educrm/educrm-backend/internal/usecase/teacher"
	"github.com/educrm/educrm-backend/internal/usecase/user"
	jwtpkg "github.com/educrm/educrm-backend/pkg/jwt"
	"gorm.io/gorm"
)

// Container is the composition root for dependency injection. Handlers and
// services receive only the dependencies they need; this type holds shared
// singletons for wiring in one place.
type Container struct {
	Config              *config.Config
	Log                 *slog.Logger
	DB                  *gorm.DB
	JWTManager          *jwtpkg.Manager
	AuthService         *auth.Service
	UserService         *user.Service
	TeacherService      *teacher.Service
	GroupService        *group.Service
	ScheduleService     *schedule.Service
	AttendanceService   *attendance.Service
	GradeService        *grade.Service
	RoomService         *room.Service
	PaymentService      *payment.Service
	DashboardService    *dashboard.Service
	FileService         *file.Service
	NotificationService *notification.Service
	AIAnalyticsService  *aianalytics.Service
}

// NewContainer builds the application container.
func NewContainer(cfg *config.Config, log *slog.Logger, db *gorm.DB, jwtMgr *jwtpkg.Manager, authSvc *auth.Service, userSvc *user.Service, teacherSvc *teacher.Service, groupSvc *group.Service, scheduleSvc *schedule.Service, attendanceSvc *attendance.Service, gradeSvc *grade.Service, roomSvc *room.Service, paymentSvc *payment.Service, dashboardSvc *dashboard.Service, fileSvc *file.Service, notificationSvc *notification.Service, aiAnalyticsSvc *aianalytics.Service) *Container {
	return &Container{
		Config:              cfg,
		Log:                 log,
		DB:                  db,
		JWTManager:          jwtMgr,
		AuthService:         authSvc,
		UserService:         userSvc,
		TeacherService:      teacherSvc,
		GroupService:        groupSvc,
		ScheduleService:     scheduleSvc,
		AttendanceService:   attendanceSvc,
		GradeService:        gradeSvc,
		RoomService:         roomSvc,
		PaymentService:      paymentSvc,
		DashboardService:    dashboardSvc,
		FileService:         fileSvc,
		NotificationService: notificationSvc,
		AIAnalyticsService:  aiAnalyticsSvc,
	}
}
