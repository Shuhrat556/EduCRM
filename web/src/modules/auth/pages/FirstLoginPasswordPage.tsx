import { Navigate, useNavigate } from "react-router-dom";

import {
  dashboardPathForRole,
  getDefaultDashboardPath,
  loginPathForRole,
  primaryRole,
  useAuth,
  type Role,
} from "@/modules/auth";
import { ChangePasswordCard } from "@/modules/auth/ui/change-password-card";
import { AppLoading } from "@/shared/ui/feedback/app-loading";

type FirstLoginPasswordPageProps = {
  role: Role;
};

export function FirstLoginPasswordPage({ role }: FirstLoginPasswordPageProps) {
  const { status, user } = useAuth();
  const navigate = useNavigate();
  const loginPath = loginPathForRole(role);

  if (status === "loading") {
    return <AppLoading fullPage label="Checking your account…" />;
  }

  if (status === "unauthenticated") {
    return <Navigate to={loginPath} replace />;
  }

  const pr = primaryRole(user?.roles);
  if (pr !== role) {
    return <Navigate to={getDefaultDashboardPath(user?.roles)} replace />;
  }

  if (!user?.mustChangePassword) {
    return <Navigate to={dashboardPathForRole(role)} replace />;
  }

  return (
    <div className="flex min-h-dvh items-center justify-center bg-gradient-to-b from-background via-background to-muted/40 p-4">
      <ChangePasswordCard
        mode="first"
        title="Set a new password"
        description="For security, you must choose a new password before continuing."
        onSuccess={() => navigate(dashboardPathForRole(role), { replace: true })}
      />
    </div>
  );
}
