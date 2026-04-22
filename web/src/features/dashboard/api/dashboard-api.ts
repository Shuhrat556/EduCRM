import { apiClient } from "@/shared/api/client";
import type { DashboardOverview } from "@/features/dashboard/model/types";

/**
 * EduCRM: `GET /api/v1/dashboard/summary` → envelope unwrap → `DashboardOverview`.
 * (Legacy `/dashboard/overview` removed; align with backend OpenAPI.)
 */
export const dashboardApi = {
  getOverview: () => apiClient.get<DashboardOverview>("/dashboard/summary"),
};
