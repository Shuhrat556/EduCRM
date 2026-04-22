import type { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";

import {
  getDefaultDashboardPath,
  loginPathForAttemptedRoute,
  useAuth,
  type Role,
} from "@/modules/auth";
import { AppLoading } from "@/shared/ui/feedback/app-loading";

type RoleRouteProps = {
  roles: Role | Role[];
  children: ReactNode;
  loginPath?: string;
};

export function RoleRoute({ roles, children, loginPath }: RoleRouteProps) {
  const { hasRole, status, user } = useAuth();
  const location = useLocation();
  const toLogin = loginPath ?? loginPathForAttemptedRoute(location.pathname);

  if (status === "loading") {
    return <AppLoading fullPage label="Checking permissions…" />;
  }

  if (status === "unauthenticated") {
    return <Navigate to={toLogin} replace state={{ from: location }} />;
  }

  if (!hasRole(roles)) {
    const dest = user?.roles?.length
      ? getDefaultDashboardPath(user.roles)
      : "/unauthorized";
    return <Navigate to={dest} replace state={{ from: location.pathname }} />;
  }

  return <>{children}</>;
}
