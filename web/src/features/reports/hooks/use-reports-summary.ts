import { useQuery } from "@tanstack/react-query";

import { reportsApi } from "@/features/reports/api/reports-api";
import type { ReportsFilters } from "@/features/reports/model/types";
import { queryKeys } from "@/shared/api/query-keys";

export function useReportsSummary(filters: ReportsFilters) {
  return useQuery({
    queryKey: queryKeys.reports.summary(filters),
    queryFn: () => reportsApi.getSummary(filters),
    staleTime: 30_000,
  });
}
