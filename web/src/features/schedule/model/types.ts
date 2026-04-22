import type { Weekday } from "@/features/schedule/lib/weekdays";

export interface ScheduleLesson {
  id: string;
  weekday: Weekday;
  /** 24h `HH:mm`. */
  startTime: string;
  /** 24h `HH:mm`. */
  endTime: string;
  roomId: string;
  teacherId: string;
  groupId: string;
  /** Optional label shown on the card. */
  title: string | null;
  createdAt: string;
  updatedAt: string;
  room: { id: string; name: string };
  teacher: { id: string; fullName: string };
  group: { id: string; name: string };
}

export type LessonCreatePayload = {
  weekday: Weekday;
  startTime: string;
  endTime: string;
  roomId: string;
  teacherId: string;
  groupId: string;
  title?: string | null;
};

export type LessonUpdatePayload = Partial<LessonCreatePayload>;
