import type { AttendanceBatchPayload } from "@/features/attendance/model/types";
import type { AttendanceStatus } from "@/features/attendance/model/types";

const STORAGE_KEY = "educrm_mock_attendance_v1";

export type AttendanceEntry = {
  status: AttendanceStatus;
  comment: string | null;
  lessonGrade: number | null;
  weeklyRating: number | null;
  updatedAt: string;
};

type StoreShape = {
  sessions: Record<string, Record<string, AttendanceEntry>>;
};

function nowIso() {
  return new Date().toISOString();
}

function sessionKey(lessonId: string, sessionDate: string) {
  return `${lessonId}::${sessionDate}`;
}

function readStore(): StoreShape {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return { sessions: {} };
    const parsed = JSON.parse(raw) as StoreShape;
    return parsed?.sessions ? parsed : { sessions: {} };
  } catch {
    return { sessions: {} };
  }
}

function writeStore(store: StoreShape) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(store));
}

function delay(ms = 200) {
  return new Promise((r) => setTimeout(r, ms));
}

/** Synchronous read for analytics (e.g. Reports) without network delay. */
export function getAttendanceSessionsSnapshotSync(): Record<
  string,
  Record<string, AttendanceEntry>
> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return {};
    const parsed = JSON.parse(raw) as StoreShape;
    return parsed?.sessions ?? {};
  } catch {
    return {};
  }
}

export const mockAttendanceStore = {
  async getEntries(
    lessonId: string,
    sessionDate: string,
  ): Promise<Record<string, AttendanceEntry>> {
    await delay(120);
    const store = readStore();
    const k = sessionKey(lessonId, sessionDate);
    return { ...(store.sessions[k] ?? {}) };
  },

  async saveBatch(payload: AttendanceBatchPayload): Promise<void> {
    await delay(280);
    const store = readStore();
    const k = sessionKey(payload.lessonId, payload.sessionDate);
    const cur = { ...(store.sessions[k] ?? {}) };
    const t = nowIso();
    for (const r of payload.rows) {
      cur[r.studentId] = {
        status: r.status,
        comment: r.comment?.trim() || null,
        lessonGrade: r.lessonGrade,
        weeklyRating: r.weeklyRating,
        updatedAt: t,
      };
    }
    store.sessions[k] = cur;
    writeStore(store);
  },
};
