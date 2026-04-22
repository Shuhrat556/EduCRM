import { canAccessPath } from "@/modules/auth/lib/app-access";
import type { Role } from "@/modules/auth/model/types";
import { primaryRole } from "@/modules/auth/lib/role-labels";
import {
  firstLoginChangePasswordPath,
  getDefaultDashboardPath,
  loginPathForRole,
} from "@/modules/auth/lib/portal-routes";

export function resolvePostLoginPath(roles: Role[]): string {
  return getDefaultDashboardPath(roles);
}

function isLoginPath(path: string): boolean {
  if (path === "/login") return true;
  return ["/admin/login", "/teacher/login", "/super-admin/login"].some(
    (p) => path === p || path.startsWith(`${p}/`),
  );
}

/**
 * Prefers returning the user to a deep link when safe; otherwise the dashboard.
 * If the user must change their password, they are sent to the first-login flow first.
 */
export function resolvePostLoginRedirect(
  fromPath: string | undefined,
  roles: Role[],
  options?: { mustChangePassword?: boolean },
): string {
  const pr = primaryRole(roles);
  if (options?.mustChangePassword && pr) {
    return firstLoginChangePasswordPath(pr);
  }

  const home = resolvePostLoginPath(roles);
  if (!fromPath || isLoginPath(fromPath)) {
    return home;
  }
  if (fromPath === "/" || fromPath === "") {
    return home;
  }
  if (canAccessPath(fromPath, roles)) {
    return fromPath;
  }
  return home;
}

/** Login path to offer on “wrong portal” / sign-out flows for a role. */
export function resolveLoginPathForRoles(roles: Role[] | undefined): string {
  const pr = primaryRole(roles);
  if (!pr) return "/login";
  return loginPathForRole(pr);
}
