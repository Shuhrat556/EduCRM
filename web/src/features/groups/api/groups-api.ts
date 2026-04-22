import { apiClient } from "@/shared/api/client";
import type { GroupRow } from "@/features/groups/lib/mock-groups-store";
import { mockGroupsStore } from "@/features/groups/lib/mock-groups-store";
import type {
  Group,
  GroupCreatePayload,
  GroupUpdatePayload,
  GroupsListParams,
  GroupsListResponse,
} from "@/features/groups/model/types";
import { mockRoomsStore } from "@/features/rooms/lib/mock-rooms-store";
import { mockSubjectsStore } from "@/features/subjects/lib/mock-subjects-store";
import { mockTeachersStore } from "@/features/teachers/lib/mock-teachers-store";

function useGroupsDemo() {
  return import.meta.env.VITE_GROUPS_DEMO === "true";
}

function expandDemoGroup(row: GroupRow): Group {
  const teacher = row.teacherId
    ? {
        id: row.teacherId,
        fullName:
          mockTeachersStore.getNameSync(row.teacherId) ?? "Unknown teacher",
      }
    : null;
  const subject = row.subjectId
    ? {
        id: row.subjectId,
        name:
          mockSubjectsStore.getNameSync(row.subjectId) ?? "Unknown subject",
      }
    : null;
  let room: Group["room"] = null;
  if (row.roomId) {
    const r = mockRoomsStore.getByIdSync(row.roomId);
    if (r) {
      room = { id: r.id, name: r.name, capacity: r.capacity };
    }
  }
  const subjName = subject?.name ?? "Subject";
  const roomLabel = room?.name ?? "Room";
  return {
    ...row,
    teacher,
    subject,
    room,
    schedulePreview: row.schedulePreview.map((s) => ({
      ...s,
      subjectName: subjName,
      roomName: roomLabel,
    })),
  };
}

/**
 * REST (adjust to your backend):
 * - GET /groups?page&pageSize&search&status
 * - GET /groups/:id
 * - POST /groups
 * - PUT /groups/:id
 * - DELETE /groups/:id
 */
export const groupsApi = {
  list: async (params: GroupsListParams): Promise<GroupsListResponse> => {
    if (useGroupsDemo()) {
      const { items, total, page, pageSize } =
        await mockGroupsStore.listRows(params);
      return {
        items: items.map(expandDemoGroup),
        total,
        page,
        pageSize,
      };
    }
    const { data } = await apiClient.get<GroupsListResponse>("/groups", {
      params: {
        page: params.page,
        pageSize: params.pageSize,
        search: params.search || undefined,
        status: params.status === "all" ? undefined : params.status,
      },
    });
    return data;
  },

  get: async (id: string): Promise<Group> => {
    if (useGroupsDemo()) {
      const row = await mockGroupsStore.getRowById(id);
      if (!row) throw new Error("Not found");
      return expandDemoGroup(row);
    }
    const { data } = await apiClient.get<Group>(`/groups/${id}`);
    return data;
  },

  create: async (payload: GroupCreatePayload): Promise<Group> => {
    if (useGroupsDemo()) {
      const row = await mockGroupsStore.create(payload);
      return expandDemoGroup(row);
    }
    const { data } = await apiClient.post<Group>("/groups", payload);
    return data;
  },

  update: async (
    id: string,
    payload: GroupUpdatePayload,
  ): Promise<Group> => {
    if (useGroupsDemo()) {
      const row = await mockGroupsStore.update(id, payload);
      return expandDemoGroup(row);
    }
    const { data } = await apiClient.put<Group>(`/groups/${id}`, payload);
    return data;
  },

  remove: async (id: string): Promise<void> => {
    if (useGroupsDemo()) return mockGroupsStore.remove(id);
    await apiClient.delete(`/groups/${id}`);
  },
};
