package rbac

import (
	"github.com/educrm/educrm-backend/internal/domain"
)

// Matrix documents role → permissions. super_admin is treated as full access in Granted().
//
// Role capability summary:
//
//	super_admin — all permissions (implicit).
//	admin       — same route permissions as super_admin; super-only business rules remain in services (e.g. payment free/discount).
//	teacher     — attendance + grades portal + teacher AI + notification inbox.
//	student     — grades portal + own payments read + notification inbox + student AI warnings.
var matrix = map[domain.Role]map[Permission]struct{}{
	domain.RoleAdmin: {
		PermUsersManage:         {},
		PermTeachersManage:    {},
		PermRoomsManage:       {},
		PermGroupsManage:      {},
		PermSchedulesManage:   {},
		PermAttendanceManage:  {},
		PermPaymentsStaff:     {},
		PermDashboardRead:     {},
		PermFilesManage:       {},
		PermNotificationsCreate: {},
		PermAIStaffAnalytics:  {},
		PermGradesAccess:      {},
		PermPaymentsReadOwn:   {},
		PermNotificationsInbox: {},
		PermAITeacherRecommend: {},
		PermAIStudentWarnings: {},
	},
	domain.RoleTeacher: {
		PermAttendanceManage:   {},
		PermGradesAccess:       {},
		PermNotificationsInbox: {},
		PermAITeacherRecommend: {},
	},
	domain.RoleStudent: {
		PermGradesAccess:         {},
		PermPaymentsReadOwn:      {},
		PermNotificationsInbox: {},
		PermAIStudentWarnings:  {},
	},
}

// Granted reports whether role may use a permission at the HTTP boundary.
func Granted(role domain.Role, p Permission) bool {
	if role == domain.RoleSuperAdmin {
		return true
	}
	set, ok := matrix[role]
	if !ok {
		return false
	}
	_, ok = set[p]
	return ok
}

// GrantedAny returns true if role has at least one of the permissions.
func GrantedAny(role domain.Role, perms ...Permission) bool {
	for _, p := range perms {
		if Granted(role, p) {
			return true
		}
	}
	return false
}

// IsStaff is true for admin and super_admin (route-level staff bucket).
func IsStaff(role domain.Role) bool {
	return role == domain.RoleAdmin || role == domain.RoleSuperAdmin
}

// RequireStaff returns a standard forbidden error if role is not staff.
func RequireStaff(role domain.Role) error {
	if IsStaff(role) {
		return nil
	}
	return ErrForbiddenStaff
}
