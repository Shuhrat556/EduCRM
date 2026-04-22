/** Whole currency units (e.g. monthly tuition), locale-formatted without decimals. */
export function formatGroupFee(amount: number) {
  try {
    return new Intl.NumberFormat(undefined, {
      maximumFractionDigits: 0,
    }).format(amount);
  } catch {
    return String(amount);
  }
}
