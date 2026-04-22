import type { RouteObject } from "react-router-dom";
import { Navigate } from "react-router-dom";

import { TeacherLayout } from "@/layouts/TeacherLayout";
import { DashboardPage } from "@/features/dashboard";
import { FirstLoginPasswordPage } from "@/modules/auth/pages/FirstLoginPasswordPage";
import { ProfileChangePasswordPage } from "@/modules/auth/pages/ProfileChangePasswordPage";
import { TeacherGradesPage } from "@/modules/teacher/pages/TeacherGradesPage";
import { TeacherStudentsPage } from "@/modules/teacher/pages/TeacherStudentsPage";
import { AttendancePage } from "@/pages/app/attendance/AttendancePage";
import { LessonAttendancePage } from "@/pages/app/attendance/LessonAttendancePage";
import { SchedulePage } from "@/pages/app/schedule/SchedulePage";
import { SettingsPage } from "@/pages/app/settings/SettingsPage";

export const teacherRoutes: RouteObject = {
  path: "teacher",
  children: [
    {
      path: "first-login/change-password",
      element: <FirstLoginPasswordPage role="teacher" />,
    },
    {
      element: <TeacherLayout />,
      children: [
        {
          index: true,
          element: <Navigate to="dashboard" replace />,
        },
        { path: "dashboard", element: <DashboardPage /> },
        { path: "students", element: <TeacherStudentsPage /> },
        { path: "grades", element: <TeacherGradesPage /> },
        { path: "attendance", element: <AttendancePage /> },
        {
          path: "attendance/lesson/:lessonId",
          element: <LessonAttendancePage />,
        },
        { path: "schedule", element: <SchedulePage /> },
        { path: "settings", element: <SettingsPage /> },
        {
          path: "profile/change-password",
          element: <ProfileChangePasswordPage role="teacher" />,
        },
      ],
    },
  ],
};
