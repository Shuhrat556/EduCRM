export function formatPercentRate(rate: number | null | undefined): string {
  if (rate === null || rate === undefined) return "—";
  return `${Math.round(rate * 1000) / 10}%`;
}
