import { apiClient } from "@/shared/api/client";
import { mockAttendanceStore } from "@/features/attendance/lib/mock-attendance-store";
import type {
  AttendanceBatchPayload,
  AttendanceSessionView,
  AttendanceStudentRow,
} from "@/features/attendance/model/types";
import { scheduleApi } from "@/features/schedule/api/schedule-api";
import { studentsApi } from "@/features/students/api/students-api";

function useAttendanceDemo() {
  return import.meta.env.VITE_ATTENDANCE_DEMO === "true";
}

const DEFAULT_STATUS = "present" as const;

function validateBatchRow(
  r: AttendanceBatchPayload["rows"][number],
): void {
  if (!["present", "absent", "late"].includes(r.status)) {
    throw new Error("Invalid attendance status");
  }
  if (
    r.lessonGrade !== null &&
    (r.lessonGrade < 0 ||
      r.lessonGrade > 100 ||
      !Number.isFinite(r.lessonGrade))
  ) {
    throw new Error("Grade must be between 0 and 100");
  }
  if (
    r.weeklyRating !== null &&
    (r.weeklyRating < 1 ||
      r.weeklyRating > 5 ||
      !Number.isInteger(r.weeklyRating))
  ) {
    throw new Error("Weekly rating must be 1–5");
  }
}

/**
 * REST (adjust to backend):
 * - GET `/attendance/sessions?lessonId=&sessionDate=` → `AttendanceSessionView`
 * - PUT `/attendance/sessions` body: `AttendanceBatchPayload`
 */
export const attendanceApi = {
  getSession: async (
    lessonId: string,
    sessionDate: string,
  ): Promise<AttendanceSessionView> => {
    if (!useAttendanceDemo()) {
      const { data } = await apiClient.get<AttendanceSessionView>(
        "/attendance/sessions",
        {
          params: { lessonId, sessionDate },
        },
      );
      return data;
    }

    const lesson = await scheduleApi.get(lessonId);
    const students = await studentsApi.list({
      page: 1,
      pageSize: 250,
      search: "",
      status: "all",
      groupId: lesson.groupId,
    });

    const stored = await mockAttendanceStore.getEntries(
      lessonId,
      sessionDate,
    );

    const rows: AttendanceStudentRow[] = students.items.map((s) => {
      const e = stored[s.id];
      return {
        studentId: s.id,
        studentName: s.fullName,
        status: e?.status ?? DEFAULT_STATUS,
        comment: e?.comment ?? null,
        lessonGrade: e?.lessonGrade ?? null,
        weeklyRating: e?.weeklyRating ?? null,
      };
    });

    return {
      lesson,
      sessionDate,
      rows,
    };
  },

  saveBatch: async (payload: AttendanceBatchPayload): Promise<void> => {
    for (const r of payload.rows) {
      validateBatchRow(r);
    }
    if (useAttendanceDemo()) {
      await mockAttendanceStore.saveBatch(payload);
      return;
    }
    await apiClient.put("/attendance/sessions", payload);
  },
};
