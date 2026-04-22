import type { LoginCredentials } from "@/modules/auth/model/types";

const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

function normalizePhone(raw: string) {
  return raw.replace(/[\s\-()]/g, "");
}

/**
 * EduCRM `POST /api/v1/auth/login` body: `{ "login": "email|phone", "password" }`.
 * Email is normalized to lowercase; phone strips spaces/parentheses.
 */
export function buildLoginRequestBody(credentials: LoginCredentials) {
  const id = credentials.identifier.trim();
  const login = EMAIL_RE.test(id) ? id.toLowerCase() : normalizePhone(id);
  return { login, password: credentials.password };
}
