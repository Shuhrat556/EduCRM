import { apiClient } from "@/shared/api/client";
import { mockGroupsStore } from "@/features/groups/lib/mock-groups-store";
import { MOCK_GROUPS } from "@/features/students/lib/mock-groups";
import type { StudentGroupOption } from "@/features/students/model/types";

function useStudentsDemo() {
  return import.meta.env.VITE_STUDENTS_DEMO === "true";
}

function useGroupsDemo() {
  return import.meta.env.VITE_GROUPS_DEMO === "true";
}

/**
 * Dropdown options for assigning a student to a group.
 * Live API: `GET /groups/options` → `{ id, name }[]`.
 */
export async function fetchGroupOptions(): Promise<StudentGroupOption[]> {
  if (useGroupsDemo()) {
    return mockGroupsStore.getAssignableOptionsSync();
  }
  if (useStudentsDemo()) {
    return MOCK_GROUPS;
  }
  const { data } = await apiClient.get<StudentGroupOption[]>("/groups/options");
  return data;
}
