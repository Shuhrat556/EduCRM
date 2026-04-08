package handler

import (
	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
	"github.com/educrm/educrm-backend/internal/middleware"
	"github.com/educrm/educrm-backend/internal/rbac"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ParseActor returns role and user id after AuthRequired.
func ParseActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	return middleware.ParseActor(c)
}

// ActorRole returns only the authenticated role.
func ActorRole(c *gin.Context) (domain.Role, error) {
	r, _, err := middleware.ParseActor(c)
	return r, err
}

// RequireStaff returns role if admin or super_admin.
func RequireStaff(c *gin.Context) (domain.Role, error) {
	role, _, err := middleware.ParseActor(c)
	if err != nil {
		return "", err
	}
	if err := rbac.RequireStaff(role); err != nil {
		return "", err
	}
	return role, nil
}

// RequirePaymentsReadActor allows staff payment APIs or a student reading own payments.
func RequirePaymentsReadActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return "", uuid.Nil, err
	}
	if !rbac.GrantedAny(role, rbac.PermPaymentsStaff, rbac.PermPaymentsReadOwn) {
		return "", uuid.Nil, apperror.Forbidden("insufficient permissions")
	}
	return role, uid, nil
}

// RequireNotificationInboxActor allows roles that may use the notifications inbox API.
func RequireNotificationInboxActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return "", uuid.Nil, err
	}
	if !rbac.Granted(role, rbac.PermNotificationsInbox) {
		return "", uuid.Nil, apperror.Forbidden("insufficient permissions")
	}
	return role, uid, nil
}

// RequireAITeacherRecommendationsActor allows roles with teacher-recommendation AI permission.
func RequireAITeacherRecommendationsActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return "", uuid.Nil, err
	}
	if !rbac.Granted(role, rbac.PermAITeacherRecommend) {
		return "", uuid.Nil, apperror.Forbidden("insufficient permissions")
	}
	return role, uid, nil
}

// RequireAIStudentWarningsActor allows staff or student for student warning AI endpoint.
func RequireAIStudentWarningsActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return "", uuid.Nil, err
	}
	if !rbac.GrantedAny(role, rbac.PermAIStudentWarnings, rbac.PermAIStaffAnalytics) {
		return "", uuid.Nil, apperror.Forbidden("insufficient permissions")
	}
	return role, uid, nil
}

// RequireGradesActor validates roles allowed to call grades HTTP handlers (service enforces row rules).
func RequireGradesActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return "", uuid.Nil, err
	}
	if !rbac.Granted(role, rbac.PermGradesAccess) {
		return "", uuid.Nil, apperror.Forbidden("insufficient permissions")
	}
	return role, uid, nil
}

// RequireAttendanceActor validates roles for attendance APIs.
func RequireAttendanceActor(c *gin.Context) (domain.Role, uuid.UUID, error) {
	role, uid, err := middleware.ParseActor(c)
	if err != nil {
		return "", uuid.Nil, err
	}
	if !rbac.Granted(role, rbac.PermAttendanceManage) {
		return "", uuid.Nil, apperror.Forbidden("insufficient permissions")
	}
	return role, uid, nil
}
