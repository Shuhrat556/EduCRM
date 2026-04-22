import type { PaymentType } from "@/features/payments/model/types";

export const PAYMENT_TYPE_OPTIONS: { value: PaymentType; label: string }[] = [
  { value: "cash", label: "Cash" },
  { value: "card", label: "Card" },
  { value: "bank_transfer", label: "Bank transfer" },
  { value: "other", label: "Other" },
];

export function paymentTypeLabel(t: PaymentType): string {
  return PAYMENT_TYPE_OPTIONS.find((o) => o.value === t)?.label ?? t;
}
