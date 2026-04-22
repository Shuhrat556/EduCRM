import { PortalWorkspace } from "@/modules/shared/routing/portal-workspace";
import type { PortalShellProps } from "@/modules/shared/layout/PortalShell";
import { teacherNav } from "@/modules/teacher/teacher-nav";

const teacherShell: PortalShellProps = {
  navItems: teacherNav,
  dashboardPath: "/teacher/dashboard",
  settingsPath: "/teacher/settings",
  changePasswordPath: "/teacher/profile/change-password",
  logoutRedirect: "/teacher/login",
};

export function TeacherLayout() {
  return <PortalWorkspace role="teacher" shell={teacherShell} />;
}
