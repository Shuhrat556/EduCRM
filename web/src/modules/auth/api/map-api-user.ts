import type { AuthUser, Role } from "@/modules/auth/model/types";

function isRole(s: string): s is Role {
  return (
    s === "super_admin" ||
    s === "admin" ||
    s === "teacher" ||
    s === "student"
  );
}

function normalizeRoles(raw: unknown): Role[] {
  if (raw == null) return [];
  if (typeof raw === "string" && isRole(raw)) return [raw];
  if (Array.isArray(raw)) {
    return raw.filter(
      (x): x is Role => typeof x === "string" && isRole(x),
    );
  }
  return [];
}

/**
 * Maps `GET /auth/me` payload (snake_case / `role` string) to `AuthUser`.
 */
export function mapApiUserToAuthUser(raw: unknown): AuthUser {
  if (!raw || typeof raw !== "object") {
    throw new Error("Invalid user payload from API");
  }
  const o = raw as Record<string, unknown>;
  const roleField = o.role ?? o.roles;
  const roles = normalizeRoles(roleField);

  const displayNameRaw =
    o.display_name ?? o.displayName ?? o.full_name ?? o.fullName ?? "";

  const mustChangePassword =
    o.must_change_password === true ||
    o.mustChangePassword === true ||
    o.require_password_change === true;

  return {
    id: String(o.id ?? ""),
    email: String(o.email ?? ""),
    phone: o.phone != null ? String(o.phone) : undefined,
    displayName: String(displayNameRaw || o.email || "User"),
    roles,
    mustChangePassword,
  };
}
