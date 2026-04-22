import type { ReactNode } from "react";
import { Navigate, useLocation } from "react-router-dom";

import { loginPathForAttemptedRoute, useAuth } from "@/modules/auth";
import { AppLoading } from "@/shared/ui/feedback/app-loading";

type ProtectedRouteProps = {
  children: ReactNode;
  /** Override login redirect (default: inferred from current path portal). */
  loginPath?: string;
};

export function ProtectedRoute({ children, loginPath }: ProtectedRouteProps) {
  const { status } = useAuth();
  const location = useLocation();
  const to = loginPath ?? loginPathForAttemptedRoute(location.pathname);

  if (status === "loading") {
    return <AppLoading fullPage label="Signing you in…" />;
  }

  if (status === "unauthenticated") {
    return <Navigate to={to} replace state={{ from: location }} />;
  }

  return <>{children}</>;
}
