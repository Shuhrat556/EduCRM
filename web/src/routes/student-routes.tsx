import type { RouteObject } from "react-router-dom";
import { Navigate } from "react-router-dom";

import { StudentLayout } from "@/layouts/StudentLayout";
import { FirstLoginPasswordPage } from "@/modules/auth/pages/FirstLoginPasswordPage";
import { ProfileChangePasswordPage } from "@/modules/auth/pages/ProfileChangePasswordPage";
import { StudentDashboardPage } from "@/modules/student/pages/StudentDashboardPage";
import { StudentGradesPage } from "@/modules/student/pages/StudentGradesPage";
import { SchedulePage } from "@/pages/app/schedule/SchedulePage";
import { SettingsPage } from "@/pages/app/settings/SettingsPage";

export const studentRoutes: RouteObject = {
  path: "student",
  children: [
    {
      path: "first-login/change-password",
      element: <FirstLoginPasswordPage role="student" />,
    },
    {
      element: <StudentLayout />,
      children: [
        {
          index: true,
          element: <Navigate to="dashboard" replace />,
        },
        { path: "dashboard", element: <StudentDashboardPage /> },
        { path: "grades", element: <StudentGradesPage /> },
        { path: "schedule", element: <SchedulePage /> },
        { path: "settings", element: <SettingsPage /> },
        {
          path: "profile/change-password",
          element: <ProfileChangePasswordPage role="student" />,
        },
      ],
    },
  ],
};
