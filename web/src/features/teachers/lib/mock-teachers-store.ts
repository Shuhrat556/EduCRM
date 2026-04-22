import { mockGroupsStore } from "@/features/groups/lib/mock-groups-store";
import { MOCK_GROUPS } from "@/features/students/lib/mock-groups";
import { FALLBACK_SUBJECT_OPTIONS } from "@/features/subjects/lib/fallback-subject-options";
import { mockSubjectsStore } from "@/features/subjects/lib/mock-subjects-store";
import type {
  Teacher,
  TeacherCreatePayload,
  TeacherStatus,
  TeacherUpdatePayload,
  TeachersListParams,
  TeachersListResponse,
} from "@/features/teachers/model/types";

const STORAGE_KEY = "educrm_mock_teachers_v1";

type Row = {
  id: string;
  fullName: string;
  phone: string;
  email: string;
  status: TeacherStatus;
  groupIds: string[];
  subjectIds: string[];
  photoUrl: string | null;
  createdAt: string;
  updatedAt: string;
};

function nowIso() {
  return new Date().toISOString();
}

function subjectCatalog() {
  if (import.meta.env.VITE_SUBJECTS_DEMO === "true") {
    return mockSubjectsStore.getOptionsSync();
  }
  return FALLBACK_SUBJECT_OPTIONS;
}

function teacherGroupCatalog() {
  if (import.meta.env.VITE_GROUPS_DEMO === "true") {
    return mockGroupsStore.getOptionsSync();
  }
  return MOCK_GROUPS;
}

function expand(r: Row): Teacher {
  const catalog = subjectCatalog();
  const groupList = teacherGroupCatalog();
  const groups = r.groupIds
    .map((id) => groupList.find((g) => g.id === id))
    .filter(Boolean)
    .map((g) => ({ id: g!.id, name: g!.name }));
  const subjects = r.subjectIds
    .map((id) => catalog.find((s) => s.id === id))
    .filter(Boolean)
    .map((s) => ({ id: s!.id, name: s!.name }));
  return {
    ...r,
    groupIds: [...r.groupIds],
    subjectIds: [...r.subjectIds],
    groups,
    subjects,
  };
}

function seed(): Row[] {
  const t = nowIso();
  return [
    {
      id: "t_seed_1",
      fullName: "Samira Karimova",
      phone: "+998901111001",
      email: "samira.k@school.edu",
      status: "active",
      groupIds: ["g1", "g2"],
      subjectIds: ["sub1", "sub2"],
      photoUrl: null,
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "t_seed_2",
      fullName: "James Okonkwo",
      phone: "+998902222002",
      email: "james.o@school.edu",
      status: "active",
      groupIds: ["g3"],
      subjectIds: ["sub3", "sub5"],
      photoUrl: null,
      createdAt: t,
      updatedAt: t,
    },
    {
      id: "t_seed_3",
      fullName: "Elena Petrov",
      phone: "+998903333003",
      email: "elena.p@school.edu",
      status: "on_leave",
      groupIds: [],
      subjectIds: ["sub4"],
      photoUrl: null,
      createdAt: t,
      updatedAt: t,
    },
  ];
}

function readAll(): Row[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seed();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as Row[];
    return Array.isArray(parsed) ? parsed : seed();
  } catch {
    return seed();
  }
}

function writeAll(rows: Row[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function delay(ms = 260) {
  return new Promise((r) => setTimeout(r, ms));
}

function readFileDataUrl(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const r = new FileReader();
    r.onload = () => resolve(String(r.result));
    r.onerror = () => reject(new Error("read failed"));
    r.readAsDataURL(file);
  });
}

export const mockTeachersStore = {
  async list(params: TeachersListParams): Promise<TeachersListResponse> {
    await delay();
    let rows = readAll();
    const q = params.search.trim().toLowerCase();
    if (q) {
      rows = rows.filter(
        (s) =>
          s.fullName.toLowerCase().includes(q) ||
          s.phone.replace(/\s/g, "").includes(q.replace(/\s/g, "")) ||
          s.email.toLowerCase().includes(q),
      );
    }
    if (params.status !== "all") {
      rows = rows.filter((s) => s.status === params.status);
    }
    const total = rows.length;
    const start = (params.page - 1) * params.pageSize;
    const items = rows.slice(start, start + params.pageSize).map(expand);
    return {
      items,
      total,
      page: params.page,
      pageSize: params.pageSize,
    };
  },

  async getById(id: string): Promise<Teacher | null> {
    await delay(160);
    const r = readAll().find((s) => s.id === id);
    return r ? expand(r) : null;
  },

  async create(payload: TeacherCreatePayload): Promise<Teacher> {
    await delay();
    const rows = readAll();
    const row: Row = {
      id: crypto.randomUUID(),
      fullName: payload.fullName,
      phone: payload.phone,
      email: payload.email,
      status: payload.status,
      groupIds: [...(payload.groupIds ?? [])],
      subjectIds: [...(payload.subjectIds ?? [])],
      photoUrl: payload.photoUrl?.trim() || null,
      createdAt: nowIso(),
      updatedAt: nowIso(),
    };
    rows.unshift(row);
    writeAll(rows);
    return expand(row);
  },

  async update(id: string, payload: TeacherUpdatePayload): Promise<Teacher> {
    await delay();
    const rows = readAll();
    const i = rows.findIndex((s) => s.id === id);
    if (i === -1) throw new Error("Teacher not found");
    const cur = rows[i]!;
    const next: Row = {
      ...cur,
      fullName: payload.fullName ?? cur.fullName,
      phone: payload.phone ?? cur.phone,
      email: payload.email ?? cur.email,
      status: (payload.status ?? cur.status) as TeacherStatus,
      groupIds:
        payload.groupIds !== undefined
          ? [...payload.groupIds]
          : [...cur.groupIds],
      subjectIds:
        payload.subjectIds !== undefined
          ? [...payload.subjectIds]
          : [...cur.subjectIds],
      photoUrl:
        payload.photoUrl !== undefined
          ? payload.photoUrl
          : cur.photoUrl,
      updatedAt: nowIso(),
    };
    rows[i] = next;
    writeAll(rows);
    return expand(next);
  },

  async remove(id: string): Promise<void> {
    await delay();
    writeAll(readAll().filter((s) => s.id !== id));
  },

  async uploadPhoto(id: string, file: File): Promise<{ photoUrl: string }> {
    await delay(200);
    const dataUrl = await readFileDataUrl(file);
    const rows = readAll();
    const i = rows.findIndex((s) => s.id === id);
    if (i === -1) throw new Error("Teacher not found");
    rows[i]!.photoUrl = dataUrl;
    rows[i]!.updatedAt = nowIso();
    writeAll(rows);
    return { photoUrl: dataUrl };
  },

  getNameSync(id: string): string | undefined {
    return readAll().find((r) => r.id === id)?.fullName;
  },
};
