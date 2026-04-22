/**
 * MVP roles for RBAC. Align with your auth API / JWT claims (`super_admin`,
 * `admin`, `teacher`, `student`).
 */
export type Role = "super_admin" | "admin" | "teacher" | "student";

export interface AuthUser {
  id: string;
  email: string;
  phone?: string;
  displayName: string;
  roles: Role[];
  /** When true, user must set a new password before using the portal. */
  mustChangePassword?: boolean;
}

export interface LoginCredentials {
  identifier: string;
  password: string;
  rememberMe: boolean;
}

export interface AuthSession {
  user: AuthUser;
  accessToken: string;
}
