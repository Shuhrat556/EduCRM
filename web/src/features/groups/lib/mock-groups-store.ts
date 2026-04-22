import type {
  GroupCreatePayload,
  GroupStatus,
  GroupUpdatePayload,
  GroupsListParams,
  ScheduleSlot,
} from "@/features/groups/model/types";

const STORAGE_KEY = "educrm_mock_groups_v1";

export type GroupRow = {
  id: string;
  name: string;
  teacherId: string | null;
  subjectId: string | null;
  roomId: string | null;
  monthlyFee: number;
  startDate: string;
  endDate: string | null;
  status: GroupStatus;
  schedulePreview: ScheduleSlot[];
  createdAt: string;
  updatedAt: string;
};

function nowIso() {
  return new Date().toISOString();
}

function seedPreview(
  subjectLabel: string,
  roomLabel: string,
): ScheduleSlot[] {
  return [
    {
      weekday: "Monday",
      startTime: "09:00",
      endTime: "10:30",
      subjectName: subjectLabel,
      roomName: roomLabel,
    },
    {
      weekday: "Wednesday",
      startTime: "09:00",
      endTime: "10:30",
      subjectName: subjectLabel,
      roomName: roomLabel,
    },
    {
      weekday: "Friday",
      startTime: "11:00",
      endTime: "12:30",
      subjectName: subjectLabel,
      roomName: roomLabel,
    },
  ];
}

function seed(): GroupRow[] {
  const t = nowIso();
  return [
    {
      id: "g1",
      name: "Grade 10-A",
      teacherId: "t_seed_1",
      subjectId: "sub1",
      roomId: "room_seed_2",
      monthlyFee: 1_200_000,
      startDate: "2025-09-01",
      endDate: "2026-06-15",
      status: "active",
      schedulePreview: seedPreview("Mathematics", "Room 201"),
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "g2",
      name: "Grade 10-B",
      teacherId: "t_seed_1",
      subjectId: "sub2",
      roomId: "room_seed_1",
      monthlyFee: 1_200_000,
      startDate: "2025-09-01",
      endDate: "2026-06-15",
      status: "active",
      schedulePreview: seedPreview("English", "Lab A"),
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "g3",
      name: "Grade 11 Sciences",
      teacherId: "t_seed_2",
      subjectId: "sub3",
      roomId: "room_seed_1",
      monthlyFee: 1_450_000,
      startDate: "2025-09-01",
      endDate: null,
      status: "active",
      schedulePreview: seedPreview("Physics", "Lab A"),
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "g4",
      name: "Grade 9 — English track",
      teacherId: "t_seed_2",
      subjectId: "sub2",
      roomId: "room_seed_3",
      monthlyFee: 980_000,
      startDate: "2025-09-01",
      endDate: "2026-06-15",
      status: "draft",
      schedulePreview: seedPreview("English", "Library seminar"),
      createdAt: t,
      updatedAt: t,
    },
  ];
}

function readAll(): GroupRow[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seed();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as GroupRow[];
    return Array.isArray(parsed) ? parsed : seed();
  } catch {
    return seed();
  }
}

function writeAll(rows: GroupRow[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function delay(ms = 240) {
  return new Promise((r) => setTimeout(r, ms));
}

export const mockGroupsStore = {
  readAllSync(): GroupRow[] {
    return readAll();
  },

  getNameSync(id: string): string | undefined {
    return readAll().find((g) => g.id === id)?.name;
  },

  getRowByIdSync(id: string): GroupRow | undefined {
    return readAll().find((g) => g.id === id);
  },

  /** All groups (e.g. teacher profile labels). */
  getOptionsSync(): { id: string; name: string }[] {
    return readAll().map((g) => ({ id: g.id, name: g.name }));
  },

  getAssignableOptionsSync(): { id: string; name: string }[] {
    return readAll()
      .filter((g) => g.status === "active" || g.status === "draft")
      .map((g) => ({ id: g.id, name: g.name }));
  },

  async listRows(params: GroupsListParams): Promise<{
    items: GroupRow[];
    total: number;
    page: number;
    pageSize: number;
  }> {
    await delay();
    let rows = readAll();
    const q = params.search.trim().toLowerCase();
    if (q) {
      rows = rows.filter((g) => g.name.toLowerCase().includes(q));
    }
    if (params.status !== "all") {
      rows = rows.filter((g) => g.status === params.status);
    }
    const total = rows.length;
    const start = (params.page - 1) * params.pageSize;
    const items = rows.slice(start, start + params.pageSize);
    return {
      items,
      total,
      page: params.page,
      pageSize: params.pageSize,
    };
  },

  async getRowById(id: string): Promise<GroupRow | null> {
    await delay(120);
    return readAll().find((g) => g.id === id) ?? null;
  },

  async create(payload: GroupCreatePayload): Promise<GroupRow> {
    await delay();
    const rows = readAll();
    const preview =
      payload.schedulePreview?.length ? payload.schedulePreview : seedPreview(
        "Course",
        "TBD",
      );
    const row: GroupRow = {
      id: crypto.randomUUID(),
      name: payload.name.trim(),
      teacherId: payload.teacherId ?? null,
      subjectId: payload.subjectId ?? null,
      roomId: payload.roomId ?? null,
      monthlyFee: payload.monthlyFee,
      startDate: payload.startDate,
      endDate: payload.endDate ?? null,
      status: payload.status,
      schedulePreview: preview,
      createdAt: nowIso(),
      updatedAt: nowIso(),
    };
    rows.unshift(row);
    writeAll(rows);
    return row;
  },

  async update(id: string, payload: GroupUpdatePayload): Promise<GroupRow> {
    await delay();
    const rows = readAll();
    const i = rows.findIndex((g) => g.id === id);
    if (i === -1) throw new Error("Group not found");
    const cur = rows[i]!;
    const next: GroupRow = {
      ...cur,
      name: payload.name !== undefined ? payload.name.trim() : cur.name,
      teacherId:
        payload.teacherId !== undefined ? payload.teacherId : cur.teacherId,
      subjectId:
        payload.subjectId !== undefined ? payload.subjectId : cur.subjectId,
      roomId: payload.roomId !== undefined ? payload.roomId : cur.roomId,
      monthlyFee:
        payload.monthlyFee !== undefined ? payload.monthlyFee : cur.monthlyFee,
      startDate:
        payload.startDate !== undefined ? payload.startDate : cur.startDate,
      endDate: payload.endDate !== undefined ? payload.endDate : cur.endDate,
      status: (payload.status ?? cur.status) as GroupStatus,
      schedulePreview:
        payload.schedulePreview !== undefined
          ? payload.schedulePreview
          : cur.schedulePreview,
      updatedAt: nowIso(),
    };
    rows[i] = next;
    writeAll(rows);
    return next;
  },

  async remove(id: string): Promise<void> {
    await delay();
    writeAll(readAll().filter((g) => g.id !== id));
  },
};
