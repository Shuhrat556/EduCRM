import type { RouteObject } from "react-router-dom";
import { Navigate } from "react-router-dom";

import { AdminLayout } from "@/layouts/AdminLayout";
import { DashboardPage } from "@/features/dashboard";
import { FirstLoginPasswordPage } from "@/modules/auth/pages/FirstLoginPasswordPage";
import { ProfileChangePasswordPage } from "@/modules/auth/pages/ProfileChangePasswordPage";
import { StudentsPage } from "@/pages/app/students/StudentsPage";
import { TeacherDetailPage } from "@/pages/app/teachers/TeacherDetailPage";
import { TeachersPage } from "@/pages/app/teachers/TeachersPage";
import { SettingsPage } from "@/pages/app/settings/SettingsPage";
import { GroupsPage } from "@/pages/app/groups/GroupsPage";
import { SubjectsPage } from "@/pages/app/subjects/SubjectsPage";
import { RoomsPage } from "@/pages/app/rooms/RoomsPage";
import { SchedulePage } from "@/pages/app/schedule/SchedulePage";
import { AttendancePage } from "@/pages/app/attendance/AttendancePage";
import { LessonAttendancePage } from "@/pages/app/attendance/LessonAttendancePage";
import { ReportsPage } from "@/pages/app/reports/ReportsPage";

export const adminRoutes: RouteObject = {
  path: "admin",
  children: [
    {
      path: "first-login/change-password",
      element: <FirstLoginPasswordPage role="admin" />,
    },
    {
      element: <AdminLayout />,
      children: [
        {
          index: true,
          element: <Navigate to="dashboard" replace />,
        },
        { path: "dashboard", element: <DashboardPage /> },
        { path: "groups", element: <GroupsPage /> },
        // TODO: add /admin/groups/:groupId when needed
        { path: "subjects", element: <SubjectsPage /> },
        { path: "rooms", element: <RoomsPage /> },
        { path: "schedule", element: <SchedulePage /> },
        { path: "attendance", element: <AttendancePage /> },
        { path: "attendance/:lessonId", element: <LessonAttendancePage /> },
        { path: "teachers", element: <TeachersPage /> },
        { path: "teachers/:teacherId", element: <TeacherDetailPage /> },
        { path: "students", element: <StudentsPage /> },
        { path: "reports", element: <ReportsPage /> },
        { path: "settings", element: <SettingsPage /> },
        {
          path: "profile/change-password",
          element: <ProfileChangePasswordPage role="admin" />,
        },
      ],
    },
  ],
};
