package user

import (
	"testing"

	"github.com/educrm/educrm-backend/internal/apperror"
	"github.com/educrm/educrm-backend/internal/domain"
)

func TestAssertCanCreateTargetRole(t *testing.T) {
	tests := []struct {
		name    string
		actor   domain.Role
		target  domain.Role
		wantErr bool
	}{
		{"super creates admin", domain.RoleSuperAdmin, domain.RoleAdmin, false},
		{"super creates super", domain.RoleSuperAdmin, domain.RoleSuperAdmin, false},
		{"admin creates teacher", domain.RoleAdmin, domain.RoleTeacher, false},
		{"admin creates student", domain.RoleAdmin, domain.RoleStudent, false},
		{"admin cannot create admin", domain.RoleAdmin, domain.RoleAdmin, true},
		{"admin cannot create super", domain.RoleAdmin, domain.RoleSuperAdmin, true},
		{"teacher forbidden", domain.RoleTeacher, domain.RoleStudent, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AssertCanCreateTargetRole(tt.actor, tt.target)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAssertActorCanAccessTargetUser(t *testing.T) {
	tests := []struct {
		name       string
		actor      domain.Role
		targetRole domain.Role
		wantErr    bool
	}{
		{"super any", domain.RoleSuperAdmin, domain.RoleAdmin, false},
		{"admin teacher", domain.RoleAdmin, domain.RoleTeacher, false},
		{"admin student", domain.RoleAdmin, domain.RoleStudent, false},
		{"admin cannot admin", domain.RoleAdmin, domain.RoleAdmin, true},
		{"teacher forbidden", domain.RoleTeacher, domain.RoleStudent, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AssertActorCanAccessTargetUser(tt.actor, tt.targetRole)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAssertCanAssignRoleOnUpdate(t *testing.T) {
	tests := []struct {
		name    string
		actor   domain.Role
		newRole domain.Role
		wantErr bool
		kind    apperror.Kind
	}{
		{"super assigns admin", domain.RoleSuperAdmin, domain.RoleAdmin, false, ""},
		{"admin assigns teacher", domain.RoleAdmin, domain.RoleTeacher, false, ""},
		{"admin cannot assign admin", domain.RoleAdmin, domain.RoleAdmin, true, apperror.KindForbidden},
		{"teacher forbidden", domain.RoleTeacher, domain.RoleStudent, true, apperror.KindForbidden},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AssertCanAssignRoleOnUpdate(tt.actor, tt.newRole)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				if ae, ok := apperror.AsError(err); ok && tt.kind != "" && ae.Kind != tt.kind {
					t.Fatalf("kind %v, want %v", ae.Kind, tt.kind)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
