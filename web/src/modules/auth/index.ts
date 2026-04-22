export type { AuthUser, LoginCredentials, Role } from "./model/types";
export {
  loginFormSchema,
  type LoginFormValues,
} from "./model/login-schema";
export {
  resolvePostLoginPath,
  resolvePostLoginRedirect,
  resolveLoginPathForRoles,
} from "./lib/post-login-redirect";
export {
  ACCESS,
  ALL_APP_ROLES,
  canAccessPath,
  getDefaultAppPath,
} from "./lib/app-access";
export {
  dashboardPathForRole,
  firstLoginChangePasswordPath,
  getDefaultDashboardPath,
  loginPathForAttemptedRoute,
  loginPathForRole,
  portalBasePath,
  portalSlugForRole,
  profileChangePasswordPath,
  requiredRoleForPath,
  roleForPortalSlug,
} from "./lib/portal-routes";
export {
  assertPrimaryRole,
  canPermission,
  rolesWithPermission,
  type Permission,
} from "./lib/permissions";
export { formatRoleLabel, primaryRole } from "./lib/role-labels";
export { usePermissions } from "./hooks/usePermissions";
export { AuthProvider, AuthContext, type AuthContextValue } from "./context/auth-context";
export { useAuth } from "./hooks/useAuth";
export { authApi } from "./api/auth-api";
