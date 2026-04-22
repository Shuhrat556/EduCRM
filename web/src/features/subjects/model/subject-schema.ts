import { z } from "zod";

import type { SubjectCreatePayload } from "@/features/subjects/model/types";

export const subjectFormSchema = z.object({
  name: z.string().trim().min(1, "Name is required").max(160),
  code: z.string().trim().max(32, "Code is too long"),
  description: z.string().trim().max(500, "Description is too long"),
});

export type SubjectFormValues = z.infer<typeof subjectFormSchema>;

export function subjectFormToPayload(
  values: SubjectFormValues,
): SubjectCreatePayload {
  const code = values.code.trim();
  const description = values.description.trim();
  return {
    name: values.name,
    code: code === "" ? null : code,
    description: description === "" ? null : description,
  };
}
