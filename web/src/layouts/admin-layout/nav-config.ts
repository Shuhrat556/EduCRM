import type { LucideIcon } from "lucide-react";
import {
  BarChart3,
  BookOpen,
  Building2,
  CalendarDays,
  ClipboardCheck,
  CreditCard,
  GraduationCap,
  LayoutDashboard,
  Settings,
  Sparkles,
  UserRound,
  Users,
} from "lucide-react";

import { ACCESS } from "@/modules/auth/lib/app-access";
import type { Role } from "@/modules/auth/model/types";

export type NavItem = {
  label: string;
  to: string;
  icon: LucideIcon;
  /** If set, item is shown only when the user has one of these roles. */
  roles?: Role[];
};

/**
 * Primary sidebar navigation for EduCRM admin.
 * `roles` must stay aligned with `canAccessPath` in `app-access.ts`.
 */
export const adminNav: NavItem[] = [
  { label: "Dashboard", to: "/app", icon: LayoutDashboard },
  {
    label: "Students",
    to: "/app/students",
    icon: GraduationCap,
    roles: ACCESS.adminDesk,
  },
  {
    label: "Teachers",
    to: "/app/teachers",
    icon: UserRound,
    roles: ACCESS.adminDesk,
  },
  {
    label: "Groups",
    to: "/app/groups",
    icon: Users,
    roles: ACCESS.withTeacher,
  },
  {
    label: "Subjects",
    to: "/app/subjects",
    icon: BookOpen,
    roles: ACCESS.adminDesk,
  },
  {
    label: "Rooms",
    to: "/app/rooms",
    icon: Building2,
    roles: ACCESS.adminDesk,
  },
  { label: "Schedule", to: "/app/schedule", icon: CalendarDays, roles: ACCESS.all },
  {
    label: "Attendance",
    to: "/app/attendance",
    icon: ClipboardCheck,
    roles: ACCESS.all,
  },
  {
    label: "Payments",
    to: "/app/payments",
    icon: CreditCard,
    roles: ACCESS.withStudentBilling,
  },
  {
    label: "Reports",
    to: "/app/reports",
    icon: BarChart3,
    roles: ACCESS.reports,
  },
  {
    label: "AI Analytics",
    to: "/app/ai-analytics",
    icon: Sparkles,
    roles: ACCESS.adminDesk,
  },
  { label: "Settings", to: "/app/settings", icon: Settings, roles: ACCESS.all },
];

const pathToLabel = new Map(adminNav.map((item) => [item.to, item.label]));

/** Resolve a human label for breadcrumb trail; falls back to Title Case segment. */
export function labelForAdminPath(pathname: string): string | undefined {
  const normalized = pathname.replace(/\/$/, "") || "/app";
  if (pathToLabel.has(normalized)) return pathToLabel.get(normalized);
  return undefined;
}
