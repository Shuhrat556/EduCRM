package rbac

import (
	"testing"

	"github.com/educrm/educrm-backend/internal/domain"
)

func TestGranted_superAdminAlwaysTrue(t *testing.T) {
	for _, p := range AllPermissions() {
		if !Granted(domain.RoleSuperAdmin, p) {
			t.Errorf("super_admin should grant %q", p)
		}
	}
}

func TestGranted_matrix(t *testing.T) {
	tests := []struct {
		role   domain.Role
		perm   Permission
		wantOK bool
	}{
		{domain.RoleAdmin, PermUsersManage, true},
		{domain.RoleAdmin, PermPaymentsStaff, true},
		{domain.RoleAdmin, PermAITeacherRecommend, true},
		{domain.RoleTeacher, PermAttendanceManage, true},
		{domain.RoleTeacher, PermGradesAccess, true},
		{domain.RoleTeacher, PermUsersManage, false},
		{domain.RoleTeacher, PermPaymentsStaff, false},
		{domain.RoleStudent, PermGradesAccess, true},
		{domain.RoleStudent, PermPaymentsReadOwn, true},
		{domain.RoleStudent, PermPaymentsStaff, false},
		{domain.RoleStudent, PermAttendanceManage, false},
		{domain.Role("unknown"), PermGradesAccess, false},
	}
	for _, tt := range tests {
		got := Granted(tt.role, tt.perm)
		if got != tt.wantOK {
			t.Errorf("Granted(%s, %s) = %v, want %v", tt.role, tt.perm, got, tt.wantOK)
		}
	}
}

func TestGrantedAny(t *testing.T) {
	if !GrantedAny(domain.RoleStudent, PermPaymentsReadOwn, PermPaymentsStaff) {
		t.Fatal("student should match read_own")
	}
	if GrantedAny(domain.RoleStudent, PermPaymentsStaff, PermUsersManage) {
		t.Fatal("student should not match staff perms only")
	}
}

func TestIsStaff(t *testing.T) {
	tests := []struct {
		role domain.Role
		want bool
	}{
		{domain.RoleSuperAdmin, true},
		{domain.RoleAdmin, true},
		{domain.RoleTeacher, false},
		{domain.RoleStudent, false},
	}
	for _, tt := range tests {
		if got := IsStaff(tt.role); got != tt.want {
			t.Errorf("IsStaff(%s) = %v", tt.role, got)
		}
	}
}

func TestRequireStaff(t *testing.T) {
	if err := RequireStaff(domain.RoleAdmin); err != nil {
		t.Fatal(err)
	}
	if err := RequireStaff(domain.RoleTeacher); err == nil {
		t.Fatal("expected error for teacher")
	}
}
