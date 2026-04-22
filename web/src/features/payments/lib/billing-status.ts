import type { BillingStatus } from "@/features/payments/model/types";

export const BILLING_STATUS_OPTIONS: {
  value: BillingStatus | "all";
  label: string;
}[] = [
  { value: "all", label: "All statuses" },
  { value: "paid", label: "Paid" },
  { value: "partial", label: "Partial" },
  { value: "unpaid", label: "Unpaid" },
  { value: "overdue", label: "Overdue" },
  { value: "waived", label: "Tuition free" },
  { value: "no_fee", label: "No group fee" },
];

export function billingStatusLabel(s: BillingStatus): string {
  const map: Record<BillingStatus, string> = {
    paid: "Paid",
    partial: "Partial",
    unpaid: "Unpaid",
    overdue: "Overdue",
    waived: "Tuition free",
    no_fee: "No fee",
  };
  return map[s] ?? s;
}

export function billingStatusBadgeVariant(
  s: BillingStatus,
): "default" | "secondary" | "outline" | "destructive" {
  switch (s) {
    case "paid":
      return "secondary";
    case "partial":
      return "outline";
    case "unpaid":
      return "outline";
    case "overdue":
      return "destructive";
    case "waived":
      return "default";
    case "no_fee":
      return "secondary";
    default:
      return "outline";
  }
}
