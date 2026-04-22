import { z } from "zod";

const envSchema = z.object({
  VITE_API_BASE_URL: z
    .string()
    .min(1, "VITE_API_BASE_URL is required")
    .refine(
      (v) => !v.endsWith("/"),
      "Remove trailing slash from VITE_API_BASE_URL",
    ),
  VITE_APP_NAME: z.string().optional().default("EduCRM Admin"),
  VITE_ENABLE_QUERY_DEVTOOLS: z
    .enum(["true", "false"])
    .optional()
    .transform((v) => v === "true"),
});

export type AppEnv = z.infer<typeof envSchema>;

function readRawEnv(): Record<string, string | undefined> {
  return {
    VITE_API_BASE_URL: import.meta.env.VITE_API_BASE_URL,
    VITE_APP_NAME: import.meta.env.VITE_APP_NAME,
    VITE_ENABLE_QUERY_DEVTOOLS: import.meta.env.VITE_ENABLE_QUERY_DEVTOOLS,
  };
}

let cached: AppEnv | null = null;

export function getEnv(): AppEnv {
  if (cached) return cached;
  const parsed = envSchema.safeParse(readRawEnv());
  if (!parsed.success) {
    const msg = parsed.error.errors.map((e) => e.message).join("; ");
    throw new Error(`Invalid environment: ${msg}`);
  }
  cached = parsed.data;
  return cached;
}

export function isQueryDevtoolsEnabled(): boolean {
  if (import.meta.env.DEV) return true;
  return Boolean(getEnv().VITE_ENABLE_QUERY_DEVTOOLS);
}
