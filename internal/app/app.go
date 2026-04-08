package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/educrm/educrm-backend/internal/ai"
	"github.com/educrm/educrm-backend/internal/config"
	"github.com/educrm/educrm-backend/internal/database"
	httpdelivery "github.com/educrm/educrm-backend/internal/delivery/http"
	"github.com/educrm/educrm-backend/internal/delivery/http/handler"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/notify"
	"github.com/educrm/educrm-backend/internal/repository/postgres"
	"github.com/educrm/educrm-backend/internal/storage"
	aianalyticssvc "github.com/educrm/educrm-backend/internal/usecase/aianalytics"
	attendancesvc "github.com/educrm/educrm-backend/internal/usecase/attendance"
	authsvc "github.com/educrm/educrm-backend/internal/usecase/auth"
	dashboardsvc "github.com/educrm/educrm-backend/internal/usecase/dashboard"
	filesvc "github.com/educrm/educrm-backend/internal/usecase/file"
	gradesvc "github.com/educrm/educrm-backend/internal/usecase/grade"
	groupsvc "github.com/educrm/educrm-backend/internal/usecase/group"
	notificationsvc "github.com/educrm/educrm-backend/internal/usecase/notification"
	paymentsvc "github.com/educrm/educrm-backend/internal/usecase/payment"
	roomsvc "github.com/educrm/educrm-backend/internal/usecase/room"
	schedulesvc "github.com/educrm/educrm-backend/internal/usecase/schedule"
	teachersvc "github.com/educrm/educrm-backend/internal/usecase/teacher"
	usersvc "github.com/educrm/educrm-backend/internal/usecase/user"
	jwtpkg "github.com/educrm/educrm-backend/pkg/jwt"
	"gorm.io/gorm"
)

// App coordinates HTTP server lifecycle and shared dependencies.
type App struct {
	cfg        *config.Config
	log        *slog.Logger
	db         *gorm.DB
	container  *Container
	httpServer *http.Server
}

