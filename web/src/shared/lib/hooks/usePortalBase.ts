import { useLocation } from "react-router-dom";

const PORTAL_SEGMENTS = new Set([
  "student",
  "admin",
  "teacher",
  "super-admin",
]);

/**
 * Returns `/${segment}` for the current portal (first path segment), e.g. `/admin`.
 * Use for building module links inside feature screens.
 */
export function usePortalBase(): string {
  const { pathname } = useLocation();
  const seg = pathname.split("/").filter(Boolean)[0] ?? "student";
  if (PORTAL_SEGMENTS.has(seg)) return `/${seg}`;
  return "/student";
}
