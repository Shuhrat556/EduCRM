import { Navigate, useLocation } from "react-router-dom";

import {
  firstLoginChangePasswordPath,
  getDefaultDashboardPath,
  loginPathForRole,
  primaryRole,
  useAuth,
  type Role,
} from "@/modules/auth";
import {
  PortalShell,
  type PortalShellProps,
} from "@/modules/shared/layout/PortalShell";
import { AppLoading } from "@/shared/ui/feedback/app-loading";

type PortalWorkspaceProps = {
  role: Role;
  shell: PortalShellProps;
};

/**
 * Authenticated shell for a single-role portal: wrong role → own dashboard;
 * must change password → first-login flow (except that route lives outside this tree).
 */
export function PortalWorkspace({ role, shell }: PortalWorkspaceProps) {
  const { status, user } = useAuth();
  const location = useLocation();
  const loginPath = loginPathForRole(role);

  if (status === "loading") {
    return <AppLoading fullPage label="Loading workspace…" />;
  }

  if (status === "unauthenticated") {
    return <Navigate to={loginPath} replace state={{ from: location }} />;
  }

  const pr = primaryRole(user?.roles);
  if (pr !== role) {
    return (
      <Navigate
        to={getDefaultDashboardPath(user?.roles)}
        replace
        state={{ from: location.pathname }}
      />
    );
  }

  if (user?.mustChangePassword) {
    return (
      <Navigate
        to={firstLoginChangePasswordPath(role)}
        replace
        state={{ from: location.pathname }}
      />
    );
  }

  return <PortalShell {...shell} />;
}
