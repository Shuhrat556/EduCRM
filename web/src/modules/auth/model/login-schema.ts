import { z } from "zod";

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function isValidPhone(val: string) {
  const compact = val.replace(/[\s\-()]/g, "");
  const digits = compact.startsWith("+") ? compact.slice(1) : compact;
  return /^[0-9]{8,15}$/.test(digits);
}

export const loginFormSchema = z.object({
  identifier: z
    .string()
    .trim()
    .min(1, "Enter your email or phone number")
    .refine(
      (val) => EMAIL_RE.test(val) || isValidPhone(val),
      "Enter a valid email or phone number",
    ),
  password: z
    .string()
    .min(1, "Enter your password")
    .min(8, "Password must be at least 8 characters"),
  rememberMe: z.boolean(),
});

export type LoginFormValues = z.infer<typeof loginFormSchema>;
