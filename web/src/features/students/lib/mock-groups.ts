import type { StudentGroupOption } from "@/features/students/model/types";

/** Static fallback when `VITE_GROUPS_DEMO` is off (student assignment pickers). */
export const MOCK_GROUPS: StudentGroupOption[] = [
  { id: "g1", name: "Grade 10-A" },
  { id: "g2", name: "Grade 10-B" },
  { id: "g3", name: "Grade 11 Sciences" },
  { id: "g4", name: "Grade 9 — English track" },
];
