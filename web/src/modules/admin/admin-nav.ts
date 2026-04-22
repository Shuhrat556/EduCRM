import {
  BookOpen,
  CalendarDays,
  GraduationCap,
  LayoutDashboard,
  Layers,
  ListChecks,
  School,
  Settings,
  UserRound,
} from "lucide-react";

import type { PortalNavItem } from "@/modules/shared/layout/nav-types";

export const adminPortalNav: PortalNavItem[] = [
  { label: "Dashboard", to: "/admin/dashboard", icon: LayoutDashboard },
  { label: "Groups", to: "/admin/groups", icon: Layers },
  { label: "Subjects", to: "/admin/subjects", icon: BookOpen },
  { label: "Rooms", to: "/admin/rooms", icon: School },
  { label: "Schedule", to: "/admin/schedule", icon: CalendarDays },
  { label: "Attendance", to: "/admin/attendance", icon: ListChecks },
  { label: "Teachers", to: "/admin/teachers", icon: UserRound },
  { label: "Students", to: "/admin/students", icon: GraduationCap },
  { label: "Reports", to: "/admin/reports", icon: LayoutDashboard },
  { label: "Settings", to: "/admin/settings", icon: Settings },
];
