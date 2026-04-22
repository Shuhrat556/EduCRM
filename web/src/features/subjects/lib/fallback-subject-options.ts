import type { SubjectOption } from "@/features/subjects/model/types";

/** Static catalog when API is off but teacher demo still resolves subject names. */
export const FALLBACK_SUBJECT_OPTIONS: SubjectOption[] = [
  { id: "sub1", name: "Mathematics" },
  { id: "sub2", name: "English" },
  { id: "sub3", name: "Physics" },
  { id: "sub4", name: "History" },
  { id: "sub5", name: "Computer science" },
];
