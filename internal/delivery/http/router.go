package httpdelivery

import (
	"log/slog"
	"path/filepath"
	"strings"

	"github.com/educrm/educrm-backend/internal/config"
	"github.com/educrm/educrm-backend/internal/delivery/http/handler"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	_ "github.com/educrm/educrm-backend/docs" // swagger generated
)

// NewRouter builds the Gin engine with global middleware and route groups.
func NewRouter(cfg *config.Config, log *slog.Logger, db *gorm.DB, deps *RouteDeps) *gin.Engine {
	if cfg != nil && cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	if cfg != nil {
		if proxies := cfg.HTTP.TrustedProxyList(); len(proxies) > 0 {
			_ = r.SetTrustedProxies(proxies)
		}
	}

	skipLog := []string(nil)
	if cfg != nil {
		skipLog = cfg.LogHTTPSkipPrefixes()
	}
	r.Use(
		middleware.Recovery(log),
		middleware.RequestID(),
	)
	if cfg != nil {
		origins := cfg.EffectiveCORSOrigins()
		r.Use(middleware.CORS(origins, cfg.CORS.AllowCredentials))
	} else {
		r.Use(middleware.CORS([]string{"*"}, false))
	}
	r.Use(
		middleware.RequestLogger(log, skipLog),
		middleware.ErrorHandler(),
	)
	if cfg != nil && cfg.RateLimit.Enabled {
		rps := cfg.RateLimit.RPS
		if rps <= 0 {
			rps = 100
		}
		r.Use(middleware.RateLimit(rps, cfg.RateLimit.Burst, skipLog))
	}

	if cfg == nil || cfg.SwaggerEnabled {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Local storage: public URLs under STORAGE_PUBLIC_BASE_URL must include this path prefix.
	if cfg != nil {
		p := strings.ToLower(strings.TrimSpace(cfg.Storage.Provider))
		if p == "" || p == "local" {
			if abs, err := filepath.Abs(cfg.Storage.LocalDir); err == nil {
				r.Static("/static/files", abs)
			}
		}
	}

	health := handler.NewHealthHandler(db)
	r.GET("/health", health.Live)

	v1 := r.Group("/api/v1")
	{
		v1.GET("/health", health.Ready)
	}

	if deps != nil && deps.Auth != nil {
		auth := v1.Group("/auth")
		{
			auth.POST("/login", deps.Auth.Login)
			auth.POST("/refresh", deps.Auth.Refresh)
			auth.POST("/logout", deps.Auth.Logout)
			if deps.AuthMiddleware != nil {
				auth.GET("/me", deps.AuthMiddleware, deps.Auth.Me)
				auth.POST("/change-password", deps.AuthMiddleware, deps.Auth.ChangePassword)
				auth.POST("/first-login/change-password", deps.AuthMiddleware, deps.Auth.FirstLoginChangePassword)
			}
		}
	}

	if deps != nil && deps.User != nil && deps.AuthMiddleware != nil {
		u := v1.Group("/users")
		u.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermUsersManage))
		{
			u.POST("", deps.User.Create)
			u.GET("", deps.User.List)
			u.GET("/:id", deps.User.GetByID)
			u.PATCH("/:id/status", deps.User.SetStatus)
			u.PATCH("/:id", deps.User.Update)
			u.DELETE("/:id", deps.User.Delete)
		}
	}

	if deps != nil && deps.Teacher != nil && deps.AuthMiddleware != nil {
		t := v1.Group("/teachers")
		t.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermTeachersManage))
		{
			t.POST("", deps.Teacher.Create)
			t.GET("", deps.Teacher.List)
			t.GET("/:id", deps.Teacher.GetByID)
			t.PATCH("/:id/photo", deps.Teacher.PatchPhoto)
			t.PATCH("/:id", deps.Teacher.Update)
			t.DELETE("/:id", deps.Teacher.Delete)
		}
	}

	if deps != nil && deps.Room != nil && deps.AuthMiddleware != nil {
		rm := v1.Group("/rooms")
		rm.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermRoomsManage))
		{
			rm.POST("", deps.Room.Create)
			rm.GET("", deps.Room.List)
			rm.GET("/:id", deps.Room.GetByID)
			rm.PATCH("/:id", deps.Room.Update)
			rm.DELETE("/:id", deps.Room.Delete)
		}
	}

	if deps != nil && deps.Group != nil && deps.AuthMiddleware != nil {
		g := v1.Group("/groups")
		g.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermGroupsManage))
		{
			g.POST("", deps.Group.Create)
			g.GET("", deps.Group.List)
			g.GET("/:id", deps.Group.GetByID)
			g.PATCH("/:id", deps.Group.Update)
			g.DELETE("/:id", deps.Group.Delete)
		}
	}

	if deps != nil && deps.Schedule != nil && deps.AuthMiddleware != nil {
		sch := v1.Group("/schedules")
		sch.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermSchedulesManage))
		{
			sch.POST("", deps.Schedule.Create)
			sch.GET("", deps.Schedule.List)
			sch.GET("/:id", deps.Schedule.GetByID)
			sch.PATCH("/:id", deps.Schedule.Update)
			sch.DELETE("/:id", deps.Schedule.Delete)
		}
	}

	if deps != nil && deps.Attendance != nil && deps.AuthMiddleware != nil {
		at := v1.Group("/attendance")
		at.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermAttendanceManage))
		{
			at.POST("", deps.Attendance.Create)
			at.GET("", deps.Attendance.List)
			at.GET("/:id", deps.Attendance.GetByID)
			at.PATCH("/:id", deps.Attendance.Update)
		}
	}

	if deps != nil && deps.Grade != nil && deps.AuthMiddleware != nil {
		gr := v1.Group("/grades")
		gr.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermGradesAccess))
		{
			gr.POST("", deps.Grade.Create)
			gr.GET("", deps.Grade.List)
			gr.GET("/:id", deps.Grade.GetByID)
			gr.PATCH("/:id", deps.Grade.Update)
			gr.DELETE("/:id", deps.Grade.Delete)
		}
	}

	if deps != nil && deps.Dashboard != nil && deps.AuthMiddleware != nil {
		dash := v1.Group("/dashboard")
		dash.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermDashboardRead))
		{
			dash.GET("/summary", deps.Dashboard.Summary)
		}
	}

	if deps != nil && deps.File != nil && deps.AuthMiddleware != nil {
		fl := v1.Group("/files")
		fl.Use(deps.AuthMiddleware, middleware.RequirePermission(rbac.PermFilesManage))
		{
			fl.POST("", deps.File.Upload)
			fl.POST("/register", deps.File.Register)
			fl.GET("", deps.File.List)
			fl.GET("/:id", deps.File.GetByID)
			fl.DELETE("/:id", deps.File.Delete)
		}
	}

	if deps != nil && deps.Notification != nil && deps.AuthMiddleware != nil {
		n := v1.Group("/notifications")
		n.Use(deps.AuthMiddleware)
		{
			n.POST("", middleware.RequirePermission(rbac.PermNotificationsCreate), deps.Notification.Create)
			inbox := middleware.RequirePermission(rbac.PermNotificationsInbox)
			n.GET("", inbox, deps.Notification.List)
			n.PATCH("/:id/read", inbox, deps.Notification.MarkRead)
			n.GET("/:id", inbox, deps.Notification.GetByID)
			n.DELETE("/:id", inbox, deps.Notification.Delete)
		}
	}

	if deps != nil && deps.AIAnalytics != nil && deps.AuthMiddleware != nil {
		ai := v1.Group("/ai/analytics")
		ai.Use(deps.AuthMiddleware)
		{
			staffAI := middleware.RequirePermission(rbac.PermAIStaffAnalytics)
			ai.POST("/debtors-summary", staffAI, deps.AIAnalytics.DebtorsSummary)
			ai.POST("/low-attendance", staffAI, deps.AIAnalytics.LowAttendance)
			ai.POST("/admin-daily-summary", staffAI, deps.AIAnalytics.AdminDailySummary)
			ai.POST("/teacher-recommendations", middleware.RequirePermission(rbac.PermAITeacherRecommend), deps.AIAnalytics.TeacherRecommendations)
			ai.POST("/student-warnings", middleware.RequireAnyPermission(rbac.PermAIStudentWarnings, rbac.PermAIStaffAnalytics), deps.AIAnalytics.StudentWarnings)
		}
	}

	if deps != nil && deps.TeacherPortal != nil && deps.AuthMiddleware != nil {
		tp := v1.Group("/teacher")
		tp.Use(deps.AuthMiddleware, middleware.RequireRoles(domain.RoleTeacher))
		{
			tp.GET("/assignments", deps.TeacherPortal.ListAssignments)
			tp.GET("/students", deps.TeacherPortal.ListAssignedStudents)
			tp.GET("/schedule", deps.TeacherPortal.MySchedule)
		}
	}

	if deps != nil && deps.StudentPortal != nil && deps.AuthMiddleware != nil {
		sp := v1.Group("/student")
		sp.Use(deps.AuthMiddleware, middleware.RequireRoles(domain.RoleStudent))
		{
			sp.GET("/grades", deps.StudentPortal.MyGrades)
			sp.GET("/schedule", deps.StudentPortal.MySchedule)
			sp.GET("/attendance", deps.StudentPortal.MyAttendance)
		}
	}

	if deps != nil && deps.Payment != nil && deps.AuthMiddleware != nil {
		staffPay := middleware.RequirePermission(rbac.PermPaymentsStaff)
		payRead := middleware.RequireAnyPermission(rbac.PermPaymentsStaff, rbac.PermPaymentsReadOwn)
		p := v1.Group("/payments")
		p.Use(deps.AuthMiddleware)
		{
			p.POST("", staffPay, deps.Payment.Create)
			p.GET("", staffPay, deps.Payment.List)
			p.GET("/history", payRead, deps.Payment.History)
			p.GET("/:id", payRead, deps.Payment.GetByID)
			p.PATCH("/:id", staffPay, deps.Payment.Update)
			p.DELETE("/:id", staffPay, deps.Payment.Delete)
		}
	}

	return r
}
