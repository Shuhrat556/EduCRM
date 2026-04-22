import { useQuery } from "@tanstack/react-query";

import { aiAnalyticsApi } from "@/features/ai-analytics/api/ai-analytics-api";
import { queryKeys } from "@/shared/api/query-keys";

export function useAiAnalyticsSnapshot() {
  return useQuery({
    queryKey: queryKeys.aiAnalytics.snapshot(),
    queryFn: () => aiAnalyticsApi.getSnapshot(),
    staleTime: 60_000,
  });
}
