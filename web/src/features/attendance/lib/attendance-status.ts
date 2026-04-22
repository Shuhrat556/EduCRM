import type { AttendanceStatus } from "@/features/attendance/model/types";

export const ATTENDANCE_STATUS_OPTIONS: {
  value: AttendanceStatus;
  label: string;
}[] = [
  { value: "present", label: "Present" },
  { value: "absent", label: "Absent" },
  { value: "late", label: "Late" },
];

export function attendanceStatusLabel(s: AttendanceStatus): string {
  return ATTENDANCE_STATUS_OPTIONS.find((o) => o.value === s)?.label ?? s;
}
