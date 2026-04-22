import type {
  Room,
  RoomCreatePayload,
  RoomUpdatePayload,
  RoomsListParams,
  RoomsListResponse,
} from "@/features/rooms/model/types";

const STORAGE_KEY = "educrm_mock_rooms_v1";

function nowIso() {
  return new Date().toISOString();
}

function seed(): Room[] {
  const t = nowIso();
  return [
    {
      id: "room_seed_1",
      name: "Lab A",
      building: "Science wing",
      capacity: 24,
      notes: "Projector, chemistry kit storage",
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "room_seed_2",
      name: "Room 201",
      building: "Main hall",
      capacity: 32,
      notes: null,
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "room_seed_3",
      name: "Library seminar",
      building: "Library",
      capacity: 16,
      notes: "Quiet zone",
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "room_seed_4",
      name: "Gymnasium",
      building: "Sports block",
      capacity: 120,
      notes: null,
      createdAt: t,
      updatedAt: t,
    },
  ];
}

function readAll(): Room[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seed();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as Room[];
    return Array.isArray(parsed) ? parsed : seed();
  } catch {
    return seed();
  }
}

function writeAll(rows: Room[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function delay(ms = 240) {
  return new Promise((r) => setTimeout(r, ms));
}

export const mockRoomsStore = {
  getByIdSync(id: string): Room | undefined {
    return readAll().find((r) => r.id === id);
  },

  async list(params: RoomsListParams): Promise<RoomsListResponse> {
    await delay();
    let rows = readAll();
    const q = params.search.trim().toLowerCase();
    if (q) {
      rows = rows.filter(
        (r) =>
          r.name.toLowerCase().includes(q) ||
          (r.building?.toLowerCase().includes(q) ?? false) ||
          (r.notes?.toLowerCase().includes(q) ?? false) ||
          String(r.capacity).includes(q),
      );
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

  async getById(id: string): Promise<Room | null> {
    await delay(150);
    return readAll().find((r) => r.id === id) ?? null;
  },

  async create(payload: RoomCreatePayload): Promise<Room> {
    await delay();
    const rows = readAll();
    const r: Room = {
      id: crypto.randomUUID(),
      name: payload.name.trim(),
      building: payload.building?.trim() || null,
      capacity: payload.capacity,
      notes: payload.notes?.trim() || null,
      createdAt: nowIso(),
      updatedAt: nowIso(),
    };
    rows.unshift(r);
    writeAll(rows);
    return r;
  },

  async update(id: string, payload: RoomUpdatePayload): Promise<Room> {
    await delay();
    const rows = readAll();
    const i = rows.findIndex((r) => r.id === id);
    if (i === -1) throw new Error("Room not found");
    const cur = rows[i]!;
    const next: Room = {
      ...cur,
      name: payload.name !== undefined ? payload.name.trim() : cur.name,
      building:
        payload.building !== undefined
          ? payload.building?.trim() || null
          : cur.building,
      capacity:
        payload.capacity !== undefined ? payload.capacity : cur.capacity,
      notes:
        payload.notes !== undefined
          ? payload.notes?.trim() || null
          : cur.notes,
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
