import type { Role } from "@/modules/auth/model/types";
import { primaryRole } from "@/modules/auth/lib/role-labels";

/** URL segment for static routing (matches React Router paths). */
export type PortalSlug = "student" | "admin" | "teacher" | "super-admin";

const ROLE_TO_SLUG: Record<Role, PortalSlug> = {
  student: "student",
  admin: "admin",
  teacher: "teacher",
  super_admin: "super-admin",
};

const SLUG_TO_ROLE: Record<PortalSlug, Role> = {
  student: "student",
  admin: "admin",
  teacher: "teacher",
  "super-admin": "super_admin",
};

export function portalSlugForRole(role: Role): PortalSlug {
  return ROLE_TO_SLUG[role];
}

export function roleForPortalSlug(slug: string | undefined): Role | null {
  if (!slug) return null;
  if (slug in SLUG_TO_ROLE) return SLUG_TO_ROLE[slug as PortalSlug];
  return null;
}

export function portalBasePath(role: Role): string {
  return `/${portalSlugForRole(role)}`;
}

/** Student portal uses `/login`; others use `/{portal}/login` pattern except super-admin. */
export function loginPathForRole(role: Role): string {
  if (role === "student") return "/login";
  return `${portalBasePath(role)}/login`;
}

export function dashboardPathForRole(role: Role): string {
  return `${portalBasePath(role)}/dashboard`;
}

export function firstLoginChangePasswordPath(role: Role): string {
  return `${portalBasePath(role)}/first-login/change-password`;
}

export function profileChangePasswordPath(role: Role): string {
  return `${portalBasePath(role)}/profile/change-password`;
}

export function getDefaultDashboardPath(roles: Role[] | undefined): string {
  const pr = primaryRole(roles);
  if (!pr) return "/login";
  return dashboardPathForRole(pr);
}

/** First path segment must match the user's primary role portal. */
export function requiredRoleForPath(pathname: string): Role | null {
  const seg = pathname.split("/").filter(Boolean)[0];
  return roleForPortalSlug(seg);
}

/**
 * When an unauthenticated user hits a protected route, send them to the matching
 * role login (falls back to student `/login`).
 */
export function loginPathForAttemptedRoute(pathname: string | undefined): string {
  if (!pathname) return "/login";
  const role = requiredRoleForPath(pathname);
  if (role) return loginPathForRole(role);
  return "/login";
}
