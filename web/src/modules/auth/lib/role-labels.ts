import type { Role } from "@/modules/auth/model/types";

const ROLE_LABELS: Record<Role, string> = {
  super_admin: "Super Admin",
  admin: "Admin",
  teacher: "Teacher",
  student: "Student",
};

export function formatRoleLabel(role: Role): string {
  return ROLE_LABELS[role] ?? role;
}

/** Highest-priority role for display (e.g. badge). */
export function primaryRole(roles: Role[] | undefined): Role | null {
  if (!roles?.length) return null;
  const order: Role[] = ["super_admin", "admin", "teacher", "student"];
  for (const r of order) {
    if (roles.includes(r)) return r;
  }
  return roles[0] ?? null;
}
