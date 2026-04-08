package httpdelivery

import (
	"github.com/educrm/educrm-backend/internal/delivery/http/handler"
	"github.com/gin-gonic/gin"
)

// RouteDeps carries HTTP handlers and middleware wired in the application layer.
type RouteDeps struct {
	Auth           *handler.AuthHandler
	User           *handler.UserHandler
	Teacher        *handler.TeacherHandler
	Room           *handler.RoomHandler
	Group          *handler.GroupHandler
	Schedule       *handler.ScheduleHandler
	Attendance     *handler.AttendanceHandler
	Grade          *handler.GradeHandler
	Payment        *handler.PaymentHandler
	Dashboard      *handler.DashboardHandler
	File           *handler.FileHandler
	Notification   *handler.NotificationHandler
	AIAnalytics    *handler.AIAnalyticsHandler
	AuthMiddleware gin.HandlerFunc
}
