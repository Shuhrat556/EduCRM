import type { Role } from "@/modules/auth/model/types";
import { primaryRole } from "@/modules/auth/lib/role-labels";

/**
 * Central permission helpers for UI and routing. Use `hasPermission` from hooks
 * that read `useAuth().user` — avoid sprinkling raw role checks in pages.
 */
export type Permission =
  | "super_admin:manage_admins"
  | "super_admin:manage_teachers"
  | "admin:manage_teachers"
  | "admin:manage_students"
  | "teacher:view_assigned_students"
  | "teacher:assign_grades"
  | "teacher:mark_attendance"
  | "teacher:view_schedule"
  | "student:view_own_grades"
  | "student:view_own_schedule";

const PERMISSION_MATRIX: Record<Permission, Role[]> = {
  "super_admin:manage_admins": ["super_admin"],
  "super_admin:manage_teachers": ["super_admin"],
  "admin:manage_teachers": ["admin"],
  "admin:manage_students": ["admin"],
  "teacher:view_assigned_students": ["teacher"],
  "teacher:assign_grades": ["teacher"],
  "teacher:mark_attendance": ["teacher"],
  "teacher:view_schedule": ["teacher"],
  "student:view_own_grades": ["student"],
  "student:view_own_schedule": ["student"],
};

export function rolesWithPermission(permission: Permission): readonly Role[] {
  return PERMISSION_MATRIX[permission];
}

export function canPermission(
  roles: Role[] | undefined,
  permission: Permission,
): boolean {
  if (!roles?.length) return false;
  const allowed = PERMISSION_MATRIX[permission];
  return roles.some((r) => allowed.includes(r));
}

/** Use primary role for portal routing; permissions still use full `roles`. */
export function assertPrimaryRole(
  roles: Role[] | undefined,
  expected: Role,
): boolean {
  return primaryRole(roles) === expected;
}
