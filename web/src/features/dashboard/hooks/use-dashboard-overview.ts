import { useQuery } from "@tanstack/react-query";

import { dashboardApi } from "@/features/dashboard/api/dashboard-api";
import { emptyDashboardOverview } from "@/features/dashboard/lib/empty-overview";
import { mockDashboardOverview } from "@/features/dashboard/lib/mock-overview";
import { queryKeys } from "@/shared/api/query-keys";

function isDemoMode() {
  return import.meta.env.VITE_DASHBOARD_DEMO === "true";
}

export function useDashboardOverview() {
  return useQuery({
    queryKey: queryKeys.dashboard.overview(),
    queryFn: async () => {
      if (isDemoMode()) {
        await new Promise((r) => setTimeout(r, 450));
        return mockDashboardOverview;
      }
      try {
        const { data } = await dashboardApi.getOverview();
        return data;
      } catch {
        return emptyDashboardOverview;
      }
    },
    staleTime: 60_000,
  });
}
