package domain

// Role is persisted as a string and enforced in application code.
type Role string

const (
	RoleSuperAdmin Role = "super_admin"
	RoleAdmin      Role = "admin"
	RoleTeacher    Role = "teacher"
	RoleStudent    Role = "student"
)

var validRoles = map[Role]struct{}{
	RoleSuperAdmin: {},
	RoleAdmin:      {},
	RoleTeacher:    {},
	RoleStudent:    {},
}

// ParseRole returns an error if s is not a known role.
func ParseRole(s string) (Role, error) {
	r := Role(s)
	if _, ok := validRoles[r]; !ok {
		return "", ErrInvalidRole
	}
	return r, nil
}
