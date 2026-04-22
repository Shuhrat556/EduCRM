import { PortalWorkspace } from "@/modules/shared/routing/portal-workspace";
import type { PortalShellProps } from "@/modules/shared/layout/PortalShell";
import { adminPortalNav } from "@/modules/admin/admin-nav";

const adminShell: PortalShellProps = {
  navItems: adminPortalNav,
  dashboardPath: "/admin/dashboard",
  settingsPath: "/admin/settings",
  changePasswordPath: "/admin/profile/change-password",
  logoutRedirect: "/admin/login",
};

export function AdminLayout() {
  return <PortalWorkspace role="admin" shell={adminShell} />;
}
