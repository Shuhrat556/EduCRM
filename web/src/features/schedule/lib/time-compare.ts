/** Compare `HH:mm` / `H:mm` strings for chronological ordering. */
export function compareTimeStrings(a: string, b: string): number {
  const ma = /^(\d{1,2}):(\d{2})$/.exec(a.trim());
  const mb = /^(\d{1,2}):(\d{2})$/.exec(b.trim());
  if (!ma || !mb) return a.localeCompare(b);
  const ta = Number(ma[1]) * 60 + Number(ma[2]);
  const tb = Number(mb[1]) * 60 + Number(mb[2]);
  return ta - tb;
}
