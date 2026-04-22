import { apiClient } from "@/shared/api/client";
import { emptyAiAnalyticsSnapshot, mockAiAnalyticsSnapshot } from "@/features/ai-analytics/lib/mock-snapshot";
import type { AiAnalyticsSnapshot } from "@/features/ai-analytics/model/types";

function useAiAnalyticsDemo() {
  return import.meta.env.VITE_AI_ANALYTICS_DEMO === "true";
}

function delay(ms: number) {
  return new Promise<void>((r) => setTimeout(r, ms));
}

/**
 * `GET /ai/analytics/snapshot` — replace path and shape to match your AI service.
 */
export const aiAnalyticsApi = {
  getSnapshot: async (): Promise<AiAnalyticsSnapshot> => {
    if (useAiAnalyticsDemo()) {
      await delay(900);
      return {
        ...mockAiAnalyticsSnapshot,
        generatedAt: new Date().toISOString(),
      };
    }

    try {
      const { data } = await apiClient.get<AiAnalyticsSnapshot>(
        "/ai/analytics/snapshot",
      );
      return { ...data, source: "api" };
    } catch {
      await delay(500);
      return emptyAiAnalyticsSnapshot(
        "Connect your AI analytics endpoint or set VITE_AI_ANALYTICS_DEMO=true for sample data.",
      );
    }
  },
};
