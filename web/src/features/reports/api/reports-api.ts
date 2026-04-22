import { apiClient } from "@/shared/api/client";
import { aggregateReportsLocal } from "@/features/reports/lib/aggregate-reports-local";
import type { ReportsFilters, ReportsSummary } from "@/features/reports/model/types";

function serializeParams(filters: ReportsFilters) {
  return {
    date_from: filters.dateFrom,
    date_to: filters.dateTo,
    group_id: filters.groupId ?? undefined,
    teacher_id: filters.teacherId ?? undefined,
    student_status:
      filters.studentStatus === "all" ? undefined : filters.studentStatus,
  };
}

/**
 * Prefer `GET /reports/summary` when the backend implements it; otherwise
 * aggregate from students, teachers, schedule, attendance, and payments APIs.
 */
export const reportsApi = {
  getSummary: async (filters: ReportsFilters): Promise<ReportsSummary> => {
    try {
      const { data } = await apiClient.get<ReportsSummary>("/reports/summary", {
        params: serializeParams(filters),
      });
      return { ...data, source: "api" };
    } catch {
      return aggregateReportsLocal(filters);
    }
  },
};
