import { z } from "zod";

import type { StudentStatus } from "@/features/students/model/types";

const statusEnum = z.enum([
  "active",
  "inactive",
  "graduated",
  "suspended",
]) satisfies z.ZodType<StudentStatus>;

export const studentFormSchema = z.object({
  fullName: z.string().trim().min(1, "Name is required").max(200),
  phone: z
    .string()
    .trim()
    .min(5, "Phone is too short")
    .max(32, "Phone is too long"),
  email: z.string().trim().email("Invalid email"),
  status: statusEnum,
  groupId: z.string().nullable().optional(),
  photoUrl: z
    .string()
    .optional()
    .refine(
      (v) => {
        const t = (v ?? "").trim();
        return !t || /^https?:\/\/.+/i.test(t);
      },
      "Enter a valid image URL",
    ),
});

export type StudentFormValues = z.infer<typeof studentFormSchema>;
