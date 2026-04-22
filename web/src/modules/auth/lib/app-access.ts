import type { Role } from "@/modules/auth/model/types";
import { primaryRole } from "@/modules/auth/lib/role-labels";
import {
  portalBasePath,
  requiredRoleForPath,
} from "@/modules/auth/lib/portal-routes";

/** All MVP roles; use for “every authenticated user” route sets. */
export const ALL_APP_ROLES: Role[] = [
  "super_admin",
  "admin",
  "teacher",
  "student",
];

export const ACCESS = {
  adminDesk: ["super_admin", "admin"] as Role[],
  withTeacher: ["super_admin", "admin", "teacher"] as Role[],
  reports: ["super_admin", "admin", "teacher"] as Role[],
  withStudentBilling: ["super_admin", "admin", "student"] as Role[],
  superAdminDesk: ["super_admin"] as Role[],
  adminOnly: ["admin"] as Role[],
  teacherDesk: ["teacher"] as Role[],
  studentDesk: ["student"] as Role[],
  all: ALL_APP_ROLES,
} as const;

function normalizePath(pathname: string): string {
  const raw = pathname.split("?")[0] ?? pathname;
  return raw.replace(/\/$/, "") || "/";
}

function portalAllowsPath(path: string, role: Role): boolean {
  const b = portalBasePath(role);

  if (
    path === `${b}/dashboard` ||
    path === `${b}/settings` ||
    path === `${b}/first-login/change-password` ||
    path === `${b}/profile/change-password`
  ) {
    return true;
  }

  switch (role) {
    case "student":
      return path === `${b}/grades` || path === `${b}/schedule`;
    case "admin":
      return (
        path.startsWith(`${b}/teachers`) || path.startsWith(`${b}/students`)
      );
    case "teacher":
      return (
        path.startsWith(`${b}/students`) ||
        path.startsWith(`${b}/grades`) ||
        path.startsWith(`${b}/attendance`) ||
        path === `${b}/schedule`
      );
    case "super_admin":
      return (
        path.startsWith(`${b}/admins`) ||
        path.startsWith(`${b}/teachers`)
      );
    default:
      return false;
  }
}

/**
 * Whether the user may open `pathname` in the SPA. Enforces portal prefix
 * (first URL segment) against the user's primary role, then sub-route rules.
 */
export function canAccessPath(pathname: string, roles: Role[] | undefined): boolean {
  if (!roles?.length) return false;
  const path = normalizePath(pathname);
  const portalRole = requiredRoleForPath(path);
  if (!portalRole) return false;
  const pr = primaryRole(roles);
  if (!pr || pr !== portalRole) return false;
  return portalAllowsPath(path, pr);
}

/** @deprecated Use `getDefaultDashboardPath` from `@/modules/auth/lib/portal-routes`. */
export function getDefaultAppPath(roles: Role[]): string {
  const pr = primaryRole(roles);
  if (!pr) return "/login";
  return `${portalBasePath(pr)}/dashboard`;
}
