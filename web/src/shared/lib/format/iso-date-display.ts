/** Format `YYYY-MM-DD` (and similar) for locale display without timezone shift. */
export function formatIsoDateDisplay(ymd: string): string {
  if (!ymd) return "—";
  const x = new Date(`${ymd}T12:00:00`);
  return Number.isNaN(x.getTime()) ? ymd : x.toLocaleDateString();
}
