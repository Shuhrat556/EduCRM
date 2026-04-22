import { LayoutDashboard, Settings, Shield, UserRound } from "lucide-react";

import type { PortalNavItem } from "@/modules/shared/layout/nav-types";

export const superAdminNav: PortalNavItem[] = [
  { label: "Dashboard", to: "/super-admin/dashboard", icon: LayoutDashboard },
  { label: "Admins", to: "/super-admin/admins", icon: Shield },
  { label: "Teachers", to: "/super-admin/teachers", icon: UserRound },
  { label: "Settings", to: "/super-admin/settings", icon: Settings },
];
