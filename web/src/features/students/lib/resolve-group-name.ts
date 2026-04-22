import { mockGroupsStore } from "@/features/groups/lib/mock-groups-store";
import { MOCK_GROUPS } from "@/features/students/lib/mock-groups";

export function resolveGroupName(groupId: string | null): string | null {
  if (!groupId) return null;
  if (import.meta.env.VITE_GROUPS_DEMO === "true") {
    return mockGroupsStore.getNameSync(groupId) ?? null;
  }
  return MOCK_GROUPS.find((g) => g.id === groupId)?.name ?? null;
}
