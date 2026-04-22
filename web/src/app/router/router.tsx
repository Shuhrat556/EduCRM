import { createBrowserRouter, Navigate } from "react-router-dom";

import {
  adminRoutes,
  authRoutes,
  studentRoutes,
  superAdminRoutes,
  teacherRoutes,
} from "@/routes";

export const appRouter = createBrowserRouter([
  ...authRoutes,
  studentRoutes,
  adminRoutes,
  teacherRoutes,
  superAdminRoutes,
  { path: "*", element: <Navigate to="/login" replace /> },
]);
