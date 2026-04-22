import type { StudentStatus } from "@/features/students/model/types";

export const STUDENT_STATUS_OPTIONS: { value: StudentStatus; label: string }[] =
  [
    { value: "active", label: "Active" },
    { value: "inactive", label: "Inactive" },
    { value: "graduated", label: "Graduated" },
    { value: "suspended", label: "Suspended" },
  ];

export function studentStatusLabel(status: StudentStatus): string {
  return STUDENT_STATUS_OPTIONS.find((o) => o.value === status)?.label ?? status;
}
