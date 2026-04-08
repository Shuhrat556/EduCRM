package rbac

import "github.com/educrm/educrm-backend/internal/apperror"

// ErrForbiddenStaff is returned when an action requires admin or super_admin.
var ErrForbiddenStaff = apperror.New(apperror.KindForbidden, "staff_only", "only admin or super_admin may perform this action")
