// Package rbac defines reusable permissions and a role→permission matrix for EduCRM.
//
// # Permission matrix (HTTP layer)
//
//	super_admin — all permissions.
//	admin       — same as super_admin at routes; privileged fields (e.g. free tuition) stay super_admin-only in services.
//	teacher     — attendance.manage, grades.access, notifications.inbox, ai.teacher_recommendations.
//	student     — grades.access, payments.read_own, notifications.inbox, ai.student_warnings.
//
// Row-level authorization (own student id, assigned group, etc.) is enforced in use cases
// even when the HTTP middleware allows the role.

package rbac
