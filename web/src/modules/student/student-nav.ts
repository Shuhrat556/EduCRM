import { BookOpen, CalendarDays, LayoutDashboard, Settings } from "lucide-react";

import type { PortalNavItem } from "@/modules/shared/layout/nav-types";

export const studentNav: PortalNavItem[] = [
  { label: "Dashboard", to: "/student/dashboard", icon: LayoutDashboard },
  { label: "My grades", to: "/student/grades", icon: BookOpen },
  { label: "My schedule", to: "/student/schedule", icon: CalendarDays },
  { label: "Settings", to: "/student/settings", icon: Settings },
];
