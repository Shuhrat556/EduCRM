package rbac

// Permission names coarse HTTP/route capabilities. Resource-level rules stay in use cases.
type Permission string

const (
	// Staff directory & catalog
	PermUsersManage          Permission = "users.manage"
	PermTeachersManage       Permission = "teachers.manage"
	PermRoomsManage          Permission = "rooms.manage"
	PermGroupsManage         Permission = "groups.manage"
	PermSchedulesManage      Permission = "schedules.manage"
	PermAttendanceManage     Permission = "attendance.manage"
	PermPaymentsStaff        Permission = "payments.staff"
	PermDashboardRead        Permission = "dashboard.read"
	PermFilesManage          Permission = "files.manage"
	PermNotificationsCreate  Permission = "notifications.create"
	PermAIStaffAnalytics     Permission = "ai.staff_analytics"

	// Shared portals (HTTP allows; services enforce row scope)
	PermGradesAccess         Permission = "grades.access"
	PermPaymentsReadOwn      Permission = "payments.read_own"
	PermNotificationsInbox   Permission = "notifications.inbox"
	PermAITeacherRecommend   Permission = "ai.teacher_recommendations"
	PermAIStudentWarnings    Permission = "ai.student_warnings"
)

// AllPermissions lists every defined permission (for tests and audits).
func AllPermissions() []Permission {
	return []Permission{
		PermUsersManage,
		PermTeachersManage,
		PermRoomsManage,
		PermGroupsManage,
		PermSchedulesManage,
		PermAttendanceManage,
		PermPaymentsStaff,
		PermDashboardRead,
		PermFilesManage,
		PermNotificationsCreate,
		PermAIStaffAnalytics,
		PermGradesAccess,
		PermPaymentsReadOwn,
		PermNotificationsInbox,
		PermAITeacherRecommend,
		PermAIStudentWarnings,
	}
}