// New creates an App with database connection and dependency container.
func New(cfg *config.Config, log *slog.Logger) (*App, error) {
	db, err := database.NewPostgres(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	if cfg.AutoMigrate {
		if err := postgres.AutoMigrate(db); err != nil {
			return nil, fmt.Errorf("migrate: %w", err)
		}
	} else {
		log.Info("gorm_auto_migrate_skipped", "hint", "apply SQL migrations with: go run ./cmd/migrate up")
	}

	storageProvider, err := storage.NewFromConfig(cfg.Storage)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	aiProvider, err := ai.NewProviderFromConfig(cfg.AI)
	if err != nil {
		return nil, fmt.Errorf("ai provider: %w", err)
	}

	jwtMgr := jwtpkg.NewManager(cfg.JWT.Secret, cfg.JWT.AccessExpiration, cfg.JWT.Issuer)
	userRepo := postgres.NewUserRepository(db)
	refreshRepo := postgres.NewRefreshTokenRepository(db)
	teacherRepo := postgres.NewTeacherRepository(db)
	groupRepo := postgres.NewGroupRepository(db)
	subjectRepo := postgres.NewSubjectRepository(db)
	scheduleRepo := postgres.NewScheduleRepository(db)
	attendanceRepo := postgres.NewAttendanceRepository(db)
	gradeRepo := postgres.NewGradeRepository(db)
	membershipRepo := postgres.NewStudentMembershipRepository(db)
	teacherLinkRepo := postgres.NewUserTeacherLinkRepository(db)
	roomRepo := postgres.NewRoomRepository(db)
	paymentRepo := postgres.NewPaymentRepository(db)
	dashboardStatsRepo := postgres.NewDashboardStatsRepository(db)
	fileMetadataRepo := postgres.NewFileMetadataRepository(db)
	notificationRepo := postgres.NewNotificationRepository(db)
	aiContextRepo := postgres.NewAIAnalyticsContextRepository(db)
	outbound := notify.NewOutbound(nil, nil)
	authService := authsvc.NewService(userRepo, refreshRepo, jwtMgr, cfg.JWT.RefreshExpiration)
	userService := usersvc.NewService(userRepo)
	teacherService := teachersvc.NewService(teacherRepo)
	groupService := groupsvc.NewService(groupRepo, subjectRepo, teacherRepo, roomRepo)
	scheduleService := schedulesvc.NewService(scheduleRepo, groupRepo, teacherRepo, roomRepo)
	attendanceService := attendancesvc.NewService(attendanceRepo, groupRepo, userRepo, membershipRepo, teacherLinkRepo)
	gradeService := gradesvc.NewService(gradeRepo, groupRepo, userRepo, membershipRepo, teacherLinkRepo)
	roomService := roomsvc.NewService(roomRepo)
	paymentService := paymentsvc.NewService(paymentRepo, groupRepo, userRepo, membershipRepo)
	dashboardService := dashboardsvc.NewService(dashboardStatsRepo)
	fileService := filesvc.NewService(fileMetadataRepo, storageProvider, userRepo, teacherRepo, cfg.Storage.MaxUploadBytes)
	notificationService := notificationsvc.NewService(notificationRepo, userRepo, outbound)
	aiPrompts := ai.NewPromptCatalog(cfg.AI.PromptsDir)
	aiAnalyticsService := aianalyticssvc.NewService(aiProvider, aiPrompts, aiContextRepo, dashboardStatsRepo, userRepo, teacherLinkRepo)
	ctr := NewContainer(cfg, log, db, jwtMgr, authService, userService, teacherService, groupService, scheduleService, attendanceService, gradeService, roomService, paymentService, dashboardService, fileService, notificationService, aiAnalyticsService)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	teacherHandler := handler.NewTeacherHandler(teacherService)
	groupHandler := handler.NewGroupHandler(groupService)
	scheduleHandler := handler.NewScheduleHandler(scheduleService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)
	gradeHandler := handler.NewGradeHandler(gradeService)
	roomHandler := handler.NewRoomHandler(roomService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	fileHandler := handler.NewFileHandler(fileService, cfg.Storage)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	aiAnalyticsHandler := handler.NewAIAnalyticsHandler(aiAnalyticsService)
	engine := httpdelivery.NewRouter(cfg, log, db, &httpdelivery.RouteDeps{
		Auth:           authHandler,
		User:           userHandler,
		Teacher:        teacherHandler,
		Room:           roomHandler,
		Group:          groupHandler,
		Schedule:       scheduleHandler,
		Attendance:     attendanceHandler,
		Grade:          gradeHandler,
		Payment:        paymentHandler,
		Dashboard:      dashboardHandler,
		File:           fileHandler,
		Notification:   notificationHandler,
		AIAnalytics:    aiAnalyticsHandler,
		AuthMiddleware: middleware.AuthRequired(jwtMgr),
	})
	srv := &http.Server{
		Addr:         cfg.HTTP.Addr(),
		Handler:      engine,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}
	return &App{
		cfg:        cfg,
		log:        log,
		db:         db,
		container:  ctr,
		httpServer: srv,
	}, nil
}

// Container returns the DI container for future route registration.
func (a *App) Container() *Container {
	return a.container
}

// Run starts the HTTP server and blocks until a shutdown signal is received.
func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		a.log.Info("http_server_listen", "addr", a.httpServer.Addr)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		a.log.Info("shutdown_signal_received")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), a.cfg.ShutdownTimeout)
		defer cancel()
		if err := a.httpServer.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http shutdown: %w", err)
		}
		if err := <-errCh; err != nil {
			return err
		}
		return nil
	case err := <-errCh:
		return err
	}
}

// Close releases resources (database pool).
func (a *App) Close() error {
	if a.db == nil {
		return nil
	}
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
