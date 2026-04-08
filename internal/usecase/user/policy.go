package user

import (
	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
)

// CanManageUsers is true for admin and super_admin (user management APIs).
func CanManageUsers(actor domain.Role) bool {
	return actor == domain.RoleSuperAdmin || actor == domain.RoleAdmin
}

// AdminManagesOnlyStaff is true when role is teacher or student (admin-visible users).
func AdminManagesOnlyStaff(r domain.Role) bool {
	return r == domain.RoleTeacher || r == domain.RoleStudent
}

// AssertActorCanManageUsers returns an error if the actor cannot use user APIs at all.
func AssertActorCanManageUsers(actor domain.Role) error {
	if !CanManageUsers(actor) {
		return apperror.Forbidden("user management is restricted to administrators")
	}
	return nil
}

// AssertCanCreateTargetRole enforces:
// - only super_admin may create users with role admin or super_admin;
// - admin may create teacher and student only.
func AssertCanCreateTargetRole(actor, target domain.Role) error {
	if err := AssertActorCanManageUsers(actor); err != nil {
		return err
	}
	switch target {
	case domain.RoleSuperAdmin, domain.RoleAdmin:
		if actor != domain.RoleSuperAdmin {
			return apperror.Forbidden("only super_admin can create users with role admin or super_admin")
		}
		return nil
	case domain.RoleTeacher, domain.RoleStudent:
		return nil
	default:
		return apperror.Validation("role", "invalid role")
	}
}

// AssertActorCanAccessTargetUser checks read/update/delete/status against target role.
// super_admin may access any user; admin may access only teacher and student.
func AssertActorCanAccessTargetUser(actor domain.Role, targetRole domain.Role) error {
	if err := AssertActorCanManageUsers(actor); err != nil {
		return err
	}
	if actor == domain.RoleSuperAdmin {
		return nil
	}
	if actor == domain.RoleAdmin && AdminManagesOnlyStaff(targetRole) {
		return nil
	}
	return apperror.Forbidden("insufficient permissions for this user")
}

// AssertCanAssignRoleOnUpdate validates role transitions on update.
// super_admin may assign any role; admin may only assign teacher or student.
func AssertCanAssignRoleOnUpdate(actor, newRole domain.Role) error {
	if actor == domain.RoleSuperAdmin {
		if _, err := domain.ParseRole(string(newRole)); err != nil {
			return apperror.Validation("role", "invalid role")
		}
		return nil
	}
	if actor == domain.RoleAdmin {
		if newRole != domain.RoleTeacher && newRole != domain.RoleStudent {
			return apperror.Forbidden("only super_admin may assign admin or super_admin roles")
		}
		return nil
	}
	return apperror.Forbidden("insufficient permissions")
}
