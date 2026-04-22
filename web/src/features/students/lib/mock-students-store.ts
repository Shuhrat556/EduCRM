import { resolveGroupName } from "@/features/students/lib/resolve-group-name";
import type {
  Student,
  StudentCreatePayload,
  StudentStatus,
  StudentUpdatePayload,
  StudentsListParams,
  StudentsListResponse,
} from "@/features/students/model/types";

const STORAGE_KEY = "educrm_mock_students_v1";

function nowIso() {
  return new Date().toISOString();
}

function seed(): Student[] {
  const t = nowIso();
  return [
    {
      id: "s_seed_1",
      fullName: "Dilnoza Rahimova",
      phone: "+998901112233",
      email: "dilnoza@student.edu",
      status: "active",
      groupId: "g1",
      groupName: "Grade 10-A",
      photoUrl: null,
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "s_seed_2",
      fullName: "Rustam Toshmatov",
      phone: "+998907654321",
      email: "rustam.t@student.edu",
      status: "active",
      groupId: "g2",
      groupName: "Grade 10-B",
      photoUrl: null,
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "s_seed_3",
      fullName: "Madina Karimova",
      phone: "+998931234567",
      email: "madina.k@student.edu",
      status: "inactive",
      groupId: null,
      groupName: null,
      photoUrl: null,
      createdAt: t,
      updatedAt: t,
    },
  ];
}

function readAll(): Student[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seed();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as Student[];
    return Array.isArray(parsed) ? parsed : seed();
  } catch {
    return seed();
  }
}

function writeAll(rows: Student[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function expandStudent(s: Student): Student {
  return {
    ...s,
    groupName: resolveGroupName(s.groupId),
  };
}

function delay(ms = 280) {
  return new Promise((r) => setTimeout(r, ms));
}

export const mockStudentsStore = {
  async list(params: StudentsListParams): Promise<StudentsListResponse> {
    await delay();
    let rows = readAll();
    const q = params.search.trim().toLowerCase();
    if (q) {
      rows = rows.filter(
        (s) =>
          s.fullName.toLowerCase().includes(q) ||
          s.phone.replace(/\s/g, "").includes(q.replace(/\s/g, "")),
      );
    }
    if (params.status !== "all") {
      rows = rows.filter((s) => s.status === params.status);
    }
    if (params.groupId) {
      rows = rows.filter((s) => s.groupId === params.groupId);
    }
    const total = rows.length;
    const start = (params.page - 1) * params.pageSize;
    const items = rows.slice(start, start + params.pageSize).map(expandStudent);
    return {
      items,
      total,
      page: params.page,
      pageSize: params.pageSize,
    };
  },

  async getById(id: string): Promise<Student | null> {
    await delay(180);
    const s = readAll().find((r) => r.id === id);
    return s ? expandStudent(s) : null;
  },

  async create(payload: StudentCreatePayload): Promise<Student> {
    await delay();
    const rows = readAll();
    const gid = payload.groupId?.trim() || null;
    const s: Student = {
      id: crypto.randomUUID(),
      fullName: payload.fullName,
      phone: payload.phone,
      email: payload.email,
      status: payload.status,
      groupId: gid,
      groupName: null,
      photoUrl: payload.photoUrl?.trim() || null,
      createdAt: nowIso(),
      updatedAt: nowIso(),
    };
    rows.unshift(s);
    writeAll(rows);
    return expandStudent(s);
  },

  async update(id: string, payload: StudentUpdatePayload): Promise<Student> {
    await delay();
    const rows = readAll();
    const i = rows.findIndex((s) => s.id === id);
    if (i === -1) throw new Error("Student not found");
    const cur = rows[i]!;
    let groupId = cur.groupId;
    if (payload.groupId !== undefined) {
      groupId = payload.groupId?.trim() || null;
    }
    const next: Student = {
      ...cur,
      fullName: payload.fullName ?? cur.fullName,
      phone: payload.phone ?? cur.phone,
      email: payload.email ?? cur.email,
      status: (payload.status ?? cur.status) as StudentStatus,
      groupId,
      groupName: null,
      photoUrl:
        payload.photoUrl !== undefined
          ? payload.photoUrl?.trim() || null
          : cur.photoUrl,
      updatedAt: nowIso(),
    };
    rows[i] = next;
    writeAll(rows);
    return expandStudent(next);
  },

  async unassignAllByGroup(groupId: string): Promise<void> {
    await delay(200);
    const rows = readAll().map((s) =>
      s.groupId === groupId
        ? { ...s, groupId: null, groupName: null, updatedAt: nowIso() }
        : s,
    );
    writeAll(rows);
  },

  async remove(id: string): Promise<void> {
    await delay();
    const rows = readAll().filter((s) => s.id !== id);
    writeAll(rows);
  },
};
