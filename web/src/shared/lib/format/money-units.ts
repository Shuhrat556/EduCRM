/** Whole-currency amounts (e.g. monthly fee fields), not minor units. */
export function formatMoneyUnits(amount: number, currencyCode: string): string {
  try {
    return new Intl.NumberFormat(undefined, {
      style: "currency",
      currency: currencyCode,
      maximumFractionDigits: 0,
    }).format(amount);
  } catch {
    return `${amount.toFixed(0)} ${currencyCode}`;
  }
}
