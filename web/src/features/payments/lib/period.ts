/** `YYYY-MM` for the current calendar month (UTC date portion). */
export function currentPeriodMonth(): string {
  return new Date().toISOString().slice(0, 7);
}

/** Last day of the billing month as `YYYY-MM-DD`. */
export function dueDateForPeriodMonth(periodMonth: string): string {
  const [yStr, mStr] = periodMonth.split("-");
  const y = Number(yStr);
  const m = Number(mStr);
  if (!y || !m || m < 1 || m > 12) return `${periodMonth}-28`;
  const last = new Date(Date.UTC(y, m, 0));
  return last.toISOString().slice(0, 10);
}

export function isPastDue(dueDateYmd: string): boolean {
  const due = new Date(dueDateYmd + "T23:59:59.999Z");
  return Date.now() > due.getTime();
}
