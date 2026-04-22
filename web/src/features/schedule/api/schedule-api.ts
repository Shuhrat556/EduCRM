import { apiClient } from "@/shared/api/client";
import { compareTimeStrings } from "@/features/schedule/lib/time-compare";
import { WEEKDAY_ORDER } from "@/features/schedule/lib/weekdays";
import type { LessonRow } from "@/features/schedule/lib/mock-schedule-store";
import { mockScheduleStore } from "@/features/schedule/lib/mock-schedule-store";
import type {
  LessonCreatePayload,
  LessonUpdatePayload,
  ScheduleLesson,
} from "@/features/schedule/model/types";
import { mockGroupsStore } from "@/features/groups/lib/mock-groups-store";
import { mockRoomsStore } from "@/features/rooms/lib/mock-rooms-store";
import { mockTeachersStore } from "@/features/teachers/lib/mock-teachers-store";

function useScheduleDemo() {
  return import.meta.env.VITE_SCHEDULE_DEMO === "true";
}

function expandLesson(row: LessonRow): ScheduleLesson {
  const room = mockRoomsStore.getByIdSync(row.roomId);
  const teacherName = mockTeachersStore.getNameSync(row.teacherId);
  const groupName = mockGroupsStore.getNameSync(row.groupId);
  return {
    ...row,
    room: room
      ? { id: room.id, name: room.name }
      : { id: row.roomId, name: "Unknown room" },
    teacher: {
      id: row.teacherId,
      fullName: teacherName ?? "Unknown teacher",
    },
    group: {
      id: row.groupId,
      name: groupName ?? "Unknown group",
    },
  };
}

function sortLessons(rows: ScheduleLesson[]): ScheduleLesson[] {
  const order = new Map(WEEKDAY_ORDER.map((d, i) => [d, i]));
  return [...rows].sort((a, b) => {
    const da = order.get(a.weekday) ?? 99;
    const db = order.get(b.weekday) ?? 99;
    if (da !== db) return da - db;
    return compareTimeStrings(a.startTime, b.startTime);
  });
}

/**
 * REST (adjust to your backend):
 * - GET /schedule/lessons → full week list
 * - GET /schedule/lessons/:id
 * - POST /schedule/lessons
 * - PUT /schedule/lessons/:id
 * - DELETE /schedule/lessons/:id
 */
export const scheduleApi = {
  listWeek: async (): Promise<ScheduleLesson[]> => {
    if (useScheduleDemo()) {
      const rows = await mockScheduleStore.listAll();
      return sortLessons(rows.map(expandLesson));
    }
    const { data } = await apiClient.get<ScheduleLesson[]>("/schedule/lessons");
    return sortLessons(data);
  },

  get: async (id: string): Promise<ScheduleLesson> => {
    if (useScheduleDemo()) {
      const row = await mockScheduleStore.getRowById(id);
      if (!row) throw new Error("Not found");
      return expandLesson(row);
    }
    const { data } = await apiClient.get<ScheduleLesson>(
      `/schedule/lessons/${id}`,
    );
    return data;
  },

  create: async (payload: LessonCreatePayload): Promise<ScheduleLesson> => {
    if (useScheduleDemo()) {
      const row = await mockScheduleStore.create(payload);
      return expandLesson(row);
    }
    const { data } = await apiClient.post<ScheduleLesson>(
      "/schedule/lessons",
      payload,
    );
    return data;
  },

  update: async (
    id: string,
    payload: LessonUpdatePayload,
  ): Promise<ScheduleLesson> => {
    if (useScheduleDemo()) {
      const row = await mockScheduleStore.update(id, payload);
      return expandLesson(row);
    }
    const { data } = await apiClient.put<ScheduleLesson>(
      `/schedule/lessons/${id}`,
      payload,
    );
    return data;
  },

  remove: async (id: string): Promise<void> => {
    if (useScheduleDemo()) return mockScheduleStore.remove(id);
    await apiClient.delete(`/schedule/lessons/${id}`);
  },
};
