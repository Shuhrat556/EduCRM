import { isPastDue } from "@/features/payments/lib/period";
import type { BillingStatus } from "@/features/payments/model/types";

export function computeBillingStatus(input: {
  tuitionExempt: boolean;
  expectedAmount: number;
  paidAmount: number;
  balance: number;
  dueDateYmd: string;
}): BillingStatus {
  if (input.tuitionExempt) return "waived";
  if (input.expectedAmount <= 0) return "no_fee";
  if (input.balance <= 0) return "paid";
  if (isPastDue(input.dueDateYmd)) return "overdue";
  if (input.paidAmount <= 0) return "unpaid";
  return "partial";
}
