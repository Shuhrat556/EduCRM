import type { GroupStatus } from "@/features/groups/model/types";

export const GROUP_STATUS_OPTIONS: {
  value: GroupStatus;
  label: string;
}[] = [
  { value: "draft", label: "Draft" },
  { value: "active", label: "Active" },
  { value: "paused", label: "Paused" },
  { value: "completed", label: "Completed" },
  { value: "archived", label: "Archived" },
];

export function groupStatusLabel(s: GroupStatus): string {
  return GROUP_STATUS_OPTIONS.find((o) => o.value === s)?.label ?? s;
}
