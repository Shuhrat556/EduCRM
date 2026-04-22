import type { RouteObject } from "react-router-dom";
import { Navigate } from "react-router-dom";

import { RoleLoginPage } from "@/modules/auth/pages/RoleLoginPage";
import { UnauthorizedPage } from "@/pages/unauthorized/UnauthorizedPage";

export const authRoutes: RouteObject[] = [
  { path: "/", element: <Navigate to="/login" replace /> },

  {
    path: "/login",
    element: (
      <RoleLoginPage
        expectedRole="student"
        subtitle="Sign in to the student portal with your school email or phone."
      />
    ),
  },
  {
    path: "/admin/login",
    element: (
      <RoleLoginPage
        expectedRole="admin"
        subtitle="Administrator access — manage teachers and students."
      />
    ),
  },
  {
    path: "/teacher/login",
    element: (
      <RoleLoginPage
        expectedRole="teacher"
        subtitle="Teacher portal — your classes, grades, and attendance."
      />
    ),
  },
  {
    path: "/super-admin/login",
    element: (
      <RoleLoginPage
        expectedRole="super_admin"
        subtitle="Super Admin — manage institution admins and teachers."
      />
    ),
  },

  { path: "/unauthorized", element: <UnauthorizedPage /> },
];
