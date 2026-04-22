/**
 * EduCRM API JSON envelope (`docs`/Swagger).
 * Success: `{ success: true, data: T }`
 * Error: `{ success: false, error: { code, message, kind } }`
 */

export type ApiErrorShape = {
  code: string;
  message: string;
  kind: string;
};

export type ApiSuccessEnvelope<T> = {
  success: true;
  data: T;
};

export type ApiFailureEnvelope = {
  success: false;
  error: ApiErrorShape;
};

const KIND_TO_HTTP: Record<string, number> = {
  validation: 400,
  unauthorized: 401,
  forbidden: 403,
  not_found: 404,
  conflict: 409,
  too_many_requests: 429,
};

export function httpStatusFromErrorKind(kind: string | undefined): number {
  if (!kind) return 0;
  return KIND_TO_HTTP[kind] ?? 0;
}

export function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}
