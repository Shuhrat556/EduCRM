import { FALLBACK_SUBJECT_OPTIONS } from "@/features/subjects/lib/fallback-subject-options";
import type {
  Subject,
  SubjectCreatePayload,
  SubjectOption,
  SubjectUpdatePayload,
  SubjectsListParams,
  SubjectsListResponse,
} from "@/features/subjects/model/types";

const STORAGE_KEY = "educrm_mock_subjects_v1";

function nowIso() {
  return new Date().toISOString();
}

function seedRows(): Subject[] {
  const t = nowIso();
  return FALLBACK_SUBJECT_OPTIONS.map((o) => ({
    id: o.id,
    name: o.name,
    code: null,
    description: null,
    createdAt: t,
    updatedAt: t,
  }));
}

function readAll(): Subject[] {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      const initial = seedRows();
      localStorage.setItem(STORAGE_KEY, JSON.stringify(initial));
      return initial;
    }
    const parsed = JSON.parse(raw) as Subject[];
    return Array.isArray(parsed) ? parsed : seedRows();
  } catch {
    return seedRows();
  }
}

function writeAll(rows: Subject[]) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(rows));
}

function delay(ms = 240) {
  return new Promise((r) => setTimeout(r, ms));
}

export const mockSubjectsStore = {
  getNameSync(id: string): string | undefined {
    return readAll().find((s) => s.id === id)?.name;
  },

  getOptionsSync(): SubjectOption[] {
    return readAll().map((s) => ({ id: s.id, name: s.name }));
  },

  async getOptions(): Promise<SubjectOption[]> {
    await delay(120);
    return readAll().map((s) => ({ id: s.id, name: s.name }));
  },

  async list(params: SubjectsListParams): Promise<SubjectsListResponse> {
    await delay();
    let rows = readAll();
    const q = params.search.trim().toLowerCase();
    if (q) {
      rows = rows.filter(
        (s) =>
          s.name.toLowerCase().includes(q) ||
          (s.code?.toLowerCase().includes(q) ?? false) ||
          (s.description?.toLowerCase().includes(q) ?? false),
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

  async getById(id: string): Promise<Subject | null> {
    await delay(150);
    return readAll().find((s) => s.id === id) ?? null;
  },

  async create(payload: SubjectCreatePayload): Promise<Subject> {
    await delay();
    const rows = readAll();
    const s: Subject = {
      id: crypto.randomUUID(),
      name: payload.name.trim(),
      code: payload.code?.trim() || null,
      description: payload.description?.trim() || null,
      createdAt: nowIso(),
      updatedAt: nowIso(),
    };
    rows.unshift(s);
    writeAll(rows);
    return s;
  },

  async update(id: string, payload: SubjectUpdatePayload): Promise<Subject> {
    await delay();
    const rows = readAll();
    const i = rows.findIndex((s) => s.id === id);
    if (i === -1) throw new Error("Subject not found");
    const cur = rows[i]!;
    const next: Subject = {
      ...cur,
      name: payload.name !== undefined ? payload.name.trim() : cur.name,
      code:
        payload.code !== undefined
          ? payload.code?.trim() || null
          : cur.code,
      description:
        payload.description !== undefined
          ? payload.description?.trim() || null
          : cur.description,
      updatedAt: nowIso(),
    };
    rows[i] = next;
    writeAll(rows);
    return next;
  },

  async remove(id: string): Promise<void> {
    await delay();
    writeAll(readAll().filter((s) => s.id !== id));
  },
};
