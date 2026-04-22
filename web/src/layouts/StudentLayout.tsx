import { PortalWorkspace } from "@/modules/shared/routing/portal-workspace";
import type { PortalShellProps } from "@/modules/shared/layout/PortalShell";
import { studentNav } from "@/modules/student/student-nav";

const studentShell: PortalShellProps = {
  navItems: studentNav,
  dashboardPath: "/student/dashboard",
  settingsPath: "/student/settings",
  changePasswordPath: "/student/profile/change-password",
  logoutRedirect: "/login",
};

export function StudentLayout() {
  return <PortalWorkspace role="student" shell={studentShell} />;
}
