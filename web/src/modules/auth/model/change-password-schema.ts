import { z } from "zod";

const baseFields = {
  newPassword: z
    .string()
    .min(8, "Use at least 8 characters")
    .max(128, "Password is too long"),
  confirmPassword: z.string(),
};

function matchPasswords<T extends { newPassword: string; confirmPassword: string }>(
  data: T,
  ctx: z.RefinementCtx,
) {
  if (data.confirmPassword !== data.newPassword) {
    ctx.addIssue({
      code: z.ZodIssueCode.custom,
      message: "Passwords do not match",
      path: ["confirmPassword"],
    });
  }
}

export const changePasswordSchema = z
  .object({
    currentPassword: z.string().optional(),
    ...baseFields,
  })
  .superRefine(matchPasswords);

export const profilePasswordSchema = z
  .object({
    currentPassword: z.string().min(1, "Current password is required"),
    ...baseFields,
  })
  .superRefine(matchPasswords);

export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>;
