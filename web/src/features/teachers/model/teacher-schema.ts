import { z } from "zod";

import type { TeacherStatus } from "@/features/teachers/model/types";

const statusEnum = z.enum([
  "active",
  "inactive",
  "on_leave",
]) satisfies z.ZodType<TeacherStatus>;

export const teacherFormSchema = z.object({
  fullName: z.string().trim().min(1, "Name is required").max(200),
  phone: z
    .string()
    .trim()
    .min(5, "Phone is too short")
    .max(32, "Phone is too long"),
  email: z.string().trim().email("Invalid email"),
  status: statusEnum,
  groupIds: z.array(z.string()),
  subjectIds: z.array(z.string()),
});

export type TeacherFormValues = z.infer<typeof teacherFormSchema>;
