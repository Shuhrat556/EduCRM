import type { Weekday } from "@/features/schedule/lib/weekdays";
import type {
  LessonCreatePayload,
  LessonUpdatePayload,
} from "@/features/schedule/model/types";

const STORAGE_KEY = "educrm_mock_schedule_v1";

export type LessonRow = {
  id: string;
  weekday: Weekday;
  startTime: string;
  endTime: string;
  roomId: string;
  teacherId: string;
  groupId: string;
  title: string | null;
  createdAt: string;
  updatedAt: string;
};

function nowIso() {
  return new Date().toISOString();
}

function seed(): LessonRow[] {
  const t = nowIso();
  return [
    {
      id: "les_seed_1",
      weekday: "monday",
      startTime: "09:00",
      endTime: "10:30",
      roomId: "room_seed_2",
      teacherId: "t_seed_1",
      groupId: "g1",
      title: "Mathematics",
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "les_seed_2",
      weekday: "monday",
      startTime: "11:00",
      endTime: "12:30",
      roomId: "room_seed_1",
      teacherId: "t_seed_2",
      groupId: "g3",
      title: "Physics lab",
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "les_seed_3",
      weekday: "wednesday",
      startTime: "09:00",
      endTime: "10:30",
      roomId: "room_seed_2",
      teacherId: "t_seed_1",
      groupId: "g1",
      title: "Mathematics",
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "les_seed_4",
      weekday: "friday",
      startTime: "11:00",
      endTime: "12:30",
      roomId: "room_seed_2",
      teacherId: "t_seed_1",
      groupId: "g2",
      title: "English",
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "les_seed_5",
      weekday: "tuesday",
      startTime: "14:00",
      endTime: "15:30",
      roomId: "room_seed_3",
      teacherId: "t_seed_2",
      groupId: "g4",
      title: "English track",
      createdAt: t,
      updatedAt: t,
    },
  ];
}

function readAll(): LessonRow[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seed();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as LessonRow[];
    return Array.isArray(parsed) ? parsed : seed();
  } catch {
    return seed();
  }
}

function writeAll(rows: LessonRow[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function delay(ms = 220) {
  return new Promise((r) => setTimeout(r, ms));
}

export const mockScheduleStore = {
  async listAll(): Promise<LessonRow[]> {
    await delay();
    return readAll();
  },

  async getRowById(id: string): Promise<LessonRow | null> {
    await delay(100);
    return readAll().find((r) => r.id === id) ?? null;
  },

  async create(payload: LessonCreatePayload): Promise<LessonRow> {
    await delay();
    const rows = readAll();
    const row: LessonRow = {
      id: crypto.randomUUID(),
      weekday: payload.weekday,
      startTime: payload.startTime,
      endTime: payload.endTime,
      roomId: payload.roomId,
      teacherId: payload.teacherId,
      groupId: payload.groupId,
      title: payload.title?.trim() || null,
      createdAt: nowIso(),
      updatedAt: nowIso(),
    };
    rows.push(row);
    writeAll(rows);
    return row;
  },

  async update(id: string, payload: LessonUpdatePayload): Promise<LessonRow> {
    await delay();
    const rows = readAll();
    const i = rows.findIndex((r) => r.id === id);
    if (i === -1) throw new Error("Lesson not found");
    const cur = rows[i]!;
    const next: LessonRow = {
      ...cur,
      weekday: payload.weekday ?? cur.weekday,
      startTime: payload.startTime ?? cur.startTime,
      endTime: payload.endTime ?? cur.endTime,
      roomId: payload.roomId ?? cur.roomId,
      teacherId: payload.teacherId ?? cur.teacherId,
      groupId: payload.groupId ?? cur.groupId,
      title:
        payload.title !== undefined
          ? payload.title?.trim() || null
          : cur.title,
      updatedAt: nowIso(),
    };
    rows[i] = next;
    writeAll(rows);
    return next;
  },

  async remove(id: string): Promise<void> {
    await delay();
    writeAll(readAll().filter((r) => r.id !== id));
  },
};
