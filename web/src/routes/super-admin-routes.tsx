import type { RouteObject } from "react-router-dom";
import { Navigate } from "react-router-dom";

import { SuperAdminLayout } from "@/layouts/SuperAdminLayout";
import { DashboardPage } from "@/features/dashboard";
import { FirstLoginPasswordPage } from "@/modules/auth/pages/FirstLoginPasswordPage";
import { ProfileChangePasswordPage } from "@/modules/auth/pages/ProfileChangePasswordPage";
import { AdminsPage } from "@/modules/super-admin/pages/AdminsPage";
import { SettingsPage } from "@/pages/app/settings/SettingsPage";
import { TeacherDetailPage } from "@/pages/app/teachers/TeacherDetailPage";
import { TeachersPage } from "@/pages/app/teachers/TeachersPage";

export const superAdminRoutes: RouteObject = {
  path: "super-admin",
  children: [
    {
      path: "first-login/change-password",
      element: <FirstLoginPasswordPage role="super_admin" />,
    },
    {
      element: <SuperAdminLayout />,
      children: [
        {
          index: true,
          element: <Navigate to="dashboard" replace />,
        },
        { path: "dashboard", element: <DashboardPage /> },
        { path: "admins", element: <AdminsPage /> },
        { path: "teachers", element: <TeachersPage /> },
        { path: "teachers/:teacherId", element: <TeacherDetailPage /> },
        { path: "settings", element: <SettingsPage /> },
        {
          path: "profile/change-password",
          element: <ProfileChangePasswordPage role="super_admin" />,
        },
      ],
    },
  ],
};
