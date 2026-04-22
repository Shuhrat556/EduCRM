import { Navigate } from "react-router-dom";

import {
  getDefaultDashboardPath,
  loginPathForRole,
  primaryRole,
  useAuth,
  type Role,
} from "@/modules/auth";
import { ChangePasswordCard } from "@/modules/auth/ui/change-password-card";
import { AppLoading } from "@/shared/ui/feedback/app-loading";
import { PageHeader } from "@/shared/ui/layout/page-header";

type ProfileChangePasswordPageProps = {
  role: Role;
};

export function ProfileChangePasswordPage({
  role,
}: ProfileChangePasswordPageProps) {
  const { status, user } = useAuth();

  if (status === "loading") {
    return <AppLoading fullPage label="Loading…" />;
  }

  if (status === "unauthenticated") {
    return <Navigate to={loginPathForRole(role)} replace />;
  }

  const pr = primaryRole(user?.roles);
  if (pr !== role) {
    return <Navigate to={getDefaultDashboardPath(user?.roles)} replace />;
  }

  return (
    <div className="space-y-8">
      <PageHeader
        title="Change password"
        description="Update your account password. Choose a strong password you have not used here before."
      />
      <ChangePasswordCard
        mode="profile"
        title="Update password"
        description="Enter your current password, then your new password twice."
        onSuccess={() => {
          /* optional: toast — staying on page is fine */
        }}
      />
    </div>
  );
}
