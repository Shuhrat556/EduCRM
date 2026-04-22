import type { TeacherStatus } from "@/features/teachers/model/types";

export const TEACHER_STATUS_OPTIONS: {
  value: TeacherStatus;
  label: string;
}[] = [
  { value: "active", label: "Active" },
  { value: "inactive", label: "Inactive" },
  { value: "on_leave", label: "On leave" },
];

export function teacherStatusLabel(status: TeacherStatus): string {
  return TEACHER_STATUS_OPTIONS.find((o) => o.value === status)?.label ?? status;
}
