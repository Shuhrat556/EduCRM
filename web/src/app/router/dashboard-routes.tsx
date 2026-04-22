import type { ReactElement } from "react";

import { guardedElement } from "@/app/router/guarded-element";
import { ACCESS } from "@/modules/auth/lib/app-access";
import type { Role } from "@/modules/auth/model/types";
import { AttendancePage } from "@/pages/app/attendance/AttendancePage";
import { LessonAttendancePage } from "@/pages/app/attendance/LessonAttendancePage";
import { AiAnalyticsPage } from "@/pages/app/ai-analytics/AiAnalyticsPage";
import { GroupDetailPage } from "@/pages/app/groups/GroupDetailPage";
import { GroupsPage } from "@/pages/app/groups/GroupsPage";
import { ModulePlaceholderPage } from "@/pages/app/ModulePlaceholderPage";
import { ReportsPage } from "@/pages/app/reports/ReportsPage";
import { RoomsPage } from "@/pages/app/rooms/RoomsPage";
import { SchedulePage } from "@/pages/app/schedule/SchedulePage";
import { SettingsPage } from "@/pages/app/settings/SettingsPage";
import { StudentsPage } from "@/pages/app/students/StudentsPage";
import { SubjectsPage } from "@/pages/app/subjects/SubjectsPage";
import { TeacherDetailPage } from "@/pages/app/teachers/TeacherDetailPage";
import { TeachersPage } from "@/pages/app/teachers/TeachersPage";

function guard(roles: Role | Role[], page: ReactElement) {
  return guardedElement(roles, page);
}

/** Central place to mount module screens and their RBAC wrappers. */
export const dashboardRoutes = [
  {
    path: "students",
    element: guard(ACCESS.adminDesk, <StudentsPage />),
  },
  {
    path: "teachers/:teacherId",
    element: guard(ACCESS.adminDesk, <TeacherDetailPage />),
  },
  {
    path: "teachers",
    element: guard(ACCESS.adminDesk, <TeachersPage />),
  },
  {
    path: "groups/:groupId",
    element: guard(ACCESS.withTeacher, <GroupDetailPage />),
  },
  {
    path: "groups",
    element: guard(ACCESS.withTeacher, <GroupsPage />),
  },
  {
    path: "subjects",
    element: guard(ACCESS.adminDesk, <SubjectsPage />),
  },
  {
    path: "rooms",
    element: guard(ACCESS.adminDesk, <RoomsPage />),
  },
  {
    path: "schedule",
    element: guard(ACCESS.all, <SchedulePage />),
  },
  {
    path: "attendance/lesson/:lessonId",
    element: guard(ACCESS.all, <LessonAttendancePage />),
  },
  {
    path: "attendance",
    element: guard(ACCESS.all, <AttendancePage />),
  },
  {
    path: "payments",
    element: guard(
      ACCESS.withStudentBilling,
      <ModulePlaceholderPage title="Payments" />,
    ),
  },
  {
    path: "reports",
    element: guard(ACCESS.reports, <ReportsPage />),
  },
  {
    path: "ai-analytics",
    element: guard(ACCESS.adminDesk, <AiAnalyticsPage />),
  },
  {
    path: "settings",
    element: guard(ACCESS.all, <SettingsPage />),
  },
] as const;
