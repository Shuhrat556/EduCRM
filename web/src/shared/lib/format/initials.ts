/** Two-letter initials for avatars and compact labels. */
export function initialsFromName(name: string, maxChars = 2): string {
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) return "?";
  if (parts.length === 1) {
    return parts[0]!.slice(0, maxChars).toUpperCase();
  }
  return parts
    .map((p) => p[0])
    .join("")
    .slice(0, maxChars)
    .toUpperCase();
}
