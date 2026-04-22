import {
  ClipboardCheck,
  ClipboardList,
  LayoutDashboard,
  CalendarDays,
  Settings,
  Users,
} from "lucide-react";

import type { PortalNavItem } from "@/modules/shared/layout/nav-types";

export const teacherNav: PortalNavItem[] = [
  { label: "Dashboard", to: "/teacher/dashboard", icon: LayoutDashboard },
  { label: "My students", to: "/teacher/students", icon: Users },
  { label: "Grades", to: "/teacher/grades", icon: ClipboardList },
  { label: "Attendance", to: "/teacher/attendance", icon: ClipboardCheck },
  { label: "Schedule", to: "/teacher/schedule", icon: CalendarDays },
  { label: "Settings", to: "/teacher/settings", icon: Settings },
];
