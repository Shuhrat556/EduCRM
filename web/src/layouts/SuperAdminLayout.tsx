import { PortalWorkspace } from "@/modules/shared/routing/portal-workspace";
import type { PortalShellProps } from "@/modules/shared/layout/PortalShell";
import { superAdminNav } from "@/modules/super-admin/super-admin-nav";

const superAdminShell: PortalShellProps = {
  navItems: superAdminNav,
  dashboardPath: "/super-admin/dashboard",
  settingsPath: "/super-admin/settings",
  changePasswordPath: "/super-admin/profile/change-password",
  logoutRedirect: "/super-admin/login",
};

export function SuperAdminLayout() {
  return <PortalWorkspace role="super_admin" shell={superAdminShell} />;
}
