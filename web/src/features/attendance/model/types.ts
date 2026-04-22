import type { ScheduleLesson } from "@/features/schedule/model/types";

export type AttendanceStatus = "present" | "absent" | "late";

export interface AttendanceStudentRow {
  studentId: string;
  studentName: string;
  status: AttendanceStatus;
  comment: string | null;
  /** Lesson grade 0–100, optional. */
  lessonGrade: number | null;
  /** Optional weekly rating 1–5. */
  weeklyRating: number | null;
}

export interface AttendanceSessionView {
  lesson: ScheduleLesson;
  sessionDate: string;
  rows: AttendanceStudentRow[];
}

export type AttendanceBatchRowPayload = {
  studentId: string;
  status: AttendanceStatus;
  comment: string | null;
  lessonGrade: number | null;
  weeklyRating: number | null;
};

export type AttendanceBatchPayload = {
  lessonId: string;
  sessionDate: string;
  rows: AttendanceBatchRowPayload[];
};
